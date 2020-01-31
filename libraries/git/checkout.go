package git

import (
	"context"
	"time"

	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/runner"
)

// TODO: this is not portable ...
func Checkout(parentContext context.Context, projectName string) {
	ctx, cancel := context.WithTimeout(parentContext, 30*time.Second)
	defer cancel()

	_, err := runner.Run(ctx, "checkout-project", projectName)
	logging.Panic(err)
}
