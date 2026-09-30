package mat3

import (
	"github.com/qeedquan/go-media/math/ga"
	"github.com/qeedquan/go-media/math/ga/vec3"
)

func Rotation[T ga.Signed](v ga.Vec3[T], r T) ga.Mat3[T] {
	s, c := ga.Sincos(r)
	w := vec3.Normalize(v)

	m := ga.Mat3[T]{}
	m[0][0] = c + w.X*w.X*(1-c)
	m[0][1] = w.X*w.Y*(1-c) - w.Z*s
	m[0][2] = w.Y*s + w.X*w.Z*(1-c)
	m[1][0] = w.Z*s + w.X*w.Y*(1-c)
	m[1][1] = c + w.Y*w.Y*(1-c)
	m[1][2] = -w.X*s + w.Y*w.Z*(1-c)
	m[2][0] = -w.Y*s + w.X*w.Z*(1-c)
	m[2][1] = w.X*s + w.Y*w.Z*(1-c)
	m[2][2] = c + w.Z*w.Z*(1-c)
	return m
}

func Scale[T ga.Number](x, y, z T) ga.Mat3[T] {
	return ga.Mat3[T]{
		{x, 0, 0},
		{0, y, 0},
		{0, 0, z},
	}
}

func Translation[T ga.Number](x, y T) ga.Mat3[T] {
	return ga.Mat3[T]{
		{1, 0, x},
		{0, 1, y},
		{0, 0, 1},
	}
}

func Diagonal[T ga.Number](v T) ga.Mat3[T] {
	return ga.Mat3[T]{
		{v, 0, 0},
		{0, v, 0},
		{0, 0, v},
	}
}

func Muls[T ga.Number](m *ga.Mat3[T], s T) {
	for i := range m {
		for j := range m[i] {
			m[i][j] *= s
		}
	}
}

func Add[T ga.Number](m, a, b *ga.Mat3[T]) *ga.Mat3[T] {
	for i := range m {
		for j := range m[i] {
			m[i][j] = a[i][j] + b[i][j]
		}
	}
	return m
}

func Sub[T ga.Number](m, a, b *ga.Mat3[T]) *ga.Mat3[T] {
	for i := range m {
		for j := range m[i] {
			m[i][j] = a[i][j] - b[i][j]
		}
	}
	return m
}

func Mul[T ga.Number](m, a, b *ga.Mat3[T]) *ga.Mat3[T] {
	var p ga.Mat3[T]
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

func Trace[T ga.Number](m *ga.Mat3[T]) T {
	return m[0][0] + m[1][1] + m[2][2]
}

func Det[T ga.Number](m *ga.Mat3[T]) T {
	m00 := m[0][0]
	m01 := m[0][1]
	m02 := m[0][2]
	m10 := m[1][0]
	m11 := m[1][1]
	m12 := m[1][2]
	m20 := m[2][0]
	m21 := m[2][1]
	m22 := m[2][2]
	c00 := m11*m22 - m12*m21
	c01 := m12*m20 - m10*m22
	c02 := m10*m21 - m11*m20
	d := m00*c00 + m01*c01 + m02*c02
	return d

}

func Adj[T ga.Number](r, m *ga.Mat3[T]) *ga.Mat3[T] {
	m00 := m[0][0]
	m01 := m[0][1]
	m02 := m[0][2]
	m10 := m[1][0]
	m11 := m[1][1]
	m12 := m[1][2]
	m20 := m[2][0]
	m21 := m[2][1]
	m22 := m[2][2]

	r[0][0] = m11*m22 - m12*m21
	r[0][1] = m02*m21 - m01*m22
	r[0][2] = m01*m12 - m02*m11
	r[1][0] = m12*m20 - m10*m22
	r[1][1] = m00*m22 - m02*m20
	r[1][2] = m02*m10 - m00*m12
	r[2][0] = m10*m21 - m11*m20
	r[2][1] = m01*m20 - m00*m21
	r[2][2] = m00*m11 - m01*m10

	return r
}

func Inv[T ga.Number](r, m *ga.Mat3[T]) T {
	var p ga.Mat3[T]
	Adj(&p, m)
	d := Det(m)
	if d == 0 {
		return d
	}
	Muls(&p, 1/d)
	*r = p
	return d
}

func Transpose[T ga.Number](r, m *ga.Mat3[T]) *ga.Mat3[T] {
	var p ga.Mat3[T]
	for i := range m {
		for j := range m[i] {
			p[j][i] = m[i][j]
		}
	}
	*m = p
	return m
}

func Apply[T ga.Number](m *ga.Mat3[T], v ga.Vec3[T]) ga.Vec3[T] {
	return ga.Vec3[T]{
		m[0][0]*v.X + m[0][1]*v.Y + m[0][2]*v.Z,
		m[1][0]*v.X + m[1][1]*v.Y + m[1][2]*v.Z,
		m[2][0]*v.X + m[2][1]*v.Y + m[2][2]*v.Z,
	}
}

func Applyc[T ga.Number](m *ga.Mat3[T], v ga.Vec3[T]) ga.Vec3[T] {
	return ga.Vec3[T]{
		v.X*m[0][0] + v.Y*m[1][0] + v.Z*m[2][0],
		v.X*m[0][1] + v.Y*m[1][1] + v.Z*m[2][1],
		v.X*m[0][2] + v.Y*m[1][2] + v.Z*m[2][2],
	}
}

func Basis[T ga.Number](x, y, z ga.Vec3[T]) ga.Mat3[T] {
	return ga.Mat3[T]{
		{x.X, y.X, z.X},
		{x.Y, y.Y, z.Y},
		{x.Z, y.Z, z.Z},
	}
}
