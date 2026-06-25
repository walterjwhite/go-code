package ocr

import (
	"github.com/otiai10/gosseract/v2"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

func Text(data []byte) (string, error) {
	client := gosseract.NewClient()
	defer closeResource(client)

	err := client.SetImageFromBytes(data)
	if err != nil {
		return "", err
	}

	return client.Text()
}

func closeResource(client *gosseract.Client) {
	logging.Warn(client.Close(), "close")
}
