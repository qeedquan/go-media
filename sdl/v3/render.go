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

func (r *Renderer) Clear() {
	C.SDL_RenderClear((*C.SDL_Renderer)(r))
}

func (r *Renderer) Present() {
	C.SDL_RenderPresent((*C.SDL_Renderer)(r))
}
