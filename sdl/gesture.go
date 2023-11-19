package sdl

/*
#include "gosdl.h"
*/
import "C"

type GestureID C.SDL_GestureID

func RecordGesture(touchId TouchID) int {
	return int(C.SDL_RecordGesture(C.SDL_TouchID(touchId)))
}
