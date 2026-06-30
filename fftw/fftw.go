package fftw

/*
#include <fftw3.h>

#cgo pkg-config: fftw3

*/
import "C"
import (
	"reflect"
	"unsafe"
)

type (
	Plan    C.fftw_plan
	R2RKind C.fftw_r2r_kind
)

const (
	ESTIMATE       = C.FFTW_ESTIMATE
	MEASURE        = C.FFTW_MEASURE
	PATIENT        = C.FFTW_PATIENT
	EXHAUSTIVE     = C.FFTW_EXHAUSTIVE
	WISDOM_ONLY    = C.FFTW_WISDOM_ONLY
	DESTROY_INPUT  = C.FFTW_DESTROY_INPUT
	PRESERVE_INPUT = C.FFTW_PRESERVE_INPUT
	UNALIGNED      = C.FFTW_UNALIGNED
)

const (
	R2HC    = C.FFTW_R2HC
	HC2R    = C.FFTW_HC2R
	DHT     = C.FFTW_DHT
	REDFT00 = C.FFTW_REDFT00
	REDFT10 = C.FFTW_REDFT10
	REDFT01 = C.FFTW_REDFT01
	REDFT11 = C.FFTW_REDFT11
	RODFT00 = C.FFTW_RODFT00
	RODFT10 = C.FFTW_RODFT10
	RODFT01 = C.FFTW_RODFT01
	RODFT11 = C.FFTW_RODFT11
)

func AllocReal(n int) []float64 {
	p := C.fftw_alloc_real(C.size_t(n))
	b := unsafe.Slice((*float64)(p), n)
	for i := range b {
		b[i] = 0
	}
	return b
}

func AllocComplex(n int) []complex128 {
	p := C.fftw_alloc_complex(C.size_t(n))
	b := unsafe.Slice((*complex128)(unsafe.Pointer(p)), n)
	for i := range b {
		b[i] = 0
	}
	return b
}

func Free(p any) {
	rv := reflect.ValueOf(p)
	C.fftw_free(unsafe.Pointer(rv.Pointer()))
}

func PlanDFT_1D(n0 int, in, out []complex128, sign int, flags uint) Plan {
	return Plan(C.fftw_plan_dft_1d(C.int(n0), (*C.fftw_complex)(unsafe.Pointer(&in[0])), (*C.fftw_complex)(unsafe.Pointer(&out[0])), C.int(sign), C.uint(flags)))
}

func PlanDFT_2D(n0, n1 int, in, out []complex128, sign int, flags uint) Plan {
	return Plan(C.fftw_plan_dft_2d(C.int(n0), C.int(n1), (*C.fftw_complex)(unsafe.Pointer(&in[0])), (*C.fftw_complex)(unsafe.Pointer(&out[0])), C.int(sign), C.uint(flags)))
}

func PlanDFT_3D(n0, n1, n2 int, in, out []complex128, sign int, flags uint) Plan {
	return Plan(C.fftw_plan_dft_3d(C.int(n0), C.int(n1), C.int(n2), (*C.fftw_complex)(unsafe.Pointer(&in[0])), (*C.fftw_complex)(unsafe.Pointer(&out[0])), C.int(sign), C.uint(flags)))
}

func PlanDFT(rank int, n []int32, in, out []complex128, sign int, flags uint) Plan {
	return Plan(C.fftw_plan_dft(C.int(rank), (*C.int)(unsafe.Pointer(&n[0])), (*C.fftw_complex)(unsafe.Pointer(&in[0])), (*C.fftw_complex)(unsafe.Pointer(&out[0])), C.int(sign), C.uint(flags)))
}

func PlanDFT_R2C_1D(n0 int, in []float64, out []complex128, flags uint) Plan {
	return Plan(C.fftw_plan_dft_r2c_1d(C.int(n0), (*C.double)(unsafe.Pointer(&in[0])), (*C.fftw_complex)(unsafe.Pointer(&out[0])), C.uint(flags)))
}

func PlanDFT_R2C_2D(n0, n1 int, in []float64, out []complex128, flags uint) Plan {
	return Plan(C.fftw_plan_dft_r2c_2d(C.int(n0), C.int(n1), (*C.double)(unsafe.Pointer(&in[0])), (*C.fftw_complex)(unsafe.Pointer(&out[0])), C.uint(flags)))
}

func PlanDFT_R2C_3D(n0, n1, n2 int, in []float64, out []complex128, flags uint) Plan {
	return Plan(C.fftw_plan_dft_r2c_3d(C.int(n0), C.int(n1), C.int(n2), (*C.double)(unsafe.Pointer(&in[0])), (*C.fftw_complex)(unsafe.Pointer(&out[0])), C.uint(flags)))
}

func PlanDFT_R2C(rank int, n []int32, in []float64, out []complex128, flags uint) Plan {
	return Plan(C.fftw_plan_dft_r2c(C.int(rank), (*C.int)(unsafe.Pointer(&n[0])), (*C.double)(unsafe.Pointer(&in[0])), (*C.fftw_complex)(unsafe.Pointer(&out[0])), C.uint(flags)))
}

func PlanR2R_1D(n0 int, in []float64, out []complex128, kind R2RKind, flags uint) Plan {
	return Plan(C.fftw_plan_r2r_1d(C.int(n0), (*C.double)(unsafe.Pointer(&in[0])), (*C.double)(unsafe.Pointer(&out[0])), C.fftw_r2r_kind(kind), C.uint(flags)))
}

func PlanR2R_2D(n0, n1 int, in, out []float64, kind0, kind1 R2RKind, flags uint) Plan {
	return Plan(C.fftw_plan_r2r_2d(C.int(n0), C.int(n1), (*C.double)(unsafe.Pointer(&in[0])), (*C.double)(unsafe.Pointer(&out[0])), C.fftw_r2r_kind(kind0), C.fftw_r2r_kind(kind1), C.uint(flags)))
}

func PlanR2R_3D(n0, n1, n2 int, in, out []float64, kind0, kind1, kind2 R2RKind, flags uint) Plan {
	return Plan(C.fftw_plan_r2r_3d(C.int(n0), C.int(n1), C.int(n2), (*C.double)(unsafe.Pointer(&in[0])), (*C.double)(unsafe.Pointer(&out[0])), C.fftw_r2r_kind(kind0), C.fftw_r2r_kind(kind1), C.fftw_r2r_kind(kind2), C.uint(flags)))
}

func PlanR2R(rank int, n []int32, in, out []float64, kind []R2RKind, flags uint) Plan {
	return Plan(C.fftw_plan_r2r(C.int(rank), (*C.int)(unsafe.Pointer(&n[0])), (*C.double)(unsafe.Pointer(&in[0])), (*C.double)(unsafe.Pointer(&out[0])), (*C.fftw_r2r_kind)(unsafe.Pointer(&kind[0])), C.uint(flags)))
}

func Execute(plan Plan) {
	C.fftw_execute(C.fftw_plan(plan))
}

func DestroyPlan(plan Plan) {
	C.fftw_destroy_plan(C.fftw_plan(plan))
}

func Cleanup() {
	C.fftw_cleanup()
}

func Cost(plan Plan) float64 {
	return float64(C.fftw_cost(C.fftw_plan(plan)))
}

func SetTimeLimit(seconds float64) {
	C.fftw_set_timelimit(C.double(seconds))
}
