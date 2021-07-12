package rng

import (
	"math"
	"math/rand"
)

type RNG interface {
	Int() int
	Float64() float64
}

// https://en.wikipedia.org/wiki/Poisson_distribution
// Poisson generates Poisson-distributed random variables
func Poisson(lambda float64) float64 {
	const step = 500

	l := lambda
	k := 0.0
	p := 1.0
	for {
		k++
		u := rand.Float64()
		if u == 0 {
			u += 1e-3
		}
		p *= u
		for p < 1 && l > 0 {
			if l > step {
				p *= math.Exp(step)
				l -= step
			} else {
				p *= math.Exp(l)
				l = 0
			}
		}
		if !(p > 1) {
			break
		}
	}
	return k - 1
}

// https://en.wikipedia.org/wiki/Box%E2%80%93Muller_transform
// BoxMuller transforms two uniform sample into two normally distributed samples
func BoxMuller(u1, u2 float64) (z0, z1 float64) {
	r := math.Sqrt(-2 * math.Log(u1))
	t := 2 * math.Pi * u2
	z0 = r * math.Cos(t)
	z1 = r * math.Sin(t)
	return
}

func Float32v(n int) []float32 {
	p := make([]float32, n)
	for i := range p {
		p[i] = rand.Float32()
	}
	return p
}

func Float32vn(n int, s float32) []float32 {
	p := Float32v(n)
	for i := range p {
		p[i] *= s
	}
	return p
}

func Float64v(n int) []float64 {
	p := make([]float64, n)
	for i := range p {
		p[i] = rand.Float64()
	}
	return p
}

func Float64vn(n int, s float64) []float64 {
	p := Float64v(n)
	for i := range p {
		p[i] *= s
	}
	return p
}

func Complex64() complex64 {
	return complex(rand.Float32(), rand.Float32())
}

func Complex64n(n float32) complex64 {
	return complex(rand.Float32()*n, rand.Float32()*n)
}

func Complex64v(n int) []complex64 {
	p := make([]complex64, n)
	for i := range p {
		p[i] = Complex64()
	}
	return p
}

func Complex64vn(n int, s float32) []complex64 {
	p := make([]complex64, n)
	for i := range p {
		p[i] = Complex64n(s)
	}
	return p
}

func Complex128() complex128 {
	return complex(rand.Float64(), rand.Float64())
}

func Complex128n(n float64) complex128 {
	return complex(rand.Float64()*n, rand.Float64()*n)
}

func Complex128v(n int) []complex128 {
	p := make([]complex128, n)
	for i := range p {
		p[i] = Complex128()
	}
	return p
}

func Complex128vn(n int, s float64) []complex128 {
	p := make([]complex128, n)
	for i := range p {
		p[i] = Complex128n(s)
	}
	return p
}
