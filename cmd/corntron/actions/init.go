package actions

import "corntron/cmd/corntron"

var mActions = make([]corntron.CmdAction, 0)

func appendAction(act corntron.CmdAction) {
	mActions = append(mActions, act)
}

func ActionMap() map[string]corntron.CmdAction {
	result := make(map[string]corntron.CmdAction)
	for _, v := range mActions {
		result[v.ActionName()] = v
	}
	return result
}
