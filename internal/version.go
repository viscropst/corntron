package internal

var Version = "staging"

func VersionInfo() string {
	return "version: " + Version +
		" " + PlatID("built for: ")
}

func AgentVersionInfo(suffix string) string {
	format := Version +
		"(" +
		OS() + ";" +
		ArchID("")
	if len(suffix) == 0 {
		return format + ")"
	}
	return format + ";" + suffix
}
