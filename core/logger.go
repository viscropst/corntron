package core

import (
	stdlog "log"

	"github.com/rs/zerolog"
)

type xLogger struct {
	logger zerolog.Logger
	level  zerolog.Level
}

func (l xLogger) With(options ...zerolog.Context) xLogger {
	tmp := &l
	tmp.logger = options[0].Logger()
	return *tmp
}

func (l xLogger) Write(p []byte) (n int, err error) {
	n = len(p)
	if n > 0 && p[n-1] == '\n' {
		// Trim CR added by stdlog.
		p = p[0 : n-1]
	}
	l.logger.WithLevel(l.level).Msg(string(p))
	return
}

func zeroLogger(level ...zerolog.Level) xLogger {
	zLog := xLogger{}
	if len(level) > 0 {
		zLog.level = level[0]
	} else {
		zLog.level = zerolog.NoLevel
	}
	cw := zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.PartsExclude = []string{zerolog.TimestampFieldName}
		if len(level) == 0 {
			w.PartsExclude = append(w.PartsExclude, zerolog.LevelFieldName)
		}
	})
	zerolog.FormattedLevels[zerolog.ErrorLevel] = "error:"
	zerolog.FormattedLevels[zerolog.WarnLevel] = "WARN:>"
	zerolog.FormattedLevels[zerolog.FatalLevel] = "fatal:"
	zerolog.FormattedLevels[zerolog.DebugLevel] = "debug:"
	zLog.logger = zerolog.New(cw)
	return zLog
}

func LogCLI(level ...zerolog.Level) *stdlog.Logger {
	_logger := stdlog.Default()
	_logger.SetOutput(
		zeroLogger(level...))
	_logger.SetFlags(0)
	return _logger
}
