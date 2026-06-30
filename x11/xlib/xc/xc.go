package xc

/*
#include <X11/cursorfont.h>

#cgo pkg-config: x11
*/
import "C"

const (
	Xcursor = C.XC_X_cursor
	Heart   = C.XC_heart
	Man     = C.XC_man
	Pirate  = C.XC_pirate
	UrAngle = C.XC_ur_angle
	Watch   = C.XC_watch
	Xterm   = C.XC_xterm
)
