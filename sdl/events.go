package sdl

/*
#include "gosdl.h"
*/
import "C"

import (
	"fmt"
	"unsafe"
)

const (
	RELEASED = C.SDL_RELEASED
	PRESSED  = C.SDL_PRESSED
)

type (
	EventType C.SDL_EventType
	Event     interface{}
)

const (
	FIRSTEVENT      EventType = C.SDL_FIRSTEVENT
	QUIT            EventType = C.SDL_QUIT
	APP_TERMINATING EventType = C.SDL_APP_TERMINATING

	APP_LOWMEMORY           EventType = C.SDL_APP_LOWMEMORY
	APP_WILLENTERBACKGROUND EventType = C.SDL_APP_WILLENTERBACKGROUND
	APP_DIDENTERBACKGROUND  EventType = C.SDL_APP_DIDENTERBACKGROUND
	APP_WILLENTERFOREGROUND EventType = C.SDL_APP_WILLENTERFOREGROUND
	APP_DIDENTERFOREGROUND  EventType = C.SDL_APP_DIDENTERFOREGROUND

	WINDOWEVENT EventType = C.SDL_WINDOWEVENT
	SYSWMEVENT  EventType = C.SDL_SYSWMEVENT

	KEYDOWN       EventType = C.SDL_KEYDOWN
	KEYUP         EventType = C.SDL_KEYUP
	TEXTEDITING   EventType = C.SDL_TEXTEDITING
	TEXTINPUT     EventType = C.SDL_TEXTINPUT
	KEYMAPCHANGED EventType = C.SDL_KEYMAPCHANGED

	MOUSEMOTION     EventType = C.SDL_MOUSEMOTION
	MOUSEBUTTONDOWN EventType = C.SDL_MOUSEBUTTONDOWN
	MOUSEBUTTONUP   EventType = C.SDL_MOUSEBUTTONUP
	MOUSEWHEEL      EventType = C.SDL_MOUSEWHEEL

	JOYAXISMOTION    EventType = C.SDL_JOYAXISMOTION
	JOYBALLMOTION    EventType = C.SDL_JOYBALLMOTION
	JOYHATMOTION     EventType = C.SDL_JOYHATMOTION
	JOYBUTTONDOWN    EventType = C.SDL_JOYBUTTONDOWN
	JOYBUTTONUP      EventType = C.SDL_JOYBUTTONUP
	JOYDEVICEADDED   EventType = C.SDL_JOYDEVICEADDED
	JOYDEVICEREMOVED EventType = C.SDL_JOYDEVICEREMOVED

	CONTROLLERAXISMOTION     EventType = C.SDL_CONTROLLERAXISMOTION
	CONTROLLERBUTTONDOWN     EventType = C.SDL_CONTROLLERBUTTONDOWN
	CONTROLLERBUTTONUP       EventType = C.SDL_CONTROLLERBUTTONUP
	CONTROLLERDEVICEADDED    EventType = C.SDL_CONTROLLERDEVICEADDED
	CONTROLLERDEVICEREMOVED  EventType = C.SDL_CONTROLLERDEVICEREMOVED
	CONTROLLERDEVICEREMAPPED EventType = C.SDL_CONTROLLERDEVICEREMAPPED

	FINGERDOWN   EventType = C.SDL_FINGERDOWN
	FINGERUP     EventType = C.SDL_FINGERUP
	FINGERMOTION EventType = C.SDL_FINGERMOTION

	DOLLARGESTURE EventType = C.SDL_DOLLARGESTURE
	DOLLARRECORD  EventType = C.SDL_DOLLARRECORD
	MULTIGESTURE  EventType = C.SDL_MULTIGESTURE

	CLIPBOARDUPDATE EventType = C.SDL_CLIPBOARDUPDATE

	DROPFILE EventType = C.SDL_DROPFILE

	AUDIODEVICEADDED   EventType = C.SDL_AUDIODEVICEADDED
	AUDIODEVICEREMOVED EventType = C.SDL_AUDIODEVICEREMOVED

	RENDER_TARGETS_RESET EventType = C.SDL_RENDER_TARGETS_RESET
	RENDER_DEVICE_RESET  EventType = C.SDL_RENDER_DEVICE_RESET

	USEREVENT EventType = C.SDL_USEREVENT

	LASTEVENT EventType = C.SDL_LASTEVENT
)

type CommonEvent struct {
	Type      EventType
	Timestamp uint32
}

type WindowEvent struct {
	CommonEvent
	WindowID uint32
	Event    WindowEventID
	Data     [2]int32
}

