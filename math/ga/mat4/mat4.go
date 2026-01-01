package mat4

import (
	"github.com/qeedquan/go-media/math/ga"
	"github.com/qeedquan/go-media/math/ga/vec3"
)

func Rotation[T ga.Signed](v ga.Vec3[T], r T) ga.Mat4[T] {
	s, c := ga.Sincos(r)
	w := vec3.Normalize(v)

	m := ga.Mat4[T]{}
	m[0][0] = c + w.X*w.X*(1-c)
	m[0][1] = w.X*w.Y*(1-c) - w.Z*s
	m[0][2] = w.Y*s + w.X*w.Z*(1-c)
	m[1][0] = w.Z*s + w.X*w.Y*(1-c)
	m[1][1] = c + w.Y*w.Y*(1-c)
	m[1][2] = -w.X*s + w.Y*w.Z*(1-c)
	m[2][0] = -w.Y*s + w.X*w.Z*(1-c)
	m[2][1] = w.X*s + w.Y*w.Z*(1-c)
	m[2][2] = c + w.Z*w.Z*(1-c)
	m[3][3] = 1
	return m
}

func Scale[T ga.Number](x, y, z T) ga.Mat4[T] {
	return ga.Mat4[T]{
		{x, 0, 0, 0},
		{0, y, 0, 0},
		{0, 0, z, 0},
		{0, 0, 0, 1},
	}
}

func Translation[T ga.Number](x, y, z T) ga.Mat4[T] {
	return ga.Mat4[T]{
		{1, 0, 0, x},
		{0, 1, 0, y},
		{0, 0, 1, z},
		{0, 0, 0, 1},
	}
}

func Diagonal[T ga.Number](v T) ga.Mat4[T] {
	return ga.Mat4[T]{
		{v, 0, 0, 0},
		{0, v, 0, 0},
		{0, 0, v, 0},
		{0, 0, 0, v},
	}
}

func LookAt[T ga.Signed](eye, center, up ga.Vec3[T]) ga.Mat4[T] {
	f := vec3.Normalize(vec3.Sub(center, eye))
	up = vec3.Normalize(up)
	s := vec3.Normalize(vec3.Cross(f, up))
	u := vec3.Cross(s, f)

	var m ga.Mat4[T]
	m[0][0] = s.X
	m[0][1] = s.Y
	m[0][2] = s.Z
	m[1][0] = u.X
	m[1][1] = u.Y
	m[1][2] = u.Z
	m[2][0] = -f.X
	m[2][1] = -f.Y
	m[2][2] = -f.Z
	m[3][3] = 1

	m[0][3] = -eye.X*s.X - eye.Y*s.Y - eye.Z*s.Z
	m[1][3] = -eye.X*u.X - eye.Y*u.Y - eye.Z*u.Z
	m[2][3] = eye.X*f.X + eye.Y*f.Y + eye.Z*f.Z

	return m
}

func Frustum[T ga.Signed](l, r, b, t, n, f T) ga.Mat4[T] {
	A := (r + l) / (r - l)
	B := (t + b) / (t - b)
	C := -(f + n) / (f - n)
	D := -2 * f * n / (f - n)
	return ga.Mat4[T]{
		{2 * n / (r - l), 0, A, 0},
		{0, 2 * n / (t - b), B, 0},
		{0, 0, C, D},
		{0, 0, -1, 0},
	}
}

func InfinitePerspective[T ga.Signed](fovy, aspect, near T) ga.Mat4[T] {
	const zp = 0

	f := 1 / ga.Tan(fovy/2)
	return ga.Mat4[T]{
		{f / aspect, 0, 0, 0},
		{0, f, 0, 0},
		{0, 0, -(1 - zp), -near * (1 - zp)},
		{0, 0, -1, 0},
	}
}

func Perspective[T ga.Signed](fovy, aspect, near, far T) ga.Mat4[T] {
	ymax := near * ga.Tan(fovy/2)
	xmax := ymax * aspect
	return Frustum[T](-xmax, xmax, -ymax, ymax, near, far)
}

func Ortho[T ga.Signed](l, r, b, t, n, f T) ga.Mat4[T] {
	sx := 2 / (r - l)
	sy := 2 / (t - b)
	sz := -2 / (f - n)

	tx := -(r + l) / (r - l)
	ty := -(t + b) / (t - b)
	tz := -(f + n) / (f - n)

	return ga.Mat4[T]{
		{sx, 0, 0, tx},
		{0, sy, 0, ty},
		{0, 0, sz, tz},
		{0, 0, 0, 1},
	}
}

