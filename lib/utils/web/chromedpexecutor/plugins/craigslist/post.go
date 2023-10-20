package craigslist

import (
	"context"
	"flag"
	"github.com/chromedp/chromedp"

	"github.com/walterjwhite/go-code/lib/time/delay"
	
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/session"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/session/remote"

	"github.com/rs/zerolog/log"
	"time"
)

const (
	craigslistBasePostUrl = "https://post.craigslist.org/c/"
)

var (
	minimumDelayBetweenActionsFlag = flag.Int("CraigslistMinimumDelayBetweenActions", 250, "Minimum Delay between actions (ms)")
	deviationBetweenActionsFlag    = flag.Int("CraigslisDeviationBetweenActions", 5000, "Deviation between actions (ms)")
	//devToolsWsUrlFlag              = flag.String("DevToolsWsUrl", "", "Dev Tools WS URL")

	//delayBetweenActions     time.Duration
	delayer delay.Delayer
)

func init() {
	//var err error

	delayer = delay.NewRandom(time.Duration(*minimumDelayBetweenActionsFlag)*time.Millisecond, time.Duration(*deviationBetweenActionsFlag)*time.Millisecond)

	//delayBetweenActions, err = time.ParseDuration(*delayBetweenActionsFlag)
	//logging.Panic(err)
}

func (p *CraigslistPost) Create(ctx context.Context) {
	log.Info().Msgf("post: %v", p)

	p.session = remote.New(ctx)

	session.Execute(p.session, chromedp.Navigate(craigslistBasePostUrl + p.Region))

	session.Execute(p.session, p.doForSaleBy()...)
	session.Execute(p.session, p.doCategory()...)

	session.Execute(p.session, p.doPostDetails()...)
	session.Execute(p.session, p.doPhone()...)
	session.Execute(p.session, p.doMedia()...)
	session.Execute(p.session, p.publish()...)
}

func (p *CraigslistPost) publish() []chromedp.Action {
	return []chromedp.Action{chromedp.Click("//*[@id=\"publish_top\"]/button")}
}

func Delay() {
	delayer.Delay()
}
