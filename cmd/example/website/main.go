package main

import (
	"context"

	"flag"

	"time"

	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/net/email"
	"github.com/walterjwhite/go-code/lib/net/email/flusher"

	"github.com/robfig/cron/v3"
	"github.com/walterjwhite/go-code/lib/io/pipe"
)

var (
	portFlag            = flag.Int("p", 8080, "Port to listen on, 8080 by default")
	hostFlag            = flag.String("h", "localhost", "Host to listen on, localhost by default")
	dailyExportCronFlag = flag.String("e", "@daily", "cron expression for daily export schedule, '@daily' by default")

	emailAccount           = &email.EmailAccount{}
	requestLogEmailFlusher = &flusher.EmailFlusher{}
	pipeReader             = &pipe.Reader{}
)

func init() {
	application.Configure(emailAccount, requestLogEmailFlusher, pipeReader)

	log.Warn().Msgf("email account loaded: %v", emailAccount)

	requestLogEmailFlusher.Account = emailAccount
	pipeReader.Flusher = requestLogEmailFlusher

	log.Debug().Msg("initialized")
}

func main() {
	server := serve()
	defer shutdown(server)

	go pipeReader.Start()
	initRequestLog()

	go func() {
		log.Warn().Msgf("starting server on %s", server.Addr)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.Panic(err)
		}
	}()

	application.Wait()
}

func shutdown(server *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	logging.Panic(server.Shutdown(ctx))
}

func initRequestLog() {
	c := cron.New(cron.WithLocation(time.Local))
	_, err := c.AddFunc(*dailyExportCronFlag, func() {
		log.Info().Msg("daily export job starting")
		logging.Warn(pipeReader.Flush(), "daily export job - failed to flush pipe reader")
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to schedule daily export job")
		return
	}

	c.Start()
}
