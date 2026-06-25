package graphical

import (
	"context"
	"errors"
	"fmt"
	"github.com/andreyvit/locateimage"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/utils/ui"
	"math"

	"image"
)

type ImageMatch struct {
	Ctx            context.Context
	Image          image.Image
	Match          *locateimage.Match
	MatchRegion    *MatchRegion
	MatchThreshold float64

	Controller ui.Controller
}

type MatchRegion struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
}

const (
	MaxCoordinateValue = 1 << 16 // 65536 pixels
)

func (i *ImageMatch) Matches() (bool, error) {
	if err := i.validate(); err != nil {
		return false, err
	}

	var screenshotImg image.Image
	var err error
	if i.MatchRegion == nil {
		screenshotImg, err = i.Controller.Screenshot(i.Ctx)
	} else {
		if i.MatchRegion.X > MaxCoordinateValue || i.MatchRegion.Y > MaxCoordinateValue ||
			i.MatchRegion.Width > MaxCoordinateValue || i.MatchRegion.Height > MaxCoordinateValue {
			return false, fmt.Errorf("coordinates exceed maximum allowed value of %d", MaxCoordinateValue)
		}
		screenshotImg, err = i.Controller.ScreenshotOf(i.Ctx, i.MatchRegion.X, i.MatchRegion.Y, i.MatchRegion.Width, i.MatchRegion.Height)
	}

	if err != nil {
		return false, err
	}

	match, err := locateimage.Find(i.contextForMatch(), locateimage.Convert(screenshotImg), i.Image, i.MatchThreshold, locateimage.Fastest)

	if errors.Is(err, locateimage.ErrNotFound) {
		i.debug(screenshotImg)

		return false, nil
	}

	log.Debug().Msg("match end - err")

	if err != nil {
		log.Warn().Msgf("Matches.err - %v", err)
		return false, err
	}

	i.Match = &match
	return i.Match != nil, nil
}

func (i *ImageMatch) contextForMatch() context.Context {
	if i.Ctx != nil {
		return i.Ctx
	}

	return context.Background()
}

func (i *ImageMatch) validate() error {
	if i.Controller == nil {
		return errors.New("matches validation failed: controller is nil")
	}

	if i.Image == nil {
		return errors.New("matches validation failed: search image is nil")
	}

	if math.IsNaN(i.MatchThreshold) || math.IsInf(i.MatchThreshold, 0) || i.MatchThreshold <= 0 || i.MatchThreshold > 1 {
		return fmt.Errorf("matches validation failed: threshold must be in (0, 1], got %v", i.MatchThreshold)
	}

	if i.MatchRegion != nil {
		if math.IsNaN(i.MatchRegion.X) || math.IsNaN(i.MatchRegion.Y) || math.IsNaN(i.MatchRegion.Width) || math.IsNaN(i.MatchRegion.Height) {
			return errors.New("matches validation failed: match region contains NaN")
		}

		if math.IsInf(i.MatchRegion.X, 0) || math.IsInf(i.MatchRegion.Y, 0) || math.IsInf(i.MatchRegion.Width, 0) || math.IsInf(i.MatchRegion.Height, 0) {
			return errors.New("matches validation failed: match region contains infinity")
		}

		if i.MatchRegion.Width <= 0 || i.MatchRegion.Height <= 0 {
			return fmt.Errorf("matches validation failed: region width/height must be > 0, got %v x %v", i.MatchRegion.Width, i.MatchRegion.Height)
		}

		if i.MatchRegion.X < 0 || i.MatchRegion.Y < 0 {
			return fmt.Errorf("matches validation failed: region coordinates must be >= 0, got X: %v, Y: %v", i.MatchRegion.X, i.MatchRegion.Y)
		}
	}

	return nil
}

func (i *ImageMatch) NotMatches() (bool, error) {
	matches, err := i.Matches()
	return !matches, err
}
