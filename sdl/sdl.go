package sdl

/*
#include "gosdl.h"
*/
import "C"

const (
	INIT_TIMER      = C.SDL_INIT_TIMER
	INIT_AUDIO      = C.SDL_INIT_AUDIO
	INIT_VIDEO      = C.SDL_INIT_VIDEO
	INIT_JOYSTICK   = C.SDL_INIT_JOYSTICK
	INIT_HAPTIC     = C.SDL_INIT_HAPTIC
	INIT_EVENTS     = C.SDL_INIT_EVENTS
	INIT_EVERYTHING = C.SDL_INIT_EVERYTHING
)

const (
	ENABLE  = C.SDL_ENABLE
	IGNORE  = C.SDL_IGNORE
	DISABLE = C.SDL_DISABLE
	QUERY   = C.SDL_QUERY
)

func Init(flags uint32) error {
	return ek(C.SDL_Init(C.Uint32(flags)))
}

func InitSubSystem(flags uint32) error {
	return ek(C.SDL_InitSubSystem(C.Uint32(flags)))
}

func WasInit(flags uint32) uint32 {
	return uint32(C.SDL_WasInit(C.Uint32(flags)))
}

func Quit() {
	C.SDL_Quit()
}
