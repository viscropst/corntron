package corntron

import (
	"corntron/core"
	"corntron/internal"
	"log"
)

func CliLog(levels ...internal.LogLevel) *log.Logger {
	return core.LogCLI(levels...)
}
