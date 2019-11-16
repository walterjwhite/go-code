package idle

import (
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/runner"

	"bytes"
	"context"
	"time"
	"strconv"
	"strings"
)


func IdleTime(ctx context.Context) time.Duration {
	cmd := runner.Prepare(ctx, "xprintidle")
	
	buffer := new(bytes.Buffer)

	runner.WithWriter(cmd, buffer)
	
	logging.Panic(runner.Start(cmd))
	logging.Panic(runner.Wait(cmd))
	
	idleTimeInMillis := strings.TrimSuffix(buffer.String(), "\n")

	i, err := strconv.ParseInt(idleTimeInMillis, 10, 64)
	logging.Panic(err)
	
	return time.Duration(i)*time.Millisecond
}
