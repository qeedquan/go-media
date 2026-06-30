package vec4

import (
	"math"

	"github.com/qeedquan/go-media/math/ga"
)

func Add[T ga.Number](a, b ga.Vec4[T]) ga.Vec4[T] {
	return ga.Vec4[T]{
		a.X + b.X,
		a.Y + b.Y,
		a.Z + b.Z,
		a.W + b.W,
	}
}

func Sub[T ga.Number](a, b ga.Vec4[T]) ga.Vec4[T] {
	return ga.Vec4[T]{
		a.X - b.X,
		a.Y - b.Y,
		a.Z - b.Z,
		a.W - b.W,
	}
}

func Len[T ga.Signed](a ga.Vec4[T]) T {
	return T(math.Sqrt(float64(Dot(a, a))))
}

func Normalize[T ga.Signed](a ga.Vec4[T]) ga.Vec4[T] {
	l := Len(a)
	if l == 0 {
		return a
	}
	return Scale(a, 1/l)
}

func Proj[T ga.Signed](a, b ga.Vec4[T]) ga.Vec4[T] {
	b = Normalize(b)
	return Scale(b, Dot(a, b))
}

func Rej[T ga.Signed](a, b ga.Vec4[T]) ga.Vec4[T] {
	return Sub(a, Proj(a, b))
}

func Hadamard[T ga.Number](a, b ga.Vec4[T]) ga.Vec4[T] {
	return ga.Vec4[T]{
		a.X * b.X,
		a.Y * b.Y,
		a.Z * b.Z,
		a.W * b.W,
	}
}

func Scale[T ga.Number](a ga.Vec4[T], s T) ga.Vec4[T] {
	return ga.Vec4[T]{
		a.X * s,
		a.Y * s,
		a.Z * s,
		a.W * s,
	}
}

func Fmas[T ga.Number](a, b ga.Vec4[T], c T) ga.Vec4[T] {
	return Add(a, Scale(b, c))
}

func Fma[T ga.Number](a, b, c ga.Vec4[T]) ga.Vec4[T] {
	return Add(a, Hadamard(b, c))
}

func Neg[T ga.Number](a ga.Vec4[T]) ga.Vec4[T] {
	return ga.Vec4[T]{-a.X, -a.Y, -a.Z, -a.W}
}

func Dot[T ga.Number](a, b ga.Vec4[T]) T {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z + a.W*b.W
}

func Vec3[T ga.Number](a ga.Vec4[T]) ga.Vec3[T] {
	return ga.Vec3[T]{a.X, a.Y, a.Z}
}

func Element[T ga.Number](a ga.Vec4[T]) (x, y, z, w T) {
	return a.X, a.Y, a.Z, a.W
}

func Fill[T ga.Number](v T) ga.Vec4[T] {
	return ga.Vec4[T]{v, v, v, v}
}

func Distance[T ga.Signed](a, b ga.Vec4[T]) T {
	return Len(Sub(a, b))
}

func Angle[T ga.Signed](a, b ga.Vec4[T]) T {
	return ga.Acos(Dot(a, b) / (Len(a) * Len(b)))
}

func Lerp[T ga.Number](t T, a, b ga.Vec4[T]) ga.Vec4[T] {
	return ga.Vec4[T]{
		ga.Lerp(t, a.X, b.X),
		ga.Lerp(t, a.Y, b.Y),
		ga.Lerp(t, a.Z, b.Z),
		ga.Lerp(t, a.W, b.W),
	}
}

func Unlerp[T ga.Number](t T, a, b ga.Vec4[T]) ga.Vec4[T] {
	return ga.Vec4[T]{
		ga.Unlerp(t, a.X, b.X),
		ga.Unlerp(t, a.Y, b.Y),
		ga.Unlerp(t, a.Z, b.Z),
		ga.Unlerp(t, a.W, b.W),
	}
}

func Round[T ga.Signed](a ga.Vec4[T]) ga.Vec4[T] {
	return ga.Vec4[T]{
		ga.Round(a.X),
		ga.Round(a.Y),
		ga.Round(a.Z),
		ga.Round(a.W),
	}
}

func Sign[T ga.Signed](a ga.Vec4[T]) ga.Vec4[T] {
	return ga.Vec4[T]{
		ga.Sign(a.X),
		ga.Sign(a.Y),
		ga.Sign(a.Z),
		ga.Sign(a.W),
	}
}

func Sin[T ga.Signed](a ga.Vec4[T]) ga.Vec4[T] {
	return ga.Vec4[T]{
		ga.Sin(a.X),
		ga.Sin(a.Y),
		ga.Sin(a.Z),
		ga.Sin(a.W),
	}
}

func Cos[T ga.Signed](a ga.Vec4[T]) ga.Vec4[T] {
	return ga.Vec4[T]{
		ga.Cos(a.X),
		ga.Cos(a.Y),
		ga.Cos(a.Z),
		ga.Cos(a.W),
	}
}

func Tan[T ga.Signed](a ga.Vec4[T]) ga.Vec4[T] {
	return ga.Vec4[T]{
		ga.Tan(a.X),
		ga.Tan(a.Y),
		ga.Tan(a.Z),
		ga.Tan(a.W),
	}
}

func Asin[T ga.Signed](a ga.Vec4[T]) ga.Vec4[T] {
	return ga.Vec4[T]{
		ga.Asin(a.X),
		ga.Asin(a.Y),
		ga.Asin(a.Z),
		ga.Asin(a.W),
	}
}

func Acos[T ga.Signed](a ga.Vec4[T]) ga.Vec4[T] {
	return ga.Vec4[T]{
		ga.Acos(a.X),
		ga.Acos(a.Y),
		ga.Acos(a.Z),
		ga.Acos(a.W),
	}
}

func Atan[T ga.Signed](a ga.Vec4[T]) ga.Vec4[T] {
	return ga.Vec4[T]{
		ga.Atan(a.X),
		ga.Atan(a.Y),
		ga.Atan(a.Z),
		ga.Atan(a.W),
	}
}

func Exp[T ga.Signed](a ga.Vec4[T]) ga.Vec4[T] {
	return ga.Vec4[T]{
		ga.Exp(a.X),
		ga.Exp(a.Y),
		ga.Exp(a.Z),
		ga.Exp(a.W),
	}
}
