package logging

import (
	"compress/zlib"
	"flag"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"os"
)

var (
	logDateTimeFormat = flag.String("LogDateTimeFormat", "2006/01/02 15:04:05 Z07:00", "LogDateTimeFormat")

	logLevel    = flag.String("LogLevel", "info", "LogLevel")
	logStdOut   = flag.Bool("LogStdOut", true, "LogStdOut")
	logFile     = flag.String("LogFile", "", "LogFile")
	logCompress = flag.Bool("LogCompress", false, "LogCompress")
)

// 1. set time format
// 2. set output & format
func Configure() {
	zerolog.TimeFieldFormat = *logDateTimeFormat

	var f io.Writer = getWriter()
	log.Logger = zerolog.New(f).With().Timestamp().Logger()

	setLogLevel()
}

func getWriter() io.Writer {
	if len(*logFile) > 0 {
		return prepareFile()
	}
	if *logStdOut {
		return zerolog.ConsoleWriter{Out: os.Stdout}
	}

	return zerolog.ConsoleWriter{Out: os.Stderr}
}

func setLogLevel() {
	if logLevel != nil {
		zlogLevel, err := zerolog.ParseLevel(*logLevel)
		Panic(err)

		zerolog.SetGlobalLevel(zlogLevel)
	}
}

func prepareFile() io.WriteCloser {
	var f io.WriteCloser
	f, err := os.OpenFile(*logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	Panic(err)
	defer func() {
		Panic(f.Close())
	}()

	if *logCompress {
		f = zlib.NewWriter(f)
		defer func() {
			Panic(f.Close())
		}()
	}

	return f
}

func Panic(err error) {
	if err != nil {
		panic(err)
	}
}

/*
func Warn(err error) {
	if err != nil {
		log.Warn().Msg(err)
	}
}
*/
