package notification

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasttemplate"

	"github.com/walterjwhite/go/lib/external/spot/weather"

	"strconv"
	"strings"
)

func (c *Notification) prepareTemplateContext() {
	log.Debug().Msgf("record: %v", c.Record)

	latitude := strconv.FormatFloat(c.Record.Latitude, 'f', -1, 64)
	longitude := strconv.FormatFloat(c.Record.Longitude, 'f', -1, 64)

	weatherUrl, weatherReport := weather.Get(c.Record)

	c.Context["DateTime"] = c.Record.DateTime.Format(dateFormat)
	c.Context["Latitude"] = latitude
	c.Context["Longitude"] = longitude
	c.Context["WeatherUrl"] = weatherUrl
	c.Context["WeatherReport"] = weatherReport
	c.Context["GoogleMaps"] = fmt.Sprintf("https://maps.google.com/maps?f=q&geocode=&q=%s,%s&ll=%s,%s", latitude, longitude, latitude, longitude)
}

func (c *Notification) doTemplate() {
	t := fasttemplate.New(c.EmailMessage.Body, "{{", "}}")
	c.EmailMessage.Body = t.ExecuteString(c.Context)
	c.EmailMessage.Body = strings.ReplaceAll(c.EmailMessage.Body, "\n", "\r\n")
}
