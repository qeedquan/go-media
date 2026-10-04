package dds

import (
	"encoding/binary"
	"io"
)

type Header struct {
	Size              uint32
	Flags             uint32
	Height            uint32
	Width             uint32
	PitchOrLinearSize uint32
	Depth             uint32
	MipMapCount       uint32
	Reserved          [11]uint32
	Spf               PixelFormat
	Caps              uint32
	Caps2             uint32
	Caps3             uint32
	Caps4             uint32
	Reserved2         uint32
}

type PixelFormat struct {
	Size        uint32
	Flags       uint32
	FourCC      uint32
	RGBBitCount uint32
	RBitMask    uint32
	GBitMask    uint32
	BBitMask    uint32
	ABitMask    uint32
}

type DXT10 struct {
	Format            uint32
	ResourceDimension uint32
	MiscFlag          uint32
	ArraySize         uint32
	MiscFlags2        uint32
}

type Image struct {
	Header Header
	DXT10  DXT10
	Pix    []byte
}

var (
	MAGIC = [4]byte{'D', 'D', 'S', ' '}
)

const (
	HEADER_SIZE      = 124
	PIXELFORMAT_SIZE = 32
)

const (
	FOURCC_DXT1  = 0x31545844
	FOURCC_DXT2  = 0x32545844
	FOURCC_DXT3  = 0x33545844
	FOURCC_DXT4  = 0x34545844
	FOURCC_DXT5  = 0x35545844
	FOURCC_DXT10 = 0x30315844
)

const (
	DDSD_CAPS        = 0x1
	DDSD_HEIGHT      = 0x2
	DDSD_WIDTH       = 0x4
	DDSD_PITCH       = 0x8
	DDSD_PIXELFORMAT = 0x1000
	DDSD_MIPMAPCOUNT = 0x20000
	DDSD_LINEARSIZE  = 0x80000
	DDSD_DEPTH       = 0x800000
)

const (
	DDSCAPS_COMPLEX = 0x8
	DDSCAPS_MIPMAP  = 0x400000
	DDSCAPS_TEXTURE = 0x1000
)

const (
	DDSCAPS2_CUBEMAP           = 0x200
	DDSCAPS2_CUBEMAP_POSITIVEX = 0x400
	DDSCAPS2_CUBEMAP_NEGATIVEX = 0x800
	DDSCAPS2_CUBEMAP_POSITIVEY = 0x1000
	DDSCAPS2_CUBEMAP_NEGATIVEY = 0x2000
	DDSCAPS2_CUBEMAP_POSITIVEZ = 0x4000
	DDSCAPS2_CUBEMAP_NEGATIVEZ = 0x8000
	DDSCAPS2_VOLUME            = 0x200000
)

const (
	DDPF_ALPHAPIXELS = 0x1
	DDPF_ALPHA       = 0x2
	DDPF_FOURCC      = 0x4
	DDPF_RGB         = 0x40
	DDPF_YUV         = 0x200
	DDPF_LUMINANCE   = 0x20000
)

func Encode(w io.Writer, m *Image) error {
	binary.Write(w, binary.LittleEndian, MAGIC)
	binary.Write(w, binary.LittleEndian, m.Header)
	_, err := w.Write(m.Pix)
	return err
}
