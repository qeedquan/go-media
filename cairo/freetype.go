package cairo

/*
#include <cairo.h>
#include <cairo-ft.h>
#include <ft2build.h>
#include FT_FREETYPE_H
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type FT_Library struct {
	l C.FT_Library
}

type FT_Face struct {
	f C.FT_Face
}

type FT_Error C.FT_Error

func (e FT_Error) Error() string {
	return fmt.Sprintf("FT_Error(%d)", e)
}

func (l *FT_Library) Init() error {
	rc := C.FT_Init_FreeType((*C.FT_Library)(&l.l))
	if rc == 0 {
		return nil
	}
	return FT_Error(rc)
}

func (l *FT_Library) NewFace(filename string, faceIndex int) (*FT_Face, error) {
	f := &FT_Face{}
	cs := C.CString(filename)
	rc := C.FT_New_Face(l.l, cs, C.FT_Long(faceIndex), &f.f)
	C.free(unsafe.Pointer(cs))
	if rc == 0 {
		return f, nil
	}
	return f, FT_Error(rc)
}

func (f *FT_Face) Create(loadFlags int) *FontFace {
	return (*FontFace)(C.cairo_ft_font_face_create_for_ft_face(f.f, C.int(loadFlags)))
}

func (f *FT_Face) Done() {
	C.FT_Done_Face(f.f)
}
