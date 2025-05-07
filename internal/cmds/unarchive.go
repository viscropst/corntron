package cmds

import (
	"cryphtron/internal/utils"
	"flag"
	"os"
)

const UnzipCmdName = "uzip"
const UntarCmdName = "utar"

func init() {
	AppendCmd(CmdName(UnzipCmdName), UnzipCmd)
	AppendCmd(CmdName(UntarCmdName), UntarCmd)
}

type unArchiveFlags struct {
	*flag.FlagSet
	IncludeFiles string
	OutputFile   string
	SourceFile   string
}

func unarchiveFlags(cmdName string) *unArchiveFlags {
	result := unArchiveFlags{}
	result.FlagSet = flag.NewFlagSet(cmdName, flag.ContinueOnError)
	result.StringVar(&result.IncludeFiles, "include-files", "all", "which files to unarchive")
	result.StringVar(&result.SourceFile, "src", "", "source archive file")
	result.StringVar(&result.OutputFile, "out", "", "output path")
	return &result
}

func UnzipCmd(args []string) error {
	return UnArchiveCmd(UnzipCmdName, args)
}

func UntarCmd(args []string) error {
	return UnArchiveCmd(UntarCmdName, args)
}

func UnArchiveCmd(cmdName string, args []string) error {
	flag := unarchiveFlags(cmdName)
	err := flag.Parse(args)
	if err != nil {
		return err
	}
	src := utils.NormalizePath(flag.SourceFile)
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	switch flag.Name() {
	case UntarCmdName:
		return untar(srcFile, flag)
	case UnzipCmdName:
		return unzip(srcFile, flag)
	}
	return nil
}

func unzip(src *os.File, flags *unArchiveFlags) error {
	return utils.UnZipFromFile(src, flags.OutputFile)
}

func untar(src *os.File, flags *unArchiveFlags) error {
	return utils.UnTarFromFile(src, flags.OutputFile)
}
