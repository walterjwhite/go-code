package client

type Message struct {
	Id             int
	MessengerId    string
	MessengerName  string
	UnixTime       int64
	MessageType    MessageType
	Latitude       float64
	Longitude      float64
	ModelId        ModelId
	ShowCustomMsg  string
	DateTime       SpotTime
	Hidden         int
	MessageContent string
	BatteryState   string
	Altitude       float64
}
