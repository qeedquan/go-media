package sdl

/*
#include "gosdl.h"
*/
import "C"

type (
	Cursor C.SDL_Cursor
)

func CursorVisible() bool {
	return bool(C.SDL_CursorVisible())
}

func ShowCursor() bool {
	return bool(C.SDL_ShowCursor())
}

func HideCursor() bool {
	return bool(C.SDL_HideCursor())
}
