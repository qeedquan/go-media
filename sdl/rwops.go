package sdl

/*
#include "gosdl.h"
*/
import "C"

import (
	"io"
	"unsafe"
)

type RWOps struct {
	io.ReadWriteSeeker
	io.Closer
	Ops unsafe.Pointer
	mem unsafe.Pointer
}

func RWFromFile(file, mode string) (*RWOps, error) {
	cfile := C.CString(file)
	cmode := C.CString(mode)
	defer C.free(unsafe.Pointer(cfile))
	defer C.free(unsafe.Pointer(cmode))
	ops := C.SDL_RWFromFile(cfile, cmode)
	if ops == nil {
		return nil, GetError()
	}
	return &RWOps{Ops: unsafe.Pointer(ops)}, nil
}

func RWFromConstMem(mem []byte) (*RWOps, error) {
	ptr := C.malloc(C.size_t(len(mem)))
	C.memmove(ptr, unsafe.Pointer(&mem[0]), C.size_t(len(mem)))
	ops := C.SDL_RWFromConstMem(ptr, C.int(len(mem)))
	if ops == nil {
		C.free(ptr)
		return nil, GetError()
	}
	return &RWOps{Ops: unsafe.Pointer(ops), mem: ptr}, nil
}

func (rw *RWOps) Free() {
	C.SDL_FreeRW((*C.SDL_RWops)(rw.Ops))
	if rw.mem != nil {
		C.free(rw.mem)
	}
}
