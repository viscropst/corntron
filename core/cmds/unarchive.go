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

func detectCompressionFormat(filePath string, format string) (bool, error) {
	file, err := internal.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return false, err
	}

	switch format {
	case "zip":
		// ZIP format: Magic number 0x50 0x4B 0x03 0x04
		return len(buffer) > 3 && buffer[0] == 0x50 && buffer[1] == 0x4B && buffer[2] == 0x03 && buffer[3] == 0x04, nil
	case "tar":
		// TAR format: Magic number at offset 257
		return len(buffer) > 257 && string(buffer[257:262]) == "ustar", nil
	case "gz":
		// GZIP format: Magic number 0x1F 0x8B
		return len(buffer) > 1 && buffer[0] == 0x1F && buffer[1] == 0x8B, nil
	case "xz":
		// XZ format: Magic number 0xFD 0x37 0x7A 0x58 0x5A 0x00
		return len(buffer) > 5 && buffer[0] == 0xFD && buffer[1] == 0x37 && buffer[2] == 0x7A && buffer[3] == 0x58 && buffer[4] == 0x5A && buffer[5] == 0x00, nil
	case "bz2":
		// BZIP2 format: Magic number 0x42 0x5A
		return len(buffer) > 1 && buffer[0] == 0x42 && buffer[1] == 0x5A, nil
	default:
		return false, nil
	}
}

func UnArchiveCmd(cmdName string, args []string) error {
	flags := unarchiveFlags(cmdName)
	includes, err := flags.normalizeFlags(args)
	if err != nil {
		return err
	}

	// Detect file format
	var isValid bool
	switch cmdName {
	case UnzipCmdName:
		// i-uzip only supports ZIP format
		isValid, err = detectCompressionFormat(flags.SourceFile, "zip")
		if err != nil {
			return err
		}
		if !isValid {
			return internal.Error("file format does not match expected ZIP format")
		}
	case UntarCmdName:
		// i-utar supports TAR, GZ, XZ, BZ2 and nested formats (tar.gz, tar.xz, tar.bz2)
		isValid, err = detectCompressionFormat(flags.SourceFile, "tar")
		if err != nil {
			return err
		}
		if !isValid {
			// Check for compressed formats
			isGz, _ := detectCompressionFormat(flags.SourceFile, "gz")
			isXz, _ := detectCompressionFormat(flags.SourceFile, "xz")
			isBz2, _ := detectCompressionFormat(flags.SourceFile, "bz2")
			if !isGz && !isXz && !isBz2 {
				return internal.Error("file format does not match expected TAR, GZ, XZ, or BZ2 format")
			}
		}
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
