package windows

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/walterjwhite/go-code/lib/utils/ui/graphical"
	"image"
	"time"
)

const (
	QueryTimeout = 5 * time.Second
)

func (a *Application) IsUp(pctx context.Context) (bool, *graphical.ImageMatch, error) {
	log.Debug().Msgf("application.%s.IsDown", a)
	return a.isApplicationIconPresent(pctx, a.TaskBarIconUpImage)
}

func (a *Application) IsDown(pctx context.Context) (bool, *graphical.ImageMatch, error) {
	log.Debug().Msgf("application.%s.IsDown", a)
	return a.isApplicationIconPresent(pctx, a.TaskBarIconDownImage)
}

func (a *Application) isApplicationIconPresent(pctx context.Context, image image.Image) (bool, *graphical.ImageMatch, error) {
	width, _, err := a.WindowsConf.Controller.GetScreenSize(pctx)
	if err != nil {
		return false, nil, err
	}

	i := &graphical.ImageMatch{Ctx: pctx, Image: image, MatchRegion: &graphical.MatchRegion{X: 0, Y: a.WindowsConf.StartButtonHeight,
		Width: float64(width), Height: a.WindowsConf.StartButtonHeight}, Controller: a.WindowsConf.Controller}

	matched, err := i.WaitUntilMatched(pctx, QueryTimeout)
	return matched, i, err
}

func (a *Application) IsRunning(pctx context.Context) (bool, error) {
	log.Debug().Msgf("application.%s.IsRunning", a)

	matches, _, err := a.IsUp(pctx)
	if err != nil {
		return false, err
	}

	if matches {
		return true, nil
	}

	matches, _, err = a.IsDown(pctx)
	return matches, err
}

func (a *Application) SwitchTo(pctx context.Context) (bool, error) {
	log.Debug().Msgf("application.%s.SwitchTo", a)
	matches, _, err := a.IsUp(pctx)
	if err != nil {
		return false, err
	}

	if matches {
		return true, nil
	}

	var i *graphical.ImageMatch
	matches, i, err = a.IsDown(pctx)
	if err != nil {
		return false, err
	}

	if matches {
		offset := a.WindowsConf.StartButtonHeight / 2.0

		ctx, cancel := context.WithTimeout(pctx, QueryTimeout)
		defer cancel()

		return true, a.WindowsConf.Controller.Click(ctx, float64(i.Match.Rect.Min.X)+offset, float64(i.Match.Rect.Min.Y)+offset)
	}

	return false, err
}

func (a *Application) Maximize(pctx context.Context) (bool, error) {
	log.Debug().Msgf("application.%s.Maximize", a)
	success, err := a.SwitchTo(pctx)
	if err != nil {
		return false, err
	}



	return success, nil
}
