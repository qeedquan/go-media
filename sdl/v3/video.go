package sdl

/*
#include "gosdl.h"
*/
import "C"
import "unsafe"

type (
	Window      C.SDL_Window
	WindowFlags = C.SDL_WindowFlags
	WindowID    = C.SDL_WindowID
)

const (
	WINDOW_FULLSCREEN          WindowFlags = C.SDL_WINDOW_FULLSCREEN
	WINDOW_OPENGL              WindowFlags = C.SDL_WINDOW_OPENGL
	WINDOW_OCCLUDED            WindowFlags = C.SDL_WINDOW_OCCLUDED
	WINDOW_HIDDEN              WindowFlags = C.SDL_WINDOW_HIDDEN
	WINDOW_BORDERLESS          WindowFlags = C.SDL_WINDOW_BORDERLESS
	WINDOW_RESIZABLE           WindowFlags = C.SDL_WINDOW_RESIZABLE
	WINDOW_MINIMIZED           WindowFlags = C.SDL_WINDOW_MINIMIZED
	WINDOW_MAXIMIZED           WindowFlags = C.SDL_WINDOW_MAXIMIZED
	WINDOW_MOUSE_GRABBED       WindowFlags = C.SDL_WINDOW_MOUSE_GRABBED
	WINDOW_INPUT_FOCUS         WindowFlags = C.SDL_WINDOW_INPUT_FOCUS
	WINDOW_MOUSE_FOCUS         WindowFlags = C.SDL_WINDOW_MOUSE_FOCUS
	WINDOW_EXTERNAL            WindowFlags = C.SDL_WINDOW_EXTERNAL
	WINDOW_MODAL               WindowFlags = C.SDL_WINDOW_MODAL
	WINDOW_HIGH_PIXEL_DENSITY  WindowFlags = C.SDL_WINDOW_HIGH_PIXEL_DENSITY
	WINDOW_MOUSE_CAPTURE       WindowFlags = C.SDL_WINDOW_MOUSE_CAPTURE
	WINDOW_MOUSE_RELATIVE_MODE WindowFlags = C.SDL_WINDOW_MOUSE_RELATIVE_MODE
	WINDOW_ALWAYS_ON_TOP       WindowFlags = C.SDL_WINDOW_ALWAYS_ON_TOP
	WINDOW_UTILITY             WindowFlags = C.SDL_WINDOW_UTILITY
	WINDOW_TOOLTIP             WindowFlags = C.SDL_WINDOW_TOOLTIP
	WINDOW_POPUP_MENU          WindowFlags = C.SDL_WINDOW_POPUP_MENU
	WINDOW_KEYBOARD_GRABBED    WindowFlags = C.SDL_WINDOW_KEYBOARD_GRABBED
	WINDOW_VULKAN              WindowFlags = C.SDL_WINDOW_VULKAN
	WINDOW_METAL               WindowFlags = C.SDL_WINDOW_METAL
	WINDOW_TRANSPARENT         WindowFlags = C.SDL_WINDOW_TRANSPARENT
	WINDOW_NOT_FOCUSABLE       WindowFlags = C.SDL_WINDOW_NOT_FOCUSABLE
)

const (
	WINDOWPOS_UNDEFINED = C.SDL_WINDOWPOS_UNDEFINED
	WINDOWPOS_CENTERED  = C.SDL_WINDOWPOS_CENTERED
)

func (w *Window) Flags() WindowFlags {
	return C.SDL_GetWindowFlags((*C.SDL_Window)(w))
}

func (w *Window) Show() bool {
	return bool(C.SDL_ShowWindow((*C.SDL_Window)(w)))
}

func (w *Window) Hide() bool {
	return bool(C.SDL_HideWindow((*C.SDL_Window)(w)))
}

func (w *Window) Raise() bool {
	return bool(C.SDL_RaiseWindow((*C.SDL_Window)(w)))
}

func (w *Window) SetTitle(title string) {
	ctitle := append([]byte(title), 0)
	C.SDL_SetWindowTitle((*C.SDL_Window)(w), (*C.char)(unsafe.Pointer(&ctitle[0])))
}

func (w *Window) Title() string {
	return C.GoString(C.SDL_GetWindowTitle((*C.SDL_Window)(w)))
}

func ScreenSaverEnabled() bool {
	return bool(C.SDL_ScreenSaverEnabled())
}

func EnableScreenSaver() bool {
	return bool(C.SDL_EnableScreenSaver())
}

func DisableScreenSaver() bool {
	return bool(C.SDL_DisableScreenSaver())
}
