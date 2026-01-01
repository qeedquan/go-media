// ported from http://staffwww.itn.liu.se/~stegu/aqsis/aqsis-newnoise/simplexnoise1234.cpp
// http://staffwww.itn.liu.se/~stegu/simplexnoise/simplexnoise.pdf
package simplx

func Noise1D(x float64) float64 {
	i0 := floor(x)
	i1 := i0 + 1
	x0 := x - float64(i0)
	x1 := x0 - 1

	t0 := 1 - x0*x0
	t0 *= t0
	n0 := t0 * t0 * grad1(perm[i0&0xff], x0)

	t1 := 1 - x1*x1
	t1 *= t1
	n1 := t1 * t1 * grad1(perm[i1&0xff], x1)
	return 0.395 * (n0 + n1)
}

func Noise2D(x, y float64) float64 {
	const F2 = 0.366025403
	const G2 = 0.211324865

	s := (x + y) * F2
	xs := x + s
	ys := y + s
	i := floor(xs)
	j := floor(ys)

	t := float64(i+j) * G2
	X0 := float64(i) - t
	Y0 := float64(j) - t
	x0 := x - X0
	y0 := y - Y0

	var i1, j1 int
	if x0 > y0 {
		i1 = 1
		j1 = 0
	} else {
		i1 = 0
		j1 = 1
	}

	x1 := x0 - float64(i1) + G2
	y1 := y0 - float64(j1) + G2
	x2 := x0 - 1 + 2*G2
	y2 := y0 - 1 + 2*G2

	ii := i & 0xff
	jj := j & 0xff

	var n0, n1, n2 float64

	t0 := 0.5 - x0*x0 - y0*y0
	if t0 < 0 {
		n0 = 0
	} else {
		t0 *= t0
		n0 = t0 * t0 * grad2(perm[ii+perm[jj]], x0, y0)
	}

	t1 := 0.5 - x1*x1 - y1*y1
	if t1 < 0 {
		n1 = 0
	} else {
		t1 *= t1
		n1 = t1 * t1 * grad2(perm[ii+i1+perm[jj+j1]], x1, y1)
	}

	t2 := 0.5 - x2*x2 - y2*y2
	if t2 < 0 {
		n2 = 0
	} else {
		t2 *= t2
		n2 = t2 * t2 * grad2(perm[ii+1+perm[jj+1]], x2, y2)
	}

	return 40 * (n0 + n1 + n2)
}

func Noise3D(x, y, z float64) float64 {
	const F3 = 0.333333333
	const G3 = 0.166666667

	s := (x + y + z) * F3
	xs := x + s
	ys := y + s
	zs := z + s
	i := floor(xs)
	j := floor(ys)
	k := floor(zs)

	t := float64(i+j+k) * G3
	X0 := float64(i) - t
	Y0 := float64(j) - t
	Z0 := float64(k) - t
	x0 := x - X0
	y0 := y - Y0
	z0 := z - Z0

	var i1, j1, k1 int
	var i2, j2, k2 int

	if x0 >= y0 {
		if y0 >= z0 {
			i1 = 1
			j1 = 0
			k1 = 0
			i2 = 1
			j2 = 1
			k2 = 0
		} else if x0 >= z0 {
			i1 = 1
			j1 = 0
			k1 = 0
			i2 = 1
			j2 = 0
			k2 = 1
		} else {
			i1 = 0
			j1 = 0
			k1 = 1
			i2 = 1
			j2 = 0
			k2 = 1
		}
	} else {
		if y0 < z0 {
			i1 = 0
			j1 = 0
			k1 = 1
			i2 = 0
			j2 = 1
			k2 = 1
		} else if x0 < z0 {
			i1 = 0
			j1 = 1
			k1 = 0
			i2 = 0
			j2 = 1
			k2 = 1
		} else {
			i1 = 0
			j1 = 1
			k1 = 0
			i2 = 1
			j2 = 1
			k2 = 0
		}
	}

	x1 := x0 - float64(i1) + G3
	y1 := y0 - float64(j1) + G3
	z1 := z0 - float64(k1) + G3
	x2 := x0 - float64(i2) + 2.0*G3
	y2 := y0 - float64(j2) + 2.0*G3
	z2 := z0 - float64(k2) + 2.0*G3
	x3 := x0 - 1.0 + 3.0*G3
	y3 := y0 - 1.0 + 3.0*G3
	z3 := z0 - 1.0 + 3.0*G3

	ii := i & 0xff
	jj := j & 0xff
	kk := k & 0xff

	var n0, n1, n2, n3 float64

	t0 := 0.6 - x0*x0 - y0*y0 - z0*z0
	if t0 < 0.0 {
		n0 = 0.0
	} else {
		t0 *= t0
		n0 = t0 * t0 * grad3(perm[ii+perm[jj+perm[kk]]], x0, y0, z0)
	}

	t1 := 0.6 - x1*x1 - y1*y1 - z1*z1
	if t1 < 0.0 {
		n1 = 0.0
	} else {
		t1 *= t1
		n1 = t1 * t1 * grad3(perm[ii+i1+perm[jj+j1+perm[kk+k1]]], x1, y1, z1)
	}

	t2 := 0.6 - x2*x2 - y2*y2 - z2*z2
	if t2 < 0.0 {
		n2 = 0.0
	} else {
		t2 *= t2
		n2 = t2 * t2 * grad3(perm[ii+i2+perm[jj+j2+perm[kk+k2]]], x2, y2, z2)
	}

	t3 := 0.6 - x3*x3 - y3*y3 - z3*z3
	if t3 < 0.0 {
		n3 = 0.0
	} else {
		t3 *= t3
		n3 = t3 * t3 * grad3(perm[ii+1+perm[jj+1+perm[kk+1]]], x3, y3, z3)
	}

	return 32.0 * (n0 + n1 + n2 + n3)
}

