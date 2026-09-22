package sdl

/*
#include "gosdl.h"
*/
import "C"
import "unsafe"

type (
	Joystick           C.SDL_Joystick
	JoystickID         C.SDL_JoystickID
	JoystickGUID       C.SDL_JoystickGUID
	JoystickPowerLevel C.SDL_JoystickPowerLevel
)

const (
	JOYSTICK_POWER_UNKNOWN JoystickPowerLevel = C.SDL_JOYSTICK_POWER_UNKNOWN
	JOYSTICK_POWER_EMPTY   JoystickPowerLevel = C.SDL_JOYSTICK_POWER_EMPTY
	JOYSTICK_POWER_LOW     JoystickPowerLevel = C.SDL_JOYSTICK_POWER_LOW
	JOYSTICK_POWER_MEDIUM  JoystickPowerLevel = C.SDL_JOYSTICK_POWER_MEDIUM
	JOYSTICK_POWER_FULL    JoystickPowerLevel = C.SDL_JOYSTICK_POWER_FULL
	JOYSTICK_POWER_WIRED   JoystickPowerLevel = C.SDL_JOYSTICK_POWER_WIRED
	JOYSTICK_POWER_MAX                        = C.SDL_JOYSTICK_POWER_MAX
)

func NumJoysticks() int {
	return int(C.SDL_NumJoysticks())
}

func JoystickNameForIndex(deviceIndex int) string {
	return C.GoString(C.SDL_JoystickNameForIndex(C.int(deviceIndex)))
}

func JoystickOpen(deviceIndex int) *Joystick {
	return (*Joystick)(C.SDL_JoystickOpen(C.int(deviceIndex)))
}

func JoystickFromInstanceID(joyid JoystickID) *Joystick {
	return (*Joystick)(C.SDL_JoystickFromInstanceID(C.SDL_JoystickID(joyid)))
}

func (j *Joystick) Name() string {
	return C.GoString(C.SDL_JoystickName((*C.SDL_Joystick)(j)))
}

func JoystickGetDeviceGUID(deviceIndex int) JoystickGUID {
	return JoystickGUID(C.SDL_JoystickGetDeviceGUID(C.int(deviceIndex)))
}

func (j *Joystick) GUID() JoystickGUID {
	return JoystickGUID(C.SDL_JoystickGetGUID((*C.SDL_Joystick)(j)))
}

func JoystickGetGUIDString(guid JoystickGUID) string {
	var pszGUID [64]byte
	C.SDL_JoystickGetGUIDString(C.SDL_JoystickGUID(guid), (*C.char)(unsafe.Pointer(&pszGUID[0])), C.int(len(pszGUID)))
	return C.GoString((*C.char)(unsafe.Pointer(&pszGUID[0])))
}

func JoystickGetGUIDFromString(pch string) JoystickGUID {
	cpch := C.CString(pch)
	defer C.free(unsafe.Pointer(cpch))
	return JoystickGUID(C.SDL_JoystickGetGUIDFromString(cpch))
}

func (j *Joystick) Attached() bool {
	return C.SDL_JoystickGetAttached((*C.SDL_Joystick)(j)) != 0
}

func (j *Joystick) ID() JoystickID {
	return JoystickID(C.SDL_JoystickInstanceID((*C.SDL_Joystick)(j)))
}

func (j *Joystick) NumAxes() int {
	return int(C.SDL_JoystickNumAxes((*C.SDL_Joystick)(j)))
}

func (j *Joystick) NumBalls() int {
	return int(C.SDL_JoystickNumBalls((*C.SDL_Joystick)(j)))
}

func (j *Joystick) NumHats() int {
	return int(C.SDL_JoystickNumHats((*C.SDL_Joystick)(j)))
}

func (j *Joystick) NumButtons() int {
	return int(C.SDL_JoystickNumButtons((*C.SDL_Joystick)(j)))
}

func JoystickUpdate() {
	C.SDL_JoystickUpdate()
}

func JoystickEventState(state int) int {
	return int(C.SDL_JoystickEventState(C.int(state)))
}

func (j *Joystick) Axis(axis int) int16 {
	return int16(C.SDL_JoystickGetAxis((*C.SDL_Joystick)(j), C.int(axis)))
}

const (
	HAT_CENTERED  = C.SDL_HAT_CENTERED
	HAT_UP        = C.SDL_HAT_UP
	HAT_RIGHT     = C.SDL_HAT_RIGHT
	HAT_DOWN      = C.SDL_HAT_DOWN
	HAT_LEFT      = C.SDL_HAT_LEFT
	HAT_RIGHTUP   = C.SDL_HAT_RIGHTUP
	HAT_RIGHTDOWN = C.SDL_HAT_RIGHTDOWN
	HAT_LEFTUP    = C.SDL_HAT_LEFTUP
	HAT_LEFTDOWN  = C.SDL_HAT_LEFTDOWN
)

func (j *Joystick) Hat(hat int) uint8 {
	return uint8(C.SDL_JoystickGetHat((*C.SDL_Joystick)(j), C.int(hat)))
}

func (j *Joystick) Ball(ball int) (dx, dy int, err error) {
	var cx, cy C.int
	rc := C.SDL_JoystickGetBall((*C.SDL_Joystick)(j), C.int(ball), &cx, &cy)
	return int(cx), int(cy), ek(rc)
}

func (j *Joystick) Button(button int) uint8 {
	return uint8(C.SDL_JoystickGetButton((*C.SDL_Joystick)(j), C.int(button)))
}

func (j *Joystick) Close() {
	C.SDL_JoystickClose((*C.SDL_Joystick)(j))
}

func (j *Joystick) PowerLevel() JoystickPowerLevel {
	return JoystickPowerLevel(C.SDL_JoystickCurrentPowerLevel((*C.SDL_Joystick)(j)))
}