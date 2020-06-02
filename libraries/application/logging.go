package application

import (
	"compress/zlib"
	"flag"
	//"fmt"
	"github.com/rs/zerolog"
	//"github.com/rs/zerolog/diode"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	//"time"
	"github.com/walterjwhite/go-application/libraries/logging"
)

// TODO: support properties?
var (
	logDateTimeFormat = flag.String("LogDateTimeFormat", "2006/01/02 15:04:05 -0700", "LogDateTimeFormat")
	logNoColor        = flag.Bool("LogNoColor", false, "LogNoColor")

	logLevel    = flag.String("LogLevel", "error", "LogLevel")
	logStdOut   = flag.Bool("LogStdOut", true, "LogStdOut")
	logFile     = flag.String("LogFile", "", "LogFile")
	logCompress = flag.Bool("LogCompress", false, "LogCompress")
)

// 1. set time format
// 2. set output & format
func configureLogging() {
	zerolog.TimeFieldFormat = *logDateTimeFormat

	var f io.Writer = getWriter()
	//log.Logger = zerolog.New(diode.NewWriter(f, 1000, 10*time.Millisecond, onMissedMessages)).With().Timestamp().Logger()
	log.Logger = zerolog.New(zerolog.SyncWriter(f)).With().Timestamp().Logger()

	setLogLevel()
}

func getWriter() io.Writer {
	if len(*logFile) > 0 {
		return prepareFile()
	}

	return zerolog.ConsoleWriter{Out: getOutputStream(), NoColor: *logNoColor, TimeFormat: *logDateTimeFormat}
}

/*
func onMissedMessages(missed int) {
	fmt.Printf("Logger Dropped %d messages", missed)
}
*/

func getOutputStream() *os.File {
	if *logStdOut {
		return os.Stdout
	}

	return os.Stderr
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

	if *logCompress {
		f = zlib.NewWriter(f)
		defer func() {
			logging.Panic(f.Close())
		}()
	}

	return f
}
