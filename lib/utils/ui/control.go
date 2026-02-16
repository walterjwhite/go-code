package ui

import (
	"context"
	"image"
)

type Controller interface {
	Type(ctx context.Context, keys string) error
	Click(ctx context.Context, x, y float64) error

	Screenshot(ctx context.Context) (image.Image, error)
	ScreenshotOf(ctx context.Context, x, y, width, height float64) (image.Image, error)
	GetScreenSize(ctx context.Context) (int, int, error)
}

type Action interface {
	Do(ctx context.Context) error
}