func Noise4D(x, y, z, w float64) float64 {
	const F4 = 0.309016994
	const G4 = 0.138196601

	s := (x + y + z + w) * F4
	xs := x + s
	ys := y + s
	zs := z + s
	ws := w + s
	i := floor(xs)
	j := floor(ys)
	k := floor(zs)
	l := floor(ws)

	t := float64(i+j+k+l) * G4
	X0 := float64(i) - t
	Y0 := float64(j) - t
	Z0 := float64(k) - t
	W0 := float64(l) - t

	x0 := x - X0
	y0 := y - Y0
	z0 := z - Z0
	w0 := w - W0

	var c1, c2, c3, c4, c5, c6 int
	if x0 > y0 {
		c1 = 32
	}
	if x0 > z0 {
		c2 = 16
	}
	if y0 > z0 {
		c3 = 8
	}
	if x0 > w0 {
		c4 = 4
	}
	if y0 > w0 {
		c5 = 2
	}
	if z0 > w0 {
		c6 = 1
	}
	c := c1 + c2 + c3 + c4 + c5 + c6

	var i1, j1, k1, l1 int
	var i2, j2, k2, l2 int
	var i3, j3, k3, l3 int

	if simplex[c][0] >= 3 {
		i1 = 1
	}
	if simplex[c][1] >= 3 {
		j1 = 1
	}
	if simplex[c][2] >= 3 {
		k1 = 1
	}
	if simplex[c][3] >= 3 {
		l1 = 1
	}

	if simplex[c][0] >= 2 {
		i2 = 1
	}
	if simplex[c][1] >= 2 {
		j2 = 1
	}
	if simplex[c][2] >= 2 {
		k2 = 1
	}
	if simplex[c][3] >= 2 {
		l2 = 1
	}

	if simplex[c][0] >= 1 {
		i3 = 1
	}
	if simplex[c][1] >= 1 {
		j3 = 1
	}
	if simplex[c][2] >= 1 {
		k3 = 1
	}
	if simplex[c][3] >= 1 {
		l3 = 1
	}

	x1 := x0 - float64(i1) + G4
	y1 := y0 - float64(j1) + G4
	z1 := z0 - float64(k1) + G4
	w1 := w0 - float64(l1) + G4
	x2 := x0 - float64(i2) + 2.0*G4
	y2 := y0 - float64(j2) + 2.0*G4
	z2 := z0 - float64(k2) + 2.0*G4
	w2 := w0 - float64(l2) + 2.0*G4
	x3 := x0 - float64(i3) + 3.0*G4
	y3 := y0 - float64(j3) + 3.0*G4
	z3 := z0 - float64(k3) + 3.0*G4
	w3 := w0 - float64(l3) + 3.0*G4
	x4 := x0 - 1.0 + 4.0*G4
	y4 := y0 - 1.0 + 4.0*G4
	z4 := z0 - 1.0 + 4.0*G4
	w4 := w0 - 1.0 + 4.0*G4

	ii := i & 0xff
	jj := j & 0xff
	kk := k & 0xff
	ll := l & 0xff

	var n0, n1, n2, n3, n4 float64
	t0 := 0.6 - x0*x0 - y0*y0 - z0*z0 - w0*w0
	if t0 < 0.0 {
		n0 = 0.0
	} else {
		t0 *= t0
		n0 = t0 * t0 * grad4(perm[ii+perm[jj+perm[kk+perm[ll]]]], x0, y0, z0, w0)
	}

	t1 := 0.6 - x1*x1 - y1*y1 - z1*z1 - w1*w1
	if t1 < 0.0 {
		n1 = 0.0
	} else {
		t1 *= t1
		n1 = t1 * t1 * grad4(perm[ii+i1+perm[jj+j1+perm[kk+k1+perm[ll+l1]]]], x1, y1, z1, w1)
	}

	t2 := 0.6 - x2*x2 - y2*y2 - z2*z2 - w2*w2
	if t2 < 0.0 {
		n2 = 0.0
	} else {
		t2 *= t2
		n2 = t2 * t2 * grad4(perm[ii+i2+perm[jj+j2+perm[kk+k2+perm[ll+l2]]]], x2, y2, z2, w2)
	}

	t3 := 0.6 - x3*x3 - y3*y3 - z3*z3 - w3*w3
	if t3 < 0.0 {
		n3 = 0.0
	} else {
		t3 *= t3
		n3 = t3 * t3 * grad4(perm[ii+i3+perm[jj+j3+perm[kk+k3+perm[ll+l3]]]], x3, y3, z3, w3)
	}

	t4 := 0.6 - x4*x4 - y4*y4 - z4*z4 - w4*w4
	if t4 < 0.0 {
		n4 = 0.0
	} else {
		t4 *= t4
		n4 = t4 * t4 * grad4(perm[ii+1+perm[jj+1+perm[kk+1+perm[ll+1]]]], x4, y4, z4, w4)
	}

	return 27.0 * (n0 + n1 + n2 + n3 + n4)
}

