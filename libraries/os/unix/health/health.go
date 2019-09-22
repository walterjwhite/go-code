package health

import (
	"context"
	"log"
)

const (
	HEALTH_UNCONFIGURED = -1
	HEALTH_INITIALIZED  = -2
	HEALTH_GOOD         = 0
	HEALTH_ERRORS       = 1
	HEALTH_BAD          = 2
)

func Health(ctx context.Context, check string, arguments []string) int {
	if len(check) == 0 {
		return HEALTH_UNCONFIGURED
	}

	switch {
	case "ping" == check:
		return Ping()

	case "chrony-sources" == check:
		return ChronySources(ctx)

	case "http" == check:
		return Http(arguments[0], arguments[1])

	case "shorewall" == check:
		return IsShorewallRunning()

	case "dig" == check:
		return Dig(arguments[0], arguments[1])
	}

	log.Printf("Unrecognized check: %v", check)
	return HEALTH_BAD
}
