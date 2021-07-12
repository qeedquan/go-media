package f32

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"math/cmplx"
	"math/rand"
	"sort"
)

type Vec2 struct {
	X, Y float32
}

func (p Vec2) Add(q Vec2) Vec2 {
	return Vec2{p.X + q.X, p.Y + q.Y}
}

func (p Vec2) Sub(q Vec2) Vec2 {
	return Vec2{p.X - q.X, p.Y - q.Y}
}

func (p Vec2) AddScale(q Vec2, t float32) Vec2 {
	return p.Add(q.Scale(t))
}

func (p Vec2) SubScale(q Vec2, t float32) Vec2 {
	return p.Sub(q.Scale(t))
}

func (p Vec2) Neg() Vec2 {
	return Vec2{-p.X, -p.Y}
}

func (p Vec2) Abs() Vec2 {
	return Vec2{Abs(p.X), Abs(p.Y)}
}

func (p Vec2) Dot(q Vec2) float32 {
	return p.X*q.X + p.Y*q.Y
}

func (p Vec2) Projection(q Vec2) Vec2 {
	qn := q.Normalize()
	p1 := p.Dot(qn)
	return qn.Scale(p1)
}

func (p Vec2) Rejection(q Vec2) Vec2 {
	return p.Sub(p.Projection(q))
}

func (p Vec2) Perp() Vec2 {
	return Vec2{-p.Y, p.X}
}

func (p Vec2) MinComp() float32 {
	return Min(p.X, p.Y)
}

func (p Vec2) MaxComp() float32 {
	return Max(p.X, p.Y)
}

func (p Vec2) Max(q Vec2) Vec2 {
	return Vec2{
		Max(p.X, q.X),
		Max(p.Y, q.Y),
	}
}

func (p Vec2) Min(q Vec2) Vec2 {
	return Vec2{
		Min(p.X, q.X),
		Min(p.Y, q.Y),
	}
}

func (p Vec2) Floor() Vec2 {
	return Vec2{
		Floor(p.X),
		Floor(p.Y),
	}
}

func (p Vec2) Ceil() Vec2 {
	return Vec2{
		Ceil(p.X),
		Ceil(p.Y),
	}
}

func (p Vec2) Len() float32 {
	return Sqrt(p.X*p.X + p.Y*p.Y)
}

func (p Vec2) LenSquared() float32 {
	return p.Dot(p)
}

func (p Vec2) Normalize() Vec2 {
	l := p.Len()
	if l == 0 {
		return Vec2{}
	}
	return Vec2{p.X / l, p.Y / l}
}

func (p Vec2) Scale2(q Vec2) Vec2 {
	return Vec2{p.X * q.X, p.Y * q.Y}
}

func (p Vec2) Scale(k float32) Vec2 {
	return Vec2{p.X * k, p.Y * k}
}

func (p Vec2) Shear(k float32) Vec2 {
	return Vec2{p.X + k*p.Y, p.Y + k*p.X}
}

func (p Vec2) Shearv(q Vec2) Vec2 {
	return Vec2{p.X + q.X*p.Y, p.Y + q.Y*p.X}
}

func (p Vec2) Lerp(t float32, q Vec2) Vec2 {
	return Vec2{
		Lerp(t, p.X, q.X),
		Lerp(t, p.Y, q.Y),
	}
}

func (p Vec2) Lerp2(t, q Vec2) Vec2 {
	return Vec2{
		Lerp(t.X, p.X, q.X),
		Lerp(t.Y, p.Y, q.Y),
	}
}

func (p Vec2) Wrap(s, e float32) Vec2 {
	p.X = Wrap(p.X, s, e)
	p.Y = Wrap(p.Y, s, e)
	return p
}

func (p Vec2) Distance(q Vec2) float32 {
	return p.Sub(q).Len()
}

func (p Vec2) DistanceSquared(q Vec2) float32 {
	r := p.Sub(q)
	return r.Dot(r)
}

func (p Vec2) Polar() Polar {
	return Polar{p.Len(), Atan2(p.Y, p.X)}
}

func (p Vec2) MinScalar(k float32) Vec2 {
	return Vec2{
		Min(p.X, k),
		Min(p.Y, k),
	}
}

func (p Vec2) MaxScalar(k float32) Vec2 {
	return Vec2{
		Max(p.X, k),
		Max(p.Y, k),
	}
}

func (p Vec2) In(r Rectangle) bool {
	return r.Min.X <= p.X && p.X < r.Max.X &&
		r.Min.Y <= p.Y && p.Y < r.Max.Y
}

func (p Vec2) Rotate(r float32) Vec2 {
	si, co := Sincos(r)
	return Vec2{
		p.X*co - p.Y*si,
		p.X*si + p.Y*co,
	}
}

func (p Vec2) RotateAround(q Vec2, r float32) Vec2 {
	p = p.Sub(q)
	p = p.Rotate(r)
	p = p.Add(q)
	return p
}

func (p Vec2) Shrink(k float32) Vec2 {
	return Vec2{p.X / k, p.Y / k}
}

func (p Vec2) Shrink2(q Vec2) Vec2 {
	return Vec2{p.X / q.X, p.Y / q.Y}
}

func (p Vec2) YX() Vec2 {
	return Vec2{p.Y, p.X}
}

func (p Vec2) Clamp(s, e float32) Vec2 {
	p.X = Clamp(p.X, s, e)
	p.Y = Clamp(p.Y, s, e)
	return p
}

func (p Vec2) Clamp2(s, e Vec2) Vec2 {
	p.X = Clamp(p.X, s.X, e.X)
	p.Y = Clamp(p.Y, s.Y, e.Y)
	return p
}

func (p Vec2) Equals(q Vec2, eps float32) bool {
	return Abs(p.X-q.X) <= eps && Abs(p.Y-q.Y) <= eps
}

func (p Vec2) OnLine(a, b Vec2) bool {
	sx := Min(a.X, b.X)
	sy := Min(a.Y, b.Y)
	ex := Max(a.X, b.X)
	ey := Max(a.Y, b.Y)
	return sx <= p.X && p.X <= ex &&
		sy <= p.Y && p.Y <= ey
}

func (p Vec2) Angle(q Vec2) float32 {
	a := p.Len()
	b := q.Len()
	d := p.Dot(q)
	return float32(math.Acos(float64(d / (a * b))))
}

func (p Vec2) Point3() Vec3 {
	return Vec3{p.X, p.Y, 1}
}

func (p Vec2) Vec3() Vec3 {
	return Vec3{p.X, p.Y, 0}
}

func (p Vec2) String() string {
	return fmt.Sprintf(`Vec2(%0.3f, %0.3f)`, p.X, p.Y)
}

type Vec3 struct {
	X, Y, Z float32
}

func (p Vec3) Add(q Vec3) Vec3 {
	return Vec3{p.X + q.X, p.Y + q.Y, p.Z + q.Z}
}

func (p Vec3) Sub(q Vec3) Vec3 {
	return Vec3{p.X - q.X, p.Y - q.Y, p.Z - q.Z}
}

func (p Vec3) AddScale(q Vec3, t float32) Vec3 {
	return p.Add(q.Scale(t))
}

func (p Vec3) SubScale(q Vec3, t float32) Vec3 {
	return p.Sub(q.Scale(t))
}

func (p Vec3) Dot(q Vec3) float32 {
	return p.X*q.X + p.Y*q.Y + p.Z*q.Z
}

func (p Vec3) Projection(q Vec3) Vec3 {
	qn := q.Normalize()
	p1 := p.Dot(qn)
	return qn.Scale(p1)
}

func (p Vec3) Rejection(q Vec3) Vec3 {
	return p.Sub(p.Projection(q))
}

func (p Vec3) Cross(q Vec3) Vec3 {
	return Vec3{
		p.Y*q.Z - p.Z*q.Y,
		p.Z*q.X - p.X*q.Z,
		p.X*q.Y - p.Y*q.X,
	}
}

