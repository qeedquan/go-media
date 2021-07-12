package mathutil

func Min8(a, b uint8) uint8 {
	if a < b {
		return a
	}
	return b
}

func Min8I(a, b int8) int8 {
	if a < b {
		return a
	}
	return b
}

func Min16(a, b uint16) uint16 {
	if a < b {
		return a
	}
	return b
}

func Min16I(a, b int16) int16 {
	if a < b {
		return a
	}
	return b
}

func Min32(a, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}

func Min32I(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func Min64(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}

func Min64I(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func MinU(a, b uint) uint {
	if a < b {
		return a
	}
	return b
}

func Max8(a, b uint8) uint8 {
	if a > b {
		return a
	}
	return b
}

func Max8I(a, b int8) int8 {
	if a > b {
		return a
	}
	return b
}

func Max16(a, b uint16) uint16 {
	if a > b {
		return a
	}
	return b
}

func Max16I(a, b int16) int16 {
	if a > b {
		return a
	}
	return b
}

func Max32(a, b uint32) uint32 {
	if a > b {
		return a
	}
	return b
}

func Max32I(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func Max64(a, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}

func Max64i(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Maxu(a, b uint) uint {
	if a > b {
		return a
	}
	return b
}

func Min8v(a uint8, b ...uint8) uint8 {
	x := a
	for _, y := range b {
		if y < x {
			x = y
		}
	}
	return x
}

func Min8iv(a int8, b ...int8) int8 {
	x := a
	for _, y := range b {
		if y < x {
			x = y
		}
	}
	return x
}

func Min16v(a uint16, b ...uint16) uint16 {
	x := a
	for _, y := range b {
		if y < x {
			x = y
		}
	}
	return x
}

func Min16iv(a int16, b ...int16) int16 {
	x := a
	for _, y := range b {
		if y < x {
			x = y
		}
	}
	return x
}

func Min32v(a uint32, b ...uint32) uint32 {
	x := a
	for _, y := range b {
		if y < x {
			x = y
		}
	}
	return x
}

func Min32iv(a int32, b ...int32) int32 {
	x := a
	for _, y := range b {
		if y < x {
			x = y
		}
	}
	return x
}

func Min64v(a uint64, b ...uint64) uint64 {
	x := a
	for _, y := range b {
		if y < x {
			x = y
		}
	}
	return x
}

func Min64iv(a int64, b ...int64) int64 {
	x := a
	for _, y := range b {
		if y < x {
			x = y
		}
	}
	return x
}

func Minv(a int, b ...int) int {
	x := a
	for _, y := range b {
		if y < x {
			x = y
		}
	}
	return x
}

func Minuv(a uint, b ...uint) uint {
	x := a
	for _, y := range b {
		if y < x {
			x = y
		}
	}
	return x
}

func Max8v(a uint8, b ...uint8) uint8 {
	x := a
	for _, y := range b {
		if y > x {
			x = y
		}
	}
	return x
}

func Max8iv(a int8, b ...int8) int8 {
	x := a
	for _, y := range b {
		if y > x {
			x = y
		}
	}
	return x
}

func Max16v(a uint16, b ...uint16) uint16 {
	x := a
	for _, y := range b {
		if y > x {
			x = y
		}
	}
	return x
}

func Max16iv(a int16, b ...int16) int16 {
	x := a
	for _, y := range b {
		if y > x {
			x = y
		}
	}
	return x
}

func Max32v(a uint32, b ...uint32) uint32 {
	x := a
	for _, y := range b {
		if y > x {
			x = y
		}
	}
	return x
}

func Max32iv(a int32, b ...int32) int32 {
	x := a
	for _, y := range b {
		if y > x {
			x = y
		}
	}
	return x
}

func Max64v(a uint64, b ...uint64) uint64 {
	x := a
	for _, y := range b {
		if y > x {
			x = y
		}
	}
	return x
}

func Max64iv(a int64, b ...int64) int64 {
	x := a
	for _, y := range b {
		if y > x {
			x = y
		}
	}
	return x
}

func Maxv(a int, b ...int) int {
	x := a
	for _, y := range b {
		if y > x {
			x = y
		}
	}
	return x
}

func Maxuv(a uint, b ...uint) uint {
	x := a
	for _, y := range b {
		if y > x {
			x = y
		}
	}
	return x
}
