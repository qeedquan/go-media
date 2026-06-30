// ported from devil image library

package psd

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"image/color"
	"io"
)

const psdHeader = "8BPS"

type header struct {
	Sig      [4]byte
	Version  uint16
	Reserved [6]byte
	Channels uint16
	Height   uint32
	Width    uint32
	Depth    uint16
	Mode     uint16
}

type decoder struct {
	header
	r   io.Reader
	img image.Image
	bps int
}

func Decode(r io.Reader) (image.Image, error) {
	d := &decoder{r: r}

	err := d.checkHeader()
	if err != nil {
		return nil, err
	}

	err = d.decode()
	if err != nil {
		return nil, err
	}

	return d.img, nil
}

func DecodeConfig(r io.Reader) (image.Config, error) {
	d := &decoder{r: r}
	err := d.checkHeader()
	if err != nil {
		return image.Config{}, err
	}

	return image.Config{
		ColorModel: color.RGBAModel,
		Width:      int(d.Width),
		Height:     int(d.Height),
	}, nil
}

func (d *decoder) checkHeader() error {
	var h header

	err := d.rb(&h)
	if err != nil {
		return err
	}
	if bytes.Compare(h.Sig[:], []byte(psdHeader)) != 0 {
		return errors.New("invalid psd signature")
	}

	if h.Version != 1 {
		return fmt.Errorf("unsupported psd version %d", h.Version)
	}

	for _, v := range h.Reserved {
		if v != 0 {
			return errors.New("reserved field is not all zeroes")
		}
	}

	if h.Channels < 1 || h.Channels > 24 {
		return fmt.Errorf("unsupported number of channels: %d", h.Channels)
	}

	if h.Height < 1 || h.Width < 1 {
		return fmt.Errorf("invalid dimension %dx%d", h.Width, h.Height)
	}

	if h.Depth != 1 && h.Depth != 8 && h.Depth != 16 {
		return fmt.Errorf("unsupported depth: %v", h.Depth)
	}

	d.header = h
	return nil
}

func (d *decoder) decode() error {
	var err error

	switch d.Mode {
	case 1: // grayscale
		err = d.decodeGray()
	case 2: // indexed
		err = d.decodeIndexed()
	case 3: // rgb
		err = d.decodeRGB()
	case 4: // cmyk
		err = d.decodeCMYK()
	default:
		err = fmt.Errorf("mode not supported: %d", d.Mode)
	}

	return err
}

func (d *decoder) decodeGray() error {
	return errors.New("mode gray not supported yet")
}

func (d *decoder) decodeIndexed() error {
	return errors.New("mode indexed not supported yet")
}

func (d *decoder) decodeRGB() error {
	var colorMode, resourceSize, miscInfo uint32
	var compressed uint16

	err := d.rb(&colorMode)
	if err != nil {
		return err
	}

	// skip over color mode data section
	err = nopRead(d.r, int64(colorMode))
	if err != nil {
		return err
	}

	// resources
	err = d.rb(&resourceSize)
	if err != nil {
		return err
	}
	resources := make([]byte, resourceSize)
	_, err = io.ReadAtLeast(d.r, resources, int(resourceSize))
	if err != nil {
		return err
	}

	// misc info
	err = d.rb(&miscInfo)
	if err != nil {
		return err
	}
	err = nopRead(d.r, int64(miscInfo))

	err = d.rb(&compressed)
	if err != nil {
		return err
	}

	var bpc int
	rgba := image.NewRGBA(image.Rect(0, 0, int(d.Width), int(d.Height)))

	// any channels beyond 5 combines into an alpha
	switch c := d.Channels; {
	case c >= 3:
		bpc = 4
	default:
		return fmt.Errorf("unsupported number of rgb channels: %d", d.Channels)
	}

	switch d.Depth {
	case 8, 16:
	default:
		return fmt.Errorf("unsupported rgb depth: %d", d.Depth)
	}

	if compressed == 0 {
	} else {
		err := d.readDataCompressed(rgba.Pix, bpc)
		if err != nil {
			return err
		}
	}

	d.img = rgba
	return nil
}

