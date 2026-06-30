package coff

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
)

const (
	F_RELFLG   = 0x1
	F_EXEC     = 0x2
	F_LNNO     = 0x4
	F_LSYMS    = 0x8
	F_LITTLE   = 0x100
	F_BIG      = 0x200
	F_SYMMERGE = 0x1000
)

const (
	TMS470       = 0x97
	TMS320C5400  = 0x98
	TMS320C6000  = 0x99
	TMS320C5500  = 0x9c
	TMS320C2800  = 0x9d
	MPS430       = 0xa0
	TMS320C5500P = 0xa1
)

const (
	STYP_REG    = 0x0
	STYP_DSECT  = 0x1
	STYP_NOLOAD = 0x2
	STYP_GROUP  = 0x4
	STYP_PAD    = 0x8
	STYP_COPY   = 0x10
	STYP_TEXT   = 0x20
	STYP_DATA   = 0x40
	STYP_BSS    = 0x80
	STYP_BLOCK  = 0x1000
	STYP_PASS   = 0x2000
	STYP_CLINK  = 0x4000
	STYP_VECTOR = 0x8000
	STYP_PADDED = 0x10000
)

const (
	RE_ADD    = 0x4000
	RE_SUB    = 0x4001
	RE_NEG    = 0x4002
	RE_MPY    = 0x4003
	RE_DIV    = 0x4004
	RE_MOD    = 0x4005
	RE_SR     = 0x4006
	RE_ASR    = 0x4007
	RE_SL     = 0x4008
	RE_AND    = 0x4009
	RE_OR     = 0x400a
	RE_XOR    = 0x400b
	RE_NOTB   = 0x400c
	RE_ULDFLD = 0x400d
	RE_SLDFLD = 0x400e
	RE_USTFLD = 0x400f
	RE_SSTFLD = 0x4010
	RE_PUSH   = 0x4011
	RE_PUSHSK = 0x4012
	RE_PUSHUK = 0x4013
	RE_PUSHPC = 0x4014
	RE_DUP    = 0x4015
	RE_XSTFLD = 0x4016
	RE_PUSHV  = 0xc011
)

const (
	R_60ABS     = 0x0000
	R_60RELBYTE = 0x000f
	R_60RELWORD = 0x0010
	R_60RELLONG = 0x0011
	R_C60BASE   = 0x0050
	R_C60DIR15  = 0x0051
	R_C60PCR21  = 0x0052
	R_C60PCR10  = 0x0053
	R_C60LO16   = 0x0054
	R_C60HI16   = 0x0055
	R_C60SECT   = 0x0056
	R_C60S16    = 0x0057
	R_C60PCR7   = 0x0070
	R_C60PCR12  = 0x0071
)

const (
	R_C28ABS       = 0x0000
	R_C28RELBYTE   = 0x000f
	R_C28RELWORD   = 0x0010
	R_C28RELLONG   = 0x0011
	R_C28PARTLS7   = 0x0028
	R_C28PARTLS6   = 0x005d
	R_C28PARTMID10 = 0x005e
	R_C28REL22     = 0x005f
	R_C28PARTMS6   = 0x0060
	R_C28PARTS16   = 0x0061
	R_C28PCR16     = 0x0062
	R_C28PCR8      = 0x0063
	R_C28PTR       = 0x0064
	R_C28HI16      = 0x0065
	R_C28LOPTR     = 0x0066
	R_C28NWORD     = 0x0067
	R_C28NBYTE     = 0x0068
	R_C28HIBYTE    = 0x0069
	R_C28RELS13    = 0x006a
)

const (
	R_C54ABS     = 0x0000
	R_C54REL24   = 0x0005
	R_C54RELBYTE = 0x017
	R_C54RELWORD = 0x0020
	R_C54RELLONG = 0x0021
	R_C54PARTSL7 = 0x0028
	R_C54PARTSM9 = 0x0029
	R_C54REL13   = 0x002a
)

const (
	R_C55ABS        = 0x0000
	R_C55REL24      = 0x0005
	R_C55RELBYTE    = 0x0017
	R_C55RELWORD    = 0x0020
	R_C55RELLONG    = 0x0021
	R_C55_LD3_DMA   = 0x0170
	R_C55_LD3_MDP   = 0x0172
	R_C55_LD3_PDP   = 0x0173
	R_C55_LD3_REL23 = 0x0174
	R_C55_LD3_k8    = 0x0210
	R_C55_LD3_k16   = 0x0211
	R_C55_LD3_K8    = 0x0212
	R_C55_LD3_K16   = 0x0213
	R_C55_LD3_l8    = 0x0214
	R_C55_LD3_l16   = 0x0215
	R_C55_LD3_L8    = 0x0216
	R_C55_LD3_L16   = 0x0217
	R_C55_LD3_k4    = 0x0220
	R_C55_LD3_k5    = 0x0221
	R_C55_LD3_K5    = 0x0222
	R_C55_LD3_k6    = 0x0223
	R_C55_LD3_k12   = 0x0224
)

