package sdl

/*
#include "gosdl.h"
*/
import "C"

func truth(x bool) C.SDL_bool {
	if x {
		return C.SDL_TRUE
	}
	return C.SDL_FALSE
}
