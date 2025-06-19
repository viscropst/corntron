package internal

import (
	"corntron/internal/log"
	stdlog "log"
)

type LogLevel = log.Level

func LogCLI(level ...log.Level) *stdlog.Logger {
	_logger := stdlog.New(log.ZeroLogger(level...), "", 0)
	return _logger
}
