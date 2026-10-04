// +build linux darwin openbsd netbsd freebsd dragonfly solaris

package sdlimage

/*
#include <SDL.h>

#cgo pkg-config: sdl2
*/
import "C"
