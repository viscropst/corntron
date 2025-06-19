package core

import (
	"corntron/internal"
	"log"
)

func LogCLI(lv ...internal.LogLevel) *log.Logger {
	return internal.LogCLI(lv...)
}
