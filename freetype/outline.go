package freetype

/*
#include <ft2build.h>
#include FT_FREETYPE_H
#include FT_OUTLINE_H
*/
import "C"

type Outline C.FT_Outline

func (o *Outline) Translate(xOffset, yOffset Pos) {
	C.FT_Outline_Translate((*C.FT_Outline)(o), C.FT_Pos(xOffset), C.FT_Pos(yOffset))
}
