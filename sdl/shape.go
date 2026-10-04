package sdl

/*
#include "gosdl.h"
*/
import "C"

const (
	NONSHAPEABLE_WINDOW    = C.SDL_NONSHAPEABLE_WINDOW
	INVALID_SHAPE_ARGUMENT = C.SDL_INVALID_SHAPE_ARGUMENT
	WINDOW_LACKS_SHAPE     = C.SDL_WINDOW_LACKS_SHAPE
)

func (w *Window) IsShaped() bool {
	return C.SDL_IsShapedWindow((*C.SDL_Window)(w)) != 0
}

type WindowShapeMode C.WindowShapeMode

const (
	ShapeModeDefault              WindowShapeMode = C.ShapeModeDefault
	ShapeModeBinarizeAlpha        WindowShapeMode = C.ShapeModeBinarizeAlpha
	ShapeModeReverseBinarizeAlpha WindowShapeMode = C.ShapeModeReverseBinarizeAlpha
	ShapeModeColorKey             WindowShapeMode = C.ShapeModeColorKey
)
