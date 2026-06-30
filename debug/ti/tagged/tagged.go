package tagged

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type File struct {
	Ident   string
	Header  *Header
	Records []Record
}

type Header struct {
	Size uint16
	Name string
}

type Record struct {
	Addr uint64
	Data []byte
}

func Decode(r io.Reader) (*File, error) {
	f := &File{}
	s := bufio.NewScanner(r)
	l := 1

	var (
		addr     uint64
		tag      int
		hdr      Header
		checksum uint16
		data     []byte
		buf      [2]byte
	)

	errf := func(format string, args ...interface{}) error {
		pfx := fmt.Sprintf("tagged: #%d: ", l)
		return fmt.Errorf(pfx+format, args...)
	}

	flushdata := func() {
		if addr != 0 {
			f.Records = append(f.Records, Record{addr, data})
			addr += uint64(len(data))
			data = nil
		}
	}

scan:
	for ; s.Scan(); l++ {
		t := strings.TrimSpace(s.Text())
		if len(t) == 0 {
			continue
		}
		r := strings.NewReader(t)

		for {
			_, err := fmt.Fscanf(r, "%c", &tag)
			if err == io.EOF {
				continue scan
			}

			switch tag {
			case '*':
				_, err := fmt.Fscanf(r, "%02x", &buf[0])
				if err != nil {
					return nil, errf("failed to read data byte: %v", err)
				}
				data = append(data, buf[0])

			case ':':
				break scan

			case '0':
				_, err := fmt.Fscanf(r, "%04x%8s", &hdr.Size, &hdr.Name)
				if err != nil {
					return nil, errf("failed to read file header: %v", err)
				}
				f.Header = &hdr

			case '7':
				pos, _ := r.Seek(0, io.SeekCurrent)
				calcsum := Checksum([]byte(t[:pos]))
				_, err := fmt.Fscanf(r, "%04x", &checksum)
				if err != nil {
					return nil, errf("failed to read checksum: %v", err)
				}
				if calcsum != checksum {
					return nil, errf("checksum mismatch: %#x != %#x", calcsum, checksum)
				}

			case '8':
				_, err := fmt.Fscanf(r, "%04x", &checksum)
				if err != nil {
					return nil, errf("failed to read dummy checksum: %v", err)
				}

			case '9':
				flushdata()
				_, err := fmt.Fscanf(r, "%04x", &addr)
				if err != nil {
					return nil, errf("failed to read address: %v", err)
				}

			case 'B':
				_, err := fmt.Fscanf(r, "%02x%02x", &buf[0], &buf[1])
				if err != nil {
					return nil, errf("failed to read binary data word: %v", err)
				}
				data = append(data, buf[0], buf[1])

			case 'F':
				flushdata()

			case 'K':
				var length int
				_, err := fmt.Fscanf(r, "%04x", &length)
				if err != nil {
					return nil, errf("invalid program identifier length %#x", length)
				}
				if length == 0 {
					length = 5 + 8
				}

				if length -= 5; length > 0 {
					format := fmt.Sprintf("%%%ds", length)
					_, err = fmt.Fscanf(r, format, &f.Ident)
					if err != nil {
						return nil, errf("failed to read program identifier: %v", err)
					}
				}

			case ' ':
				continue scan

			default:
				return nil, errf("unknown tag %q", tag)
			}
		}
	}

	return f, nil
}

func Encode(w io.Writer, f *File) error {
	b := bufio.NewWriter(w)
	c := &checksumWriter{}
	c.Reset()

	fmt.Fprintf(c, "K%04X", 5+len(f.Ident))
	for _, r := range f.Records {
		fmt.Fprintf(c, "9%04X", r.Addr)
		for i := 0; i < len(r.Data); {
			if i+2 < len(r.Data) {
				fmt.Fprintf(c, "B%02X%02X", r.Data[i], r.Data[i+1])
				i += 2
			} else {
				fmt.Fprintf(c, "*%02X", r.Data[i])
				i += 1
			}
		}
		fmt.Fprintf(c, "7")
		b.Write(c.Data())
		fmt.Fprintf(b, "%04X", c.Checksum())
		fmt.Fprintf(b, "F\n")
		c.Reset()
	}
	fmt.Fprintf(b, ":\n")
	return b.Flush()
}

func Checksum(b []byte) uint16 {
	var s uint16
	for i := range b {
		s += uint16(b[i])
	}
	return ^s + 1
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

type checksumWriter struct {
	data []byte
	sum  uint16
}

func (c *checksumWriter) Reset() {
	c.data = []byte{}
	c.sum = 0
}

func (c *checksumWriter) Data() []byte {
	return c.data
}

func (c *checksumWriter) Write(b []byte) (n int, err error) {
	for i := range b {
		c.data = append(c.data, b[i])
		c.sum += uint16(b[i])
	}
	return len(b), nil
}

func (c *checksumWriter) Checksum() uint16 {
	return ^c.sum + 1
}
