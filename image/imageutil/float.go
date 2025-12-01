package imageutil

import (
	"image"
	"image/color"

	"github.com/qeedquan/go-media/math/ga"
)

const (
	WrapClamp = iota
	WrapRepeat
)

const (
	OpConv = iota
	OpCorr
)

type FilterOptions struct {
	Op   int
	Wrap int
}

type Float struct {
	Pix    [][4]float64
	Stride int
	Rect   image.Rectangle
}

func NewFloat(r image.Rectangle) *Float {
	return &Float{
		Pix:    make([][4]float64, r.Dx()*r.Dy()),
		Stride: r.Dx(),
		Rect:   r,
	}
}

func (f *Float) Bounds() image.Rectangle {
	return f.Rect
}

func (f *Float) FloatAt(x, y int) [4]float64 {
	r := f.Bounds()
	if !image.Pt(x, y).In(r) {
		return [4]float64{}
	}

	x -= r.Min.X
	y -= r.Min.Y

	n := y*f.Stride + x
	if 0 <= n && n < len(f.Pix) {
		return f.Pix[n]
	}
	return [4]float64{}
}

func (f *Float) SetFloat(x, y int, c [4]float64) {
	r := f.Bounds()
	if !image.Pt(x, y).In(r) {
		return
	}

	x -= r.Min.X
	y -= r.Min.Y

	n := y*f.Stride + x
	if 0 <= n && n < len(f.Pix) {
		f.Pix[n] = c
	}
}

func (f *Float) ToRGB() *image.RGBA {
	r := f.Rect
	m := image.NewRGBA(r)
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			cf := f.FloatAt(x, y)
			cr := color.RGBA{
				uint8(ga.Clamp(cf[0], 0, 255)),
				uint8(ga.Clamp(cf[1], 0, 255)),
				uint8(ga.Clamp(cf[2], 0, 255)),
				255,
			}
			m.SetRGBA(x, y, cr)
		}
	}
	return m
}

func (f *Float) ToFloat() *Float {
	return &Float{
		Pix:    append([][4]float64{}, f.Pix...),
		Stride: f.Stride,
		Rect:   f.Rect,
	}
}

func (f *Float) ToRGBA() *image.RGBA {
	r := f.Rect
	m := image.NewRGBA(r)
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			cf := f.FloatAt(x, y)
			cr := color.RGBA{
				uint8(ga.Clamp(cf[0]*255, 0, 255)),
				uint8(ga.Clamp(cf[1]*255, 0, 255)),
				uint8(ga.Clamp(cf[2]*255, 0, 255)),
				uint8(ga.Clamp(cf[3]*255, 0, 255)),
			}
			m.SetRGBA(x, y, cr)
		}
	}
	return m
}

func (f *Float) Filter(kr [][]float64, o *FilterOptions) {
	if o == nil {
		o = &FilterOptions{
			Op:   OpConv,
			Wrap: WrapRepeat,
		}
	}

	if len(kr) == 0 || len(kr[0]) == 0 {
		return
	}
	a := len(kr)
	b := len(kr[0])
	r := f.Bounds()
	for i := r.Min.Y; i < r.Max.Y; i++ {
		for j := r.Min.X; j < r.Max.X; j++ {
			var s [4]float64
			for k := -a / 2; k <= a/2; k++ {
				for l := -b / 2; l <= b/2; l++ {
					y := i - k
					x := j - l
					if o.Op == OpCorr {
						y = i + k
						x = j + l
					}

					if o.Wrap == WrapRepeat {
						x = clamp(x, r.Min.X, r.Max.X-1)
						y = clamp(y, r.Min.Y, r.Max.Y-1)
					}

					c := f.FloatAt(x, y)
					for n := range s {
						s[n] += c[n] * kr[a/2+k][b/2+l]
					}
				}
			}
			f.SetFloat(j, i, s)
		}
	}
}

func ImageToFloat(m image.Image) *Float {
	r := m.Bounds()
	f := NewFloat(r)
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			cr := color.RGBAModel.Convert(m.At(x, y)).(color.RGBA)
			cf := [4]float64{float64(cr.R) / 255, float64(cr.G) / 255, float64(cr.B) / 255, float64(cr.A) / 255}
			f.SetFloat(x, y, cf)
		}
	}
	return f
}

func Filter(m image.Image, kr [][]float64, o *FilterOptions) *Float {
	f := ImageToFloat(m)
	f.Filter(kr, o)
	return f
}
