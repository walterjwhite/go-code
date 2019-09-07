package screenshot

import (
	"github.com/vova616/screenshot"
	"github.com/walterjwhite/go-application/libraries/path"
	"image/png"
	"log"
)

// TODO: currently only taking PNG screenshots
// support other formats

func Take(label string, detail string) {
	img, err := screenshot.CaptureScreen()
	if err != nil {
		panic(err)
	}

	file := path.GetFile(label, "png", detail)

	defer file.Close()
	err = png.Encode(file, img)
	if err != nil {
		panic(err)
	}

	log.Printf("Captured screenshot: %v / %v", label, file.Name())
}
