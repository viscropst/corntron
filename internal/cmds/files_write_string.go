package cmds

import (
	"cryphtron/internal/utils"
	"cryphtron/internal/utils/log"
	"errors"
	"flag"
	"io"
	"strings"
)

const WriteStringCmdId = "wstr"

var WriteStringCmdName = CmdName(WriteStringCmdId)

func init() {
	AppendCmd(CmdName(WriteStringCmdId), WriteStringCmd)
}

type writeStringFlagSet struct {
	*flag.FlagSet
	Content    string
	OutputPath string
}

func writeStringFlags() *writeStringFlagSet {
	result := writeStringFlagSet{}
	result.FlagSet = flag.NewFlagSet(WriteStringCmdName, flag.ContinueOnError)
	result.StringVar(&result.OutputPath, "out", "", "output path of target file")
	result.StringVar(&result.Content, "content", "", "content of target file")
	return &result
}

func (f *writeStringFlagSet) normalizeFlags(args []string) (io.Reader, error) {
	err := f.Parse(args)
	if err != nil {
		return nil, err
	}
	output := utils.NormalizePath(f.OutputPath)
	if len(output) == 0 {
		return nil, errors.New("no output file specified")
	}
	content := strings.TrimSpace(f.Content)
	return strings.NewReader(content), nil
}

func WriteStringCmd(args []string) error {
	flags := writeStringFlags()
	from, err := flags.normalizeFlags(args)
	if err != nil {
		return err
	}
	return createCmd(from, flags.OutputPath)
}

func createCmd(from io.Reader, output string) error {
	utils.LogCLI(log.InfoLevel).Println("Writing to target:", output)
	return utils.IOToFile(from, output, nil)
}
