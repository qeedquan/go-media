package sdl

/*
#include <SDL.h>
*/
import "C"

import (
	"image"
	"unsafe"
)

type Point struct {
	X, Y int32
}

type Rect struct {
	X, Y, W, H int32
}

type FPoint struct {
	X, Y float32
}

type FRect struct {
	X, Y, W, H float32
}

func (p Point) Add(q Point) Point {
	return Point{p.X + q.X, p.Y + q.Y}
}

func (p Point) Sub(q Point) Point {
	return Point{p.X - q.X, p.Y - q.Y}
}

func (p Point) In(r Rect) bool {
	return C.SDL_PointInRect((*C.SDL_Point)(unsafe.Pointer(&p)), (*C.SDL_Rect)(unsafe.Pointer(&r))) != 0
}

func (r Rect) Empty() bool {
	return C.SDL_RectEmpty((*C.SDL_Rect)(unsafe.Pointer(&r))) != 0
}

func (r Rect) Equal(p Rect) bool {
	return C.SDL_RectEquals((*C.SDL_Rect)(unsafe.Pointer(&r)), (*C.SDL_Rect)(unsafe.Pointer(&p))) != 0
}

func (r Rect) Enclose(p []Point) Rect {
	var res C.SDL_Rect
	C.SDL_EnclosePoints((*C.SDL_Point)(unsafe.Pointer(&p[0])), C.int(len(p)), (*C.SDL_Rect)(unsafe.Pointer(&r)), &res)
	return Rect{int32(res.x), int32(res.y), int32(res.w), int32(res.h)}
}

func (r Rect) Intersect(s Rect) Rect {
	a := image.Rect(int(r.X), int(r.Y), int(r.X+r.W), int(r.Y+r.H))
	b := image.Rect(int(s.X), int(s.Y), int(s.X+s.W), int(s.Y+s.H))
	c := a.Intersect(b)
	return Rect{
		X: int32(c.Min.X),
		Y: int32(c.Min.Y),
		W: int32(c.Dx()),
		H: int32(c.Dy()),
	}
}

func (r Rect) Collide(s Rect) bool {
	p := r.Intersect(s)
	return p.W != 0 && p.H != 0
}

func (r Rect) CenterX() int {
	return int(r.X + r.W/2)
}

func (r Rect) CenterY() int {
	return int(r.Y + r.H/2)
}

func (r Rect) Int() image.Rectangle {
	return image.Rect(int(r.X), int(r.Y), int(r.X+r.W), int(r.Y+r.H))
}

func Recti(r image.Rectangle) Rect {
	r = r.Canon()
	return Rect{
		int32(r.Min.X),
		int32(r.Min.Y),
		int32(r.Dx()),
		int32(r.Dy()),
	}
}