type KeyboardEvent struct {
	CommonEvent
	WindowID uint32
	State    uint8
	Repeat   bool
	Keysym
}

type (
	KeyDownEvent KeyboardEvent
	KeyUpEvent   KeyboardEvent
)

type TextEditingEvent struct {
	CommonEvent
	WindowID uint32
	Text     string
	Start    int32
	Length   int32
}

type TextInputEvent struct {
	CommonEvent
	WindowID uint32
	Text     string
}

type MouseMotionEvent struct {
	CommonEvent
	WindowID   uint32
	Which      uint32
	State      uint32
	X, Y       int32
	Xrel, Yrel int32
}

type MouseButtonEvent struct {
	CommonEvent
	WindowID uint32
	Which    uint32
	Button   uint8
	State    uint8
	Clicks   uint8
	X, Y     int32
}

type MouseButtonDownEvent MouseButtonEvent
type MouseButtonUpEvent MouseButtonEvent

type MouseWheelEvent struct {
	CommonEvent
	WindowID  uint32
	Which     uint32
	X, Y      int32
	Direction uint32
}

type JoyAxisEvent struct {
	CommonEvent
	Which JoystickID
	Axis  uint8
	Value int16
}

type JoyBallEvent struct {
	CommonEvent
	Which JoystickID
	Ball  uint8
	Xrel  int16
	Yrel  int16
}

type JoyButtonEvent struct {
	CommonEvent
	Which  JoystickID
	Button uint8
	State  uint8
}

type (
	JoyButtonDownEvent JoyButtonEvent
	JoyButtonUpEvent   JoyButtonEvent
)

type JoyHatEvent struct {
	CommonEvent
	Which JoystickID
	Hat   uint8
	Value uint8
}

type JoyDeviceEvent struct {
	CommonEvent
	Which JoystickID
}

type (
	JoyDeviceAddedEvent   JoyDeviceEvent
	JoyDeviceRemovedEvent JoyDeviceEvent
)

type ControllerAxisEvent struct {
	CommonEvent
	Which JoystickID
	Axis  uint8
	Value int16
}

type ControllerButtonEvent struct {
	CommonEvent
	Which  JoystickID
	Button uint8
	State  uint8
}

type (
	ControllerButtonDownEvent ControllerButtonEvent
	ControllerButtonUpEvent   ControllerButtonEvent
)

type ControllerDeviceEvent struct {
	CommonEvent
	Which JoystickID
}

type (
	ControllerDeviceAddedEvent   ControllerDeviceEvent
	ControllerDeviceRemovedEvent ControllerDeviceEvent
)

type QuitEvent struct {
	CommonEvent
}

type AudioDeviceEvent struct {
	CommonEvent
	Which     uint32
	IsCapture bool
}

type (
	AudioDeviceAddedEvent   AudioDeviceEvent
	AudioDeviceRemovedEvent AudioDeviceEvent
)

type TouchFingerEvent struct {
	CommonEvent
	TouchID                TouchID
	FingerID               FingerID
	X, Y, Dx, Dy, Pressure float32
}

type (
	FingerMotionEvent TouchFingerEvent
	FingerDownEvent   TouchFingerEvent
	FingerUpEvent     TouchFingerEvent
)

type MultiGestureEvent struct {
	CommonEvent
	TouchID    TouchID
	Theta      float32
	Distance   float32
	X, Y       float32
	NumFingers uint16
}

type DollarGestureEvent struct {
	CommonEvent
	TouchID    TouchID
	GestureID  GestureID
	NumFingers uint32
	Error      float32
	X, Y       float32
}

type DollarRecordEvent DollarGestureEvent

type DropEvent struct {
	CommonEvent
	File string
}

type OSEvent struct {
	CommonEvent
}

type UserEvent struct {
	CommonEvent
	WindowID uint32
	Code     int32
}

type SysWMEvent struct {
	CommonEvent
}

type ClipboardEvent struct {
	CommonEvent
	Text string
}

func PollEvent() Event {
	var ev C.SDL_Event
	if C.SDL_PollEvent(&ev) == 0 {
		return nil
	}
	return evCommon(ev)
}

func WaitEvent() Event {
	var ev C.SDL_Event
	if C.SDL_WaitEvent(&ev) == 0 {
		return nil
	}
	return evCommon(ev)
}

func WaitEventTimeout(timeout int) Event {
	var ev C.SDL_Event
	if C.SDL_WaitEventTimeout(&ev, C.int(timeout)) == 0 {
		return nil
	}
	return evCommon(ev)
}

