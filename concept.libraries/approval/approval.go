package approval

import (
	"time"
)

type Request struct {
	Id      string
	Details string

	Deadline time.Time
}

type ApprovalAction int

const (
	Approve ApprovalAction = iota
	Deny
)

type Action struct {
	Request *Request

	Action   ApprovalAction
	Comments string

	Username string

	Other []string
}

type ApprovalWriter interface {
	Write(action *Action)
}
