package citrix

import (
	"context"

	"fmt"

	"github.com/chromedp/chromedp"

	"github.com/rs/zerolog/log"

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

func (i *Instance) closePermissionPrompts() {
	permissionPromptCtx, permissionPromptCancel := context.WithTimeout(i.ctx, 1*time.Second)
	defer permissionPromptCancel()

	i.closePermissionPrompt(permissionPromptCtx, multipleMonitorPrompt)
}

func (i *Instance) closePermissionPrompt(ctx context.Context, innerText string) {
	log.Info().Msgf("%v - closePermissionPrompt - closing permission prompt if it exists: %s", i, innerText)

	var exists bool

	javascriptPromptScriptFormatted := fmt.Sprintf(javascriptPromptScript, strings.ReplaceAll(innerText, "'", "\\'"))

	log.Debug().Msgf("%v - closePermissionPrompt - running javascript: %s", i, javascriptPromptScriptFormatted)
	err := chromedp.Run(ctx,
		chromedp.Evaluate(javascriptPromptScriptFormatted, &exists),
	)
	if err != nil {
		log.Warn().Msgf("%v - closePermissionPrompt - error closing prompt - %s", i, err.Error())
	}

	log.Info().Msgf("%v - closePermissionPrompt - prompt closed: %v", i, exists)
}
