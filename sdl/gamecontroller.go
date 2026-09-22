package sdl

/*
#include "gosdl.h"

static int addMappingsFromFile(const char *file) {
	return SDL_GameControllerAddMappingsFromFile(file);
}
*/
import "C"
import "unsafe"

type (
	GameController           C.SDL_GameController
	GameControllerBindType   C.SDL_GameControllerBindType
	GameControllerAxis       C.SDL_GameControllerAxis
	GameControllerButton     C.SDL_GameControllerButton
	GameControllerButtonBind C.SDL_GameControllerButtonBind
)

const (
	CONTROLLER_BINDTYPE_NONE   GameControllerBindType = C.SDL_CONTROLLER_BINDTYPE_NONE
	CONTROLLER_BINDTYPE_BUTTON GameControllerBindType = C.SDL_CONTROLLER_BINDTYPE_BUTTON
	CONTROLLER_BINDTYPE_AXIS   GameControllerBindType = C.SDL_CONTROLLER_BINDTYPE_AXIS
	CONTROLLER_BINDTYPE_HAT    GameControllerBindType = C.SDL_CONTROLLER_BINDTYPE_HAT
)

func GameControllerAddMappingsFromFile(file string) error {
	cfile := C.CString(file)
	defer C.free(unsafe.Pointer(cfile))
	return ek(C.addMappingsFromFile(cfile))
}

func GameControllerAddMapping(mapping string) error {
	cmapping := C.CString(mapping)
	defer C.free(unsafe.Pointer(cmapping))
	return ek(C.SDL_GameControllerAddMapping(cmapping))
}

func (guid JoystickGUID) Mapping() string {
	cmap := C.SDL_GameControllerMappingForGUID(C.SDL_JoystickGUID(guid))
	defer C.free(unsafe.Pointer(cmap))
	return C.GoString(cmap)
}

func (c *GameController) Mapping() string {
	cmap := C.SDL_GameControllerMapping((*C.SDL_GameController)(c))
	defer C.free(unsafe.Pointer(cmap))
	return C.GoString(cmap)
}

func IsGameController(joystickIndex int) bool {
	return C.SDL_IsGameController(C.int(joystickIndex)) != 0
}

func GameControllerNameForIndex(joystickIndex int) string {
	return C.GoString(C.SDL_GameControllerNameForIndex(C.int(joystickIndex)))
}

func GameControllerOpen(joystickIndex int) (*GameController, error) {
	c := C.SDL_GameControllerOpen(C.int(joystickIndex))
	if c == nil {
		return nil, GetError()
	}
	return (*GameController)(c), nil
}

func GameControllerFromInstanceID(joyid JoystickID) *GameController {
	return (*GameController)(C.SDL_GameControllerFromInstanceID(C.SDL_JoystickID(joyid)))
}

func (c *GameController) Name() string {
	return C.GoString(C.SDL_GameControllerName((*C.SDL_GameController)(c)))
}

func (c *GameController) Attached() bool {
	return C.SDL_GameControllerGetAttached((*C.SDL_GameController)(c)) != 0
}

func (c *GameController) Joystick() *Joystick {
	return (*Joystick)(C.SDL_GameControllerGetJoystick((*C.SDL_GameController)(c)))
}

func GameControllerEventState(state int) int {
	return int(C.SDL_GameControllerEventState(C.int(state)))
}

func GameControllerUpdate() {
	C.SDL_GameControllerUpdate()
}

const (
	CONTROLLER_AXIS_INVALID      GameControllerAxis = C.SDL_CONTROLLER_AXIS_INVALID
	CONTROLLER_AXIS_LEFTX        GameControllerAxis = C.SDL_CONTROLLER_AXIS_LEFTX
	CONTROLLER_AXIS_LEFTY        GameControllerAxis = C.SDL_CONTROLLER_AXIS_LEFTY
	CONTROLLER_AXIS_RIGHTX       GameControllerAxis = C.SDL_CONTROLLER_AXIS_RIGHTX
	CONTROLLER_AXIS_RIGHTY       GameControllerAxis = C.SDL_CONTROLLER_AXIS_RIGHTY
	CONTROLLER_AXIS_TRIGGERLEFT  GameControllerAxis = C.SDL_CONTROLLER_AXIS_TRIGGERLEFT
	CONTROLLER_AXIS_TRIGGERRIGHT GameControllerAxis = C.SDL_CONTROLLER_AXIS_TRIGGERRIGHT
	CONTROLLER_AXIS_MAX                             = C.SDL_CONTROLLER_AXIS_MAX
)

