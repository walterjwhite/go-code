package action

import (
	"context"

	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"
)

const (
	fullscreenJavaScript = `(
		function(){
			if(document.documentElement.requestFullscreen) {
				document.documentElement.requestFullscreen();
				return true;
			}
			
			return false;
		}
	)()`

	exitFullscreenJavaScript = `(
		return document.exitFullscreen().then(() => true)
			.catch(err => {
				console.error(err);
				return false;
			});
		}
		)()`
)

func Fullscreen(ctx context.Context) bool {
	var successful bool

	err := chromedp.Run(ctx, chromedp.Evaluate(fullscreenJavaScript, &successful))
	if err != nil || !successful {
		log.Warn().Msg("attempting F11")
		err = chromedp.Run(ctx, chromedp.KeyEvent("\uE03B"))

		return err == nil
	}

	return true
}

func ExitFullscreen(ctx context.Context) bool {
	var successful bool

	err := chromedp.Run(ctx, chromedp.Evaluate(exitFullscreenJavaScript, &successful))
	if err != nil || !successful {
		log.Warn().Msg("attempting F11")
		err = chromedp.Run(ctx, chromedp.KeyEvent("\uE03B"))

		return err == nil
	}

	return true
}
