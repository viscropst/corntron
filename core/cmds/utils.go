package cmds

import (
	"corntron/internal"
	"corntron/internal/log"
	"io"
)

type cmdError string

func (s cmdError) Error() string {
	return string(s)
}

func cmdErrors(err ...string) cmdError {
	msg := ""
	for _, e := range err {
		msg += e
	}
	return cmdError(msg)
}

func isNeedHelpAtFirst(args []string) bool {
	if len(args) == 0 {
		return true
	}
	firstArg := args[0]
	if firstArg == "-h" || firstArg == "--help" || firstArg == "-help" {
		return true
	}
	return false
}

var LogDebug = func(v ...any) {
	internal.LogCLI(log.DebugLevel).Println(v...)
}

var LogInfo = func(v ...any) {
	internal.LogCLI(log.InfoLevel).Println(v...)
}

var LogError = func(v ...any) {
	internal.LogCLI(log.ErrorLevel).Println(v...)
}

var FSAppendFlag = func() int {
	return internal.FSAppendFlag()
}

var MkDir = func(path string) error {
	return internal.Mkdir(path)
}

var GetWorkDir = func(alt ...string) string {
	return internal.GetWorkDir(alt...)
}

var NormalizeFilePath = func(p string) string {
	return internal.NormalizePath(p)
}

var RemoveFileAndFolders = func(p string) error {
	return internal.Remove(p)
}

var CopyFile = func(from internal.FileInfo, to string, excludes ...string) error {
	return internal.CopyToFile(from, to, excludes...)
}

var OpenFile = func(src string) (internal.File, error) {
	return internal.Open(src)
}

var CloseFileAndFinishBar = func(file io.Closer, bar internal.ProgressBar) {
	internal.CloseFileAndFinishBar(file, bar)
}

var ReaderToFile = func(from io.Reader, to string, bar internal.ProgressBar, flags ...int) error {
	return internal.IOToFile(from, to, bar, flags...)
}

var StatFilePath = func(path string) (*internal.FileStatInfo, error) {
	return internal.StatPath(path)
}

type UnArchiveLogicPath = internal.UnarLogicPath
type UnArchiveFlag = internal.UnarchiveFlag

const (
	UnArchiveTar = internal.UnTar
	UnArchiveZip = internal.UnZip
)

var UnArchiveFile = func(src internal.File, flags UnArchiveFlag, includes ...string) error {
	return internal.Unarchive(src, flags, includes...)
}

var HttpRequestBytesWithAgentSuffix = func(url string, agentSuffix string, others ...string) ([]byte, error) {
	return internal.HttpRequestBytesWithAgentSuffix(url, agentSuffix, others...)
}

var HttpRequestFileWithAgentSuffix = func(url, agentSuffix, filename string, others ...string) error {
	return internal.HttpRequestFileWithAgentSuffix(url, agentSuffix, filename, others...)
}

var JoinFilePath = func(p ...string) string { return internal.JoinFilePath(p...) }

var BaseFileName = func(src string) string { return internal.BaseFileName(src) }
