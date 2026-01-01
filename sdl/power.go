package sdl

/*
#include "gosdl.h"
*/
import "C"

type PowerState C.SDL_PowerState

const (
	POWERSTATE_UNKNOWN    PowerState = C.SDL_POWERSTATE_UNKNOWN
	POWERSTATE_ON_BATTERY PowerState = C.SDL_POWERSTATE_ON_BATTERY
	POWERSTATE_NO_BATTERY PowerState = C.SDL_POWERSTATE_NO_BATTERY
	POWERSTATE_CHARGING   PowerState = C.SDL_POWERSTATE_CHARGING
	POWERSTATE_CHARGED    PowerState = C.SDL_POWERSTATE_CHARGED
)

func GetPowerInfo() (state PowerState, secs, pct int) {
	var csecs, cpct C.int
	cstate := C.SDL_GetPowerInfo(&csecs, &cpct)
	return PowerState(cstate), int(csecs), int(cpct)
}
