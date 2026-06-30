package extekhex

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type File struct {
	Records []Record
}

type Record struct {
	Addr uint64
	Data []byte
}

func Decode(r io.Reader) (*File, error) {
	f := &File{}
	s := bufio.NewScanner(r)
	l := 1

	errf := func(format string, args ...interface{}) error {
		pfx := fmt.Sprintf("extekhex: #%d: ", l)
		return fmt.Errorf(pfx+format, args...)
	}

loop:
	for ; s.Scan(); l++ {
		t := strings.TrimSpace(s.Text())
		r := strings.NewReader(t)

		var (
			hid      byte
			blen     int
			btyp     int
			checksum uint8
			addrsz   int
			addr     uint64
			data     []byte
			val      uint8
		)

		_, err := fmt.Fscanf(r, "%c%02x%1x%2x%1x", &hid, &blen, &btyp, &checksum, &addrsz)
		if err != nil {
			return nil, errf("failed to read line header info")
		}
		if hid != '%' {
			continue
		}

		format := fmt.Sprintf("%%%dx", addrsz)
		_, err = fmt.Fscanf(r, format, &addr)
		if err != nil {
			return nil, errf("failed to read address")
		}

		for {
			n, _ := fmt.Fscanf(r, "%2x", &val)
			if n != 1 {
				break
			}
			data = append(data, val)
		}
		f.Records = append(f.Records, Record{addr, data})

		switch btyp {
		case 6:
		case 8:
			break loop
		default:
			return nil, errf("invalid block type %q", btyp)
		}
	}
	return f, nil
}

func Encode(w io.Writer, f *File) error {
	b := bufio.NewWriter(w)
	for i, r := range f.Records {
		blen := 2 + 1 + 2 + 1 + 8 + len(r.Data)
		btyp := 8
		if i+1 < len(f.Records) {
			btyp = 6
		}

		str := fmt.Sprintf("%02X%1X8%8X", blen, btyp, r.Addr)
		for i := range r.Data {
			str += fmt.Sprintf("%2X", r.Data[i])
		}
		checksum := Checksum([]byte(str))

		fmt.Fprintf(b, "%02X%1X%02X8%8X", blen, btyp, checksum, r.Addr)
		for i := range r.Data {
			fmt.Fprintf(b, "%02X", r.Data[i])
		}
	}

	return b.Flush()
}

func Checksum(b []byte) uint8 {
	var s uint8
	for i := range b {
		s += b[i]
	}
	return s
}

func (f *File) Binary() []byte {
	if len(f.Records) == 0 {
		return nil
	}
	m := uint64(0)
	for _, r := range f.Records {
		l := r.Addr + uint64(len(r.Data))
		if m < l {
			m = l
		}
	}
	b := make([]byte, m)
	for _, r := range f.Records {
		copy(b[r.Addr:], r.Data)
	}
	return b
}