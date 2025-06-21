package nipponichi

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type Slide struct{}

func (c *Slide) Encode2(src []byte) ([]byte, error) {
	b := new(bytes.Buffer)
	binary.Write(b, binary.LittleEndian, uint32(0x0))
	binary.Write(b, binary.LittleEndian, uint32(0x0))
	binary.Write(b, binary.LittleEndian, uint32(0x4))
	binary.Write(b, binary.LittleEndian, uint32(0x0))
	binary.Write(b, binary.LittleEndian, uint32(0x0))

	dn := 0
	sn := 0x14
	for i := 0; i < len(src); {
		n := len(src) - i
		if n >= 64 {
			n = 64
		}
		b.WriteByte(byte(n))
		b.Write(src[i : i+n])
		dn += n
		sn += n + 1
		i += n
	}

	pad := (4 - (b.Len() % 4)) % 4
	for i := 0; i < pad; i++ {
		b.WriteByte(0)
	}

	p := b.Bytes()
	binary.LittleEndian.PutUint32(p[0xc:], uint32(sn))
	binary.LittleEndian.PutUint32(p[0x10:], uint32(dn))

	return p, nil
}

func (c *Slide) copy(dst, src []byte, dp, sp, n uint32) error {
	if uint64(dp) >= uint64(len(dst)) || uint64(sp) >= uint64(len(src)) || uint64(len(src))-uint64(sp) < uint64(n) {
		return io.ErrUnexpectedEOF
	}
	copy(dst[dp:], src[sp:sp+n])
	return nil
}

func (c *Slide) Decode2(src []byte, limit uint32) ([]byte, error) {
	enc := binary.LittleEndian.Uint32(src[0x8:])
	if enc != 0x4 {
		return nil, fmt.Errorf("decode: invalid encoding")
	}

	dn := binary.LittleEndian.Uint32(src[0x10:])
	dn = (dn + 3) &^ 3
	if limit < dn {
		dn = limit
	}
	dst := make([]byte, dn)

	sn := binary.LittleEndian.Uint32(src[0xc:])
	dp := uint32(0)
	sp := uint32(0x14)

loop:
	for {
		if uint64(sp) >= uint64(len(src)) || sp >= sn || uint64(dp) >= uint64(len(dst)) {
			break
		}

		op := uint32(src[sp])
		sp++

		var (
			tp  uint32
			n   uint32
			err error
		)
		switch v := int32(op); {
		case v&0x80 == 0:
			n = uint32(v)
			err = c.copy(dst, src, dp, sp, n)

			dp += n
			sp += n

		case v&0x40 == 0:
			tp = dp - (op & 0xf) - 1
			n = (op&0x30)>>4 + 1

			err = c.copy(dst, dst, dp, tp, n)
			dp += n

		case v&0x20 != 0:
			if sp+2 >= uint32(len(src)) {
				break loop
			}

			hi := (op & 0x1f) << 4
			lo := uint32((src[sp] & 0xf0) >> 4)
			n = (lo | hi) + 3

			hi = (uint32(src[sp]) << 8) & 0xf00
			lo = uint32(src[sp+1])
			tp = dp - (hi + lo) - 1

			err = c.copy(dst, dst, dp, tp, n)
			dp += n
			sp += 2

		default:
			if sp >= uint32(len(src)) {
				break loop
			}

			n = (op & 0x1f) + 2
			tp = dp - uint32(src[sp]) - 1

			err = c.copy(dst, dst, dp, tp, n)
			dp += n
			sp++
		}

		if err != nil {
			return dst, fmt.Errorf("decode: malformed data: op %x dp %d sp %d tp %d n %d dst %d src %d",
				op, dp, sp, tp, n, len(dst), len(src))
		}
	}
	return dst, nil
}
