package ga

import (
	"math"
	"math/cmplx"

	"golang.org/x/exp/constraints"
)

const (
	TAU     = 2 * math.Pi
	PI_2    = math.Pi / 2
	PI_4    = math.Pi / 4
	SQRT3   = 1.732050807568877293527446341505872366942805253810380628055806
	SQRT5   = 2.23606797749978969640917366873127623544061835961152572427089
	SQRT1_2 = 1 / math.Sqrt2
	SQRT1_3 = 1 / SQRT3
	SQRT1_5 = 1 / SQRT5
)

func Rad2Deg[T constraints.Float](x T) T {
	return x * 180 / math.Pi
}

func Deg2Rad[T constraints.Float](x T) T {
	return x * math.Pi / 180
}

func Clamp[T Ordinal](x, a, b T) T {
	return min(max(x, a), b)
}

func Lerp[T Number](t, a, b T) T {
	return a + t*(b-a)
}

func Unlerp[T Number](t, a, b T) T {
	return (t - a) / (b - a)
}

func LinearRemap[T Number](x, a, b, c, d T) T {
	return Lerp(Unlerp(x, a, b), c, d)
}

func Sign[T Signed](a T) T {
	if a < 0 {
		return -1
	}
	if a > 0 {
		return 1
	}
	return 0
}

func SignStrict[T Signed](a T) T {
	if a >= 0 {
		return 1
	}
	return -1
}

func Saturate[T Ordinal](x T) T {
	return Clamp(x, 0, 1)
}

func Abs[T Ordinal](x T) T {
	if x < 0 {
		x = -x
	}
	return x
}

func Align[T constraints.Integer](a, m T) T {
	return ((a + m - 1) / m) * m
}

func Multiple[T constraints.Float](a, m T) T {
	return Ceil(a/m) * m
}

func Gcd[T constraints.Integer](a, b T) T {
	a = Abs(a)
	b = Abs(b)
	k := max(a, b)
	m := min(a, b)
	for m != 0 {
		k, m = m, k%m
	}
	return k
}

func Lcm[T constraints.Signed](a, b T) T {
	return Abs(a*b) / Gcd(a, b)
}

func IsPow2[T constraints.Integer](x T) bool {
	return x&(x-1) == 0
}

func Ceil2[T constraints.Integer](x T) T {
	p := T(1)
	for p < x {
		p <<= 1
	}
	return p
}

func Floor2[T constraints.Integer](x T) T {
	v := Ceil2(x) >> 1
	if v == 0 {
		v = 1
	}
	return v
}

func Mod[T Signed](x, y T) T {
	return T(math.Mod(float64(x), float64(y)))
}

func Hypot[T Signed](x, y T) T {
	return T(math.Hypot(float64(x), float64(y)))
}

func Sqrt[T Signed](x T) T {
	return T(math.Sqrt(float64(x)))
}

func Cbrt[T Signed](x T) T {
	return T(math.Cbrt(float64(x)))
}

func Sin[T Signed](x T) T {
	return T(math.Sin(float64(x)))
}

func Cos[T Signed](x T) T {
	return T(math.Cos(float64(x)))
}

func Tan[T Signed](x T) T {
	return T(math.Tan(float64(x)))
}

func Sec[T Signed](x T) T {
	return 1 / Cos(x)
}

func Csc[T Signed](x T) T {
	return 1 / Sin(x)
}

func Cot[T Signed](x T) T {
	return 1 / Tan(x)
}

func Asin[T Signed](x T) T {
	return T(math.Asin(float64(x)))
}

func Acos[T Signed](x T) T {
	return T(math.Acos(float64(x)))
}

func Atan[T Signed](x T) T {
	return T(math.Atan(float64(x)))
}

func Atan2[T Signed](y, x T) T {
	return T(math.Atan2(float64(y), float64(x)))
}

func Sincos[T Signed](x T) (T, T) {
	s, c := math.Sincos(float64(x))
	return T(s), T(c)
}

func Cosh[T Signed](x T) T {
	return T(math.Cosh(float64(x)))
}

func Sinh[T Signed](x T) T {
	return T(math.Sinh(float64(x)))
}

func Tanh[T Signed](x T) T {
	return T(math.Tanh(float64(x)))
}

func Sech[T Signed](x T) T {
	return 1 / Cosh(x)
}

func Csch[T Signed](x T) T {
	return 1 / Sinh(x)
}

func Coth[T Signed](x T) T {
	return 1 / Tanh(x)
}

func Acosh[T Signed](x T) T {
	return T(math.Acosh(float64(x)))
}

func Asinh[T Signed](x T) T {
	return T(math.Asinh(float64(x)))
}

func Atanh[T Signed](x T) T {
	return T(math.Atanh(float64(x)))
}

func Copysign[T Signed](y, x T) T {
	return T(math.Copysign(float64(x), float64(y)))
}

func Exp[T Signed](x T) T {
	return T(math.Exp(float64(x)))
}

