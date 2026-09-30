package xrm

/*
#include <X11/Xlib.h>
#include <X11/Xresource.h>

#cgo pkg-config: x11
*/
import "C"

type (
	Database C.XrmDatabase
)
