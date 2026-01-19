package property

import (
	"github.com/rs/zerolog/log"
	"regexp"
	"strings"

	"github.com/vrischmann/envconfig"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/utils/typename"
)

func LoadEnv(config interface{}) {
	prefix := getShortPrefix(config)
	log.Debug().Msgf("Loading environment variables with prefix: %s", prefix)
	logging.Warn(envconfig.InitWithOptions(config, envconfig.Options{Prefix: prefix, AllOptional: true}), "LoadEnv - InitWithTypePrefix failed")
}

func getShortPrefix(config interface{}) string {
	return sanitizeEnvKey(typename.Get(config))
}

func sanitizeEnvKey(key string) string {
	reg := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	sanitized := reg.ReplaceAllString(key, "_")
	return strings.ToUpper(sanitized)
}
