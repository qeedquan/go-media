package debug

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
)

type StringType int

const (
	UTF8 StringType = iota
)

type String struct {
	Type   StringType
	Name   string
	Value  string
	Offset uint64
}

func (s StringType) String() string {
	switch s {
	case UTF8:
		return "utf8"
	default:
		return "unknown"
	}
}

func (s String) String() string {
	return fmt.Sprintf("(%s %s+%#0x:%q)", StringType(s.Type), s.Name, s.Offset, s.Value)
}

func FindStringsBuffer(buf []byte, minWidth, maxWidth int, ascii bool) []String {
	return FindStringsReader(bytes.NewReader(buf), minWidth, maxWidth, ascii)
}

func FindStringsReader(r io.Reader, minWidth, maxWidth int, ascii bool) []String {
	var tab []String

	in := bufio.NewReader(r)
	str := make([]rune, 0, maxWidth)
	pos := int64(0)
	off := uint64(0)
	print := func() {
		if len(str) >= minWidth {
			tab = append(tab, String{
				Type:   UTF8,
				Value:  string(str),
				Offset: off,
			})
		}
		str = str[:0]
	}

	for {
		var (
			r   rune
			wid int
			err error
		)

		for ; ; pos += int64(wid) {
			r, wid, err = in.ReadRune()
			if err != nil {
				return tab
			}
			if !strconv.IsPrint(r) || ascii && r >= 0xFF {
				print()
				continue
			}
			if len(str) >= maxWidth {
				print()
			}
			str = append(str, r)
			off += uint64(wid)
		}
	}
	return tab
}
