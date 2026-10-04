package sdl

/*
#include "gosdl.h"
*/
import "C"

type Version struct {
	Major, Minor, Patch uint8
}

const (
	MAJOR_VERSION = C.SDL_MAJOR_VERSION
	MINOR_VERSION = C.SDL_MINOR_VERSION
	PATCHLEVEL    = C.SDL_PATCHLEVEL
)

func GetVersion() Version {
	var v C.SDL_version
	C.SDL_GetVersion(&v)
	return Version{
		uint8(v.major),
		uint8(v.minor),
		uint8(v.patch),
	}
}

func GetRevision() string {
	return C.GoString(C.SDL_GetRevision())
}

