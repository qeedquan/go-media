package rect2

import "github.com/qeedquan/go-media/math/ga"

func Add[T ga.Number](r ga.Rect2[T], p ga.Vec2[T]) ga.Rect2[T] {
	return ga.Rect2[T]{
		ga.Vec2[T]{r.Min.X + p.X, r.Min.Y + p.Y},
		ga.Vec2[T]{r.Max.X + p.X, r.Max.Y + p.Y},
	}
}

func Sub[T ga.Number](r ga.Rect2[T], p ga.Vec2[T]) ga.Rect2[T] {
	return ga.Rect2[T]{
		ga.Vec2[T]{r.Min.X - p.X, r.Min.Y - p.Y},
		ga.Vec2[T]{r.Max.X - p.X, r.Max.Y - p.Y},
	}
}

func Empty[T ga.Ordinal](r ga.Rect2[T]) bool {
	return r.Min.X >= r.Max.X || r.Min.Y >= r.Max.Y
}

func In[T ga.Ordinal](r, s ga.Rect2[T]) bool {
	if Empty(r) {
		return true
	}
	return s.Min.X <= r.Min.X && r.Max.X <= s.Max.X &&
		s.Min.Y <= r.Min.Y && r.Max.Y <= s.Max.Y
}

func Intersect[T ga.Ordinal](r, s ga.Rect2[T]) ga.Rect2[T] {
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
	if Empty(r) {
		return ga.Rect2[T]{}
	}
	return r
}

func Union[T ga.Ordinal](r, s ga.Rect2[T]) ga.Rect2[T] {
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
	if r.Max.X < s.Max.X {
		r.Max.X = s.Max.X
	}
	if r.Max.Y < s.Max.Y {
		r.Max.Y = s.Max.Y
	}
	return r
}

func Inset[T ga.Ordinal](r ga.Rect2[T], n T) ga.Rect2[T] {
	w, h := Dim(r)
	if w < 2*n {
		r.Min.X = (r.Min.X + r.Max.X) / 2
		r.Max.X = r.Min.X
	} else {
		r.Min.X += n
		r.Max.X -= n
	}
	if h < 2*n {
		r.Min.Y = (r.Min.Y + r.Max.Y) / 2
		r.Max.Y = r.Min.Y
	} else {
		r.Min.Y += n
		r.Max.Y -= n
	}
	return r
}

func Dim[T ga.Number](r ga.Rect2[T]) (T, T) {
	return r.Max.X - r.Min.X, r.Max.Y - r.Min.Y
}

func Canon[T ga.Ordinal](r ga.Rect2[T]) ga.Rect2[T] {
	if r.Max.X < r.Min.X {
		r.Min.X, r.Max.X = r.Max.X, r.Min.X
	}
	if r.Max.Y < r.Min.Y {
		r.Min.Y, r.Max.Y = r.Max.Y, r.Min.Y
	}
	return r
}

func Overlaps[T ga.Ordinal](r, s ga.Rect2[T]) bool {
	return !Empty(r) && !Empty(s) &&
		r.Min.X < s.Max.X && s.Min.X < r.Max.X &&
		r.Min.Y < s.Max.Y && s.Min.Y < r.Max.Y
}

func Midpoint[T ga.Number](r ga.Rect2[T]) ga.Vec2[T] {
	return ga.Vec2[T]{
		(r.Min.X + r.Max.X) / 2,
		(r.Min.Y + r.Max.Y) / 2,
	}
}
