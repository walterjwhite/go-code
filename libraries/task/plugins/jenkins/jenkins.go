package jenkins

import (
	"time"
)

type Jenkins struct {
	Url     string
	Job     string
	Timeout time.Duration
}
