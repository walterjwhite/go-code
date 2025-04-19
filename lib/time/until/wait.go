package until

import (
	"time"
)

func Wait(hour int) bool {
	currentTime := time.Now()
	if currentTime.Hour() < hour {
		targetTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), hour, 0, 0, 0, currentTime.Location())
		duration := targetTime.Sub(currentTime)

		time.Sleep(duration)
		return true
	}

	return false
}
