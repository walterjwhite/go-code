package run

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/notification"
)

// TODO: the channel is only sending the matching line
// perhaps we should instead send a notification on the channel
// then we can configure how we want to receive the notifications here, OS notification, email, sms, etc.
func monitorChannel(ctx context.Context, application string, channel chan *string) {
	select {
	case applicationStartedLine := <-channel:
		log.Info().Msgf("Application Started: %v\n", applicationStartedLine)
		notification.NotifierInstance.Notify(notification.Notification{Title: fmt.Sprintf("run: %v", application), Description: *applicationStartedLine, Type: notification.Info})
	case <-ctx.Done():
		close(channel)
		return
	}
}
