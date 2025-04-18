package core

import (
	"cryphtron/internal/utils"
	"log"
)

func LogCLI(lv ...utils.LogLevel) *log.Logger {
	return utils.LogCLI(lv...)
}
