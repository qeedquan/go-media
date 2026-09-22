package sdl

/*
#include <SDL.h>
*/
import "C"
import (
	"unsafe"
)

type (
	MessageBoxFlags       C.SDL_MessageBoxFlags
	MessageBoxButtonFlags C.SDL_MessageBoxButtonFlags

	MessageBoxButtonData struct {
		Flags    MessageBoxButtonFlags
		ButtonID int
		Text     string
	}

	MessageBoxColor     struct{ R, G, B uint8 }
	MessageBoxColorType C.SDL_MessageBoxColorType

	MessageBoxColorScheme struct {
		Colors [MESSAGEBOX_COLOR_MAX]MessageBoxColor
	}

	MessageBoxData struct {
		Flags       MessageBoxFlags
		Window      *Window
		Title       string
		Message     string
		Buttons     []MessageBoxButtonData
		ColorScheme *MessageBoxColorScheme
	}
)

const (
	MESSAGEBOX_ERROR       MessageBoxFlags = C.SDL_MESSAGEBOX_ERROR
	MESSAGEBOX_WARNING     MessageBoxFlags = C.SDL_MESSAGEBOX_WARNING
	MESSAGEBOX_INFORMATION MessageBoxFlags = C.SDL_MESSAGEBOX_INFORMATION
)

const (
	MESSAGEBOX_BUTTON_RETURNKEY_DEFAULT MessageBoxButtonFlags = C.SDL_MESSAGEBOX_BUTTON_RETURNKEY_DEFAULT
	MESSAGEBOX_BUTTON_ESCAPEKEY_DEFAULT MessageBoxButtonFlags = C.SDL_MESSAGEBOX_BUTTON_ESCAPEKEY_DEFAULT
)

const (
	MESSAGEBOX_COLOR_BACKGROUND        MessageBoxColorType = C.SDL_MESSAGEBOX_COLOR_BACKGROUND
	MESSAGEBOX_COLOR_TEXT              MessageBoxColorType = C.SDL_MESSAGEBOX_COLOR_TEXT
	MESSAGEBOX_COLOR_BUTTON_BORDER     MessageBoxColorType = C.SDL_MESSAGEBOX_COLOR_BUTTON_BORDER
	MESSAGEBOX_COLOR_BUTTON_BACKGROUND MessageBoxColorType = C.SDL_MESSAGEBOX_COLOR_BUTTON_BACKGROUND
	MESSAGEBOX_COLOR_BUTTON_SELECTED   MessageBoxColorType = C.SDL_MESSAGEBOX_COLOR_BUTTON_SELECTED
	MESSAGEBOX_COLOR_MAX                                   = C.SDL_MESSAGEBOX_COLOR_MAX
)

func ShowSimpleMessageBox(flags MessageBoxFlags, title, message string, window *Window) error {
	ctitle := C.CString(title)
	cmessage := C.CString(message)
	defer C.free(unsafe.Pointer(ctitle))
	defer C.free(unsafe.Pointer(cmessage))
	return ek(C.SDL_ShowSimpleMessageBox(C.Uint32(flags), ctitle, cmessage, (*C.SDL_Window)(window)))
}

func ShowMessageBox(data *MessageBoxData) (int, error) {
	cdata := C.SDL_MessageBoxData{
		flags:      C.Uint32(data.Flags),
		window:     (*C.SDL_Window)(data.Window),
		title:      C.CString(data.Title),
		message:    C.CString(data.Message),
		numbuttons: C.int(len(data.Buttons)),
	}
	defer C.free(unsafe.Pointer(cdata.title))
	defer C.free(unsafe.Pointer(cdata.message))

	if size := len(data.Buttons); size > 0 {
		rsize := C.size_t(unsafe.Sizeof(C.SDL_MessageBoxButtonData{}))
		cdata.buttons = (*C.SDL_MessageBoxButtonData)(C.malloc(C.size_t(size) * rsize))
		defer C.free(unsafe.Pointer(cdata.buttons))

		ptr := (*[1 << 27]C.SDL_MessageBoxButtonData)(unsafe.Pointer(cdata.buttons))
		for i := 0; i < size; i++ {
			x := &((*ptr)[i])
			y := &data.Buttons[i]
			x.flags = C.Uint32(y.Flags)
			x.buttonid = C.int(y.ButtonID)
			x.text = C.CString(y.Text)
			defer C.free(unsafe.Pointer(x.text))
		}
	}

	if data.ColorScheme != nil {
		var scheme C.SDL_MessageBoxColorScheme
		cdata.colorScheme = &scheme
		for i := range data.ColorScheme.Colors {
			x := &data.ColorScheme.Colors[i]
			y := &cdata.colorScheme.colors[i]
			y.r, y.g, y.b = C.Uint8(x.R), C.Uint8(x.G), C.Uint8(x.B)
		}
	}

	var buttonID C.int
	rc := C.SDL_ShowMessageBox(&cdata, &buttonID)
	return int(buttonID), ek(rc)
}
