package mat2

import (
	"math/cmplx"

	"github.com/qeedquan/go-media/math/ga"
	"golang.org/x/exp/constraints"
)

func Rotation[T ga.Signed](t T) ga.Mat2[T] {
	s, c := ga.Sincos(t)
	return ga.Mat2[T]{
		{c, -s},
		{s, c},
	}
}

func Shear[T ga.Number](x, y T) ga.Mat2[T] {
	return ga.Mat2[T]{
		{1, x},
		{y, 1},
	}
}

func Scale[T ga.Number](x, y T) ga.Mat2[T] {
	return ga.Mat2[T]{
		{x, 0},
		{0, y},
	}
}

func Diagonal[T ga.Number](v T) ga.Mat2[T] {
	return ga.Mat2[T]{
		{v, 0},
		{0, v},
	}
}

func Muls[T ga.Number](m *ga.Mat2[T], s T) {
	for i := range m {
		for j := range m[i] {
			m[i][j] *= s
		}
	}
}

func Add[T ga.Number](m, a, b *ga.Mat2[T]) *ga.Mat2[T] {
	for i := range m {
		for j := range m[i] {
			m[i][j] = a[i][j] + b[i][j]
		}
	}
	return m
}

func Sub[T ga.Number](m, a, b *ga.Mat2[T]) *ga.Mat2[T] {
	for i := range m {
		for j := range m[i] {
			m[i][j] = a[i][j] - b[i][j]
		}
	}
	return m
}

func Mul[T ga.Number](m, a, b *ga.Mat2[T]) *ga.Mat2[T] {
	var p ga.Mat2[T]
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

func Trace[T ga.Number](m *ga.Mat2[T]) T {
	return m[0][0] + m[1][1]
}

func Det[T ga.Number](m *ga.Mat2[T]) T {
	return m[0][0]*m[1][1] - m[0][1]*m[1][0]
}

func Adj[T ga.Number](r, m *ga.Mat2[T]) *ga.Mat2[T] {
	m00 := m[0][0]
	m01 := m[0][1]
	m10 := m[1][0]
	m11 := m[1][1]

	r[0][0] = m11
	r[0][1] = -m01
	r[1][0] = -m10
	r[1][1] = m00

	return r
}

func Inv[T ga.Number](r, m *ga.Mat2[T]) T {
	var p ga.Mat2[T]
	Adj(&p, m)
	d := Det(m)
	if d == 0 {
		return d
	}
	Muls(&p, 1/d)
	*r = p
	return d
}

func Transpose[T ga.Number](r, m *ga.Mat2[T]) *ga.Mat2[T] {
	var p ga.Mat2[T]
	for i := range m {
		for j := range m[i] {
			p[j][i] = m[i][j]
		}
	}
	*m = p
	return m
}

func Apply[T ga.Number](m *ga.Mat2[T], v ga.Vec2[T]) ga.Vec2[T] {
	return ga.Vec2[T]{
		m[0][0]*v.X + m[0][1]*v.Y,
		m[1][0]*v.X + m[1][1]*v.Y,
	}
}

func Applyc[T ga.Number](m *ga.Mat2[T], v ga.Vec2[T]) ga.Vec2[T] {
	return ga.Vec2[T]{
		v.X*m[0][0] + v.Y*m[1][0],
		v.X*m[0][1] + v.Y*m[1][1],
	}
}

func Basis[T ga.Number](x, y ga.Vec2[T]) ga.Mat2[T] {
	return ga.Mat2[T]{
		{x.X, y.X},
		{x.Y, y.Y},
	}
}

func Eigc[T constraints.Complex](m *ga.Mat2[T]) (ev [2]T, ew [2]ga.Vec2[T]) {
	ev = ga.Roots2c(1, -Trace(m), Det(m))
	ew[0] = eigv(m, ev[0])
	ew[1] = eigv(m, ev[1])
	return
}

func eigv[T constraints.Complex](m *ga.Mat2[T], ev T) ga.Vec2[T] {
	a, b := m[0][0], m[0][1]
	if b == 0 {
		a, b = m[1][0], m[1][1]
	}

	t := cmplx.Atan(complex128((ev - a) / b))
	return ga.Vec2[T]{
		T(cmplx.Cos(t)),
		T(cmplx.Sin(t)),
	}
}
