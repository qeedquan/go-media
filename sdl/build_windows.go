//go:build windows

package sdl

/*
#include "gosdl.h"

#cgo LDFLAGS: -lSDL2
*/
import "C"
