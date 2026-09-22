package peutil

import (
	"bufio"
	"bytes"
	"debug/pe"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"reflect"
	"sort"
	"strings"

	"github.com/qeedquan/go-media/debug"
	"github.com/qeedquan/go-media/math/ga"
)

const (
	IMAGE_SCN_TYPE_NO_PAD            = 0x00000008
	IMAGE_SCN_CNT_CODE               = 0x00000020
	IMAGE_SCN_CNT_INITIALIZED_DATA   = 0x00000040
	IMAGE_SCN_CNT_UNINITIALIZED_DATA = 0x00000080
	IMAGE_SCN_LNK_OTHER              = 0x00000100
	IMAGE_SCN_LNK_INFO               = 0x00000200
	IMAGE_SCN_LNK_REMOVE             = 0x00000800
	IMAGE_SCN_LNK_COMDAT             = 0x00001000
	IMAGE_SCN_GPREL                  = 0x00008000
	IMAGE_SCN_MEM_PURGEABLE          = 0x00020000
	IMAGE_SCN_MEM_16BIT              = 0x00020000
	IMAGE_SCN_MEM_LOCKED             = 0x00040000
	IMAGE_SCN_MEM_PRELOAD            = 0x00080000
	IMAGE_SCN_ALIGN_1BYTES           = 0x00100000
	IMAGE_SCN_ALIGN_2BYTES           = 0x00200000
	IMAGE_SCN_ALIGN_8BYTES           = 0x00400000
	IMAGE_SCN_ALIGN_16BYTES          = 0x00500000
	IMAGE_SCN_ALIGN_32BYTES          = 0x00600000
	IMAGE_SCN_ALIGN_64BYTES          = 0x00700000
	IMAGE_SCN_ALIGN_128BYTES         = 0x00800000
	IMAGE_SCN_ALIGN_256BYTES         = 0x00900000
	IMAGE_SCN_ALIGN_512BYTES         = 0x00A00000
	IMAGE_SCN_ALIGN_1024BYTES        = 0x00B00000
	IMAGE_SCN_ALIGN_2048BYTES        = 0x00C00000
	IMAGE_SCN_ALIGN_4096BYTES        = 0x00D00000
	IMAGE_SCN_ALIGN_8192BYTES        = 0x00E00000
	IMAGE_SCN_LNK_NRELOC_OVFL        = 0x01000000
	IMAGE_SCN_MEM_DISCARDABLE        = 0x02000000
	IMAGE_SCN_MEM_NOT_CACHED         = 0x04000000
	IMAGE_SCN_MEM_NOT_PAGED          = 0x08000000
	IMAGE_SCN_MEM_SHARED             = 0x10000000
	IMAGE_SCN_MEM_EXECUTE            = 0x20000000
	IMAGE_SCN_MEM_READ               = 0x40000000
	IMAGE_SCN_MEM_WRITE              = 0x80000000
)

const (
	IMAGE_FILE_RELOCS_STRIPPED         = 0x0001
	IMAGE_FILE_EXECUTABLE_IMAGE        = 0x0002
	IMAGE_FILE_LINE_NUMS_STRIPPED      = 0x0004
	IMAGE_FILE_LOCAL_SYMS_STRIPPED     = 0x0008
	IMAGE_FILE_AGGRESSIVE_WS_TRIM      = 0x0010
	IMAGE_FILE_LARGE_ADDRESS_AWARE     = 0x0020
	IMAGE_FILE_BYTES_REVERSED_LO       = 0x0080
	IMAGE_FILE_32BIT_MACHINE           = 0x0100
	IMAGE_FILE_DEBUG_STRIPPED          = 0x0200
	IMAGE_FILE_REMOVABLE_RUN_FROM_SWAP = 0x0400
	IMAGE_FILE_NET_RUN_FROM_SWAP       = 0x0800
	IMAGE_FILE_SYSTEM                  = 0x1000
	IMAGE_FILE_DLL                     = 0x2000
	IMAGE_FILE_UP_SYSTEM_ONLY          = 0x4000
	IMAGE_FILE_BYTES_REVERSED_HI       = 0x8000
)

const (
	IMAGE_SUBSYSTEM_UNKNOWN                  = 0
	IMAGE_SUBSYSTEM_NATIVE                   = 1
	IMAGE_SUBSYSTEM_WINDOWS_GUI              = 2
	IMAGE_SUBSYSTEM_WINDOWS_CUI              = 3
	IMAGE_SUBSYSTEM_OS2_CUI                  = 5
	IMAGE_SUBSYSTEM_POSIX_CUI                = 7
	IMAGE_SUBSYSTEM_NATIVE_WINDOWS           = 8
	IMAGE_SUBSYSTEM_WINDOWS_CE_GUI           = 9
	IMAGE_SUBSYSTEM_EFI_APPLICATION          = 10
	IMAGE_SUBSYSTEM_EFI_BOOT_SERVICE_DRIVER  = 11
	IMAGE_SUBSYSTEM_EFI_RUNTIME_DRIVER       = 12
	IMAGE_SUBSYSTEM_EFI_ROM                  = 13
	IMAGE_SUBSYSTEM_XBOX                     = 14
	IMAGE_SUBSYSTEM_WINDOWS_BOOT_APPLICATION = 16
)

