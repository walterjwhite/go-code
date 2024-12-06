package run

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/notification"
)

func (a *Application) monitorChannel(ctx context.Context, channel chan *string) {
	select {
	case applicationStartedLine := <-channel:
		log.Info().Msgf("Application Started: %v\n", *applicationStartedLine)
		notification.NotifierInstance.Notify(notification.Notification{Title: fmt.Sprintf("run: %v", a.Name), Description: *applicationStartedLine, Type: notification.Info})
	case <-ctx.Done():
		close(channel)
		return
	}
}
