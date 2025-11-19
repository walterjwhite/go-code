package timeofday

import (
	"testing"
	"time"
)

func TestTill(t *testing.T) {
	now := time.Now()
	futureTime := now.Add(1 * time.Hour)
	pastTime := now.Add(-1 * time.Hour)

	todFuture := &TimeOfDay{Hour: futureTime.Hour(), Minute: futureTime.Minute()}
	todPast := &TimeOfDay{Hour: pastTime.Hour(), Minute: pastTime.Minute()}

	if todFuture.Till() < 0 {
		t.Errorf("Till() for a future time should be positive")
	}

	if todPast.Till() > 0 {
		t.Errorf("Till() for a past time should be negative")
	}
}





