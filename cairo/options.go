package cairo

/*
#include <cairo.h>
*/
import "C"

func CreateFontOptions() *FontOptions {
	return (*FontOptions)(C.cairo_font_options_create())
}

func (f *FontOptions) Destroy() {
	C.cairo_font_options_destroy((*C.cairo_font_options_t)(f))
}

func (f *FontOptions) Status() Status {
	return Status(C.cairo_font_options_status((*C.cairo_font_options_t)(f)))
}

func (f *FontOptions) Hash() uint {
	return uint(C.cairo_font_options_hash((*C.cairo_font_options_t)(f)))
}

func (f *FontOptions) SetAntiAlias(antialias AntiAlias) {
	C.cairo_font_options_set_antialias((*C.cairo_font_options_t)(f), C.cairo_antialias_t(antialias))
}

func (f *FontOptions) SetHintStyle(style HintStyle) {
	C.cairo_font_options_set_hint_style((*C.cairo_font_options_t)(f), C.cairo_hint_style_t(style))
}
