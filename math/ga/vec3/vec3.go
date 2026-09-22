package vec3

import (
	"math"

	"github.com/qeedquan/go-media/math/ga"
)

func Add[T ga.Number](a, b ga.Vec3[T]) ga.Vec3[T] {
	return ga.Vec3[T]{
		a.X + b.X,
		a.Y + b.Y,
		a.Z + b.Z,
	}
}

func Sub[T ga.Number](a, b ga.Vec3[T]) ga.Vec3[T] {
	return ga.Vec3[T]{
		a.X - b.X,
		a.Y - b.Y,
		a.Z - b.Z,
	}
}

func Len[T ga.Signed](a ga.Vec3[T]) T {
	return T(math.Sqrt(float64(Dot(a, a))))
}

func Normalize[T ga.Signed](a ga.Vec3[T]) ga.Vec3[T] {
	l := Len(a)
	if l == 0 {
		return a
	}
	return Scale(a, 1/l)
}

func Proj[T ga.Signed](a, b ga.Vec3[T]) ga.Vec3[T] {
	b = Normalize(b)
	return Scale(b, Dot(a, b))
}

func Rej[T ga.Signed](a, b ga.Vec3[T]) ga.Vec3[T] {
	return Sub(a, Proj(a, b))
}

func Hadamard[T ga.Number](a, b ga.Vec3[T]) ga.Vec3[T] {
	return ga.Vec3[T]{
		a.X * b.X,
		a.Y * b.Y,
		a.Z * b.Z,
	}
}

func Scale[T ga.Number](a ga.Vec3[T], s T) ga.Vec3[T] {
	return ga.Vec3[T]{
		a.X * s,
		a.Y * s,
		a.Z * s,
	}
}

func Fmas[T ga.Number](a, b ga.Vec3[T], c T) ga.Vec3[T] {
	return Add(a, Scale(b, c))
}

func Fma[T ga.Number](a, b, c ga.Vec3[T]) ga.Vec3[T] {
	return Add(a, Hadamard(b, c))
}

func Neg[T ga.Number](a ga.Vec3[T]) ga.Vec3[T] {
	return ga.Vec3[T]{-a.X, -a.Y, -a.Z}
}

