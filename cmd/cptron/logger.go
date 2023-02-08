package cptron

import (
	"cryphtron/core"
	"log"

	"github.com/skerkour/rz"
)

func CliLog(levels ...rz.LogLevel) *log.Logger {
	return core.LogCLI(levels...)
}
