package graphical

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"

	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"

	"image"
	"image/png"
	"os"
	"path/filepath"
)

func (i *ImageMatch) debug(image image.Image) {
	log.Debug().Msg("match end - no matches")

	if log.Debug().Enabled() && os.Getenv("IMAGE_MATCH_DEBUG_WRITE") == "1" {
		logging.Warn(i.write(image, "screenshot"), "Matches.debug.write - screenshot")
		logging.Warn(i.write(i.Image, "search"), "Matches.debug.write - image")
	}
}

func secureTempName(prefix string) (string, error) {
	bytes := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return fmt.Sprintf("%s-%s.png", prefix, hex.EncodeToString(bytes)), nil
}

func (i *ImageMatch) write(image image.Image, fileNamePrefix string) error {
	if image == nil {
		return errors.New("write image failed: image is nil")
	}

	fileName, err := secureTempName(fileNamePrefix)
	if err != nil {
		return err
	}

	tempDir := os.TempDir()
	tempPath := filepath.Join(tempDir, fileName)

	tempFile, err := os.OpenFile(tempPath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0600)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}

	defer func() {
		_ = tempFile.Close()
		_ = os.Remove(tempPath)
	}()

	var imgBytes []byte
	imgBytes, err = ImageToBytes(image)
	if err != nil {
		return err
	}

	_, err = tempFile.Write(imgBytes)
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
