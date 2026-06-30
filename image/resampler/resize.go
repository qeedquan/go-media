package resampler

import (
	"image"
	"image/color"
	"image/draw"
	"math"

	"github.com/qeedquan/go-media/math/ga"
)

var (
	srgb   [256]float64
	linear [4096]uint8
)

func init() {
	const gamma = 1.75
	for i := range srgb {
		srgb[i] = math.Pow(float64(i)/255, gamma)
	}

	for i := range linear {
		k := 255*math.Pow(float64(i)/float64(len(linear)), 1/gamma) + .5
		k = ga.Clamp(k, 0, 255)
		linear[i] = uint8(k)
	}
}

func ResizeImage(m image.Image, p draw.Image, o *Options) {
	var (
		resamplers [4]*Resampler
		samples    [4][]float64
	)
	dr := p.Bounds()
	sr := m.Bounds()
	sn := image.Pt(sr.Dx(), sr.Dy())
	dn := image.Pt(dr.Dx(), dr.Dy())
	for i := range resamplers {
		resamplers[i] = New(dn, sn, o)
		samples[i] = make([]float64, sn.X)
	}

	dy := 0
	for y := sr.Min.Y; y < sr.Max.Y; y++ {
		for x := sr.Min.X; x < sr.Max.X; x++ {
			c := color.RGBAModel.Convert(m.At(x, y)).(color.RGBA)
			samples[0][x] = srgb[c.R]
			samples[1][x] = srgb[c.G]
			samples[2][x] = srgb[c.B]
			samples[3][x] = float64(c.A) / 255
		}

		for i, rp := range resamplers {
			rp.PutLine(samples[i])
		}

	loop:
		for ; ; dy++ {
			var out [4][]float64
			for i := range resamplers {
				out[i] = resamplers[i].GetLine()
				if out[i] == nil {
					break loop
				}
			}

			for dx := 0; dx < dn.X; dx++ {
				c := color.RGBA{
					linear2srgb(out[0][dx]),
					linear2srgb(out[1][dx]),
					linear2srgb(out[2][dx]),
					linear2alpha(out[3][dx]),
				}
				p.Set(dx, dy, c)
			}
		}
	}
}

func linear2srgb(x float64) uint8 {
	i := float64(len(linear))*x + .5
	i = ga.Clamp(i, 0, float64(len(linear)-1))
	return linear[int(i)]
}

func linear2alpha(x float64) uint8 {
	return uint8(ga.Clamp(255*x+.5, 0, 255))
}
