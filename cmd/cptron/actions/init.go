package actions

import "cryphtron/cmd/cptron"

var mActions = make([]cptron.CmdAction, 0)

func appendAction(act cptron.CmdAction) {
	mActions = append(mActions, act)
}

func ActionMap() map[string]cptron.CmdAction {
	result := make(map[string]cptron.CmdAction)
	for _, v := range mActions {
		result[v.ActionName()] = v
	}
	return result
}
