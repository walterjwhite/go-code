package craigslist

import (
	"context"

	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/session"
)

func Accept(ctx context.Context, url string, session session.ChromeDPSession) {
	log.Info().Msgf("accept post: %v", url)

	s := chromedpexecutor.New(ctx)
	defer s.Cancel()

	s.Execute(chromedp.Navigate(url))
	s.Execute(chromedp.Click("//*[@id=\"new-edit\"]/div/div[2]/div[1]/button"))
}
