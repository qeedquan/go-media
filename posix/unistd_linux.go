// +build linux
package posix

import (
	"syscall"
)

func Geteuid() int {
	return syscall.Geteuid()
}
