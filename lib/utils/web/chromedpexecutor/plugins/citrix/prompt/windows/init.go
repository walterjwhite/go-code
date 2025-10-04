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

	windows10TermsAcceptanceButtonData  []byte
	windows10TermsAcceptanceButtonImage image.Image

	windows11TermsAcceptanceButtonData  []byte
	windows11TermsAcceptanceButtonImage image.Image
)

func init() {
	windows10PngData, err := png.Decode(bytes.NewReader(windows10StartButtonData))
	logging.Panic(err)

	windows10StartButtonImage = locateimage.Convert(windows10PngData)

	windows11PngData, err := png.Decode(bytes.NewReader(windows11StartButtonData))
	logging.Panic(err)

	windows11StartButtonImage = locateimage.Convert(windows11PngData)

	windows10TermsAcceptancePngData, err := png.Decode(bytes.NewReader(windows10TermsAcceptanceButtonData))
	logging.Panic(err)

	windows10TermsAcceptanceButtonImage = locateimage.Convert(windows10TermsAcceptancePngData)

	windows11TermsAcceptancePngData, err := png.Decode(bytes.NewReader(windows11TermsAcceptanceButtonData))
	logging.Panic(err)

	windows11TermsAcceptanceButtonImage = locateimage.Convert(windows11TermsAcceptancePngData)
}
