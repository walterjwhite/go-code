package citrix

import (
	"context"

	"fmt"

	"github.com/chromedp/chromedp"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"

	"strings"
	"time"
)

const (
	multipleMonitorPrompt = "Allow pop-ups from the browser Settings and restart the virtual app or desktop session to use multiple monitors."
	dontAskMeAgainText    = "Don't ask me again"
	okButtonText          = "OK"

	javascriptPromptScript = `(function(){
		matchingDialog = Array.from(document.querySelectorAll('div'))
		.find(d => d.textContent.trim() === '%s');
		if(matchingDialog != null && matchingDialog.parentElement != null) {
		  const childDiv = Array.from(matchingDialog.parentElement.children)
		  .find(el => el.textContent.trim() === 'OK');

		  button = Array.from(childDiv.children)
		  .find(el => el.textContent.trim() === 'OK');
		  button.click();
		  return true;
		}
		return false;
	})()`
)

func closePermissionPrompts(ctx context.Context) {
	permissionPromptCtx, permissionPromptCancel := context.WithTimeout(ctx, 1*time.Second)
	defer permissionPromptCancel()

	closePermissionPrompt(permissionPromptCtx, multipleMonitorPrompt)
}

func closePermissionPrompt(ctx context.Context, innerText string) {
	log.Info().Msgf("closePermissionPrompt - closing permission prompt if it exists: %s", innerText)

	var exists bool

	javascriptPromptScriptFormatted := fmt.Sprintf(javascriptPromptScript, strings.ReplaceAll(innerText, "'", "\\'"))

	log.Debug().Msgf("closePermissionPrompt - running javascript: %s", javascriptPromptScriptFormatted)
	err := chromedp.Run(ctx,
		chromedp.Evaluate(javascriptPromptScriptFormatted, &exists),
	)
	logging.Warn(err, false, "closePermissionPrompt - error closing prompt")

	log.Info().Msgf("closePermissionPrompt - prompt closed: %v", exists)
}
