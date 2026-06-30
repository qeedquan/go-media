package freetype

/*
#include <ft2build.h>
#include FT_FREETYPE_H
#include FT_STROKER_H
#include FT_LCD_FILTER_H
*/
import "C"
import (
	"unsafe"
)

type Library struct {
	l C.FT_Library
}

type Face struct {
	f C.FT_Face
}

type Stroker struct {
	s C.FT_Stroker
}

type Size struct {
	s C.FT_Size
}

type SizeMetrics struct {
	sm C.FT_Size_Metrics
}

type GlyphSlot struct {
	gs C.FT_GlyphSlot
}

type Glyph struct {
	g C.FT_Glyph
}

type Bitmap struct {
	b C.FT_Bitmap
}

type BitmapGlyph struct {
	bg C.FT_BitmapGlyph
}

type Pos C.FT_Pos
type F26Dot6 C.FT_F26Dot6
type Fixed C.FT_Fixed
type LCDFilter C.FT_LcdFilter

func NewLibrary() (Library, error) {
	var library C.FT_Library

	err := C.FT_Init_FreeType(&library)
	if err != 0 {
		return Library{}, Error(err)
	}

	return Library{
		l: library,
	}, nil
}

func (l Library) NewFace(filename string, ptSize int) (Face, error) {
	var face C.FT_Face

	cs := C.CString(filename)
	defer C.free(unsafe.Pointer(cs))
	err := C.FT_New_Face(l.l, cs, C.FT_Long(ptSize), &face)
	if err != 0 {
		return Face{}, Error(err)
	}

	return Face{
		f: face,
	}, nil
}

func (l Library) NewStroker() (*Stroker, error) {
	var stroker C.FT_Stroker

	err := C.FT_Stroker_New(l.l, &stroker)
	if err != 0 {
		return nil, Error(err)
	}

	return &Stroker{stroker}, nil
}

func (l Library) SetLCDFilter(filter LCDFilter) {
	C.FT_Library_SetLcdFilter(l.l, C.FT_LcdFilter(filter))
}

func (l Library) SetLCDFilterWeights(weights []byte) {
	C.FT_Library_SetLcdFilterWeights(l.l, (*C.uchar)(unsafe.Pointer(&weights[0])))
}

func (l Library) Done() error {
	if l.l == nil {
		return nil
	}
	err := C.FT_Done_FreeType(l.l)
	if err != 0 {
		return Error(err)
	}
	return nil
}

func (f Face) SetPixelSizes(pixelWidth, pixelHeight uint) {
	C.FT_Set_Pixel_Sizes(f.f, C.FT_UInt(pixelWidth), C.FT_UInt(pixelHeight))
}

func (f Face) LoadChar(charCode rune, loadFlags int32) error {
	err := C.FT_Load_Char(f.f, C.FT_ULong(charCode), C.FT_Int32(loadFlags))
	if err != 0 {
		return Error(err)
	}
	return nil
}

func (f Face) LoadGlyph(glyphIndex uint, loadFlags int32) error {
	err := C.FT_Load_Glyph(f.f, C.FT_UInt(glyphIndex), C.FT_Int32(loadFlags))
	if err != 0 {
		return Error(err)
	}
	return nil
}

func (f Face) CharIndex(charCode uint) uint {
	return uint(C.FT_Get_Char_Index(f.f, C.FT_ULong(charCode)))
}

func (f Face) SelectCharmap(charmap Encoding) error {
	err := C.FT_Select_Charmap(f.f, C.FT_Encoding(charmap))
	if err != 0 {
		return Error(err)
	}
	return nil
}

func (f Face) SelectSize(strikeIndex int) error {
	err := C.FT_Select_Size(f.f, C.FT_Int(strikeIndex))
	if err != 0 {
		return Error(err)
	}
	return nil
}

func (f Face) SetCharSize(charWidth, charHeight F26Dot6, horzResolution, vertResolution uint) error {
	err := C.FT_Set_Char_Size(f.f, C.FT_F26Dot6(charWidth), C.FT_F26Dot6(charHeight), C.FT_UInt(horzResolution), C.FT_UInt(vertResolution))
	if err != 0 {
		return Error(err)
	}
	return nil
}

func (f Face) SetTransform(matrix *Matrix, vector *Vector) {
	C.FT_Set_Transform(f.f, (*C.FT_Matrix)(unsafe.Pointer(matrix)), (*C.FT_Vector)(unsafe.Pointer(vector)))
}

