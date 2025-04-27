package sdl

/*
#include "gosdl.h"
*/
import "C"
import "unsafe"

func SetClipboardText(text string) error {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	return ek(C.SDL_SetClipboardText(ctext))
}

func GetClipboardText() string {
	ctext := C.SDL_GetClipboardText()
	defer C.free(unsafe.Pointer(ctext))
	return C.GoString(ctext)
}

func HasClipboardText() bool {
	return C.SDL_HasClipboardText() != 0
}
