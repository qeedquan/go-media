package sdl

/*
#include "gosdl.h"
*/
import "C"

type (
	EventType C.SDL_EventType
	Event     any
)

const (
	QUIT EventType = C.SDL_EVENT_QUIT
)

type CommonEvent struct {
	Type      EventType
	Reserved  uint32
	Timestamp uint64
}

type QuitEvent struct {
	CommonEvent
}

func PollEvent() Event {
	return nil
}

func PumpEvents() {
	C.SDL_PumpEvents()
}
