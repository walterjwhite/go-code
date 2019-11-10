package git

import (
	"context"
	"fmt"
	"time"

	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/runner"
)

func AmendDate(parentContext context.Context, dateString string) {
	ctx, cancel := context.WithTimeout(parentContext, 30*time.Second)
	defer cancel()

	_, err := runner.Run( /*application.Context*/ ctx, "git", "commit", "--amend", "--no-edit",
		fmt.Sprintf("--date=\"%v\"", dateString))
	logging.Panic(err)
}
