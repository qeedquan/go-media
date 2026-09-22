package resampler

import "math"

type Filter struct {
	Name    string
	Sample  func(float64) float64
	Support float64
}

func Box(t float64) float64 {
	if -0.5 <= t && t < 0.5 {
		return 1
	}
	return 0
}

func Tent(t float64) float64 {
	t = math.Abs(t)
	if t < 1 {
		return 1 - t
	}
	return 0
}

func Bell(t float64) float64 {
	t = math.Abs(t)
	if t < .5 {
		return .75 - t*t
	}
	if t < 1.5 {
		t -= 1.5
		return .5 * t * t
	}
	return 0
}

func Bspline(t float64) float64 {
	t = math.Abs(t)
	if t < 1 {
		tt := t * t
		return .5*tt*t - tt + 2/3.0
	}
	if t < 2 {
		t = 2 - t
		return (1.0 / 6) * t * t * t
	}
	return 0
}

// Dodgson, N., "Quadratic Interpolation for Image Resampling"
func quadratic(t, R float64) float64 {
	t = math.Abs(t)
	if t < 1.5 {
		tt := t * t
		if t <= .5 {
			return -2*R*tt + .5*(R+1)
		}
		return R*tt + (-2*R-.5)*t + 3.0/4*(R+1)
	}
	return 0
}

func Quadratic(t float64) float64 {
	return quadratic(t, 1)
}

func QuadraticApprox(t float64) float64 {
	return quadratic(t, .5)
}

func QuadraticMix(t float64) float64 {
	return quadratic(t, .8)
}

// Mitchell, D. and A. Netravali, "Reconstruction Filters in Computer Graphics."
// Computer Graphics, Vol. 22, No. 4, pp. 221-228.
// (B, C)
// (1/3, 1/3)  - Defaults recommended by Mitchell and Netravali
// (1, 0)	   - Equivalent to the Cubic B-Spline
// (0, 0.5)		- Equivalent to the Catmull-Rom Spline
// (0, C)		- The family of Cardinal Cubic Splines
// (B, 0)		- Duff's tensioned B-Splines.
func mitchell(t, B, C float64) float64 {
	tt := t * t
	t = math.Abs(t)

	if t < 1.0 {
		t = (((12.0 - 9.0*B - 6.0*C) * (t * tt)) +
			((-18.0 + 12.0*B + 6.0*C) * tt) +
			(6.0 - 2.0*B))

		return (t / 6.0)
	}

	if t < 2.0 {
		t = (((-1.0*B - 6.0*C) * (t * tt)) +
			((6.0*B + 30.0*C) * tt) +
			((-12.0*B - 48.0*C) * t) +
			(8.0*B + 24.0*C))

		return (t / 6.0)
	}

	return 0
}

func Mitchell(t float64) float64 {
	return mitchell(t, 1.0/3, 1.0/3)
}

func CatmullRom(t float64) float64 {
	return mitchell(t, 0, .5)
}

func sinc(x float64) float64 {
	x *= math.Pi
	if x < 0.01 && x > -0.01 {
		return 1 + x*x*((-1.0/6)+x*x*1.0/120)
	}

	return math.Sin(x) / x
}

func clean(t float64) float64 {
	const EPS = 0.0000125
	if math.Abs(t) < EPS {
		return 0
	}
	return t
}

func blackman(x float64) float64 {
	return 0.42659071 + 0.49656062*math.Cos(math.Pi*x) + 0.07684867*math.Cos(2*math.Pi*x)
}

func Blackman(t float64) float64 {
	t = math.Abs(t)
	if t < 3 {
		return clean(sinc(t) * blackman(t/3))
	}
	return 0
}

func Gaussian(t float64) float64 {
	t = math.Abs(t)
	if t < 1.25 {
		return clean(math.Exp(-2.0*t*t) * math.Sqrt(2/math.Pi) * blackman(t/1.25))
	}
	return 0
}

// Windowed sinc -- see "Jimm Blinn's Corner: Dirty Pixels" pg. 26.
func Lanczos(a float64) func(float64) float64 {
	return func(t float64) float64 {
		t = math.Abs(t)
		if t < a {
			return clean(sinc(t) * sinc(t/a))
		}
		return 0
	}
}

func kaiser(alpha, halfWidth, x float64) float64 {
	ratio := (x / halfWidth)
	return math.J0(alpha*math.Sqrt(1-ratio*ratio)) / math.J0(alpha)
}

func Kaiser(t float64) float64 {
	t = math.Abs(t)
	if t < 3 {
		att := 40.0
		alpha := (math.Exp(math.Log(0.58417*(att-20.96))*0.4) + 0.07886*(att-20.96))
		return clean(sinc(t) * kaiser(alpha, 3, t))
	}
	return 0
}

var Filters = []Filter{
	{"box", Box, 0.5},
	{"tent", Tent, 1},
	{"bell", Bell, 1.5},
	{"bspline", Bspline, 2},
	{"mithcell", Mitchell, 2},
	{"catmullrom", CatmullRom, 2},
	{"lanczos3", Lanczos(3), 3},
	{"lanczos4", Lanczos(4), 4},
	{"lanczos6", Lanczos(6), 6},
	{"lanczos12", Lanczos(12), 12},
	{"blackman", Blackman, 3},
	{"kaiser", Kaiser, 3},
	{"gaussian", Gaussian, 1.25},
	{"quadratic", Quadratic, 1.5},
	{"quadratic_approx", QuadraticApprox, 1.5},
	{"quadratic_mix", QuadraticMix, 1.5},
}

func GetFilter(name string) Filter {
	for _, f := range Filters {
		if f.Name == name {
			return f
		}
	}
	return Filter{}
}

func RegisterFilter(f Filter) {
	Filters = append(Filters, f)
}
