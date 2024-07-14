package core

import (
	"cryphtron/internal/utils"
	"log"

	"github.com/rs/zerolog"
)

func LogCLI(lv ...zerolog.Level) *log.Logger {
	return utils.LogCLI(lv...)
}
