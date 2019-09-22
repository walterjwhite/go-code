package health

import (
	"github.com/sparrc/go-ping"
)

const TARGET = "9.9.9.9"

func Ping() int {
	pinger, err := ping.NewPinger(TARGET)
	if err != nil {
		return HEALTH_BAD
	}

	pinger.Count = 3
	pinger.Run()
	stats := pinger.Statistics()

	if stats.PacketsRecv > 0 {
		if stats.PacketsRecv == stats.PacketsSent {
			return HEALTH_GOOD
		}

		return HEALTH_ERRORS
	}

	return HEALTH_BAD
}
