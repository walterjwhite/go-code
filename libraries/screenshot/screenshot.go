package screenshot

import (
	"github.com/rs/zerolog/log"
	"github.com/vova616/screenshot"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/path"
	"image/png"
)

// TODO: currently only taking PNG screenshots
// support other formats

func Take(label string, detail string) {
	img, err := screenshot.CaptureScreen()
	logging.Panic(err)

	file := path.GetFile(label, "png", detail)

	defer file.Close()

	logging.Panic(png.Encode(file, img))

	log.Debug().Msgf("Captured screenshot: %v / %v", label, file.Name())
}
