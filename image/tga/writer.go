package tga

import (
	"bufio"
	"encoding/binary"
	"image"
	"io"
)

func Encode(w io.Writer, m image.Image) (err error) {
	bw := bufio.NewWriter(w)
	defer func() {
		xerr := bw.Flush()
		if err == nil {
			err = xerr
		}
	}()

	width, height := m.Bounds().Dx(), m.Bounds().Dy()
	head := header{
		Type:   2,
		Bpp:    32,
		Width:  uint16(width),
		Height: uint16(height),
		Desc:   0x20,
	}

	if err := binary.Write(bw, binary.LittleEndian, head); err != nil {
		return err
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := m.At(x, y).RGBA()
			err := binary.Write(bw, binary.LittleEndian,
				[]byte{byte(b), byte(g), byte(r), byte(a)})
			if err != nil {
				return err
			}
		}
	}

	return nil
}
