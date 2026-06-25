package agent

import (
	"bytes"

	_ "embed"

	"github.com/andreyvit/locateimage"

	"image"
	"image/png"

	"github.com/walterjwhite/go-code/lib/application/logging"
)

var (
	microsoft2FAHeaderData  []byte
	microsoft2FAHeaderImage image.Image

	edgeIconData  []byte
	edgeIconImage image.Image





)

func init() {
	microsoft2FAHeaderPNGData, err := png.Decode(bytes.NewReader(microsoft2FAHeaderData))
	logging.Error(err)

	microsoft2FAHeaderImage = locateimage.Convert(microsoft2FAHeaderPNGData)

	edgeIconPNGData, err := png.Decode(bytes.NewReader(edgeIconData))
	logging.Error(err)

	edgeIconImage = locateimage.Convert(edgeIconPNGData)








}
