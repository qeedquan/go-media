package sdl

/*
#include "gosdl.h"
*/
import "C"

type (
	Cursor           C.SDL_Cursor
	MouseButtonFlags C.SDL_MouseButtonFlags
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

func HasMouse() bool {
	return bool(C.SDL_HasMouse())
}

func GetMouseState() (MouseButtonFlags, float32, float32) {
	var x, y C.float
	flags := C.SDL_GetMouseState(&x, &y)
	return MouseButtonFlags(flags), float32(x), float32(y)
}

func WarpMouseInWindow(window *Window, x, y float32) {
	C.SDL_WarpMouseInWindow((*C.SDL_Window)(window), C.float(x), C.float(y))
}

func WarpMouseGlobal(x, y float32) bool {
	return bool(C.SDL_WarpMouseGlobal(C.float(x), C.float(y)))
}
