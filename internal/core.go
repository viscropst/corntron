package internal

import (
	"os"
	"runtime"
	"strings"
)

type Core struct {
	*ValueScope
	Environ    map[string]string
	ProfileDir string
}

func (c *Core) fillEnviron() {
	if c.Environ == nil {
		c.Environ = make(map[string]string)
	}
	for _, s := range os.Environ() {
		pairs := strings.SplitN(s, "=", 2)
		if pairs[1] == "" {
			continue
		}
		key := ""
		switch runtime.GOOS {
		case "windows":
			key = strings.ToUpper(pairs[0])
		default:
			key = pairs[0]
		}
		c.Environ[key] = pairs[1]
	}
}

const PathPlaceHolder = "+{PATH}"

func (c *Core) assignWithEnviron(key string) {
	if v, ok := c.Environ[key]; key != "" && ok {
		if key == "PATH" {
			c.Env[key] = PathPlaceHolder
			return
		}
		c.Env[key] = v
	}
}

func (c *Core) assertWithEnviron(args ...string) bool {
	if len(args) == 0 {
		return false
	}
	key := args[0]
	v, ok := c.Environ[key]
	if len(args) == 1 {
		return ok
	} else {
		value := args[1]
		return ok && v == value
	}
}

func (c *Core) Prepare() {
	if c.Environ != nil {
		return
	}
	c.fillEnviron()
	c.assignWithEnviron("PATH")

	switch runtime.GOOS {
	case "windows":
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
	case "linux", "freebsd", "openbsd", "macos", "darwin", "ios", "android":
		c.assignWithEnviron("SSH_AUTH_SOCK")
		c.assignWithEnviron("SSH_ASKPASS")
		c.assignWithEnviron("USER")
		c.assignWithEnviron("PWD")
		c.assignWithEnviron("LANG")
		c.assignWithEnviron("TMPDIR")
		c.assignWithEnviron("TERM")
		c.assignWithEnviron("DBUS_SESSION_BUS_ADDRESS")
		if c.ProfileDir != "" {
			c.Env["HOME"] = c.ProfileDir
		} else {
			c.assignWithEnviron("HOME")
		}
		if c.assertWithEnviron("DESKTOP_SESSION") {
			c.unixWithDesktop()
		}

	default:
	}

}

func (c *Core) unixWithDesktop() {
	c.assignWithEnviron("DISPLAY")
	c.assignWithEnviron("SESSION_MANAGER")
	c.assignWithEnviron("XDG_DATA_DIRS")
	c.assignWithEnviron("XDG_CONFIG_DIRS")
	c.assignWithEnviron("XDG_CONFIG_HOME")
	c.assignWithEnviron("XDG_RUNTIME_DIR")
	c.assignWithEnviron("XDG_CACHE_HOME")
	c.assignWithEnviron("XDG_SESSION_TYPE")
	c.assignWithEnviron("XDG_SESSION_PATH")
	c.assignWithEnviron("XDG_SESSION_CLASS")
	c.assignWithEnviron("XDG_SEAT_PATH")
	c.assignWithEnviron("XDG_CURRENT_DESKTOP")
	c.assignWithEnviron("XDG_SESSION_DESKTOP")
	c.assignWithEnviron("ICEAUTHORITY")
	c.assignWithEnviron("XAUTHORITY")
	c.assignWithEnviron("GTK_RC_FILES")
	c.assignWithEnviron("GTK2_RC_FILES")
	if c.assertWithEnviron("DESKTOP_SESSION", "plasma") {
		c.assignWithEnviron("PAM_KWALLET5_LOGIN")
		c.assignWithEnviron("QT_WAYLAND_DECORATIONS")
	}

	if c.assertWithEnviron("XDG_SESSION_TYPE", "wayland") {
		c.assignWithEnviron("WAYLAND_DISPLAY")
		c.assignWithEnviron("QT_WAYLAND_RECONNECT")
	}

	if c.assertWithEnviron("XDG_SESSION_TYPE", "x11") {
		c.assignWithEnviron("INPUT_METHOD")
		c.assignWithEnviron("XMODIFIERS")
		c.assignWithEnviron("QT_IM_MODULE")
		c.assignWithEnviron("GTK_IM_MODULE")
	}
}
