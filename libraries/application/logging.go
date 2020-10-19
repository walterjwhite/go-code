package application

import (
	"compress/zlib"
	"flag"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-application/libraries/application/logging"
	"io"
	"os"
)

const (
	logDateTimeFormat = "2006/01/02 15:04:05 -0700"
)

var (
	logLevel = flag.String("log-level", "error", "log level")
	logFile  = flag.String("log-file", "", "log file, if empty, stdout is used")
)

func configureLogging() {
	zerolog.TimeFieldFormat = logDateTimeFormat

	var f io.Writer = getWriter()
	log.Logger = zerolog.New(zerolog.SyncWriter(f)).With().Timestamp().Logger()

	setLogLevel()
}

func getWriter() io.Writer {
	if len(*logFile) > 0 {
		return prepareFile()
	}

	return zerolog.ConsoleWriter{Out: os.Stdout /*NoColor: false,*/, TimeFormat: logDateTimeFormat}
}

func setLogLevel() {
	if logLevel != nil {
		zlogLevel, err := zerolog.ParseLevel(*logLevel)
		logging.Panic(err)

		zerolog.SetGlobalLevel(zlogLevel)
	}
}

func prepareFile() io.WriteCloser {
	var f io.WriteCloser
	f, err := os.OpenFile(*logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	logging.Panic(err)
	defer func() {
		logging.Panic(f.Close())
	}()

	f = zlib.NewWriter(f)
	defer func() {
		logging.Panic(f.Close())
	}()

	return f
}
