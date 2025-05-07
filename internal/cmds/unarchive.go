package cmds

import (
	"cryphtron/internal/utils"
	"cryphtron/internal/utils/log"
	"errors"
	"flag"
	"os"
	"path/filepath"
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
	RemoveSrc    bool
}

func unarchiveFlags(cmdName string) *unArchiveFlags {
	result := unArchiveFlags{}
	result.FlagSet = flag.NewFlagSet(cmdName, flag.ContinueOnError)
	result.StringVar(&result.IncludeFiles, "include-files", "all", "which files to unarchive")
	result.StringVar(&result.SourceFile, "src", "", "source archive file")
	result.StringVar(&result.OutputFile, "out", "", "output path")
	result.BoolVar(&result.RemoveSrc, "remove-src", false, "remove source file after unarchiving")
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
	if len(src) == 0 {
		return errors.New("no source file specified")
	}
	srcStat, _ := os.Stat(src)
	if srcStat == nil {
		return errors.New("source file does not exist")
	}
	if srcStat.IsDir() {
		return errors.New("source file is a directory")
	}
	out := utils.NormalizePath(flag.OutputFile)
	if len(out) == 0 {
		out = filepath.Dir(src)
	}
	utils.LogCLI(log.InfoLevel).Println(cmdName, ":", "Unarchiving", src, "to", out)
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	switch flag.Name() {
	case UntarCmdName:
		err = untar(srcFile, flag)
	case UnzipCmdName:
		err = unzip(srcFile, flag)
	default:
		err = errors.New("unknown command")
	}
	if err != nil {
		return err
	}
	utils.CloseFileAndFinishBar(srcFile, nil)
	if flag.RemoveSrc {
		return os.RemoveAll(utils.NormalizePath(srcFile.Name()))
	}
	return nil
}

func unzip(src *os.File, flags *unArchiveFlags) error {
	return utils.UnZipFromFile(src, flags.OutputFile)
}

func untar(src *os.File, flags *unArchiveFlags) error {
	return utils.UnTarFromFile(src, flags.OutputFile)
}
