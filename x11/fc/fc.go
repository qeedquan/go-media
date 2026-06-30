package fc

/*
#include <stdlib.h>
#include <X11/Xft/Xft.h>

#cgo pkg-config: fontconfig
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type (
	Cache     C.FcCache
	Config    C.FcConfig
	CharSet   C.FcCharSet
	FontSet   C.FcFontSet
	MatchKind C.FcMatchKind
	Pattern   C.FcPattern
	Result    C.FcResult
)

const (
	CHARSET      = C.FC_CHARSET
	PIXEL_SIZE   = C.FC_PIXEL_SIZE
	SCALABLE     = C.FC_SCALABLE
	SIZE         = C.FC_SIZE
	SLANT        = C.FC_SLANT
	SLANT_ITALIC = C.FC_SLANT_ITALIC
	SLANT_ROMAN  = C.FC_SLANT_ROMAN
	WEIGHT       = C.FC_WEIGHT
	WEIGHT_BOLD  = C.FC_WEIGHT_BOLD
)

const (
	MatchPattern MatchKind = C.FcMatchPattern
)

const (
	ResultMatch        Result = C.FcResultMatch
	ResultTypeMismatch Result = C.FcResultTypeMismatch
)

func (r Result) Error() string {
	switch r {
	case ResultMatch:
		return "result matched"
	}
	return "unknown fontconfig error"
}

func Init() error {
	if C.FcInit() == 0 {
		return fmt.Errorf("failed to init fontconfig")
	}
	return nil
}

func NameParse(name string) *Pattern {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return (*Pattern)(C.FcNameParse((*C.FcChar8)(unsafe.Pointer(cname))))
}

func PatternGetInteger(pattern *Pattern, object string, n int) (r Result, i int) {
	var ci C.int
	cobject := C.CString(object)
	defer C.free(unsafe.Pointer(cobject))
	rc := C.FcPatternGetInteger((*C.FcPattern)(pattern), cobject, C.int(n), &ci)
	return Result(rc), int(ci)
}

func PatternGetDouble(pattern *Pattern, object string, n int) (r Result, d float64) {
	var cd C.double
	cobject := C.CString(object)
	defer C.free(unsafe.Pointer(cobject))
	rc := C.FcPatternGetDouble((*C.FcPattern)(pattern), cobject, C.int(n), &cd)
	return Result(rc), float64(cd)
}

func PatternDel(pattern *Pattern, object string) {
	cobject := C.CString(object)
	defer C.free(unsafe.Pointer(cobject))
	C.FcPatternDel((*C.FcPattern)(pattern), cobject)
}

func PatternAddInteger(pattern *Pattern, object string, i int) bool {
	cobject := C.CString(object)
	defer C.free(unsafe.Pointer(cobject))
	return C.FcPatternAddInteger((*C.FcPattern)(pattern), cobject, C.int(i)) != 0
}

func PatternAddBool(pattern *Pattern, object string, b bool) bool {
	cobject := C.CString(object)
	defer C.free(unsafe.Pointer(cobject))
	return C.FcPatternAddInteger((*C.FcPattern)(pattern), cobject, cbool(b)) != 0
}

func PatternAddDouble(pattern *Pattern, object string, d float64) bool {
	cobject := C.CString(object)
	defer C.free(unsafe.Pointer(cobject))
	return C.FcPatternAddDouble((*C.FcPattern)(pattern), cobject, C.double(d)) != 0
}

func PatternAddCharSet(pattern *Pattern, object string, c *CharSet) bool {
	cobject := C.CString(object)
	defer C.free(unsafe.Pointer(cobject))
	return C.FcPatternAddCharSet((*C.FcPattern)(pattern), cobject, (*C.FcCharSet)(c)) != 0
}

func PatternDuplicate(pattern *Pattern) *Pattern {
	return (*Pattern)(C.FcPatternDuplicate((*C.FcPattern)(pattern)))
}

func PatternDestroy(pattern *Pattern) {
	C.FcPatternDestroy((*C.FcPattern)(pattern))
}

func FontSetMatch(config *Config, sets []*FontSet, pattern *Pattern) (*Pattern, Result) {
	var r Result
	p := C.FcFontSetMatch((*C.FcConfig)(config), (**C.FcFontSet)((unsafe.Pointer)(&sets[0])), C.int(len(sets)), (*C.FcPattern)(pattern), (*C.FcResult)(&r))
	return (*Pattern)(p), r
}

func CharSetDestroy(fcs *CharSet) {
	C.FcCharSetDestroy((*C.FcCharSet)(fcs))
}

func FontMatch(config *Config, pattern *Pattern) (*Pattern, Result) {
	var cresult C.FcResult
	cpattern := C.FcFontMatch((*C.FcConfig)(config), (*C.FcPattern)(pattern), &cresult)
	return (*Pattern)(cpattern), Result(cresult)
}

func ConfigSubstitute(config *Config, pattern *Pattern, kind MatchKind) bool {
	return C.FcConfigSubstitute((*C.FcConfig)(config), (*C.FcPattern)(pattern), C.FcMatchKind(kind)) != 0
}

func DefaultSubstitute(pattern *Pattern) {
	C.FcDefaultSubstitute((*C.FcPattern)(pattern))
}

func FontSetDestroy(fontset *FontSet) {
	C.FcFontSetDestroy((*C.FcFontSet)(fontset))
}

func CharSetCreate() *CharSet {
	return (*CharSet)(C.FcCharSetCreate())
}

func CharSetAddChar(fcs *CharSet, ucs4 rune) bool {
	return C.FcCharSetAddChar((*C.FcCharSet)(fcs), C.FcChar32(ucs4)) != 0
}

func FontSort(config *Config, pattern *Pattern, trim bool) (set *FontSet, csp *CharSet, result Result) {
	set = (*FontSet)(C.FcFontSort((*C.FcConfig)(config), (*C.FcPattern)(pattern), cbool(trim), (**C.FcCharSet)(unsafe.Pointer(&csp)), (*C.FcResult)(&result)))
	return
}

func cbool(b bool) C.int {
	if b {
		return 1
	}
	return 0
}