func (p Vec3) CrossNormalize(q Vec3) Vec3 {
	return p.Cross(q).Normalize()
}

func (p Vec3) Neg() Vec3 {
	return Vec3{-p.X, -p.Y, -p.Z}
}

func (p Vec3) Reflect(q Vec3) Vec3 {
	q = q.Scale(2 * p.Dot(q))
	return p.Sub(q)
}

func (p Vec3) Refract(q Vec3, eta float32) Vec3 {
	x := p.Dot(q)
	k := 1 - eta*eta*(1-x*x)
	if k < 0 {
		return Vec3{}
	}
	a := q.Scale(eta)
	b := p.Scale(eta*x + Sqrt(k))
	return a.Sub(b)
}

func (p Vec3) Len() float32 {
	return Sqrt(p.X*p.X + p.Y*p.Y + p.Z*p.Z)
}

func (p Vec3) LenSquared() float32 {
	return p.Dot(p)
}

func (p Vec3) Scale3(q Vec3) Vec3 {
	return Vec3{p.X * q.X, p.Y * q.Y, p.Z * q.Z}
}

func (p Vec3) Scale(k float32) Vec3 {
	return Vec3{p.X * k, p.Y * k, p.Z * k}
}

func (p Vec3) Shrink(k float32) Vec3 {
	return Vec3{
		p.X / k,
		p.Y / k,
		p.Z / k,
	}
}

func (p Vec3) Shrink3(q Vec3) Vec3 {
	return Vec3{
		p.X / q.X,
		p.Y / q.Y,
		p.Z / q.Z,
	}
}

func (p Vec3) Normalize() Vec3 {
	l := p.Len()
	if l == 0 {
		return Vec3{}
	}
	return Vec3{p.X / l, p.Y / l, p.Z / l}
}

func (p Vec3) Distance(q Vec3) float32 {
	return p.Sub(q).Len()
}

func (p Vec3) DistanceSquared(q Vec3) float32 {
	r := p.Sub(q)
	return r.Dot(r)
}

func (p Vec3) Lerp(t float32, q Vec3) Vec3 {
	return Vec3{
		Lerp(t, p.X, q.X),
		Lerp(t, p.Y, q.Y),
		Lerp(t, p.Z, q.Z),
	}
}

func (p Vec3) Abs() Vec3 {
	return Vec3{
		Abs(p.X),
		Abs(p.Y),
		Abs(p.Z),
	}
}

func (p Vec3) MaxScalar(k float32) Vec3 {
	return Vec3{
		Max(p.X, k),
		Max(p.Y, k),
		Max(p.Z, k),
	}
}

func (p Vec3) MinScalar(k float32) Vec3 {
	return Vec3{
		Min(p.X, k),
		Min(p.Y, k),
		Min(p.Z, k),
	}
}

func (p Vec3) Max(q Vec3) Vec3 {
	return Vec3{
		Max(p.X, q.X),
		Max(p.Y, q.Y),
		Max(p.Z, q.Z),
	}
}

func (p Vec3) Min(q Vec3) Vec3 {
	return Vec3{
		Min(p.X, q.X),
		Min(p.Y, q.Y),
		Min(p.Z, q.Z),
	}
}

func (p Vec3) Point4() Vec4 {
	return Vec4{p.X, p.Y, p.Z, 1}
}

func (p Vec3) Vec4() Vec4 {
	return Vec4{p.X, p.Y, p.Z, 0}
}

func (p Vec3) MinComp() float32 {
	return Min(p.X, Min(p.Y, p.Z))
}

func (p Vec3) MaxComp() float32 {
	return Max(p.X, Max(p.Y, p.Z))
}

func (p Vec3) Rotate(a Vec3, r float32) Vec3 {
	s, c := Sincos(r)
	v1 := p.Scale(c)
	v2 := p.Cross(a)
	v2 = v2.Scale(s)
	v3 := a.Scale(p.Dot(a))
	v3 = v3.Scale(1 - c)
	return v1.Add(v2).Add(v3)
}

func (p Vec3) ToRGBA() color.RGBA {
	if 0 <= p.X && p.X <= 1 {
		p.X *= 255
	}
	if 0 <= p.Y && p.Y <= 1 {
		p.Y *= 255
	}
	if 0 <= p.Z && p.Z <= 1 {
		p.Z *= 255
	}
	return color.RGBA{
		uint8(Clamp(p.X, 0, 255)),
		uint8(Clamp(p.Y, 0, 255)),
		uint8(Clamp(p.Z, 0, 255)),
		255,
	}

}

func (p Vec3) RGBA() (r, g, b, a uint32) {
	c := p.ToRGBA()
	return c.RGBA()
}

func (p Vec3) Spherical() Spherical {
	l := p.Len()
	return Spherical{
		R: l,
		T: Acos(p.Z / l),
		P: Atan2(p.Y, p.X),
	}
}

func (p Vec3) Equals(q Vec3, eps float32) bool {
	return Abs(p.X-q.X) <= eps && Abs(p.Y-q.Y) <= eps &&
		Abs(p.Z-q.Z) <= eps
}

func (p Vec3) XY() Vec2 { return Vec2{p.X, p.Y} }
func (p Vec3) XZ() Vec2 { return Vec2{p.X, p.Z} }
func (p Vec3) YX() Vec2 { return Vec2{p.Y, p.X} }
func (p Vec3) YZ() Vec2 { return Vec2{p.Y, p.Z} }
func (p Vec3) ZX() Vec2 { return Vec2{p.Z, p.X} }
func (p Vec3) ZY() Vec2 { return Vec2{p.Z, p.Y} }

func (p Vec3) Clamp(s, e float32) Vec3 {
	p.X = Clamp(p.X, s, e)
	p.Y = Clamp(p.Y, s, e)
	p.Z = Clamp(p.Z, s, e)
	return p
}

func (p Vec3) Clamp3(s, e Vec3) Vec3 {
	p.X = Clamp(p.X, s.X, e.X)
	p.Y = Clamp(p.Y, s.Y, e.Y)
	p.Z = Clamp(p.Z, s.Z, e.Z)
	return p
}

func (p Vec3) Angle(q Vec3) float32 {
	a := p.Len()
	b := q.Len()
	d := p.Dot(q)
	return float32(math.Acos(float64(d / (a * b))))
}

func (p Vec3) String() string {
	return fmt.Sprintf(`Vec3(%0.3f, %0.3f, %0.3f)`, p.X, p.Y, p.Z)
}

type Vec4 struct {
	X, Y, Z, W float32
}

func (p Vec4) Add(q Vec4) Vec4 {
	return Vec4{p.X + q.X, p.Y + q.Y, p.Z + q.Z, p.W + q.W}
}

func (p Vec4) Sub(q Vec4) Vec4 {
	return Vec4{p.X - q.X, p.Y - q.Y, p.Z - q.Z, p.W - q.W}
}

func (p Vec4) AddScale(q Vec4, k float32) Vec4 {
	return p.Add(q.Scale(k))
}

func (p Vec4) SubScale(q Vec4, k float32) Vec4 {
	return p.Sub(q.Scale(k))
}

func (p Vec4) Scale(k float32) Vec4 {
	return Vec4{p.X * k, p.Y * k, p.Z * k, p.W}
}

func (p Vec4) Scale4(q Vec4) Vec4 {
	return Vec4{p.X * q.X, p.Y * q.Y, p.Z * q.Z, p.W * q.W}
}

func (p Vec4) Shrink(k float32) Vec4 {
	return Vec4{p.X / k, p.Y / k, p.Z / k, p.W}
}

func (p Vec4) Shrink4(q Vec4) Vec4 {
	return Vec4{p.X / q.X, p.Y / q.Y, p.Z / q.Z, p.W}
}

func (p Vec4) Dot(q Vec4) float32 {
	return p.X*q.X + p.Y*q.Y + p.Z*q.Z + p.W*q.W
}

