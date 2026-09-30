// +build ignore
//go: generate go tool cgo -godefs defs.go

package main

/*
#include "gosdl.h"
*/
import "C"

type (
	Window                C.SDL_Window
	Renderer              C.SDL_Renderer
	Color                 C.SDL_Color
	Rect                  C.SDL_Rect
	DisplayMode           C.SDL_DisplayMode
	Event                 C.SDL_Event
	CommonEvent           C.SDL_CommonEvent
	WindowEvent           C.SDL_WindowEvent
	KeyboardEvent         C.SDL_KeyboardEvent
	Keysym                C.SDL_Keysym
	Scancode              C.SDL_Scancode
	Haptic                C.SDL_Haptic
	Joystick              C.SDL_Joystick
	GameController        C.SDL_GameController
	Texture               C.SDL_Texture
	Surface               C.SDL_Surface
	PixelFormat           C.SDL_PixelFormat
	Palette               C.SDL_Palette
	BlitMap               C.struct_SDL_BlitMap
	AudioFormat           C.SDL_AudioFormat
	AudioSpec             C.SDL_AudioSpec
	MessageBoxColor       C.SDL_MessageBoxColor
	MessageBoxColorType   C.SDL_MessageBoxColorType
	MessageBoxColorScheme C.SDL_MessageBoxColorScheme
	MessageBoxData        C.SDL_MessageBoxData
)
