package sdl

/*
#include "gosdl.h"
*/
import "C"

type (
	Renderer C.SDL_Renderer
)

func (r *Renderer) SetDrawColor(red, green, blue, alpha uint8) {
	C.SDL_SetRenderDrawColor((*C.SDL_Renderer)(r), C.Uint8(red), C.Uint8(green), C.Uint8(blue), C.Uint8(alpha))
}

func (r *Renderer) SetRenderScale(scalex, scaley float32) bool {
	return bool(C.SDL_SetRenderScale((*C.SDL_Renderer)(r), C.float(scalex), C.float(scaley)))
}

func (r *Renderer) Clear() bool {
	return bool(C.SDL_RenderClear((*C.SDL_Renderer)(r)))
}

func (r *Renderer) Flush() bool {
	return bool(C.SDL_FlushRenderer((*C.SDL_Renderer)(r)))
}

func (r *Renderer) Present() bool {
	return bool(C.SDL_RenderPresent((*C.SDL_Renderer)(r)))
}

func (r *Renderer) RenderPoint(x, y float32) bool {
	return bool(C.SDL_RenderPoint((*C.SDL_Renderer)(r), C.float(x), C.float(y)))
}

func (r *Renderer) RenderLine(x1, y1, x2, y2 float32) bool {
	return bool(C.SDL_RenderLine((*C.SDL_Renderer)(r), C.float(x1), C.float(y1), C.float(x2), C.float(y2)))
}

func (r *Renderer) SetRenderColorScale(scale float32) bool {
	return bool(C.SDL_SetRenderColorScale((*C.SDL_Renderer)(r), C.float(scale)))
}

func (r *Renderer) SetRenderVSync(vsync int) bool {
	return bool(C.SDL_SetRenderVSync((*C.SDL_Renderer)(r), C.int(vsync)))
}

func (r *Renderer) GetRenderColorScale() (float32, bool) {
	var scale C.float
	status := C.SDL_GetRenderColorScale((*C.SDL_Renderer)(r), &scale)
	return float32(scale), bool(status)
}
