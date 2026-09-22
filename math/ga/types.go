package ga

import (
	"bytes"
	"fmt"

	"golang.org/x/exp/constraints"
)

type Ordinal interface {
	constraints.Integer | constraints.Float
}

type Signed interface {
	constraints.Signed | constraints.Float
}

type Number interface {
	Ordinal | constraints.Complex
}

type Vec2[T Number] struct{ X, Y T }
type Vec3[T Number] struct{ X, Y, Z T }
type Vec4[T Number] struct{ X, Y, Z, W T }

type Vec2i = Vec2[int]
type Vec2f = Vec2[float32]
type Vec2d = Vec2[float64]
type Vec2z = Vec2[complex128]

type Vec3i = Vec3[int]
type Vec3f = Vec3[float32]
type Vec3d = Vec3[float64]
type Vec3z = Vec3[complex128]

type Vec4i = Vec4[int]
type Vec4f = Vec4[float32]
type Vec4d = Vec4[float64]
type Vec4z = Vec4[complex128]

type Quat[T Number] Vec4[T]

type Quati = Quat[int]
type Quatf = Quat[float32]
type Quatd = Quat[float64]
type Quatz = Quat[complex128]

type Mat2[T Number] [2][2]T
type Mat3[T Number] [3][3]T
type Mat4[T Number] [4][4]T

type Mat2i = Mat2[int]
type Mat2f = Mat2[float32]
type Mat2d = Mat2[float64]
type Mat2z = Mat2[complex128]

type Mat3i = Mat3[int]
type Mat3f = Mat3[float32]
type Mat3d = Mat3[float64]
type Mat3z = Mat3[complex128]

type Mat4i = Mat4[int]
type Mat4f = Mat4[float32]
type Mat4d = Mat4[float64]
type Mat4z = Mat4[complex128]

type Rect2[T Number] struct {
	Min, Max Vec2[T]
}

type Rect2i = Rect2[int]
type Rect2f = Rect2[float32]
type Rect2d = Rect2[float64]
type Rect2z = Rect2[complex128]

type Rect3[T Number] struct {
	Min, Max Vec3[T]
}

type Rect3i = Rect3[int]
type Rect3f = Rect3[float32]
type Rect3d = Rect3[float64]
type Rect3z = Rect3[complex128]

func (v Vec2[T]) String() string {
	return fmt.Sprintf(strfmt("Vec2", 1, 2, v.X), v.X, v.Y)
}

func (v Vec3[T]) String() string {
	return fmt.Sprintf(strfmt("Vec3", 1, 3, v.X), v.X, v.Y, v.Z)
}

func (v Vec4[T]) String() string {
	return fmt.Sprintf(strfmt("Vec4", 1, 4, v.X), v.X, v.Y, v.Z, v.W)
}

func (q Quat[T]) String() string {
	return fmt.Sprintf(strfmt("Quat", 1, 4, q.X), q.X, q.Y, q.Z, q.W)
}

func (m Mat2[T]) String() string {
	return fmt.Sprintf(strfmt("Mat2", 2, 2, m[0][0]),
		m[0][0], m[0][1],
		m[1][0], m[1][1])
}

func (m Mat3[T]) String() string {
	return fmt.Sprintf(strfmt("Mat3", 3, 3, m[0][0]),
		m[0][0], m[0][1], m[0][2],
		m[1][0], m[1][1], m[1][2],
		m[2][0], m[2][1], m[2][2])
}

func (m Mat4[T]) String() string {
	return fmt.Sprintf(strfmt("Mat4", 4, 4, m[0][0]),
		m[0][0], m[0][1], m[0][2], m[0][3],
		m[1][0], m[1][1], m[1][2], m[1][3],
		m[2][0], m[2][1], m[2][2], m[2][3],
		m[3][0], m[3][1], m[3][2], m[3][3])
}

func strfmt(prefix string, rows, cols int, v interface{}) string {
	format := "% v"
	switch v.(type) {
	case float32, float64:
		format = "% 0.3f"
	}

	w := new(bytes.Buffer)

	fmt.Fprintf(w, "%s[", prefix)
	for i := 0; i < rows; i++ {
		if i > 0 {
			fmt.Fprintf(w, "%*s", len(prefix)+1, " ")
		}
		for j := 0; j < cols; j++ {
			fmt.Fprintf(w, "%s", format)
			if j+1 < cols {
				fmt.Fprintf(w, ",")
			}
		}
		if i+1 < rows {
			fmt.Fprintf(w, "\n")
		}
	}
	fmt.Fprintf(w, "]")
	return w.String()
}