const (
	IMAGE_DLLCHARACTERISTICS_HIGH_ENTROPY_VA       = 0x0020
	IMAGE_DLLCHARACTERISTICS_DYNAMIC_BASE          = 0x0040
	IMAGE_DLLCHARACTERISTICS_FORCE_INTEGRITY       = 0x0080
	IMAGE_DLLCHARACTERISTICS_NX_COMPAT             = 0x0100
	IMAGE_DLLCHARACTERISTICS_NO_ISOLATION          = 0x0200
	IMAGE_DLLCHARACTERISTICS_NO_SEH                = 0x0400
	IMAGE_DLLCHARACTERISTICS_NO_BIND               = 0x0800
	IMAGE_DLLCHARACTERISTICS_APPCONTAINER          = 0x1000
	IMAGE_DLLCHARACTERISTICS_WDM_DRIVER            = 0x2000
	IMAGE_DLLCHARACTERISTICS_GUARD_CF              = 0x4000
	IMAGE_DLLCHARACTERISTICS_TERMINAL_SERVER_AWARE = 0x8000
)

const (
	IMAGE_DEBUG_TYPE_UNKNOWN       = 0
	IMAGE_DEBUG_TYPE_COFF          = 1
	IMAGE_DEBUG_TYPE_CODEVIEW      = 2
	IMAGE_DEBUG_TYPE_FPO           = 3
	IMAGE_DEBUG_TYPE_MISC          = 4
	IMAGE_DEBUG_TYPE_EXCEPTION     = 5
	IMAGE_DEBUG_TYPE_FIXUP         = 6
	IMAGE_DEBUG_TYPE_OMAP_TO_SRC   = 7
	IMAGE_DEBUG_TYPE_OMAP_FROM_SRC = 8
	IMAGE_DEBUG_TYPE_BORLAND       = 9
	IMAGE_DEBUG_TYPE_RESERVED10    = 10
	IMAGE_DEBUG_TYPE_CLSID         = 11
	IMAGE_DEBUG_TYPE_REPRO         = 16
)

const (
	IMAGE_REL_AMD64_ABSOLUTE = 0x0000
	IMAGE_REL_AMD64_ADDR64   = 0x0001
	IMAGE_REL_AMD64_ADDR32   = 0x0002
	IMAGE_REL_AMD64_ADDR32NB = 0x0003
	IMAGE_REL_AMD64_REL32    = 0x0004
	IMAGE_REL_AMD64_REL32_1  = 0x0005
	IMAGE_REL_AMD64_REL32_2  = 0x0006
	IMAGE_REL_AMD64_REL32_3  = 0x0007
	IMAGE_REL_AMD64_REL32_4  = 0x0008
	IMAGE_REL_AMD64_REL32_5  = 0x0009
	IMAGE_REL_AMD64_SECTION  = 0x000A
	IMAGE_REL_AMD64_SECREL   = 0x000B
	IMAGE_REL_AMD64_SECREL7  = 0x000C
	IMAGE_REL_AMD64_TOKEN    = 0x000D
	IMAGE_REL_AMD64_SREL32   = 0x000E
	IMAGE_REL_AMD64_PAIR     = 0x000F
	IMAGE_REL_AMD64_SSPAN32  = 0x0010
)

const (
	IMAGE_REL_I386_ABSOLUTE = 0x0000
	IMAGE_REL_I386_DIR16    = 0x0001
	IMAGE_REL_I386_REL16    = 0x0002
	IMAGE_REL_I386_DIR32    = 0x0006
	IMAGE_REL_I386_DIR32NB  = 0x0007
	IMAGE_REL_I386_SEG12    = 0x0009
	IMAGE_REL_I386_SECTION  = 0x000A
	IMAGE_REL_I386_SECREL   = 0x000B
	IMAGE_REL_I386_TOKEN    = 0x000C
	IMAGE_REL_I386_SECREL7  = 0x000D
	IMAGE_REL_I386_REL32    = 0x0014
)

type DOSHeader struct {
	Magic      uint16 // MZ
	LastSize   uint16 // image size mod 512, number of bytes on last page
	NumBlocks  uint16 // number of 512-byte pages in images
	NumRelocs  uint16 // count of relocation entries
	HeaderSize uint16 // size of header in paragraphs
	MinAlloc   uint16 // min required memory
	MaxAlloc   uint16 // max required memory
	SS         uint16 // stack seg offset in load module
	SP         uint16 // initial sp value
	Checksum   uint16 // one complement sum of all word in exe file
	IP         uint16 // initial ip value
	CS         uint16 // cs offset in load module
	RelocPos   uint16 // offset of first reloc item
	NoOverlay  uint16 // overlay number
	_          [4]uint16
	OEMID      uint16
	OEMInfo    uint16
	_          [10]uint16
	LFANew     uint32 // offset to pe header in windows
}

type ImportDescriptor struct {
	OriginalFirstThunk uint32
	TimeDateStamp      uint32
	ForwarderChain     uint32
	Name               uint32
	FirstThunk         uint32
}

