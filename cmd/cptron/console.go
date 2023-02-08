//go:build !windows && !appengine
// +build !windows,!appengine

package main

import (
	"os"

	"github.com/mattn/go-isatty"
)

func IsInTerminal() bool {
	fd := os.Stdout.Fd()
	return isatty.IsTerminal(fd)
}
