package pnm

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"image/color"
	"io"
)

var (
	ErrFormat = errors.New("pnm: unsupported format")
)

func Decode(r io.Reader) (image.Image, error) {
	d := decoder{r: r}
	err := d.decodeHeader()
	if err != nil {
		return nil, err
	}

	b := image.Rect(0, 0, d.w, d.h)
	m := image.NewRGBA(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			m.Set(x, y, d.readColor())
		}
	}

	if d.err != nil {
		return nil, fmt.Errorf("pnm: %v", d.err)
	}

	return m, nil
}

func DecodeConfig(r io.Reader) (image.Config, error) {
	d := decoder{r: r}
	err := d.decodeHeader()
	if err != nil {
		return image.Config{}, err
	}

	return image.Config{
		ColorModel: color.RGBAModel,
		Width:      d.w,
		Height:     d.h,
	}, nil
}

type header struct {
	format int
	maxval int
	w, h   int
}

type decoder struct {
	r    io.Reader
	b    *bufio.Reader
	err  error
	bits uint8
	bw   uint8
	header
}

func (d *decoder) decodeHeader() error {
	d.b = bufio.NewReader(d.r)

	var sig [2]byte
	fmt.Fscanf(d.b, "%c%c", &sig[0], &sig[1])
	switch string(sig[:]) {
	case "P1", "P2", "P3", "P4", "P5", "P6":
		d.format = int(sig[1] - '0')
	default:
		return ErrFormat
	}

	d.w = d.readInt()
	d.h = d.readInt()
	d.maxval = 255
	if d.format != 1 && d.format != 4 {
		d.maxval = d.readInt()
	}
	if d.maxval <= 0 {
		d.maxval = 255
	}

	if d.err != nil {
		return fmt.Errorf("pnm: %v", d.err)
	}

	d.skipws()

	return nil
}

func (d *decoder) peek() uint8 {
	if d.err != nil {
		return 0
	}

	p, err := d.b.Peek(1)
	if err != nil {
		d.err = err
		return 0
	}
	return p[0]
}

func (d *decoder) getch() byte {
	if d.err != nil {
		return 0
	}
	ch, err := d.b.ReadByte()
	if err != nil {
		d.err = err
		return 0
	}
	return ch
}

func (d *decoder) skipws() {
	for {
		switch d.peek() {
		case ' ', '\t', '\n', '\r':
			d.getch()
			return
		case '#':
			for {
				ch := d.getch()
				if ch == '\n' || ch == 0 {
					break
				}
			}
		default:
			return
		}
	}
}

func (d *decoder) readInt() int {
	if d.err != nil {
		return 0
	}

	d.skipws()
	var n int
	_, d.err = fmt.Fscan(d.b, &n)
	return n
}

func (d *decoder) readColor() color.RGBA {
	var (
		c color.RGBA
		y uint8
	)

	if d.err != nil {
		return c
	}

	switch d.format {
	case 1, 2:
		_, d.err = fmt.Fscanf(d.b, "%d", &y)
		if d.format == 1 && y != 0 {
			y = 255
		}
		c = color.RGBA{y, y, y, 255}
	case 3:
		_, d.err = fmt.Fscanf(d.b, "%d %d %d", &c.R, &c.G, &c.B)
		c.A = 255
	case 4:
		if d.bits == 0 {
			d.err = binary.Read(d.b, binary.LittleEndian, &d.bw)
			d.bits = 8
		} else {
			y := uint8(0)
			if d.bw&(1<<(7-(d.bits-1))) != 0 {
				y = 255
			}
			c = color.RGBA{y, y, y, 255}
			d.bits--
		}
	case 5:
		d.err = binary.Read(d.b, binary.LittleEndian, &y)
		c = color.RGBA{y, y, y, 255}
	case 6:
		var p [6]uint8
		if d.maxval < 256 {
			d.err = binary.Read(d.b, binary.LittleEndian, p[:3])
			c = color.RGBA{p[0], p[1], p[2], 255}
		} else {
			d.err = binary.Read(d.b, binary.LittleEndian, p[:6])
			c = color.RGBA{p[0], p[2], p[4], 255}
		}
	}

	s := 255 / float64(d.maxval)
	cr := float64(c.R) * s
	cg := float64(c.G) * s
	cb := float64(c.B) * s
	return color.RGBA{clamp(cr), clamp(cg), clamp(cb), c.A}
}

func clamp(x float64) uint8 {
	if x > 255 {
		return 255
	}
	return uint8(x)
}

func init() {
	image.RegisterFormat("pbm", "P1", Decode, DecodeConfig)
	image.RegisterFormat("pgm", "P2", Decode, DecodeConfig)
	image.RegisterFormat("ppm", "P3", Decode, DecodeConfig)
	image.RegisterFormat("pbm", "P4", Decode, DecodeConfig)
	image.RegisterFormat("pgm", "P5", Decode, DecodeConfig)
	image.RegisterFormat("ppm", "P6", Decode, DecodeConfig)
}
