package vec2

import (
	"math"

	"github.com/qeedquan/go-media/math/ga"
)

func Add[T ga.Number](a, b ga.Vec2[T]) ga.Vec2[T] {
	return ga.Vec2[T]{
		a.X + b.X,
		a.Y + b.Y,
	}
}

func Sub[T ga.Number](a, b ga.Vec2[T]) ga.Vec2[T] {
	return ga.Vec2[T]{
		a.X - b.X,
		a.Y - b.Y,
	}
}

func Len[T ga.Signed](a ga.Vec2[T]) T {
	return T(math.Sqrt(float64(Dot(a, a))))
}

func Normalize[T ga.Signed](a ga.Vec2[T]) ga.Vec2[T] {
	l := Len(a)
	if l == 0 {
		return a
	}
	return Scale(a, 1/l)
}

func Proj[T ga.Signed](a, b ga.Vec2[T]) ga.Vec2[T] {
	b = Normalize(b)
	return Scale(b, Dot(a, b))
}

func Rej[T ga.Signed](a, b ga.Vec2[T]) ga.Vec2[T] {
	return Sub(a, Proj(a, b))
}

func Hadamard[T ga.Number](a, b ga.Vec2[T]) ga.Vec2[T] {
	return ga.Vec2[T]{
		a.X * b.X,
		a.Y * b.Y,
	}
}

func Scale[T ga.Number](a ga.Vec2[T], s T) ga.Vec2[T] {
	return ga.Vec2[T]{
		a.X * s,
		a.Y * s,
	}
}

func Fmas[T ga.Number](a, b ga.Vec2[T], c T) ga.Vec2[T] {
	return Add(a, Scale(b, c))
}

func Fma[T ga.Number](a, b, c ga.Vec2[T]) ga.Vec2[T] {
	return Add(a, Hadamard(b, c))
}

func Neg[T ga.Number](a ga.Vec2[T]) ga.Vec2[T] {
	return ga.Vec2[T]{-a.X, -a.Y}
}

func Perp[T ga.Number](a ga.Vec2[T]) ga.Vec2[T] {
	return ga.Vec2[T]{-a.Y, a.X}
}

func Dot[T ga.Number](a, b ga.Vec2[T]) T {
	return a.X*b.X + a.Y*b.Y
}

func Wedge[T ga.Number](a, b ga.Vec2[T]) T {
	return a.X*b.Y - a.Y*b.X
}

func Reflect[T ga.Number](a, b ga.Vec2[T]) ga.Vec2[T] {
	s := 2 * Dot(b, a) / Dot(a, a)
	return Sub(b, Scale(a, s))
}

func Min[T ga.Ordinal](a ga.Vec2[T], b ...ga.Vec2[T]) ga.Vec2[T] {
	for i := range b {
		a.X = min(a.X, b[i].X)
		a.Y = min(a.Y, b[i].Y)
	}
	return a
}

func Max[T ga.Ordinal](a ga.Vec2[T], b ...ga.Vec2[T]) ga.Vec2[T] {
	for i := range b {
		a.X = max(a.X, b[i].X)
		a.Y = max(a.Y, b[i].Y)
	}
	return a
}

func Clamp[T ga.Ordinal](x, a, b ga.Vec2[T]) ga.Vec2[T] {
	return ga.Vec2[T]{
		ga.Clamp(x.X, a.X, b.X),
		ga.Clamp(x.Y, a.Y, b.Y),
	}
}

func Lerp[T ga.Number](t T, a, b ga.Vec2[T]) ga.Vec2[T] {
	return ga.Vec2[T]{
		ga.Lerp(t, a.X, b.X),
		ga.Lerp(t, a.Y, b.Y),
	}
}

func Unlerp[T ga.Number](t T, a, b ga.Vec2[T]) ga.Vec2[T] {
	return ga.Vec2[T]{
		ga.Unlerp(t, a.X, b.X),
		ga.Unlerp(t, a.Y, b.Y),
	}
}

func Round[T ga.Signed](a ga.Vec2[T]) ga.Vec2[T] {
	return ga.Vec2[T]{
		ga.Round(a.X),
		ga.Round(a.Y),
	}
}

func Sign[T ga.Signed](a ga.Vec2[T]) ga.Vec2[T] {
	return ga.Vec2[T]{
		ga.Sign(a.X),
		ga.Sign(a.Y),
	}
}

func Saturate[T ga.Ordinal](a ga.Vec2[T]) ga.Vec2[T] {
	return ga.Vec2[T]{
		ga.Saturate(a.X),
		ga.Saturate(a.Y),
	}
}

func MinComp[T ga.Ordinal](a ga.Vec2[T]) T {
	return min(a.X, a.Y)
}

func MaxComp[T ga.Ordinal](a ga.Vec2[T]) T {
	return max(a.X, a.Y)
}

func Rotate[T ga.Signed](a ga.Vec2[T], r T) ga.Vec2[T] {
	s, c := ga.Sincos(r)
	return ga.Vec2[T]{
		a.X*c - a.Y*s,
		a.X*s + a.Y*c,
	}
}

func Vec3[T ga.Number](a ga.Vec2[T]) ga.Vec3[T] {
	return ga.Vec3[T]{a.X, a.Y, 1}
}

func Element[T ga.Number](a ga.Vec2[T]) (x, y T) {
	return a.X, a.Y
}

func Fill[T ga.Number](v T) ga.Vec2[T] {
	return ga.Vec2[T]{v, v}
}

func Distance[T ga.Signed](a, b ga.Vec2[T]) T {
	return Len(Sub(a, b))
}

func Angle[T ga.Signed](a, b ga.Vec2[T]) T {
	return ga.Acos(Dot(a, b) / (Len(a) * Len(b)))
}

func CartesianToPolar[T ga.Signed](a ga.Vec2[T]) ga.Vec2[T] {
	return ga.Vec2[T]{
		Len(a),
		ga.Atan2(a.Y, a.X),
	}
}

func PolarToCartesian[T ga.Signed](a ga.Vec2[T]) ga.Vec2[T] {
	r := a.X
	p := a.Y
	return ga.Vec2[T]{
		r * ga.Cos(p),
		r * ga.Sin(p),
	}
}

func Sin[T ga.Signed](a ga.Vec2[T]) ga.Vec2[T] {
	return ga.Vec2[T]{
		ga.Sin(a.X),
		ga.Sin(a.Y),
	}
}

func Cos[T ga.Signed](a ga.Vec2[T]) ga.Vec2[T] {
	return ga.Vec2[T]{
		ga.Cos(a.X),
		ga.Cos(a.Y),
	}
}

func Tan[T ga.Signed](a ga.Vec2[T]) ga.Vec2[T] {
	return ga.Vec2[T]{
		ga.Tan(a.X),
		ga.Tan(a.Y),
	}
}

func Asin[T ga.Signed](a ga.Vec2[T]) ga.Vec2[T] {
	return ga.Vec2[T]{
		ga.Asin(a.X),
		ga.Asin(a.Y),
	}
}

func Acos[T ga.Signed](a ga.Vec2[T]) ga.Vec2[T] {
	return ga.Vec2[T]{
		ga.Acos(a.X),
		ga.Acos(a.Y),
	}
}

func Atan[T ga.Signed](a ga.Vec2[T]) ga.Vec2[T] {
	return ga.Vec2[T]{
		ga.Atan(a.X),
		ga.Atan(a.Y),
	}
}

func Exp[T ga.Signed](a ga.Vec2[T]) ga.Vec2[T] {
	return ga.Vec2[T]{
		ga.Exp(a.X),
		ga.Exp(a.Y),
	}
}
