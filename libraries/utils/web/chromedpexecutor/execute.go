package chromedpexecutor

import (
	"context"
	"errors"
	"flag"
	"github.com/chromedp/chromedp"
	"github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/application/logging"

	"github.com/walterjwhite/go-application/libraries/time/sleep"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

type ChromeDPSession struct {
	Context context.Context
	Waiter  sleep.Waiter

	CancelAllocator context.CancelFunc
	CancelContext   context.CancelFunc
}

var (
	devToolsWsUrlFlag          = flag.String("u", "", "Dev Tools WS URL")
	devToolsWsFileFlag         = flag.String("s", "~/.remote-browser-sessions", "Remote Browser Session File")
	launchRemoteBrowserCmdFlag = flag.String("browser-remote-cmd", "launch-remote-control-browser", "Command to launch remote browser")

	// TODO: add flags to tweak the deviation and minimum wait times
	// OR if a fixed delay is preferred
	waiter sleep.Waiter
	cmd    *exec.Cmd
)

func init() {
	waiter = &sleep.RandomDelay{MinimumDelay: 250, Deviation: 5000}
}

func IsRemoteBrowserRunning() bool {
	getURLFromFile()
	return len(*devToolsWsUrlFlag) > 0
}

func LaunchRemoteBrowser(ctx context.Context) *ChromeDPSession {
	if cmd != nil {
		logging.Panic(cmd.Process.Kill())
		cmd = nil
	}

	if !IsRemoteBrowserRunning() {
		cmd = exec.CommandContext(ctx, *launchRemoteBrowserCmdFlag)
		logging.Panic(cmd.Start())

		// ensure browser has had enough time to initialize
		time.Sleep(5 * time.Second)
	}

	return New(ctx)
}

func New(ctx context.Context) *ChromeDPSession {
	if len(*devToolsWsFileFlag) > 0 {
		getURLFromFile()
	}

	if len(*devToolsWsUrlFlag) == 0 {
		logging.Panic(errors.New("Please specify the DevToolsWSUrl"))
	}

	actxt, cancelActxt := chromedp.NewRemoteAllocator(ctx, *devToolsWsUrlFlag)

	// TODO: we might want to use an existing tab ...
	// create new tab
	ctx, cancelCtxt := chromedp.NewContext(actxt)

	return &ChromeDPSession{Context: ctx, CancelAllocator: cancelActxt, CancelContext: cancelCtxt, Waiter: waiter}
}

func getURLFromFile() {
	// check if the file exists
	f, err := homedir.Expand(*devToolsWsFileFlag)
	logging.Panic(err)

	_, err = os.Stat(f)
	if err == nil {
		data, err := ioutil.ReadFile(f)
		logging.Panic(err)

		dataString := strings.TrimSpace(string(data))
		devToolsWsUrlFlag = &dataString

		log.Debug().Msgf("URL: %v", dataString)
		log.Debug().Msgf("URL: %v", *devToolsWsUrlFlag)
	} else {
		logging.Panic(err)
	}
}

func (s *ChromeDPSession) Execute(actions ...chromedp.Action) {
	for i, action := range actions {
		log.Info().Msgf("running %v", action)
		logging.Panic(chromedp.Run(s.Context, action))

		s.doWait(i, actions...)
	}
}

func (s *ChromeDPSession) doWait(i int, actions ...chromedp.Action) {
	if s.Waiter != nil {
		if i < (len(actions) - 1) {
			s.Waiter.Wait()
		}
	}
}

func (s *ChromeDPSession) Cancel() {
	s.CancelAllocator()
	s.CancelContext()
}
