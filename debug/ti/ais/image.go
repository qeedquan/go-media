// http://www.ti.com.cn/cn/lit/an/sprabb1c/sprabb1c.pdf

package ais

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/qeedquan/go-media/debug/elfutil"
	"gopkg.in/ini.v1"
)

type Image struct {
	Cmds []Cmd
}

type Cmd struct {
	Op uint32

	Addr    uint32
	Size    uint32
	Type    uint32
	Pattern uint32
	Data    []byte

	Func uint32
	Args []uint32
}

func Open(name string) (*Image, error) {
	r, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return NewFile(r)
}

func NewFile(r io.Reader) (*Image, error) {
	b := bufio.NewReader(r)

	var sig uint32
	binary.Read(b, binary.LittleEndian, &sig)
	if sig != MAGIC {
		return nil, fmt.Errorf("ais: invalid header signature %#x", sig)
	}

	var (
		op   uint32
		buf  [8]uint32
		cmds []Cmd
	)
	for {
		err := binary.Read(b, binary.LittleEndian, &op)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("ais: %v", err)
		}

		switch op {
		case FUNCTION_EXEC:
			binary.Read(b, binary.LittleEndian, &buf[0])
			narg := (buf[0] >> 16) & 0xffff
			fun := buf[0] & 0xffff
			args := make([]uint32, narg)
			binary.Read(b, binary.LittleEndian, args)
			cmds = append(cmds, Cmd{
				Op:   op,
				Func: fun,
				Args: args,
			})

		case SECTION_LOAD:
			binary.Read(b, binary.LittleEndian, buf[:2])
			data := make([]byte, buf[1])
			binary.Read(b, binary.LittleEndian, data)

			cmds = append(cmds, Cmd{
				Op:   op,
				Addr: buf[0],
				Size: buf[1],
				Data: data,
			})

		case SECTION_FILL:
			binary.Read(b, binary.LittleEndian, buf[:4])
			cmds = append(cmds, Cmd{
				Op:      op,
				Addr:    buf[0],
				Size:    buf[1],
				Type:    buf[2],
				Pattern: buf[3],
			})

		case JUMP, JUMP_CLOSE:
			binary.Read(b, binary.LittleEndian, &buf[0])
			cmds = append(cmds, Cmd{Op: op, Addr: buf[0]})

		default:
			cmds = append(cmds, Cmd{Op: op})
		}
	}

	return &Image{
		Cmds: cmds,
	}, nil
}

func Format(m *Image, w io.Writer) error {
	b := bufio.NewWriter(w)

	sig := uint32(MAGIC)
	binary.Write(b, binary.LittleEndian, &sig)
	for _, c := range m.Cmds {
		binary.Write(b, binary.LittleEndian, c.Op)
		switch c.Op {
		case FUNCTION_EXEC:
			fxnarg := c.Func | uint32(len(c.Args))<<16
			binary.Write(b, binary.LittleEndian, fxnarg)
			for _, a := range c.Args {
				binary.Write(b, binary.LittleEndian, &a)
			}

		case SECTION_LOAD:
			size := uint32(len(c.Data))
			binary.Write(b, binary.LittleEndian, &c.Addr)
			binary.Write(b, binary.LittleEndian, &size)
			binary.Write(b, binary.LittleEndian, c.Data)

		case SECTION_FILL:
			binary.Write(b, binary.LittleEndian, &c.Addr)
			binary.Write(b, binary.LittleEndian, &c.Size)
			binary.Write(b, binary.LittleEndian, &c.Type)
			binary.Write(b, binary.LittleEndian, &c.Pattern)

		case JUMP, JUMP_CLOSE:
			binary.Write(b, binary.LittleEndian, &c.Addr)

		default:
			return fmt.Errorf("ais: unknown op(%#x)", c.Op)
		}
	}

	err := b.Flush()
	if err != nil {
		return fmt.Errorf("ais: %v", err)
	}
	return nil
}

func Disasm(c *Cmd) string {
	var s string
	switch c.Op {
	case SECTION_LOAD:
		s = fmt.Sprintf("section_load(addr = %#x-%#x, size = %#x)",
			c.Addr, c.Addr+c.Size-1, c.Size)

	case SECTION_FILL:
		s = fmt.Sprintf("section_fill(addr = %#x-%#x, size = %#x, type = %d, pattern = %#x)",
			c.Addr, c.Addr+c.Size-1, c.Size, c.Type, c.Pattern)

	case FUNCTION_EXEC:
		s = fmt.Sprintf("function_exec(code = %#x, args = [", c.Func)
		for i, a := range c.Args {
			s += fmt.Sprintf("%#x", a)
			if i+1 < len(c.Args) {
				s += fmt.Sprintf(", ")
			}
		}
		s += fmt.Sprintf("])")

	case JUMP:
		s = fmt.Sprintf("jump(addr = %#x)", c.Addr)

	case JUMP_CLOSE:
		s = fmt.Sprintf("jump_close(addr = %#x)", c.Addr)

	case ENABLE_CRC:
		s = fmt.Sprintf("enable_crc()")

	case DISABLE_CRC:
		s = fmt.Sprintf("disable_crc()")

	default:
		s = fmt.Sprintf("unknown(%#x)", c.Op)
	}
	return s
}