type ImportDescription struct {
	ImportDescriptor
	DLLName string
	Symbols []Symbol
}

type ImportDescriptionList []ImportDescription

type ExportDirectory struct {
	Characteristics       uint32
	TimeDateStamp         uint32
	MajorVersion          uint16
	MinorVersion          uint16
	Name                  uint32
	Base                  uint32
	NumberOfFunctions     uint32
	NumberOfNames         uint32
	AddressOfFunctions    uint32
	AddressOfNames        uint32
	AddressOfNameOrdinals uint32
}

type Symbol struct {
	pe.Symbol
	DLLName          string
	ForwardedAddress uint64
	DLLNameRVA       uint64
	NameRVA          uint64
	OriginalThunkRVA uint64
	ThunkRVA         uint64
	IddIdx           int
	Auxillary        interface{}
}

type Section struct {
	*pe.Section
	Data []byte
}

type File struct {
	*pe.File
	ImageBase        uint64
	RawSizeOfHeaders uint64
	SizeOfHeaders    uint64
	SizeOfImage      uint64
	WordSize         uint64
	FileAlignment    uint64
	SectionAlignment uint64
	DOSHeader        DOSHeader
	DOSStub          []byte
	Sections         []*Section
	Strings          []string
	r                io.ReaderAt
}

type Word uint64

// these values are common across many exe files
var DOSHdr = DOSHeader{
	Magic:      0x5a4d,
	LastSize:   0x90,
	NumBlocks:  0x03,
	HeaderSize: 0x04,
	MaxAlloc:   0xffff,
	SP:         0xb8,
	RelocPos:   0x40,
	LFANew:     0x40 + uint32(len(DOSStub)),
}

var DOSStub = []byte{
	// push cs
	0x0E,
	// pop ds
	0x1F,
	// mov dx, 0xe
	0xBA, 0x0E, 0x00,
	// mov ah, 0x9
	0xB4, 0x09,
	// int 0x21
	0xCD, 0x21,
	// mov ax, 0x4c01
	0xB8, 0x01, 0x4C,
	// int 0x21
	0xCD, 0x21,
	// "This program cannot be run in DOS Mode"
	0x54, 0x68, 0x69, 0x73, 0x20, 0x70, 0x72, 0x6F,
	0x67, 0x72, 0x61, 0x6D, 0x20, 0x63, 0x61, 0x6E,
	0x6E, 0x6F, 0x74, 0x20, 0x62, 0x65, 0x20, 0x72,
	0x75, 0x6E, 0x20, 0x69, 0x6E, 0x20, 0x44, 0x4F,
	0x53, 0x20, 0x6D, 0x6F, 0x64, 0x65, 0x2E, 0x0D,
	0x0D, 0x0A, 0x24, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00,
}

func Open(name string) (*File, error) {
	p, err := pe.Open(name)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	return newFile(p, f)
}

func NewFile(r io.ReaderAt) (*File, error) {
	p, err := pe.NewFile(r)
	if err != nil {
		return nil, err
	}
	return newFile(p, r)
}

func newFile(p *pe.File, r io.ReaderAt) (*File, error) {
	f := &File{
		File:      p,
		DOSHeader: DOSHdr,
		DOSStub:   DOSStub,
	}
	for i := 4; i < len(f.StringTable); {
		str, err := f.StringTable.String(uint32(i))
		if err != nil {
			break
		}
		f.Strings = append(f.Strings, str)
		i += len(str) + 1
	}

	var dh DOSHeader
	sr := io.NewSectionReader(r, 0, math.MaxInt64)
	err := binary.Read(sr, binary.LittleEndian, &dh)
	if err == nil {
		f.DOSHeader = dh
	}
	stub := make([]byte, f.DOSHeader.LFANew-0x40)
	_, err = io.ReadAtLeast(sr, stub, len(stub))
	if err == nil {
		f.DOSStub = stub
	}

	for _, s := range f.File.Sections {
		p := &Section{}
		p.Section = s
		p.Data, err = s.Data()
		if err != nil {
			return nil, err
		}
		f.Sections = append(f.Sections, p)
	}

	switch f.Machine {
	case pe.IMAGE_FILE_MACHINE_AMD64:
		f.WordSize = 8
	case pe.IMAGE_FILE_MACHINE_I386:
		f.WordSize = 4
	default:
		return nil, fmt.Errorf("unsupported machine type %d", f.Machine)
	}

	pagesize := uint64(os.Getpagesize())
	switch h := f.OptionalHeader.(type) {
	case *pe.OptionalHeader32:
		f.ImageBase = uint64(h.ImageBase)
		f.FileAlignment = uint64(h.FileAlignment)
		f.SectionAlignment = uint64(h.SectionAlignment)
		f.SizeOfImage = uint64(h.SizeOfImage)
	case *pe.OptionalHeader64:
		f.ImageBase = uint64(h.ImageBase)
		f.FileAlignment = uint64(h.FileAlignment)
		f.SectionAlignment = uint64(h.SectionAlignment)
		f.SizeOfImage = uint64(h.SizeOfImage)
	default:
		f.ImageBase = uint64(0x400000)
		f.FileAlignment = 512
		f.SectionAlignment = pagesize
	}
	if f.FileAlignment == 0 {
		f.FileAlignment = 512
	}
	if f.SectionAlignment < pagesize {
		f.FileAlignment = f.SectionAlignment
	}

	f.updateImageHeaderSize()
	return f, nil
}

