package sdlgfx

/*
#include "SDL2_gfxPrimitives.h"
#include "SDL2_framerate.h"
#include "gfx.h"
*/
import "C"

import (
	"errors"
	"unsafe"

	"github.com/qeedquan/go-media/sdl"
)

func Pixel(re *sdl.Renderer, x, y int, c sdl.Color) error {
	return ek(C.pixelRGBA((*C.SDL_Renderer)(re), C.Sint16(x), C.Sint16(y), C.Uint8(c.R), C.Uint8(c.G), C.Uint8(c.B), C.Uint8(c.A)))
}

func Hline(re *sdl.Renderer, x1, x2, y int, c sdl.Color) error {
	return ek(C.hlineRGBA((*C.SDL_Renderer)(re), C.Sint16(x1), C.Sint16(x2), C.Sint16(y), C.Uint8(c.R), C.Uint8(c.G), C.Uint8(c.B), C.Uint8(c.A)))
}

func Vline(re *sdl.Renderer, x, y1, y2 int, c sdl.Color) error {
	return ek(C.vlineRGBA((*C.SDL_Renderer)(re), C.Sint16(x), C.Sint16(y1), C.Sint16(y2), C.Uint8(c.R), C.Uint8(c.G), C.Uint8(c.B), C.Uint8(c.A)))
}

func Rectangle(re *sdl.Renderer, x1, y1, x2, y2 int, c sdl.Color) error {
	return ek(C.rectangleRGBA((*C.SDL_Renderer)(re), C.Sint16(x1), C.Sint16(y1), C.Sint16(x2), C.Sint16(y2), C.Uint8(c.R), C.Uint8(c.G), C.Uint8(c.B), C.Uint8(c.A)))
}

func RoundedRectangle(re *sdl.Renderer, x1, y1, x2, y2, rad int, c sdl.Color) error {
	return ek(C.roundedRectangleRGBA((*C.SDL_Renderer)(re), C.Sint16(x1), C.Sint16(y1), C.Sint16(x2), C.Sint16(y2), C.Sint16(rad), C.Uint8(c.R), C.Uint8(c.G), C.Uint8(c.B), C.Uint8(c.A)))
}

func Box(re *sdl.Renderer, x1, y1, x2, y2 int, c sdl.Color) error {
	return ek(C.boxRGBA((*C.SDL_Renderer)(re), C.Sint16(x1), C.Sint16(y1), C.Sint16(x2), C.Sint16(y2), C.Uint8(c.R), C.Uint8(c.G), C.Uint8(c.B), C.Uint8(c.A)))
}

func RoundedBox(re *sdl.Renderer, x1, y1, x2, y2, rad int, c sdl.Color) error {
	return ek(C.roundedBoxRGBA((*C.SDL_Renderer)(re), C.Sint16(x1), C.Sint16(y1), C.Sint16(x2), C.Sint16(y2), C.Sint16(rad), C.Uint8(c.R), C.Uint8(c.G), C.Uint8(c.B), C.Uint8(c.A)))
}

func Line(re *sdl.Renderer, x1, y1, x2, y2 int, c sdl.Color) error {
	return ek(C.lineRGBA((*C.SDL_Renderer)(re), C.Sint16(x1), C.Sint16(y1), C.Sint16(x2), C.Sint16(y2), C.Uint8(c.R), C.Uint8(c.G), C.Uint8(c.B), C.Uint8(c.A)))
}

func AALine(re *sdl.Renderer, x1, y1, x2, y2 int, c sdl.Color) error {
	return ek(C.aalineRGBA((*C.SDL_Renderer)(re), C.Sint16(x1), C.Sint16(y1), C.Sint16(x2), C.Sint16(y2), C.Uint8(c.R), C.Uint8(c.G), C.Uint8(c.B), C.Uint8(c.A)))
}

func ThickLine(re *sdl.Renderer, x1, y1, x2, y2, width int, c sdl.Color) error {
	return ek(C.thickLineRGBA((*C.SDL_Renderer)(re), C.Sint16(x1), C.Sint16(y1), C.Sint16(x2), C.Sint16(y2), C.Uint8(width), C.Uint8(c.R), C.Uint8(c.G), C.Uint8(c.B), C.Uint8(c.A)))
}

func Circle(re *sdl.Renderer, x, y, rad int, c sdl.Color) error {
	return ek(C.goCircle((*C.SDL_Renderer)(re), C.int(x), C.int(y), C.int(rad), C.Uint8(c.R), C.Uint8(c.G), C.Uint8(c.B), C.Uint8(c.A)))
}

func Arc(re *sdl.Renderer, x, y, rad, start, end int, c sdl.Color) error {
	return ek(C.arcRGBA((*C.SDL_Renderer)(re), C.Sint16(x), C.Sint16(y), C.Sint16(rad), C.Sint16(start), C.Sint16(end), C.Uint8(c.R), C.Uint8(c.G), C.Uint8(c.B), C.Uint8(c.A)))
}

