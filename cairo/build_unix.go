// +build linux darwin openbsd netbsd freebsd dragonfly solaris

package cairo

/*
#include <cairo.h>

#cgo pkg-config: cairo freetype2
*/
import "C"
