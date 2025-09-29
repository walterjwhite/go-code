package windows

import (
	"context"
	"github.com/chromedp/cdproto/input"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
	"github.com/rs/zerolog/log"
	"time"
)

func SendCtrlAltDelete(ctx context.Context) error {
	log.Warn().Msg("sending ctrl+alt+delete")

	return chromedp.Run(ctx,
		input.DispatchKeyEvent(input.KeyDown).WithKey(kb.Control).WithModifiers(input.ModifierCtrl),
		chromedp.Sleep(200*time.Millisecond),

		input.DispatchKeyEvent(input.KeyDown).WithKey(kb.Alt).WithModifiers(input.ModifierAlt+input.ModifierCtrl),
		chromedp.Sleep(200*time.Millisecond),

		input.DispatchKeyEvent(input.KeyDown).WithKey(kb.Delete).WithModifiers(input.ModifierAlt+input.ModifierCtrl),
		chromedp.Sleep(500*time.Millisecond),

		input.DispatchKeyEvent(input.KeyUp).WithKey(kb.Delete).WithModifiers(input.ModifierAlt+input.ModifierCtrl),
		input.DispatchKeyEvent(input.KeyUp).WithKey(kb.Alt).WithModifiers(input.ModifierAlt+input.ModifierCtrl),
		input.DispatchKeyEvent(input.KeyUp).WithKey(kb.Control).WithModifiers(input.ModifierCtrl))

}

func Unlock(ctx context.Context, password string, delayBeforeEnteringPassword time.Duration) error {
	err := SendCtrlAltDelete(ctx)
	if err != nil {
		return err
	}

	time.Sleep(delayBeforeEnteringPassword)

	err = chromedp.Run(ctx,
		chromedp.KeyEvent(password),
		chromedp.KeyEvent(kb.Enter))
	if err != nil {
		return err
	}

	return nil
}

func Lock(ctx context.Context, delayBeforeLocking time.Duration) error {
	err := SendCtrlAltDelete(ctx)
	if err != nil {
		return err
	}

	time.Sleep(delayBeforeLocking)

	err = chromedp.Run(ctx,
		chromedp.KeyEvent(kb.Enter))
	if err != nil {
		return err
	}

	return nil
}

func ToggleWindowsStartMenu(ctx context.Context) error {
	return chromedp.Run(ctx, chromedp.KeyEvent(kb.Meta))
}
