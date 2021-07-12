package chroma

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/qeedquan/go-media/math/f64"
)

type HSL struct {
	H, S, L float64
}

type HSV struct {
	H, S, V float64
}

type Float4 [4]float64

func (f Float4) RGBA() (r, g, b, a uint32) {
	c := color.RGBA{uint8(f[0]), uint8(f[1]), uint8(f[2]), uint8(f[3])}
	return c.RGBA()
}

var (
	HSVModel    = color.ModelFunc(hsvModel)
	HSLModel    = color.ModelFunc(hslModel)
	Vec3dModel  = color.ModelFunc(vec3dModel)
	Vec4dModel  = color.ModelFunc(vec4dModel)
	Float4Model = color.ModelFunc(float4Model)
)

func float4Model(c color.Color) color.Color {
	n := color.RGBAModel.Convert(c).(color.RGBA)
	return Float4{
		float64(n.R),
		float64(n.G),
		float64(n.B),
		float64(n.A),
	}
}

func vec3dModel(c color.Color) color.Color {
	n := color.RGBAModel.Convert(c).(color.RGBA)
	return f64.Vec3{
		float64(n.R) / 255,
		float64(n.G) / 255,
		float64(n.B) / 255,
	}
}

func vec4dModel(c color.Color) color.Color {
	n := color.RGBAModel.Convert(c).(color.RGBA)
	return f64.Vec4{
		float64(n.R) / 255,
		float64(n.G) / 255,
		float64(n.B) / 255,
		float64(n.A) / 255,
	}
}

func hsvModel(c color.Color) color.Color {
	b := color.RGBAModel.Convert(c).(color.RGBA)
	return RGB2HSV(b)
}

func hslModel(c color.Color) color.Color {
	b := hsvModel(c).(HSV)
	return HSV2HSL(b)
}

func (h HSV) RGBA() (r, g, b, a uint32) {
	c := HSV2RGB(h)
	return color.RGBA{c.R, c.G, c.B, c.A}.RGBA()
}

func (h HSL) RGBA() (r, g, b, a uint32) {
	c := HSL2HSV(h)
	return c.RGBA()
}

func HSV2VEC4(c HSV) f64.Vec4 {
	h := c.H * 360
	s := c.S
	v := c.V

	hi := int(h/60.0) % 6
	f := (h / 60.0) - float64(hi)
	p := v * (1.0 - s)
	q := v * (1.0 - s*f)
	t := v * (1.0 - s*(1.0-f))

	var r, g, b float64
	switch hi {
	case 0:
		r, g, b = v, t, p
	case 1:
		r, g, b = q, v, p
	case 2:
		r, g, b = p, v, t
	case 3:
		r, g, b = p, q, v
	case 4:
		r, g, b = t, p, v
	case 5:
		r, g, b = v, p, q
	}
	return f64.Vec4{r, g, b, 1}
}

func HSV2RGB(c HSV) color.RGBA {
	return HSV2VEC4(c).ToRGBA()
}

func VEC42HSV(c f64.Vec4) HSV {
	r := c.X
	g := c.Y
	b := c.Z

	max := math.Max(r, math.Max(g, b))
	min := math.Min(r, math.Min(g, b))

	var h, s, v float64

	v = max
	if max == 0 || max == min {
		s = 0
		h = 0
	} else {
		s = (max - min) / max

		if max == r {
			h = 60*((g-b)/(max-min)) + 0
		} else if max == g {
			h = 60*((b-r)/(max-min)) + 120
		} else {
			h = 60*((r-g)/(max-min)) + 240
		}
	}
	if h < 0 {
		h += 360
	}
	h /= 360

	return HSV{h, s, v}
}

func RGB2HSV(c color.RGBA) HSV {
	return VEC42HSV(RGBA2VEC4(c))
}

func HSV2HSL(c HSV) HSL {
	h := c.H
	l := (2 - c.S) * c.V
	s := c.S * c.V
	if l <= 1 {
		s /= l
	} else {
		s /= 2 - l
	}
	l /= 2
	return HSL{h, s, l}
}

func HSL2HSV(c HSL) HSV {
	h := c.H
	l := c.L * 2
	s := c.S
	if l <= 1 {
		s *= l
	} else {
		s *= 2 - l
	}
	v := (l + s) / 2
	s = 2 * s / (l + s)
	return HSV{h, s, v}
}

