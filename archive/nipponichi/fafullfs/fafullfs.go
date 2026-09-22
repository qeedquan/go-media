package fafullfs

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

type Archive struct {
	Header
	Entries []Entry
}

type Header struct {
	Magic      [8]byte
	NumEntries uint64
	NameOff    uint64
	NameSize   uint64
	EntryOff   uint64
}

type Entry struct {
	Checksum uint64
	NameOff  uint64
	Zero     uint64
	Size     uint64
	Off      uint64
	Time     uint64
	Name     string
	Data     []byte
}

const sig = "FAFULLFS"

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

	sr.Seek(int64(hdr.EntryOff), io.SeekStart)
	for i := uint64(0); i < hdr.NumEntries; i++ {
		var (
			val [6]uint64
		)

		err = binary.Read(sr, binary.LittleEndian, val[:])
		if err != nil {
			return nil, fmt.Errorf("failed to read entry %v header: %v", i, err)
		}

		entries = append(entries, Entry{
			Checksum: val[0],
			NameOff:  val[1],
			Zero:     val[2],
			Size:     val[3],
			Off:      val[4],
			Time:     val[5],
			Data:     make([]byte, val[3]),
		})
	}

	for i := range entries {
		entry := &entries[i]

		sr.Seek(int64(hdr.NameOff+entry.NameOff), io.SeekStart)
		entry.Name, err = readstrz(sr)

		if err != nil {
			return nil, fmt.Errorf("failed to read entry %v name: %v", i, err)
		}

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

func readstrz(r io.Reader) (string, error) {
	str := ""
	br := bufio.NewReader(r)
	for {
		c, err := br.ReadByte()
		if c == 0 {
			break
		}

		if err != nil {
			return str, err
		}
		str += string(c)
	}
	return str, nil
}
