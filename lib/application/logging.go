package application

import (
	"compress/zlib"
	"flag"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"io"
	"log/syslog"
	"os"
)

const (
	logDateTimeFormat = "2006/01/02 00:00:00 -0700"
)

var (
	logLevel  = flag.String("log-level", "warn", "log level")
	logTarget = flag.String("log-target", "", "log file, if empty, stderr is used, if SYSLOG, syslog is used")
)

func configureLogging() {
	zerolog.TimeFieldFormat = logDateTimeFormat

	var f = getWriter()
	log.Logger = zerolog.New(zerolog.SyncWriter(f)).With().Timestamp().Logger()

	setLogLevel()
}

func getWriter() io.Writer {
	if len(*logTarget) > 0 {
		if *logTarget == "SYSLOG" {
			syslogger, err := syslog.New(syslog.LOG_KERN|syslog.LOG_EMERG|syslog.LOG_ERR|syslog.LOG_INFO|syslog.LOG_CRIT|syslog.LOG_WARNING|syslog.LOG_NOTICE|syslog.LOG_DEBUG, ApplicationName)
			logging.Panic(err)

			return zerolog.ConsoleWriter{Out: syslogger, TimeFormat: logDateTimeFormat, NoColor: true}
		}

		return prepareFile()
	}

	return zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: logDateTimeFormat}
}

func setLogLevel() {
	zlogLevel, err := zerolog.ParseLevel(*logLevel)
	logging.Panic(err)

	zerolog.SetGlobalLevel(zlogLevel)
}

func prepareFile() io.WriteCloser {
	var f io.WriteCloser
	f, err := os.OpenFile(*logTarget, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
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
