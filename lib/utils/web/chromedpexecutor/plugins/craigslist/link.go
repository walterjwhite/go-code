package craigslist

import (
	"context"

	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/session"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/session/remote"
)

func Accept(ctx context.Context, url string) {
	log.Info().Msgf("accept post: %v", url)

	s := remote.New(ctx)
	defer s.Cancel()

	session.Execute(s, chromedp.Navigate(url))
	session.Execute(s, chromedp.Click("//*[@id=\"new-edit\"]/div/div[2]/div[1]/button"))
}
