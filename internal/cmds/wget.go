package cmds

import (
	"cryphtron/internal/utils"
	"cryphtron/internal/utils/log"
	"errors"
	"flag"
	"os"
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
	wd, _ := os.Getwd()
	dst := wd
	if len(args) > 1 {
		dst = utils.NormalizePath(args[1])
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
	utils.LogCLI(log.InfoLevel).Println(WgetCmdName+":", "Downloading", dst, "from", srcURL)
	return utils.HttpRequestFile(srcURL, dst)
}
