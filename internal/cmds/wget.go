package cmds

import (
	"errors"
	"flag"
	"net/http"
	"os"
)

const WgetCmdId = "wgt"

var WgetCmdName = CmdName(WgetCmdId)

func init() {
	AppendCmd(CmdName(WgetCmdId), WgetCmd)
}

func wgetFlags() *flag.FlagSet {
	result := flag.NewFlagSet(WgetCmdName, flag.ContinueOnError)
	return result
}

func WgetCmd(args []string) error {

	if len(args) < 1 {
		return errors.New(
			" correct usage was: " + WgetCmdId + " src [options]")
	}

	srcUrl := args[0]

	client := http.DefaultClient

	if len(args) > 1 {
		tmp := wgetFlags()
		tmp.Parse(args[1:])
	}
	method := http.MethodGet
	var err error
	var req *http.Request
	req, err = http.NewRequest(method, srcUrl, nil)
	if err != nil {
		return err
	}
	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		return err
	}

	wd, _ := os.Getwd()
	return ioToFile(resp.Body, wd)
}
