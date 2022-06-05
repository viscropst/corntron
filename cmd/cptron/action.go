package cptron

import (
	"cryphtron"
	"cryphtron/core"
)

type CmdAction interface {
	ActionName() string
	ParseArg(info FlagInfo) error
	BeforeCore(coreConfig *core.MainConfig) error
	Exec(core *cryphtron.Core) error
}
