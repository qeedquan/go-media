// +build linux darwin openbsd netbsd freebsd dragonfly solaris

package sdlimage

/*
#include <SDL.h>

#cgo LDFLAGS: -lSDL2
*/
import "C"
