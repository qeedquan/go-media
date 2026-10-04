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

func Float32n(min, max float32) float32 {
	return min + (max-min)*rand.Float32()
}

func Float64n(min, max float64) float64 {
	return min + (max-min)*rand.Float64()
}
