package download

import (
	"github.com/walterjwhite/go-application/libraries/application/logging"
	"io"
	"net/http"
	"os"
)

type HttpDownload struct {
	LocalFilepath string
	Url           string
}

func (h *HttpDownload) Fetch() {
	response, err := http.Get(h.Url)
	logging.Panic(err)

	defer response.Body.Close()

	out, err := os.Create(h.LocalFilepath)
	logging.Panic(err)

	defer out.Close()

	_, err = io.Copy(out, response.Body)
	logging.Panic(err)
}