func grad1(hash int, x float64) float64 {
	h := hash & 15
	grad := 1 + float64(h&7)
	if h&8 != 0 {
		grad = -grad
	}
	return grad * x
}

func grad2(hash int, x, y float64) float64 {
	h := hash & 7
	u, v := x, y
	if !(h < 4) {
		u, v = y, x
	}
	if h&1 != 0 {
		u = -u
	}
	if h&2 != 0 {
		v = -v
	}
	return u + 2*v
}

func grad3(hash int, x, y, z float64) float64 {
	h := hash & 15
	u, v := x, y
	if !(h < 8) {
		u = y
	}
	if !(h < 4) {
		if h == 12 || h == 14 {
			v = x
		} else {
			v = z
		}
	}
	if h&1 != 0 {
		u = -u
	}
	if h&2 != 0 {
		v = -v
	}
	return u + v
}

func grad4(hash int, x, y, z, t float64) float64 {
	h := hash & 31
	u, v, w := x, y, z
	if !(h < 24) {
		u = y
	}
	if !(h < 16) {
		v = z
	}
	if !(h < 8) {
		w = t
	}

	if h&1 != 0 {
		u = -u
	}

	if h&2 != 0 {
		v = -v
	}

	if h&4 != 0 {
		w = -w
	}

	return u + v + w
}

func floor(x float64) int {
	if x > 0 {
		return int(x)
	}
	return int(x) - 1
}

