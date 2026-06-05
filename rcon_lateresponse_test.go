package rcon

import (
	"bufio"
	"encoding/binary"
	"io"
	"net"
	"testing"
	"time"

	"github.com/karasiq69/squad-rcon-go-v2-noping/internal/utils"
)

// readFramed consumes one length-prefixed RCON packet (4-byte little-endian
// size, then that many bytes) from the client and discards it.
func readFramed(br *bufio.Reader) error {
	var sz [4]byte
	if _, err := io.ReadFull(br, sz[:]); err != nil {
		return err
	}
	body := make([]byte, binary.LittleEndian.Uint32(sz[:]))
	_, err := io.ReadFull(br, body)
	return err
}

// writeResponse emits a command response packet followed by the 7-byte
// end-of-response marker the client's byteReader keys on.
func writeResponse(conn net.Conn, body string) {
	conn.Write(utils.Encode(serverDataResponse, executeCommandID, body))
	conn.Write([]byte{0, 1, 0, 0, 0, 0, 0})
}

// TestExecute_LateResponseDoesNotWedgeReader reproduces the production zombie:
// a command (cmd1) whose response arrives AFTER Execute's timeout must neither
// (a) block byteReader forever on the unbuffered executeChan, nor (b) get
// mis-delivered to the next command (cmd2). After the fix the stale response is
// dropped and cmd2 receives its own response.
func TestExecute_LateResponseDoesNotWedgeReader(t *testing.T) {
	prev := ExecuteTimeout
	ExecuteTimeout = 200 * time.Millisecond
	defer func() { ExecuteTimeout = prev }()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	done := make(chan struct{})
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		br := bufio.NewReader(conn)
		_ = readFramed(br) // auth packet
		// cmd1: two packets (command + empty), then respond LATE.
		_ = readFramed(br)
		_ = readFramed(br)
		time.Sleep(400 * time.Millisecond) // > ExecuteTimeout: cmd1 times out
		writeResponse(conn, "LATE1")
		// cmd2: respond promptly.
		_ = readFramed(br)
		_ = readFramed(br)
		writeResponse(conn, "OK2")
		<-done       // hold the connection until the test is done asserting
		conn.Close() // clean server-side close lets byteReader exit on its own
	}()

	host, port, _ := net.SplitHostPort(ln.Addr().String())
	r, err := NewRcon(RconConfig{Host: host, Port: port, Password: "x", AutoReconnect: false})
	if err != nil {
		t.Fatal(err)
	}

	if got := r.Execute("cmd1"); got != "" {
		t.Fatalf("cmd1: expected empty (timeout), got %q", got)
	}
	// Let the late cmd1 response arrive while nobody is waiting.
	time.Sleep(300 * time.Millisecond)
	if got := r.Execute("cmd2"); got != "OK2" {
		t.Fatalf("cmd2: expected \"OK2\", got %q (reader wedged or stale response mis-delivered)", got)
	}
	close(done)
}
