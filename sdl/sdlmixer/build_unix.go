// +build linux darwin freebsd openbsd netbsd dragonfly solaris

package sdlmixer

/*
#cgo pkg-config: sdl2
#cgo LDFLAGS: -lSDL2_mixer
*/
import "C"
