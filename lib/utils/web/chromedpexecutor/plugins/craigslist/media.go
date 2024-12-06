package craigslist

import (
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"
	"time"
)

const (
	uploadTimePerImage = 5 * time.Second
)

func (p *CraigslistPost) doMedia() []chromedp.Action {
	var actions []chromedp.Action

	if len(p.Media) > 0 {
		log.Info().Msgf("has %v images to upload", len(p.Media))

		actions = append(actions, chromedp.SetUploadFiles("//input[@type = 'file']", p.Media, chromedp.NodeVisible))

		actions = append(actions, chromedp.Sleep(p.sleepTime(len(p.Media))))

		actions = append(actions, chromedp.Click("/html/body/article/section/form/button"))
	}

	return actions
}

func (p *CraigslistPost) sleepTime(imageCount int) time.Duration {
	return uploadTimePerImage * time.Duration(imageCount)
}
