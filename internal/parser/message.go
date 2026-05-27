package parser

import (
	"regexp"
	"strings"

	"github.com/karasiq69/squad-rcon-go-v2-noping/rconEvents"
	"github.com/karasiq69/squad-rcon-go-v2-noping/rconTypes"
)

func message(line string) (event string, data interface{}) {
	re := regexp.MustCompile(`\[(ChatAll|ChatTeam|ChatSquad|ChatAdmin)] \[Online IDs:EOS: ([0-9a-f]{32}) steam: (\d{17})\] (.+?) : (.*)`)
	matches := re.FindStringSubmatch(line)

	if matches != nil {
		return rconEvents.CHAT_MESSAGE, rconTypes.Message{
			Raw:        line,
			ChatType:   matches[1],
			EosID:      matches[2],
			SteamID:    matches[3],
			PlayerName: strings.TrimSpace(matches[4]),
			Message:    matches[5],
		}
	}

	return rconEvents.CHAT_MESSAGE, nil
}