const (
	R_MTRELLONG = 0x0011
	R_MTPCR23H  = 0x0061
	R_MTPCR24W  = 0x0017
)

const (
	N_UNDEF = 0
	N_ABS   = 1
	N_DEBUG = 2
)

const (
	C_EFCN    = -1
	C_NULL    = 0
	C_AUTO    = 1
	C_EXT     = 2
	C_STAT    = 3
	C_REG     = 4
	C_EXTDEF  = 5
	C_LABEL   = 6
	C_ULABEL  = 7
	C_MOS     = 8
	C_ARG     = 9
	C_STRTAG  = 10
	C_MOU     = 11
	C_UNTAG   = 12
	C_TPDEF   = 13
	C_USTATIC = 14
	C_ENTAG   = 15
	C_MOE     = 16
	C_REGPARM = 17
	C_FIELD   = 18
	C_BLOCK   = 100
	C_FCN     = 101
	C_EOS     = 102
	C_FILE    = 103
	C_LINE    = 104
	C_ALIAS   = 105
	C_HIDDEN  = 106
	C_SHADOW  = 107
	C_WEAKEXT = 108
)

type File struct {
	Header
	OptionalHeader *OptionalHeader
	Sections       []*Section
	Symbols        []*Symbol
	Strings        []byte
	r              io.ReaderAt
	closer         io.Closer
}

type Header struct {
	Version            uint16
	NumSections        uint16
	Timestamp          uint32
	SymbolOff          uint32
	NumSymbols         uint32
	OptionalHeaderSize uint16
	Flags              uint16
	TargetID           uint16
}

type OptionalHeader struct {
	Magic    uint16
	Version  uint16
	TextSize uint32
	DataSize uint32
	BSSSize  uint32
	Entry    uint32
	TextAddr uint32
	DataAddr uint32
}

type Section struct {
	SectionHeader
	Data   []byte
	Relocs []Reloc
}

type SectionHeader struct {
	Name      string
	NameOff   int64
	PhysAddr  uint32
	VirtAddr  uint32
	Size      uint32
	DataOff   uint32
	RelocOff  uint32
	LineOff   uint32
	NumRelocs uint32
	NumLines  uint32
	Flags     uint32
	Page      uint16
}

type SectionHeader1 struct {
	Name      [8]byte
	PhysAddr  uint32
	VirtAddr  uint32
	Size      uint32
	DataOff   uint32
	RelocOff  uint32
	LineOff   uint32
	NumRelocs uint16
	NumLines  uint16
	Flags     uint16
	_         byte
	Page      uint8
}

type SectionHeader2 struct {
	Name      [8]byte
	PhysAddr  uint32
	VirtAddr  uint32
	Size      uint32
	DataOff   uint32
	RelocOff  uint32
	LineOff   uint32
	NumRelocs uint32
	NumLines  uint32
	Flags     uint32
	_         [2]byte
	Page      uint16
}

type Line struct {
	Loc uint32
	Num uint16
}

type Reloc struct {
	VirtAddr    uint64
	SymbolIndex uint64
	Extended    uint16
	Type        uint16
}

type Reloc1 struct {
	VirtAddr    uint32
	SymbolIndex uint32
	Extended    uint16
	Type        uint16
}

type Reloc2 struct {
	VirtAddr    uint32
	SymbolIndex uint16
	_           [2]byte
	Type        uint16
}

type Symbol struct {
	Name    string
	NameOff int64
	Value   uint32
	Section uint16
	Type    uint16
	Class   uint8
	Aux     uint8
}

type Symbol1 struct {
	Name    [8]byte
	Value   uint32
	Section uint16
	Type    uint16
	Class   uint8
	Aux     uint8
}

