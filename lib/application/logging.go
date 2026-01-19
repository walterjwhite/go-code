package application

import (
	"flag"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/walterjwhite/go-code/lib/application/logging"
	"github.com/walterjwhite/go-code/lib/application/logging/pubsub"
	"github.com/walterjwhite/go-code/lib/application/property"
	"github.com/walterjwhite/go-code/lib/net/google"
	"io"
	"log/syslog"
	"os"
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

	log.Logger = zerolog.New(zerolog.SyncWriter(getLogWriter())).With().Timestamp().Logger()
	setLogLevel()
}

func getLogWriter() io.Writer {
	rw := setupPubsubLogging()
	w := getWriter()

	if rw != nil {
		return zerolog.MultiLevelWriter(rw, w)
	}

	return w
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

func setupPubsubLogging() io.Writer {
	w := &pubsub.PubsubWriter{}
	property.LoadFile(ApplicationName, w)

	if len(w.TopicName) > 0 {
		conf := &google.Conf{}
		w.Init(Context, conf)
		cw := zerolog.ConsoleWriter{Out: w, TimeFormat: logDateTimeFormat, NoColor: true}
		l := &LevelWriter{Writer: cw}


		if len(w.Level) > 0 {
			level, err := zerolog.ParseLevel(w.Level)
			if err == nil {
				l.Level = level
			}
		}

		return l
	}

	return nil
}

type LevelWriter struct {
	io.Writer
	Level zerolog.Level
}

func (lw *LevelWriter) WriteLevel(l zerolog.Level, p []byte) (n int, err error) {
	if l >= lw.Level {
		return lw.Write(p)
	}

	return len(p), nil
}