func evCommon(ev C.SDL_Event) Event {
	cev := (*C.SDL_CommonEvent)(unsafe.Pointer(&ev))
	switch cev._type {
	case C.SDL_APP_TERMINATING,
		C.SDL_APP_LOWMEMORY,
		C.SDL_APP_WILLENTERBACKGROUND,
		C.SDL_APP_DIDENTERBACKGROUND,
		C.SDL_APP_WILLENTERFOREGROUND,
		C.SDL_APP_DIDENTERFOREGROUND,
		C.SDL_KEYMAPCHANGED,
		C.SDL_RENDER_TARGETS_RESET,
		C.SDL_RENDER_DEVICE_RESET:
		return CommonEvent{EventType(cev._type), uint32(cev.timestamp)}

	case C.SDL_QUIT:
		pev := (*C.SDL_QuitEvent)(unsafe.Pointer(&ev))
		return QuitEvent{
			CommonEvent{EventType(pev._type), uint32(pev.timestamp)},
		}

	case C.SDL_WINDOWEVENT:
		pev := (*C.SDL_WindowEvent)(unsafe.Pointer(&ev))
		return WindowEvent{
			CommonEvent{EventType(pev._type), uint32(pev.timestamp)},
			uint32(pev.windowID),
			WindowEventID(pev.event),
			[2]int32{int32(pev.data1), int32(pev.data2)},
		}

	case C.SDL_KEYDOWN:
		pev := (*C.SDL_KeyboardEvent)(unsafe.Pointer(&ev))
		return KeyDownEvent(evKey(pev))

	case C.SDL_KEYUP:
		pev := (*C.SDL_KeyboardEvent)(unsafe.Pointer(&ev))
		return KeyUpEvent(evKey(pev))

	case C.SDL_TEXTEDITING:
		pev := (*C.SDL_TextEditingEvent)(unsafe.Pointer(&ev))
		return TextEditingEvent{
			CommonEvent{EventType(pev._type), uint32(pev.timestamp)},
			uint32(pev.windowID),
			C.GoString(&pev.text[0]),
			int32(pev.start),
			int32(pev.length),
		}

	case C.SDL_TEXTINPUT:
		pev := (*C.SDL_TextInputEvent)(unsafe.Pointer(&ev))
		return TextInputEvent{
			CommonEvent{EventType(pev._type), uint32(pev.timestamp)},
			uint32(pev.windowID),
			C.GoString(&pev.text[0]),
		}

	case C.SDL_MOUSEMOTION:
		pev := (*C.SDL_MouseMotionEvent)(unsafe.Pointer(&ev))
		return MouseMotionEvent{
			CommonEvent{EventType(pev._type), uint32(pev.timestamp)},
			uint32(pev.windowID),
			uint32(pev.which),
			uint32(pev.state),
			int32(pev.x), int32(pev.y),
			int32(pev.xrel), int32(pev.yrel),
		}

	case C.SDL_MOUSEBUTTONDOWN:
		pev := (*C.SDL_MouseButtonEvent)(unsafe.Pointer(&ev))
		return MouseButtonDownEvent(evMouseButton(pev))

	case C.SDL_MOUSEBUTTONUP:
		pev := (*C.SDL_MouseButtonEvent)(unsafe.Pointer(&ev))
		return MouseButtonUpEvent(evMouseButton(pev))

	case C.SDL_MOUSEWHEEL:
		pev := (*C.SDL_MouseWheelEvent)(unsafe.Pointer(&ev))
		return MouseWheelEvent{
			CommonEvent{EventType(pev._type), uint32(pev.timestamp)},
			uint32(pev.windowID),
			uint32(pev.which),
			int32(pev.x),
			int32(pev.y),
			uint32(pev.direction),
		}

	case C.SDL_JOYAXISMOTION:
		pev := (*C.SDL_JoyAxisEvent)(unsafe.Pointer(&ev))
		return JoyAxisEvent{
			CommonEvent{EventType(pev._type), uint32(pev.timestamp)},
			JoystickID(pev.which),
			uint8(pev.axis),
			int16(pev.value),
		}

	case C.SDL_JOYBALLMOTION:
		pev := (*C.SDL_JoyBallEvent)(unsafe.Pointer(&ev))
		return JoyBallEvent{
			CommonEvent{EventType(pev._type), uint32(pev.timestamp)},
			JoystickID(pev.which),
			uint8(pev.ball),
			int16(pev.xrel),
			int16(pev.yrel),
		}

	case C.SDL_JOYHATMOTION:
		pev := (*C.SDL_JoyHatEvent)(unsafe.Pointer(&ev))
		return JoyHatEvent{
			CommonEvent{EventType(pev._type), uint32(pev.timestamp)},
			JoystickID(pev.which),
			uint8(pev.hat),
			uint8(pev.value),
		}

	case C.SDL_JOYBUTTONDOWN:
		pev := (*C.SDL_JoyButtonEvent)(unsafe.Pointer(&ev))
		return JoyButtonDownEvent(evJoyButton(pev))

	case C.SDL_JOYBUTTONUP:
		pev := (*C.SDL_JoyButtonEvent)(unsafe.Pointer(&ev))
		return JoyButtonUpEvent(evJoyButton(pev))

	case C.SDL_JOYDEVICEADDED:
		pev := (*C.SDL_JoyDeviceEvent)(unsafe.Pointer(&ev))
		return JoyDeviceAddedEvent(evJoyDevice(pev))

	case C.SDL_JOYDEVICEREMOVED:
		pev := (*C.SDL_JoyDeviceEvent)(unsafe.Pointer(&ev))
		return JoyDeviceRemovedEvent(evJoyDevice(pev))

	case C.SDL_CONTROLLERAXISMOTION:
		pev := (*C.SDL_ControllerAxisEvent)(unsafe.Pointer(&ev))
		return ControllerAxisEvent{
			CommonEvent{EventType(pev._type), uint32(pev.timestamp)},
			JoystickID(pev.which),
			uint8(pev.axis),
			int16(pev.value),
		}

	case C.SDL_CONTROLLERBUTTONDOWN:
		pev := (*C.SDL_ControllerButtonEvent)(unsafe.Pointer(&ev))
		return ControllerButtonDownEvent(evControllerButton(pev))

	case C.SDL_CONTROLLERBUTTONUP:
		pev := (*C.SDL_ControllerButtonEvent)(unsafe.Pointer(&ev))
		return ControllerButtonUpEvent(evControllerButton(pev))

	case C.SDL_CONTROLLERDEVICEADDED:
		pev := (*C.SDL_ControllerDeviceEvent)(unsafe.Pointer(&ev))
		return ControllerDeviceAddedEvent(evControllerDevice(pev))

	case C.SDL_CONTROLLERDEVICEREMOVED:
		pev := (*C.SDL_ControllerDeviceEvent)(unsafe.Pointer(&ev))
		return ControllerDeviceRemovedEvent(evControllerDevice(pev))

	case C.SDL_AUDIODEVICEADDED:
		pev := (*C.SDL_AudioDeviceEvent)(unsafe.Pointer(&ev))
		return AudioDeviceAddedEvent(evAudioDevice(pev))

	case C.SDL_AUDIODEVICEREMOVED:
		pev := (*C.SDL_AudioDeviceEvent)(unsafe.Pointer(&ev))
		return AudioDeviceRemovedEvent(evAudioDevice(pev))

	case C.SDL_FINGERMOTION:
		pev := (*C.SDL_TouchFingerEvent)(unsafe.Pointer(&ev))
		return FingerMotionEvent(evTouchFinger(pev))

	case C.SDL_FINGERDOWN:
		pev := (*C.SDL_TouchFingerEvent)(unsafe.Pointer(&ev))
		return FingerDownEvent(evTouchFinger(pev))

	case C.SDL_FINGERUP:
		pev := (*C.SDL_TouchFingerEvent)(unsafe.Pointer(&ev))
		return FingerUpEvent(evTouchFinger(pev))

	case C.SDL_MULTIGESTURE:
		pev := (*C.SDL_MultiGestureEvent)(unsafe.Pointer(&ev))
		return MultiGestureEvent{
			CommonEvent{EventType(pev._type), uint32(pev.timestamp)},
			TouchID(pev.touchId),
			float32(pev.dTheta),
			float32(pev.dDist),
			float32(pev.x),
			float32(pev.y),
			uint16(pev.numFingers),
		}

	case C.SDL_DOLLARGESTURE:
		pev := (*C.SDL_DollarGestureEvent)(unsafe.Pointer(&ev))
		return evDollarGesture(pev)

	case C.SDL_DOLLARRECORD:
		pev := (*C.SDL_DollarGestureEvent)(unsafe.Pointer(&ev))
		return DollarRecordEvent(evDollarGesture(pev))

	case C.SDL_DROPFILE:
		pev := (*C.SDL_DropEvent)(unsafe.Pointer(&ev))
		defer C.free(unsafe.Pointer(pev.file))
		return DropEvent{
			CommonEvent{EventType(pev._type), uint32(pev.timestamp)},
			C.GoString(pev.file),
		}

	case C.SDL_CLIPBOARDUPDATE:
		pev := (*C.SDL_CommonEvent)(unsafe.Pointer(&ev))
		return ClipboardEvent{
			CommonEvent{EventType(pev._type), uint32(pev.timestamp)},
			GetClipboardText(),
		}

	default:
		if C.SDL_USEREVENT <= cev._type && cev._type <= C.SDL_LASTEVENT-1 {
			pev := (*C.SDL_UserEvent)(unsafe.Pointer(&ev))
			return UserEvent{
				CommonEvent{EventType(pev._type), uint32(pev.timestamp)},
				uint32(pev.windowID),
				int32(pev.code),
			}
		}

		fmt.Printf("unknown event type %#x\n", cev._type)
		return nil
	}
}

