package windows

import (
	"bytes"

	_ "embed"

	"github.com/andreyvit/locateimage"

	"image"
	"image/png"

	"github.com/walterjwhite/go-code/lib/application/logging"
)

var (
	windows10StartButtonData  []byte
	windows10StartButtonImage image.Image

	windows11StartButtonData  []byte
	windows11StartButtonImage image.Image
)

func init() {
	windows10PngData, err := png.Decode(bytes.NewReader(windows10StartButtonData))
	logging.Panic(err)

	windows10StartButtonImage = locateimage.Convert(windows10PngData)

	windows11PngData, err := png.Decode(bytes.NewReader(windows11StartButtonData))
	logging.Panic(err)

	windows11StartButtonImage = locateimage.Convert(windows11PngData)
}
