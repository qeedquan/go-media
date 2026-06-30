package car

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"strings"
)

var (
	magic = [3]byte{'r', 'a', 'c'}
)

const (
	TAG_UNKNOWN  = 0x0
	TAG_CONTENTS = 0x1
	TAG_DATA     = 0x2
	TAG_EOF      = 0xFFFFFFFF
)

type Archive struct {
	Entries []Entry
}

type Entry struct {
	Type   uint32
	Offset uint32
	Name   string
	Data   []byte
}

func Decode(r io.ReaderAt) (ar *Archive, err error) {
	sr := io.NewSectionReader(r, 0, math.MaxInt32)
	err = readheader(sr)
	if err != nil {
		return
	}

	elems, err := readcontents(sr)
	if err != nil {
		return
	}

	var (
		entries []Entry
		entry   Entry
	)
	for i := uint32(0); i < elems; i++ {
		entry, err = readentry(sr)
		if err != nil {
			return
		}
		entries = append(entries, entry)
	}

	for i := range entries {
		e := &entries[i]
		e.Data, err = readdata(r, int64(e.Offset))
		if err != nil {
			return
		}
	}

	return &Archive{Entries: entries}, nil
}

func readheader(r io.Reader) error {
	var sig [4]byte
	err := binary.Read(r, binary.LittleEndian, sig[:])
	if err != nil {
		return err
	}

	if bytes.Compare(sig[:3], magic[:]) != 0 {
		return fmt.Errorf("invalid magic signature: %q", sig[:3])
	}

	if sig[3] != 1 {
		return fmt.Errorf("unsupported version: %v", sig[3])
	}

	return nil
}

func readcontents(r io.Reader) (elems uint32, err error) {
	var (
		val [3]uint32
	)
	err = binary.Read(r, binary.LittleEndian, &val)
	if err != nil {
		return
	}

	if val[0] != TAG_CONTENTS {
		err = fmt.Errorf("expected content tag, got %v", val[0])
		return
	}

	elems = val[2]
	return
}

func readentry(r io.Reader) (entry Entry, err error) {
	var (
		val  [2]uint32
		name string
	)

	err = binary.Read(r, binary.LittleEndian, &val)
	if err != nil {
		return
	}

	name, err = readstr(r)
	if err != nil {
		return
	}

	entry = Entry{
		Type:   val[0],
		Offset: val[1],
		Name:   name,
	}
	return
}

func readstr(r io.Reader) (str string, err error) {
	var (
		buf    []byte
		length uint32
	)

	err = binary.Read(r, binary.LittleEndian, &length)
	if err != nil {
		return
	}

	length = (length + 4) &^ 3
	buf = make([]byte, length)
	err = binary.Read(r, binary.LittleEndian, buf)
	if err != nil {
		return
	}

	str = strings.Trim(string(buf), "\x00")
	return
}

func readdata(r io.ReaderAt, off int64) (data []byte, err error) {
	var val [3]uint32

	sr := io.NewSectionReader(r, off, math.MaxInt32)
	err = binary.Read(sr, binary.LittleEndian, &val)
	if err != nil {
		return
	}

	if val[0] != TAG_DATA {
		err = fmt.Errorf("expected data tag, but got %v", val[0])
		return
	}

	length := val[2]
	data = make([]byte, length)
	err = binary.Read(sr, binary.LittleEndian, data)

	return
}
