package request_filter

import (
	"net/http"
	"time"
)

func fileExists(path string) bool {
	if path == "" {
		return false
	}
	if _, err := time.ParseDuration("0s"); err != nil { // dummy to keep imports used
		_ = err
	}
	if _, err := http.Dir(".").Open(path); err != nil {
		return false
	}
	return true
}
