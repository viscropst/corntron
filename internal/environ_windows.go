//go:build windows
// +build windows

package internal

func (c *Core) prepareEnvsByEnviron() {
	c.assignWithEnviron("PATH")
	c.assignWithEnviron("USERNAME")
	c.assignWithEnviron("APPDATA")
	c.assignWithEnviron("TEMP")
	c.assignWithEnviron("TMP")
	c.assignWithEnviron("WINDIR")
	c.assignWithEnviron("OS")
	c.assignWithEnviron("LOCALAPPDATA")
	if c.ProfileDir != "" {
		c.Env["USERPROFILE"] = c.ProfileDir
	} else {
		c.assignWithEnviron("USERPROFILE")
	}
	c.assignWithEnviron("PROGRAMW6432")
	c.assignWithEnviron("PATHEXT")
	c.assignWithEnviron("SYSTEMDRIVE")
	c.assignWithEnviron("PROGRAMDATA")
	c.assignWithEnviron("PROCESSOR_ARCHITECTURE")
	c.Env["ProgramFiles(x86)"] = c.Environ["PROGRAMFILES"] + " (x86)"
	c.Env["PSExecutionPolicyPreference"] = "RemoteSigned"
}
