package cptron

import (
	"cryphtron/core"
	"cryphtron/internal/utils"
	"log"
)

func CliLog(levels ...utils.LogLevel) *log.Logger {
	return core.LogCLI(levels...)
}
