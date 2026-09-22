package xfixes

/*
#include <X11/extensions/Xfixes.h>

#cgo pkg-config: x11
#cgo LDFLAGS: -lXfixes
*/
import "C"
import "github.com/qeedquan/go-media/x11/xlib"

func QueryExtension(display *xlib.Display) (supported bool, event_base, error_base int) {
	var cevent_base, cerror_base C.int
	rc := C.XFixesQueryExtension((*C.Display)(display), &cevent_base, &cerror_base)
	return rc != 0, int(cevent_base), int(cerror_base)
}

func ChangeSaveSet(display *xlib.Display, window xlib.Window, mode, target, map_ int) {
	C.XFixesChangeSaveSet((*C.Display)(display), C.Window(window), C.int(mode), C.int(target), C.int(map_))
}

func HideCursor(display *xlib.Display, window xlib.Window) {
	C.XFixesHideCursor((*C.Display)(display), C.Window(window))
}

func ShowCursor(display *xlib.Display, window xlib.Window) {
	C.XFixesShowCursor((*C.Display)(display), C.Window(window))
}
