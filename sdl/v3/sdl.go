package sdl

/*
#include "gosdl.h"
*/
import "C"
import "errors"

const (
	INIT_AUDIO    = C.SDL_INIT_AUDIO
	INIT_VIDEO    = C.SDL_INIT_VIDEO
	INIT_JOYSTICK = C.SDL_INIT_JOYSTICK
	INIT_HAPTIC   = C.SDL_INIT_HAPTIC
	INIT_GAMEPAD  = C.SDL_INIT_GAMEPAD
	INIT_EVENTS   = C.SDL_INIT_EVENTS
	INIT_SENSOR   = C.SDL_INIT_SENSOR
	INIT_CAMERA   = C.SDL_INIT_CAMERA
)

func Init(flags uint32) error {
	if !C.SDL_Init(C.Uint32(flags)) {
		return errors.New(GetError())
	}
	return nil
}

func Quit() {
	C.SDL_Quit()
}

func GetError() string {
	return C.GoString(C.SDL_GetError())
}
