package nipponichi

type Frand struct {
	rnum           uint64
	rn, r1, r2, r3 uint64
}

func divqr64(a, b uint64) (quo, rem uint64) {
	if b == 0 {
		return 0, 0
	}
	return a / b, a % b
}

func mulshi128(a, b uint64) uint64 {
	_, hi := mul128(a, b)
	return hi >> 31
}

func mul128(a, b uint64) (lo, hi uint64) {
	a_lo := uint64(uint32(a))
	a_hi := uint64(a >> 32)
	b_lo := uint64(uint32(b))
	b_hi := uint64(b >> 32)

	a_x_b_hi := a_hi * b_hi
	a_x_b_mid := a_hi * b_lo
	b_x_a_mid := b_hi * a_lo
	a_x_b_lo := a_lo * b_lo

	carry_bit := (uint64(uint32(a_x_b_mid)) +
		uint64(uint32(b_x_a_mid)) +
		(a_x_b_lo >> 32)) >> 32

	lo = a * b
	hi = a_x_b_hi +
		(a_x_b_mid >> 32) + (b_x_a_mid >> 32) +
		carry_bit

	return
}

func (fr *Frand) Random(limit uint32) uint64 {
	const K = 0x8000000080000001

	var H uint64

	if fr.rn++; fr.rn <= 2 {
		goto L3
	}
	fr.rn = 0
L1:
	H = 0x8000000080000001
L2:
	fr.rnum = fr.rnum*0x36d + 0xe021
	H = mulshi128(fr.rnum, H)
	fr.r2 = fr.rnum - ((H << 32) - H)
	fr.rnum = fr.rnum*0x36d + 0xe021

	H = mulshi128(fr.rnum, K)
	fr.r3 = fr.rnum - ((H << 32) - H)

	if limit == 0 {
		return 0
	}
	_, H = divqr64(fr.r1+fr.r2+fr.r3, uint64(limit))
	return H
L3:
	if fr.rn == 0 {
		goto L1
	}
	fr.rnum = fr.rnum*0x36d + 0xe021
	H = mulshi128(fr.rnum, K)
	H = fr.rnum - ((H << 32) - H)
	fr.r1 = H
	if fr.rn == 0x1 {
		goto L2
	}

	fr.rnum = fr.rnum*0x36d + 0xe021
	H = mulshi128(fr.rnum, K)
	fr.r2 = fr.rnum - ((H << 32) - H)

	if limit == 0 {
		return 0
	}
	_, H = divqr64(fr.r1+fr.r2+fr.r3, uint64(limit))
	return H
}

// Seed is a linear sequence of the form
// y = a*x + b, to see this, do x_n+1 - x_n
// to get the constants.
func (fr *Frand) Seed(seed uint64) {
	const K = 0x8000000080000001

	fr.rn = 0

	A := seed*0x36d + 0xe021
	H := mulshi128(A, K)
	fr.r1 = A - ((H << 32) - H)

	A = (A * 0x36d) + 0xe021
	H = mulshi128(A, K)
	fr.r2 = A - ((H << 32) - H)

	A = (A * 0x36d) + 0xe021
	fr.rnum = A

	H = mulshi128(A, K)
	fr.r3 = A - ((H << 32) - H)
}

func (fr *Frand) GetSeed() uint64 {
	return fr.rnum
}
