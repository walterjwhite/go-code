package run

import (
	"fmt"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/security/secrets"
	"github.com/walterjwhite/go-code/lib/utils/token/providers/stdin"
	"strconv"
	"strings"
	"time"
)

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
	if len(line) == 0 || strings.HasPrefix(line, "#") {
		log.Debug().Msg("line is a comment or empty")
		return nil
	}

	arguments := strings.Split(line, ",")

	switch arguments[0] {
	case "navigate":
		return chromedp.Navigate(arguments[1])
	case "click":
		return chromedp.Click(arguments[1])
	case "mouseClick":
		x, err := strconv.ParseFloat(arguments[1], 64)
		logging.Error(err)

		y, err := strconv.ParseFloat(arguments[2], 64)
		logging.Error(err)

		return chromedp.MouseClickXY(x, y)
	case "sendKeys":
		return chromedp.SendKeys(arguments[1], arguments[2])
	case "sendKeysSecret":
		return chromedp.SendKeys(arguments[1], secret(arguments[2]))
	case "clear":
		return chromedp.Clear(arguments[1])
	case "key":
		return chromedp.KeyEvent(process(arguments[1]))
	case "exec":
		return &Exec{arguments[1], arguments[2:]}
	case "setValue":
		return chromedp.SetValue(arguments[1], arguments[2])
	case "scrollIntoView":
		return chromedp.ScrollIntoView(arguments[1])
	case "sleep":
		d, err := time.ParseDuration(arguments[1])
		logging.Error(err)

		return chromedp.Sleep(d)
	case "tickle":
		/*d*/ _, err := time.ParseDuration(arguments[1])
		logging.Error(err)


	case "waitVisible":
		return chromedp.WaitVisible(arguments[1])
	case "submit":
		return chromedp.Submit(arguments[1])
	default:
		if strings.Index(arguments[0], "#") == 0 || len(arguments[0]) == 0 {
			log.Debug().Msgf("Ignoring comment: %v", arguments[0])
			return nil
		}

		logging.Error(fmt.Errorf("unsupported action: %v", arguments[0]))
	}

	return nil
}

func process(value string) string {
	return getKeyFromString(stdIn(value))
}

func secret(key string) string {
	return strings.TrimSpace(secrets.Get(key))
}

func stdIn(value string) string {
	if strings.Index(value, "stdin:") == 0 {
		promptMessage := value[6:]
		s := stdin.StdInReader{PromptMessage: promptMessage}
		return s.Get()
	}

	return value
}

func getKeyFromString(key string) string {
	switch key {
	case "META":
		return kb.Meta
	case "SUPER":
		return kb.Super
	case "CONTROL":
		return kb.Control
	case "ALT":
		return kb.Alt
	case "SHIFT":
		return kb.Shift
	case "ENTER":
		return kb.Enter
	case "TAB":
		return kb.Tab
	default:
		return key
	}
}