func (d *decoder) decodeCMYK() error {
	return errors.New("mode CYMK not supported yet")
}

func (d *decoder) readByte() (b byte, err error) {
	var buf [1]byte
	_, err = d.r.Read(buf[:])
	b = buf[0]
	return
}

func (d *decoder) readInt8() (v int8, err error) {
	b, err := d.readByte()
	v = int8(b)
	return
}

func (d *decoder) rb(v interface{}) error {
	return binary.Read(d.r, binary.BigEndian, v)
}

func (d *decoder) numChan() int {
	// needed to correctly load grayscale and paletted images
	switch d.Mode {
	case 1, 2:
		return 1
	}
	return 3
}

func (d *decoder) readDataCompressed(pix []byte, bpc int) error {
	// TODO: add support for this, but very rare format
	if d.Depth == 16 {
		return errors.New("compressed format at depth 16 not supported")
	}

	_, err := d.getCompressedChanLen()
	if err != nil {
		return err
	}

	numChan := d.numChan()
	w, h := int(d.Width), int(d.Height)
	size := w * h
	channel := make([]byte, w*h*int(d.Depth/8))
	c := 0
	for ; c < numChan; c++ {
		err := d.readCompressedChan(size, channel)
		if err != nil {
			return err
		}

		i := 0
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				j := (y*w + x) * bpc
				pix[j+c] = channel[i]
				if channel[i] != 0 {
					pix[j+3] = 255
				}
				i++
			}
		}
	}

	// TODO: this needs to be changed for grayscale images
	if d.Channels >= 4 {
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				i := (y*w + x) * bpc
				pix[i+3] = 255
			}
		}

		for ; c < int(d.Channels); c++ {
			err := d.readCompressedChan(size, channel)
			if err != nil {
				return err
			}

			i := 0
			for y := 0; y < h; y++ {
				for x := 0; x < w; x++ {
					j := (y*w + x) * bpc
					cv := byteToFloat32(pix[j+3])
					nv := byteToFloat32(channel[i])
					pix[j+3] = float32ToByte(cv * nv)
					i++
				}
			}
		}
	}
	return nil
}

func (d *decoder) readCompressedChan(size int, channel []byte) error {
	for i := 0; i < size; {
		b8, err := d.readInt8()
		if err != nil {
			return err
		}
		hb := int(b8)

		if hb >= 0 {
			if i+hb > size {
				return errors.New("corrupted compressed file")
			}

			_, err := io.ReadAtLeast(d.r, channel[i:i+hb+1], hb+1)
			if err != nil {
				return err
			}

			i += hb + 1
		} else if hb >= -127 && hb <= -1 {
			run, err := d.readByte()
			if err != nil {
				return err
			}
			if i+(-hb+1) > size {
				return errors.New("corrupted compressed file")
			}

			for j := 0; j < -hb+1; j++ {
				channel[i+j] = run
			}
			i += -hb + 1
		}
	}
	return nil
}

// this is just a table for fast lookup, not really needed
// if we just want to be simple about reading rle data
// but we still need to read it in because we cant seek
func (d *decoder) getCompressedChanLen() ([]uint, error) {
	tab := make([]uint16, d.Height*uint32(d.Channels))
	cl := make([]uint, d.Channels)
	for i := range tab {
		err := d.rb(&tab[i])
		if err != nil {
			return nil, err
		}
	}
	for c := 0; c < int(d.Channels); c++ {
		j := c * int(d.Height)
		for i := 0; i < int(d.Height); i++ {
			cl[c] = uint(tab[i+j])
		}
	}
	return cl, nil
}

func nopRead(r io.Reader, length int64) error {
	var buf [32768]byte

	for length > 0 {
		n := int64(len(buf))
		if n > length {
			n = length
		}

		_, err := io.ReadAtLeast(r, buf[:n], int(n))
		if err != nil {
			return err
		}
		length -= n
	}

	return nil
}

func byteToFloat32(v byte) float32 { return float32(v) / 255 }
func float32ToByte(v float32) byte { return byte(v * 255) }

func init() {
	image.RegisterFormat("psd", psdHeader, Decode, DecodeConfig)
}
