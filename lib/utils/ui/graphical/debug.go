package graphical

import (
	"bytes"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"

	"image"
	"image/png"
	"os"
)

func (i *ImageMatch) debug(image image.Image) {
	log.Debug().Msg("match end - no matches")

	if log.Debug().Enabled() {
		logging.Warn(i.write(image, "screenshot-*.png"), false, "Matches.debug.write - screenshot")
		logging.Warn(i.write(i.Image, "search-*.png"), false, "Matches.debug.write - image")
	}
}

func (i *ImageMatch) write(image image.Image, fileNameTemplate string) error {
	tempFile, err := os.CreateTemp("", fileNameTemplate)
	if err != nil {
		return err
	}

	var bytes []byte
	bytes, err = ImageToBytes(image)
	if err != nil {
		return err
	}

	err = os.WriteFile(tempFile.Name(), bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func ImageToBytes(img image.Image) ([]byte, error) {
	buf := new(bytes.Buffer)

	err := png.Encode(buf, img)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func BytesToImage(imgBytes []byte) (image.Image, error) {
	reader := bytes.NewReader(imgBytes)

	img, err := png.Decode(reader)
	if err != nil {
		return nil, err
	}

	return img, nil
}