func (f *File) Section(name string) *Section {
	for _, s := range f.Sections {
		if s.Name == name {
			return s
		}
	}
	return nil
}

func (f *File) updateImageHeaderSize() {
	f.RawSizeOfHeaders, f.SizeOfHeaders, f.SizeOfImage = f.calcImageSizes()
	switch h := f.OptionalHeader.(type) {
	case *pe.OptionalHeader32:
		h.SizeOfHeaders = uint32(f.SizeOfHeaders)
		h.SizeOfImage = uint32(f.SizeOfImage)
	case *pe.OptionalHeader64:
		h.SizeOfHeaders = uint32(f.SizeOfHeaders)
		h.SizeOfImage = uint32(f.SizeOfImage)
	}
}

func (f *File) calcImageSizes() (rawSizeOfHeaders, sizeOfHeaders, sizeOfImage uint64) {
	rawSizeOfHeaders += uint64(reflect.TypeOf(f.DOSHeader).Size())
	rawSizeOfHeaders += uint64(len(f.DOSStub))
	// PE signature
	rawSizeOfHeaders += 4
	rawSizeOfHeaders += uint64(reflect.TypeOf(f.FileHeader).Size())

	switch h := f.OptionalHeader.(type) {
	case *pe.OptionalHeader32:
		sizeOfHeaders = uint64(h.SizeOfHeaders)
		rawSizeOfHeaders += uint64(reflect.TypeOf(*h).Size())
	case *pe.OptionalHeader64:
		sizeOfHeaders = uint64(h.SizeOfHeaders)
		rawSizeOfHeaders += uint64(reflect.TypeOf(*h).Size())
	}

	for _, s := range f.Sections {
		sh := f.sectionHeader32(&s.SectionHeader)
		rawSizeOfHeaders += uint64(reflect.TypeOf(sh).Size())
	}
	if rawSizeOfHeaders > sizeOfHeaders {
		sizeOfHeaders = rawSizeOfHeaders
	}
	sizeOfHeaders = ga.Align(sizeOfHeaders, f.FileAlignment)

	sizeOfImage = sizeOfHeaders
	for _, s := range f.Sections {
		sizeOfImage += uint64(ga.Align(uint64(s.VirtualSize), f.SectionAlignment))
	}
	sizeOfImage = ga.Align(sizeOfImage, f.FileAlignment)
	return
}

func (f *File) CreateSection(name string, va, size, flags uint64) (*Section, error) {
	a, _, _ := f.LookupVirtualAddress(va)
	b, _, _ := f.LookupVirtualAddress(va + size)
	c := f.Section(name)
	if va != 0 && (a != nil || b != nil || c != nil) {
		return nil, fmt.Errorf("section %s already exist", name)
	}

	var off uint32
	var dynva bool
	if va == 0 {
		dynva = true
	}
	for _, s := range f.Sections {
		nva := uint64(s.VirtualAddress + s.Size)
		if dynva && va < nva {
			va = nva
		}
		if off < s.Offset+s.Size {
			off = s.Offset + s.Size
		}
	}

	p := &Section{}
	p.Section = new(pe.Section)
	p.Name = name
	p.Data = make([]byte, size)
	p.Offset = uint32(off)
	p.Size = uint32(size)
	p.VirtualAddress = uint32(ga.Align(va, f.SectionAlignment))
	p.VirtualSize = uint32(ga.Align(uint64(p.Size), f.SectionAlignment))
	p.Characteristics = uint32(flags)
	f.Sections = append(f.Sections, p)
	f.NumberOfSections++
	sort.Slice(f.Sections, func(i, j int) bool {
		return f.Sections[i].VirtualAddress < f.Sections[j].VirtualAddress
	})

	oldSizeOfHeaders := f.SizeOfHeaders
	f.updateImageHeaderSize()
	for _, s := range f.Sections {
		s.Offset += uint32(f.SizeOfHeaders - oldSizeOfHeaders)
	}

	return p, nil
}

func (f *File) DeleteSection(name string) error {
	s := f.Section(name)
	if s == nil {
		return fmt.Errorf("section %s does not exist", name)
	}

	n := -1
	for i, p := range f.Sections {
		if p.Name == name {
			n = i
		}
		if p.Offset >= s.Offset {
			p.Offset -= s.Size
		}
	}

	copy(f.Sections[n:], f.Sections[n+1:])
	f.Sections = f.Sections[:len(f.Sections)-1]
	f.NumberOfSections--
	f.updateImageHeaderSize()
	return nil
}

