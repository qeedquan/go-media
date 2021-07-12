package mathutil

func GCD(a, b int) int {
	k := Max(a, b)
	m := Min(a, b)
	for m != 0 {
		r := k % m
		k = m
		m = r
	}
	return k
}

func LCM(a, b int) int {
	return Abs(a*b) / GCD(a, b)
}

func Divceil(n, d int) int         { return (n + d - 1) / d }
func Divceil16(n, d uint16) uint16 { return (n + d - 1) / d }
func Divceil32(n, d uint32) uint32 { return (n + d - 1) / d }
func Divceil64(n, d uint64) uint64 { return (n + d - 1) / d }

func Multiple(a, m int) int         { return ((a + m - 1) / m) * m }
func Multiple16(a, m uint16) uint16 { return ((a + m - 1) / m) * m }
func Multiple32(a, m uint32) uint32 { return ((a + m - 1) / m) * m }
func Multiple64(a, m uint64) uint64 { return ((a + m - 1) / m) * m }

func IsPow2(v int) bool {
	return v&(v-1) == 0
}

func NextPow2(v int) int {
	x := 1
	for x < v {
		x <<= 1
	}
	return x
}

func NearestPow2(v int) int {
	if v <= 0 {
		return -1
	}
	for i := 1; ; {
		if v == 1 {
			return i
		}
		if v == 3 {
			return i * 4
		}
		v >>= 1
		i *= 2
	}
}

func Log2(v int) int {
	if v <= 0 {
		return -1
	}

	for i := 0; ; i++ {
		if v&1 != 0 {
			if v != 1 {
				return -1
			}
			return i
		}
		v >>= 1
	}
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func Abs8(x int8) int8 {
	if x < 0 {
		return -x
	}
	return x
}

func Abs16(x int16) int16 {
	if x < 0 {
		return -x
	}
	return x
}

func Abs32(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}

func Abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func Clamp(x, a, b int) int {
	if x < a {
		x = a
	}
	if x > b {
		x = b
	}
	return x
}
