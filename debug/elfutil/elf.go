package elfutil

import (
	"bytes"
	"debug/elf"
	"encoding/binary"
	"errors"
	"io"
	"math"
	"os"

	"github.com/qeedquan/go-media/debug"
)

type File struct {
	*elf.File
	Phoff     uint64
	Shoff     uint64
	Flags     uint32
	Phentsize uint64
	Shentsize uint64
	Sections  []*Section
	Progs     []*Prog
}

type Section struct {
	*elf.Section
	Data []byte
}

type Prog struct {
	*elf.Prog
	Data []byte
}

var (
	ErrInvalidClass = errors.New("elf: invalid class")
)

func Open(name string) (*File, error) {
	r, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	ef, err := NewFile(r)
	if err != nil {
		r.Close()
	}
	return ef, err
}

func NewFile(r io.ReaderAt) (*File, error) {
	e, err := elf.NewFile(r)
	if err != nil {
		return nil, err
	}
	return newFile(e, r)
}

func newFile(e *elf.File, r io.ReaderAt) (*File, error) {
	f := &File{
		File: e,
	}
	for _, s := range e.Sections {
		var (
			b   []byte
			err error
		)
		if s.Type != elf.SHT_NOBITS {
			b, err = io.ReadAll(s.Open())
			if err != nil {
				return nil, err
			}
		}

		f.Sections = append(f.Sections, &Section{
			Section: s,
			Data:    b,
		})
	}
	for _, p := range e.Progs {
		b, err := io.ReadAll(p.Open())
		if err != nil {
			return nil, err
		}
		f.Progs = append(f.Progs, &Prog{
			Prog: p,
			Data: b,
		})
	}

	var err error
	sr := io.NewSectionReader(r, 0, math.MaxInt32)
	switch f.Class {
	case elf.ELFCLASS32:
		var h elf.Header32
		err = binary.Read(sr, f.ByteOrder, &h)
		f.Phoff = uint64(h.Phoff)
		f.Shoff = uint64(h.Shoff)
		f.Phentsize = uint64(h.Phentsize)
		f.Shentsize = uint64(h.Shentsize)
		f.Flags = h.Flags
	case elf.ELFCLASS64:
		var h elf.Header64
		err = binary.Read(sr, f.ByteOrder, &h)
		f.Phoff = h.Phoff
		f.Shoff = h.Shoff
		f.Phentsize = uint64(h.Phentsize)
		f.Shentsize = uint64(h.Shentsize)
		f.Flags = h.Flags
	default:
		err = ErrInvalidClass
	}
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (f *File) Section(name string) (int, *Section) {
	for i, s := range f.Sections {
		if name == s.Name {
			return i, s
		}
	}
	return -1, nil
}

func Format(f *File, w io.Writer) error {
	b := new(debug.Patchset)

	shstrndx, strtab := f.Section(".shstrtab")
	if shstrndx < 0 {
		shstrndx = 0
	}

	switch f.Class {
	case elf.ELFCLASS32:
		binary.Write(b, f.ByteOrder, &elf.Header32{
			Ident: [elf.EI_NIDENT]byte{
				0x7f, 'E', 'L', 'F',
				byte(f.Class), byte(f.Data), byte(f.Version),
				byte(f.OSABI), f.ABIVersion,
			},
			Type:      uint16(f.Type),
			Machine:   uint16(f.Machine),
			Version:   uint32(f.Version),
			Entry:     uint32(f.Entry),
			Phoff:     uint32(f.Phoff),
			Shoff:     uint32(f.Shoff),
			Ehsize:    0x34,
			Flags:     0,
			Phentsize: uint16(f.Phentsize),
			Phnum:     uint16(len(f.Progs)),
			Shentsize: uint16(f.Shentsize),
			Shnum:     uint16(len(f.Sections)),
			Shstrndx:  uint16(shstrndx),
		})

	case elf.ELFCLASS64:
		binary.Write(b, f.ByteOrder, &elf.Header64{
			Ident: [elf.EI_NIDENT]byte{
				0x7f, 'E', 'L', 'F',
				byte(f.Class), byte(f.Data), byte(f.Version),
				byte(f.OSABI), f.ABIVersion,
			},
			Type:      uint16(f.Type),
			Machine:   uint16(f.Machine),
			Version:   uint32(f.Version),
			Entry:     f.Entry,
			Phoff:     f.Phoff,
			Shoff:     f.Shoff,
			Ehsize:    0x40,
			Flags:     0,
			Phentsize: uint16(f.Phentsize),
			Phnum:     uint16(len(f.Progs)),
			Shentsize: uint16(f.Shentsize),
			Shnum:     uint16(len(f.Sections)),
			Shstrndx:  uint16(shstrndx),
		})
	}

	for i, p := range f.Progs {
		b.Seek(int64(f.Phoff+uint64(i)*f.Phentsize), io.SeekStart)
		switch f.Class {
		case elf.ELFCLASS32:
			binary.Write(b, f.ByteOrder, &elf.Prog32{
				Type:   uint32(p.Type),
				Off:    uint32(p.Off),
				Vaddr:  uint32(p.Vaddr),
				Paddr:  uint32(p.Paddr),
				Filesz: uint32(p.Filesz),
				Memsz:  uint32(p.Memsz),
				Flags:  uint32(p.Flags),
				Align:  uint32(p.Align),
			})
		case elf.ELFCLASS64:
			binary.Write(b, f.ByteOrder, &elf.Prog64{
				Type:   uint32(p.Type),
				Off:    uint64(p.Off),
				Vaddr:  uint64(p.Vaddr),
				Paddr:  uint64(p.Paddr),
				Filesz: uint64(p.Filesz),
				Memsz:  uint64(p.Memsz),
				Flags:  uint32(p.Flags),
				Align:  uint64(p.Align),
			})
		}
	}

	for _, p := range f.Progs {
		b.Seek(int64(p.Off), io.SeekStart)
		b.Write(p.Data)
	}

	for i, s := range f.Sections {
		nameoff := 0
		if strtab != nil {
			off := bytes.Index(strtab.Data, append([]byte(s.Name), 0))
			if off > 0 {
				nameoff = off
			}
		}
		b.Seek(int64(f.Shoff+uint64(i)*f.Shentsize), io.SeekStart)
		switch f.Class {
		case elf.ELFCLASS32:
			binary.Write(b, f.ByteOrder, &elf.Section32{
				Name:      uint32(nameoff),
				Type:      uint32(s.Type),
				Flags:     uint32(s.Flags),
				Addr:      uint32(s.Addr),
				Off:       uint32(s.Offset),
				Size:      uint32(s.Size),
				Link:      uint32(s.Link),
				Info:      uint32(s.Info),
				Addralign: uint32(s.Addralign),
				Entsize:   uint32(s.Entsize),
			})
		case elf.ELFCLASS64:
			binary.Write(b, f.ByteOrder, &elf.Section64{
				Name:      uint32(nameoff),
				Type:      uint32(s.Type),
				Flags:     uint64(s.Flags),
				Addr:      uint64(s.Addr),
				Off:       uint64(s.Offset),
				Size:      uint64(s.Size),
				Link:      uint32(s.Link),
				Info:      uint32(s.Info),
				Addralign: uint64(s.Addralign),
				Entsize:   uint64(s.Entsize),
			})
		}
	}

	for _, s := range f.Sections {
		b.Seek(int64(s.Offset), io.SeekStart)
		b.Write(s.Data)
	}

	b.Merge()

	_, err := w.Write(b.Data)
	return err
}
