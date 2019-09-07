package screenshot

import (
	"github.com/vova616/screenshot"
	"image/png"
	"github.com/walterjwhite/go-application/libraries/path"
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
	png.Encode(file, img)

	log.Printf("Captured screenshot: %v / %v", label, file.Name())
}
