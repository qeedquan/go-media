package dsarcfl

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"strings"
)

type Archive struct {
	Header
	Entries []Entry
	Version int
}

type Header struct {
	Magic      [8]byte
	NumEntries uint32
	Zero       uint32
}

type Entry struct {
	Name  string
	Namez []byte
	Size  uint32
	Off   uint32
	Data  []byte
}

const sig = "DSARC FL"

func Decode(r io.ReaderAt, version int) (*Archive, error) {
	var (
		hdr     Header
		entries []Entry
	)

	sr := io.NewSectionReader(r, 0, math.MaxInt64)
	err := binary.Read(sr, binary.LittleEndian, &hdr)
	if err != nil {
		return nil, err
	}

	if string(hdr.Magic[:]) != sig {
		return nil, fmt.Errorf("invalid header magic: %q", hdr.Magic)
	}

	for i := uint32(0); i < hdr.NumEntries; i++ {
		var (
			name []byte
			val  []uint32
		)

		switch {
		case version > 1:
			name = make([]byte, 116)
			val = make([]uint32, 3)
		default:
			name = make([]byte, 40)
			val = make([]uint32, 2)
		}

		err := binary.Read(sr, binary.LittleEndian, name[:])
		xerr := binary.Read(sr, binary.LittleEndian, val[:])
		if err != nil || xerr != nil {
			return nil, fmt.Errorf("failed to read entry %v header", i)
		}

		entries = append(entries, Entry{
			Name:  strings.TrimRight(string(name[:]), "\x00"),
			Namez: name[:],
			Size:  val[0],
			Off:   val[1],
			Data:  make([]byte, val[0]),
		})
	}

	for i := range entries {
		e := &entries[i]
		sr.Seek(int64(e.Off), io.SeekStart)
		_, err := io.ReadAtLeast(sr, e.Data, len(e.Data))
		if err != nil {
			return nil, fmt.Errorf("failed to read entry %v data: %v", i, err)
		}
	}

	return &Archive{
		Header:  hdr,
		Entries: entries,
		Version: version,
	}, nil
}
