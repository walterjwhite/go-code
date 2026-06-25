package chromedpexecutor

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/walterjwhite/go-code/lib/utils/ui/graphical"
	"github.com/walterjwhite/go-code/lib/utils/web/chromedpexecutor/action"
	"image"
)

type ChromeDPController struct {
}

func (c *ChromeDPController) Type(ctx context.Context, keys string) error {
	return action.Execute(ctx,
		chromedp.KeyEvent(keys))
}

func (c *ChromeDPController) Click(ctx context.Context, x, y float64) error {
	return action.Execute(ctx,
		chromedp.MouseClickXY(x, y))
}

func (c *ChromeDPController) Screenshot(ctx context.Context) (image.Image, error) {
	bytes, err := action.TakeScreenshot(ctx)

	if err != nil {
		return nil, err
	}

	return graphical.BytesToImage(bytes)
}

func (c *ChromeDPController) ScreenshotOf(ctx context.Context, x, y, width, height float64) (image.Image, error) {
	bytes, err := action.TakeScreenshotOf(ctx, x, y, width, height)

	if err != nil {
		return nil, err
	}

	return graphical.BytesToImage(bytes)
}

func (c *ChromeDPController) GetScreenSize(ctx context.Context) (int, int, error) {
	size, err := action.GetWindowSize(ctx)
	return size.Width, size.Height, err
}