func evKey(pev *C.SDL_KeyboardEvent) KeyboardEvent {
	return KeyboardEvent{
		CommonEvent{EventType(pev._type), uint32(pev.timestamp)},
		uint32(pev.windowID),
		uint8(pev.state),
		pev.repeat != 0,
		Keysym{
			Scancode(pev.keysym.scancode),
			Keycode(pev.keysym.sym),
			uint16(pev.keysym.mod),
		},
	}
}

func evMouseButton(pev *C.SDL_MouseButtonEvent) MouseButtonEvent {
	return MouseButtonEvent{
		CommonEvent{EventType(pev._type), uint32(pev.timestamp)},
		uint32(pev.windowID),
		uint32(pev.which),
		uint8(pev.button),
		uint8(pev.state),
		uint8(pev.clicks),
		int32(pev.x),
		int32(pev.y),
	}
}

func evJoyButton(pev *C.SDL_JoyButtonEvent) JoyButtonEvent {
	return JoyButtonEvent{
		CommonEvent{EventType(pev._type), uint32(pev.timestamp)},
		JoystickID(pev.which),
		uint8(pev.button),
		uint8(pev.state),
	}
}

func evJoyDevice(pev *C.SDL_JoyDeviceEvent) JoyDeviceEvent {
	return JoyDeviceEvent{
		CommonEvent{EventType(pev._type), uint32(pev.timestamp)},
		JoystickID(pev.which),
	}
}

