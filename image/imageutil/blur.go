package imageutil

import (
	"image"
	"math"

	"github.com/qeedquan/go-media/math/ga"
)

func GaussianBlur(m image.Image, s float64) *image.RGBA {
	n := gaussKernelSize(s)
	k := gaussCoeff(n, s)
	p := Filter(m, k, &FilterOptions{
		Op:   OpConv,
		Wrap: WrapRepeat,
	})
	return p.ToRGBA()
}

func gaussKernelSize(s float64) int {
	return 1 + int(2*math.Ceil(math.Sqrt(-2*s*s*math.Log(0.005))))
}

func gaussCoeff(n int, s float64) [][]float64 {
	w := genMatrix(n, n)
	f := func(x, y float64) float64 {
		return gauss2D(x, y, 0, 0, s)
	}

	const N = 50
	wn := 0.0
	for i := 0; i < n; i++ {
		y0 := float64(i-n/2-1) + 0.5
		y1 := y0 + 1
		for j := 0; j < n; j++ {
			x0 := float64(j-n/2-1) + 0.5
			x1 := x0 + 1
			w[i][j] = ga.Simpson2D(f, x0, x1, y0, y1, N, N)
			wn += w[i][j]
		}
	}
	for i := range w {
		for j := range w[i] {
			w[i][j] /= wn
		}
	}
	return w
}

func gauss2D(x, y, mx, my, s float64) float64 {
	dx := x - mx
	dy := y - my
	n := 1 / (2 * math.Pi * s * s)
	return n * math.Exp(-(dx*dx+dy*dy)/(2*s*s))
}

func genMatrix(r, c int) [][]float64 {
	p := make([][]float64, r)
	q := make([]float64, r*c)
	for i := range p {
		p[i] = q[i*c : (i+1)*c]
	}
	return p
}
