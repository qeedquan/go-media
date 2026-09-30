package imy

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

type Archive struct {
	Offs    []uint32
	Entries []Entry
}

type Header struct {
	Magic     [4]byte
	CompStart uint16
	Unknown1  uint16
	StreamOff uint16
	CompType  uint8
	Unknown2  uint8
	ChunkSize uint16
	Unknown4  uint16
	Pad       [16]byte
}

type Entry struct {
	Header
	Data []byte
}

var MAGIC = [4]byte{'I', 'M', 'Y', 0}

func Decode(r io.ReaderAt) (*Archive, error) {
	var (
		sig [4]byte
		dat [4]byte
	)
	sr := io.NewSectionReader(r, 0, math.MaxInt64)
	err := binary.Read(sr, binary.LittleEndian, &sig)
	if err != nil {
		return nil, err
	}

	ar := &Archive{
		Offs:    []uint32{0},
		Entries: make([]Entry, 1),
	}
	if sig != MAGIC {
		nentries := binary.LittleEndian.Uint32(sig[:])
		err = binary.Read(sr, binary.LittleEndian, &dat)
		if err != nil {
			return nil, err
		}

		ar = &Archive{
			Offs:    make([]uint32, nentries),
			Entries: make([]Entry, nentries),
		}
		for i := range ar.Offs {
			err = binary.Read(sr, binary.LittleEndian, &ar.Offs[i])
			if err != nil {
				return nil, err
			}
		}
	}

	for i := range ar.Offs {
		sr.Seek(int64(ar.Offs[i]), io.SeekStart)
		err = decodeIMY(sr, &ar.Entries[i])
		if err != nil {
			return nil, err
		}
	}

	return ar, nil
}

func decodeIMY(sr *io.SectionReader, entry *Entry) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("corrupted file")
		}
	}()

	err = binary.Read(sr, binary.LittleEndian, &entry.Header)
	if err != nil {
		return
	}

	if entry.Magic != MAGIC {
		return fmt.Errorf("invalid header signature")
	}

	switch typ := (entry.CompType >> 4) & 0xf; typ {
	case 1:
		return decompress1(sr, entry)
	default:
		return fmt.Errorf("unsupported compression scheme %d", typ)
	}
}

func decompress1(sr *io.SectionReader, entry *Entry) error {
	var (
		numops  uint16
		startop int64
		endop   int64
		op      uint8
		val     uint16
	)
	binary.Read(sr, binary.LittleEndian, &numops)
	decsize := int64(entry.StreamOff) * int64(entry.ChunkSize)
	startop, _ = sr.Seek(0, io.SeekCurrent)
	endop = startop + int64(numops)
	dataoff := endop

	lut := [4]int64{}
	lut[0] = 2
	lut[1] = int64(entry.StreamOff)
	lut[2] = lut[1] + 2
	lut[3] = lut[1] - 2

	out := make([]byte, decsize*2)
	pos := int64(0)
	for pos < decsize && startop < endop {
		sr.Seek(startop, io.SeekStart)
		binary.Read(sr, binary.LittleEndian, &op)
		startop += 1
		if op&0xF0 != 0 {
			if op&0x80 != 0 && op&0x40 != 0 {
				index := (op & 0x30) >> 4
				for range (op & 0x0F) + 1 {
					val = binary.LittleEndian.Uint16(out[pos-lut[index]:])
					binary.LittleEndian.PutUint16(out[pos:], val)
					pos += 2
				}
			} else {
				lookback_bytes := (int64(op)-16)*2 + 2
				sr.Seek(dataoff-lookback_bytes, io.SeekStart)
				binary.Read(sr, binary.LittleEndian, &val)
				binary.LittleEndian.PutUint16(out[pos:], val)
				pos += 2
			}
		} else {
			copy_bytes := (int64(op) + 1) * 2
			sr.Seek(dataoff, io.SeekStart)
			sr.Read(out[pos : pos+copy_bytes])
			pos += copy_bytes
			dataoff += copy_bytes
		}
	}

	entry.Data = out[:decsize]
	return nil
}
