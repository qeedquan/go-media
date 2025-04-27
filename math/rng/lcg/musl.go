package lcg

import "math"

type Musl struct {
	val uint64
}

func (l *Musl) Seed(seed uint64) {
	l.val = seed - 1
}

func (l *Musl) Int31() int32 {
	l.val = 6364136223846793005*l.val + 1
	return int32(l.val >> 33)
}

func (l *Musl) Int() int {
	return int(l.Int31())
}

func (l *Musl) Float64() float64 {
	return float64(l.Int31()) / math.MaxInt32
}
