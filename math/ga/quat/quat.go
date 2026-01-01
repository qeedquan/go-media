package quat

import (
	"math"

	"github.com/qeedquan/go-media/math/ga"
	"github.com/qeedquan/go-media/math/ga/vec3"
	"golang.org/x/exp/constraints"
)

func Add[T ga.Number](a, b ga.Quat[T]) ga.Quat[T] {
	return ga.Quat[T]{
		a.X + b.X,
		a.Y + b.Y,
		a.Z + b.Z,
		a.W + b.W,
	}
}

func Sub[T ga.Number](a, b ga.Quat[T]) ga.Quat[T] {
	return ga.Quat[T]{
		a.X - b.X,
		a.Y - b.Y,
		a.Z - b.Z,
		a.W - b.W,
	}
}

func Neg[T ga.Signed](a ga.Quat[T]) ga.Quat[T] {
	return ga.Quat[T]{-a.X, -a.Y, -a.Z, -a.W}
}

func Dot[T ga.Number](a, b ga.Quat[T]) T {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z + a.W*b.W
}

func Mul[T ga.Number](a, b ga.Quat[T]) ga.Quat[T] {
	return ga.Quat[T]{
		a.X*b.W + a.W*b.X + a.Y*b.Z - a.Z*b.Y,
		a.Y*b.W + a.W*b.Y + a.Z*b.X - a.X*b.Z,
		a.Z*b.W + a.W*b.Z + a.X*b.Y - a.Y*b.X,
		a.W*b.W - a.X*b.X - a.Y*b.Y - a.Z*b.Z,
	}
}

func Mulv[T ga.Number](q ga.Quat[T], r ga.Vec4[T]) ga.Quat[T] {
	w := -q.X*r.X - q.Y*r.Y - q.Z*r.Z
	x := q.W*r.X + q.Y*r.Z - q.Z*r.Y
	y := q.W*r.Y + q.Z*r.X - q.X*r.Z
	z := q.W*r.Z + q.X*r.Y - q.Y*r.X
	return ga.Quat[T]{x, y, z, w}
}

func Scale[T ga.Number](a ga.Quat[T], s T) ga.Quat[T] {
	return ga.Quat[T]{
		a.X * s,
		a.Y * s,
		a.Z * s,
		a.W * s,
	}
}

func Len[T ga.Signed](a ga.Quat[T]) T {
	return ga.Sqrt(a.X*a.X + a.Y*a.Y + a.Z*a.Z + a.W*a.W)
}

func Normalize[T ga.Signed](a ga.Quat[T]) ga.Quat[T] {
	l := Len(a)
	if l == 0 {
		return ga.Quat[T]{}
	}
	return ga.Quat[T]{
		a.X / l,
		a.Y / l,
		a.Z / l,
		a.W / l,
	}
}

func Conj[T ga.Number](a ga.Quat[T]) ga.Quat[T] {
	return ga.Quat[T]{-a.X, -a.Y, -a.Z, a.W}
}

func Inv[T ga.Signed](a ga.Quat[T]) ga.Quat[T] {
	l2 := Dot(a, a)
	return ga.Quat[T]{
		-a.X / l2,
		-a.Y / l2,
		-a.Z / l2,
		a.W / l2,
	}
}

func Powu[T ga.Signed](a ga.Quat[T], p T) ga.Quat[T] {
	return Exp(Scale(Log(a), p))
}

func Exp[T ga.Signed](a ga.Quat[T]) ga.Quat[T] {
	l := Len(a)
	if l == 0 {
		return ga.Quat[T]{0, 0, 0, 1}
	}
	v := ga.Quat[T]{a.X / l, a.Y / l, a.Z / l, 0}

	s, c := ga.Sincos(l)
	return ga.Quat[T]{
		v.X * s,
		v.Y * s,
		v.Z * s,
		ga.Exp(a.W) * c,
	}
}

func Log[T ga.Signed](a ga.Quat[T]) ga.Quat[T] {
	q := Len(a)
	v := Len(ga.Quat[T]{a.X, a.Y, a.Z, 0})
	if q == 0 || v == 0 {
		return ga.Quat[T]{}
	}
	c := ga.Acos(a.W/q) / v

	return ga.Quat[T]{
		c * a.X,
		c * a.Y,
		c * a.Z,
		ga.Log(q),
	}
}

func AxisAngle[T ga.Signed](a ga.Vec3[T], r T) ga.Quat[T] {
	w := vec3.Normalize(a)
	s, c := ga.Sincos(r)
	return ga.Quat[T]{
		w.X * s,
		w.Y * s,
		w.Z * s,
		c,
	}
}

func Euler[T ga.Signed](a ga.Vec3[T]) ga.Quat[T] {
	y, p, r := a.X, a.Y, a.Z
	sy, cy := ga.Sincos(y / 2)
	sp, cp := ga.Sincos(p / 2)
	sr, cr := ga.Sincos(r / 2)
	return ga.Quat[T]{
		sr*cp*cy - cr*sp*sy,
		cr*sp*cy + sr*cp*sy,
		cr*cp*sy - sr*sp*cy,
		cr*cp*cy + sr*sp*sy,
	}
}

