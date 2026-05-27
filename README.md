# squad-rcon-go-v2-noping

Fork of [SquadGO/squad-rcon-go](https://github.com/SquadGO/squad-rcon-go) v2.0.5 with the internal `ping()` goroutine removed.

Upstream spawns a goroutine in `connect()` that calls `Execute("PING_CONNECTION")` every 10s, bypassing any external serialization. That races with consumer `Execute` calls and corrupts the FIFO response correlation, blowing up the TCP connection. Auto-reconnect masks the symptom but the connection ends up flapping.

This fork removes that goroutine. Keeping the socket warm is now the consumer's responsibility — issue `PING_CONNECTION` via your own serialized command path.

Everything else is identical to upstream v2.0.5.

## Install

```text
go get -u github.com/karasiq69/squad-rcon-go-v2-noping
```

## Quick start example

```golang
import (
  "fmt"
  rcon "github.com/karasiq69/squad-rcon-go-v2-noping"
  "github.com/karasiq69/squad-rcon-go-v2-noping/rconEvents"
  "github.com/karasiq69/squad-rcon-go-v2-noping/rconTypes"
)

func main() {
  r, err := rcon.NewRcon(rcon.RconConfig{Host: "127.0.0.1", Password: "123456", Port: "27165", AutoReconnect: true, AutoReconnectDelay: 5})
  if err != nil {
    fmt.Println(err)
    return
  }

  defer r.Close()

  fmt.Println("[RCON] Connection successful")

  /* Listeners works after first initialization */

  r.Emitter.On(rconEvents.CONNECTED, func(_ interface{}) {
    fmt.Println("[RCON] Connection successful")
  })

  r.Emitter.On(rconEvents.RECONNECTING, func(_ interface{}) {
    fmt.Println("[RCON] Reconnecting")
  })

  r.Emitter.On(rconEvents.CLOSE, func(_ interface{}) {
    fmt.Println("[RCON] Connection closed")
  })

  r.Emitter.On(rconEvents.ERROR, func(err interface{}) {
    fmt.Println(err)
  })

  r.Emitter.On(rconEvents.DATA, func(data interface{}) {
    fmt.Println("Data: ", data)
  })

  r.Emitter.On(rconEvents.CHAT_MESSAGE, func(data interface{}) {
    if v, ok := data.(rconTypes.Message); ok {
      fmt.Println("Message: ", v.Message)
    }
  })

  r.Emitter.On(rconEvents.LIST_PLAYERS, func(data interface{}) {
    if v, ok := data.(rconTypes.Players); ok {
      fmt.Println("Players: ", v)
    }
  })

  r.Execute(rconEvents.LIST_PLAYERS)

  // Use to prevent the program from ending
  select {}
}
```

## Listeners

| Listener                     | Returns           |
| ---------------------------- | ----------------- |
| **CONNECTED**                | **nil**           |
| **RECONNECTING**             | **nil**           |
| **CLOSE**                    | **nil**           |
| **ERROR**                    | **Error**         |
| **DATA**                     | **String**        |
| **CHAT_MESSAGE**             | **Message**       |
| **SQUAD_CREATED**            | **SquadCreated**  |
| **PLAYER_WARNED**            | **Warn**          |
| **PLAYER_KICKED**            | **Kick**          |
| **PLAYER_BANNED**            | **Ban**           |
| **POSSESSED_ADMIN_CAMERA**   | **PosAdminCam**   |
| **UNPOSSESSED_ADMIN_CAMERA** | **UnposAdminCam** |
| **LIST_PLAYERS**             | **Players**       |
| **LIST_SQUADS**              | **Squads**        |
| **SHOW_SERVER_INFO**         | **ServerInfo**    |
| **SHOW_CURRENT_MAP**         | **CurrentMap**    |
| **SHOW_NEXT_MAP**            | **NextMap**       |