func (f Face) Kerning(leftGlyph, rightGlyph uint, kernMode KerningMode) (Vector, error) {
	var akerning C.FT_Vector
	err := C.FT_Get_Kerning(f.f, C.FT_UInt(leftGlyph), C.FT_UInt(rightGlyph), C.FT_UInt(kernMode), &akerning)
	if err != 0 {
		return Vector{Pos(akerning.x), Pos(akerning.y)}, Error(err)
	}
	return Vector{Pos(akerning.x), Pos(akerning.y)}, nil
}

func (f Face) Done() error {
	if f.f == nil {
		return nil
	}
	err := C.FT_Done_Face(f.f)
	if err != 0 {
		return Error(err)
	}
	return nil
}

func (s Stroker) Set(radius Fixed, linecap StrokerLineCap, linejoin StrokerLineJoin, miterLimit Fixed) {
	C.FT_Stroker_Set(s.s, C.FT_Fixed(radius), C.FT_Stroker_LineCap(linecap), C.FT_Stroker_LineJoin(linejoin), C.FT_Fixed(miterLimit))
}

func (s Stroker) Done() {
	if s.s == nil {
		return
	}
	C.FT_Stroker_Done(s.s)
}

func (f Face) UnderlinePosition() int  { return int(f.f.underline_position) }
func (f Face) UnderlineThickness() int { return int(f.f.underline_thickness) }
func (f Face) Size() Size              { return Size{f.f.size} }
func (f Face) Glyph() GlyphSlot        { return GlyphSlot{f.f.glyph} }

func (s Size) Metrics() SizeMetrics { return SizeMetrics{s.s.metrics} }

func (sm SizeMetrics) Height() Pos    { return Pos(sm.sm.height) }
func (sm SizeMetrics) Descender() Pos { return Pos(sm.sm.descender) }
func (sm SizeMetrics) Ascender() Pos  { return Pos(sm.sm.ascender) }

func (gs GlyphSlot) Bitmap() Bitmap  { return Bitmap{gs.gs.bitmap} }
func (gs GlyphSlot) Advance() Vector { return Vector{Pos(gs.gs.advance.x), Pos(gs.gs.advance.y)} }
func (gs GlyphSlot) BitmapTop() int  { return int(gs.gs.bitmap_top) }
func (gs GlyphSlot) BitmapLeft() int { return int(gs.gs.bitmap_left) }

func (gs GlyphSlot) Glyph() (Glyph, error) {
	var glyph C.FT_Glyph
	err := C.FT_Get_Glyph(gs.gs, &glyph)
	if err != 0 {
		return Glyph{}, Error(err)
	}
	return Glyph{glyph}, nil
}

func (g Glyph) Stroke(stroker *Stroker, destroy bool) error {
	var b C.FT_Bool
	if destroy {
		b = 1
	}
	err := C.FT_Glyph_Stroke(&g.g, stroker.s, b)
	if err != 0 {
		return Error(err)
	}
	return nil
}

func (g Glyph) StrokeBorder(stroker *Stroker, inside, destroy bool) error {
	var i, b C.FT_Bool
	if inside {
		i = 1
	}
	if destroy {
		b = 1
	}

	err := C.FT_Glyph_StrokeBorder(&g.g, stroker.s, i, b)
	if err != 0 {
		return Error(err)
	}
	return nil
}

func (g Glyph) ToBitmap(renderMode RenderMode, origin *Vector, destroy bool) error {
	var b C.FT_Bool
	if destroy {
		b = 1
	}

	err := C.FT_Glyph_To_Bitmap(&g.g, C.FT_Render_Mode(renderMode), (*C.FT_Vector)(unsafe.Pointer(origin)), b)
	if err != 0 {
		return Error(err)
	}
	return nil
}

func (g Glyph) BitmapGlyph() BitmapGlyph {
	return BitmapGlyph{(C.FT_BitmapGlyph)(unsafe.Pointer(g.g))}
}

func (g Glyph) Done() {
	if g.g == nil {
		return
	}
	C.FT_Done_Glyph(g.g)
}

func (bg BitmapGlyph) Bitmap() Bitmap { return Bitmap{bg.bg.bitmap} }
func (bg BitmapGlyph) Top() int       { return int(bg.bg.top) }
func (bg BitmapGlyph) Left() int      { return int(bg.bg.left) }

func (b Bitmap) Width() int { return int(b.b.width) }
func (b Bitmap) Rows() int  { return int(b.b.rows) }
func (b Bitmap) Pitch() int { return int(b.b.pitch) }

func (b Bitmap) Buffer() []byte {
	return unsafe.Slice((*byte)(b.b.buffer), b.Rows()*b.Pitch())
}
