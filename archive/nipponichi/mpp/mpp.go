package mpp

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

type Header struct {
	NumTextures uint16
	HasNormals  uint16
	NumGeoms    uint16
	Unk         uint16
	TotalSize   uint32
}

type Entry struct {
	Name string
	Off  uint32
	Data []byte
}

type Archive struct {
	Header
	Entries []Entry
}

func Decode(r io.ReaderAt) (*Archive, error) {
	var (
		hdr     Header
		entries []Entry
	)

	sr := io.NewSectionReader(r, 0, math.MaxInt64)
	err := binary.Read(sr, binary.LittleEndian, &hdr)
	if err != nil {
		return nil, err
	}

	numNormals := 0
	if hdr.HasNormals != 0 {
		numNormals = int(hdr.NumTextures)
	}
	offs := make([]uint32, int(hdr.NumTextures)+numNormals+int(hdr.NumGeoms))
	err = binary.Read(sr, binary.LittleEndian, offs[:])
	if err != nil {
		return nil, err
	}

	var n0, n1, n2 int
	for i := range offs {
		size := int(hdr.TotalSize) - int(offs[i])
		if i+1 < len(offs) {
			size = int(offs[i+1]) - int(offs[i])
		}
		if size < 0 {
			return nil, fmt.Errorf("invalid size at entry %d", i)
		}

		var name string
		switch {
		case i < int(hdr.NumTextures):
			name = fmt.Sprintf("TEXTURE%d.TX2", n0)
			n0 += 1
		case hdr.HasNormals != 0 && i < int(hdr.NumTextures)+numNormals:
			name = fmt.Sprintf("NORMAL%d.TX2", n1)
			n1 += 1
		default:
			name = fmt.Sprintf("GEOMETRY%d.GEO", n2)
			n2 += 1
		}
		entry := Entry{
			Name: name,
			Off:  offs[i],
			Data: make([]byte, size),
		}

		sr.Seek(int64(offs[i]), io.SeekStart)
		_, err = io.ReadAtLeast(sr, entry.Data, len(entry.Data))
		if err != nil {
			return nil, err
		}

		entries = append(entries, entry)
	}

	return &Archive{
		Header:  hdr,
		Entries: entries,
	}, nil
}