func (f *File) ReadImportTable() (ImportDescriptionList, error) {
	idd := f.DataDirectory(pe.IMAGE_DIRECTORY_ENTRY_IMPORT)
	if idd == nil {
		return nil, fmt.Errorf("import table does not exist")
	}

	var dt ImportDescriptionList
	var d ImportDescription
	addr := uint64(idd.VirtualAddress)
	for {
		err := f.ReadVirtualAddress(addr, &d.ImportDescriptor)
		if err != nil {
			return nil, fmt.Errorf("failed to read import descriptor: %v", err)
		}
		if d.OriginalFirstThunk == 0 {
			break
		}
		f.ReadVirtualAddress(uint64(d.Name), &d.DLLName)
		dt = append(dt, d)

		addr += uint64(reflect.TypeOf(d.ImportDescriptor).Size())
	}

	for i := range dt {
		d := &dt[i]
		oftaddr := uint64(d.OriginalFirstThunk)
		ftaddr := uint64(d.FirstThunk)
		for {
			var na Word
			err := f.ReadVirtualAddress(oftaddr, &na)
			if err != nil {
				return nil, fmt.Errorf("failed to read from dll import %q", d.DLLName)
			}
			if na == 0 {
				break
			}
			s := Symbol{}
			s.DLLName = d.DLLName
			s.DLLNameRVA = uint64(d.Name)
			s.OriginalThunkRVA = oftaddr
			s.ThunkRVA = ftaddr

			var mask Word
			if f.WordSize == 8 {
				mask = 1 << 63
			} else {
				mask = 1 << 31
			}
			s.NameRVA = uint64(na)
			if na&mask == 0 {
				f.ReadVirtualAddress(s.NameRVA+2, &s.Name)
			} else {
				s.Name = fmt.Sprintf("%#x", na)
			}

			d.Symbols = append(d.Symbols, s)
			oftaddr += f.WordSize
			ftaddr += f.WordSize
		}
	}

	return dt, nil
}

func (f *File) ExportedSymbols() ([]Symbol, error) {
	var d ExportDirectory
	idd := f.DataDirectory(pe.IMAGE_DIRECTORY_ENTRY_EXPORT)
	if idd == nil {
		return nil, nil
	}
	err := f.ReadVirtualAddress(uint64(idd.VirtualAddress), &d)
	if err != nil {
		return nil, err
	}

	_, fp, _ := f.LookupVirtualAddress(uint64(d.AddressOfFunctions))
	_, od, _ := f.LookupVirtualAddress(uint64(d.AddressOfNameOrdinals))
	_, na, _ := f.LookupVirtualAddress(uint64(d.AddressOfNames))
	if fp == nil || od == nil {
		return nil, nil
	}
	no := uint64(0)

	var dllName string
	var syms []Symbol

	f.ReadVirtualAddress(uint64(d.Name), &d.Name)
	for i := uint32(0); i < d.NumberOfFunctions && len(od) >= 4; i, od = i+1, od[2:] {
		fn := binary.LittleEndian.Uint16(od) * 2
		if fn >= uint16(len(fp)) {
			continue
		}

		var name string
		var fwd uint64

		va := binary.LittleEndian.Uint32(fp[fn:])
		if idd.VirtualAddress <= va && va < idd.VirtualAddress+idd.Size {
			fwd = uint64(va)
		}

		an := uint64(d.AddressOfNames) + no
		sz := uint64(4)
		if uint64(len(na)) < sz {
			name = fmt.Sprintf("%s+%#x", dllName, va)
			an = uint64(va)
		} else {
			f.ReadVirtualAddress(f.readData(na, sz), &name)
			na = na[sz:]
			no += sz
		}
		p := Symbol{
			Symbol: pe.Symbol{
				Name: name,
			},
			DLLName:          dllName,
			ForwardedAddress: fwd,
			DLLNameRVA:       uint64(d.Name),
			NameRVA:          an,
			Auxillary:        idd,
		}
		syms = append(syms, p)
	}
	return syms, nil
}

func (f *File) readData(b []byte, size uint64) uint64 {
	var v uint64
	switch size {
	case 8:
		v = binary.LittleEndian.Uint64(b)
	case 4:
		v = uint64(binary.LittleEndian.Uint32(b))
	default:
		panic(fmt.Errorf("unsupported word size %d", f.WordSize))
	}
	return v
}

func (f *File) memoryAccessCrossesSection(va uint64, v interface{}) (*Section, bool) {
	p := reflect.TypeOf(v)
	if p.Kind() == reflect.Ptr {
		p = reflect.TypeOf(v).Elem()
	}
	n := reflect.TypeOf(p).Size()

	var t *Section
	for i := uint64(0); i < uint64(n); i++ {
		s, _, _ := f.LookupVirtualAddress(va + i)
		if t == nil {
			t = s
		}
		if t != s {
			return t, true
		}
	}
	return t, false
}

func (f *File) LookupVirtualAddress(va uint64) (*Section, []byte, int) {
	for _, s := range f.Sections {
		if uint64(s.VirtualAddress) <= va && va < uint64(s.VirtualAddress+s.VirtualSize) {
			off := int(va - uint64(s.VirtualAddress))
			if off >= len(s.Data) {
				return s, nil, off
			}
			return s, s.Data[off:], off
		}
	}
	return nil, nil, 0
}

