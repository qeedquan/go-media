package pnm

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"io"
)

type Options struct {
	Format int
}

func Encode(w io.Writer, m image.Image, o *Options) error {
	if o == nil {
		o = &Options{Format: 3}
	}

	b := bufio.NewWriter(w)
	r := m.Bounds()
	bits := uint(0)
	bw := uint8(0)

	fmt.Fprintf(w, "P%d\n", o.Format)
	fmt.Fprintf(w, "%d %d\n", r.Dx(), r.Dy())
	switch o.Format {
	case 1, 4:
	case 2, 3, 5, 6:
		fmt.Fprintf(w, "255\n")
	default:
		return ErrFormat
	}

	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			p := m.At(x, y)
			switch o.Format {
			case 1, 2:
				c := (color.GrayModel.Convert(p)).(color.Gray)
				if o.Format == 1 && c.Y != 0 {
					c.Y = 1
				}
				fmt.Fprintf(w, "%d ", c.Y)

			case 3:
				c := color.RGBAModel.Convert(p).(color.RGBA)
				fmt.Fprintf(w, "%d %d %d", c.R, c.G, c.B)

			case 4:
				c := color.RGBAModel.Convert(p).(color.Gray)
				if c.Y != 0 {
					bw |= 1 << (7 - bits)
				}

				if bits++; bits == 8 {
					binary.Write(w, binary.LittleEndian, bw)
					bits = 0
					bw = 0
				}

			case 5:
				c := (color.GrayModel.Convert(p)).(color.Gray)
				binary.Write(w, binary.LittleEndian, c.Y)

			case 6:
				c := color.RGBAModel.Convert(p).(color.RGBA)
				p := [3]uint8{c.R, c.G, c.B}
				binary.Write(w, binary.LittleEndian, p)
			}

			if x+1 < r.Max.X {
				fmt.Fprintf(w, " ")
			}
		}

		switch o.Format {
		case 1, 2, 3:
			fmt.Fprintf(w, "\n")
		}
	}

	if bits != 0 {
		binary.Write(w, binary.LittleEndian, bw)
	}

	err := b.Flush()
	if err != nil {
		return fmt.Errorf("pnm: %v", err)
	}
	return nil
}