func (p Vec4) Len() float32 {
	return Sqrt(p.Dot(p))
}

func (p Vec4) XYZ() Vec3 { return Vec3{p.X, p.Y, p.Z} }
func (p Vec4) XZY() Vec3 { return Vec3{p.X, p.Z, p.Y} }
func (p Vec4) YXZ() Vec3 { return Vec3{p.Y, p.X, p.Z} }
func (p Vec4) YZX() Vec3 { return Vec3{p.Y, p.Z, p.X} }

func (p Vec4) Normalize() Vec4 {
	l := p.Len()
	if l == 0 {
		return Vec4{0, 0, 0, p.W}
	}
	return Vec4{
		p.X / l,
		p.Y / l,
		p.Z / l,
		p.W / l,
	}
}

func (p Vec4) ToRGBA() color.RGBA {
	if 0 <= p.X && p.X <= 1 {
		p.X *= 255
	}
	if 0 <= p.Y && p.Y <= 1 {
		p.Y *= 255
	}
	if 0 <= p.Z && p.Z <= 1 {
		p.Z *= 255
	}
	c := color.RGBA{
		uint8(Clamp(p.X, 0, 255)),
		uint8(Clamp(p.Y, 0, 255)),
		uint8(Clamp(p.Z, 0, 255)),
		uint8(Clamp(p.W, 0, 255)),
	}
	return c
}

func (p Vec4) RGBA() (r, g, b, a uint32) {
	c := p.ToRGBA()
	return c.RGBA()
}

func (p Vec4) Angle(q Vec4) float32 {
	a := p.Len()
	b := q.Len()
	d := p.Dot(q)
	return float32(math.Acos(float64(d / (a * b))))
}

func (p Vec4) PerspectiveDivide() Vec4 {
	if p.W == 0 {
		return p
	}
	return Vec4{
		p.X / p.W,
		p.Y / p.W,
		p.Z / p.W,
		1,
	}
}

func (p Vec4) Lerp(t float32, q Vec4) Vec4 {
	return Vec4{
		Lerp(t, p.X, q.X),
		Lerp(t, p.Y, q.Y),
		Lerp(t, p.Z, q.Z),
		Lerp(t, p.W, q.W),
	}
}

func (p Vec4) String() string {
	return fmt.Sprintf(`Vec4(%0.3f, %0.3f, %0.3f, %0.3f)`, p.X, p.Y, p.Z, p.W)
}

type Mat2 [2][2]float32

func (m *Mat2) Identity() *Mat2 {
	*m = Mat2{
		{1, 0},
		{0, 1},
	}
	return m
}

func (m *Mat2) Scale(s float32) *Mat2 {
	return m.Scale2(s, s)
}

func (m *Mat2) Scale2(sx, sy float32) *Mat2 {
	*m = Mat2{
		{sx, 0},
		{0, sy},
	}
	return m
}

func (m *Mat2) Rotate(r float64) *Mat2 {
	s, c := math.Sincos(r)
	*m = Mat2{
		{float32(c), float32(-s)},
		{float32(s), float32(c)},
	}
	return m
}

