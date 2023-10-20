package session

import (
	"context"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"

	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/session"
)

func Exists(s session.ChromeDPSession, selector interface{}) (bool, error) {
	var existingNodes []*cdp.Node

	err := chromedp.Run(s.Context(),
		chromedp.Query(selector, chromedp.AtLeast(0),
			chromedp.After(func(i context.Context, executionId runtime.ExecutionContextID, n ...*cdp.Node) error {
				existingNodes = n
				return nil
			}),
		))

	if err != nil {
		return false, err
	}

	return len(existingNodes) > 0, nil
}