func evControllerDevice(pev *C.SDL_ControllerDeviceEvent) ControllerDeviceEvent {
	return ControllerDeviceEvent{
		CommonEvent{EventType(pev._type), uint32(pev.timestamp)},
		JoystickID(pev.which),
	}
}

func evControllerButton(pev *C.SDL_ControllerButtonEvent) ControllerButtonEvent {
	return ControllerButtonEvent{
		CommonEvent{EventType(pev._type), uint32(pev.timestamp)},
		JoystickID(pev.which),
		uint8(pev.button),
		uint8(pev.state),
	}
}

func evAudioDevice(pev *C.SDL_AudioDeviceEvent) AudioDeviceEvent {
	return AudioDeviceEvent{
		CommonEvent{EventType(pev._type), uint32(pev.timestamp)},
		uint32(pev.which),
		pev.iscapture != 0,
	}
}

func evTouchFinger(pev *C.SDL_TouchFingerEvent) TouchFingerEvent {
	return TouchFingerEvent{
		CommonEvent{EventType(pev._type), uint32(pev.timestamp)},
		TouchID(pev.touchId),
		FingerID(pev.fingerId),
		float32(pev.x),
		float32(pev.y),
		float32(pev.dx),
		float32(pev.dy),
		float32(pev.pressure),
	}
}

func evDollarGesture(pev *C.SDL_DollarGestureEvent) DollarGestureEvent {
	return DollarGestureEvent{
		CommonEvent{EventType(pev._type), uint32(pev.timestamp)},
		TouchID(pev.touchId),
		GestureID(pev.gestureId),
		uint32(pev.numFingers),
		float32(pev.error),
		float32(pev.x),
		float32(pev.y),
	}
}

func PumpEvents() {
	C.SDL_PumpEvents()
}

func EventState(typ EventType, state int) uint8 {
	return uint8(C.SDL_EventState(C.Uint32(typ), C.int(state)))
}

func RegisterEvents(numEvents int) uint32 {
	return uint32(C.SDL_RegisterEvents(C.int(numEvents)))
}
