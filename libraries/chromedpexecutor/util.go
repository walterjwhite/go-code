package chromedpexecutor

import (
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/walterjwhite/go-application/libraries/logging"
	"strings"
	"time"
)

// TODO: craigslist
func ParseActions(lines ...string) []chromedp.Action {
	var actions []chromedp.Action

	for _, line := range lines {
		action := ParseAction(line)
		if action != nil {
			actions = append(actions, action)
		}
	}

	return actions
}

func ParseAction(line string) chromedp.Action {
	arguments := strings.Split(line, ",")

	if arguments[0] == "Click" {
		return chromedp.Click(arguments[1])
	} else if arguments[0] == "SendKeys" {
		return chromedp.SendKeys(arguments[1], arguments[2])
	} else if arguments[0] == "Clear" {
		return chromedp.Clear(arguments[1])
	} else if arguments[0] == "Key" {
		return chromedp.KeyEvent(arguments[1])
	} else if arguments[0] == "SetValue" {
		return chromedp.SetValue(arguments[1], arguments[2])
	} else if arguments[0] == "ScrollIntoView" {
		return chromedp.ScrollIntoView(arguments[1])
	} else if arguments[0] == "Sleep" {
		d, err := time.ParseDuration(arguments[1])
		logging.Panic(err)

		return chromedp.Sleep(d)
	}

	// TODO: + Mouse ...

	logging.Panic(fmt.Errorf("Unsupported action: %v", arguments[0]))
	return nil
}
