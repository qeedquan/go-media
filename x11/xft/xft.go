package xft

/*
#include <X11/Xlib.h>
#include <X11/Xft/Xft.h>

#cgo pkg-config: xft fontconfig
*/
import "C"
import (
	"unsafe"

	"github.com/qeedquan/go-media/x11/xlib"
	"github.com/qeedquan/go-media/x11/xlib/xrender"
)

type (
	CharFontSpec  C.XftCharFontSpec
	CharSpec      C.XftCharSpec
	Color         C.XftColor
	Draw          C.XftDraw
	Endian        C.FcEndian
	Font          C.XftFont
	FontInfo      C.XftFontInfo
	GlyphFontSpec C.XftGlyphFontSpec
	GlyphInfo     C.XGlyphInfo
	GlyphSpec     C.XftGlyphSpec
	Pattern       C.XftPattern
	Picture       C.Picture
	Result        C.XftResult
)

const (
	ResultMatch   Result = C.XftResultMatch
	ResultNoMatch Result = C.XftResultNoMatch
)

func InitFtLibrary() bool {
	return C.XftInitFtLibrary() != 0
}

func Init(config string) bool {
	var cconfig *C.char
	if config != "" {
		cconfig = C.CString(config)
		defer C.free(unsafe.Pointer(cconfig))
	}
	return C.XftInit(cconfig) != 0
}

func GetVersion() int {
	return int(C.XftGetVersion())
}

func CharExists(display *xlib.Display, pub *Font, ucs4 rune) bool {
	return C.XftCharIndex((*C.Display)(display), (*C.XftFont)(pub), C.FcChar32(ucs4)) != 0
}

func CharIndex(display *xlib.Display, pub *Font, ucs4 rune) uint {
	return uint(C.XftCharIndex((*C.Display)(display), (*C.XftFont)(pub), C.FcChar32(ucs4)))
}

func DrawCreate(display *xlib.Display, drawable xlib.Drawable, vis *xlib.Visual, cmap xlib.Colormap) *Draw {
	return (*Draw)(C.XftDrawCreate((*C.Display)(display), C.Drawable(drawable), (*C.Visual)(unsafe.Pointer(vis)), C.Colormap(cmap)))
}

func DrawChange(draw *Draw, drawable xlib.Drawable) {
	C.XftDrawChange((*C.XftDraw)(draw), C.Drawable(drawable))
}

func DrawRect(draw *Draw, color *Color, x, y, w, h int) {
	C.XftDrawRect((*C.XftDraw)(draw), (*C.XftColor)(color), C.int(x), C.int(y), C.uint(w), C.uint(h))
}

func DrawPicture(draw *Draw) Picture {
	return Picture(C.XftDrawPicture((*C.XftDraw)(draw)))
}

func DrawSetClipRectangles(draw *Draw, xOrigin, yOrigin int, rects []xlib.Rectangle) {
	C.XftDrawSetClipRectangles((*C.XftDraw)(draw), C.int(xOrigin), C.int(yOrigin), (*C.XRectangle)(unsafe.Pointer(&rects[0])), C.int(len(rects)))
}

func XlfdParse(xlfd string, ignore_scalable, complete bool) *Pattern {
	cxlfd := C.CString(xlfd)
	defer C.free(unsafe.Pointer(cxlfd))
	return (*Pattern)(C.XftXlfdParse(cxlfd, xbool(ignore_scalable), xbool(complete)))
}

func ColorAllocValue(display *xlib.Display, visual *xlib.Visual, cmap xlib.Colormap, color *xrender.Color, result *Color) bool {
	return C.XftColorAllocValue((*C.Display)(display), (*C.Visual)(unsafe.Pointer(visual)), C.Colormap(cmap), (*C.XRenderColor)(unsafe.Pointer(color)), (*C.XftColor)(result)) != 0
}

func ColorAllocName(display *xlib.Display, visual *xlib.Visual, cmap xlib.Colormap, name string, result *Color) bool {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return C.XftColorAllocName((*C.Display)(display), (*C.Visual)(unsafe.Pointer(visual)), C.Colormap(cmap), cname, (*C.XftColor)(result)) != 0
}

func ColorFree(display *xlib.Display, visual *xlib.Visual, cmap xlib.Colormap, col *Color) {
	C.XftColorFree((*C.Display)(display), (*C.Visual)(unsafe.Pointer(visual)), C.Colormap(cmap), (*C.XftColor)(col))
}

