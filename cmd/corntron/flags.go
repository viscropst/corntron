package corntron

import (
	"corntron/internal/log"
	"errors"
	"flag"
	"os"
	"strings"
)

type FlagInfo struct {
	Name     string
	Index    int
	CmdName  string
	TotalLen int
	Args     []string
}

type CmdFlag struct {
	Host             *flag.FlagSet
	flagLen          int
	argLen           int
	osArgLen         int
	args             []string
	NoWaiting        bool
	ConfigBase       string
	RunningBase      string
	MirrorType       string
	PassToCornConfig bool
}

func (f CmdFlag) Prepare() *CmdFlag {
	result := &CmdFlag{}
	result.Host = flag.CommandLine
	return result
}

func (f *CmdFlag) Parse() (*FlagInfo, error) {
	startIdx := 1
	f.args = os.Args
	if len(f.args) > 2 {
		if len(f.args[1]) < 2 {
			startIdx += 1
		}
		if sp := strings.Split(strings.TrimSpace(f.args[startIdx]), " "); len(sp) > 1 {
			f.args = make([]string, 1)
			f.args[0] = os.Args[0]
			f.args = append(f.args, sp...)
			f.args = append(f.args, os.Args[startIdx+1:]...)
		}
	}
	err := f.Host.Parse(f.args[startIdx:])
	if err != nil && !f.PassToCornConfig {
		return nil, err
	}
	f.flagLen = f.Host.NFlag() * 2
	f.argLen = f.Host.NArg()
	f.osArgLen = len(f.args) - 1
	if (f.osArgLen-f.argLen) < 0 || (f.argLen+f.flagLen) == 0 {
		return nil, errors.New("invalid length of args,use '-help' for usage")
	}
	idxArgAct := startIdx - 1
	idxArgAct = idxArgAct + f.flagLen + 1
	if (f.osArgLen + 1) < (idxArgAct + f.argLen) {
		idxArgAct -= 1
	}
	info := FlagInfo{
		Name:     os.Args[idxArgAct],
		Index:    idxArgAct,
		CmdName:  f.Host.Name(),
		TotalLen: f.osArgLen,
		Args:     f.args,
	}
	CliLog(log.DebugLevel).Println("startIndex", startIdx)
	return &info, nil
}
