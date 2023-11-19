package stbhw

/*

#define STB_HEXWAVE_IMPLEMENTATION
#include "stb_hexwave.h"

#cgo LDFLAGS: -lm

*/
import "C"
import "unsafe"

type (
	HexWave C.HexWave
)

func Init(width, oversample int) {
	C.hexwave_init(C.int(width), C.int(oversample), nil)
}

func Shutdown() {
	C.hexwave_shutdown(nil)
}

func Create(reflect int, peak_time, half_height, zero_wait float64) *HexWave {
	h := &HexWave{}
	C.hexwave_create((*C.HexWave)(h), C.int(reflect), C.float(peak_time), C.float(half_height), C.float(zero_wait))
	return h
}

func Change(h *HexWave, reflect int, peak_time, half_height, zero_wait float64) {
	C.hexwave_change((*C.HexWave)(h), C.int(reflect), C.float(peak_time), C.float(half_height), C.float(zero_wait))
}

func Gen(output []float32, num_samples int, h *HexWave, freq float64) {
	C.hexwave_generate_samples((*C.float)(unsafe.Pointer((&output[0]))), C.int(num_samples), (*C.HexWave)(h), C.float(freq))
}