func VEC42RGBA(c f64.Vec4) color.RGBA {
	const eps = 1.001
	if c.X <= eps {
		c.X *= 255
	}
	if c.Y <= eps {
		c.Y *= 255
	}
	if c.Z <= eps {
		c.Z *= 255
	}
	if c.W <= eps {
		c.W *= 255
	}
	c.X = f64.Clamp(c.X, 0, 255)
	c.Y = f64.Clamp(c.Y, 0, 255)
	c.Z = f64.Clamp(c.Z, 0, 255)
	c.W = f64.Clamp(c.W, 0, 255)
	return color.RGBA{
		uint8(c.X),
		uint8(c.Y),
		uint8(c.Z),
		uint8(c.W),
	}
}

func RGBA2VEC4(c color.RGBA) f64.Vec4 {
	return f64.Vec4{
		float64(c.R) / 255.0,
		float64(c.G) / 255.0,
		float64(c.B) / 255.0,
		float64(c.A) / 255.0,
	}
}

func ParseRGBA(s string) (color.RGBA, error) {
	var r, g, b, a uint8
	n, _ := fmt.Sscanf(s, "rgb(%v,%v,%v)", &r, &g, &b)
	if n == 3 {
		return color.RGBA{r, g, b, 255}, nil
	}

	n, _ = fmt.Sscanf(s, "rgba(%v,%v,%v,%v)", &r, &g, &b, &a)
	if n == 4 {
		return color.RGBA{r, g, b, a}, nil
	}

	n, _ = fmt.Sscanf(s, "#%02x%02x%02x%02x", &r, &g, &b, &a)
	if n == 4 {
		return color.RGBA{r, g, b, a}, nil
	}

	n, _ = fmt.Sscanf(s, "#%02x%02x%02x", &r, &g, &b)
	if n == 3 {
		return color.RGBA{r, g, b, 255}, nil
	}

	n, _ = fmt.Sscanf(s, "#%02x", &r)
	if n == 1 {
		return color.RGBA{r, r, r, 255}, nil
	}

	var h HSV
	n, _ = fmt.Sscanf(s, "hsv(%v,%v,%v)", &h.H, &h.S, &h.V)
	if n == 3 {
		return HSV2RGB(h), nil
	}

	return color.RGBA{}, fmt.Errorf("failed to parse color %q, unknown format", s)
}

func RandRGB() color.RGBA {
	return color.RGBA{
		uint8(rand.Intn(256)),
		uint8(rand.Intn(256)),
		uint8(rand.Intn(256)),
		255,
	}
}

func RandRGBA() color.RGBA {
	return color.RGBA{
		uint8(rand.Intn(256)),
		uint8(rand.Intn(256)),
		uint8(rand.Intn(256)),
		uint8(rand.Intn(256)),
	}
}

func RandHSV() HSV {
	return HSV{
		H: rand.Float64(),
		S: rand.Float64(),
		V: rand.Float64(),
	}
}

func MixRGBA(a, b color.RGBA, t float64) color.RGBA {
	return color.RGBA{
		uint8(float64(a.R)*(1-t) + t*float64(b.R)),
		uint8(float64(a.G)*(1-t) + t*float64(b.G)),
		uint8(float64(a.B)*(1-t) + t*float64(b.B)),
		uint8(float64(a.A)*(1-t) + t*float64(b.A)),
	}
}

func MixHSL(a, b HSL, t float64) HSL {
	return HSL{
		a.H*(1-t) + t*b.H,
		a.S*(1-t) + t*b.S,
		a.L*(1-t) + t*b.L,
	}
}

func RGBA32(c color.RGBA) uint32 {
	return uint32(c.R) | uint32(c.G)<<8 | uint32(c.B)<<16 | uint32(c.A)<<24
}

func BGRA32(c color.RGBA) uint32 {
	return uint32(c.B) | uint32(c.G)<<8 | uint32(c.R)<<16 | uint32(c.A)<<24
}

func AlphaBlendRGBA(a, b color.RGBA) color.RGBA {
	t := float64(b.A) / 255.0
	return color.RGBA{
		uint8(f64.Lerp(t, float64(a.R), float64(b.R))),
		uint8(f64.Lerp(t, float64(a.G), float64(b.G))),
		uint8(f64.Lerp(t, float64(a.B), float64(b.B))),
		255,
	}
}
