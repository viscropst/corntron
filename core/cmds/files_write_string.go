package cmds

import (
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
	output := NormalizeFilePath(f.OutputPath)
	if len(output) == 0 {
		return nil, cmdError("no output file specified")
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
	LogInfo("Writing to target:", output)
	return ReaderToFile(from, output, nil)
}

func appendCmd(from io.Reader, output string) error {
	LogInfo("Appending to target:", output)
	_, err := StatFilePath(output)
	if err != nil {
		return err
	}
	return ReaderToFile(from, output, nil, FSAppendFlag())
}
