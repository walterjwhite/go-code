package data

import (
	"github.com/walterjwhite/go/lib/external/spot/client"
	"time"
)

type Record struct {
	Id           int
	Latitude     float64
	Longitude    float64
	Altitude     float64
	DateTime     time.Time
	BatteryState string
	Message      string
	MessageType  client.MessageType
}

func New(message *client.Message) *Record {
	return &Record{Id: message.Id,
		Latitude:     message.Latitude,
		Longitude:    message.Longitude,
		Altitude:     message.Altitude,
		DateTime:     time.Unix(message.UnixTime, 0),
		Message:      message.MessageContent,
		MessageType:  message.MessageType,
		BatteryState: message.BatteryState,
	}
}
