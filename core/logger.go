package core

import (
	stdlog "log"

	"github.com/skerkour/rz"
	"github.com/skerkour/rz/log"
)

type xLogger struct {
	logger rz.Logger
	level  rz.LogLevel
}

func (l xLogger) With(options ...rz.LoggerOption) xLogger {
	tmp := &l
	tmp.logger = tmp.logger.With(options...)
	return *tmp
}

func (l xLogger) Write(p []byte) (n int, err error) {
	n = len(p)
	if n > 0 && p[n-1] == '\n' {
		// Trim CR added by stdlog.
		p = p[0 : n-1]
	}
	l.logger.LogWithLevel(l.level, string(p))
	return
}

func rzLogger(level ...rz.LogLevel) xLogger {
	rzLog := xLogger{logger: log.Logger()}
	if len(level) > 0 {
		rzLog.level = level[0]
	} else {
		rzLog.level = rz.NoLevel
	}
	return rzLog
}

func LogCLI(level ...rz.LogLevel) *stdlog.Logger {
	_logger := stdlog.Default()
	_logger.SetOutput(
		rzLogger(level...).
			With(rz.Formatter(rz.FormatterCLI())))
	_logger.SetFlags(0)
	return _logger
}