func (f *File) ReadVirtualAddress(va uint64, v interface{}) error {
	var b []byte
	switch v := v.(type) {
	case nil:
		return nil
	case *uint8:
		b = make([]byte, 1)
	case *uint16:
		b = make([]byte, 2)
	case *uint32:
		b = make([]byte, 4)
	case *uint64:
		b = make([]byte, 8)
	case *Word:
		switch f.WordSize {
		case 8:
			var w uint64
			err := f.ReadVirtualAddress(va, &w)
			*v = Word(w)
			return err
		case 4:
			var w uint32
			err := f.ReadVirtualAddress(va, &w)
			*v = Word(w)
			return err
		default:
			panic(fmt.Errorf("unsupported word size %d", f.WordSize))
		}
	case *string:
		*v = ""
		for i := 0; ; i++ {
			_, p, _ := f.LookupVirtualAddress(va + uint64(i))
			if len(p) == 0 {
				return fmt.Errorf("invalid read of unmapped address %#x", va+uint64(i))
			}
			if p[0] == 0 {
				break
			}
			*v += string(p[0])
		}
		return nil
	case *ExportDirectory:
		n := reflect.TypeOf(ExportDirectory{}).Size()
		b = make([]byte, int(n))
	case *ImportDescriptor:
		n := reflect.TypeOf(ImportDescriptor{}).Size()
		b = make([]byte, int(n))
	default:
		panic(fmt.Errorf("unsupported type %T", v))
	}

	for i := range b {
		_, p, _ := f.LookupVirtualAddress(va + uint64(i))
		if p == nil {
			return fmt.Errorf("invalid read of unmapped address %#x", va+uint64(i))
		}
		b[i] = p[0]
	}

	r := bytes.NewReader(b)
	return binary.Read(r, binary.LittleEndian, v)
}

func (f *File) ReadVirtualAddressBound(va uint64, v interface{}) (*Section, error) {
	s, crossed := f.memoryAccessCrossesSection(va, v)
	if crossed {
		return nil, fmt.Errorf("data read of address %#x crosses section", va)
	}
	err := f.ReadVirtualAddress(va, v)
	return s, err
}

func (f *File) WriteVirtualAddress(va uint64, v interface{}) error {
	var b []byte
	switch v := v.(type) {
	case nil:
		return nil
	case []byte:
		b = v
	case Word:
		switch f.WordSize {
		case 8:
			b = make([]byte, 8)
			binary.LittleEndian.PutUint64(b, uint64(v))
		case 4:
			b = make([]byte, 4)
			binary.LittleEndian.PutUint32(b, uint32(v))
		default:
			panic(fmt.Errorf("unsupported word size %d", f.WordSize))
		}
	case string:
		b = []byte(v)
		b = append(b, 0)
	default:
		w := new(bytes.Buffer)
		err := binary.Write(w, binary.LittleEndian, v)
		if err != nil {
			return err
		}
		b = w.Bytes()
	}
	for i := range b {
		_, p, _ := f.LookupVirtualAddress(va + uint64(i))
		if len(p) == 0 {
			return fmt.Errorf("invalid read of unmapped address %#x", va+uint64(i))
		}
		p[0] = b[i]
	}
	return nil
}

func (f *File) WriteVirtualAddressBound(va uint64, v interface{}) (*Section, error) {
	s, crossed := f.memoryAccessCrossesSection(va, v)
	if crossed {
		return nil, fmt.Errorf("data write of address %#x crosses section", va)
	}
	err := f.WriteVirtualAddress(va, v)
	return s, err
}

func (f *File) CreateDataDirectory(index int) *pe.DataDirectory {
	defer f.updateImageHeaderSize()
	if f.OptionalHeader == nil {
		if f.WordSize == 8 {
			f.OptionalHeader = &pe.OptionalHeader64{}
		} else {
			f.OptionalHeader = &pe.OptionalHeader32{}
		}
	}

	i := uint32(index)
	switch h := f.OptionalHeader.(type) {
	case *pe.OptionalHeader64:
		if h.NumberOfRvaAndSizes < i {
			h.NumberOfRvaAndSizes = i + 1
		}
		return &h.DataDirectory[i]
	case *pe.OptionalHeader32:
		if h.NumberOfRvaAndSizes < i {
			h.NumberOfRvaAndSizes = i + 1
		}
		return &h.DataDirectory[i]
	}
	return nil
}

func (f *File) DataDirectory(index int) *pe.DataDirectory {
	var dirlen uint32
	var idd *pe.DataDirectory
	switch h := f.OptionalHeader.(type) {
	case *pe.OptionalHeader64:
		dirlen = h.NumberOfRvaAndSizes
		idd = &h.DataDirectory[index]
	case *pe.OptionalHeader32:
		dirlen = h.NumberOfRvaAndSizes
		idd = &h.DataDirectory[index]
	}
	if dirlen < uint32(index)+1 {
		return nil
	}
	return idd
}