var perm = [512]int{
	151, 160, 137, 91, 90, 15,
	131, 13, 201, 95, 96, 53, 194, 233, 7, 225, 140, 36, 103, 30, 69, 142, 8, 99, 37, 240, 21, 10, 23,
	190, 6, 148, 247, 120, 234, 75, 0, 26, 197, 62, 94, 252, 219, 203, 117, 35, 11, 32, 57, 177, 33,
	88, 237, 149, 56, 87, 174, 20, 125, 136, 171, 168, 68, 175, 74, 165, 71, 134, 139, 48, 27, 166,
	77, 146, 158, 231, 83, 111, 229, 122, 60, 211, 133, 230, 220, 105, 92, 41, 55, 46, 245, 40, 244,
	102, 143, 54, 65, 25, 63, 161, 1, 216, 80, 73, 209, 76, 132, 187, 208, 89, 18, 169, 200, 196,
	135, 130, 116, 188, 159, 86, 164, 100, 109, 198, 173, 186, 3, 64, 52, 217, 226, 250, 124, 123,
	5, 202, 38, 147, 118, 126, 255, 82, 85, 212, 207, 206, 59, 227, 47, 16, 58, 17, 182, 189, 28, 42,
	223, 183, 170, 213, 119, 248, 152, 2, 44, 154, 163, 70, 221, 153, 101, 155, 167, 43, 172, 9,
	129, 22, 39, 253, 19, 98, 108, 110, 79, 113, 224, 232, 178, 185, 112, 104, 218, 246, 97, 228,
	251, 34, 242, 193, 238, 210, 144, 12, 191, 179, 162, 241, 81, 51, 145, 235, 249, 14, 239, 107,
	49, 192, 214, 31, 181, 199, 106, 157, 184, 84, 204, 176, 115, 121, 50, 45, 127, 4, 150, 254,
	138, 236, 205, 93, 222, 114, 67, 29, 24, 72, 243, 141, 128, 195, 78, 66, 215, 61, 156, 180,
	151, 160, 137, 91, 90, 15,
	131, 13, 201, 95, 96, 53, 194, 233, 7, 225, 140, 36, 103, 30, 69, 142, 8, 99, 37, 240, 21, 10, 23,
	190, 6, 148, 247, 120, 234, 75, 0, 26, 197, 62, 94, 252, 219, 203, 117, 35, 11, 32, 57, 177, 33,
	88, 237, 149, 56, 87, 174, 20, 125, 136, 171, 168, 68, 175, 74, 165, 71, 134, 139, 48, 27, 166,
	77, 146, 158, 231, 83, 111, 229, 122, 60, 211, 133, 230, 220, 105, 92, 41, 55, 46, 245, 40, 244,
	102, 143, 54, 65, 25, 63, 161, 1, 216, 80, 73, 209, 76, 132, 187, 208, 89, 18, 169, 200, 196,
	135, 130, 116, 188, 159, 86, 164, 100, 109, 198, 173, 186, 3, 64, 52, 217, 226, 250, 124, 123,
	5, 202, 38, 147, 118, 126, 255, 82, 85, 212, 207, 206, 59, 227, 47, 16, 58, 17, 182, 189, 28, 42,
	223, 183, 170, 213, 119, 248, 152, 2, 44, 154, 163, 70, 221, 153, 101, 155, 167, 43, 172, 9,
	129, 22, 39, 253, 19, 98, 108, 110, 79, 113, 224, 232, 178, 185, 112, 104, 218, 246, 97, 228,
	251, 34, 242, 193, 238, 210, 144, 12, 191, 179, 162, 241, 81, 51, 145, 235, 249, 14, 239, 107,
	49, 192, 214, 31, 181, 199, 106, 157, 184, 84, 204, 176, 115, 121, 50, 45, 127, 4, 150, 254,
	138, 236, 205, 93, 222, 114, 67, 29, 24, 72, 243, 141, 128, 195, 78, 66, 215, 61, 156, 180,
}

var simplex = [64][4]int{
	{0, 1, 2, 3}, {0, 1, 3, 2}, {0, 0, 0, 0}, {0, 2, 3, 1}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {1, 2, 3, 0},
	{0, 2, 1, 3}, {0, 0, 0, 0}, {0, 3, 1, 2}, {0, 3, 2, 1}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {1, 3, 2, 0},
	{0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0},
	{1, 2, 0, 3}, {0, 0, 0, 0}, {1, 3, 0, 2}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {2, 3, 0, 1}, {2, 3, 1, 0},
	{1, 0, 2, 3}, {1, 0, 3, 2}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {2, 0, 3, 1}, {0, 0, 0, 0}, {2, 1, 3, 0},
	{0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0},
	{2, 0, 1, 3}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {3, 0, 1, 2}, {3, 0, 2, 1}, {0, 0, 0, 0}, {3, 1, 2, 0},
	{2, 1, 0, 3}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {3, 1, 0, 2}, {0, 0, 0, 0}, {3, 2, 0, 1}, {3, 2, 1, 0},
}
