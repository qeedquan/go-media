package glew

/*
#include <GL/glew.h>

#cgo pkg-config: glew
*/
import "C"
import (
	"errors"
	"unsafe"
)

func Init() error {
	C.glewExperimental = 1
	rc := C.glewInit()
	if rc != C.GLEW_OK {
		return errorString(int(rc))
	}
	return nil
}

func errorString(rc int) error {
	cstr := C.glewGetErrorString(C.GLenum(rc))
	str := C.GoString((*C.char)(unsafe.Pointer(cstr)))
	return errors.New(str)
}
