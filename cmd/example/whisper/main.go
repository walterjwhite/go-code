package main

import (
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"regexp"
	"time"
)

var (
	transcriber = &Transcriber{}
)

func init() {
	application.Configure(transcriber)

	transcriber.MessageHandlers = []Conf{Patterns: }
	log.Debug().Msgf("conf: %v", transcriber)

	logging.Panic(transcriber.init())
}

func main() {
	defer transcriber.cleanup()

	log.Info().Msg("starting")
	logging.Panic(transcriber.Start())

	for {
		time.Sleep(transcriber.ChunkDuration)
		transcriber.Process()
		log.Info().Msg("processed")
	}
}

type MessageHandler interface {
	OnMessage(message string)
}

type Conf struct {
	Patterns []string
}

func (c *Conf) OnMessage(message string) {
	log.Info().Msgf("message: %s", message)

	matchingPattern := c.matches(message)
	if len(matchingPattern) > 0 {
		log.Info().Msgf("matched message: %s -> %s", message, matchingPattern)
	}
}

func (c *Conf) matches(message string) string {
	for _, pattern := range c.Patterns {
		re, err := regexp.Compile(pattern)
		logging.Panic(err)

		if re.MatchString(message) {
			return pattern
		}
	}

	return ""
}
