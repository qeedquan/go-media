package cairo

/*
#include <cairo.h>
*/
import "C"
import "unsafe"

type Cairo C.cairo_t
type LineCap C.cairo_line_cap_t
type LineJoin C.cairo_line_join_t
type FillRule C.cairo_fill_rule_t
type AntiAlias C.cairo_antialias_t
type HintStyle C.cairo_hint_style_t
type SubPixelOrder C.cairo_subpixel_order_t
type Operator C.cairo_operator_t

const (
	OPERATOR_CLEAR          Operator = C.CAIRO_OPERATOR_CLEAR
	OPERATOR_SOURCE         Operator = C.CAIRO_OPERATOR_SOURCE
	OPERATOR_OVER           Operator = C.CAIRO_OPERATOR_OVER
	OPERATOR_IN             Operator = C.CAIRO_OPERATOR_IN
	OPERATOR_OUT            Operator = C.CAIRO_OPERATOR_OUT
	OPERATOR_ATOP           Operator = C.CAIRO_OPERATOR_ATOP
	OPERATOR_DEST           Operator = C.CAIRO_OPERATOR_DEST
	OPERATOR_DEST_OVER      Operator = C.CAIRO_OPERATOR_DEST_OVER
	OPERATOR_DEST_IN        Operator = C.CAIRO_OPERATOR_DEST_IN
	OPERATOR_DEST_OUT       Operator = C.CAIRO_OPERATOR_DEST_OUT
	OPERATOR_DEST_ATOP      Operator = C.CAIRO_OPERATOR_DEST_ATOP
	OPERATOR_XOR            Operator = C.CAIRO_OPERATOR_XOR
	OPERATOR_ADD            Operator = C.CAIRO_OPERATOR_ADD
	OPERATOR_SATURATE       Operator = C.CAIRO_OPERATOR_SATURATE
	OPERATOR_MULTIPLY       Operator = C.CAIRO_OPERATOR_MULTIPLY
	OPERATOR_SCREEN         Operator = C.CAIRO_OPERATOR_SCREEN
	OPERATOR_OVERLAY        Operator = C.CAIRO_OPERATOR_OVERLAY
	OPERATOR_DARKEN         Operator = C.CAIRO_OPERATOR_DARKEN
	OPERATOR_LIGHTEN        Operator = C.CAIRO_OPERATOR_LIGHTEN
	OPERATOR_COLOR_DODGE    Operator = C.CAIRO_OPERATOR_COLOR_DODGE
	OPERATOR_COLOR_BURN     Operator = C.CAIRO_OPERATOR_COLOR_BURN
	OPERATOR_HARD_LIGHT     Operator = C.CAIRO_OPERATOR_HARD_LIGHT
	OPERATOR_SOFT_LIGHT     Operator = C.CAIRO_OPERATOR_SOFT_LIGHT
	OPERATOR_DIFFERENCE     Operator = C.CAIRO_OPERATOR_DIFFERENCE
	OPERATOR_EXCLUSION      Operator = C.CAIRO_OPERATOR_EXCLUSION
	OPERATOR_HSL_HUE        Operator = C.CAIRO_OPERATOR_HSL_HUE
	OPERATOR_HSL_SATURATION Operator = C.CAIRO_OPERATOR_HSL_SATURATION
	OPERATOR_HSL_COLOR      Operator = C.CAIRO_OPERATOR_HSL_COLOR
	OPERATOR_HSL_LUMINOSITY Operator = C.CAIRO_OPERATOR_HSL_LUMINOSITY
)

const (
	HINT_STYLE_DEFAULT HintStyle = C.CAIRO_HINT_STYLE_DEFAULT
	HINT_STYLE_FULL    HintStyle = C.CAIRO_HINT_STYLE_FULL
)

