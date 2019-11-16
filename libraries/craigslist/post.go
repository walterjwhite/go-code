package craigslist

import (
	"context"
	"flag"
	"github.com/chromedp/chromedp"
	"github.com/walterjwhite/go-application/libraries/logging"

	"github.com/rs/zerolog/log"
	//"os"
)

const (
	craigslistBasePostUrl = "https://post.craigslist.org/c/"
)

var (
	minimumDelayBetweenActionsFlag = flag.Int("CraigslistMinimumDelayBetweenActions", 250, "Minimum Delay between actions (ms)")
	deviationBetweenActionsFlag = flag.Int("CraigslisDeviationBetweenActions", 2500, "Deviation between actions (ms)")
	
	//delayBetweenActions     time.Duration
	delay *RandomDelay
)

func init() {
	//var err error

	delay = &RandomDelay{MinimumDelay: *minimumDelayBetweenActionsFlag, Deviation: *deviationBetweenActionsFlag}
	
	//delayBetweenActions, err = time.ParseDuration(*delayBetweenActionsFlag)
	//logging.Panic(err)
}

func (p *CraigslistPost) Create(ctx context.Context) {
	log.Info().Msgf("post: %v", p)
	
	p.execute(ctx, chromedp.Navigate(craigslistBasePostUrl+p.Region))

	p.execute(ctx, p.doForSaleBy()...)
	p.execute(ctx, p.doCategory()...)

	p.execute(ctx, p.doPostDetails()...)
	p.execute(ctx, p.doPhone()...)
	p.execute(ctx, p.doMedia()...)
	p.execute(ctx, p.publish()...)
}

func (p *CraigslistPost) publish() []chromedp.Action {
	return []chromedp.Action{chromedp.Click("//*[@id=\"publish_top\"]/button")}
}

// TODO: this is generic code, unrelated to craigslist
// move this out into chromedp helper ...
func (p *CraigslistPost) execute(ctx context.Context, actions ...chromedp.Action) {

	for i, action := range actions {
		log.Info().Msgf("running %v", action)
		logging.Panic(chromedp.Run(ctx, action))

		if i < (len(actions) - 1) {
			delay.Wait()
		}
	}
}

func (p *CraigslistPost) HasDefault() bool {
	return false
}

func (p *CraigslistPost) Refreshable() bool {
	return false
}

func (p *CraigslistPost) EncryptedFields() []string {
	return nil
}
