package window

import (
	"context"
	"time"
)

type TimeOfDay struct {
	TargetHour   int
	TargetMinute int
}

func New(parent_ctx context.Context, t *TimeOfDay) (context.Context, context.CancelFunc) {
	now := time.Now()

	targetTime := time.Date(now.Year(), now.Month(), now.Day(), t.TargetHour, t.TargetMinute, 0, 0, now.Location())
	if targetTime.Before(now) {
		targetTime = targetTime.Add(24 * time.Hour)
	}

	durationUntilTarget := targetTime.Sub(now)

	return context.WithTimeout(context.Background(), durationUntilTarget)
}


