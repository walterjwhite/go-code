package build

import (
	"context"
	"log"
	"os"
	//	"github.com/walterjwhite/go-application/libraries/notify"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/maven"
	"github.com/walterjwhite/go-application/libraries/runner"
)

func Build(ctx context.Context, debug *bool /*, notifications notify.Notifier*/) {
	command, arguments := maven.GetCommandLine([]string{"clean", "install", "-Dmaven.test.skip=true", "-Dorg.slf4j.simpleLogger.log.org.apache.maven.cli.transfer.Slf4jMavenTransferListener=warn", "-B"}, debug)
	/*exitcode*/ _, err := runner.Run(ctx, command, arguments...)

	// TODO: rather than panic here, we should allow this to:
	// send an error to a channel which may raise a Windows 10 notification, dbus notification, etc.
	// finally panic
	logging.Panic(err)

	//	notification.Notify()
}

/*
func BuildNotification() notify.Notification {
	return notify.Notification{Id: "build", Title: "build complete", Details: fmt.Sprintf("build complete: %v", getBuildDirectory())}
}
*/

func GetBuildDirectory() string {
	dir, err := os.Getwd()
	logging.Panic(err)

	return dir
}
