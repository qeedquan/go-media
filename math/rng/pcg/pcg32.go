package pcg

type Rand32 struct {
	state uint64
	inc   uint64
}

func (r *Rand32) Seed(initstate, initseq uint64) {
	r.state = 0
	r.inc = initseq<<1 | 1
	r.Uint32()
	r.state += initstate
	r.Uint32()
}

func (r *Rand32) Int() int {
	return int(r.Int31())
}

func (r *Rand32) Int31() int32 {
	return int32(r.Uint32() >> 1)
}

func (r *Rand32) Uint32() uint32 {
	oldstate := r.state
	r.state = oldstate*6364136223846793005 + r.inc

	xorshifted := ((oldstate >> 18) ^ oldstate) >> 27
	rot := oldstate >> 59

	return uint32((xorshifted >> rot) | (xorshifted << ((-rot) & 31)))
}

func (r *Rand32) Intn(n int) int {
	threshold := -n % n
	for {
		if v := r.Int(); v >= threshold {
			return v % n
		}
	}
}
