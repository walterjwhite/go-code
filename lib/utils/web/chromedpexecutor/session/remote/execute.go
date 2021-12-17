package remote

import (
	"context"
	"flag"
	"os/exec"
	"time"

	"github.com/walterjwhite/go-code/lib/application/logging"
)

var (
	launchRemoteBrowserCmdFlag = flag.String("browser-remote-cmd", "launch-remote-chrome", "Command to launch remote browser")

	browserRemoteCmd *exec.Cmd
)

func launchRemoteBrowser(ctx context.Context) {
	killRemoteBrowser()

	if !isRemoteBrowserRunning() {
		browserRemoteCmd = exec.CommandContext(ctx, *launchRemoteBrowserCmdFlag)
		logging.Panic(browserRemoteCmd.Start())

		// ensure browser has had enough time to initialize
		time.Sleep(5 * time.Second)
	}
}

func killRemoteBrowser() {
	if browserRemoteCmd != nil {
		logging.Panic(browserRemoteCmd.Process.Kill())
		browserRemoteCmd = nil
	}
}
