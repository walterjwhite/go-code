package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/io/pipe"
	"github.com/walterjwhite/go-code/lib/net/email"
	"github.com/walterjwhite/go-code/lib/net/email/flusher"
	googlepubsub "github.com/walterjwhite/go-code/lib/net/google"
)

var (
	portFlag            = flag.Int("p", 8080, "Port to listen on")
	hostFlag            = flag.String("h", "localhost", "Host to listen on")
	dailyExportCronFlag = flag.String("e", "@daily", "cron expression for daily export")

	emailAccount           = &email.EmailAccount{}
	requestLogEmailFlusher = &flusher.EmailFlusher{}
	pipeReader             = &pipe.Reader{}

	chatPubSubConf = &ChatPubSubConfiguration{
		PubSubConf: &googlepubsub.Conf{},
	}
)

type ChatPubSubConfiguration struct {
	PublishTopic         string
	ResponseTopic        string
	ResponseSubscription string

	PubSubConf *googlepubsub.Conf
}

func (c *ChatPubSubConfiguration) PostLoad(ctx context.Context) error {
	if c == nil || c.PubSubConf == nil || c.PublishTopic == "" || c.ResponseTopic == "" || c.ResponseSubscription == "" {
		log.Info().Msgf("pubsub chat cfg: %s | %s | %s", c.PublishTopic, c.ResponseTopic, c.ResponseSubscription)
		log.Info().Msg("chat pubsub disabled — chatClient will remain nil")
		return nil
	}

	if err := c.validate(); err != nil {
		return err
	}

	if err := c.PubSubConf.Init(ctx); err != nil {
		return fmt.Errorf("init GCP pubsub client: %w", err)
	}

	chatClient = newPubsubChatService(
		c.PubSubConf,
		c.PublishTopic,
		c.ResponseTopic,
		c.ResponseSubscription,
	)

	log.Info().Msgf("pubsub chat: %s | %s | %s", c.PublishTopic, c.ResponseTopic, c.ResponseSubscription)

	return nil
}

func (c *ChatPubSubConfiguration) validate() error {
	if c.PublishTopic == "" {
		return fmt.Errorf("CHAT_PUBLISH_TOPIC is required")
	}
	if c.ResponseTopic == "" {
		return fmt.Errorf("CHAT_RESPONSE_TOPIC is required")
	}
	if c.ResponseSubscription == "" {
		return fmt.Errorf("CHAT_RESPONSE_SUBSCRIPTION is required")
	}
	if c.PubSubConf == nil {
		return fmt.Errorf("pubsub config is required")
	}
	return nil
}

func isRunningTests() bool {
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-test.") {
			return true
		}
	}
	return false
}

func init() {
	gin.SetMode(gin.ReleaseMode)

	if isRunningTests() {
		return
	}

	application.Configure(emailAccount, requestLogEmailFlusher, pipeReader, chatPubSubConf)

	if emailAccount != nil && emailAccount.EmailAddress != nil {
		log.Info().Str("address", emailAccount.EmailAddress.Address).Msg("email account loaded")
	} else {
		log.Warn().Msg("email account loaded but address is not configured")
	}

	requestLogEmailFlusher.Account = emailAccount
	pipeReader.Flusher = requestLogEmailFlusher

	log.Debug().Msg("initialized")
}

func main() {
	defer application.OnPanic()
	server := serve()
	defer shutdown(server)

	go pipeReader.Start()
	initRequestLog()

	go func() {
		log.Info().Str("addr", server.Addr).Msg("starting server")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.Error(err)
		}
	}()

	application.Wait()
}

func shutdown(server *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	logging.Error(server.Shutdown(ctx))
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