func NewFile(r io.ReaderAt) (*File, error) {
	f := &File{r: r}
	sr := io.NewSectionReader(r, 0, math.MaxUint32)
	err := binary.Read(sr, binary.LittleEndian, &f.Header)
	if err != nil {
		return nil, fmt.Errorf("coff: failed to read header: %v", err)
	}

	switch f.Version {
	case 0xc1, 0xc2:
	default:
		return nil, fmt.Errorf("coff: unsupported version %#x", f.Version)
	}

	if f.OptionalHeaderSize != 0 {
		f.OptionalHeader = &OptionalHeader{}
		err = binary.Read(sr, binary.LittleEndian, f.OptionalHeader)
		if err != nil {
			return nil, fmt.Errorf("coff: failed to read optional header: %v", err)
		}
	}

	for i := uint16(0); i < f.NumSections; i++ {
		s := &Section{}
		switch f.Version {
		case 0xc1:
			var h SectionHeader1
			err = binary.Read(sr, binary.LittleEndian, &h)
			s.SectionHeader = SectionHeader{
				Name:      string(h.Name[:]),
				NameOff:   f.nameoff(h.Name[:]),
				PhysAddr:  h.PhysAddr,
				VirtAddr:  h.VirtAddr,
				Size:      h.Size,
				DataOff:   h.DataOff,
				RelocOff:  h.RelocOff,
				LineOff:   h.LineOff,
				NumRelocs: uint32(h.NumRelocs),
				NumLines:  uint32(h.NumLines),
				Flags:     uint32(h.Flags),
				Page:      uint16(h.Page),
			}

		default:
			var h SectionHeader2
			err = binary.Read(sr, binary.LittleEndian, &h)
			s.SectionHeader = SectionHeader{
				Name:      string(h.Name[:]),
				NameOff:   f.nameoff(h.Name[:]),
				PhysAddr:  h.PhysAddr,
				VirtAddr:  h.VirtAddr,
				Size:      h.Size,
				DataOff:   h.DataOff,
				RelocOff:  h.RelocOff,
				LineOff:   h.LineOff,
				NumRelocs: h.NumRelocs,
				NumLines:  h.NumLines,
				Flags:     h.Flags,
				Page:      h.Page,
			}
		}

		if err != nil {
			return nil, fmt.Errorf("coff: failed to read section header %d: %v", i+1, err)
		}
		f.Sections = append(f.Sections, s)
	}

	for i, s := range f.Sections {
		if s.DataOff != 0 {
			s.Data = make([]byte, s.Size)
			sr = io.NewSectionReader(r, int64(s.DataOff), math.MaxUint32)
			err = binary.Read(sr, binary.LittleEndian, s.Data)
			if err != nil {
				return nil, fmt.Errorf("coff: failed to read section %d data: %v", i+1, err)
			}
		}

		if s.RelocOff != 0 {
			sr = io.NewSectionReader(r, int64(s.RelocOff), math.MaxUint32)
			for j := uint32(0); j < s.NumRelocs; j++ {
				switch f.Version {
				case 0xc1:
					var rl Reloc1
					err = binary.Read(sr, binary.LittleEndian, &rl)
					s.Relocs = append(s.Relocs, Reloc{
						VirtAddr:    uint64(rl.VirtAddr),
						SymbolIndex: uint64(rl.SymbolIndex),
						Extended:    rl.Extended,
						Type:        rl.Type,
					})

				default:
					var rl Reloc2
					err = binary.Read(sr, binary.LittleEndian, &rl)
					s.Relocs = append(s.Relocs, Reloc{
						VirtAddr:    uint64(rl.VirtAddr),
						SymbolIndex: uint64(rl.SymbolIndex),
						Type:        rl.Type,
					})
				}
				if err != nil {
					return nil, fmt.Errorf("coff: failed to read section %d relocation entry %d: %v", i+1, j+1, err)
				}
			}
		}
	}

	sr = io.NewSectionReader(r, int64(f.SymbolOff), math.MaxUint32)
	for i := uint32(0); i < f.NumSymbols; i++ {
		var y Symbol1
		err = binary.Read(sr, binary.LittleEndian, &y)
		if err != nil {
			return nil, fmt.Errorf("coff: failed to read symbol %d: %v", i+1, err)
		}

		name := ""
		nameoff := f.nameoff(y.Name[:])
		if nameoff == 0 {
			name = f.mkname(y.Name[:], nameoff)
		}

		f.Symbols = append(f.Symbols, &Symbol{
			Name:    name,
			NameOff: nameoff,
			Value:   y.Value,
			Section: y.Section,
			Type:    y.Type,
			Class:   y.Class,
			Aux:     y.Aux,
		})
	}

	var numstr uint32
	err = binary.Read(sr, binary.LittleEndian, &numstr)
	if err != nil {
		return nil, fmt.Errorf("coff: failed to read string table length: %v", err)
	}
	sr.Seek(-4, io.SeekCurrent)

	f.Strings = make([]byte, numstr)
	err = binary.Read(sr, binary.LittleEndian, f.Strings)
	if err != nil {
		return nil, fmt.Errorf("coff: failed to read string table of length %d: %v", numstr, err)
	}

	for _, s := range f.Sections {
		s.Name = f.mkname([]byte(s.Name), s.NameOff)
	}
	for _, y := range f.Symbols {
		if y.NameOff != 0 {
			y.Name = f.mkname(nil, y.NameOff)
		}
	}

	return f, nil
}

