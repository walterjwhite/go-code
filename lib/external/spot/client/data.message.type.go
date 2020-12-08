package client

type MessageType string

const (
	OK              MessageType = "OK"
	TRACK           MessageType = "TRACK"
	EXTREME_TRACK   MessageType = "EXTREME-TRACK"
	UNLIMITED_TRACK MessageType = "UNLIMITED-TRACK"
	NEWMOVEMENT     MessageType = "NEWMOVEMENT"
	HELP            MessageType = "HELP"
	HELP_CANCEL     MessageType = "HELP-CANCEL"
	CUSTOM          MessageType = "CUSTOM"
	POI             MessageType = "POI"
	STOP            MessageType = "STOP"
)
