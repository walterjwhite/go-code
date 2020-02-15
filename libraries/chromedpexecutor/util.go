package chromedpexecutor

import (
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/walterjwhite/go-application/libraries/logging"
	"strings"
)

// TODO: craigslist
func GetScript(line string) chromedp.Action {
	arguments := strings.Split(line, ",")

	if arguments[0] == "Click" {
		return chromedp.Click(arguments[1])
	} else if arguments[0] == "SendKeys" {
		return chromedp.SendKeys(arguments[1], arguments[2])
	}

	logging.Panic(fmt.Errorf("Unsupported action: %v", arguments[0]))
	return nil
}
