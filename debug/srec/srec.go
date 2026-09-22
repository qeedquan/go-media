package srec

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
)

type File struct {
	Version int
	Start   uint64
	Records []Record
}

type Record struct {
	Addr uint64
	Data []byte
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

func Decode(r io.Reader) (*File, error) {
	f := &File{Version: 3}
	s := bufio.NewScanner(r)
	l := 1

	errf := func(format string, args ...interface{}) error {
		pfx := fmt.Sprintf("srec: #%d: ", l)
		return fmt.Errorf(pfx+format, args...)
	}

	for ; s.Scan(); l++ {
		t := strings.TrimSpace(s.Text())
		if !strings.HasPrefix(t, "S") {
			continue
		}

		if len(t) < 4 {
			return nil, errf("line too short")
		}

		b, err := hex.DecodeString(t[2:])
		if err != nil {
			return nil, errf("%v", err)
		}

		if int(b[0]) != len(b)-1 {
			return nil, errf("byte count mismatch, expected %d, got %d", b[0], len(b))
		}

		u := Checksum(b[:len(b)-1])
		if u != b[len(b)-1] {
			return nil, errf("checksum mismatch, expected %#x, got %#x", u, b[len(b)-1])
		}
		b = b[1 : len(b)-1]

		var (
			rec Record
			eof bool
		)
		switch t[1] {
		case '0', '4', '5', '6':
		case '1', '2', '3':
			f.Version = int(t[1]) - '0'
			rec.Addr, b, err = rtv(t[1], b)
			rec.Data = b
			f.Records = append(f.Records, rec)
		case '7', '8', '9':
			f.Version = '9' - int(t[1]) + 1
			f.Start, b, err = rtv(t[1], b)
			eof = true
		default:
			err = fmt.Errorf("unsupported record type %q", t[1])
		}

		if err != nil {
			return nil, errf("%v", err)
		}

		if eof {
			break
		}
	}

	return f, nil
}

func Encode(w io.Writer, f *File) error {
	bw := bufio.NewWriter(w)
	wrc(bw, 0, []byte("\x00\x00HDR"))
	for _, r := range f.Records {
		wrc(bw, f.Version, append(wtv(f.Version, r.Addr), r.Data...))
	}
	wrc(bw, 10-f.Version, wtv(f.Version, f.Start))

	err := bw.Flush()
	if err != nil {
		return fmt.Errorf("srec: %v", err)
	}
	return nil
}

func Checksum(b []byte) byte {
	var s byte
	for i := range b {
		s += b[i]
	}
	return ^s
}

func rtv(t byte, b []byte) (uint64, []byte, error) {
	if t > '3' {
		t = '9' - t + '1'
	}

	n := int(t) - '0' + 1
	if len(b) < n {
		return 0, b, io.ErrShortBuffer
	}

	v := uint64(0)
	for i := 0; i < n; i++ {
		v |= uint64(b[i]) << (8 * (n - i - 1))
	}
	return v, b[n:], nil
}

func wtv(t int, v uint64) []byte {
	switch t {
	case 1:
		return []byte{byte(v >> 8), byte(v)}
	case 2:
		return []byte{byte(v >> 16), byte(v >> 8), byte(v)}
	case 3:
		return []byte{byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)}
	}
	return nil
}

func wrc(bw *bufio.Writer, t int, b []byte) {
	b = append([]byte{byte(len(b) + 1)}, b...)
	u := Checksum(b)
	h := hex.EncodeToString(b)
	h = strings.ToUpper(h)
	fmt.Fprintf(bw, "S%d%s%02X\n", t, h, u)
}
