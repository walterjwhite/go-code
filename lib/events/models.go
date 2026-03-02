package events

type Action struct {
	ActionID     int    `json:"action_id"`
	Message      string `json:"message"`
	SupportsArgs bool   `json:"supports_args"`
}

type Event struct {
	EventID          int      `json:"event_id"`
	Details          string   `json:"details"`
	SupportedActions []Action `json:"supported_actions"`
}

type Response struct {
	EventID  int      `json:"event_id"`
	ActionID int      `json:"action_id"`
	Args     []string `json:"args,omitempty"`
}
