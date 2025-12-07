package agent

import (
	"context"
	"github.com/rs/zerolog/log"

	"github.com/mitchellh/go-homedir"
	"github.com/walterjwhite/go-code/lib/time/until"

	"bufio"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/utils/ui/windows"
	"os"
	"time"
)


func (c *Conf) PostLoad(ctx context.Context) error {
	c.read()
	return nil
}

func (c *Conf) read() {
	path, err := homedir.Expand(c.QuestionFile)
	logging.Panic(err)

	file, err := os.Open(path)
	logging.Panic(err)

	defer close(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		c.questions = append(c.questions, scanner.Text())
	}

	logging.Panic(scanner.Err())
}

func close(file *os.File) {
	logging.Warn(file.Close(), "read.close")
}

func (c *Conf) Init(ctx context.Context, headless bool, contextuals ...interface{}) error {
	log.Info().Msgf("agent.Init.ctx: %v", ctx)
	c.contextuals = contextuals
	c.processContextuals()

	log.Info().Msg("agent.Init.processedContextuals")

	err := c.launchBrowser(ctx)
	if err != nil {
		return err
	}

	log.Info().Msg("agent.Init.launchBrowser")
	err = c.navigateToUrl(ctx)
	if err != nil {
		return err
	}

	log.Info().Msg("agent.Init.navigateToUrl")
	return c.waitFor2FAToComplete(ctx)
}

func (c *Conf) processContextuals() {
	for _, ctx := range c.contextuals {
		if conf, ok := ctx.(*windows.WindowsConf); ok {
			c.handleWindowsConf(conf)
		}
	}
}

func (c *Conf) handleWindowsConf(conf *windows.WindowsConf) {
	c.WindowsConf = conf
}

func (c *Conf) waitForWindowsToLoad(pctx context.Context) error {
	ctx, cancel := context.WithTimeout(pctx, 2*time.Minute)
	defer cancel()

	i, err := c.WindowsConf.WindowsStartButtonMatcher(ctx)
	if err != nil {
		return err
	}

	return until.Until(ctx, 1*time.Second, i.Matches)
}
