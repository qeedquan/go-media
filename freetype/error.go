package freetype

/*
#include <ft2build.h>
#include FT_FREETYPE_H
*/
import "C"

type Error C.FT_Error

var errStrings = [...]string{
	// generic errors
	0x00: "no error",
	0x01: "cannot open resource",
	0x02: "unknown file format",
	0x03: "broken file",
	0x04: "invalid freetype version",
	0x05: "module version is too low",
	0x06: "invalid argument",
	0x07: "unimplemented feature",
	0x08: "broken table",
	0x09: "broken offset within table",
	0x0A: "array allocation size too large",
	0x0B: "missing module",
	0x0C: "missing property",

	// glyph/character errors
	0x10: "invalid glyph index",
	0x11: "invalid character code",
	0x12: "unsupported glyph image format",
	0x13: "cannot render this glyph format",
	0x14: "invalid outline",
	0x15: "invalid composite glyph",
	0x16: "too many hints",
	0x17: "invalid pixel size",

	// handle errors
	0x20: "invalid object handle",
	0x21: "invalid library handle",
	0x22: "invalid module handle",
	0x23: "invalid face handle",
	0x24: "invalid size handle",
	0x25: "invalid glyph slot handle",
	0x26: "invalid charmap handle",
	0x27: "invalid cache manager handle",
	0x28: "invalid stream handle",

	// driver errors
	0x30: "too many modules",
	0x31: "too many extensions",

	// memory errors
	0x40: "out of memory",
	0x41: "unlisted object",

	// stream errors
	0x51: "cannot open stream",
	0x52: "invalid stream seek",
	0x53: "invalid stream skip",

	0xBA: "font glyphs corrupted or missing fields",
}

func (e Error) Error() string {
	const err = "freetype: unknown error"
	if !(0 <= e && e <= Error(len(errStrings))) {
		return err
	}

	errstr := errStrings[e]
	if errstr == "" {
		return err
	}

	return "freetype: " + errstr
}
