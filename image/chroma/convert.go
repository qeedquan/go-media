package chroma

import (
	"image/color"
	"math"

	"github.com/qeedquan/go-media/math/ga"
	"golang.org/x/exp/constraints"
)

type HSL struct {
	H, S, L float64
}

type HSV struct {
	H, S, V float64
}

func (h HSV) RGBA() (r, g, b, a uint32) {
	c := HSVToRGBA(h)
	return color.RGBA{c.R, c.G, c.B, c.A}.RGBA()
}

func (h HSL) RGBA() (r, g, b, a uint32) {
	c := HSLToHSV(h)
	return c.RGBA()
}

func HSVToRGBA(c HSV) color.RGBA {
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
	return Vec4ToRGBA(ga.Vec4d{r, g, b, 1})
}

func RGBAToHSV(c color.RGBA) HSV {
	r := float64(c.R) / 255
	g := float64(c.G) / 255
	b := float64(c.B) / 255

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

func HSVToHSL(c HSV) HSL {
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

func HSLToHSV(c HSL) HSV {
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

func Vec3ToRGBA[T constraints.Float](v ga.Vec3[T]) color.RGBA {
	r := ga.Clamp(v.X*255, 0, 255)
	g := ga.Clamp(v.Y*255, 0, 255)
	b := ga.Clamp(v.Z*255, 0, 255)
	return color.RGBA{uint8(r), uint8(g), uint8(b), 255}
}

func Vec4ToRGBA[T constraints.Float](v ga.Vec4[T]) color.RGBA {
	r := ga.Clamp(v.X*255, 0, 255)
	g := ga.Clamp(v.Y*255, 0, 255)
	b := ga.Clamp(v.Z*255, 0, 255)
	a := ga.Clamp(v.W*255, 0, 255)
	return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
}
