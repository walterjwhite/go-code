package screenshot

import (
	"github.com/rs/zerolog/log"
	"github.com/vova616/screenshot"
	"github.com/walterjwhite/go-application/libraries/logging"
	"github.com/walterjwhite/go-application/libraries/path"
	// poor performance ~ 9 seconds on Linux?
	//"image/png"
	"bytes"
	"flag"
	"image/jpeg"
)

var (
	jpegQualityFlag = flag.Int("ScreenshotJpegQuality", 90, "JPEG Quality (1 - 100)")
)

// TODO: currently only taking PNG screenshots
// support other formats

func Take(label string, detail string) {
	img, err := screenshot.CaptureScreen()
	logging.Panic(err)

	file := path.GetFile(label, "jpg", detail)

	defer file.Close()

	buffer := new(bytes.Buffer)
	logging.Panic( /*png*/ jpeg.Encode(buffer, img, &jpeg.Options{Quality: *jpegQualityFlag}))

	//logging.Panic(png.Encode(file, img))
	_, err = file.Write(buffer.Bytes())
	logging.Panic(err)

	log.Debug().Msgf("Captured screenshot: %v / %v", label, file.Name())
}