const (
	SUBPIXEL_ORDER_DEFAULT SubPixelOrder = C.CAIRO_SUBPIXEL_ORDER_DEFAULT
	SUBPIXEL_ORDER_RGB     SubPixelOrder = C.CAIRO_SUBPIXEL_ORDER_RGB
	SUBPIXEL_ORDER_BGR     SubPixelOrder = C.CAIRO_SUBPIXEL_ORDER_BGR
	SUBPIXEL_ORDER_VRGB    SubPixelOrder = C.CAIRO_SUBPIXEL_ORDER_VRGB
	SUBPIXEL_ORDER_VBGR    SubPixelOrder = C.CAIRO_SUBPIXEL_ORDER_VBGR
)

const (
	ANTIALIAS_DEFAULT AntiAlias = C.CAIRO_ANTIALIAS_DEFAULT

	ANTIALIAS_NONE     AntiAlias = C.CAIRO_ANTIALIAS_NONE
	ANTIALIAS_GRAY     AntiAlias = C.CAIRO_ANTIALIAS_GRAY
	ANTIALIAS_SUBPIXEL AntiAlias = C.CAIRO_ANTIALIAS_SUBPIXEL

	ANTIALIAS_FAST AntiAlias = C.CAIRO_ANTIALIAS_FAST
	ANTIALIAS_GOOD AntiAlias = C.CAIRO_ANTIALIAS_GOOD
	ANTIALIAS_BEST AntiAlias = C.CAIRO_ANTIALIAS_BEST
)

const (
	FILL_RULE_WINDING  FillRule = C.CAIRO_FILL_RULE_WINDING
	FILL_RULL_EVEN_ODD FillRule = C.CAIRO_FILL_RULE_EVEN_ODD
)

const (
	LINE_CAP_BUTT   LineCap = C.CAIRO_LINE_CAP_BUTT
	LINE_CAP_ROUND  LineCap = C.CAIRO_LINE_CAP_ROUND
	LINE_CAP_SQUARE LineCap = C.CAIRO_LINE_CAP_SQUARE
)

const (
	LINE_JOIN_MITER LineJoin = C.CAIRO_LINE_JOIN_MITER
	LINE_JOIN_BEVEL LineJoin = C.CAIRO_LINE_JOIN_BEVEL
	LINE_JOIN_ROUND LineJoin = C.CAIRO_LINE_JOIN_ROUND
)

func Create(target *Surface) *Cairo {
	return (*Cairo)(C.cairo_create((*C.cairo_surface_t)(target)))
}

func (c *Cairo) Ref() *Cairo {
	return (*Cairo)(C.cairo_reference((*C.cairo_t)(c)))
}

func (c *Cairo) Destroy() {
	C.cairo_destroy((*C.cairo_t)(c))
}

func (c *Cairo) Status() Status {
	return Status(C.cairo_status((*C.cairo_t)(c)))
}

func (c *Cairo) Save() {
	C.cairo_save((*C.cairo_t)(c))
}

func (c *Cairo) Restore() {
	C.cairo_restore((*C.cairo_t)(c))
}

func (c *Cairo) Target() *Surface {
	return (*Surface)(C.cairo_get_target((*C.cairo_t)(c)))
}

func (c *Cairo) PushGroup() {
	C.cairo_push_group((*C.cairo_t)(c))
}

func (c *Cairo) PopGroup() {
	C.cairo_pop_group((*C.cairo_t)(c))
}

func (c *Cairo) SetSourceRGB(red, green, blue float64) {
	C.cairo_set_source_rgb((*C.cairo_t)(c), C.double(red), C.double(green), C.double(blue))
}

func (c *Cairo) SetSourceRGBA(red, green, blue, alpha float64) {
	C.cairo_set_source_rgba((*C.cairo_t)(c), C.double(red), C.double(green), C.double(blue), C.double(alpha))
}

func (c *Cairo) SetOperator(op Operator) {
	C.cairo_set_operator((*C.cairo_t)(c), C.cairo_operator_t(op))
}

func (c *Cairo) Operator() Operator {
	return Operator(C.cairo_get_operator((*C.cairo_t)(c)))
}

