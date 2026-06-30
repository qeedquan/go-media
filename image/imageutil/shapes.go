package imageutil

import (
	"image"
	"image/color"
	"math"
)

func abs(x int) int {
	if x < 0 {
		x = -x
	}
	return x
}

func Point(m *image.RGBA, x, y, wd int, c color.RGBA) {
	for i := -wd / 2; i <= wd/2; i++ {
		for j := -wd / 2; j <= wd/2; j++ {
			m.SetRGBA(x+j, y+i, c)
		}
	}
}

func Line(m *image.RGBA, x0, y0, x1, y1 int, c color.RGBA) {
	dx := abs(x1 - x0)
	dy := -abs(y1 - y0)
	sx, sy := -1, -1
	if x0 < x1 {
		sx = 1
	}
	if y0 < y1 {
		sy = 1
	}
	e := dx + dy
	for {
		m.SetRGBA(x0, y0, c)
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * e
		if e2 >= dy {
			e += dy
			x0 += sx
		}
		if e2 <= dx {
			e += dx
			y0 += sy
		}
	}
}

func ThickLine(m *image.RGBA, x0, y0, x1, y1 int, wd float64, c color.RGBA) {
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	x2, y2 := 0, 0
	sx, sy := -1, -1
	if x0 < x1 {
		sx = 1
	}
	if y0 < y1 {
		sy = 1
	}
	e := dx - dy
	e2 := 0
	ed := 1.0
	if dx+dy != 0 {
		ed = math.Hypot(float64(dx), float64(dy))
	}

	cc := color.RGBAModel.Convert(c).(color.RGBA)
	wd = (wd + 1) / 2
	for {
		r := uint8(math.Max(0, float64(cc.R)*(math.Abs(float64(e-dx+dy))/ed-wd+1)))
		g := uint8(math.Max(0, float64(cc.G)*(math.Abs(float64(e-dx+dy))/ed-wd+1)))
		b := uint8(math.Max(0, float64(cc.B)*(math.Abs(float64(e-dx+dy))/ed-wd+1)))
		col := color.RGBA{r, g, b, 255}
		m.SetRGBA(x0, y0, col)

		e2 = e
		x2 = x0
		if 2*e2 >= -dx {
			e2 += dy
			y2 = y0
			for float64(e2) < ed*wd && (y1 != y2 || dx > dy) {
				y2 += sy
				r := uint8(math.Max(0, float64(cc.R)*(math.Abs(float64(e2))/ed-wd+1)))
				g := uint8(math.Max(0, float64(cc.G)*(math.Abs(float64(e2))/ed-wd+1)))
				b := uint8(math.Max(0, float64(cc.B)*(math.Abs(float64(e2))/ed-wd+1)))
				col := color.RGBA{r, g, b, 255}
				m.SetRGBA(x0, y2, col)
				e2 += dx
			}
			if x0 == x1 {
				break
			}
			e2 = e
			e -= dy
			x0 += sx
		}
		if 2*e2 <= dy {
			e2 = dx - e2
			for float64(e2) < ed*wd && (x1 != x2 || dx < dy) {
				x2 += sx
				r := uint8(math.Max(0, float64(cc.R)*(math.Abs(float64(e2))/ed-wd+1)))
				g := uint8(math.Max(0, float64(cc.G)*(math.Abs(float64(e2))/ed-wd+1)))
				b := uint8(math.Max(0, float64(cc.B)*(math.Abs(float64(e2))/ed-wd+1)))
				col := color.RGBA{r, g, b, 255}
				m.SetRGBA(x2, y0, col)
				e2 += dy
			}
			if y0 == y1 {
				break
			}
			e += dx
			y0 += sy
		}
	}
}

func Circle(m *image.RGBA, cx, cy, r int, c color.RGBA) {
	ThickCircle(m, cx, cy, r, 1, c)
}