func AACircle(re *sdl.Renderer, x, y, rad int, c sdl.Color) error {
	return ek(C.aacircleRGBA((*C.SDL_Renderer)(re), C.Sint16(x), C.Sint16(y), C.Sint16(rad), C.Uint8(c.R), C.Uint8(c.G), C.Uint8(c.B), C.Uint8(c.A)))
}

func StrokeCircle(re *sdl.Renderer, x, y, rad int, stroke sdl.Color, fill sdl.Color) error {
	err := FilledCircle(re, x, y, rad, fill)
	xerr := Circle(re, x, y, rad, stroke)
	if err == nil {
		err = xerr
	}
	return err
}

func AAStrokeCircle(re *sdl.Renderer, x, y, rad int, stroke sdl.Color, fill sdl.Color) error {
	err := FilledCircle(re, x, y, rad, fill)
	xerr := AACircle(re, x, y, rad, stroke)
	if err == nil {
		err = xerr
	}
	return err
}

func FilledCircle(re *sdl.Renderer, x, y, rad int, c sdl.Color) error {
	return ek(C.goFilledCircle((*C.SDL_Renderer)(re), C.int(x), C.int(y), C.int(rad), C.Uint8(c.R), C.Uint8(c.G), C.Uint8(c.B), C.Uint8(c.A)))
}

func Ellipse(re *sdl.Renderer, x, y, rx, ry int, c sdl.Color) error {
	return ek(C.ellipseRGBA((*C.SDL_Renderer)(re), C.Sint16(x), C.Sint16(y), C.Sint16(rx), C.Sint16(ry), C.Uint8(c.R), C.Uint8(c.G), C.Uint8(c.B), C.Uint8(c.A)))
}

func AAEllipse(re *sdl.Renderer, x, y, rx, ry int, c sdl.Color) error {
	return ek(C.aaellipseRGBA((*C.SDL_Renderer)(re), C.Sint16(x), C.Sint16(y), C.Sint16(rx), C.Sint16(ry), C.Uint8(c.R), C.Uint8(c.G), C.Uint8(c.B), C.Uint8(c.A)))
}

func FilledEllipse(re *sdl.Renderer, x, y, rx, ry int, c sdl.Color) error {
	return ek(C.filledEllipseRGBA((*C.SDL_Renderer)(re), C.Sint16(x), C.Sint16(y), C.Sint16(rx), C.Sint16(ry), C.Uint8(c.R), C.Uint8(c.G), C.Uint8(c.B), C.Uint8(c.A)))
}

func Pie(re *sdl.Renderer, x, y, rad, start, end int, c sdl.Color) error {
	return ek(C.pieRGBA((*C.SDL_Renderer)(re), C.Sint16(x), C.Sint16(y), C.Sint16(rad), C.Sint16(start), C.Sint16(end), C.Uint8(c.R), C.Uint8(c.G), C.Uint8(c.B), C.Uint8(c.A)))
}

func FilledPie(re *sdl.Renderer, x, y, rad, start, end int, c sdl.Color) error {
	return ek(C.filledPieRGBA((*C.SDL_Renderer)(re), C.Sint16(x), C.Sint16(y), C.Sint16(rad), C.Sint16(start), C.Sint16(end), C.Uint8(c.R), C.Uint8(c.G), C.Uint8(c.B), C.Uint8(c.A)))
}

func Trigon(re *sdl.Renderer, x1, y1, x2, y2, x3, y3 int, c sdl.Color) error {
	return ek(C.trigonRGBA((*C.SDL_Renderer)(re), C.Sint16(x1), C.Sint16(y1), C.Sint16(x2), C.Sint16(y2), C.Sint16(x3), C.Sint16(y3), C.Uint8(c.R), C.Uint8(c.G), C.Uint8(c.B), C.Uint8(c.A)))
}

func AATrigon(re *sdl.Renderer, x1, y1, x2, y2, x3, y3 int, c sdl.Color) error {
	return ek(C.aatrigonRGBA((*C.SDL_Renderer)(re), C.Sint16(x1), C.Sint16(y1), C.Sint16(x2), C.Sint16(y2), C.Sint16(x3), C.Sint16(y3), C.Uint8(c.R), C.Uint8(c.G), C.Uint8(c.B), C.Uint8(c.A)))
}

func FilledTrigon(re *sdl.Renderer, x1, y1, x2, y2, x3, y3 int, c sdl.Color) error {
	return ek(C.filledTrigonRGBA((*C.SDL_Renderer)(re), C.Sint16(x1), C.Sint16(y1), C.Sint16(x2), C.Sint16(y2), C.Sint16(x3), C.Sint16(y3), C.Uint8(c.R), C.Uint8(c.G), C.Uint8(c.B), C.Uint8(c.A)))
}

