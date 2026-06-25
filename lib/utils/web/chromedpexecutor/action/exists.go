package action

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"

	"time"
)

const (
	dimensionQuery = "(function(){\n" +
		"e = %s;if(e === null) {return false;} rect = e.getBoundingClientRect(); return rect.height > 0 && rect.width > 0;})()\n"

	existsQueryTimeout = 500 * time.Millisecond
)

func ExistsById(pctx context.Context, id string) bool {
	log.Debug().Msgf("checking for visibility of: %s via %s", id, pctx)

	ctx, cancel := context.WithTimeout(pctx, existsQueryTimeout)
	defer cancel()

	var exists bool
	err := chromedp.Run(ctx,
		chromedp.Evaluate(fmt.Sprintf("document.getElementById('%s') !== null", id), &exists),
	)

	log.Debug().Msgf("exists: %v", exists)
	return err == nil && exists
}

func ExistsByCssSelector(pctx context.Context, cssSelector string) bool {
	log.Debug().Msgf("checking for visibility of: %s via %s", cssSelector, pctx)

	ctx, cancel := context.WithTimeout(pctx, existsQueryTimeout)
	defer cancel()

	var exists bool
	err := chromedp.Run(ctx,
		chromedp.Evaluate(fmt.Sprintf("document.querySelector('%s') !== null", cssSelector), &exists),
	)

	log.Debug().Msgf("exists: %v", exists)
	return err == nil && exists
}

func ExistsByXPath(pctx context.Context, selector string) bool {
	log.Debug().Msgf("checking for visibility of: %s via %s", selector, pctx)

	ctx, cancel := context.WithTimeout(pctx, existsQueryTimeout)
	defer cancel()

	var exists bool
	err := chromedp.Run(ctx,
		chromedp.Evaluate(fmt.Sprintf("document.evaluate('%s', document, null, XPathResult.FIRST_ORDERED_NODE_TYPE, null).singleNodeValue !== null", selector), &exists),
	)

	log.Debug().Msgf("exists: %v", exists)
	return err == nil && exists
}

func HasDimensions(pctx context.Context, selector string) bool {
	log.Debug().Msgf("checking for visibility of: %s via %s", selector, pctx)

	ctx, cancel := context.WithTimeout(pctx, existsQueryTimeout)
	defer cancel()

	var exists bool
	err := chromedp.Run(ctx,
		chromedp.Evaluate(fmt.Sprintf(dimensionQuery, selector), &exists),
	)

	log.Debug().Msgf("exists: %v", exists)
	return err == nil && exists
}









