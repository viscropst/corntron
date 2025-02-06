package cmds

import (
	"cryphtron/internal/utils"
	"errors"
	"flag"
	"net/http"
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

	if len(args) < 1 {
		return errors.New(
			" correct usage was: " + WgetCmdID + " src [options]")
	}

	wd, _ := os.Getwd()
	dst := wd
	if len(args) > 1 {
		dst = utils.NormalizePath(args[1])
	}

	srcURL := args[0]

	client := http.DefaultClient

	if len(args) > 1 {
		tmp := wgetFlags()
		tmp.Parse(args[1:])
	}
	method := http.MethodGet
	var err error
	var req *http.Request
	req, err = http.NewRequest(method, srcURL, nil)
	if err != nil {
		return err
	}
	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		return err
	}

	return utils.IOToFile(resp.Body, dst)
}
