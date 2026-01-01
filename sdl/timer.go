package sdl

/*
#include "gosdl.h"
*/
import "C"

type TimerID C.SDL_TimerID

func GetTicks() uint32 {
	return uint32(C.SDL_GetTicks())
}

func GetPerformanceCounter() uint64 {
	return uint64(C.SDL_GetPerformanceCounter())
}

func GetPerformanceFrequency() uint64 {
	return uint64(C.SDL_GetPerformanceFrequency())
}

func Delay(ms uint32) {
	C.SDL_Delay(C.Uint32(ms))
}

func RemoveTimer(id TimerID) bool {
	return C.SDL_RemoveTimer(C.SDL_TimerID(id)) != 0
}
