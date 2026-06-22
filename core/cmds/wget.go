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
	result.Usage = func() {
		LogInfo(WgetCmdName, ": Downloads a File")
		LogInfo("usage was: " + WgetCmdID + " src [dst]")
		LogInfo("dst is optional, default is current working directory,the file name is the same as the last part of the src url")
	}
	return result
}

func WgetCmd(args []string) error {
	srcURL := args[0]
	flags := wgetFlags()
	if len(args) >= 1 {
		flags.Parse(args[1:])
	}
	if isNeedHelpAtFirst(args) {
		flags.Usage()
		return nil
	}
	dst := GetWorkDir()
	if len(args) > 1 {
		dst = NormalizeFilePath(args[1])
	} else {
		dst = JoinFilePath(dst, BaseFileName(srcURL))
	}

	LogInfo(WgetCmdName+":", "Downloading", dst, "from", srcURL)
	return HttpRequestFileWithAgentSuffix(srcURL, AgentName(WgetCmdID), dst)
}
