package windows

import (
	"bytes"
	"fmt"

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

	initErr error
)

func init() {
	var err error
	var pngData image.Image

	pngData, err = png.Decode(bytes.NewReader(windows10StartButtonData))
	if err != nil {
		initErr = fmt.Errorf("failed to decode windows 10 start button image: %w", err)
		logging.Error(initErr)
		return
	}
	windows10StartButtonImage = locateimage.Convert(pngData)

	pngData, err = png.Decode(bytes.NewReader(windows11StartButtonData))
	if err != nil {
		initErr = fmt.Errorf("failed to decode windows 11 start button image: %w", err)
		logging.Error(initErr)
		return
	}
	windows11StartButtonImage = locateimage.Convert(pngData)

	pngData, err = png.Decode(bytes.NewReader(windows10TermsAcceptanceButtonData))
	if err != nil {
		initErr = fmt.Errorf("failed to decode windows 10 terms acceptance button image: %w", err)
		logging.Error(initErr)
		return
	}
	windows10TermsAcceptanceButtonImage = locateimage.Convert(pngData)

	pngData, err = png.Decode(bytes.NewReader(windows11TermsAcceptanceButtonData))
	if err != nil {
		initErr = fmt.Errorf("failed to decode windows 11 terms acceptance button image: %w", err)
		logging.Error(initErr)
		return
	}
	windows11TermsAcceptanceButtonImage = locateimage.Convert(pngData)
}

func GetInitError() error {
	return initErr
}