func ThickCircle(m *image.RGBA, cx, cy, r, wd int, c color.RGBA) {
	x := -r
	y := 0
	e := 2 - 2*r
	for {
		Point(m, cx-x, cy+y, wd, c)
		Point(m, cx-y, cy-x, wd, c)
		Point(m, cx+x, cy-y, wd, c)
		Point(m, cx+y, cy+x, wd, c)
		r := e
		if r <= y {
			y++
			e += 2*y + 1
		}
		if r > x || e > y {
			x++
			e += 2*x + 1
		}
		if x > 0 {
			break
		}
	}
}

func FilledCircle(m *image.RGBA, cx, cy, r int, c color.RGBA) {
	for y := -r; y <= r; y++ {
		for x := -r; x <= r; x++ {
			if x*x+y*y <= r*r {
				m.SetRGBA(cx+x, cy+y, c)
			}
		}
	}
}

func StrokeFilledCircle(m *image.RGBA, cx, cy, r, wd int, bc, ic color.RGBA) {
	FilledCircle(m, cx, cy, r, bc)
	ThickCircle(m, cx, cy, r, wd, ic)
}

func Triangle(m *image.RGBA, x0, y0, x1, y1, x2, y2 int, c color.RGBA) {
	Line(m, x0, y0, x1, y1, c)
	Line(m, x0, y0, x2, y2, c)
	Line(m, x1, y1, x2, y2, c)
}

func FilledTriangle(m *image.RGBA, x0, y0, x1, y1, x2, y2 int, c color.RGBA) {
	if y0 > y1 {
		x0, y0, x1, y1 = x1, y1, x0, y0
	}
	if y0 > y2 {
		x0, y0, x2, y2 = x2, y2, x0, y0
	}
	if y1 > y2 {
		x1, y1, x2, y2 = x2, y2, x1, y1
	}
	filledTriangleBottom(m, x0, y0, x1, y1, x2, y2, c)
	filledTriangleTop(m, x0, y0, x1, y1, x2, y2, c)
}

func filledTriangleBottom(m *image.RGBA, x0, y0, x1, y1, x2, y2 int, c color.RGBA) {
	i1 := float64(x1-x0) / float64(y1-y0)
	i2 := float64(x2-x0) / float64(y2-y0)

	cx1 := float64(x0)
	cx2 := cx1
	for y := y0; y <= y1; y++ {
		Line(m, int(cx1), y, int(cx2), y, c)
		cx1 += i1
		cx2 += i2
	}
}

func filledTriangleTop(m *image.RGBA, x0, y0, x1, y1, x2, y2 int, c color.RGBA) {
	i1 := float64(x2-x0) / float64(y2-y0)
	i2 := float64(x2-x1) / float64(y2-y1)

	cx1 := float64(x2)
	cx2 := cx1
	for y := y2; y > y0; y-- {
		Line(m, int(cx1), y, int(cx2), y, c)
		cx1 -= i1
		cx2 -= i2
	}
}

func EllipseRect(m *image.RGBA, x0, y0, x1, y1 int, c color.RGBA) {
	a := abs(x1 - x0)
	b := abs(y1 - y0)
	b1 := b & 1
	dx := 4 * (1 - a) * b * b
	dy := 4 * (b1 + 1) * a * a
	e := dx + dy + b1*a*a
	if x0 > x1 {
		x0 = x1
		x1 += a
	}
	if y0 > y1 {
		y0 = y1
	}
	y0 += (b + 1) / 2
	y1 = y0 - b1
	a *= 8 * a
	b1 = 8 * b * b
	for {
		m.SetRGBA(x1, y0, c)
		m.SetRGBA(x0, y0, c)
		m.SetRGBA(x0, y1, c)
		m.SetRGBA(x1, y1, c)
		e2 := 2 * e
		if e2 <= dy {
			y0++
			y1--
			dy += a
			e += dy
		}
		if e2 >= dx || 2*e > dy {
			x0++
			x1--
			dx += b1
			e += dx
		}
		if x0 > x1 {
			break
		}
	}
	for y0-y1 < b {
		m.SetRGBA(x0-1, y0, c)
		m.SetRGBA(x1+1, y0, c)
		y0++
		m.SetRGBA(x0-1, y1, c)
		m.SetRGBA(x1+1, y1, c)
		y1--
	}
}
