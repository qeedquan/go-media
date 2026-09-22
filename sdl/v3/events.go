package sdl

/*
#include "gosdl.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type (
	EventType C.SDL_EventType
	Event     any
)

type (
	KeyboardID C.SDL_KeyboardID
	Scancode   C.SDL_Scancode
	Keycode    C.SDL_Keycode
	Keymod     C.SDL_Keymod
)

const (
	QUIT EventType = C.SDL_EVENT_QUIT
)

type CommonEvent struct {
	Type      EventType
	Reserved  uint32
	Timestamp uint64
}

type KeyboardEvent struct {
	CommonEvent
	WindowID WindowID
	Which    KeyboardID
	Scancode Scancode
	Key      Keycode
	Mod      Keymod
	Raw      uint16
	Down     bool
	Repeat   bool
}

type (
	KeyDownEvent KeyboardEvent
	KeyUpEvent   KeyboardEvent
)

type QuitEvent struct {
	CommonEvent
}

func PollEvent() Event {
	var ev C.SDL_Event
	if C.SDL_PollEvent(&ev) {
		return nil
	}
	return evCommon(ev)
}

func PumpEvents() {
	C.SDL_PumpEvents()
}

func evCommon(ev C.SDL_Event) Event {
	cev := (*C.SDL_CommonEvent)(unsafe.Pointer(&ev))
	switch cev._type {
	case C.SDL_EVENT_QUIT:
		pev := (*C.SDL_QuitEvent)(unsafe.Pointer(&ev))
		return QuitEvent{
			CommonEvent{EventType(pev._type), uint32(pev.reserved), uint64(pev.timestamp)},
		}

	case C.SDL_EVENT_KEY_DOWN:
		pev := (*C.SDL_KeyboardEvent)(unsafe.Pointer(&ev))
		return KeyDownEvent(evKey(pev))

	case C.SDL_EVENT_KEY_UP:
		pev := (*C.SDL_KeyboardEvent)(unsafe.Pointer(&ev))
		return KeyUpEvent(evKey(pev))

	case C.SDL_EVENT_WINDOW_RESIZED:

	default:
		fmt.Printf("unknown event type %#x\n", cev._type)
	}
	return nil
}

func evKey(pev *C.SDL_KeyboardEvent) KeyboardEvent {
	return KeyboardEvent{
		CommonEvent{EventType(pev._type), uint32(pev.reserved), uint64(pev.timestamp)},
		WindowID(pev.windowID),
		KeyboardID(pev.which),
		Scancode(pev.scancode),
		Keycode(pev.key),
		Keymod(pev.mod),
		uint16(pev.raw),
		pev.down != false,
		pev.repeat != false,
	}
}
