package freetype

/*
#include <ft2build.h>
#include FT_FREETYPE_H
#include FT_STROKER_H
#include FT_LCD_FILTER_H
*/
import "C"

type Encoding C.FT_Encoding
type StrokerLineCap C.FT_Stroker_LineCap
type StrokerLineJoin C.FT_Stroker_LineJoin
type RenderMode C.FT_Render_Mode
type KerningMode C.FT_Kerning_Mode

const (
	ENCODING_UNICODE Encoding = C.FT_ENCODING_UNICODE
)

const (
	LOAD_NO_BITMAP      = C.FT_LOAD_NO_BITMAP
	LOAD_RENDER         = C.FT_LOAD_RENDER
	LOAD_NO_HINTING     = C.FT_LOAD_NO_HINTING
	LOAD_NO_AUTOHINT    = C.FT_LOAD_NO_AUTOHINT
	LOAD_FORCE_AUTOHINT = C.FT_LOAD_FORCE_AUTOHINT
	LOAD_TARGET_LCD     = C.FT_LOAD_TARGET_LCD
)

const (
	LCD_FILTER_LIGHT LCDFilter = C.FT_LCD_FILTER_LIGHT
)

const (
	STROKER_LINECAP_ROUND StrokerLineCap = C.FT_STROKER_LINECAP_ROUND
)

const (
	STROKER_LINEJOIN_ROUND StrokerLineJoin = C.FT_STROKER_LINEJOIN_ROUND
)

const (
	RENDER_MODE_NORMAL RenderMode = C.FT_RENDER_MODE_NORMAL
	RENDER_MODE_LCD    RenderMode = C.FT_RENDER_MODE_LCD
)

const (
	KERNING_DEFAULT  KerningMode = C.FT_KERNING_DEFAULT
	KERNING_UNFITTED KerningMode = C.FT_KERNING_UNFITTED
	KERNING_UNSCALED KerningMode = C.FT_KERNING_UNSCALED
)
