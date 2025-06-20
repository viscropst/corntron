package cmds

import (
	"corntron/internal"
	"corntron/internal/log"
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
	IsAppend   bool
	IsNewLine  bool
}

func writeStringFlags() *writeStringFlagSet {
	result := writeStringFlagSet{}
	result.FlagSet = flag.NewFlagSet(WriteStringCmdName, flag.ContinueOnError)
	result.StringVar(&result.OutputPath, "out", "", "output path of target file")
	result.StringVar(&result.Content, "content", "", "content of target file")
	result.BoolVar(&result.IsAppend, "append", false, "writing file by appending")
	result.BoolVar(&result.IsNewLine, "newline", false, "writing file by appending a line")
	return &result
}

func (f *writeStringFlagSet) normalizeFlags(args []string) (io.Reader, error) {
	err := f.Parse(args)
	if err != nil {
		return nil, err
	}
	output := internal.NormalizePath(f.OutputPath)
	if len(output) == 0 {
		return nil, errors.New("no output file specified")
	}
	f.OutputPath = output
	content := strings.TrimSpace(f.Content)
	if f.IsNewLine {
		content = "\n" + content
	}
	return strings.NewReader(content), nil
}

func WriteStringCmd(args []string) error {
	flags := writeStringFlags()
	from, err := flags.normalizeFlags(args)
	if err != nil {
		return err
	}
	if flags.IsAppend {
		return appendCmd(from, flags.OutputPath)
	}
	return createCmd(from, flags.OutputPath)
}

func createCmd(from io.Reader, output string) error {
	internal.LogCLI(log.InfoLevel).Println("Writing to target:", output)
	return internal.IOToFile(from, output, nil)
}

func appendCmd(from io.Reader, output string) error {
	internal.LogCLI(log.InfoLevel).Println("Appending to target:", output)
	_, err := internal.StatPath(output)
	if err != nil {
		return err
	}
	return internal.IOToFile(from, output, nil, internal.FSAppendFlag())
}
