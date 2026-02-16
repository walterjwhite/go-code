package property

import (
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/vrischmann/envconfig"
	"github.com/walterjwhite/go-code/lib/utils/typename"
)

func LoadEnv(config interface{}) {
	prefix := getShortPrefix(config)
	log.Debug().Msgf("Loading environment variables with prefix: %s", prefix)

	err := envconfig.InitWithOptions(config, envconfig.Options{Prefix: prefix, AllOptional: true})
	if err != nil {
		if !strings.Contains(err.Error(), "unexported field") {
			log.Debug().Msgf("LoadEnv - InitWithOptions failed: %v", err)
		}
	}
}

func getShortPrefix(config interface{}) string {
	return sanitizeEnvKey(typename.Get(config))
}

func sanitizeEnvKey(key string) string {
	reg := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	sanitized := reg.ReplaceAllString(key, "_")
	return strings.ToUpper(sanitized)
}
