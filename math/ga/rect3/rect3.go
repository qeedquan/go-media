package rect3

import "github.com/qeedquan/go-media/math/ga"

func Add[T ga.Number](r ga.Rect3[T], p ga.Vec3[T]) ga.Rect3[T] {
	return ga.Rect3[T]{
		ga.Vec3[T]{r.Min.X + p.X, r.Min.Y + p.Y, r.Min.Z + p.Z},
		ga.Vec3[T]{r.Max.X + p.X, r.Max.Y + p.Y, r.Max.Z + p.Z},
	}
}

func Sub[T ga.Number](r ga.Rect3[T], p ga.Vec3[T]) ga.Rect3[T] {
	return ga.Rect3[T]{
		ga.Vec3[T]{r.Min.X - p.X, r.Min.Y - p.Y, r.Min.Z - p.Z},
		ga.Vec3[T]{r.Max.X - p.X, r.Max.Y - p.Y, r.Max.Z - p.Z},
	}
}

func Empty[T ga.Ordinal](r ga.Rect3[T]) bool {
	return r.Min.X >= r.Max.X || r.Min.Y >= r.Max.Y || r.Min.Z >= r.Max.Z
}

func In[T ga.Ordinal](r, s ga.Rect3[T]) bool {
	if Empty(r) {
		return true
	}
	return s.Min.X <= r.Min.X && r.Max.X <= s.Max.X &&
		s.Min.Y <= r.Min.Y && r.Max.Y <= s.Max.Y &&
		s.Min.Z <= r.Min.Z && r.Max.Z <= s.Max.Z
}

func Intersect[T ga.Ordinal](r, s ga.Rect3[T]) ga.Rect3[T] {
	if r.Min.X < s.Min.X {
		r.Min.X = s.Min.X
	}
	if r.Min.Y < s.Min.Y {
		r.Min.Y = s.Min.Y
	}
	if r.Min.Z < s.Min.Z {
		r.Min.Z = s.Min.Z
	}
	if r.Max.X > s.Max.X {
		r.Max.X = s.Max.X
	}
	if r.Max.Y > s.Max.Y {
		r.Max.Y = s.Max.Y
	}
	if r.Max.Z > s.Max.Z {
		r.Max.Z = s.Max.Z
	}
	if Empty(r) {
		return ga.Rect3[T]{}
	}
	return r
}

func Union[T ga.Ordinal](r, s ga.Rect3[T]) ga.Rect3[T] {
	if Empty(r) {
		return s
	}
	if Empty(s) {
		return r
	}
	if r.Min.X > s.Min.X {
		r.Min.X = s.Min.X
	}
	if r.Min.Y > s.Min.Y {
		r.Min.Y = s.Min.Y
	}
	if r.Min.Z > s.Min.Z {
		r.Min.Z = s.Min.Z
	}
	if r.Max.X < s.Max.X {
		r.Max.X = s.Max.X
	}
	if r.Max.Y < s.Max.Y {
		r.Max.Y = s.Max.Y
	}
	if r.Max.Z < s.Max.Z {
		r.Max.Z = s.Max.Z
	}
	return r
}

func Dim[T ga.Number](r ga.Rect3[T]) (T, T, T) {
	return r.Max.X - r.Min.X, r.Max.Y - r.Min.Y, r.Max.Z - r.Min.Z
}

func Canon[T ga.Ordinal](r ga.Rect3[T]) ga.Rect3[T] {
	if r.Max.X < r.Min.X {
		r.Min.X, r.Max.X = r.Max.X, r.Min.X
	}
	if r.Max.Y < r.Min.Y {
		r.Min.Y, r.Max.Y = r.Max.Y, r.Min.Y
	}
	if r.Max.Z < r.Min.Z {
		r.Min.Z, r.Max.Z = r.Max.Z, r.Min.Z
	}
	return r
}

func Overlaps[T ga.Ordinal](r, s ga.Rect3[T]) bool {
	return !Empty(r) && !Empty(s) &&
		r.Min.X < s.Max.X && s.Min.X < r.Max.X &&
		r.Min.Y < s.Max.Y && s.Min.Y < r.Max.Y &&
		r.Min.Z < s.Max.Z && s.Min.Z < r.Max.Z
}

func Midpoint[T ga.Number](r ga.Rect3[T]) ga.Vec3[T] {
	return ga.Vec3[T]{
		(r.Min.X + r.Max.X) / 2,
		(r.Min.Y + r.Max.Y) / 2,
		(r.Min.Z + r.Max.Z) / 2,
	}
}
