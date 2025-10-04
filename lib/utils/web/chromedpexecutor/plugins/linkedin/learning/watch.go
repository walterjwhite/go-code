package learning

import (
	"context"
	"errors"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"

	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"
	"time"
)

const (
	loadDelay = 10 * time.Second

	tryAgainClick = `(
	function() {
		tryAgainButton = Array.from(document.querySelectorAll('*'))
			.find(d => d.textContent.trim() === 'Try Again');

		if(tryAgainButton !== undefined) {
			tryAgainButton.click();
			return true;
		}

		return false;
	}
	)()`
)

func (s *Session) watch(course *Course) {
	log.Info().Msgf("watching: %s", course.Title)

	err := s.fetch(course)
	if err != nil {
		logging.Warn(err, false, "watch.fetch")
		return
	}

	err = s.watchContent(course)
	if err != nil {
		logging.Warn(err, false, "watch.watchContent")
		return
	}
}

func (s *Session) fetch(course *Course) error {
	ctx, cancel := context.WithTimeout(s.ctx, *s.StepTimeout)
	defer cancel()

	return action.Execute(ctx,
		chromedp.Navigate(course.Url),
		chromedp.WaitReady(watchMediaContainer))
}

func (s *Session) watchContent(course *Course) error {
	if course.Duration == 0 {
		log.Warn().Msg("watchContent - defaulting to 30m course timeout as no value was provided")
		course.Duration = 30 * time.Minute
	}

	ctx, cancel := context.WithTimeout(s.ctx, course.Duration)
	defer cancel()

	errorTimeStamps := []time.Time{}
	errorTimeWindow := 3 * time.Minute

	var err error
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("session.watchContent - courseDuration exceeded: %v -> %v", course.Title, course.Duration)
		default:
		}

		switch {
		case s.isContinueCourse():
			log.Info().Msg("watchContent.isContinueCourse")
		case s.tryAgain():
			log.Info().Msgf("watchContent.tryAgain")
		case s.isSurvey():
			log.Info().Msg("watchContent.isSurvey")
		case s.isQuiz():
			log.Info().Msg("watchContent.isQuiz")
			err = s.quizNext()
		case s.isCoderPad():
			log.Info().Msg("watchContent.isCoderPad")
			err = s.coderPadNext()
		case s.isText():
			log.Info().Msg("watchContent.isText")
			exists, err := s.textNext()
			if !exists && err == nil {
				err = errors.New("watchContent.isText.error")
			}
		case s.isVideo():
			log.Info().Msg("watchContent.isVideo")
			err = s.waitPlaying(course)
		default:
			err = errors.New("watchContent.other error")
		}

		if err != nil {
			err = s.onPlaybackError(course, &errorTimeStamps, errorTimeWindow, err)
			if err != nil {
				return err
			}
		}

		time.Sleep(watchIterationTimeout)
	}
}

func (s *Session) tryAgain() bool {
	ctx, cancel := context.WithTimeout(s.ctx, *s.StepTimeout)
	defer cancel()

	var exists bool
	err := chromedp.Run(ctx,
		chromedp.Evaluate(tryAgainClick, &exists))

	return err == nil && exists
}

func (s *Session) courseCompletion(course *Course) {
	log.Info().Msgf("watch.courseCompletion - watched: %s", course.Title)

	if s.WasCompleted(course) {
		log.Info().Msgf("watch.courseCompletion - completed: %s", course.Title)
	} else {
		log.Info().Msgf("watch.courseCompletion - NOT completed: %s", course.Title)
	}
}

func (s *Session) onPlaybackError(course *Course, errorTimeStamps *[]time.Time, errorTimeWindow time.Duration, err error) error {
	log.Warn().Msgf("watch.onPlaybackError - %s - watchContent: %v", course.Title, err)
	action.Screenshot(s.ctx, fmt.Sprintf("/tmp/linkedinlearning-error-%v-%v.png", time.Now().Unix(), course.Title))

	*errorTimeStamps = append(*errorTimeStamps, time.Now())

	cutoff := time.Now().Add(-errorTimeWindow)
	for i := 0; i < len(*errorTimeStamps); i++ {
		if (*errorTimeStamps)[i].Before(cutoff) {
			*errorTimeStamps = (*errorTimeStamps)[i+1:]
			break
		}
	}

	if len(*errorTimeStamps) > 3 {
		return fmt.Errorf("watch.onPlaybackError - %s - more than 3 errors occurred within %v", course.Title, errorTimeWindow)
	}

	return nil
}
