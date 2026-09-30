package tx2

import (
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"io"
)

type Header struct {
	Width       uint16
	Height      uint16
	Type        uint16
	Unkown      uint16
	NumColors   uint16
	NumPalettes uint16
	Pad         uint32
}

type Image struct {
	Header
	Palette [][]color.RGBA
	Pix     []byte
}

const (
	TYPE_DXT1        = 0
	TYPE_DXT4        = 1
	TYPE_DXT5        = 2
	TYPE_BGRA        = 3
	TYPE_PAL_BGRA16  = 16
	TYPE_PAL_RGBA16  = 17
	TYPE_PAL_BGRA256 = 256
	TYPE_PAL_RGBA256 = 257
	TYPE_FONT        = 4
	TYPE_ERROR       = 999
)

func Decode(r io.Reader) (*Image, error) {
	img := &Image{}
	err := binary.Read(r, binary.LittleEndian, &img.Header)
	if err != nil {
		return nil, err
	}

	palsize := img.NumPalettes
	if palsize == 0 {
		switch img.Type {
		case TYPE_BGRA, TYPE_DXT1, TYPE_DXT4,
			TYPE_DXT5, TYPE_ERROR:
		default:
			palsize = 1
		}
	}
	img.Palette = make([][]color.RGBA, palsize)

	for i := range img.Palette {
		img.Palette[i] = make([]color.RGBA, img.NumColors)
		for j := range img.NumColors {
			err := binary.Read(r, binary.LittleEndian, &img.Palette[i][j])
			if err != nil {
				return nil, err
			}
			switch img.Type {
			case TYPE_PAL_BGRA16, TYPE_PAL_BGRA256:
				img.Palette[i][j].R, img.Palette[i][j].B = img.Palette[i][j].B, img.Palette[i][j].R
			}
		}
	}

	img.Pix, err = io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func (c *Image) ToRGBA() (*image.RGBA, error) {
	width := int(c.Width)
	height := int(c.Height)
	pix := image.NewRGBA(image.Rect(0, 0, width, height))

	switch c.Type {
	case TYPE_BGRA:
		if len(c.Pix) < width*height*4 {
			return nil, fmt.Errorf("image format %s too small for output", typestr(c.Type))
		}

		i := 0
		for y := range height {
			for x := range width {
				pix.Set(x, y, color.RGBA{
					c.Pix[i+2],
					c.Pix[i+1],
					c.Pix[i],
					c.Pix[i+3],
				})
				i += 4
			}
		}

	case TYPE_PAL_BGRA16, TYPE_PAL_RGBA16:
		if len(c.Palette) < 1 || len(c.Palette[0]) != 16 {
			return nil, fmt.Errorf("image format %s has invalid palette size", typestr(c.Type))
		}
		if len(c.Pix) < width*height/2 {
			return nil, fmt.Errorf("image format %s too small for output", typestr(c.Type))
		}

		i := 0
		for y := range height {
			for x := 0; x+1 < width; x += 2 {
				lo := c.Pix[i] & 0xf
				hi := (c.Pix[i] >> 4) & 0xf
				pix.Set(x, y, c.Palette[0][lo])
				pix.Set(x+1, y, c.Palette[0][hi])
				i += 1
			}
		}

	case TYPE_PAL_BGRA256, TYPE_PAL_RGBA256:
		i := 0
		for y := range height {
			for x := range width {
				sel := c.Pix[i]
				pix.Set(x, y, c.Palette[0][sel])
				i += 1
			}
		}

	default:
		return nil, fmt.Errorf("image format %s to RGBA is unsupported", typestr(c.Type))
	}

	return pix, nil
}

func typestr(typ uint16) string {
	switch typ {
	case TYPE_DXT1:
		return "DXT1"
	case TYPE_DXT4:
		return "DXT4"
	case TYPE_DXT5:
		return "DXT5"
	case TYPE_BGRA:
		return "BGRA"
	case TYPE_PAL_BGRA16:
		return "PAL_BGRA16"
	case TYPE_PAL_RGBA16:
		return "PAL_RGBA16"
	case TYPE_PAL_BGRA256:
		return "PAL_BGRA256"
	case TYPE_PAL_RGBA256:
		return "PAL_RGBA256"
	case TYPE_FONT:
		return "FONT"
	case TYPE_ERROR:
		return "ERROR"
	}
	return "UNKNOWN"
}
