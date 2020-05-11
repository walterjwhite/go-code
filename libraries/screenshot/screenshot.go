package screenshot

import (
	"github.com/vova616/screenshot"
	"github.com/walterjwhite/go-application/libraries/logging"

	"flag"
	"image"
	"image/jpeg"
	"os"
)

var (
	jpegQualityFlag = flag.Int("ScreenshotJpegQuality", 50, "JPEG Quality (1 - 100)")
)

type Instance struct {
	filename string

	channel chan bool
	image   *image.RGBA
}

// TODO: currently only taking PNG screenshots
// support other formats
func Take(filename string) *Instance {
	i := &Instance{filename: filename}
	i.channel = make(chan bool, 1)

	img, err := screenshot.CaptureScreen()
	logging.Panic(err)

	i.image = img

	// write async as encoding is slow ...
	go i.write()

	return i
}

func (i *Instance) write() {
	file, err := os.OpenFile(i.filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	logging.Panic(err)
	defer file.Close()

	logging.Panic( /*png*/ jpeg.Encode(file, i.image, &jpeg.Options{Quality: *jpegQualityFlag}))

	i.channel <- true
}

func (i *Instance) Wait() {
	<-i.channel
	close(i.channel)
}