func (m *Image) AddINI(r io.ReadCloser) error {
	tab := []struct {
		funct   uint32
		section string
		key     string
		value   uint32
	}{
		{funct: CFG_PLL, section: "PLL0CONFIG", key: "PLL0CFG0"},
		{funct: CFG_PLL, section: "PLL0CONFIG", key: "PLL0CFG1"},

		{funct: CFG_PLL_CLK, section: "PLLANDCLOCKCONFIG", key: "PLL0CFG0"},
		{funct: CFG_PLL_CLK, section: "PLLANDCLOCKCONFIG", key: "PLL0CFG1"},
		{funct: CFG_PLL_CLK, section: "PLLANDCLOCKCONFIG", key: "PERIPHCLKCFG"},

		{funct: CFG_EMIFA_SDRAM, section: "EMIF3DDR", key: "PLL1CFG0"},
		{funct: CFG_EMIFA_SDRAM, section: "EMIF3DDR", key: "PLL1CFG1"},
		{funct: CFG_EMIFA_SDRAM, section: "EMIF3DDR", key: "DDRPHYC1R"},
		{funct: CFG_EMIFA_SDRAM, section: "EMIF3DDR", key: "SDCR"},
		{funct: CFG_EMIFA_SDRAM, section: "EMIF3DDR", key: "SDTIMR"},
		{funct: CFG_EMIFA_SDRAM, section: "EMIF3DDR", key: "SDTIMR2"},
		{funct: CFG_EMIFA_SDRAM, section: "EMIF3DDR", key: "SDRCR"},
		{funct: CFG_EMIFA_SDRAM, section: "EMIF3DDR", key: "CLK2XSRC"},

		{funct: CFG_EMIFA_CE_SPACE, section: "EMIF25ASYNC", key: "A1CR"},
		{funct: CFG_EMIFA_CE_SPACE, section: "EMIF25ASYNC", key: "A2CR"},
		{funct: CFG_EMIFA_CE_SPACE, section: "EMIF25ASYNC", key: "A3CR"},
		{funct: CFG_EMIFA_CE_SPACE, section: "EMIF25ASYNC", key: "A4CR"},
		{funct: CFG_EMIFA_CE_SPACE, section: "EMIF25ASYNC", key: "NANDFCR"},
	}

	var nargs [16]int
	var sect [16]string
	arity := map[int]int{
		CFG_PLL:            2,
		CFG_PLL_CLK:        3,
		CFG_EMIFA_SDRAM:    8,
		CFG_EMIFA_CE_SPACE: 5,
	}

	f, err := ini.Load(r)
	if err != nil {
		return err
	}

	for i := range tab {
		p := &tab[i]

		s, err := f.GetSection(p.section)
		if err != nil {
			continue
		}

		k, err := s.GetKey(p.key)
		if err != nil {
			continue
		}

		v, err := k.Uint64()
		if err != nil {
			continue
		}

		p.value = uint32(v)
		nargs[p.funct]++
		sect[p.funct] = p.section
	}

	for k, v := range arity {
		if nargs[k] > 0 && nargs[k] != v {
			return fmt.Errorf("ais: section %q has insufficient number of arguments, expected %d but got %d", sect[k], v, nargs[k])
		}
	}

	var args []uint32
	for i, p := range tab {
		if nargs[p.funct] == 0 {
			continue
		}

		args = append(args, p.value)
		if i+1 >= len(tab) || p.funct != tab[i+1].funct {
			m.Cmds = append(m.Cmds, Cmd{Op: FUNCTION_EXEC, Func: p.funct, Args: args})
			args = make([]uint32, 0)
		}
	}

	return nil
}

func (m *Image) AddELF(r io.ReaderAt) (*elfutil.File, error) {
	f, err := elfutil.NewFile(r)
	if err != nil {
		return nil, err
	}

	for _, p := range f.Progs {
		if p.Filesz == 0 {
			continue
		}
		m.Cmds = append(m.Cmds, Cmd{Op: SECTION_LOAD, Addr: uint32(p.Paddr), Data: append([]byte{}, p.Data...)})
	}

	return f, nil
}
