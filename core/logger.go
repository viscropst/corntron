package core

import (
	"corntron/internal/utils"
	"log"
)

func LogCLI(lv ...utils.LogLevel) *log.Logger {
	return utils.LogCLI(lv...)
}
