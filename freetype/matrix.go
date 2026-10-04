package freetype

/*
#include <ft2build.h>
#include FT_FREETYPE_H
*/
import "C"

type Matrix struct {
	XX, XY Fixed
	YX, YY Fixed
}

type Vector struct {
	X, Y Pos
}
