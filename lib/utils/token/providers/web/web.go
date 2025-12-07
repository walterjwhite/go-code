package web

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/security/random"
	"math/rand"
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

func NewRandomWebReader() (*WebReader, error) {
	path, err := random.String(4 + rand.Intn(4))
	if err != nil {
		return nil, err
	}

	return &WebReader{Port: 1024 + rand.Intn(64512),
		Path:    "/" + path + "/",
		channel: make(chan string, 1)}, nil
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
	_, err := fmt.Fprintf(writer, "No data received: %v", request.URL.Path)
	logging.Panic(err)
}

func (r *WebReader) handleData(writer http.ResponseWriter, request *http.Request, j int, k int) {
	data := request.URL.Path[j:]
	log.Info().Msgf("Received: %v / %v", data, request.URL.Path)
	r.channel <- data
}

func (r *WebReader) handleUnknown(writer http.ResponseWriter, request *http.Request) {
	log.Debug().Msgf("Ignored %v", request.URL.Path)
	_, err := fmt.Fprintf(writer, "Ignored %v\n", request.URL.Path)
	logging.Panic(err)
}
