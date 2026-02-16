package main

import (
	"bytes"
	"context"
	_ "embed"
	"github.com/andreyvit/locateimage"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"image/png"
	"time"
)

var (
	screenshotData []byte

	screenshotCroppedData []byte

	windowsIconData []byte
)

func main() {
	windowsIconImage, err := png.Decode(bytes.NewReader(windowsIconData))
	logging.Error(err)

	screenshotImage, err := png.Decode(bytes.NewReader(screenshotData))
	logging.Error(err)

	start := time.Now()
	m, err := locateimage.Find(context.Background(), locateimage.Convert(screenshotImage), locateimage.Convert(windowsIconImage), 0, locateimage.Fastest)
	logging.Error(err)

	runtime := time.Since(start)
	log.Info().Msgf("sample found at %v, similarity = %.*f%%, %v", m.Rect, locateimage.SimilarityDigits-2, 100*m.Similarity, runtime)

	screenshotCroppedImage, err := png.Decode(bytes.NewReader(screenshotCroppedData))
	logging.Error(err)

	start = time.Now()
	m, err = locateimage.Find(context.Background(), locateimage.Convert(screenshotCroppedImage), locateimage.Convert(windowsIconImage), 0, locateimage.Fastest)
	logging.Error(err)

	runtime = time.Since(start)
	log.Info().Msgf("sample found at %v, similarity = %.*f%%, %v", m.Rect, locateimage.SimilarityDigits-2, 100*m.Similarity, runtime)
}
