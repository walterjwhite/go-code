package web

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go/lib/application/logging"
	"github.com/walterjwhite/go/lib/security/random"
	"net/http"
	"strings"
)

/* Attempts to read the specified file, periodically polling <Interval> and waiting at most <Timeout>*/
type WebReader struct {
	Port int
	Path string

	Context context.Context
	Cancel  func()

	channel chan string
}

func NewRandomWebReader() *WebReader {
	return &WebReader{Port: 1024 + random.Of(64512), Path: "/" + random.String(4+random.Of(4)) + "/", channel: make(chan string, 1)}
}

func (r *WebReader) Get() string {
	log.Info().Msgf("Please enter token into: http://localhost:%v%v", r.Port, r.Path)

	go r.serve()

	log.Debug().Msg("Waiting for data")
	return <-r.channel
}

func (r *WebReader) serve() {
	logging.Panic(http.ListenAndServe(fmt.Sprintf(":%v", r.Port), r))
}

func (r *WebReader) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if r.isValid(request) {
		r.handleValid(writer, request)
		return
	} else {
		r.handleUnknown(writer, request)
	}
}

func (r *WebReader) isValid(request *http.Request) bool {
	return strings.Index(request.URL.Path, r.Path) == 0
}

func (r *WebReader) handleValid(writer http.ResponseWriter, request *http.Request) {
	j := len(r.Path)
	k := len(request.URL.Path)

	if k > j {
		r.handleData(writer, request, j, k)
		return
	}

	log.Debug().Msgf("No data received: %v", request.URL.Path)
	fmt.Fprintf(writer, "No data received: %v", request.URL.Path)
}

func (r *WebReader) handleData(writer http.ResponseWriter, request *http.Request, j int, k int) {
	data := request.URL.Path[j:]
	log.Info().Msgf("Received: %v / %v", data, request.URL.Path)
	r.channel <- data
}

func (r *WebReader) handleUnknown(writer http.ResponseWriter, request *http.Request) {
	log.Debug().Msgf("Ignored %v", request.URL.Path)
	fmt.Fprintf(writer, "Ignored %v\n", request.URL.Path)
}
