package utils

import (
	"cryphtron/internal/utils/log"
	stdlog "log"
)

type LogLevel = log.Level

func LogCLI(level ...log.Level) *stdlog.Logger {
	_logger := stdlog.Default()
	_logger.SetOutput(
		log.ZeroLogger(level...))
	_logger.SetFlags(0)
	return _logger
}