func (c *Cairo) Stroke() {
	C.cairo_stroke((*C.cairo_t)(c))
}

func (c *Cairo) StrokePreserve() {
	C.cairo_stroke_preserve((*C.cairo_t)(c))
}

func (c *Cairo) Fill() {
	C.cairo_fill((*C.cairo_t)(c))
}

func (c *Cairo) FillPreserve() {
	C.cairo_fill_preserve((*C.cairo_t)(c))
}

func (c *Cairo) SetSource(source *Pattern) {
	C.cairo_set_source((*C.cairo_t)(c), (*C.cairo_pattern_t)(source))
}

func (c *Cairo) SetSourceSurface(surface *Surface, x, y float64) {
	C.cairo_set_source_surface((*C.cairo_t)(c), (*C.cairo_surface_t)(surface), C.double(x), C.double(y))
}

func (c *Cairo) Source() *Pattern {
	return (*Pattern)(C.cairo_get_source((*C.cairo_t)(c)))
}

func (c *Cairo) SetDash(dashes []float64, offset float64) {
	C.cairo_set_dash((*C.cairo_t)(c), (*C.double)(unsafe.Pointer(&dashes[0])), C.int(len(dashes)), C.double(offset))
}

func (c *Cairo) DashCount() int {
	return int(C.cairo_get_dash_count((*C.cairo_t)(c)))
}

func (c *Cairo) Rectangle(x, y, w, h float64) {
	C.cairo_rectangle((*C.cairo_t)(c), C.double(x), C.double(y), C.double(w), C.double(h))
}

func (c *Cairo) Arc(xc, yc, radius, angle1, angle2 float64) {
	C.cairo_arc((*C.cairo_t)(c), C.double(xc), C.double(yc), C.double(radius), C.double(angle1), C.double(angle2))
}

func (c *Cairo) ArcNegative(xc, yc, radius, angle1, angle2 float64) {
	C.cairo_arc_negative((*C.cairo_t)(c), C.double(xc), C.double(yc), C.double(radius), C.double(angle1), C.double(angle2))
}

func (c *Cairo) SetLineWidth(w float64) {
	C.cairo_set_line_width((*C.cairo_t)(c), C.double(w))
}

func (c *Cairo) LineTo(x, y float64) {
	C.cairo_line_to((*C.cairo_t)(c), C.double(x), C.double(y))
}

func (c *Cairo) MoveTo(x, y float64) {
	C.cairo_move_to((*C.cairo_t)(c), C.double(x), C.double(y))
}

func (c *Cairo) NewPath() {
	C.cairo_new_path((*C.cairo_t)(c))
}

func (c *Cairo) NewSubPath() {
	C.cairo_new_sub_path((*C.cairo_t)(c))
}

func (c *Cairo) ClosePath() {
	C.cairo_close_path((*C.cairo_t)(c))
}

func (c *Cairo) SetLineCap(linecap LineCap) {
	C.cairo_set_line_cap((*C.cairo_t)(c), C.cairo_line_cap_t(linecap))
}

func (c *Cairo) LineCap() LineCap {
	return LineCap(C.cairo_get_line_cap((*C.cairo_t)(c)))
}

func (c *Cairo) RelLineTo(dx, dy float64) {
	C.cairo_rel_line_to((*C.cairo_t)(c), C.double(dx), C.double(dy))
}

func (c *Cairo) RelMoveTo(dx, dy float64) {
	C.cairo_rel_move_to((*C.cairo_t)(c), C.double(dx), C.double(dy))
}

func (c *Cairo) RelCurveTo(dx1, dy1, dx2, dy2, dx3, dy3 float64) {
	C.cairo_rel_curve_to((*C.cairo_t)(c), C.double(dx1), C.double(dy1), C.double(dx2), C.double(dy2), C.double(dx3), C.double(dy3))
}

