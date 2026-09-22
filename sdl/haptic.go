package sdl

/*
#include "gosdl.h"
*/
import "C"

type Haptic C.SDL_Haptic

const (
	HAPTIC_CONSTANT     = C.SDL_HAPTIC_CONSTANT
	HAPTIC_SINE         = C.SDL_HAPTIC_SINE
	HAPTIC_LEFTRIGHT    = C.SDL_HAPTIC_LEFTRIGHT
	HAPTIC_TRIANGLE     = C.SDL_HAPTIC_TRIANGLE
	HAPTIC_SAWTOOTHUP   = C.SDL_HAPTIC_SAWTOOTHUP
	HAPTIC_SAWTOOTHDOWN = C.SDL_HAPTIC_SAWTOOTHDOWN
)

func HapticOpen(deviceIndex int) (*Haptic, error) {
	h := C.SDL_HapticOpen(C.int(deviceIndex))
	if h == nil {
		return nil, GetError()
	}
	return (*Haptic)(h), nil
}

func HapticOpenFromMouse() (*Haptic, error) {
	h := C.SDL_HapticOpenFromMouse()
	if h == nil {
		return nil, GetError()
	}
	return (*Haptic)(h), nil
}

func (j *Joystick) IsHaptic() (bool, error) {
	rc := C.SDL_JoystickIsHaptic((*C.SDL_Joystick)(j))
	if rc < 0 {
		return false, GetError()
	}
	return rc != 0, nil
}

func HapticOpenFromJoystick(j *Joystick) (*Haptic, error) {
	h := C.SDL_HapticOpenFromJoystick((*C.SDL_Joystick)(j))
	if h == nil {
		return nil, GetError()
	}
	return (*Haptic)(h), nil
}

func (h *Haptic) NumEffects() (int, error) {
	rc := C.SDL_HapticNumEffects((*C.SDL_Haptic)(h))
	return int(rc), ek(rc)
}

func (h *Haptic) NumEffectsPlaying() (int, error) {
	rc := C.SDL_HapticNumEffectsPlaying((*C.SDL_Haptic)(h))
	return int(rc), ek(rc)
}

func (h *Haptic) Query() uint {
	return uint(C.SDL_HapticQuery((*C.SDL_Haptic)(h)))
}

func (h *Haptic) Close() {
	C.SDL_HapticClose((*C.SDL_Haptic)(h))
}

func (h *Haptic) Pause() error {
	return ek(C.SDL_HapticPause((*C.SDL_Haptic)(h)))
}

func (h *Haptic) Unpause() error {
	return ek(C.SDL_HapticUnpause((*C.SDL_Haptic)(h)))
}

func (h *Haptic) StopAll() error {
	return ek(C.SDL_HapticStopAll((*C.SDL_Haptic)(h)))
}

func (h *Haptic) RumbleInit() error {
	return ek(C.SDL_HapticRumbleInit((*C.SDL_Haptic)(h)))
}

func (h *Haptic) RumblePlay(strength float32, length uint32) error {
	return ek(C.SDL_HapticRumblePlay((*C.SDL_Haptic)(h), C.float(strength), C.Uint32(length)))
}

func (h *Haptic) RumbleStop() error {
	return ek(C.SDL_HapticRumbleStop((*C.SDL_Haptic)(h)))
}
