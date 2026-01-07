package sdl

/*
#include "gosdl.h"

static int setError(const char *str) {
	return SDL_SetError("%s", str);
}
*/
import "C"

import (
	"errors"
	"fmt"
	"unsafe"
)

func ek(rc C.int) error {
	if rc < 0 {
		return GetError()
	}
	return nil
}

func SetError(format string, args ...interface{}) int {
	str := C.CString(fmt.Sprintf(format, args...))
	defer C.free(unsafe.Pointer(str))
	return int(C.setError(str))
}

func GetError() error {
	return errors.New(C.GoString(C.SDL_GetError()))
}

func ClearError() {
	C.SDL_ClearError()
}