func Dot[T ga.Number](a, b ga.Vec3[T]) T {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

func Cross[T ga.Number](a, b ga.Vec3[T]) ga.Vec3[T] {
	return ga.Vec3[T]{
		a.Y*b.Z - a.Z*b.Y,
		a.Z*b.X - a.X*b.Z,
		a.X*b.Y - a.Y*b.X,
	}
}

func Wedge[T ga.Number](a, b ga.Vec3[T]) ga.Vec3[T] {
	return ga.Vec3[T]{
		a.X*b.Y - a.Y*b.X,
		a.X*b.Z - a.Z*b.X,
		a.Y*b.Z - a.Z*b.Y,
	}
}

func Reflect[T ga.Number](a, b ga.Vec3[T]) ga.Vec3[T] {
	s := 2 * Dot(b, a) / Dot(a, a)
	return Sub(b, Scale(a, s))
}

func Refract[T ga.Signed](a, b ga.Vec3[T], eta T) ga.Vec3[T] {
	x := Dot(a, b)
	k := 1 - eta*eta*(1-x*x)
	if k < 0 {
		return ga.Vec3[T]{}
	}
	u := Scale(a, eta)
	v := Scale(b, eta*x+ga.Sqrt(k))
	return Sub(u, v)
}

func Min[T ga.Ordinal](a ga.Vec3[T], b ...ga.Vec3[T]) ga.Vec3[T] {
	for i := range b {
		a.X = min(a.X, b[i].X)
		a.Y = min(a.Y, b[i].Y)
		a.Z = min(a.Z, b[i].Z)
	}
	return a
}

func Max[T ga.Ordinal](a ga.Vec3[T], b ...ga.Vec3[T]) ga.Vec3[T] {
	for i := range b {
		a.X = max(a.X, b[i].X)
		a.Y = max(a.Y, b[i].Y)
		a.Z = max(a.Z, b[i].Z)
	}
	return a
}

func Clamp[T ga.Ordinal](x, a, b ga.Vec3[T]) ga.Vec3[T] {
	return ga.Vec3[T]{
		ga.Clamp(x.X, a.X, b.X),
		ga.Clamp(x.Y, a.Y, b.Y),
		ga.Clamp(x.Z, a.Z, b.Z),
	}
}

func Lerp[T ga.Number](t T, a, b ga.Vec3[T]) ga.Vec3[T] {
	return ga.Vec3[T]{
		ga.Lerp(t, a.X, b.X),
		ga.Lerp(t, a.Y, b.Y),
		ga.Lerp(t, a.Z, b.Z),
	}
}

func Unlerp[T ga.Number](t T, a, b ga.Vec3[T]) ga.Vec3[T] {
	return ga.Vec3[T]{
		ga.Unlerp(t, a.X, b.X),
		ga.Unlerp(t, a.Y, b.Y),
		ga.Unlerp(t, a.Z, b.Z),
	}
}

func Round[T ga.Signed](a ga.Vec3[T]) ga.Vec3[T] {
	return ga.Vec3[T]{
		ga.Round(a.X),
		ga.Round(a.Y),
		ga.Round(a.Z),
	}
}

func Sign[T ga.Signed](a ga.Vec3[T]) ga.Vec3[T] {
	return ga.Vec3[T]{
		ga.Sign(a.X),
		ga.Sign(a.Y),
		ga.Sign(a.Z),
	}
}

func Saturate[T ga.Ordinal](a ga.Vec3[T]) ga.Vec3[T] {
	return ga.Vec3[T]{
		ga.Saturate(a.X),
		ga.Saturate(a.Y),
		ga.Saturate(a.Z),
	}
}

func MinComp[T ga.Ordinal](a ga.Vec3[T]) T {
	return min(a.X, a.Y, a.Z)
}

func MaxComp[T ga.Ordinal](a ga.Vec3[T]) T {
	return max(a.X, a.Y, a.Z)
}

func Vec2[T ga.Number](a ga.Vec3[T]) ga.Vec2[T] {
	return ga.Vec2[T]{a.X, a.Y}
}

func Vec4[T ga.Number](a ga.Vec3[T]) ga.Vec4[T] {
	return ga.Vec4[T]{a.X, a.Y, a.Z, 1}
}

func Element[T ga.Number](a ga.Vec3[T]) (x, y, z T) {
	return a.X, a.Y, a.Z
}

func Fill[T ga.Number](v T) ga.Vec3[T] {
	return ga.Vec3[T]{v, v, v}
}

func Distance[T ga.Signed](a, b ga.Vec3[T]) T {
	return Len(Sub(a, b))
}

func Angle[T ga.Signed](a, b ga.Vec3[T]) T {
	return ga.Acos(Dot(a, b) / (Len(a) * Len(b)))
}

func SphericalToCartesian[T ga.Signed](a ga.Vec3[T]) ga.Vec3[T] {
	return ga.Vec3[T]{
		Len(a),
		ga.Atan2(ga.Hypot(a.X, a.Y), a.Z),
		ga.Atan2(a.Y, a.X),
	}
}

func CartesianToSpherical[T ga.Signed](a ga.Vec3[T]) ga.Vec3[T] {
	r := a.X
	t := a.Y
	p := a.Z
	st, ct := ga.Sincos(t)
	sp, cp := ga.Sincos(p)
	return ga.Vec3[T]{
		r * cp * st,
		r * sp * st,
		r * ct,
	}
}

func Sin[T ga.Signed](a ga.Vec3[T]) ga.Vec3[T] {
	return ga.Vec3[T]{
		ga.Sin(a.X),
		ga.Sin(a.Y),
		ga.Sin(a.Z),
	}
}

func Cos[T ga.Signed](a ga.Vec3[T]) ga.Vec3[T] {
	return ga.Vec3[T]{
		ga.Cos(a.X),
		ga.Cos(a.Y),
		ga.Cos(a.Z),
	}
}

func Tan[T ga.Signed](a ga.Vec3[T]) ga.Vec3[T] {
	return ga.Vec3[T]{
		ga.Tan(a.X),
		ga.Tan(a.Y),
		ga.Tan(a.Z),
	}
}
func Asin[T ga.Signed](a ga.Vec3[T]) ga.Vec3[T] {
	return ga.Vec3[T]{
		ga.Asin(a.X),
		ga.Asin(a.Y),
		ga.Asin(a.Z),
	}
}

func Acos[T ga.Signed](a ga.Vec3[T]) ga.Vec3[T] {
	return ga.Vec3[T]{
		ga.Acos(a.X),
		ga.Acos(a.Y),
		ga.Acos(a.Z),
	}
}

func Atan[T ga.Signed](a ga.Vec3[T]) ga.Vec3[T] {
	return ga.Vec3[T]{
		ga.Atan(a.X),
		ga.Atan(a.Y),
		ga.Atan(a.Z),
	}
}

func Exp[T ga.Signed](a ga.Vec3[T]) ga.Vec3[T] {
	return ga.Vec3[T]{
		ga.Exp(a.X),
		ga.Exp(a.Y),
		ga.Exp(a.Z),
	}
}
