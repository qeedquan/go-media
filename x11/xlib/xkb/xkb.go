package xkb

/*
#include <X11/XKBlib.h>
#cgo pkg-config: x11
*/
import "C"
import "github.com/qeedquan/go-media/x11/xlib"

func Bell(display *xlib.Display, window xlib.Window, percent int, name xlib.Atom) bool {
	return C.XkbBell((*C.Display)(display), C.Window(window), C.int(percent), C.Atom(name)) != 0
}