func Polygon(re *sdl.Renderer, pts []sdl.Point, c sdl.Color) error {
	if len(pts) < 3 {
		return nil
	}
	return ek(C.goPolygonRGBA((*C.SDL_Renderer)(re), (*C.SDL_Point)(unsafe.Pointer(&pts[0])), C.int(len(pts)), C.Uint8(c.R), C.Uint8(c.G), C.Uint8(c.B), C.Uint8(c.A)))
}

func AAPolygon(re *sdl.Renderer, pts []sdl.Point, c sdl.Color) error {
	return ek(C.goAAPolygonRGBA((*C.SDL_Renderer)(re), (*C.SDL_Point)(unsafe.Pointer(&pts[0])), C.int(len(pts)), C.Uint8(c.R), C.Uint8(c.G), C.Uint8(c.B), C.Uint8(c.A)))
}

func FilledPolygon(re *sdl.Renderer, pts []sdl.Point, c sdl.Color) error {
	if len(pts) < 3 {
		return nil
	}
	return ek(C.goFilledPolygonRGBA((*C.SDL_Renderer)(re), (*C.SDL_Point)(unsafe.Pointer(&pts[0])), C.int(len(pts)), C.Uint8(c.R), C.Uint8(c.G), C.Uint8(c.B), C.Uint8(c.A)))
}

func TexturedPolygon(re *sdl.Renderer, pts []sdl.Point, texture *sdl.Surface, dx, dy int) error {
	if len(pts) < 3 {
		return nil
	}
	return ek(C.goTexturedPolygon((*C.SDL_Renderer)(re), (*C.SDL_Point)(unsafe.Pointer(&pts[0])), C.int(len(pts)), (*C.SDL_Surface)(unsafe.Pointer(&texture)), C.int(dx), C.int(dy)))
}

func Bezier(re *sdl.Renderer, pts []sdl.Point, s int, c sdl.Color) error {
	return ek(C.goBezierRGBA((*C.SDL_Renderer)(re), (*C.SDL_Point)(unsafe.Pointer(&pts[0])), C.int(len(pts)), C.int(s), C.Uint8(c.R), C.Uint8(c.G), C.Uint8(c.B), C.Uint8(c.A)))
}

func SetFont(font []byte, cw, ch uint32) {
	C.goGfxPrimitivesSetFont(unsafe.Pointer(&font[0]), C.Uint32(cw), C.Uint32(ch))
}

func SetFontRotation(rotation uint32) {
	C.goGfxPrimitivesSetFontRotation(C.Uint32(rotation))
}

func FontMetrics() (w, h, r int) {
	return int(C.goCharWidth), int(C.goCharHeight), int(C.goCharRotation)
}

func FontSize(str string) (w, h int) {
	curw := 0
	maxw := 0
	line := 1
	for _, ch := range str {
		if ch == '\n' {
			if maxw < curw {
				maxw = curw
			}
			curw = 0
			line++
		} else {
			curw++
		}
	}
	if maxw < curw {
		maxw = curw
	}
	w = maxw * int(C.goCharWidth)
	h = line * int(C.goCharHeight)
	return
}

func Character(re *sdl.Renderer, x, y int, c sdl.Color, r rune) error {
	return ek(C.characterRGBA((*C.SDL_Renderer)(re), C.Sint16(x), C.Sint16(y), C.char(r), C.Uint8(c.R), C.Uint8(c.G), C.Uint8(c.B), C.Uint8(c.A)))
}

func String(re *sdl.Renderer, x, y int, c sdl.Color, s string) error {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	return ek(C.stringRGBA((*C.SDL_Renderer)(re), C.Sint16(x), C.Sint16(y), cs, C.Uint8(c.R), C.Uint8(c.G), C.Uint8(c.B), C.Uint8(c.A)))
}

func ek(rc C.int) error {
	if rc < 0 {
		return errors.New("invalid parameter")
	}
	return nil
}

type FPSManager C.FPSmanager

func (m *FPSManager) Init() {
	C.SDL_initFramerate((*C.FPSmanager)(unsafe.Pointer(m)))
}

func (m *FPSManager) SetRate(rate uint32) error {
	return ek(C.SDL_setFramerate((*C.FPSmanager)(unsafe.Pointer(m)), C.Uint32(rate)))
}

func (m *FPSManager) Rate() (int, error) {
	rc := C.SDL_getFramerate((*C.FPSmanager)(unsafe.Pointer(m)))
	return int(rc), ek(rc)
}

func (m *FPSManager) Count() (int, error) {
	rc := C.SDL_getFramecount((*C.FPSmanager)(unsafe.Pointer(m)))
	return int(rc), ek(rc)
}

func (m *FPSManager) Delay() uint32 {
	return uint32(C.SDL_framerateDelay((*C.FPSmanager)(unsafe.Pointer(m))))
}
