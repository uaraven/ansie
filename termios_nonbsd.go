//go:build !darwin && !freebsd && !netbsd && !openbsd && !windows

package ansie

import (
    "golang.org/x/sys/unix"
)

const (
    getTermios = unix.TCGETS
    setTermios = unix.TCSETS
)
