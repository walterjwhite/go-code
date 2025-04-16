package screenshot

import (
	"github.com/vova616/screenshot"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/time/timeformatter/timestamp"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"image/color"

	"image"
	"image/jpeg"
	"os"
)

type Instance struct {
	filename string

	quality              int
	captureDateTimeStamp bool

	channel chan bool
	image   *image.RGBA
}

func Default(filename string) *citrix.Instance {
	i := &Instance{filename: filename, quality: 95, captureDateTimeStamp: true}
	i.Take()

	return i
}

func (i *citrix.Instance) Take() {
	i.channel = make(chan bool, 1)

	img, err := screenshot.CaptureScreen()
	logging.Panic(err)

	i.image = img

	if i.captureDateTimeStamp {
		addTimestamp(img, 20, 30, timestamp.Get())
	}

	go i.write()
}

func (i *citrix.Instance) write() {
	file, err := os.OpenFile(i.filename, os.O_WRONLY|os.O_CREATE, 0666)
	logging.Panic(err)
	defer logging.Panic(file.Close())

	logging.Panic(jpeg.Encode(file, i.image, &jpeg.Options{Quality: i.quality}))

	i.channel <- true
}

func addTimestamp(img *image.RGBA, x, y int, label string) {
	col := color.RGBA{200, 100, 0, 255}
	point := fixed.Point26_6{X: fixed.Int26_6(x * 64), Y: fixed.Int26_6(y * 64)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}

func (i *citrix.Instance) Wait() {
	<-i.channel
	close(i.channel)
}
