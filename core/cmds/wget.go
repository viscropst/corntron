package cmds

import (
	"flag"
)

const WgetCmdID = "wgt"

var WgetCmdName = CmdName(WgetCmdID)

func init() {
	AppendCmd(CmdName(WgetCmdID), WgetCmd)
}

func wgetFlags() *flag.FlagSet {
	result := flag.NewFlagSet(WgetCmdName, flag.ContinueOnError)
	return result
}

func WgetCmd(args []string) error {
	dst := GetWorkDir()
	if len(args) > 1 {
		dst = NormalizeFilePath(args[1])
	}
	srcURL := args[0]
	if len(args) > 1 {
		tmp := wgetFlags()
		tmp.Parse(args[1:])
	}
	if len(args) < 1 {
		return cmdError(
			" correct usage was: " + WgetCmdID + " src [options]")
	}
	LogInfo(WgetCmdName+":", "Downloading", dst, "from", srcURL)
	return HttpRequestFileWithAgentSuffix(srcURL, AgentName(WgetCmdID), dst)
}
