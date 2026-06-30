// ported from http://www.math.sci.hiroshima-u.ar.jp/~m-mat/MT/MT2002/CODES/mt19937ar.c
package mt19937

const (
	N           = 624
	M           = 397
	MATRIX_A_32 = 0x9908b0df
	UPPER_MASK  = 0x80000000
	LOWER_MASK  = 0x7fffffff
)

type Rand32 struct {
	mt  []uint64
	mti int
}

func New32() *Rand32 {
	return &Rand32{
		mt:  make([]uint64, N),
		mti: N + 1,
	}
}

func (r *Rand32) Seed(s uint64) {
	r.mt[0] = s & 0xffffffff
	r.mti = N
	for i := uint64(1); i < uint64(len(r.mt)); i++ {
		r.mt[i] = (1812433253*(r.mt[i-1]^(r.mt[i-1]>>30)) + i)
		r.mt[i] &= 0xffffffff
	}
}

func (r *Rand32) SeedArray(p []uint64) {
	r.Seed(19650218)

	i := uint64(1)
	j := uint64(0)
	k := len(p)
	if N > k {
		k = N
	}

	for ; k != 0; k-- {
		r.mt[i] = (r.mt[i] ^ ((r.mt[i-1] ^ (r.mt[i-1] >> 30)) * 1664525)) + p[j] + j
		r.mt[i] &= 0xffffffff
		i++
		j++
		if i >= N {
			r.mt[0] = r.mt[N-1]
			i = 1
		}
		if j >= uint64(len(p)) {
			j = 0
		}
	}

	for k = N - 1; k != 0; k-- {
		r.mt[i] = (r.mt[i] ^ ((r.mt[i-1] ^ (r.mt[i-1] >> 30)) * 1566083941)) - i
		r.mt[i] &= 0xffffffff
		i++
		if i >= N {
			r.mt[0] = r.mt[N-1]
			i = 1
		}
	}

	r.mt[0] = 0x80000000
}

func (r *Rand32) Uint32() uint32 {
	y := uint64(0)
	mag01 := [2]uint64{0, MATRIX_A_32}
	if r.mti >= N {
		if r.mti == N+1 {
			r.Seed(5489)
		}

		kk := 0
		for kk = 0; kk < N-M; kk++ {
			y = (r.mt[kk] & UPPER_MASK) | (r.mt[kk+1] & LOWER_MASK)
			r.mt[kk] = r.mt[kk+M] ^ (y >> 1) ^ mag01[y&0x1]
		}

		for ; kk < N-1; kk++ {
			y = (r.mt[kk] & UPPER_MASK) | (r.mt[kk+1] & LOWER_MASK)
			r.mt[kk] = r.mt[kk+(M-N)] ^ (y >> 1) ^ mag01[y&0x1]
		}

		y = (r.mt[N-1] & UPPER_MASK) | (r.mt[0] & LOWER_MASK)
		r.mt[N-1] = r.mt[M-1] ^ (y >> 1) ^ mag01[y&0x1]

		r.mti = 0
	}

	y = r.mt[r.mti]
	r.mti++

	y ^= (y >> 11)
	y ^= (y << 7) & 0x9d2c5680
	y ^= (y << 15) & 0xefc60000
	y ^= (y >> 18)

	return uint32(y)
}

func (r *Rand32) Int() int {
	return int(r.Int31())
}

func (r *Rand32) Int31() int32 {
	return int32(r.Uint32() >> 1)
}

func (r *Rand32) Float32() float32 {
	return float32(r.Uint32()) * (1.0 / 4294967296.0)
}

func (r *Rand32) Real1() float32 {
	return float32(r.Uint32()) * (1.0 / 4294967295.0)
}

func (r *Rand32) Real3() float32 {
	return float32((float64(r.Uint32()) + 0.5) * (1.0 / 4294967296.0))
}

func (r *Rand32) Float64() float64 {
	a := r.Uint32() >> 5
	b := r.Uint32() >> 6
	return (float64(a)*67108864.0 + float64(b)) * (1.0 / 9007199254740992.0)
}

func (r *Rand32) Complex64() complex64 {
	return complex(r.Float32(), r.Float32())
}

func (r *Rand32) Complex64n(n float32) complex64 {
	return complex(r.Float32()*n, r.Float32()*n)
}

func (r *Rand32) Complex128() complex128 {
	return complex(r.Float64(), r.Float64())
}

func (r *Rand32) Complex128n(n float64) complex128 {
	return complex(r.Float64()*n, r.Float64()*n)
}
