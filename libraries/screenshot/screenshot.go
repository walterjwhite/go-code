package screenshot

import (
	"github.com/rs/zerolog/log"
	"github.com/vova616/screenshot"
	"github.com/walterjwhite/go-application/libraries/logging"
	//"github.com/walterjwhite/go-application/libraries/path"
	// poor performance ~ 9 seconds on Linux?
	//"image/png"
	"bytes"
	"flag"
	"image/jpeg"
	"os"
)

var (
	jpegQualityFlag = flag.Int("ScreenshotJpegQuality", 90, "JPEG Quality (1 - 100)")
)

// TODO: currently only taking PNG screenshots
// support other formats

func Take(filename string) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	logging.Panic(err)
	defer file.Close()

	img, err := screenshot.CaptureScreen()
	logging.Panic(err)

	buffer := new(bytes.Buffer)
	logging.Panic( /*png*/ jpeg.Encode(buffer, img, &jpeg.Options{Quality: *jpegQualityFlag}))

	//logging.Panic(png.Encode(file, img))
	_, err = file.Write(buffer.Bytes())
	logging.Panic(err)

	log.Debug().Msgf("Captured screenshot: %v", filename)
}
