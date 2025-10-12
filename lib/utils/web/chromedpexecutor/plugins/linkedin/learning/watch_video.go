package learning

import (
	"context"
	"errors"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"

	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"
	"time"
)

const (
	hasElementWithText   = `Array.from(document.querySelectorAll('*')).find(e => e.textContent.trim() === '%s') !== undefined`
	clickElementWithText = `(
		function() {
			e = Array.from(document.querySelectorAll('*')).find(e => e.textContent.trim() === '%s');
			if(e != undefined) {
				e.click();
				return true;
			}

			return false;
		})()
	`

	watchMediaContainer = ".classroom-media-screens"
	watchMediaState     = "classroom-layout__stage"
	mediaProgress       = ".vjs-progress-holder"

	mediaProgressQuery = `Array.from(document.getElementsByClassName("vjs-progress-holder")).map(el => el.getAttribute("aria-valuenow").trim())`


	watchIterationTimeout = 5 * time.Second
)

func (s *Session) isVideo() bool {
	ctx, cancel := context.WithTimeout(s.ctx, *s.StepTimeout)
	defer cancel()

	err := action.Execute(ctx,
		chromedp.WaitReady(mediaProgress))

	return err == nil
}

func (s *Session) waitPlaying(course *Course) error {
	progress := "NOT YET FETCHED"

	for {
		currentProgress, err := s.fetchMediaProgress()
		if err != nil {
			return err
		}

		if isNotPlaying(progress, currentProgress) {
			log.Warn().Msgf("watch.video.waitPlaying: content appears to be stopped: %s | %s", currentProgress, progress)
			return fmt.Errorf("watch.video.waitPlaying: content appears to be stopped: %s | %s", currentProgress, progress)
		}

		log.Debug().Msg("watch.video.waitPlaying: playing, sleeping")
		time.Sleep(watchIterationTimeout)
		log.Debug().Msg("watch.video.waitPlaying: finished sleep")

		progress = currentProgress
	}
}

func isNotPlaying(progress, currentProgress string) bool {
	return progress == currentProgress || len(currentProgress) == 0
}

func (s *Session) fetchMediaProgress() (string, error) {
	currentProgress, err := s.extract(mediaProgressQuery)
	if err != nil {
		return "", err
	}

	if len(currentProgress) == 0 {
		return "", errors.New("media progress not found: currentProgress")
	}

	log.Debug().Msgf("currentProgress: %s", currentProgress[0])
	return currentProgress[0], nil
}


func (s *Session) extract(expression string) ([]string, error) {
	ctx, cancel := context.WithTimeout(s.ctx, extractTimeout)
	defer cancel()

	var values []string

	return values, chromedp.Run(ctx,
		chromedp.Evaluate(expression, &values))
}
