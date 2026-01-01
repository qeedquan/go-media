package lz4

/*
#include "lz4frame.h"
#include "lz4.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func Compress(src, dst []byte) (int, error) {
	n := C.LZ4_compress_default((*C.char)(unsafe.Pointer(&src[0])), (*C.char)(unsafe.Pointer(&dst[0])), C.int(len(src)), C.int(len(dst)))
	if n < 1 {
		return 0, fmt.Errorf("compression failed (%v)", n)
	}
	return int(n), nil
}

func Decompress(src, dst []byte) (int, error) {
	n := int(C.LZ4_decompress_safe((*C.char)(unsafe.Pointer(&src[0])), (*C.char)(unsafe.Pointer(&dst[0])), C.int(len(src)), C.int(len(dst))))
	if n < 0 {
		return 0, fmt.Errorf("decompression failed (%v)", n)
	}
	return n, nil
}
