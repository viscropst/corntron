//go:build !windows && !wasm && !wasi
// +build !windows,!wasm,!wasi

package internal

func (c Environ) PrepareEnvsByEnviron(profileDir ...string) {
	c.assignWithEnviron("SSH_AUTH_SOCK")
	c.assignWithEnviron("SSH_ASKPASS")
	c.assignWithEnviron("USER")
	c.assignWithEnviron("PWD")
	c.assignWithEnviron("LANG")
	c.assignWithEnviron("TMPDIR")
	c.assignWithEnviron("TERM")
	c.assignWithEnviron("DBUS_SESSION_BUS_ADDRESS")
	c.assignWithEnviron("EDITOR")
	if len(profileDir) > 0 {
		c["HOME"] = profileDir[0]
	} else {
		c.assignWithEnviron("HOME")
	}
	if AssertWithEnviron("DESKTOP_SESSION") ||
		AssertWithEnviron("XDG_SESSION_TYPE") {
		c.unixEnvWithDesktop()
	}
	if AssertWithEnviron("XDG_SEAT") ||
		AssertWithEnviron("DISPLAY") {
		c.unixEnvWithDesktop()
	}
}

func (c Environ) unixEnvWithDesktop() {
	c.assignWithEnviron("DISPLAY")
	c.assignWithEnviron("WINDOWID")
	c.assignWithEnviron("SESSION_MANAGER")
	c.assignWithEnviron("XDG_DATA_DIRS")
	c.assignWithEnviron("XDG_CONFIG_DIRS")
	c.assignWithEnviron("XDG_CONFIG_HOME")
	c.assignWithEnviron("XDG_RUNTIME_DIR")
	c.assignWithEnviron("XDG_CACHE_HOME")
	c.assignWithEnviron("XDG_SESSION_TYPE")
	c.assignWithEnviron("XDG_SESSION_PATH")
	c.assignWithEnviron("XDG_SESSION_CLASS")
	c.assignWithEnviron("XDG_SESSION_ID")
	c.assignWithEnviron("XDG_SEAT")
	c.assignWithEnviron("XDG_SEAT_PATH")
	c.assignWithEnviron("XDG_CURRENT_DESKTOP")
	c.assignWithEnviron("XDG_SESSION_DESKTOP")
	c.assignWithEnviron("ICEAUTHORITY")
	c.assignWithEnviron("XAUTHORITY")
	c.assignWithEnviron("GTK_RC_FILES")
	c.assignWithEnviron("GTK2_RC_FILES")
	c.assignWithEnviron("DCONF_PROFILE")
	c.assignWithEnviron("GDK_BACKEND")
	c.assignWithEnviron("QT_QPA_PLATFORM")
	if AssertWithEnviron("QT_ENABLE_HIGHDPI_SCALING", "1") {
		c.assignWithEnviron("QT_ENABLE_HIGHDPI_SCALING")
		c.assignWithEnviron("QT_AUTO_SCREEN_SCALE_FACTOR")
	}
	if AssertWithEnviron("DESKTOP_SESSION", "plasma") {
		c.assignWithEnviron("PAM_KWALLET5_LOGIN")
		c.assignWithEnviron("KDE_SESSION_VERSION")
		c.assignWithEnviron("KDE_FULL_SESSION")
	}

	if AssertWithEnviron("XDG_SESSION_TYPE", "wayland") {
		c.assignWithEnviron("WAYLAND_DISPLAY")
		c.assignWithEnviron("QT_WAYLAND_RECONNECT")
		c.assignWithEnviron("QT_WAYLAND_DECORATIONS")
		c.assignWithEnviron("MOZ_ENABLE_WAYLAND")
	}

	if AssertWithEnviron("XDG_SESSION_TYPE", "x11") {
		c.assignWithEnviron("INPUT_METHOD")
		c.assignWithEnviron("XMODIFIERS")
		c.assignWithEnviron("QT_IM_MODULE")
		c.assignWithEnviron("GTK_IM_MODULE")
	}
}
