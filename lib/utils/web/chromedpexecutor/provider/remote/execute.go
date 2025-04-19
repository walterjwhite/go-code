package remote

import (
	"context"
	"flag"
	"os/exec"
	"time"

	"github.com/walterjwhite/go-code/lib/application/logging"
)

var (
	launchRemoteBrowserCmdFlag = flag.String("browser-remote-cmd", "remote-web-browser", "Command to launch remote browser")

	browserRemoteCmd *exec.Cmd
)

func launchRemoteBrowser(ctx context.Context) {
	killRemoteBrowser()

	if !isRemoteBrowserRunning() {
		browserRemoteCmd = exec.CommandContext(ctx, *launchRemoteBrowserCmdFlag)
		logging.Panic(browserRemoteCmd.Start())

		time.Sleep(5 * time.Second)
	}
}

func killRemoteBrowser() {
	if browserRemoteCmd != nil {
		logging.Panic(browserRemoteCmd.Process.Kill())
		browserRemoteCmd = nil
	}
}
