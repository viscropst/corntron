package corntron

import (
	"corntron/core"
	"corntron/internal/utils"
	"log"
)

func CliLog(levels ...utils.LogLevel) *log.Logger {
	return core.LogCLI(levels...)
}
