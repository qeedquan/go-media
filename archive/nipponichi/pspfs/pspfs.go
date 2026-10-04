package pspfs

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"strings"
)

type Archive struct {
	Header
	Entries []Entry
}

type Header struct {
	Magic      [8]byte
	NumEntries uint32
	Unknown    uint32
}

type Entry struct {
	Name    string
	DecSize uint32
	Size    uint32
	Off     uint32
	Data    []byte
}

const sig = "PSPFS_V1"

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

	if string(hdr.Magic[:]) != sig {
		return nil, fmt.Errorf("invalid header magic: %q", hdr.Magic)
	}

	for i := uint32(0); i < hdr.NumEntries; i++ {
		var dat struct {
			Name [40]byte
			Val  [3]uint32
		}

		err = binary.Read(sr, binary.LittleEndian, &dat)
		if err != nil {
			return nil, fmt.Errorf("failed to read entry %v header: %v", i, err)
		}

		name := strings.TrimRight(string(dat.Name[:]), "\x00")
		name = strings.TrimRight(name, " ")
		name = strings.ToUpper(name)
		entries = append(entries, Entry{
			Name:    name,
			DecSize: dat.Val[0],
			Size:    dat.Val[1],
			Off:     dat.Val[2],
			Data:    make([]byte, dat.Val[1]),
		})
	}

	for i := range entries {
		entry := &entries[i]

		sr.Seek(int64(entry.Off), io.SeekStart)
		br := bufio.NewReader(sr)
		_, err = io.ReadAtLeast(br, entry.Data, int(entry.Size))
		if err != nil {
			return nil, fmt.Errorf("failed to read entry %v data: %v", i, err)
		}
	}

	return &Archive{
		Header:  hdr,
		Entries: entries,
	}, nil
}
