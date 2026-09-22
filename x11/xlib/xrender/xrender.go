package xrender

/*
#include <X11/extensions/Xrender.h>
*/
import "C"

type (
	Color             C.XRenderColor
	DirectFormat      C.XRenderDirectFormat
	Filters           C.XFilters
	Fixed             C.XFixed
	IndexValue        C.XIndexValue
	LineFixed         C.XLineFixed
	PictFormat        C.XRenderPictFormat
	PictureAttributes C.XRenderPictureAttributes
	PointFixed        C.XPointFixed
	RenderPictFormat  C.XRenderPictFormat
	Transform         C.XTransform
	Trapezoid         C.XTrapezoid
	Triangle          C.XTriangle
)

const (
	PictFormatID        = C.PictFormatID
	PictFormatType      = C.PictFormatType
	PictFormatDepth     = C.PictFormatDepth
	PictFormatRed       = C.PictFormatRed
	PictFormatRedMask   = C.PictFormatRedMask
	PictFormatGreen     = C.PictFormatGreen
	PictFormatGreenMask = C.PictFormatGreenMask
	PictFormatBlue      = C.PictFormatBlue
	PictFormatBlueMask  = C.PictFormatBlueMask
	PictFormatAlpha     = C.PictFormatAlpha
	PictFormatAlphaMask = C.PictFormatAlphaMask
	PictFormatColormap  = C.PictFormatColormap
)

func (c *Color) Alpha() uint16 {
	return uint16(c.alpha)
}

func (c *Color) Red() uint16 {
	return uint16(c.red)
}

func (c *Color) Blue() uint16 {
	return uint16(c.blue)
}

func (c *Color) Green() uint16 {
	return uint16(c.green)
}
func (c *Color) SetRed(red uint16) {
	c.red = C.ushort(red)
}

func (c *Color) SetGreen(green uint16) {
	c.green = C.ushort(green)
}

func (c *Color) SetBlue(blue uint16) {
	c.blue = C.ushort(blue)
}

func (c *Color) SetAlpha(alpha uint16) {
	c.alpha = C.ushort(alpha)
}
