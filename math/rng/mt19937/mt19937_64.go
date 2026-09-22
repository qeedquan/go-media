// ported from http://www.math.sci.hiroshima-u.ac.jp/~m-mat/MT/VERSIONS/C-LANG/mt19937-64.c
package mt19937

const (
	NN       = 312
	MM       = 156
	MATRIX_A = 0xB5026F5AA96619E9
	UM       = 0xFFFFFFFF80000000
	LM       = 0x7FFFFFFF
)

type Rand64 struct {
	mt  [NN]uint64
	mti int
}

func New64() *Rand64 {
	return &Rand64{
		mti: NN + 1,
	}
}

func (r *Rand64) Seed(s uint64) {
	r.mt[0] = s
	r.mti = NN
	for i := uint64(1); i < NN; i++ {
		r.mt[i] = (6364136223846793005*(r.mt[i-1]^(r.mt[i-1]>>62)) + i)
	}
}

func (r *Rand64) SeedArray(p []uint64) {
	r.Seed(19650218)
	i, j := uint64(1), uint64(0)
	k := len(p)
	if NN > k {
		k = NN
	}

	for ; k != 0; k-- {
		r.mt[i] = (r.mt[i] ^ ((r.mt[i-1] ^ (r.mt[i-1] >> 62)) * 3935559000370003845)) + p[j] + j
		i++
		j++
		if i >= NN {
			r.mt[0] = r.mt[NN-1]
			i = 1
		}
		if j >= uint64(len(p)) {
			j = 0
		}
	}

	for k = NN - 1; k != 0; k-- {
		r.mt[i] = (r.mt[i] ^ ((r.mt[i-1] ^ (r.mt[i-1] >> 62)) * 2862933555777941757)) - i
		i++
		if i >= NN {
			r.mt[0] = r.mt[NN-1]
			i = 1
		}
	}

	r.mt[0] = 1 << 63
}

func (r *Rand64) Uint64() uint64 {
	x := uint64(0)
	mag01 := [...]uint64{0, MATRIX_A}

	if r.mti >= NN {
		if r.mti == NN+1 {
			r.Seed(5489)
		}

		var i uint64
		for i = 0; i < NN-MM; i++ {
			x = (r.mt[i] & UM) | (r.mt[i+1] & LM)
			r.mt[i] = r.mt[i+MM] ^ (x >> 1) ^ mag01[x&1]
		}

		for ; i < NN-1; i++ {
			x = (r.mt[i] & UM) | (r.mt[i+1] & LM)
			r.mt[i] = r.mt[i+MM-NN] ^ (x >> 1) ^ mag01[x&1]
		}

		x = (r.mt[NN-1] & UM) | (r.mt[0] & LM)
		r.mt[NN-1] = r.mt[MM-1] ^ (x >> 1) ^ mag01[x&1]

		r.mti = 0
	}

	x = r.mt[r.mti]
	r.mti++

	x ^= (x >> 29) & 0x5555555555555555
	x ^= (x << 17) & 0x71D67FFFEDA60000
	x ^= (x << 37) & 0xFFF7EEE000000000
	x ^= (x >> 43)
	return x
}

func (r *Rand64) Int() int {
	return int(r.Int63())
}

func (r *Rand64) Int63() int64 {
	return int64(r.Uint64() >> 1)
}

func (r *Rand64) Float64() float64 {
	return float64(r.Uint64()>>11) * (1.0 / 9007199254740992.0)
}

func (r *Rand64) Real1() float64 {
	return float64(r.Uint64()>>11) * (1.0 / 9007199254740991.0)
}

func (r *Rand64) Real3() float64 {
	return (float64(r.Uint64()>>12) + 0.5) * (1.0 / 4503599627370496.0)
}

func (r *Rand64) Complex64() complex64 {
	return complex64(complex(r.Float64(), r.Float64()))
}

func (r *Rand64) Complex64n(n float32) complex64 {
	return complex64(complex(r.Float64()*float64(n), r.Float64()*float64(n)))
}

func (r *Rand64) Complex128() complex128 {
	return complex(r.Float64(), r.Float64())
}

func (r *Rand64) Complex128n(n float64) complex128 {
	return complex(r.Float64()*n, r.Float64()*n)
}
