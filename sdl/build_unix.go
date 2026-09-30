//go:build linux || darwin || openbsd || netbsd || freebsd || dragonfly || solaris

package sdl

/*
#include "gosdl.h"

#cgo pkg-config: sdl2
*/
import "C"
