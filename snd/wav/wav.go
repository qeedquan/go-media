package wav

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
)

const (
	FORMAT_UNKNOWN           = 0x0000
	FORMAT_PCM               = 0x0001
	FORMAT_ADPCM             = 0x0002
	FORMAT_IEEE_FLOAT        = 0x0003
	FORMAT_ALAW              = 0x0006
	FORMAT_MULAW             = 0x0007
	FORMAT_OKI_ADPCM         = 0x0010
	FORMAT_IMA_ADPCM         = 0x0011
	FORMAT_DIGISTD           = 0x0015
	FORMAT_DIGIFIX           = 0x0016
	FORMAT_DOLBY_AC2         = 0x0030
	FORMAT_GSM610            = 0x0031
	FORMAT_ROCKWELL_ADPCM    = 0x003b
	FORMAT_ROCKWELL_DIGITALK = 0x003c
	FORMAT_G721_ADPCM        = 0x0040
	FORMAT_G728_CELP         = 0x0041
	FORMAT_MPEG              = 0x0050
	FORMAT_MPEGLAYER3        = 0x0055
	FORMAT_G726_ADPCM        = 0x0064
	FORMAT_G722_ADPCM        = 0x0065
	FORMAT_EXTENSIBLE        = 0xfffe
)

type Header struct {
	RIFF          [4]byte
	Size          uint32
	WAVE          [4]byte
	Fmt           [4]byte
	FmtSize       uint32
	Format        uint16
	Channels      uint16
	SampleRate    uint32
	ByteRate      uint32
	BlockAlign    uint16
	BitsPerSample uint16
	Data          [4]byte
	DataSize      uint32
}

type File struct {
	Header  Header
	Samples []uint8
}

func Decode(r io.Reader) (*File, error) {
	f := &File{}
	h := &f.Header
	err := binary.Read(r, binary.LittleEndian, h)
	if err != nil {
		return nil, err
	}
	if h.RIFF != [4]byte{'R', 'I', 'F', 'F'} {
		return nil, fmt.Errorf("invalid riff header")
	}
	if h.WAVE != [4]byte{'W', 'A', 'V', 'E'} {
		return nil, fmt.Errorf("invalid wave header")
	}

	f.Samples, err = io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func Encode(f *File, w io.Writer) error {
	h := f.Header
	h.RIFF = [4]byte{'R', 'I', 'F', 'F'}
	h.WAVE = [4]byte{'W', 'A', 'V', 'E'}
	h.Fmt = [4]byte{'f', 'm', 't', ' '}
	h.Data = [4]byte{'d', 'a', 't', 'a'}

	h.FmtSize = 16
	h.DataSize = uint32(len(f.Samples))
	h.Size = 4 + (8 + h.FmtSize) + (8 + h.DataSize)
	h.ByteRate = h.SampleRate * uint32(h.Channels) * uint32(h.BitsPerSample/8)
	h.BlockAlign = h.Channels * (h.BitsPerSample / 8)

	b := bufio.NewWriter(w)
	binary.Write(b, binary.LittleEndian, &h)
	binary.Write(b, binary.LittleEndian, f.Samples)
	return b.Flush()
}
