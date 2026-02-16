package action

import (
	"context"

	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/chromedp"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

func Grant(ctx context.Context, permissions []browser.PermissionType) {
	logging.Error(chromedp.Run(ctx, browser.GrantPermissions(permissions)))
}
