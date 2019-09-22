package health

import (
	"github.com/lixiangzhong/dnsutil"
	"github.com/walterjwhite/go-application/libraries/logging"
	"time"
)

//const localDnsServer = "127.0.0.1"
const (
	goodQueryTime = 500 * time.Millisecond
	slowQueryTime = 1 * time.Second
)

func Dig(server string, target string) int {
	var dig dnsutil.Dig
	logging.Panic(dig.SetDNS(server))

	dig.SetTimeOut(1 * time.Second)

	start := time.Now()
	_, err := dig.A(target)
	end := time.Now()
	elapsed := end.Sub(start)

	return digGetStatus(elapsed, err)
}

func digGetStatus(elapsedTime time.Duration, err error) int {
	if err != nil {
		return HEALTH_BAD
	}

	if elapsedTime <= goodQueryTime {
		return HEALTH_GOOD
	} else if elapsedTime <= slowQueryTime {
		return HEALTH_ERRORS
	}

	return HEALTH_BAD
}