func (c *Color) Color() xrender.Color {
	var p xrender.Color
	p.SetRed(uint16(c.color.red))
	p.SetGreen(uint16(c.color.green))
	p.SetBlue(uint16(c.color.blue))
	p.SetAlpha(uint16(c.color.alpha))
	return p
}

func (c *Color) Pixel() uint64 {
	return uint64(c.pixel)
}

func xbool(b bool) C.Bool {
	if b {
		return 1
	}
	return 0
}

func (f *Font) Pattern() *Pattern {
	return (*Pattern)(f.pattern)
}

func FontOpenPattern(display *xlib.Display, pattern *Pattern) *Font {
	return (*Font)(C.XftFontOpenPattern((*C.Display)(display), (*C.XftPattern)(pattern)))
}

func DefaultSubstitute(display *xlib.Display, screen int, pattern *Pattern) {
	C.XftDefaultSubstitute((*C.Display)(display), C.int(screen), (*C.XftPattern)(pattern))
}

func PatternGetInteger(pattern *Pattern, object string, n int) (r Result, i int) {
	var ci C.int
	cobject := C.CString(object)
	defer C.free(unsafe.Pointer(cobject))
	rc := C.XftPatternGetInteger((*C.XftPattern)(pattern), cobject, C.int(n), &ci)
	return Result(rc), int(ci)
}

func TextExtentsUtf8(display *xlib.Display, font *Font, str string, extents *GlyphInfo) {
	cstr := C.CString(str)
	defer C.free(unsafe.Pointer(cstr))
	C.XftTextExtentsUtf8((*C.Display)(display), (*C.XftFont)(font), (*C.uchar)(unsafe.Pointer(cstr)), C.int(len(str)), (*C.XGlyphInfo)(extents))
}

func TextExtentsUtf16(display *xlib.Display, font *Font, endian Endian, str string, extents *GlyphInfo) {
	cstr := C.CString(str)
	defer C.free(unsafe.Pointer(cstr))
	C.XftTextExtentsUtf16((*C.Display)(display), (*C.XftFont)(font), (*C.uchar)(unsafe.Pointer(cstr)), C.FcEndian(endian), C.int(len(str)), (*C.XGlyphInfo)(extents))
}

func GlyphExtents(display *xlib.Display, font *Font, glyphs []int, extents *GlyphInfo) {
	cglyphs := make([]C.FT_UInt, len(glyphs))
	C.XftGlyphExtents((*C.Display)(display), (*C.XftFont)(font), &cglyphs[0], C.int(len(cglyphs)), (*C.XGlyphInfo)(extents))
	for i := range glyphs {
		glyphs[i] = int(cglyphs[i])
	}
}

func FontClose(display *xlib.Display, font *Font) {
	C.XftFontClose((*C.Display)(display), (*C.XftFont)(font))
}

func (g *GlyphFontSpec) SetFont(font *Font) {
	g.font = (*C.XftFont)(font)
}

func (g *GlyphFontSpec) SetGlyph(glyphidx uint) {
	g.glyph = C.uint(glyphidx)
}

func (g *GlyphFontSpec) SetX(x int) {
	g.x = C.short(x)
}

func (g *GlyphFontSpec) SetY(y int) {
	g.y = C.short(y)
}

func (c *Color) SetPixel(pixel uint64) {
	c.pixel = C.ulong(pixel)
}

func (c *Color) SetColor(color xrender.Color) {
	c.color = C.XRenderColor{
		red:   C.ushort(color.Red()),
		green: C.ushort(color.Green()),
		blue:  C.ushort(color.Blue()),
		alpha: C.ushort(color.Alpha()),
	}
}

func (f *Font) Ascent() int {
	return int(f.ascent)
}

func (f *Font) Descent() int {
	return int(f.descent)
}

func (f *Font) MaxAdvanceWidth() int {
	return int(f.max_advance_width)
}

func (g *GlyphInfo) XOff() int {
	return int(g.xOff)
}

func (g *GlyphInfo) YOff() int {
	return int(g.yOff)
}

func (g *GlyphInfo) X() int {
	return int(g.x)
}

func (g *GlyphInfo) Y() int {
	return int(g.y)
}

func DrawSetClip(draw *Draw, region xlib.Region) {
	C.XftDrawSetClip((*C.XftDraw)(draw), C.Region(unsafe.Pointer(region)))
}

func DrawGlyphFontSpec(draw *Draw, color *Color, spec []GlyphFontSpec) {
	C.XftDrawGlyphFontSpec((*C.XftDraw)(draw), (*C.XftColor)(color), (*C.XftGlyphFontSpec)(&spec[0]), C.int(len(spec)))
}