func Viewport[T ga.Signed](x, y, w, h T) ga.Mat4[T] {
	l := x
	b := y
	r := x + w
	t := y + h
	z := T(1)
	return ga.Mat4[T]{
		{(r - l) / 2, 0, 0, (r + l) / 2},
		{0, (t - b) / 2, 0, (t + b) / 2},
		{0, 0, z / 2, z / 2},
		{0, 0, 0, 1},
	}
}

func Add[T ga.Number](m, a, b *ga.Mat4[T]) *ga.Mat4[T] {
	for i := range m {
		for j := range m[i] {
			m[i][j] = a[i][j] + b[i][j]
		}
	}
	return m
}

func Sub[T ga.Number](m, a, b *ga.Mat4[T]) *ga.Mat4[T] {
	for i := range m {
		for j := range m[i] {
			m[i][j] = a[i][j] - b[i][j]
		}
	}
	return m
}

func Mul[T ga.Number](m, a, b *ga.Mat4[T]) *ga.Mat4[T] {
	var p ga.Mat4[T]
	for i := range a {
		for j := range a[i] {
			for k := range a[j] {
				p[i][j] += a[i][k] * b[k][j]
			}
		}
	}
	*m = p
	return m
}

func Muls[T ga.Number](m *ga.Mat4[T], s T) {
	for i := range m {
		for j := range m[i] {
			m[i][j] *= s
		}
	}
}

func Trace[T ga.Number](m *ga.Mat4[T]) T {
	return m[0][0] + m[1][1] + m[2][2] + m[3][3]
}

func Det[T ga.Number](m *ga.Mat4[T]) T {
	return m[0][3]*m[1][2]*m[2][1]*m[3][0] - m[0][2]*m[1][3]*m[2][1]*m[3][0] -
		m[0][3]*m[1][1]*m[2][2]*m[3][0] + m[0][1]*m[1][3]*m[2][2]*m[3][0] +
		m[0][2]*m[1][1]*m[2][3]*m[3][0] - m[0][1]*m[1][2]*m[2][3]*m[3][0] -
		m[0][3]*m[1][2]*m[2][0]*m[3][1] + m[0][2]*m[1][3]*m[2][0]*m[3][1] +
		m[0][3]*m[1][0]*m[2][2]*m[3][1] - m[0][0]*m[1][3]*m[2][2]*m[3][1] -
		m[0][2]*m[1][0]*m[2][3]*m[3][1] + m[0][0]*m[1][2]*m[2][3]*m[3][1] +
		m[0][3]*m[1][1]*m[2][0]*m[3][2] - m[0][1]*m[1][3]*m[2][0]*m[3][2] -
		m[0][3]*m[1][0]*m[2][1]*m[3][2] + m[0][0]*m[1][3]*m[2][1]*m[3][2] +
		m[0][1]*m[1][0]*m[2][3]*m[3][2] - m[0][0]*m[1][1]*m[2][3]*m[3][2] -
		m[0][2]*m[1][1]*m[2][0]*m[3][3] + m[0][1]*m[1][2]*m[2][0]*m[3][3] +
		m[0][2]*m[1][0]*m[2][1]*m[3][3] - m[0][0]*m[1][2]*m[2][1]*m[3][3] -
		m[0][1]*m[1][0]*m[2][2]*m[3][3] + m[0][0]*m[1][1]*m[2][2]*m[3][3]
}

