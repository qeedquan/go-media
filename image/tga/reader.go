package tga

import (
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"image/color"
	"io"
)

type header struct {
	SizeID   uint8
	ColorMap uint8
	Type     uint8

	ColorMapStart   uint16
	ColorMapEntries uint16
	ColorMapBpp     uint8

	Xstart uint16
	Ystart uint16
	Width  uint16
	Height uint16

	Bpp  uint8
	Desc uint8
}

type decoder struct {
	header
	r        io.Reader
	img      *image.RGBA
	colormap []byte
	pix      []byte
}

func Decode(r io.Reader) (image.Image, error) {
	d := &decoder{r: r}

	if err := d.checkHeader(); err != nil {
		return nil, err
	}

	if d.SizeID != 0 {
		_, err := d.readLength(int(d.SizeID))
		if err != nil {
			return nil, err
		}
	}

	if d.ColorMap&1 == 1 {
		var err error

		length := int(d.ColorMapEntries) * int(d.ColorMapBpp+1) / 8
		d.colormap, err = d.readLength(length)
		if err != nil {
			return nil, err
		}
	}

	dim := image.Rect(0, 0, int(d.Width), int(d.Height))
	d.img = image.NewRGBA(dim)

	var err error
	d.pix, err = io.ReadAll(d.r)
	if err != nil {
		return nil, err
	}

	if d.Type&^byte(3) != 0 {
		d.pix, err = d.rleUncompress(d.pix)
		if err != nil {
			return nil, err
		}
	}

	d.decode()

	return d.img, nil
}

func DecodeConfig(r io.Reader) (image.Config, error) {
	d := &decoder{r: r}

	if err := d.checkHeader(); err != nil {
		return image.Config{}, err
	}

	return image.Config{
		ColorModel: color.RGBAModel,
		Width:      int(d.Width),
		Height:     int(d.Height),
	}, nil
}

func (d *decoder) readLength(length int) ([]byte, error) {
	b := make([]byte, length)
	n, err := d.r.Read(b)
	if err != nil {
		return nil, err
	}
	if n != length {
		return nil, io.ErrUnexpectedEOF
	}

	return b, err
}

func (d *decoder) checkHeader() error {
	return binary.Read(d.r, binary.LittleEndian, &d.header)
}

func (d *decoder) rleUncompress(p []byte) ([]byte, error) {
	var b []byte

	if len(p) == 0 {
		return p, nil
	}

	bpp := int(d.Bpp) / 8
	size := int(d.Width) * int(d.Height) * bpp
	for i := 0; len(b) < size; {
		op := int(p[i])
		if i++; i+bpp >= len(p) {
			return nil, io.ErrUnexpectedEOF
		}
		pix := p[i : i+bpp]

		if op&0x80 != 0 {
			j := op - 0x7f
			for k := 0; k < j; k++ {
				b = append(b, pix...)
			}
			i += bpp
		} else {
			j := op + 1
			for k := 0; k < j; k++ {
				if i+bpp >= len(p) {
					return nil, io.ErrUnexpectedEOF
				}
				b = append(b, pix...)
				pix = p[i : i+bpp]
				i += bpp
			}
		}
	}
	return b, nil
}

func (d *decoder) decode() error {
	size := int(d.Width) * int(d.Height) * int(d.Bpp) / 8
	if len(d.pix) > size {
		return errors.New("invalid size")
	}

	y, dy, h := int(d.Height-1), -1, -1
	if d.Desc&0x20 != 0 {
		y, dy, h = 0, 1, int(d.Height)
	}

	x, dx, w := 0, 1, int(d.Width)
	if d.Desc&0x10 != 0 {
		x, dx, w = int(d.Width-1), -1, -1
	}
	rx := x

	var dec func(*decoder, []byte) color.RGBA
	var i, inc int

	switch d.Bpp {
	case 8:
		dec = decodeGray
	case 16:
		dec = decode16
	case 24:
		dec = decode24
	case 32:
		dec = decode32
	default:
		return fmt.Errorf("unsupported bpp size: %d", d.Bpp)
	}

	inc = int(d.Bpp / 8)
	for y != h {
		x = rx
		for x != w {
			d.img.SetRGBA(x, y, dec(d, d.pix[i:]))
			x += dx
			i += inc
		}
		y += dy
	}
	return nil
}

func decodeIndexed(d *decoder, p []byte) color.RGBA {
	return color.RGBA{}
}

func decodeGray(_ *decoder, p []byte) color.RGBA {
	return color.RGBA{p[0], p[0], p[0], 255}
}

func decode16(_ *decoder, p []byte) color.RGBA {
	r := (p[1] & 0x7C) << 1
	g := ((p[1] & 0x3) << 6) | ((p[0] & 0xE0) >> 2)
	b := (p[0] & 0x1F) << 3
	a := uint8((uint16(p[1]&0x80) << 1) - 1)
	return color.RGBA{r, g, b, a}
}

func decode24(_ *decoder, p []byte) color.RGBA {
	return color.RGBA{p[2], p[1], p[0], 255}
}

func decode32(_ *decoder, p []byte) color.RGBA {
	return color.RGBA{p[2], p[1], p[0], p[3]}
}
