package sdl

/*
#include "gosdl.h"
*/
import "C"

type (
	TouchID  C.SDL_TouchID
	FingerID C.SDL_FingerID
)

type Finger struct {
	ID       FingerID
	X, Y     float32
	Pressure float32
}

func GetNumTouchDevices() int {
	return int(C.SDL_GetNumTouchDevices())
}

func GetTouchDevice(index int) TouchID {
	return TouchID(C.SDL_GetTouchDevice(C.int(index)))
}

func GetNumTouchFingers(touchID TouchID) int {
	return int(C.SDL_GetNumTouchFingers(C.SDL_TouchID(touchID)))
}

func GetTouchFinger(touchID TouchID, index int) *Finger {
	f := (C.SDL_GetTouchFinger(C.SDL_TouchID(touchID), C.int(index)))
	if f == nil {
		return nil
	}
	return &Finger{
		FingerID(f.id),
		float32(f.x),
		float32(f.y),
		float32(f.pressure),
	}
}
