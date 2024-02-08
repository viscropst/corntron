package cptron

import (
	"cryphtron/core"
	"log"

	"github.com/rs/zerolog"
)

func CliLog(levels ...zerolog.Level) *log.Logger {
	return core.LogCLI(levels...)
}
