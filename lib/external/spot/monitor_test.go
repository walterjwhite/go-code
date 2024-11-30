package spot

import (
	"context"

	"testing"
	"time"
)

func TestMonitor(t *testing.T) {
	ctx := context.Background()

	*testDataPath = "./test.data.secret"

	// constant now, unable to tweak
	//minRefreshInterval = time.Duration(1 * time.Second)

	c := New("export_test")
	c.Monitor(ctx)

	// allow feed to update
	time.Sleep(1 * time.Minute)
}
