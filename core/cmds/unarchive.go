package cmds

import (
	"corntron/internal"
	"corntron/internal/log"
	"flag"
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
	LogicPath    internal.UnarLogicPath
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
	src := internal.NormalizePath(f.SourceFile)
	if len(src) == 0 {
		return nil, internal.Error("no source file specified")
	}
	srcStat, _ := internal.StatPath(src)
	if srcStat == nil {
		return nil, internal.Error("source file does not exist")
	}
	if srcStat.IsDir() {
		return nil, internal.Error("source file is a directory")
	}
	f.SourceFile = src
	out := internal.NormalizePath(f.OutputPath)
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

func (f unArchiveFlags) UnarFlags() internal.UnarchiveFlag {
	result := internal.UnarchiveFlag{
		SourceFile: f.SourceFile,
		OutputPath: f.OutputPath,
		BaseDir:    f.BaseDir,
	}
	switch f.Name() {
	case UntarCmdName:
		result.LogicType = internal.UnTar
	case UnzipCmdName:
		result.LogicType = internal.UnZip
	}
	return result
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
	internal.LogCLI(log.InfoLevel).Println(cmdName, ":", "Unarchiving", flags.SourceFile, "to", flags.OutputPath)
	srcFile, err := internal.Open(flags.SourceFile)
	if err != nil {
		return err
	}
	err = internal.Unarchive(srcFile, flags.UnarFlags(), includes...)
	if err != nil {
		return err
	}
	internal.CloseFileAndFinishBar(srcFile, nil)
	if flags.RemoveSrc {
		return internal.Remove(internal.NormalizePath(srcFile.Name()))
	}
	return nil
}
