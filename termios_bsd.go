//go:build darwin || freebsd || openbsd || netbsd

package ansie

import (
	"golang.org/x/sys/unix"
)

const (
	getTermios = unix.TIOCGETA
	setTermios = unix.TIOCSETA
)
