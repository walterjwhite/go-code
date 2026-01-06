package request_filter

import (
	"net/http"
	"time"
)

func fileExists(path string) bool {
    if path == "" {
        return false
    }
    if _, err := os.Stat(path); err != nil {
        return false
    }
		
    return true
}
