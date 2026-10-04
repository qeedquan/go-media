package lcg

import "math"

type Glibc struct {
	r [34]int32
	i int
}

func (l *Glibc) Seed(seed int32) {
	l.r[0] = seed
	for i := 1; i < 31; i++ {
		l.r[i] = int32(16807 * int64(l.r[i-1]) % 2147483647)
		if l.r[i] < 0 {
			l.r[i] += 2147483647
		}
	}

	for i := 31; i < 34; i++ {
		l.r[i] = l.r[i-31]
	}

	l.i = 34
	for i := 34; i < 344; i++ {
		l.Uint32()
	}
}

func (l *Glibc) Uint32() uint32 {
	n := len(l.r)
	i := mod(l.i-31, n)
	j := mod(l.i-3, n)
	p := mod(l.i, n)

	l.r[p] = l.r[i] + l.r[j]
	l.i = (l.i + 1) % n

	return uint32(l.r[p])
}

func (l *Glibc) Int31() int32 {
	return int32(l.Uint32() >> 1)
}

func (l *Glibc) Int() int {
	return int(l.Int31())
}

func (l *Glibc) Float64() float64 {
	return float64(l.Int()) / math.MaxInt32
}

func mod(x, m int) int {
	x %= m
	if x < 0 {
		x += m
	}
	return x
}
