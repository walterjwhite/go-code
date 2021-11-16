package movement

import (
	"fmt"

	"github.com/walterjwhite/go/lib/external/spot/notification"
)

func (c *MovementConfiguration) onTimeout() error {
	c.schedule(c.getDuration())

	n := c.buildNotification()
	n.Send()

	return nil
}

func (c *MovementConfiguration) getNotificationName() string {
	return fmt.Sprintf("movement-%s", c.getAlertLevel())
}

func (c *MovementConfiguration) getAlertLevel() string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if c.count < c.AlertAfter {
		return "warning"
	}

	return "alert"
}

func (c *MovementConfiguration) buildNotification() *notification.Notification {
	n := notification.New(c.Session,
		c.Session.LatestReceivedRecord, c.getNotificationName())

	n.Context["StartHour"] = fmt.Sprint(c.StartHour)
	n.Context["StartMinute"] = fmt.Sprint(c.StartMinute)
	n.Context["EndHour"] = fmt.Sprint(c.EndHour)
	n.Context["EndMinute"] = fmt.Sprint(c.EndMinute)

	n.Context["AlertAfter"] = fmt.Sprint(c.AlertAfter)
	n.Context["MovementDurationTimeout"] = fmt.Sprintf("%.0f minutes", c.MovementDurationTimeout.Minutes())
	n.Context["SuspendDurationTimeout"] = fmt.Sprintf("%.0f minutes", c.SuspendDurationTimeout.Minutes())

	return n
}