func Rotation[T ga.Signed](m *ga.Mat4[T]) ga.Quat[T] {
	x := 1 + m[0][0] - m[1][1] - m[2][2]
	y := 1 - m[0][0] + m[1][1] - m[2][2]
	z := 1 - m[0][0] - m[1][1] + m[2][2]
	w := 1 + m[0][0] + m[1][1] + m[2][2]
	if x < 0 {
		x = 0
	} else {
		x = ga.Sqrt(x) / 2
	}

	if y < 0 {
		y = 0
	} else {
		y = ga.Sqrt(y) / 2
	}

	if z < 0 {
		z = 0
	} else {
		z = ga.Sqrt(z) / 2
	}

	if w < 0 {
		w = 0
	} else {
		w = ga.Sqrt(w) / 2
	}

	if m[2][1]-m[1][2] < 0 {
		x = -x
	}
	if m[0][2]-m[2][0] < 0 {
		y = -y
	}
	if m[1][0]-m[0][1] < 0 {
		z = -z
	}

	return ga.Quat[T]{T(x), T(y), T(z), T(w)}
}

func Matrix[T ga.Signed](a ga.Quat[T]) ga.Mat4[T] {
	x2 := a.X * a.X
	y2 := a.Y * a.Y
	z2 := a.Z * a.Z
	xy := a.X * a.Y
	xz := a.X * a.Z
	yz := a.Y * a.Z
	wx := a.W * a.X
	wy := a.W * a.Y
	wz := a.W * a.Z
	var m ga.Mat4[T]
	m[0][0] = 1.0 - 2.0*(y2+z2)
	m[0][1] = 2.0 * (xy - wz)
	m[0][2] = 2.0 * (xz + wy)
	m[1][0] = 2.0 * (xy + wz)
	m[1][1] = 1.0 - 2.0*(x2+z2)
	m[1][2] = 2.0 * (yz - wx)
	m[2][0] = 2.0 * (xz - wy)
	m[2][1] = 2.0 * (yz + wx)
	m[2][2] = 1.0 - 2.0*(x2+y2)
	m[3][3] = 1
	return m
}

func Slerp[T constraints.Float](t T, a, b ga.Quat[T]) ga.Quat[T] {
	v0 := Normalize(a)
	v1 := Normalize(b)

	const threshold = 0.9995
	dot := Dot(v0, v1)
	if dot > threshold {
		return Normalize(Lerp(t, v0, v1))
	}

	if dot < 0 {
		v1 = Neg(v1)
		dot = -dot
	}

	dot = ga.Clamp(dot, -1, 1)
	theta0 := ga.Acos(dot)
	theta := theta0 * t

	v2 := Sub(v1, Scale(v0, dot))
	v2 = Normalize(v2)

	v3 := Scale(v0, ga.Cos(theta))
	v4 := Scale(v2, ga.Sin(theta))
	return Add(v3, v4)
}

func Lerp[T ga.Number](t T, a, b ga.Quat[T]) ga.Quat[T] {
	return ga.Quat[T]{
		ga.Lerp(t, a.X, b.X),
		ga.Lerp(t, a.Y, b.Y),
		ga.Lerp(t, a.Z, b.Z),
		ga.Lerp(t, a.W, b.W),
	}
}

func Unlerp[T ga.Number](t T, a, b ga.Quat[T]) ga.Quat[T] {
	return ga.Quat[T]{
		ga.Unlerp(t, a.X, b.X),
		ga.Unlerp(t, a.Y, b.Y),
		ga.Unlerp(t, a.Z, b.Z),
		ga.Unlerp(t, a.W, b.W),
	}
}

func Distance[T ga.Signed](a, b ga.Quat[T]) T {
	d := Dot(a, b)
	return ga.Acos(2*d*d - 1)
}

func Sqrt[T ga.Signed](a ga.Quat[T]) ga.Quat[T] {
	u := ga.Quat[T]{a.X, a.Y, a.Z, 0}
	l := Len(u)
	if l == 0 {
		return ga.Quat[T]{}
	}

	r := Len(a)
	t := ga.Acos(a.W / r)
	z := ga.Sqrt(r)

	s, c := ga.Sincos(t / 2)
	return ga.Quat[T]{
		z * s * u.X / l,
		z * s * u.Y / l,
		z * s * u.Z / l,
		z * c,
	}
}

func Angle[T constraints.Float](a ga.Quat[T]) ga.Vec3[T] {
	var r ga.Vec3[T]

	sy := 2 * (a.W*a.Z + a.X*a.Y)
	cy := 1 - 2*(a.Y*a.Y+a.Z*a.Z)
	r.X = ga.Atan2(sy, cy)

	sp := 2 * (a.W*a.Y - a.Z*a.X)
	if ga.Abs(sp) >= 1 {
		r.Y = ga.Copysign(math.Pi/2, sp)
	} else {
		r.Y = ga.Asin(sp)
	}

	sr := 2 * (a.W*a.X + a.Y*a.Z)
	cr := 1 - 2*(a.X*a.X+a.Y*a.Y)
	r.Z = ga.Atan2(sr, cr)

	return r
}
