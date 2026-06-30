package cairo

/*
#include <cairo.h>
*/
import "C"
import "unsafe"

type Pattern C.cairo_pattern_t
type Extend C.cairo_extend_t
type Filter C.cairo_filter_t
type PatternType C.cairo_pattern_type_t

const (
	EXTEND_NONE   Extend = C.CAIRO_EXTEND_NONE
	EXTEND_REPEAT Extend = C.CAIRO_EXTEND_REPEAT
	EXTEND_PAD    Extend = C.CAIRO_EXTEND_PAD
)

const (
	FILTER_FAST     Filter = C.CAIRO_FILTER_FAST
	FILTER_GOOD     Filter = C.CAIRO_FILTER_GOOD
	FILTER_BEST     Filter = C.CAIRO_FILTER_BEST
	FILTER_BILINEAR Filter = C.CAIRO_FILTER_BILINEAR
	FILTER_GAUSSIAN Filter = C.CAIRO_FILTER_GAUSSIAN
)

const (
	PATTERN_TYPE_SOLID         PatternType = C.CAIRO_PATTERN_TYPE_SOLID
	PATTERN_TYPE_SURFACE       PatternType = C.CAIRO_PATTERN_TYPE_SURFACE
	PATTERN_TYPE_LINEAR        PatternType = C.CAIRO_PATTERN_TYPE_LINEAR
	PATTERN_TYPE_RADIAL        PatternType = C.CAIRO_PATTERN_TYPE_RADIAL
	PATTERN_TYPE_MESH          PatternType = C.CAIRO_PATTERN_TYPE_MESH
	PATTERN_TYPE_RASTER_SOURCE PatternType = C.CAIRO_PATTERN_TYPE_RASTER_SOURCE
)

func (p *Pattern) Ref() *Pattern {
	return (*Pattern)(C.cairo_pattern_reference((*C.cairo_pattern_t)(p)))
}

func (p *Pattern) Destroy() {
	C.cairo_pattern_destroy((*C.cairo_pattern_t)(p))
}

func (p *Pattern) Status() Status {
	return Status(C.cairo_pattern_status((*C.cairo_pattern_t)(p)))
}

func (p *Pattern) AddColorStopRGB(offset, red, green, blue float64) {
	C.cairo_pattern_add_color_stop_rgb((*C.cairo_pattern_t)(p), C.double(offset), C.double(red), C.double(green), C.double(blue))
}

func (p *Pattern) AddColorStopRGBA(offset, red, green, blue, alpha float64) {
	C.cairo_pattern_add_color_stop_rgba((*C.cairo_pattern_t)(p), C.double(offset), C.double(red), C.double(green), C.double(blue), C.double(alpha))
}

func (p *Pattern) RGBA() (red, green, blue, alpha float64) {
	var cr, cg, cb, ca C.double
	C.cairo_pattern_get_rgba((*C.cairo_pattern_t)(p), &cr, &cg, &cb, &ca)
	red, green, blue, alpha = float64(cr), float64(cg), float64(cb), float64(ca)
	return
}

func (p *Pattern) LinearPoints() (x0, y0, x1, y1 float64, err error) {
	var a, b, c, d C.double
	err = Status(C.cairo_pattern_get_linear_points((*C.cairo_pattern_t)(p), &a, &b, &c, &d))
	x0 = float64(a)
	y0 = float64(b)
	x1 = float64(c)
	y1 = float64(d)
	return
}

func (p *Pattern) RadialCircles() (x0, y0, r0, x1, y1, r1 float64, err error) {
	var a, b, c, d, e, f C.double
	err = Status(C.cairo_pattern_get_radial_circles((*C.cairo_pattern_t)(p), &a, &b, &c, &d, &e, &f))
	x0 = float64(a)
	y0 = float64(b)
	r0 = float64(c)
	x1 = float64(d)
	y1 = float64(e)
	r1 = float64(f)
	return
}

func (p *Pattern) BeginPatch() {
	C.cairo_mesh_pattern_begin_patch((*C.cairo_pattern_t)(p))
}

func (p *Pattern) EndPatch() {
	C.cairo_mesh_pattern_end_patch((*C.cairo_pattern_t)(p))
}

func (p *Pattern) MoveTo(x, y float64) {
	C.cairo_mesh_pattern_move_to((*C.cairo_pattern_t)(p), C.double(x), C.double(y))
}

func (p *Pattern) LineTo(x, y float64) {
	C.cairo_mesh_pattern_line_to((*C.cairo_pattern_t)(p), C.double(x), C.double(y))
}

func (p *Pattern) CurveTo(x1, y1, x2, y2, x3, y3 float64) {
	C.cairo_mesh_pattern_curve_to((*C.cairo_pattern_t)(p), C.double(x1), C.double(y1), C.double(x2), C.double(y2), C.double(x3), C.double(y3))
}

func (p *Pattern) SetControlPoint(pointNum int, x, y float64) {
	C.cairo_mesh_pattern_set_control_point((*C.cairo_pattern_t)(p), C.uint(pointNum), C.double(x), C.double(y))
}

func (p *Pattern) SetExtend(extend Extend) {
	C.cairo_pattern_set_extend((*C.cairo_pattern_t)(p), (C.cairo_extend_t)(extend))
}

func (p *Pattern) Extend() Extend {
	return Extend(C.cairo_pattern_get_extend((*C.cairo_pattern_t)(p)))
}

func (p *Pattern) SetFilter(filter Filter) {
	C.cairo_pattern_set_filter((*C.cairo_pattern_t)(p), (C.cairo_filter_t)(filter))
}

func (p *Pattern) Filter() Filter {
	return Filter(C.cairo_pattern_get_filter((*C.cairo_pattern_t)(p)))
}

func (p *Pattern) SetMatrix(matrix *Matrix) {
	C.cairo_pattern_set_matrix((*C.cairo_pattern_t)(p), (*C.cairo_matrix_t)(unsafe.Pointer(matrix)))
}

func (p *Pattern) Matrix() Matrix {
	var matrix Matrix
	C.cairo_pattern_get_matrix((*C.cairo_pattern_t)(p), (*C.cairo_matrix_t)(unsafe.Pointer(&matrix)))
	return matrix
}

func CreatePatternRGB(red, green, blue float64) *Pattern {
	return (*Pattern)(C.cairo_pattern_create_rgb(C.double(red), C.double(green), C.double(blue)))
}

func CreatePatternRGBA(red, green, blue, alpha float64) *Pattern {
	return (*Pattern)(C.cairo_pattern_create_rgba(C.double(red), C.double(green), C.double(blue), C.double(alpha)))
}

func CreatePatternLinear(x0, y0, x1, y1 float64) *Pattern {
	return (*Pattern)(C.cairo_pattern_create_linear(C.double(x0), C.double(y0), C.double(x1), C.double(y1)))
}

func CreatePatternRadial(cx0, cy0, radius0, cx1, cy1, radius1 float64) *Pattern {
	return (*Pattern)(C.cairo_pattern_create_radial(C.double(cx0), C.double(cy0), C.double(radius0), C.double(cx1), C.double(cy1), C.double(radius1)))
}

func CreatePatternMesh() *Pattern {
	return (*Pattern)(C.cairo_pattern_create_mesh())
}
