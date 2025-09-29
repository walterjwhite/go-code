package learning

import (
	"context"
	"github.com/chromedp/chromedp"

	"github.com/rs/zerolog/log"

	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"
	"time"
)

const (
	searchResultsCount     = ".search-body__search-result-count"
	searchResultsContainer = "/html/body/div[8]/div[2]/div[2]/main/div[2]/div[2]/div"

	extractTimeout      = 1 * time.Second
	delayBetweenScrolls = 5 * time.Second

	coursesNotCompleted = `
		Array.from(document.getElementsByClassName('lls-card-detail-card__main'))
			.filter(el => !el.textContent.includes('Completed'))
			.map(element => {
				return {
					title: element.children[1].children[0].children[0].children[0].children[0].children[0].textContent.trim(),
					url: element.children[1].children[0].children[0].children[0].children[0].children[0].href,
					durationstring: element.children[0].children[0].children[0].children[0].children[3].textContent.trim().replace(/\s+/g, ''),
				};
			}
		)
	`

	courseQueueSize     = 10
	linkedInLearningUrl = "https://www.linkedin.com/learning"
)

func (s *Session) Search(searchTerm string) []*Course {
	log.Info().Msg("Session.Search - start")

	if len(searchTerm) == 0 {
		log.Warn().Msg("Session.Search - No search term specified")
		return []*Course{}
	}

	log.Info().Str("searchTerm", searchTerm).Msg("Session.Search - searching")

	ctx, cancel := context.WithTimeout(s.ctx, *s.StepTimeout)
	defer cancel()

	err := action.Execute(ctx,
		chromedp.Navigate(linkedInLearningUrl+"/search?entityType=COURSE&keywords="+searchTerm),
		chromedp.WaitReady(searchResultsCount))
	if err != nil {
		log.Panic().Err(err).Msg("Session.Search - error navigating to search page")
	}

	log.Info().Msg("Session.Search - entered search criteria, now, we just need to parse the results")

	return s.extractCourses()
}

func (s *Session) extractCourses() []*Course {
	lastCourseCount := 0

	var courses []*Course

	for {
		foundCourses := s.doExtractCourses()
		log.Info().Interface("courses", foundCourses).Msg("Session.extractCourses - non-completed courses")

		for _, course := range foundCourses {
			log.Info().Interface("course", &course).Msg("Session.extractCourses - found course")

			courses = append(courses, &course)
		}

		if len(courses) >= courseQueueSize {
			return courses
		}

		if len(courses) == lastCourseCount {
			log.Warn().Int("courseCount", lastCourseCount).Msg("Session.extractCourses - no more courses found")
			return courses
		}

		err := s.scrollToEnd()
		if err != nil {
			log.Warn().Err(err).Msg("Session.extractCourses - scrollToEnd - Error")
			return courses
		}

		time.Sleep(delayBetweenScrolls)
		lastCourseCount = len(courses)
	}

	return courses
}

func (s *Session) scrollToEnd() error {
	ctx, cancel := context.WithTimeout(s.ctx, extractTimeout)
	defer cancel()

	return action.End(ctx)
}

func (s *Session) doExtractCourses() []Course {
	ctx, cancel := context.WithTimeout(s.ctx, extractTimeout)
	defer cancel()

	var values []Course

	err := chromedp.Run(ctx,
		chromedp.Evaluate(coursesNotCompleted, &values),
	)

	if err != nil {
		log.Warn().Err(err).Msg("Session.doExtractCourses - error extracting courses")
		return nil
	}

	for i := range values {
		duration, err := time.ParseDuration(values[i].DurationString)
		if err != nil {
			log.Warn().Str("duration", values[i].DurationString).Msg("Session.doExtractCourses - unable to parse duration")
		} else {
			values[i].Duration = duration
		}
	}

	return values
}
