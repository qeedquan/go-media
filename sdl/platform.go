package sdl

/*
#include "gosdl.h"
*/
import "C"

func GetPlatform() string {
	return C.GoString(C.SDL_GetPlatform())
}
