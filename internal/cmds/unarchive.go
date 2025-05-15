package cmds

import (
	"cryphtron/internal/utils"
	"cryphtron/internal/utils/log"
	"errors"
	"flag"
	"os"
	"path/filepath"
	"strings"
)

const UnzipCmdId = "uzip"
const UntarCmdId = "utar"

var (
	UnzipCmdName = CmdName(UnzipCmdId)
	UntarCmdName = CmdName(UntarCmdId)
)

const UnzipIncludeAll = "all"

func init() {
	AppendCmd(UnzipCmdName, UnzipCmd)
	AppendCmd(UntarCmdName, UntarCmd)
}

type unArchiveFlags struct {
	*flag.FlagSet
	IncludeFiles string
	OutputPath   string
	SourceFile   string
	RemoveSrc    bool
	BaseDir      string
}

func unarchiveFlags(cmdName string) *unArchiveFlags {
	result := unArchiveFlags{}
	result.FlagSet = flag.NewFlagSet(cmdName, flag.ContinueOnError)
	result.StringVar(&result.IncludeFiles, "include-files", UnzipIncludeAll, "which files to unarchive")
	result.StringVar(&result.SourceFile, "src", "", "source archive file")
	result.StringVar(&result.OutputPath, "out", "", "output path")
	result.BoolVar(&result.RemoveSrc, "remove-src", false, "remove source file after unarchiving")
	result.StringVar(&result.BaseDir, "base-dir", "", "base dir in archive file")
	return &result
}

func (f *unArchiveFlags) normalizeFlags(args []string) ([]string, error) {
	err := f.Parse(args)
	if err != nil {
		return nil, err
	}
	src := utils.NormalizePath(f.SourceFile)
	if len(src) == 0 {
		return nil, errors.New("no source file specified")
	}
	srcStat, _ := os.Stat(src)
	if srcStat == nil {
		return nil, errors.New("source file does not exist")
	}
	if srcStat.IsDir() {
		return nil, errors.New("source file is a directory")
	}
	f.SourceFile = src
	out := utils.NormalizePath(f.OutputPath)
	if len(out) == 0 {
		out = filepath.Dir(src)
	}
	f.OutputPath = out
	if len(f.IncludeFiles) == 0 || f.IncludeFiles == UnzipIncludeAll {
		f.IncludeFiles = UnzipIncludeAll
		return nil, nil
	}
	f.IncludeFiles = strings.TrimSpace(f.IncludeFiles)
	return strings.Split(f.IncludeFiles, ","), nil
}

func UnzipCmd(args []string) error {
	return UnArchiveCmd(UnzipCmdName, args)
}

func UntarCmd(args []string) error {
	return UnArchiveCmd(UntarCmdName, args)
}

func UnArchiveCmd(cmdName string, args []string) error {
	flags := unarchiveFlags(cmdName)
	includes, err := flags.normalizeFlags(args)
	if err != nil {
		return err
	}
	utils.LogCLI(log.InfoLevel).Println(cmdName, ":", "Unarchiving", flags.SourceFile, "to", flags.OutputPath)
	srcFile, err := os.Open(flags.SourceFile)
	if err != nil {
		return err
	}
	err = unarchive(srcFile, flags, includes...)
	if err != nil {
		return err
	}
	utils.CloseFileAndFinishBar(srcFile, nil)
	if flags.RemoveSrc {
		return os.RemoveAll(utils.NormalizePath(srcFile.Name()))
	}
	return nil
}

func unarchive(srcFile *os.File, flags *unArchiveFlags, includes ...string) error {
	switch flags.Name() {
	case UntarCmdName:
		return utils.UnTarFromFileWithBaseDir(
			srcFile, flags.OutputPath, flags.BaseDir, includes...)
	case UnzipCmdName:
		return utils.UnZipFromFileWithBaseDir(
			srcFile, flags.OutputPath, flags.BaseDir, includes...)
	default:
		return errors.New("unknown command")
	}
}