func Adj[T ga.Number](r, m *ga.Mat4[T]) *ga.Mat4[T] {
	m00 := m[0][0]
	m01 := m[0][1]
	m02 := m[0][2]
	m03 := m[0][3]
	m10 := m[1][0]
	m11 := m[1][1]
	m12 := m[1][2]
	m13 := m[1][3]
	m20 := m[2][0]
	m21 := m[2][1]
	m22 := m[2][2]
	m23 := m[2][3]
	m30 := m[3][0]
	m31 := m[3][1]
	m32 := m[3][2]
	m33 := m[3][3]

	r[0][0] = m11*(m22*m33-m23*m32) + m21*(m13*m32-m12*m33) + m31*(m12*m23-m13*m22)
	r[0][1] = m01*(m23*m32-m22*m33) + m21*(m02*m33-m03*m32) + m31*(m03*m22-m02*m23)
	r[0][2] = m01*(m12*m33-m13*m32) + m11*(m03*m32-m02*m33) + m31*(m02*m13-m03*m12)
	r[0][3] = m01*(m13*m22-m12*m23) + m11*(m02*m23-m03*m22) + m21*(m03*m12-m02*m13)
	r[1][0] = m10*(m23*m32-m22*m33) + m20*(m12*m33-m13*m32) + m30*(m13*m22-m12*m23)
	r[1][1] = m00*(m22*m33-m23*m32) + m20*(m03*m32-m02*m33) + m30*(m02*m23-m03*m22)
	r[1][2] = m00*(m13*m32-m12*m33) + m10*(m02*m33-m03*m32) + m30*(m03*m12-m02*m13)
	r[1][3] = m00*(m12*m23-m13*m22) + m10*(m03*m22-m02*m23) + m20*(m02*m13-m03*m12)
	r[2][0] = m10*(m21*m33-m23*m31) + m20*(m13*m31-m11*m33) + m30*(m11*m23-m13*m21)
	r[2][1] = m00*(m23*m31-m21*m33) + m20*(m01*m33-m03*m31) + m30*(m03*m21-m01*m23)
	r[2][2] = m00*(m11*m33-m13*m31) + m10*(m03*m31-m01*m33) + m30*(m01*m13-m03*m11)
	r[2][3] = m00*(m13*m21-m11*m23) + m10*(m01*m23-m03*m21) + m20*(m03*m11-m01*m13)
	r[3][0] = m10*(m22*m31-m21*m32) + m20*(m11*m32-m12*m31) + m30*(m12*m21-m11*m22)
	r[3][1] = m00*(m21*m32-m22*m31) + m20*(m02*m31-m01*m32) + m30*(m01*m22-m02*m21)
	r[3][2] = m00*(m12*m31-m11*m32) + m10*(m01*m32-m02*m31) + m30*(m02*m11-m01*m12)
	r[3][3] = m00*(m11*m22-m12*m21) + m10*(m02*m21-m01*m22) + m20*(m01*m12-m02*m11)

	return r
}

func Inv[T ga.Number](r, m *ga.Mat4[T]) T {
	var p ga.Mat4[T]
	Adj(&p, m)
	d := Det(m)
	if d == 0 {
		return d
	}
	Muls(&p, 1/d)
	*r = p
	return d
}

func Transpose[T ga.Number](r, m *ga.Mat4[T]) *ga.Mat4[T] {
	var p ga.Mat4[T]
	for i := range m {
		for j := range m[i] {
			p[j][i] = m[i][j]
		}
	}
	*m = p
	return m
}

func Apply[T ga.Number](m *ga.Mat4[T], v ga.Vec4[T]) ga.Vec4[T] {
	return ga.Vec4[T]{
		m[0][0]*v.X + m[0][1]*v.Y + m[0][2]*v.Z + m[0][3]*v.W,
		m[1][0]*v.X + m[1][1]*v.Y + m[1][2]*v.Z + m[1][3]*v.W,
		m[2][0]*v.X + m[2][1]*v.Y + m[2][2]*v.Z + m[2][3]*v.W,
		m[3][0]*v.X + m[3][1]*v.Y + m[3][2]*v.Z + m[3][3]*v.W,
	}
}

func Applyc[T ga.Number](m *ga.Mat4[T], v ga.Vec4[T]) ga.Vec4[T] {
	return ga.Vec4[T]{
		v.X*m[0][0] + v.Y*m[1][0] + v.Z*m[2][0] + v.W*m[3][0],
		v.X*m[0][1] + v.Y*m[1][1] + v.Z*m[2][1] + v.W*m[3][1],
		v.X*m[0][2] + v.Y*m[1][2] + v.Z*m[2][2] + v.W*m[3][2],
		v.X*m[0][3] + v.Y*m[1][3] + v.Z*m[2][3] + v.W*m[3][3],
	}
}

func Apply3[T ga.Number](m *ga.Mat4[T], v ga.Vec3[T]) ga.Vec3[T] {
	return ga.Vec3[T]{
		m[0][0]*v.X + m[0][1]*v.Y + m[0][2]*v.Z + m[0][3],
		m[1][0]*v.X + m[1][1]*v.Y + m[1][2]*v.Z + m[1][3],
		m[2][0]*v.X + m[2][1]*v.Y + m[2][2]*v.Z + m[2][3],
	}
}

func Apply3c[T ga.Number](m *ga.Mat4[T], v ga.Vec3[T]) ga.Vec3[T] {
	return ga.Vec3[T]{
		v.X*m[0][0] + v.Y*m[1][0] + v.Z*m[2][0] + m[3][0],
		v.X*m[0][1] + v.Y*m[1][1] + v.Z*m[2][1] + m[3][1],
		v.X*m[0][2] + v.Y*m[1][2] + v.Z*m[2][2] + m[3][2],
	}
}

func Basis[T ga.Number](x, y, z, w ga.Vec4[T]) ga.Mat4[T] {
	return ga.Mat4[T]{
		{x.X, y.X, z.X, w.X},
		{x.Y, y.Y, z.Y, w.Y},
		{x.Z, y.Z, z.Z, w.Z},
		{x.W, y.W, z.W, w.W},
	}
}
