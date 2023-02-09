//go:build !windows && !appengine
// +build !windows,!appengine

package cptron

import (
	"os"

	"github.com/mattn/go-isatty"
)

func IsInTerminal() bool {
	fd := os.Stdout.Fd()
	return isatty.IsTerminal(fd)
}
