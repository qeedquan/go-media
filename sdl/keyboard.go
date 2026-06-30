package sdl

/*
#include "gosdl.h"
*/
import "C"
import (
	"unsafe"
)

type (
	Scancode C.SDL_Scancode
	Keycode  C.SDL_Keycode
)

type Keysym struct {
	Scancode Scancode
	Sym      Keycode
	Mod      uint16
}

func GetKeyboardFocus() *Window {
	return (*Window)(C.SDL_GetKeyboardFocus())
}

func GetKeyName(key Keycode) string {
	return C.GoString(C.SDL_GetKeyName(C.SDL_Keycode(key)))
}

func GetKeyFromName(name string) Keycode {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return Keycode(C.SDL_GetKeyFromName(cname))
}

func GetKeyboardState() []uint8 {
	var numkeys C.int
	state := C.SDL_GetKeyboardState(&numkeys)
	return ((*[1 << 30]uint8)(unsafe.Pointer(state)))[:numkeys:numkeys]
}

func GetModState() Keymod {
	return Keymod(C.SDL_GetModState())
}

func SetModState(modstate Keymod) {
	C.SDL_SetModState(C.SDL_Keymod(modstate))
}

func StartTextInput() {
	C.SDL_StartTextInput()
}

func IsTextInputActive() bool {
	return C.SDL_IsTextInputActive() != 0
}

func StopTextInput() {
	C.SDL_StopTextInput()
}

func SetTextInputRect(r *Rect) {
	C.SDL_SetTextInputRect((*C.SDL_Rect)(unsafe.Pointer(r)))
}

func HasScreenKeyboardSupport() bool {
	return C.SDL_HasScreenKeyboardSupport() != 0
}

func IsScreenKeyboardShown(w *Window) bool {
	return C.SDL_IsScreenKeyboardShown((*C.SDL_Window)(w)) != 0
}

func GetKeyFromScancode(scancode Scancode) Keycode {
	return Keycode(C.SDL_GetKeyFromScancode(C.SDL_Scancode(scancode)))
}

func GetScancodeFromName(name string) Scancode {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return Scancode(C.SDL_GetScancodeFromName(cname))
}
