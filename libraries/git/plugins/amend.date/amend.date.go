package git

import (
	"context"
	"fmt"
	"time"

	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/runner"
)

// TODO: can this be implemented without os/exec?
// TODO: modify this to be able to target any git hash
func AmendDate(parentContext context.Context, dateString string) {
	ctx, cancel := context.WithTimeout(parentContext, 30*time.Second)
	defer cancel()

	_, err := runner.Run(ctx, "git", "commit", "--amend", "--no-edit",
		fmt.Sprintf("--date=\"%v\"", dateString))
	logging.Panic(err)
}
