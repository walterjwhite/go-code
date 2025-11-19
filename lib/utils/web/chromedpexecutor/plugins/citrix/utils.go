package citrix

import (
	"context"

	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"
	"strings"
)

func IsContextExpired(ctx context.Context) bool {
	currentUrl := action.Location(ctx)
	if strings.HasSuffix(currentUrl, "/logout.html") {
		return true
	}

	return strings.HasSuffix(currentUrl, "LogonPoint/tmindex.html")
}