func GameControllerGetAxisFromString(pch string) GameControllerAxis {
	cpch := C.CString(pch)
	defer C.free(unsafe.Pointer(cpch))
	return GameControllerAxis(C.SDL_GameControllerGetAxisFromString(cpch))
}

func GameControllerGetStringForAxis(axis GameControllerAxis) string {
	return C.GoString(C.SDL_GameControllerGetStringForAxis(C.SDL_GameControllerAxis(axis)))
}

func (c *GameController) BindForAxis(axis GameControllerAxis) GameControllerButtonBind {
	return GameControllerButtonBind(C.SDL_GameControllerGetBindForAxis((*C.SDL_GameController)(c), C.SDL_GameControllerAxis(axis)))
}

func (c *GameController) Axis(axis GameControllerAxis) int16 {
	return int16(C.SDL_GameControllerGetAxis((*C.SDL_GameController)(c), C.SDL_GameControllerAxis(axis)))
}

const (
	CONTROLLER_BUTTON_INVALID       GameControllerButton = C.SDL_CONTROLLER_BUTTON_INVALID
	CONTROLLER_BUTTON_A             GameControllerButton = C.SDL_CONTROLLER_BUTTON_A
	CONTROLLER_BUTTON_B             GameControllerButton = C.SDL_CONTROLLER_BUTTON_B
	CONTROLLER_BUTTON_X             GameControllerButton = C.SDL_CONTROLLER_BUTTON_X
	CONTROLLER_BUTTON_Y             GameControllerButton = C.SDL_CONTROLLER_BUTTON_Y
	CONTROLLER_BUTTON_BACK          GameControllerButton = C.SDL_CONTROLLER_BUTTON_BACK
	CONTROLLER_BUTTON_GUIDE         GameControllerButton = C.SDL_CONTROLLER_BUTTON_GUIDE
	CONTROLLER_BUTTON_START         GameControllerButton = C.SDL_CONTROLLER_BUTTON_START
	CONTROLLER_BUTTON_LEFTSTICK     GameControllerButton = C.SDL_CONTROLLER_BUTTON_LEFTSTICK
	CONTROLLER_BUTTON_RIGHTSTICK    GameControllerButton = C.SDL_CONTROLLER_BUTTON_RIGHTSTICK
	CONTROLLER_BUTTON_LEFTSHOULDER  GameControllerButton = C.SDL_CONTROLLER_BUTTON_LEFTSHOULDER
	CONTROLLER_BUTTON_RIGHTSHOULDER GameControllerButton = C.SDL_CONTROLLER_BUTTON_RIGHTSHOULDER
	CONTROLLER_BUTTON_DPAD_UP       GameControllerButton = C.SDL_CONTROLLER_BUTTON_DPAD_UP
	CONTROLLER_BUTTON_DPAD_DOWN     GameControllerButton = C.SDL_CONTROLLER_BUTTON_DPAD_DOWN
	CONTROLLER_BUTTON_DPAD_LEFT     GameControllerButton = C.SDL_CONTROLLER_BUTTON_DPAD_LEFT
	CONTROLLER_BUTTON_DPAD_RIGHT    GameControllerButton = C.SDL_CONTROLLER_BUTTON_DPAD_RIGHT
	CONTROLLER_BUTTON_MAX                                = C.SDL_CONTROLLER_BUTTON_MAX
)

func GameControllerGetButtonFromString(pch string) GameControllerButton {
	cpch := C.CString(pch)
	defer C.free(unsafe.Pointer(cpch))
	return GameControllerButton(C.SDL_GameControllerGetButtonFromString(cpch))
}

func GameControllerGetStringForButton(button GameControllerButton) string {
	return C.GoString(C.SDL_GameControllerGetStringForButton(C.SDL_GameControllerButton(button)))
}

func (c *GameController) ButtonBind(button GameControllerButton) GameControllerButtonBind {
	return GameControllerButtonBind(C.SDL_GameControllerGetBindForButton((*C.SDL_GameController)(c), C.SDL_GameControllerButton(button)))
}

func (c *GameController) Button(button GameControllerButton) uint8 {
	return uint8(C.SDL_GameControllerGetButton((*C.SDL_GameController)(c), C.SDL_GameControllerButton(button)))
}

func (c *GameController) Close() {
	C.SDL_GameControllerClose((*C.SDL_GameController)(c))
}
