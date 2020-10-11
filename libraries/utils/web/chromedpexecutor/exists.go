package chromedpexecutor

import (
	"context"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

func (s *ChromeDPSession) Exists(selector interface{}) bool {
	var existingNodes []*cdp.Node

	s.Execute(
		chromedp.Query(selector, chromedp.AtLeast(0),
			chromedp.After(func(i context.Context, n ...*cdp.Node) error {
				existingNodes = n
				return nil
			}),
		))

	return len(existingNodes) > 0
}
