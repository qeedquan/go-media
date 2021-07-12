package cairo

/*
#include <cairo.h>
#include <cairo-ft.h>
#include <stdlib.h>
*/
import "C"
import "unsafe"

type FontSlant C.cairo_font_slant_t
type FontWeight C.cairo_font_weight_t
type FontOptions C.cairo_font_options_t
type FontFace C.cairo_font_face_t
type ScaledFont C.cairo_scaled_font_t

type FontExtents struct {
	Ascent      float64
	Descent     float64
	Height      float64
	MaxXAdvance float64
	MaxYAdvance float64
}

type TextExtents struct {
	XBearing float64
	YBearing float64
	Width    float64
	Height   float64
	XAdvance float64
	YAdvance float64
}

const (
	FONT_SLANT_NORMAL  FontSlant = C.CAIRO_FONT_SLANT_NORMAL
	FONT_SLANT_ITALIC  FontSlant = C.CAIRO_FONT_SLANT_ITALIC
	FONT_SLANT_OBLIQUE FontSlant = C.CAIRO_FONT_SLANT_OBLIQUE
)

const (
	FONT_WEIGHT_NORMAL FontWeight = C.CAIRO_FONT_WEIGHT_NORMAL
	FONT_WEIGHT_BOLD   FontWeight = C.CAIRO_FONT_WEIGHT_BOLD
)

func (c *Cairo) SelectFontFace(family string, slant FontSlant, weight FontWeight) {
	cs := C.CString(family)
	C.cairo_select_font_face((*C.cairo_t)(c), cs, C.cairo_font_slant_t(slant), C.cairo_font_weight_t(weight))
	C.free(unsafe.Pointer(cs))
}

func (c *Cairo) SetFontSize(size float64) {
	C.cairo_set_font_size((*C.cairo_t)(c), C.double(size))
}

func (c *Cairo) ShowText(text string) {
	cs := C.CString(text)
	C.cairo_show_text((*C.cairo_t)(c), cs)
	C.free(unsafe.Pointer(cs))
}

func (c *Cairo) TextPath(text string) {
	cs := C.CString(text)
	C.cairo_text_path((*C.cairo_t)(c), cs)
	C.free(unsafe.Pointer(cs))
}

func (f *ScaledFont) Extents() FontExtents {
	var fe FontExtents
	C.cairo_scaled_font_extents((*C.cairo_scaled_font_t)(f), (*C.cairo_font_extents_t)(unsafe.Pointer(&fe)))
	return fe
}

func (c *Cairo) TextExtents(utf8 string) TextExtents {
	var extents TextExtents
	var buf [1024]byte
	if len(utf8) < len(buf) {
		copy(buf[:], utf8)
		C.cairo_text_extents((*C.cairo_t)(c), (*C.char)(unsafe.Pointer(&buf[0])), (*C.cairo_text_extents_t)((unsafe.Pointer(&extents))))
	} else {
		cs := C.CString(utf8)
		C.cairo_text_extents((*C.cairo_t)(c), cs, (*C.cairo_text_extents_t)((unsafe.Pointer(&extents))))
		C.free(unsafe.Pointer(cs))
	}
	return extents
}

func (f *FontFace) Ref() *FontFace {
	return (*FontFace)(C.cairo_font_face_reference((*C.cairo_font_face_t)(f)))
}

func (f *FontFace) Destroy() {
	C.cairo_font_face_destroy((*C.cairo_font_face_t)(f))
}

func (f *FontFace) CreateScaledFont(matrix, ctm *Matrix, options *FontOptions) *ScaledFont {
	return (*ScaledFont)(C.cairo_scaled_font_create((*C.cairo_font_face_t)(f), (*C.cairo_matrix_t)(unsafe.Pointer(matrix)), (*C.cairo_matrix_t)(unsafe.Pointer(ctm)), (*C.cairo_font_options_t)(options)))
}
