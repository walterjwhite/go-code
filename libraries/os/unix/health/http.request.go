package health

import (
	"net/http"
	//"../data"
	"os"
	"time"
)

const requestTimeout = 3 * time.Second
const goodResponseTime = 1 * time.Second

func Http(proxy string, target string) int {
	setProxy(proxy)

	client := http.Client{Timeout: requestTimeout}

	start := time.Now()
	resp, err := client.Get(target)
	end := time.Now()
	elapsed := end.Sub(start)

	if err != nil {
		return HEALTH_BAD
	}

	return getStatusCode(resp.StatusCode, elapsed)
}

func getStatusCode(statusCode int, elapsed time.Duration) int {
	if statusCode < 400 {
		if elapsed <= goodResponseTime {
			return HEALTH_GOOD
		}

		return HEALTH_ERRORS
	}

	return HEALTH_BAD
}

// as we're using go-routines, this will muck with TOR/straight-out
func setProxy(proxy string) {
	os.Setenv("http_proxy", proxy)
	os.Setenv("https_proxy", proxy)
	os.Setenv("HTTP_PROXY", proxy)
	os.Setenv("HTTPS_PROXY", proxy)
}