func (f *File) nameoff(p []byte) int64 {
	z := binary.LittleEndian.Uint32(p[0:])
	o := binary.LittleEndian.Uint32(p[4:])
	if z == 0 {
		return int64(o)
	}
	return 0
}

func (f *File) mkname(s []byte, o int64) string {
	if o == 0 {
		return strz(s)
	}

	if o < int64(len(f.Strings)) {
		return strz(f.Strings[o:])
	}

	return ""
}

func (f *File) mkoname(p []byte, name string, nameoff int64) {
	if nameoff == 0 {
		copy(p, []byte(name))
	} else {
		binary.LittleEndian.PutUint32(p[0:], 0)
		binary.LittleEndian.PutUint32(p[4:], uint32(nameoff))
	}
}

func strz(s []byte) string {
	for i := range s {
		if s[i] == 0 {
			return string(s[:i])
		}
	}
	return string(s)
}

func Open(name string) (*File, error) {
	r, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	f, err := NewFile(r)
	if err != nil {
		return nil, err
	}
	f.closer = r
	return f, nil
}

func (f *File) Close() error {
	var err error
	if f.closer != nil {
		err = f.closer.Close()
		f.closer = nil
	}
	return err
}

func (f *File) Section(name string) *Section {
	for _, s := range f.Sections {
		if s.Name == name {
			return s
		}
	}
	return nil
}

func (f *File) LookupSymbol(name string) *Symbol {
	for _, y := range f.Symbols {
		if y.Name == name {
			return y
		}
	}
	return nil
}

type TargetID uint16

func (t TargetID) String() string {
	switch t {
	case TMS470:
		return "TMS470"
	case TMS320C5400:
		return "TMS320C5400"
	case TMS320C6000:
		return "TMS320C6000"
	case TMS320C5500:
		return "TMS320C5500"
	case TMS320C2800:
		return "TMS320C2800"
	case MPS430:
		return "MPS430"
	case TMS320C5500P:
		return "TMS320C5500P"
	}
	return "UNKNOWN"
}

func Format(f *File, w io.Writer) (n int, err error) {
	b := new(bytes.Buffer)
	binary.Write(b, binary.LittleEndian, &f.Header)
	if f.OptionalHeader != nil {
		binary.Write(b, binary.LittleEndian, f.OptionalHeader)
	}

	for _, s := range f.Sections {
		switch f.Version {
		case 0xc1:
			h := SectionHeader1{
				PhysAddr:  s.PhysAddr,
				VirtAddr:  s.VirtAddr,
				Size:      s.Size,
				DataOff:   s.DataOff,
				RelocOff:  s.RelocOff,
				LineOff:   s.LineOff,
				NumRelocs: uint16(s.NumRelocs),
				NumLines:  uint16(s.NumLines),
				Flags:     uint16(s.Flags),
				Page:      uint8(s.Page),
			}
			f.mkoname(h.Name[:], s.Name, s.NameOff)
			binary.Write(b, binary.LittleEndian, &h)

		default:
			h := SectionHeader2{
				PhysAddr:  s.PhysAddr,
				VirtAddr:  s.VirtAddr,
				Size:      s.Size,
				DataOff:   s.DataOff,
				RelocOff:  s.RelocOff,
				LineOff:   s.LineOff,
				NumRelocs: s.NumRelocs,
				NumLines:  s.NumLines,
				Flags:     s.Flags,
				Page:      s.Page,
			}
			f.mkoname(h.Name[:], s.Name, s.NameOff)
			binary.Write(b, binary.LittleEndian, &h)
		}
	}

	for _, s := range f.Sections {
		binary.Write(b, binary.LittleEndian, s.Data)
	}

	for _, y := range f.Symbols {
		yl := Symbol1{
			Value:   y.Value,
			Section: y.Section,
			Type:    y.Type,
			Class:   y.Class,
			Aux:     y.Aux,
		}
		f.mkoname(yl.Name[:], y.Name, y.NameOff)
		binary.Write(b, binary.LittleEndian, &yl)
	}

	binary.Write(b, binary.LittleEndian, f.Strings)

	return w.Write(b.Bytes())
}
