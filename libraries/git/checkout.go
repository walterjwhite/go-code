package git

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/runner"
)

// TODO: this is not portable ...
func Checkout(parentContext context.Context, projectName string /*branchName, targetDirectory*/, optionalArguments ...map[string]string) {
	ctx, cancel := context.WithTimeout(parentContext, 30*time.Second)
	defer cancel()

	arguments := []string{"clone", projectName}

	for _, optionalMap := range optionalArguments {
		for key, value := range optionalMap {
			if key == "branch" {
				arguments = append(arguments, "-b", value)
			} else if key == "targetDirectory" {
				arguments = append(arguments, value)
			} else {
				log.Warn().Msgf("Unrecognized argument: %v / %v", key, value)
			}
		}
	}

	_, err := runner.Run(ctx, "checkout-project", projectName)
	//_, err := runner.Run(ctx, "git", arguments...)
	logging.Panic(err)
}
