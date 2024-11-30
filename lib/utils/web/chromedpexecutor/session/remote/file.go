package remote

import (
	"flag"

	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

var (
	devToolsWsFileFlag = flag.String("s", "/tmp/w/walterjwhite/app/default/web-browser/remote-web-browser/instance/remote", "Remote Browser Session File")
)

func getURLFromFile() {
	// check if the file exists
	f, err := homedir.Expand(*devToolsWsFileFlag)
	logging.Panic(err)

	_, err = os.Stat(f)
	if err == nil {
		data, err := os.ReadFile(f)
		logging.Panic(err)

		log.Info().Msg("getting URL from file ...")

		// ws url is on line 2
		dataString := strings.TrimSpace(strings.Split(string(data), "\n")[0])
		devToolsWsUrlFlag = &dataString

		log.Info().Msgf("Dev Tools WS URL: %v", *devToolsWsUrlFlag)
	}
}

func isRemoteBrowserRunning() bool {
	getURLFromFile()
	return len(*devToolsWsUrlFlag) > 0
}
