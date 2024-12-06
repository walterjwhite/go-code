package spot

import (
	"context"

	"testing"
	"time"
)

func TestMonitor(t *testing.T) {
	ctx := context.Background()

	*testDataPath = "./test.data.secret"


	c := New("export_test")
	c.Monitor(ctx)

	time.Sleep(1 * time.Minute)
}
