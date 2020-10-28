package idle

import (
	"github.com/walterjwhite/go/lib/application/logging"
	"os/exec"

	"bytes"
	"context"
	"strconv"
	"strings"
	"time"
)

func IdleTime(ctx context.Context) time.Duration {
	cmd := exec.CommandContext(ctx, "xprintidle")

	buffer := new(bytes.Buffer)

	runner.WithWriter(cmd, buffer)

	logging.Panic(cmd.Start())
	logging.Panic(cmd.Wait())

	idleTimeInMillis := strings.TrimSuffix(buffer.String(), "\n")

	i, err := strconv.ParseInt(idleTimeInMillis, 10, 64)
	logging.Panic(err)

	return time.Duration(i) * time.Millisecond
}
