package parser

import (
	"regexp"

	"github.com/karasiq69/squad-rcon-go-v2-noping/rconEvents"
	"github.com/karasiq69/squad-rcon-go-v2-noping/rconTypes"
)

func showNextMap(line string) (event string, data interface{}) {
	re := regexp.MustCompile(`^Next level is (.*), layer is (.*)`)
	matches := re.FindStringSubmatch(line)

	if matches != nil {
		return rconEvents.SHOW_NEXT_MAP, rconTypes.NextMap{
			Raw:   line,
			Level: matches[1],
			Layer: matches[2],
		}
	}

	return rconEvents.SHOW_NEXT_MAP, nil
}
