package corntron

import (
	"corntron/internal"
	internalLog "corntron/internal/log"
	"log"
)

func CliLog(levels ...internal.LogLevel) *log.Logger {
	return internal.LogCLI(levels...)
}

var errorLogger = CliLog(internalLog.ErrorLevel)

func ErrorLog(values ...any) {
	var err error
	if ve, errOk := values[0].(error); errOk {
		err = ve
	}
	errorLogger.Println(values...)
	CliExit(err, !IsInTerminal())
}