func (f *File) sectionHeader32(s *pe.SectionHeader) pe.SectionHeader32 {
	h := pe.SectionHeader32{
		VirtualSize:          s.VirtualSize,
		VirtualAddress:       s.VirtualAddress,
		SizeOfRawData:        s.Size,
		PointerToRawData:     s.Offset,
		PointerToRelocations: s.PointerToRelocations,
		PointerToLineNumbers: s.PointerToLineNumbers,
		NumberOfRelocations:  s.NumberOfRelocations,
		NumberOfLineNumbers:  s.NumberOfLineNumbers,
		Characteristics:      s.Characteristics,
	}

	name := s.Name
	if len(s.Name) > len(h.Name) {
		n := bytes.Index(f.StringTable, []byte(s.Name))
		if n >= 0 {
			name = fmt.Sprintf("/%d", n+4)
		}
	}
	copy(h.Name[:], name[:])
	return h
}

func (f *File) RedirectImportTable(dt []ImportDescription) error {
	var (
		sectsize  uint64 = 512 * 1024
		strsize   uint64 = 128
		symsize   uint64 = 16 * 1024
		thunksize uint64 = 1024
		thunkoff  uint64 = 4 * 1024
		dlloff    uint64 = 32 * 1024
		symoff    uint64 = 64 * 1024
	)

	size := uint64(256 * 1024)
	for _, d := range dt {
		for _, y := range d.Symbols {
			if l := uint64(len(y.Name)) + 1; strsize < l {
				strsize = ga.Align(l, 128)
			}
			size += strsize
		}
		if thunksize < uint64(len(d.Symbols)+1) {
			thunksize = uint64(len(d.Symbols) + 1)
			thunksize = ga.Align(thunksize, 1024)
		}
		size += thunksize
		size += strsize
	}
	if sectsize < size {
		sectsize = ga.Align(size, f.SectionAlignment)
	}

	s, err := f.CreateSection(".imptbl", 0, sectsize, IMAGE_SCN_CNT_INITIALIZED_DATA|IMAGE_SCN_MEM_READ|IMAGE_SCN_MEM_WRITE)
	if err != nil {
		return err
	}

	idd := f.CreateDataDirectory(pe.IMAGE_DIRECTORY_ENTRY_IMPORT)
	idd.VirtualAddress = s.VirtualAddress
	idd.Size = uint32(len(dt)+1) * uint32(reflect.TypeOf(ImportDescriptor{}).Size())

	base := uint64(idd.VirtualAddress)
	iaddr := base
	for i, d := range dt {
		daddr := base + dlloff + strsize*uint64(i)
		f.WriteVirtualAddress(daddr, d.DLLName)

		p := d.ImportDescriptor
		p.OriginalFirstThunk = uint32(base) + uint32(thunkoff) + uint32(thunksize)*uint32(i)
		p.Name = uint32(daddr)
		f.WriteVirtualAddress(iaddr, &p)
		iaddr += uint64(reflect.TypeOf(ImportDescriptor{}).Size())

		yaddr := base + symoff + symsize*uint64(i)
		for j, y := range d.Symbols {
			taddr := uint64(p.OriginalFirstThunk) + uint64(j)*f.WordSize
			saddr := yaddr + strsize*uint64(j)
			f.WriteVirtualAddress(taddr, Word(saddr))
			f.WriteVirtualAddress(saddr+2, y.Name)
		}
	}

	return nil
}

func (f *File) SwapSectionName(a, b string) error {
	p := f.Section(a)
	q := f.Section(b)
	if p == nil {
		return fmt.Errorf("section %q does not exist", a)
	}
	if q == nil {
		return fmt.Errorf("section %q does not exist", b)
	}
	p.Name, q.Name = q.Name, p.Name
	return nil
}

func (f *File) FindStrings() []debug.String {
	var tab []debug.String
	for _, s := range f.Sections {
		st := debug.FindStringsBuffer(s.Data, 4, 256, true)
		for i := range st {
			st[i].Name = s.Name
		}
		tab = append(tab, st...)
	}
	return tab
}

func (dt *ImportDescriptionList) RenameDLL(from, to string) error {
	for i := range *dt {
		if strings.ToLower((*dt)[i].DLLName) == strings.ToLower(from) {
			(*dt)[i].DLLName = to
			return nil
		}
	}

	return fmt.Errorf("dll import %q does not exist", from)
}

func (dt *ImportDescriptionList) RenameSymbol(from, to string) error {
	found := false
	for _, d := range *dt {
		for _, y := range d.Symbols {
			if y.Name == from {
				y.Name = to
				found = true
			}
		}
	}
	if !found {
		return fmt.Errorf("import symbol %q does not exist")
	}
	return nil
}

func (dt *ImportDescriptionList) RedirectSymbol(fromDLL, toDLL, fromSym, toSym string) error {
	var i, j int
loop:
	for i = range *dt {
		if strings.ToLower((*dt)[i].DLLName) != strings.ToLower(fromDLL) {
			continue
		}
		for j = range (*dt)[i].Symbols {
			if (*dt)[i].Symbols[j].Name == fromSym {
				break loop
			}
		}
	}
	if i == len(*dt) {
		return fmt.Errorf("import symbol %q does not exist")
	}

	if j < len((*dt)[i].Symbols) {
		p := ImportDescription{}
		p.DLLName = (*dt)[i].DLLName
		p.ForwarderChain = (*dt)[i].ForwarderChain
		p.Name = (*dt)[i].Name
		p.OriginalFirstThunk = uint32((*dt)[i].Symbols[j+1].OriginalThunkRVA)
		p.FirstThunk = uint32((*dt)[i].Symbols[j+1].ThunkRVA)
		p.Symbols = append(p.Symbols, (*dt)[i].Symbols[j+1:]...)
		*dt = append(*dt, p)
	}

	y := Symbol{}
	y.Name = toSym
	p := ImportDescription{}
	p.DLLName = toDLL
	p.OriginalFirstThunk = uint32((*dt)[i].Symbols[j].OriginalThunkRVA)
	p.FirstThunk = uint32((*dt)[i].Symbols[j].ThunkRVA)
	p.Symbols = append(p.Symbols, y)
	*dt = append(*dt, p)

	(*dt)[i].Symbols = (*dt)[i].Symbols[:j]

	return nil
}

