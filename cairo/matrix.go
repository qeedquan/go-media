package cairo

/*
#include <cairo.h>
*/
import "C"
import (
	"unsafe"
)

type Matrix struct {
	XX, YX float64
	XY, YY float64
	X0, Y0 float64
}

func (m *Matrix) Init(xx, yx, xy, yy, x0, y0 float64) {
	C.cairo_matrix_init((*C.cairo_matrix_t)(unsafe.Pointer(m)), C.double(xx), C.double(yx), C.double(xy), C.double(yy), C.double(x0), C.double(y0))
}

func (m *Matrix) InitIdentity() {
	C.cairo_matrix_init_identity((*C.cairo_matrix_t)(unsafe.Pointer(m)))
}

func (m *Matrix) InitScale(sx, sy float64) {
	C.cairo_matrix_init_scale((*C.cairo_matrix_t)(unsafe.Pointer(m)), C.double(sx), C.double(sy))
}

func (m *Matrix) InitTranslate(tx, ty float64) {
	C.cairo_matrix_init_translate((*C.cairo_matrix_t)(unsafe.Pointer(m)), C.double(tx), C.double(ty))
}

func (m *Matrix) InitRotate(radians float64) {
	C.cairo_matrix_init_rotate((*C.cairo_matrix_t)(unsafe.Pointer(m)), C.double(radians))
}

func (m *Matrix) Invert() error {
	return xk(C.cairo_matrix_invert((*C.cairo_matrix_t)(unsafe.Pointer(m))))
}
