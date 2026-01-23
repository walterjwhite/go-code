package application

import (
	"flag"
	"io"
	"log/syslog"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

const (
	logDateTimeFormat = "2006/01/02 00:00:00 -0700"
)

var (
	logLevel  = flag.String("log-level", "info", "log level")
	logTarget = flag.String("log-target", "", "log file, if empty, stderr is used, if SYSLOG, syslog is used")
)

func configureLogging() {
	zerolog.TimeFieldFormat = logDateTimeFormat
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	log.Logger = zerolog.New(zerolog.SyncWriter(getWriter())).With().Timestamp().Logger()
	setLogLevel()
}

func getWriter() io.WriteCloser {
	if len(*logTarget) > 0 {
		if *logTarget == "SYSLOG" {
			return getSysLogger()
		}

		return getFileLogger()
	}

	return zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: logDateTimeFormat}
}

func setLogLevel() {
	zlogLevel, err := zerolog.ParseLevel(*logLevel)
	logging.Error(err)

	zerolog.SetGlobalLevel(zlogLevel)
}

func getSysLogger() io.WriteCloser {
	syslogger, err := syslog.New(syslog.LOG_KERN|syslog.LOG_EMERG|syslog.LOG_ERR|syslog.LOG_INFO|syslog.LOG_CRIT|syslog.LOG_WARNING|syslog.LOG_NOTICE|syslog.LOG_DEBUG, ApplicationName)
	logging.Error(err)

	return zerolog.ConsoleWriter{Out: syslogger, TimeFormat: logDateTimeFormat, NoColor: true}
}

func getFileLogger() io.WriteCloser {
	var f io.WriteCloser
	f, err := os.OpenFile(*logTarget, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	logging.Error(err)

	return f
}
