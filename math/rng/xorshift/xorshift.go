package xorshift

type Rand32 struct {
	a uint32
}

func (r *Rand32) Seed(a uint32) {
	r.a = a
}

func (r *Rand32) Uint32() uint32 {
	x := r.a
	x ^= x << 13
	x ^= x >> 17
	x ^= x << 5
	r.a = x
	return x
}

type Rand64 struct {
	a uint64
}

func (r *Rand64) Seed(a uint64) {
	r.a = a
}

func (r *Rand64) Uint64() uint64 {
	x := r.a
	x ^= x << 13
	x ^= x >> 7
	x ^= x << 17
	r.a = x
	return x
}

type Rand128 struct {
	a [4]uint32
}

func (r *Rand128) Seed(a [4]uint32) {
	r.a = a
}

func (r *Rand128) Uint32() uint32 {
	t := r.a[3]
	s := r.a[0]
	r.a[3] = r.a[2]
	r.a[2] = r.a[1]
	r.a[1] = s

	t ^= t << 11
	t ^= t >> 8
	r.a[0] = t ^ s ^ (s >> 19)

	return r.a[0]
}