func (c *Cairo) SetLineJoin(linejoin LineJoin) {
	C.cairo_set_line_join((*C.cairo_t)(c), C.cairo_line_join_t(linejoin))
}

func (c *Cairo) Paint() {
	C.cairo_paint((*C.cairo_t)(c))
}

func (c *Cairo) PaintWithAlpha(alpha float64) {
	C.cairo_paint_with_alpha((*C.cairo_t)(c), C.double(alpha))
}

func (c *Cairo) SetAntiAlias(antialias AntiAlias) {
	C.cairo_set_antialias((*C.cairo_t)(c), C.cairo_antialias_t(antialias))
}

func (c *Cairo) AntiAlias() AntiAlias {
	return AntiAlias(C.cairo_get_antialias((*C.cairo_t)(c)))
}

func (c *Cairo) SetTolerance(tolerance float64) {
	C.cairo_set_tolerance((*C.cairo_t)(c), C.double(tolerance))
}

func (c *Cairo) Tolerance() float64 {
	return float64(C.cairo_get_tolerance(((*C.cairo_t)(c))))
}

func (c *Cairo) ClipPreserve() {
	C.cairo_clip_preserve((*C.cairo_t)(c))
}

func (c *Cairo) InClip(x, y float64) bool {
	rc := C.cairo_in_clip((*C.cairo_t)(c), C.double(x), C.double(y))
	if rc == 0 {
		return false
	}
	return true
}

func (c *Cairo) InFill(x, y float64) bool {
	return C.cairo_in_fill((*C.cairo_t)(c), C.double(x), C.double(y)) != 0
}

func (c *Cairo) ResetClip() {
	C.cairo_reset_clip((*C.cairo_t)(c))
}

func (c *Cairo) FillExtents() (x1, y1, x2, y2 float64) {
	var x, y, z, w C.double
	C.cairo_fill_extents((*C.cairo_t)(c), &x, &y, &z, &w)
	x1 = float64(x)
	y1 = float64(y)
	x2 = float64(z)
	y2 = float64(w)
	return
}

func (c *Cairo) StrokeExtents() (x1, y1, x2, y2 float64) {
	var x, y, z, w C.double
	C.cairo_stroke_extents((*C.cairo_t)(c), &x, &y, &z, &w)
	x1 = float64(x)
	y1 = float64(y)
	x2 = float64(z)
	y2 = float64(w)
	return
}

func (c *Cairo) SetScaledFont(scaledFont *ScaledFont) {
	C.cairo_set_scaled_font((*C.cairo_t)(c), (*C.cairo_scaled_font_t)(scaledFont))
}

func (c *Cairo) ScaledFont() *ScaledFont {
	return (*ScaledFont)(C.cairo_get_scaled_font((*C.cairo_t)(c)))
}

func (c *Cairo) Scale(sx, sy float64) {
	C.cairo_scale((*C.cairo_t)(c), C.double(sx), C.double(sy))
}

func (c *Cairo) Rotate(angle float64) {
	C.cairo_rotate((*C.cairo_t)(c), C.double(angle))
}

func (c *Cairo) Transform(matrix *Matrix) {
	C.cairo_transform((*C.cairo_t)(c), (*C.cairo_matrix_t)(unsafe.Pointer(matrix)))
}

func (c *Cairo) FontMatrix() Matrix {
	var m Matrix
	C.cairo_get_font_matrix((*C.cairo_t)(c), (*C.cairo_matrix_t)(unsafe.Pointer(&m)))
	return m
}

func (c *Cairo) Matrix() Matrix {
	var m Matrix
	C.cairo_get_matrix((*C.cairo_t)(c), (*C.cairo_matrix_t)(unsafe.Pointer(&m)))
	return m
}

func (c *Cairo) ShowPage() {
	C.cairo_show_page((*C.cairo_t)(c))
}

func (c *Cairo) ReferenceCount() int {
	return int(C.cairo_get_reference_count((*C.cairo_t)(c)))
}
