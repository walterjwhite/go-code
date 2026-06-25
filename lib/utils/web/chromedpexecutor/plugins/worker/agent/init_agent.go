package agent

import (
	"context"
	"github.com/rs/zerolog/log"

	"github.com/mitchellh/go-homedir"
	"github.com/walterjwhite/go-code/lib/time/until"

	"bufio"
	"fmt"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/utils/ui/windows"
	"os"
	"path/filepath"
	"strings"
	"time"
)


func (c *Conf) PostLoad(ctx context.Context) error {
	c.read()
	return nil
}

func (c *Conf) read() {
	path, err := homedir.Expand(c.QuestionFile)
	logging.Error(err)

	if err := validateFilePath(path); err != nil {
		logging.Error(err)
		return
	}

	file, err := os.Open(path) // #nosec G304 - path is validated by validateFilePath
	logging.Error(err)

	defer closeResource(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		c.questions = append(c.questions, scanner.Text())
	}

	logging.Error(scanner.Err())
}

func validateFilePath(path string) error {
	if strings.Contains(path, "..") {
		return fmt.Errorf("path traversal detected in file path: %s", path)
	}

	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("cannot access file: %w", err)
	}

	if !info.Mode().IsRegular() {
		return fmt.Errorf("path is not a regular file: %s", path)
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("cannot resolve absolute path: %w", err)
	}

	homeDir, err := homedir.Dir()
	if err == nil && !strings.HasPrefix(absPath, homeDir) {
		log.Warn().Msgf("Question file is outside home directory: %s", absPath)
	}

	return nil
}

func closeResource(file *os.File) {
	logging.Warn(file.Close(), "read.close")
}

func (c *Conf) Init(ctx context.Context, headless bool, contextuals ...any) error {
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
