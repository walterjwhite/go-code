package citrix

import (
	"context"
	"fmt"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"
	"strings"
)

func validateToken(token string) {
	if len(token) != 6 {
		logging.Panic(fmt.Errorf("please enter the 6-digit token"))
	}
}

func isExpired(ctx context.Context) bool {
	currentUrl := action.Location(ctx)
	if strings.HasSuffix(currentUrl, "/logout.html") {
		return true
	}

	return strings.HasSuffix(currentUrl, "LogonPoint/tmindex.html")
}
