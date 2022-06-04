package cptron

import "cryphtron"

type CmdAction interface {
	ActionName() string
	ParseArg(info FlagInfo) error
	Exec(core *cryphtron.Core) error
}