func (dt *ImportDescriptionList) LookupSymbol(name string) *Symbol {
	for _, d := range *dt {
		for _, y := range d.Symbols {
			if y.Name == name {
				return &y
			}
		}
	}
	return nil
}

func Format(f *File, w io.Writer) error {
	b := bufio.NewWriter(w)

	binary.Write(b, binary.LittleEndian, &f.DOSHeader)
	b.Write(f.DOSStub)

	peSig := [...]byte{'P', 'E', 0x00, 0x00}
	binary.Write(b, binary.LittleEndian, peSig)
	binary.Write(b, binary.LittleEndian, &f.FileHeader)

	switch h := f.OptionalHeader.(type) {
	case *pe.OptionalHeader32:
		binary.Write(b, binary.LittleEndian, h)
	case *pe.OptionalHeader64:
		binary.Write(b, binary.LittleEndian, h)
	}

	for _, s := range f.Sections {
		sh := f.sectionHeader32(&s.SectionHeader)
		binary.Write(b, binary.LittleEndian, &sh)
	}
	pad := make([]byte, f.SizeOfHeaders-f.RawSizeOfHeaders)
	b.Write(pad)

	for _, s := range f.Sections {
		b.Write(s.Data)
	}

	return b.Flush()
}

type MachineType int

func (m MachineType) String() string {
	switch m {
	case pe.IMAGE_FILE_MACHINE_WCEMIPSV2:
		return "IMAGE_FILE_MACHINE_WCEMIPSV2"
	case pe.IMAGE_FILE_MACHINE_THUMB:
		return "IMAGE_FILE_MACHINE_THUMB"
	case pe.IMAGE_FILE_MACHINE_SH5:
		return "IMAGE_FILE_MACHINE_SH5"
	case pe.IMAGE_FILE_MACHINE_SH4:
		return "IMAGE_FILE_MACHINE_SH4"
	case pe.IMAGE_FILE_MACHINE_SH3DSP:
		return "IMAGE_FILE_MACHINE_SH3DSP"
	case pe.IMAGE_FILE_MACHINE_SH3:
		return "IMAGE_FILE_MACHINE_SH3"
	case pe.IMAGE_FILE_MACHINE_R4000:
		return "IMAGE_FILE_MACHINE_R4000"
	case pe.IMAGE_FILE_MACHINE_POWERPCFP:
		return "IMAGE_FILE_MACHINE_POWERPCFP"
	case pe.IMAGE_FILE_MACHINE_POWERPC:
		return "IMAGE_FILE_MACHINE_POWERPC"
	case pe.IMAGE_FILE_MACHINE_MIPSFPU16:
		return "IMAGE_FILE_MACHINE_MIPSFPU16"
	case pe.IMAGE_FILE_MACHINE_MIPSFPU:
		return "IMAGE_FILE_MACHINE_MIPSFPU"
	case pe.IMAGE_FILE_MACHINE_MIPS16:
		return "IMAGE_FILE_MACHINE_MIPS16"
	case pe.IMAGE_FILE_MACHINE_M32R:
		return "IMAGE_FILE_MACHINE_M32R"
	case pe.IMAGE_FILE_MACHINE_IA64:
		return "IMAGE_FILE_MACHINE_IA64"
	case pe.IMAGE_FILE_MACHINE_I386:
		return "IMAGE_FILE_MACHINE_I386"
	case pe.IMAGE_FILE_MACHINE_EBC:
		return "IMAGE_FILE_MACHINE_EBC"
	case pe.IMAGE_FILE_MACHINE_ARM64:
		return "IMAGE_FILE_MACHINE_ARM64"
	case pe.IMAGE_FILE_MACHINE_ARMNT:
		return "IMAGE_FILE_MACHINE_ARMNT"
	case pe.IMAGE_FILE_MACHINE_ARM:
		return "IMAGE_FILE_MACHINE_ARM"
	case pe.IMAGE_FILE_MACHINE_AMD64:
		return "IMAGE_FILE_MACHINE_AMD64"
	case pe.IMAGE_FILE_MACHINE_AM33:
		return "IMAGE_FILE_MACHINE_AM33"
	case pe.IMAGE_FILE_MACHINE_UNKNOWN:
		return "IMAGE_FILE_MACHINE_UNKNOWN"
	default:
		return fmt.Sprintf("MachineType(%d)", m)
	}
}