func Exp2[T Signed](x T) T {
	return T(math.Exp2(float64(x)))
}

func Gamma[T Signed](x T) T {
	return T(math.Gamma(float64(x)))
}

func Pow[T Signed](x, y T) T {
	return T(math.Pow(float64(x), float64(y)))
}

func Pow10[T Signed](x T) T {
	return T(math.Pow10(int(x)))
}

func Log[T Signed](x T) T {
	return T(math.Log(float64(x)))
}

func Log10[T Signed](x T) T {
	return T(math.Log10(float64(x)))
}

func Log2[T Signed](x T) T {
	return T(math.Log2(float64(x)))
}

func Logb[T Signed](x T) T {
	return T(math.Logb(float64(x)))
}

func Floor[T Signed](x T) T {
	return T(math.Floor(float64(x)))
}

func Ceil[T Signed](x T) T {
	return T(math.Ceil(float64(x)))
}

func Round[T Signed](x T) T {
	return T(math.Round(float64(x)))
}

func RoundToEven[T Signed](x T) T {
	return T(math.RoundToEven(float64(x)))
}

func Dim[T Signed](x, y T) T {
	return T(math.Dim(float64(x), float64(y)))
}

func Erf[T Signed](x T) T {
	return T(math.Erf(float64(x)))
}

func Erfc[T Signed](x T) T {
	return T(math.Erfc(float64(x)))
}

func Erfinv[T Signed](x T) T {
	return T(math.Erfinv(float64(x)))
}

func Erfcinv[T Signed](x T) T {
	return T(math.Erfcinv(float64(x)))
}

func FMA[T Signed](x, a, b T) T {
	return T(math.FMA(float64(x), float64(a), float64(b)))
}

func Simpson1D[T Ordinal](f func(x T) T, start, end T, n int) T {
	r := T(0)
	s := (end - start) / T(n)
	i := 0

	r += f(start)
	for j := 1; j < n; j++ {
		r += (4 - T(i<<1)) * f(start+T(j)*s)
		i = (i + 1) & 1
	}
	r += f(end)
	r *= s / 3
	return r
}

func simpsonweight[T Number](i, n int) T {
	if i == 0 || i == n {
		return 1
	}
	if i%2 != 0 {
		return 4
	}
	return 2
}

func Simpson2D[T Ordinal](f func(x, y T) T, x0, x1, y0, y1 T, m, n int) T {
	if n%2 != 0 || m%2 != 0 {
		panic("integration range must be even")
	}

	dx := (x1 - x0) / T(m)
	dy := (y1 - y0) / T(n)
	r := T(0)
	for i := 0; i <= n; i++ {
		y := y0 + T(i)*dy
		wy := simpsonweight[T](i, n)
		for j := 0; j <= m; j++ {
			x := x0 + T(j)*dx
			wx := simpsonweight[T](j, m)
			r += f(x, y) * wx * wy
		}
	}
	r *= dx * dy / (9 * T(m*n))
	return r
}

func Wrap[T Ordinal](x, a, b T) T {
	if x < a {
		x += b
	}
	if x >= b {
		x -= b
	}
	return x
}

func Smoothstep[T Signed](x, a, b T) T {
	x = Clamp((x-a)/(b-a), 0.0, 1.0)
	return x * x * (3 - 2*x)
}

func Sinc[T Signed](x T) T {
	if x == 0 {
		return 1
	}
	return Sin(x) / x
}

func Sum[T Number](a ...T) T {
	var r T
	for _, v := range a {
		r += v
	}
	return r
}

func Prod[T Number](a ...T) T {
	if len(a) == 0 {
		return 0
	}

	r := a[0]
	for _, v := range a[1:] {
		r *= v
	}
	return r
}

func Binomial[T constraints.Integer](n, k T) T {
	if k < 0 || k > n {
		return 0
	}

	if k == 0 || k == n {
		return 1
	}

	k = min(k, n-k)
	c := T(1)
	for i := T(0); i < k; i++ {
		c = c * (n - i) / (i + 1)
	}
	return c
}

func Factorial[T constraints.Integer](n T) T {
	if n < 0 {
		return 0
	}

	r := T(1)
	for i := T(2); i <= n; i++ {
		r *= i
	}
	return r
}

func Mean[T Ordinal](a []T) T {
	n := len(a)
	if n == 0 {
		return 0
	}

	r := T(0)
	for i := range a {
		r += a[i]
	}
	return r / T(n)
}

func Variance[T Ordinal](a []T, ddof T) T {
	n := len(a)
	if n == 0 {
		return 0
	}

	m := Mean(a)
	r := T(0)
	for i := range a {
		r += (a[i] - m) * (a[i] - m)
	}
	return r / (T(n) - ddof)
}

func Roots2c[T constraints.Complex](a, b, c T) [2]T {
	d := cmplx.Sqrt(complex128(b*b - 4*a*c))
	return [2]T{
		(-b + T(d)) / (2 * a),
		(-b - T(d)) / (2 * a),
	}
}
