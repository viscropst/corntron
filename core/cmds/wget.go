package cmds

import (
	"corntron/internal"
	"corntron/internal/log"
	"errors"
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
	dst := internal.GetWorkDir()
	if len(args) > 1 {
		dst = internal.NormalizePath(args[1])
	}
	srcURL := args[0]
	if len(args) > 1 {
		tmp := wgetFlags()
		tmp.Parse(args[1:])
	}
	if len(args) < 1 {
		return errors.New(
			" correct usage was: " + WgetCmdID + " src [options]")
	}
	internal.LogCLI(log.InfoLevel).Println(WgetCmdName+":", "Downloading", dst, "from", srcURL)
	return internal.HttpRequestFile(srcURL, dst)
}
