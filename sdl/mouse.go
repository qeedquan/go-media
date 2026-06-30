package sdl

/*
#include "gosdl.h"
*/
import "C"

type (
	Cursor              C.SDL_Cursor
	SystemCursor        C.SDL_SystemCursor
	MouseWheelDirection C.SDL_MouseWheelDirection
)

const (
	SYSTEM_CURSOR_ARROW     SystemCursor = C.SDL_SYSTEM_CURSOR_ARROW
	SYSTEM_CURSOR_IBEAM     SystemCursor = C.SDL_SYSTEM_CURSOR_IBEAM
	SYSTEM_CURSOR_WAIT      SystemCursor = C.SDL_SYSTEM_CURSOR_WAIT
	SYSTEM_CURSOR_CROSSHAIR SystemCursor = C.SDL_SYSTEM_CURSOR_CROSSHAIR
	SYSTEM_CURSOR_WAITARROW SystemCursor = C.SDL_SYSTEM_CURSOR_WAITARROW
	SYSTEM_CURSOR_SIZENWSE  SystemCursor = C.SDL_SYSTEM_CURSOR_SIZENWSE
	SYSTEM_CURSOR_SIZENESW  SystemCursor = C.SDL_SYSTEM_CURSOR_SIZENESW
	SYSTEM_CURSOR_SIZEWE    SystemCursor = C.SDL_SYSTEM_CURSOR_SIZEWE
	SYSTEM_CURSOR_SIZENS    SystemCursor = C.SDL_SYSTEM_CURSOR_SIZENS
	SYSTEM_CURSOR_SIZEALL   SystemCursor = C.SDL_SYSTEM_CURSOR_SIZEALL
	SYSTEM_CURSOR_NO        SystemCursor = C.SDL_SYSTEM_CURSOR_NO
	SYSTEM_CURSOR_HAND      SystemCursor = C.SDL_SYSTEM_CURSOR_HAND
	NUM_SYSTEM_CURSORS                   = C.SDL_NUM_SYSTEM_CURSORS
)

func GetMouseFocus() *Window {
	return (*Window)(C.SDL_GetMouseFocus())
}

func GetMouseState() (x, y int, button uint32) {
	var cx, cy C.int
	cbutton := C.SDL_GetMouseState(&cx, &cy)
	return int(cx), int(cy), uint32(cbutton)
}

func GetGlobalMouseState() (x, y int, button uint32) {
	var cx, cy C.int
	cbutton := C.SDL_GetGlobalMouseState(&cx, &cy)
	return int(cx), int(cy), uint32(cbutton)
}

func GetRelativeMouseState() (x, y int, button uint32) {
	var cx, cy C.int
	cbutton := C.SDL_GetRelativeMouseState(&cx, &cy)
	return int(cx), int(cy), uint32(cbutton)
}

func GetCursor() *Cursor {
	return (*Cursor)(C.SDL_GetCursor())
}

func GetDefaultCursor() *Cursor {
	return (*Cursor)(C.SDL_GetDefaultCursor())
}

func WarpMouseInWindow(w *Window, x, y int) {
	C.SDL_WarpMouseInWindow((*C.SDL_Window)(w), C.int(x), C.int(y))
}

func WarpMouseGlobal(x, y int) error {
	return ek(C.SDL_WarpMouseGlobal(C.int(x), C.int(y)))
}

func CaptureMouse(enabled bool) error {
	return ek(C.SDL_CaptureMouse(truth(enabled)))
}

func GetRelativeMouseMode() bool {
	return C.SDL_GetRelativeMouseMode() != 0
}

func CreateSystemCursor(id SystemCursor) *Cursor {
	return (*Cursor)(C.SDL_CreateSystemCursor(C.SDL_SystemCursor(id)))
}

func (c *Cursor) Free() {
	C.SDL_FreeCursor((*C.SDL_Cursor)(c))
}

func ShowCursor(toggle int) int {
	return int(C.SDL_ShowCursor(C.int(toggle)))
}

func BUTTON(x uint) uint32 {
	return 1 << (x - 1)
}

const (
	BUTTON_LEFT   = C.SDL_BUTTON_LEFT
	BUTTON_MIDDLE = C.SDL_BUTTON_MIDDLE
	BUTTON_RIGHT  = C.SDL_BUTTON_RIGHT
	BUTTON_X1     = C.SDL_BUTTON_X1
	BUTTON_X2     = C.SDL_BUTTON_X2
	BUTTON_LMASK  = C.SDL_BUTTON_LMASK
	BUTTON_MMASK  = C.SDL_BUTTON_MMASK
	BUTTON_RMASK  = C.SDL_BUTTON_RMASK
	BUTTON_X1MASK = C.SDL_BUTTON_X1MASK
	BUTTON_X2MASK = C.SDL_BUTTON_X2MASK
)

const (
	MOUSEWHEEL_NORMAL  = C.SDL_MOUSEWHEEL_NORMAL
	MOUSEWHEEL_FLIPPED = C.SDL_MOUSEWHEEL_FLIPPED
)

func SetCursor(cursor *Cursor) {
	C.SDL_SetCursor((*C.SDL_Cursor)(cursor))
}
