// https://www.scratchapixel.com/lessons/procedural-generation-virtual-worlds/procedural-patterns-noise-part-1/creating-simple-2D-noise
package valuenoise

import (
	"math"
	"math/rand"

	"github.com/qeedquan/go-media/math/ga"
)

type Rand struct {
	vals   []float64
	perm   []int
	mask   int
	smooth func(float64) float64
}

var defrand = New(512, Smooth)

func SetDefault(r *Rand) {
	defrand = r
}

func New(size int, smooth func(float64) float64) *Rand {
	vals := make([]float64, size)
	for i := range vals {
		vals[i] = rand.Float64()
	}
	return NewFromTable(vals, smooth)
}

func NewFromTable(vals []float64, smooth func(float64) float64) *Rand {
	mask := len(vals) - 1
	if len(vals)&mask != 0 {
		panic("size of table must be power of two")
	}

	perm := rand.Perm(len(vals))
	perm = append(perm, perm...)

	return &Rand{
		vals:   vals,
		perm:   perm,
		mask:   mask,
		smooth: smooth,
	}
}

func (r *Rand) Gen1D(x float64) float64 {
	x0 := int(x)
	x1 := x0 + 1
	x0 &= r.mask
	x1 &= r.mask

	_, t := math.Modf(x)
	t = r.smooth(t)

	return ga.Lerp(t, r.vals[r.perm[x0]], r.vals[r.perm[x1]])
}

func (r *Rand) Gen2D(x, y float64) float64 {
	xi := int(x)
	yi := int(y)

	_, tx := math.Modf(x)
	_, ty := math.Modf(y)

	rx0 := xi & r.mask
	rx1 := (rx0 + 1) & r.mask
	ry0 := yi & r.mask
	ry1 := (ry0 + 1) & r.mask

	c00 := r.vals[r.perm[r.perm[rx0]+ry0]]
	c10 := r.vals[r.perm[r.perm[rx1]+ry0]]
	c01 := r.vals[r.perm[r.perm[rx0]+ry1]]
	c11 := r.vals[r.perm[r.perm[rx1]+ry1]]

	sx := r.smooth(tx)
	sy := r.smooth(ty)

	nx0 := ga.Lerp(sx, c00, c10)
	nx1 := ga.Lerp(sx, c01, c11)

	return ga.Lerp(sy, nx0, nx1)
}

func Gen1D(x float64) float64    { return defrand.Gen1D(x) }
func Gen2D(x, y float64) float64 { return defrand.Gen2D(x, y) }

// http://sol.gfxile.net/interpolation
func Step(t float64) float64      { return math.Round(t + 0.5) }
func Linear(t float64) float64    { return t }
func Quadratic(t float64) float64 { return t * t }
func Smooth(t float64) float64    { return t * t * (3 - 2*t) }
func Smoother(t float64) float64  { return t * t * t * (t*(t*6-15) + 10) }
func Sin(t float64) float64       { return math.Sin(t * math.Pi / 2) }
