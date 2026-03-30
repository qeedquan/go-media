package ico

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/png"
	"io"
	"math"

	"golang.org/x/image/bmp"
)

type File struct {
	Image []image.Image
}

type header struct {
	_       uint16
	Type    uint16
	Entries uint16
}

type dirent struct {
	Width   uint8
	Height  uint8
	Palette uint8
	_       uint8
	Planes  uint16
	Bpp     uint16
	Size    uint32
	Off     uint32
}

func Encode(w io.Writer, f *File) error {
	if len(f.Image) >= math.MaxUint16 {
		return fmt.Errorf("ico: format cannot support %d images", len(f.Image))
	}

	b := bufio.NewWriter(w)
	h := header{
		Type:    1,
		Entries: uint16(len(f.Image)),
	}
	binary.Write(b, binary.LittleEndian, &h)

	const (
		headerLen = 6
		direntLen = 16
	)

	off := headerLen + direntLen*uint32(h.Entries)
	for i, m := range f.Image {
		p := new(bytes.Buffer)
		err := png.Encode(p, m)
		if err != nil {
			return err
		}

		if p.Len() >= math.MaxInt32 {
			return fmt.Errorf("ico: image %d is too big", i)
		}

		r := m.Bounds()
		if r.Dx() >= math.MaxUint8 || r.Dy() >= math.MaxUint8 {
			return fmt.Errorf("ico: image %d with dimension %dx%d is too big", i, r.Dx(), r.Dy())
		}

		d := dirent{
			Width:  uint8(r.Dx()),
			Height: uint8(r.Dy()),
			Bpp:    32,
			Size:   uint32(p.Len()),
			Off:    off,
		}
		binary.Write(b, binary.LittleEndian, &d)

		if int64(off)+int64(p.Len()) >= math.MaxUint32 {
			return fmt.Errorf("ico: too many images")
		}
		off += uint32(p.Len())
	}

	return b.Flush()
}

func Decode(r io.Reader) (*File, error) {
	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	b := bytes.NewReader(buf)

	var h header
	err = binary.Read(b, binary.LittleEndian, &h)
	if err != nil {
		return nil, err
	}

	d := make([]dirent, h.Entries)
	for i := range d {
		err = binary.Read(b, binary.LittleEndian, &d[i])
		if err != nil {
			return nil, err
		}
	}

	f := &File{}
	for i, d := range d {
		_, err = b.Seek(int64(d.Off), io.SeekStart)
		if err != nil {
			return nil, err
		}

		m, _, err := image.Decode(b)
		if err != nil {
			m, err = decodeBMP(i, &d, buf)
		}
		if err != nil {
			return nil, err
		}
		f.Image = append(f.Image, m)
	}

	return f, nil
}

func decodeBMP(n int, d *dirent, b []byte) (image.Image, error) {
	const (
		fileHeaderLen = 14
		infoHeaderLen = 40
	)

	if int64(d.Off) >= int64(len(b)) || int64(len(b))-int64(d.Off) < int64(d.Size) || len(b) <= fileHeaderLen+infoHeaderLen {
		return nil, fmt.Errorf("invalid size for bitmap %d with offset %d and size %d",
			n, d.Off, d.Size)
	}

	sz := d.Size
	b = b[d.Off : d.Off+d.Size]
	bpp := readUint16(b[28-fileHeaderLen:])
	off := fileHeaderLen + infoHeaderLen
	switch bpp {
	case 8:
		off += 256 * 4
	}

	b = append([]uint8{
		'B', 'M',
		uint8(sz), uint8(sz >> 8), uint8(sz >> 16), uint8(sz >> 32),
		0, 0,
		0, 0,
		uint8(off), uint8(off >> 8), uint8(off >> 16), uint8(off >> 32),
	}, b...)

	return bmp.Decode(bytes.NewReader(b))
}

func readUint16(b []byte) uint16 {
	return uint16(b[0]) | uint16(b[1])<<8
}
