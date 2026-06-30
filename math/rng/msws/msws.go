package msws

// https://arxiv.org/pdf/1704.00358
type Rand32 struct {
	x, w, s uint64
}

func (r *Rand32) Seed(x, w uint64) {
	r.x = x
	r.w = w
	r.s = 0xb5ad4eceda1ce2a9
}

func (r *Rand32) Uint32() uint32 {
	r.x *= r.x
	r.w += r.s
	r.x += r.w
	r.x = (r.x >> 32) | (r.x << 32)
	return uint32(r.x)
}
