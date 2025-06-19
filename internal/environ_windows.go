//go:build windows
// +build windows

package internal

func (c Environ) PrepareEnvsByEnviron(profileDir ...string) {
	c.assignWithEnviron("USERNAME")
	c.assignWithEnviron("APPDATA")
	c.assignWithEnviron("TEMP")
	c.assignWithEnviron("TMP")
	c.assignWithEnviron("WINDIR")
	c.assignWithEnviron("OS")
	c.assignWithEnviron("LOCALAPPDATA")
	if len(profileDir) > 0 {
		c["USERPROFILE"] = profileDir[0]
	} else {
		c.assignWithEnviron("USERPROFILE")
	}
	c.assignWithEnviron("PROGRAMW6432")
	c.assignWithEnviron("PATHEXT")
	c.assignWithEnviron("SYSTEMDRIVE")
	c.assignWithEnviron("PROGRAMDATA")
	c.assignWithEnviron("PROCESSOR_ARCHITECTURE")
	c["ProgramFiles(x86)"] = environ["PROGRAMFILES"] + " (x86)"
	c["PSExecutionPolicyPreference"] = "RemoteSigned"
}
