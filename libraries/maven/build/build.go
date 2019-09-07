package build

import (
	"context"
	"log"
	"os"
	//	"github.com/walterjwhite/go-application/libraries/notify"
	"github.com/walterjwhite/go-application/libraries/maven"
	"github.com/walterjwhite/go-application/libraries/runner"
)

func Build(ctx context.Context, debug *bool /*, notifications notify.Notifier*/) {
	command, arguments := maven.GetCommandLine([]string{"clean", "install", "-Dmaven.test.skip=true", "-Dorg.slf4j.simpleLogger.log.org.apache.maven.cli.transfer.Slf4jMavenTransferListener=warn", "-B"}, debug)
	runner.Run(ctx, command, arguments...)
	//	notification.Notify()
}

/*
func BuildNotification() notify.Notification {
	return notify.Notification{Id: "build", Title: "build complete", Details: fmt.Sprintf("build complete: %v", getBuildDirectory())}
}
*/

func getBuildDirectory() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return dir
}
