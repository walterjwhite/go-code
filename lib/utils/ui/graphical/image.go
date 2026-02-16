package graphical

import (
	"context"
	"errors"
	"github.com/andreyvit/locateimage"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/utils/ui"

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

func (i *ImageMatch) Matches() (bool, error) {
	var image image.Image
	var err error
	if i.MatchRegion == nil {
		image, err = i.Controller.Screenshot(i.Ctx)
	} else {
		image, err = i.Controller.ScreenshotOf(i.Ctx, i.MatchRegion.X, i.MatchRegion.Y, i.MatchRegion.Width, i.MatchRegion.Height)
	}

	if err != nil {
		return false, err
	}

	match, err := locateimage.Find(context.Background(), locateimage.Convert(image), i.Image, i.MatchThreshold, locateimage.Fastest)
	i.Match = &match

	if errors.Is(err, locateimage.ErrNotFound) {
		i.debug(image)

		return false, nil
	}

	log.Debug().Msg("match end - err")

	if err != nil {
		log.Warn().Msgf("Matches.err - %v", err)
		return false, err
	}

	return i.Match != nil, nil
}

func (i *ImageMatch) NotMatches() (bool, error) {
	matches, err := i.Matches()
	return !matches, err
}
