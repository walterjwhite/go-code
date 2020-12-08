package weather

import (
	"bytes"
	"fmt"
	"github.com/walterjwhite/go/lib/application/logging"
	"github.com/walterjwhite/go/lib/external/spot/data"
	"net/http"
	"strconv"
)

const (
	urlTemplate = "http://wttr.in/%s,%s?0&Q&T"
)

func Get(r *data.Record) (string, string) {
	url := getWeatherUrl(r)

	req, err := http.NewRequest("GET", url, nil)
	logging.Panic(err)

	req.Header.Set("User-Agent", "curl/7.72.0")

	client := &http.Client{}
	response, err := client.Do(req)

	logging.Panic(err)

	defer response.Body.Close()

	if response.ContentLength <= 0 {
		logging.Panic(fmt.Errorf("Invalid response: %s -> %v", url, response.ContentLength))
	}

	buf := bytes.NewBuffer(make([]byte, 0, response.ContentLength))
	_, err = buf.ReadFrom(response.Body)
	logging.Panic(err)

	return url, buf.String()
}

func getWeatherUrl(r *data.Record) string {
	latitude := strconv.FormatFloat(r.Latitude, 'f', -1, 64)
	longitude := strconv.FormatFloat(r.Longitude, 'f', -1, 64)

	return fmt.Sprintf(urlTemplate, latitude, longitude)
}
