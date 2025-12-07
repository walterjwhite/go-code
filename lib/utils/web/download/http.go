package download

import (
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"os"
	"time"
)

type HttpDownload struct {
	LocalFilepath string
	Url           string
}

func (h *HttpDownload) Fetch() {
	client := &http.Client{Timeout: 30 * time.Second}

	response, err := client.Get(h.Url)
	if err != nil {
		log.Error().Err(err).Msg("http fetch failed")
		return
	}
	defer func() {
		if cerr := response.Body.Close(); cerr != nil {
			log.Error().Err(cerr).Msg("response.Body.Close failed")
		}
	}()

	out, err := os.Create(h.LocalFilepath)
	if err != nil {
		log.Error().Err(err).Msg("failed to create local file for download")
		return
	}
	defer func() {
		if cerr := out.Close(); cerr != nil {
			log.Error().Err(cerr).Msg("out.Close failed")
		}
	}()

	if _, err = io.Copy(out, response.Body); err != nil {
		log.Error().Err(err).Msg("failed to copy response body to file")
		return
	}
}