func (m *Mat2) Mul(a, b *Mat2) *Mat2 {
	var p Mat2
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

func (m *Mat2) Transform(p Vec2) Vec2 {
	return Vec2{
		m[0][0]*p.X + m[0][1]*p.Y,
		m[1][0]*p.X + m[1][1]*p.Y,
	}
}

func (m *Mat2) Transpose() *Mat2 {
	var p Mat2
	for i := range m {
		for j := range m[i] {
			p[j][i] = m[i][j]
		}
	}
	*m = p
	return m
}

func (m *Mat2) Row(n int) Vec2 {
	return Vec2{m[n][0], m[n][1]}
}

func (m *Mat2) Col(n int) Vec2 {
	return Vec2{m[0][n], m[1][n]}
}

func (m *Mat2) Trace() float32 {
	return m[0][0] + m[1][1]
}

func (m *Mat2) Det() float32 {
	return m[0][0]*m[1][1] - m[0][1]*m[1][0]
}

func (m *Mat2) Inverse() *Mat2 {
	det := m.Det()
	invdet := float32(0.0)
	if det != 0 {
		invdet = 1 / det
	}
	var minv Mat2
	minv[0][0] = m[1][1] * invdet
	minv[0][1] = -m[0][1] * invdet
	minv[1][0] = -m[1][0] * invdet
	minv[1][1] = m[0][0] * invdet

	*m = minv
	return m
}

func (m Mat2) String() string {
	return fmt.Sprintf(`
Mat2[% 0.3f, % 0.3f,
     % 0.3f, % 0.3f]`,
		m[0][0], m[0][1],
		m[1][0], m[1][1])
}

type Mat3 [3][3]float32

func (m *Mat3) Identity() *Mat3 {
	*m = Mat3{
		{1, 0, 0},
		{0, 1, 0},
		{0, 0, 1},
	}
	return m
}

func (m *Mat3) Mul(a, b *Mat3) *Mat3 {
	var p Mat3
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

func (m *Mat3) Transform(p Vec3) Vec3 {
	return Vec3{
		m[0][0]*p.X + m[0][1]*p.Y + m[0][2]*p.Z,
		m[1][0]*p.X + m[1][1]*p.Y + m[1][2]*p.Z,
		m[2][0]*p.X + m[2][1]*p.Y + m[2][2]*p.Z,
	}
}

func (m *Mat3) Transpose() *Mat3 {
	var p Mat3
	for i := range m {
		for j := range m[i] {
			p[j][i] = m[i][j]
		}
	}
	*m = p
	return m
}

func (m *Mat3) Det() float32 {
	A1 := Vec3{m[0][0], m[0][1], m[0][2]}
	A2 := Vec3{m[1][1], m[1][2], m[1][0]}
	A3 := Vec3{m[2][2], m[2][0], m[2][1]}
	A4 := Vec3{m[1][2], m[1][0], m[1][1]}
	A5 := Vec3{m[2][1], m[2][2], m[2][0]}

	X := A2.Scale3(A3)
	Y := A4.Scale3(A5)
	A6 := X.Sub(Y)
	return A1.Dot(A6)
}

func (m *Mat3) Trace() float32 {
	return m[0][0] + m[1][1] + m[2][2]
}

func (m *Mat3) Mat4() Mat4 {
	return Mat4{
		{m[0][0], m[1][0], m[2][0], 0},
		{m[0][1], m[1][1], m[2][1], 0},
		{m[0][2], m[1][2], m[2][2], 0},
		{0, 0, 0, 1},
	}
}

func (m *Mat3) FromBasis(X, Y, Z Vec3) *Mat3 {
	*m = Mat3{
		{X.X, Y.X, Z.X},
		{X.Y, Y.Y, Z.Y},
		{X.Z, Y.Z, Z.Z},
	}
	return m
}

func (m *Mat3) Basis() (X, Y, Z, W Vec3) {
	X = Vec3{m[0][0], m[1][0], m[2][0]}
	Y = Vec3{m[0][1], m[1][1], m[2][1]}
	Z = Vec3{m[0][2], m[1][2], m[2][2]}
	return
}

func (m *Mat3) SetCol(n int, p Vec3) {
	m[0][n] = p.X
	m[1][n] = p.Y
	m[2][n] = p.Z
}

func (m *Mat3) SetRow(n int, p Vec3) {
	m[n][0] = p.X
	m[n][1] = p.Y
	m[n][2] = p.Z
}

func (m *Mat3) Row(n int) Vec3 {
	return Vec3{m[n][0], m[n][1], m[n][2]}
}

func (m *Mat3) Col(n int) Vec3 {
	return Vec3{m[0][n], m[1][n], m[2][n]}
}

func (m *Mat3) Inverse() *Mat3 {
	det := m[0][0]*(m[1][1]*m[2][2]-m[2][1]*m[1][2]) -
		m[0][1]*(m[1][0]*m[2][2]-m[1][2]*m[2][0]) +
		m[0][2]*(m[1][0]*m[2][1]-m[1][1]*m[2][0])
	invdet := float32(0.0)
	if det != 0 {
		invdet = 1 / det
	}

	var minv Mat3
	minv[0][0] = (m[1][1]*m[2][2] - m[2][1]*m[1][2]) * invdet
	minv[0][1] = (m[0][2]*m[2][1] - m[0][1]*m[2][2]) * invdet
	minv[0][2] = (m[0][1]*m[1][2] - m[0][2]*m[1][1]) * invdet
	minv[1][0] = (m[1][2]*m[2][0] - m[1][0]*m[2][2]) * invdet
	minv[1][1] = (m[0][0]*m[2][2] - m[0][2]*m[2][0]) * invdet
	minv[1][2] = (m[1][0]*m[0][2] - m[0][0]*m[1][2]) * invdet
	minv[2][0] = (m[1][0]*m[2][1] - m[2][0]*m[1][1]) * invdet
	minv[2][1] = (m[2][0]*m[0][1] - m[0][0]*m[2][1]) * invdet
	minv[2][2] = (m[0][0]*m[1][1] - m[1][0]*m[0][1]) * invdet
	*m = minv
	return m
}

func (m Mat3) String() string {
	return fmt.Sprintf(`
Mat3[% 0.3f, % 0.3f, % 0.3f,
     % 0.3f, % 0.3f, % 0.3f,
     % 0.3f, % 0.3f, % 0.3f]`,
		m[0][0], m[0][1], m[0][2],
		m[1][0], m[1][1], m[1][2],
		m[2][0], m[2][1], m[2][2])
}

type Mat4 [4][4]float32

func (m *Mat4) Identity() *Mat4 {
	*m = Mat4{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
	return m
}

func (m *Mat4) Add(a, b *Mat4) *Mat4 {
	for i := range a {
		for j := range a[i] {
			m[i][j] = a[i][j] + b[i][j]
		}
	}
	return m
}

func (m *Mat4) Sub(a, b *Mat4) *Mat4 {
	for i := range a {
		for j := range a[i] {
			m[i][j] = a[i][j] - b[i][j]
		}
	}
	return m
}

func (m *Mat4) Mul(a, b *Mat4) *Mat4 {
	var p Mat4
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

func (m *Mat4) Trace() float32 {
	return m[0][0] + m[1][1] + m[2][2] + m[3][3]
}

func (m *Mat4) Translate3(p Vec3) *Mat4 {
	return m.Translate(p.X, p.Y, p.Z)
}

func (m *Mat4) Translate4(p Vec4) *Mat4 {
	return m.Translate(p.X, p.Y, p.Z)
}

func (m *Mat4) Translate(tx, ty, tz float32) *Mat4 {
	*m = Mat4{
		{1, 0, 0, tx},
		{0, 1, 0, ty},
		{0, 0, 1, tz},
		{0, 0, 0, 1},
	}
	return m
}

func (m *Mat4) Scale3(p Vec3) *Mat4 {
	return m.Scale(p.X, p.Y, p.Z)
}

func (m *Mat4) Scale4(p Vec4) *Mat4 {
	return m.Scale(p.X, p.Y, p.Z)
}

func (m *Mat4) Scale(sx, sy, sz float32) *Mat4 {
	*m = Mat4{
		{sx, 0, 0, 0},
		{0, sy, 0, 0},
		{0, 0, sz, 0},
		{0, 0, 0, 1},
	}
	return m
}

func (m *Mat4) LookAt(eye, center, up Vec3) *Mat4 {
	f := center
	f = f.Sub(eye)
	f = f.Normalize()

	s := f.Cross(up)
	s = s.Normalize()
	u := s.Cross(f)

	*m = Mat4{
		{s.X, s.Y, s.Z, 0},
		{u.X, u.Y, u.Z, 0},
		{-f.X, -f.Y, -f.Z, 0},
		{0, 0, 0, 1},
	}

	var t Mat4
	t.Translate(-eye.X, -eye.Y, -eye.Z)
	m.Mul(m, &t)
	return m
}

func (m *Mat4) Frustum(l, r, b, t, n, f float32) *Mat4 {
	A := (r + l) / (r - l)
	B := (t + b) / (t - b)
	C := -(f + n) / (f - n)
	D := -2 * f * n / (f - n)
	*m = Mat4{
		{2 * n / (r - l), 0, A, 0},
		{0, 2 * n / (t - b), B, 0},
		{0, 0, C, D},
		{0, 0, -1, 0},
	}
	return m
}

func (m *Mat4) Perspective(fovy, aspect, near, far float32) *Mat4 {
	ymax := near * Tan(fovy/2)
	xmax := ymax * aspect
	m.Frustum(-xmax, xmax, -ymax, ymax, near, far)
	return m
}

func (m *Mat4) InfinitePerspective(fovy, aspect, near float32) *Mat4 {
	const zp = 0.0
	f := 1 / Tan(fovy/2)
	*m = Mat4{
		{f / aspect, 0, 0, 0},
		{0, f, 0, 0},
		{0, 0, -(1 - zp), -near * (1 - zp)},
		{0, 0, -1, 0},
	}
	return m
}

func (m *Mat4) Ortho(l, r, b, t, n, f float32) *Mat4 {
	sx := 2 / (r - l)
	sy := 2 / (t - b)
	sz := -2 / (f - n)

	tx := -(r + l) / (r - l)
	ty := -(t + b) / (t - b)
	tz := -(f + n) / (f - n)

	*m = Mat4{
		{sx, 0, 0, tx},
		{0, sy, 0, ty},
		{0, 0, sz, tz},
		{0, 0, 0, 1},
	}
	return m
}

func (m *Mat4) Viewport(x, y, w, h float32) *Mat4 {
	l := x
	b := y
	r := x + w
	t := y + h
	*m = Mat4{
		{(r - l) / 2, 0, 0, (r + l) / 2},
		{0, (t - b) / 2, 0, (t + b) / 2},
		{0, 0, 0.5, 0.5},
		{0, 0, 0, 1},
	}
	return m
}

func (m *Mat4) Inverse() *Mat4 {
	a := Vec3{m[0][0], m[1][0], m[2][0]}
	b := Vec3{m[0][1], m[1][1], m[2][1]}
	c := Vec3{m[0][2], m[1][2], m[2][2]}
	d := Vec3{m[0][3], m[1][3], m[2][3]}

	s := a.Cross(b)
	t := c.Cross(d)

	invDet := 1 / s.Dot(c)

	s = s.Scale(invDet)
	t = t.Scale(invDet)
	v := c.Scale(invDet)

	r0 := b.Cross(v)
	r1 := v.Cross(a)

	*m = Mat4{
		{r0.X, r0.Y, r0.Z, -b.Dot(t)},
		{r1.X, r1.Y, r1.Z, a.Dot(t)},
		{s.X, s.Y, s.Z, -d.Dot(s)},
		{0, 0, 0, 1},
	}
	return m
}

func (m *Mat4) Transpose() *Mat4 {
	var p Mat4
	for i := range m {
		for j := range m[i] {
			p[j][i] = m[i][j]
		}
	}
	*m = p
	return m
}

func (m *Mat4) Transform(v Vec4) Vec4 {
	return Vec4{
		m[0][0]*v.X + m[0][1]*v.Y + m[0][2]*v.Z + m[0][3]*v.W,
		m[1][0]*v.X + m[1][1]*v.Y + m[1][2]*v.Z + m[1][3]*v.W,
		m[2][0]*v.X + m[2][1]*v.Y + m[2][2]*v.Z + m[2][3]*v.W,
		m[3][0]*v.X + m[3][1]*v.Y + m[3][2]*v.Z + m[3][3]*v.W,
	}
}

func (m *Mat4) Transform3(v Vec3) Vec3 {
	s := m[3][0]*v.X + m[3][1]*v.Y + m[3][2]*v.Z + m[3][3]
	switch s {
	case 0:
		return Vec3{}
	default:
		invs := 1 / s
		p := m.Transform(Vec4{v.X, v.Y, v.Z, 1})
		return Vec3{p.X * invs, p.Y * invs, p.Z * invs}
	}

}

func (m *Mat4) FromBasis3(X, Y, Z, W Vec3) *Mat4 {
	*m = Mat4{
		{X.X, Y.X, Z.X, 0},
		{X.Y, Y.Y, Z.Y, 0},
		{X.Z, Y.Z, Z.Z, 0},
		{0, 0, 0, 1},
	}
	return m
}

func (m *Mat4) FromBasis(X, Y, Z, W Vec4) *Mat4 {
	*m = Mat4{
		{X.X, Y.X, Z.X, W.X},
		{X.Y, Y.Y, Z.Y, W.Y},
		{X.Z, Y.Z, Z.Z, W.Z},
		{X.W, Y.W, Z.W, W.W},
	}
	return m
}

func (m *Mat4) Basis3() (X, Y, Z, W Vec3) {
	X = Vec3{m[0][0], m[1][0], m[2][0]}
	Y = Vec3{m[0][1], m[1][1], m[2][1]}
	Z = Vec3{m[0][2], m[1][2], m[2][2]}
	W = Vec3{m[0][3], m[1][3], m[2][3]}
	return
}

func (m *Mat4) Basis() (X, Y, Z, W Vec4) {
	X = Vec4{m[0][0], m[1][0], m[2][0], m[3][0]}
	Y = Vec4{m[0][1], m[1][1], m[2][1], m[3][1]}
	Z = Vec4{m[0][2], m[1][2], m[2][2], m[3][2]}
	W = Vec4{m[0][3], m[1][3], m[2][3], m[3][3]}
	return
}

func (m *Mat4) SetCol(n int, p Vec4) {
	m[0][n] = p.X
	m[1][n] = p.Y
	m[2][n] = p.Z
	m[3][n] = p.W
}

func (m *Mat4) SetRow(n int, p Vec4) {
	m[n][0] = p.X
	m[n][1] = p.Y
	m[n][2] = p.Z
	m[n][3] = p.W
}

func (m *Mat4) Row(n int) Vec4 {
	return Vec4{m[n][0], m[n][1], m[n][2], m[n][3]}
}

func (m *Mat4) Col(n int) Vec4 {
	return Vec4{m[0][n], m[1][n], m[2][n], m[3][n]}
}

func (m *Mat4) Det() float32 {
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

func (m *Mat4) Adjoint() *Mat4 {
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

	var a Mat4
	a[0][0] = m11*(m22*m33-m23*m32) + m21*(m13*m32-m12*m33) + m31*(m12*m23-m13*m22)
	a[0][1] = m01*(m23*m32-m22*m33) + m21*(m02*m33-m03*m32) + m31*(m03*m22-m02*m23)
	a[0][2] = m01*(m12*m33-m13*m32) + m11*(m03*m32-m02*m33) + m31*(m02*m13-m03*m12)
	a[0][3] = m01*(m13*m22-m12*m23) + m11*(m02*m23-m03*m22) + m21*(m03*m12-m02*m13)
	a[1][0] = m10*(m23*m32-m22*m33) + m20*(m12*m33-m13*m32) + m30*(m13*m22-m12*m23)
	a[1][1] = m00*(m22*m33-m23*m32) + m20*(m03*m32-m02*m33) + m30*(m02*m23-m03*m22)
	a[1][2] = m00*(m13*m32-m12*m33) + m10*(m02*m33-m03*m32) + m30*(m03*m12-m02*m13)
	a[1][3] = m00*(m12*m23-m13*m22) + m10*(m03*m22-m02*m23) + m20*(m02*m13-m03*m12)
	a[2][0] = m10*(m21*m33-m23*m31) + m20*(m13*m31-m11*m33) + m30*(m11*m23-m13*m21)
	a[2][1] = m00*(m23*m31-m21*m33) + m20*(m01*m33-m03*m31) + m30*(m03*m21-m01*m23)
	a[2][2] = m00*(m11*m33-m13*m31) + m10*(m03*m31-m01*m33) + m30*(m01*m13-m03*m11)
	a[2][3] = m00*(m13*m21-m11*m23) + m10*(m01*m23-m03*m21) + m20*(m03*m11-m01*m13)
	a[3][0] = m10*(m22*m31-m21*m32) + m20*(m11*m32-m12*m31) + m30*(m12*m21-m11*m22)
	a[3][1] = m00*(m21*m32-m22*m31) + m20*(m02*m31-m01*m32) + m30*(m01*m22-m02*m21)
	a[3][2] = m00*(m12*m31-m11*m32) + m10*(m01*m32-m02*m31) + m30*(m02*m11-m01*m12)
	a[3][3] = m00*(m11*m22-m12*m21) + m10*(m02*m21-m01*m22) + m20*(m01*m12-m02*m11)

	*m = a
	return m
}

func (m Mat4) String() string {
	return fmt.Sprintf(`
Mat4[% 0.3f, % 0.3f, % 0.3f, % 0.3f,
     % 0.3f, % 0.3f, % 0.3f, % 0.3f,
     % 0.3f, % 0.3f, % 0.3f, % 0.3f,
     % 0.3f, % 0.3f, % 0.3f, % 0.3f]`,
		m[0][0], m[0][1], m[0][2], m[0][3],
		m[1][0], m[1][1], m[1][2], m[1][3],
		m[2][0], m[2][1], m[2][2], m[2][3],
		m[3][0], m[3][1], m[3][2], m[3][3])
}

type Polar struct {
	R, T float32
}

func (p Polar) Mul(q Polar) Polar {
	return Polar{p.R * q.R, p.T + q.T}
}

func (p Polar) Quo(q Polar) Polar {
	return Polar{p.R / q.R, p.T - q.T}
}

func (p Polar) Cartesian() Vec2 {
	s, c := Sincos(p.T)
	return Vec2{p.R * c, p.R * s}
}

type Quat struct {
	X, Y, Z, W float32
}

func (q Quat) Add(r Quat) Quat {
	return Quat{
		q.X + r.X,
		q.Y + r.Y,
		q.Z + r.Z,
		q.W + r.W,
	}
}

func (q Quat) Sub(r Quat) Quat {
	return Quat{
		q.X - r.X,
		q.Y - r.Y,
		q.Z - r.Z,
		q.W - r.W,
	}
}

func (q Quat) Mul(r Quat) Quat {
	w := q.W*r.W - q.X*r.X - q.Y*r.Y - q.Z*r.Z
	x := q.X*r.W + q.W*r.X + q.Y*r.Z - q.Z*r.Y
	y := q.Y*r.W + q.W*r.Y + q.Z*r.X - q.X*r.Z
	z := q.Z*r.W + q.W*r.Z + q.X*r.Y - q.Y*r.X
	return Quat{x, y, z, w}
}

func (q Quat) Mulv(r Vec4) Quat {
	w := -q.X*r.X - q.Y*r.Y - q.Z*r.Z
	x := q.W*r.X + q.Y*r.Z - q.Z*r.Y
	y := q.W*r.Y + q.Z*r.X - q.X*r.Z
	z := q.W*r.Z + q.X*r.Y - q.Y*r.X
	return Quat{x, y, z, w}
}

func (q Quat) Neg() Quat {
	return Quat{-q.X, -q.Y, -q.Z, -q.W}
}

func (q Quat) Dot(p Quat) float32 {
	return q.W*q.W + q.X*q.X + q.Y*q.Y + q.Z*q.Z
}

func (q Quat) Scale(k float32) Quat {
	return Quat{
		q.X * k,
		q.Y * k,
		q.Z * k,
		q.W * k,
	}
}

func (q Quat) Len() float32 {
	return Sqrt(q.X*q.X + q.Y*q.Y + q.Z*q.Z + q.W*q.W)
}

func (q Quat) Normalize() Quat {
	l := q.Len()
	if l == 0 {
		return Quat{}
	}
	return Quat{
		q.X / l,
		q.Y / l,
		q.Z / l,
		q.W / l,
	}
}

func (q Quat) Conj() Quat {
	return Quat{-q.X, -q.Y, -q.Z, q.W}
}

func (q Quat) FromAxisAngle(v Vec3, r float32) Quat {
	r *= 0.5
	vn := v.Normalize()
	si, co := Sincos(r)
	return Quat{
		vn.X * si,
		vn.Y * si,
		vn.Z * si,
		co,
	}
}

func (q Quat) FromEuler(pitch, yaw, roll float32) Quat {
	p := pitch / 2
	y := yaw / 2
	r := roll / 2

	sinp := Sin(p)
	siny := Sin(y)
	sinr := Sin(r)
	cosp := Cos(p)
	cosy := Cos(y)
	cosr := Cos(r)

	return Quat{
		sinr*cosp*cosy - cosr*sinp*siny,
		cosr*sinp*cosy + sinr*cosp*siny,
		cosr*cosp*siny - sinr*sinp*cosy,
		cosr*cosp*cosy + sinr*sinp*siny,
	}.Normalize()
}

func (q Quat) Mat3() Mat3 {
	x, y, z, w := q.X, q.Y, q.Z, q.W
	x2 := x * x
	y2 := y * y
	z2 := z * z
	xy := x * y
	xz := x * z
	yz := y * z
	wx := w * x
	wy := w * y
	wz := w * z

	m := Mat3{
		{1.0 - 2.0*(y2+z2), 2.0 * (xy - wz), 2.0 * (xz + wy)},
		{2.0 * (xy + wz), 1.0 - 2.0*(x2+z2), 2.0 * (yz - wx)},
		{2.0 * (xz - wy), 2.0 * (yz + wx), 1.0 - 2.0*(x2+y2)},
	}
	return m
}

func (q Quat) Mat4() Mat4 {
	m := q.Mat3()
	return m.Mat4()
}

func (q Quat) Transform3(v Vec3) Vec3 {
	m := q.Mat3()
	return m.Transform(v)
}

func (q Quat) Transform4(v Vec4) Vec4 {
	m := q.Mat4()
	return m.Transform(v)
}

func (q Quat) AxisAngle() (v Vec3, r float32) {
	s := Sqrt(q.X*q.X + q.Y*q.Y + q.Z*q.Z)
	v = Vec3{q.X / s, q.Y / s, q.Z / s}
	r = Acos(q.W) * 2
	return
}

func (q Quat) Lerp(t float32, p Quat) Quat {
	return q.Add(p.Sub(q).Scale(t))
}

func (q Quat) Inverse() Quat {
	l2 := q.X*q.X + q.Y*q.Y + q.Z*q.Z + q.W*q.W
	return Quat{
		-q.X / l2,
		-q.Y / l2,
		-q.Z / l2,
		q.W / l2,
	}
}

func (q Quat) Powu(p float32) Quat {
	t := Acos(q.W) * p
	u := Vec3{q.X, q.Y, q.Z}
	u = u.Normalize()
	u = u.Scale(Sin(t))
	w := Cos(t)
	return Quat{
		u.X, u.Y, u.Z, w,
	}
}

func (q Quat) Slerp(t float32, p Quat) Quat {
	v0 := q.Normalize()
	v1 := p.Normalize()

	const threshold = 0.9995
	dot := v0.Dot(v1)
	if dot > threshold {
		return v0.Lerp(t, v1).Normalize()
	}

	if dot < 0 {
		v1 = v1.Neg()
		dot = -dot
	}

	dot = Clamp(dot, -1, 1)
	theta0 := Acos(dot)
	theta := theta0 * t

	v2 := v1.Sub(v0.Scale(dot))
	v2 = v2.Normalize()

	v3 := v0.Scale(Cos(theta))
	v4 := v2.Scale(Sin(theta))
	return v3.Add(v4)
}

func (q Quat) String() string {
	return fmt.Sprintf(`Quat(%0.3f, %0.3f, %0.3f, %0.3f)`, q.X, q.Y, q.Z, q.W)
}

type Spherical struct {
	R, T, P float32
}

func (s Spherical) Euclidean() Vec3 {
	sint := Sin(s.T)
	cost := Sin(s.T)
	sinp := Sin(s.P)
	cosp := Sin(s.P)
	r := s.R

	return Vec3{
		r * sint * cosp,
		r * sint * sinp,
		r * cost,
	}
}

func Lerp(t, a, b float32) float32 {
	return a + t*(b-a)
}

func Unlerp(t, a, b float32) float32 {
	return (t - a) / (b - a)
}

func LinearRemap(x, a, b, c, d float32) float32 {
	return Lerp(Unlerp(x, a, b), c, d)
}

func Smoothstep(a, b, x float32) float32 {
	t := Clamp((x-a)/(b-a), 0, 1)
	return t * t * (3 - 2*t)
}

func CubicBezier1D(t, p0, p1, p2, p3 float32) float32 {
	it := 1 - t
	return it*it*it*p0 + 3*it*it*t*p1 + 3*it*t*t*p2 + t*t*t*p3
}

func Clamp(x, s, e float32) float32 {
	if x < s {
		x = s
	}
	if x > e {
		x = e
	}
	return x
}

func Saturate(x float32) float32 {
	return Max(0, Min(1, x))
}

func SignNZ(x float32) float32 {
	if x >= 0 {
		return 1
	}
	return -1
}

func Sign(x float32) float32 {
	if x < 0 {
		return -1
	}
	if x == 0 {
		return 0
	}
	return 1
}

func Deg2Rad(d float32) float32 {
	return d * math.Pi / 180
}

func Rad2Deg(r float32) float32 {
	return r * 180 / math.Pi
}

type Circle struct {
	X, Y, R float32
}

func (c Circle) InPoint(x, y float32) bool {
	return (x-c.X)*(x-c.X)+(y-c.Y)*(y-c.Y) <= c.R
}

func (c Circle) InRect(r Rectangle) bool {
	dx := c.X - Max(r.Min.X, Min(c.X, r.Max.X))
	dy := c.Y - Max(r.Min.Y, Min(c.Y, r.Max.Y))
	return dx*dx+dy*dy <= c.R
}

type Rectangle struct {
	Min, Max Vec2
}

func Rect(x0, y0, x1, y1 float32) Rectangle {
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	return Rectangle{Vec2{x0, y0}, Vec2{x1, y1}}
}

func (r Rectangle) Empty() bool {
	return r.Min.X >= r.Max.X || r.Min.Y >= r.Max.Y
}

func (r Rectangle) Overlaps(s Rectangle) bool {
	return !r.Empty() && !s.Empty() &&
		r.Min.X < s.Max.X && s.Min.X < r.Max.X &&
		r.Min.Y < s.Max.Y && s.Min.Y < r.Max.Y
}

func (r Rectangle) Intersect(s Rectangle) Rectangle {
	if r.Min.X < s.Min.X {
		r.Min.X = s.Min.X
	}
	if r.Min.Y < s.Min.Y {
		r.Min.Y = s.Min.Y
	}
	if r.Max.X > s.Max.X {
		r.Max.X = s.Max.X
	}
	if r.Max.Y > s.Max.Y {
		r.Max.Y = s.Max.Y
	}

	if r.Empty() {
		return Rectangle{}
	}
	return r
}

func (r Rectangle) Union(s Rectangle) Rectangle {
	if r.Empty() {
		return s
	}
	if s.Empty() {
		return r
	}
	if r.Min.X > s.Min.X {
		r.Min.X = s.Min.X
	}
	if r.Min.Y > s.Min.Y {
		r.Min.Y = s.Min.Y
	}
	if r.Max.X < s.Max.X {
		r.Max.X = s.Max.X
	}
	if r.Max.Y < s.Max.Y {
		r.Max.Y = s.Max.Y
	}
	return r
}

func (r Rectangle) Int() image.Rectangle {
	return image.Rect(int(r.Min.X), int(r.Min.Y), int(r.Max.X), int(r.Max.Y))
}

func (r Rectangle) Canon() Rectangle {
	if r.Max.X < r.Min.X {
		r.Min.X, r.Max.X = r.Max.X, r.Min.X
	}
	if r.Max.Y < r.Min.Y {
		r.Min.Y, r.Max.Y = r.Max.Y, r.Min.Y
	}
	return r
}

func (r Rectangle) Scale(s Vec2) Rectangle {
	r.Min.X *= s.X
	r.Max.X *= s.X
	r.Min.Y *= s.Y
	r.Max.Y *= s.Y
	return r
}

func (r Rectangle) Add(p Vec2) Rectangle {
	return Rectangle{
		Vec2{r.Min.X + p.X, r.Min.Y + p.Y},
		Vec2{r.Max.X + p.X, r.Max.Y + p.Y},
	}
}

func (r Rectangle) Sub(p Vec2) Rectangle {
	return Rectangle{
		Vec2{r.Min.X - p.X, r.Min.Y - p.Y},
		Vec2{r.Max.X - p.X, r.Max.Y - p.Y},
	}
}

func (r Rectangle) Size() Vec2 {
	return Vec2{
		r.Max.X - r.Min.X,
		r.Max.Y - r.Min.Y,
	}
}

func (r Rectangle) Inset(n float32) Rectangle {
	if r.Dx() < 2*n {
		r.Min.X = (r.Min.X + r.Max.X) / 2
		r.Max.X = r.Min.X
	} else {
		r.Min.X += n
		r.Max.X -= n
	}
	if r.Dy() < 2*n {
		r.Min.Y = (r.Min.Y + r.Max.Y) / 2
		r.Max.Y = r.Min.Y
	} else {
		r.Min.Y += n
		r.Max.Y -= n
	}
	return r
}

func (r Rectangle) Dx() float32 {
	return r.Max.X - r.Min.X
}

func (r Rectangle) Dy() float32 {
	return r.Max.Y - r.Min.Y
}

func (r Rectangle) Center() Vec2 {
	return Vec2{
		(r.Min.X + r.Max.X) / 2,
		(r.Min.Y + r.Max.Y) / 2,
	}
}

func (r Rectangle) Diagonal() float32 {
	x := r.Max.X - r.Min.X
	y := r.Max.Y - r.Min.Y
	return Sqrt(x*x + y*y)
}

func (r Rectangle) PosSize() (x, y, w, h float32) {
	return r.Min.X, r.Min.Y, r.Dx(), r.Dy()
}

func (r Rectangle) In(s Rectangle) bool {
	if r.Empty() {
		return true
	}

	return s.Min.X <= r.Min.X && r.Max.X <= s.Max.X &&
		s.Min.Y <= r.Min.Y && r.Max.Y <= s.Max.Y
}

func (r Rectangle) TL() Vec2 { return r.Min }
func (r Rectangle) TR() Vec2 { return Vec2{r.Max.X, r.Min.Y} }
func (r Rectangle) BL() Vec2 { return Vec2{r.Min.X, r.Max.Y} }
func (r Rectangle) BR() Vec2 { return r.Max }

func (r Rectangle) Inverted() bool {
	return r.Min.X > r.Max.X || r.Min.Y > r.Max.Y
}

func (r Rectangle) Expand(x, y float32) Rectangle {
	r.Min.X -= x
	r.Min.Y -= y
	r.Max.X += x
	r.Max.Y += y
	return r
}

func (r Rectangle) Expand2(v Vec2) Rectangle {
	return r.Expand(v.X, v.Y)
}

func RoundPrec(v float32, prec int) float32 {
	if prec < 0 {
		return v
	}

	tab := [...]float32{
		1, 1e-1, 1e-2, 1e-3, 1e-4, 1e-5, 1e-6, 1e-7, 1e-8, 1e-9, 1e-10,
	}
	step := float32(0.0)
	if prec < len(tab) {
		step = tab[prec]
	} else {
		step = Pow(10, float32(-prec))
	}

	neg := v < 0
	v = Abs(v)
	rem := Mod(v, step)
	if rem <= step*0.5 {
		v -= rem
	} else {
		v += step - rem
	}

	if neg {
		v = -v
	}

	return v
}

func Sinc(x float32) float32 {
	x *= math.Pi
	if x < 0.01 && x > -0.01 {
		return 1 + x*x*((-1.0/6)+x*x*1.0/120)
	}
	return Sin(x) / x
}

func LinearController(curpos *float32, targetpos, acc, deacc, dt float32) {
	sign := float32(1.0)
	p := float32(0.0)
	cp := *curpos
	if cp == targetpos {
		return
	}
	if targetpos < cp {
		targetpos = -targetpos
		cp = -cp
		sign = -1
	}

	// first decelerate
	if cp < 0 {
		p = cp + deacc*dt
		if p > 0 {
			p = 0
			dt = dt - p/deacc
			if dt < 0 {
				dt = 0
			}
		} else {
			dt = 0
		}
		cp = p
	}

	// now accelerate
	p = cp + acc*dt
	if p > targetpos {
		p = targetpos
	}
	*curpos = p * sign
}

func Multiple(a, m float32) float32 {
	return Ceil(a/m) * m
}

func Abs(x float32) float32 {
	return float32(math.Abs(float64(x)))
}

func Min(a, b float32) float32 {
	return float32(math.Min(float64(a), float64(b)))
}

func Max(a, b float32) float32 {
	return float32(math.Max(float64(a), float64(b)))
}

func Sin(x float32) float32 {
	return float32(math.Sin(float64(x)))
}

func Cos(x float32) float32 {
	return float32(math.Cos(float64(x)))
}

func Tan(x float32) float32 {
	return float32(math.Tan(float64(x)))
}

func Floor(x float32) float32 {
	return float32(math.Floor(float64(x)))
}

func Ceil(x float32) float32 {
	return float32(math.Ceil(float64(x)))
}

func Sqrt(x float32) float32 {
	return float32(math.Sqrt(float64(x)))
}

func Atan2(y, x float32) float32 {
	return float32(math.Atan2(float64(y), float64(x)))
}

func Sincos(x float32) (si, co float32) {
	s, c := math.Sincos(float64(x))
	return float32(s), float32(c)
}

func Acos(x float32) float32 {
	return float32(math.Acos(float64(x)))
}

func Pow(x, y float32) float32 {
	return float32(math.Pow(float64(x), float64(y)))
}

func Pow10(x int) float32 {
	return float32(math.Pow10(x))
}

func Round(x float32) float32 {
	return float32(math.Round(float64(x)))
}

func Mod(x, y float32) float32 {
	return float32(math.Mod(float64(x), float64(y)))
}

func Hypot(x, y float32) float32 {
	return Sqrt(x*x + y*y)
}

func Trunc(x float32) float32 {
	return float32(math.Trunc(float64(x)))
}

func Sinh(x float32) float32 {
	return float32(math.Sinh(float64(x)))
}

func Cosh(x float32) float32 {
	return float32(math.Cosh(float64(x)))
}

func Tanh(x float32) float32 {
	return float32(math.Tanh(float64(x)))
}

func Exp(x float32) float32 {
	return float32(math.Exp(float64(x)))
}

func Log(x float32) float32 {
	return float32(math.Log(float64(x)))
}

func Log2(x float32) float32 {
	return float32(math.Log2(float64(x)))
}

func Log10(x float32) float32 {
	return float32(math.Log10(float64(x)))
}

func Log1p(x float32) float32 {
	return float32(math.Log1p(float64(x)))
}

func Log1b(x float32) float32 {
	return float32(math.Logb(float64(x)))
}

func Wrap(x, s, e float32) float32 {
	if x < s {
		x += e
	}
	if x >= e {
		x -= e
	}
	return x
}

func TriangleBarycentric(p, a, b, c Vec2) Vec3 {
	x := Vec3{
		c.X - a.X,
		b.X - a.X,
		a.X - p.X,
	}
	y := Vec3{
		c.Y - a.Y,
		b.Y - a.Y,
		a.Y - p.Y,
	}
	u := x.Cross(y)
	if Abs(u.Z) > 1e-2 {
		return Vec3{
			1 - (u.X+u.Y)/u.Z,
			u.Y / u.Z,
			u.X / u.Z,
		}
	}
	return Vec3{-1, -1, -1}
}

func Clamp8(x, a, b float32) uint8 {
	x = Round(x)
	if x < a {
		x = a
	}
	if x > b {
		x = b
	}
	return uint8(x)
}

func ditfft2c(x, y []complex64, n, s int) {
	if n == 1 {
		y[0] = x[0]
		return
	}
	ditfft2c(x, y, n/2, 2*s)
	ditfft2c(x[s:], y[n/2:], n/2, 2*s)
	for k := 0; k < n/2; k++ {
		tf := complex64(cmplx.Rect(1, -2*math.Pi*float64(k)/float64(n))) * y[k+n/2]
		y[k], y[k+n/2] = y[k]+complex64(tf), y[k]-complex64(tf)
	}
}

func ditfft2r(x []float32, y []complex64, n, s int) {
	if n == 1 {
		y[0] = complex(x[0], 0)
		return
	}
	ditfft2r(x, y, n/2, 2*s)
	ditfft2r(x[s:], y[n/2:], n/2, 2*s)
	for k := 0; k < n/2; k++ {
		tf := complex64(cmplx.Rect(1, float64(-2*math.Pi*float32(k)/float32(n)))) * y[k+n/2]
		y[k], y[k+n/2] = y[k]+tf, y[k]-tf
	}
}

func FFT1DC(dst, src []complex64) {
	ditfft2c(src, dst, len(dst), 1)
}

func IFFT1DC(dst, src []complex64) {
	for i := range src {
		src[i] = complex64(cmplx.Conj(complex128(src[i])))
	}
	FFT1DC(dst, src)
	for i := range src {
		src[i] = complex64(cmplx.Conj(complex128(src[i])))
	}
	for i := range dst {
		dst[i] = complex64(cmplx.Conj(complex128(dst[i]))) / complex64(complex(float64(len(dst)), 0))
	}
}

func FFT1DR(dst []complex64, src []float32) {
	ditfft2r(src, dst, len(dst), 1)
}

func IFFT1DR(dst []float32, src []complex64) {
	for i := range src {
		src[i] = complex64(cmplx.Conj(complex128(src[i])))
	}
	tmp := make([]complex64, len(dst))
	FFT1DC(tmp, src)
	for i := range src {
		src[i] = complex64(cmplx.Conj(complex128(src[i])))
	}
	for i := range dst {
		dst[i] = real(tmp[i]) / float32(len(dst))
	}
}

func Simpson1D(f func(x float32) float32, start, end float32, n int) float32 {
	r := float32(0.0)
	s := (end - start) / float32(n)
	i := 0

	r += f(start)
	for j := 1; j < n; j++ {
		r += (4 - float32(i<<1)) * f(start+float32(j)*s)
		i = (i + 1) & 1
	}
	r += f(end)
	r *= s / 3
	return r
}

func simpsonweight(i, n int) float32 {
	if i == 0 || i == n {
		return 1
	}
	if i%2 != 0 {
		return 4
	}
	return 2
}

func Simpson2D(f func(x, y float32) float32, x0, x1, y0, y1 float32, m, n int) float32 {
	if n%2 != 0 || m%2 != 0 {
		panic("integration range must be even")
	}

	dx := (x1 - x0) / float32(m)
	dy := (y1 - y0) / float32(n)
	r := float32(0.0)
	for i := 0; i <= n; i++ {
		y := y0 + float32(i)*dy
		wy := simpsonweight(i, n)
		for j := 0; j <= m; j++ {
			x := x0 + float32(j)*dx
			wx := simpsonweight(j, m)
			r += f(x, y) * wx * wy
		}
	}
	r *= dx * dy / (9 * float32(m) * float32(n))
	return r
}

func Convolve1D(dst, src, coeffs []float32, shape int) []float32 {
	var m, n int
	switch shape {
	case 'f':
		m = len(src) + len(coeffs) - 1
	case 's':
		m = len(src)
		n = len(coeffs) - 2
	case 'v':
		m = len(src) - len(coeffs) + 1
		n = len(coeffs) - 1
		if m < 0 {
			m = 0
		}
	default:
		panic("unsupported convolution shape")
	}

	for k := 0; k < m; k++ {
		dst[k] = 0
		for j := range src {
			l := k + n - j
			if l < 0 || l >= len(coeffs) {
				continue
			}
			dst[k] += src[j] * coeffs[l]
		}
	}

	return dst[:m]
}

func Sample1D(f func(float32) float32, x0, x1 float32, n int) (p []float32, s float32) {
	p = make([]float32, n)
	s = (x1 - x0) / float32(n-1)
	for i := 0; i < n; i++ {
		p[i] = f(x0 + float32(i)*s)
	}
	return
}

func Sample2D(f func(x, y float32) float32, x0, x1, y0, y1 float32, nx, ny int) (p []float32, sx, sy float32) {
	p = make([]float32, nx*ny)
	sx = (x1 - x0) / float32(nx-1)
	sy = (y1 - y0) / float32(ny-1)
	for y := 0; y < ny; y++ {
		for x := 0; x < nx; x++ {
			p[y*nx+x] = f(x0+sx*float32(x), y0+sy*float32(y))
		}
	}
	return
}

func FloatToComplex(v []float32) []complex64 {
	p := make([]complex64, len(v))
	for i := range v {
		p[i] = complex(v[i], 0)
	}
	return p
}

func ComplexToFloat(v []complex64) []float32 {
	p := make([]float32, len(v))
	for i := range v {
		p[i] = float32(cmplx.Abs(complex128(v[i])))
	}
	return p
}

func PearsonCorrelation1D(x, y []float32) float32 {
	return Cov1D(x, y, 1) / (Stddev(x, 1) * Stddev(y, 1))
}

func Cov1D(x, y []float32, ddof float32) float32 {
	mx := Mean(x)
	my := Mean(y)
	s := float32(0.0)
	for i := range x {
		s += (x[i] - mx) * (y[i] * my)
	}
	return s / (float32(len(x)) - ddof)
}

func Mean(x []float32) float32 {
	s := float32(0.0)
	for i := range x {
		s += x[i]
	}
	return s / float32(len(x))
}

func Stddev(x []float32, ddof float32) float32 {
	if len(x) <= 1 {
		return 0
	}
	xm := Mean(x)
	s := float32(0.0)
	for i := range x {
		s += (x[i] - xm) * (x[i] - xm)
	}
	return Sqrt(s / (float32(len(x)) - ddof))
}

func Median(x []float32) float32 {
	if len(x) == 0 {
		return 0
	}
	sort.Slice(x, func(i, j int) bool {
		return x[i] < x[j]
	})
	return x[len(x)/2]
}

func Mins(x ...float32) float32 {
	if len(x) == 0 {
		return 0
	}

	n := x[0]
	for i := range x[1:] {
		n = Min(n, x[i])
	}
	return n
}

func Maxs(x ...float32) float32 {
	if len(x) == 0 {
		return 0
	}

	n := x[0]
	for i := range x[1:] {
		n = Max(n, x[i])
	}
	return n
}

func Iround(x float32) int {
	return int(math.Round(float64(x)))
}

func Ifloor(x float32) int {
	return int(math.Floor(float64(x)))
}

func Iceil(x float32) int {
	return int(math.Ceil(float64(x)))
}

func Randn(min, max float32) float32 {
	return min + (max-min)*float32(rand.Float64())
}
