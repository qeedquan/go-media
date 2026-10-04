package pemapfile

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
)

type File struct {
	Name      string
	Timestamp string
	Checksum  string
	BaseAddr  uint64
	Sections  []Section
	Symbols   []Symbol
	Segments  []Segment
}

type Segment struct {
	Size uint64
}

type Section struct {
	Segment uint64
	Offset  uint64
	Size    uint64
	Name    string
	Class   string
}

type Symbol struct {
	Section int
	Segment uint64
	Offset  uint64
	Size    uint64
	Name    string
	Addr    uint64
	Flags   uint64
	File    string
}

func streq(s, t []string) bool {
	for i := range s {
		if s[i] != t[i] {
			return false
		}
	}
	return true
}

func Open(name string) (*File, error) {
	const (
		NAME = iota
		TIMESTAMP
		LOAD_ADDRESS
		MATCH
		SECTION
		SYMBOL
	)

	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var t [8]string
	m := &File{}
	sc := bufio.NewScanner(f)
	st := NAME
	for sc.Scan() {
		ln := strings.TrimSpace(sc.Text())
		if ln == "" {
			if st == SECTION {
				st = MATCH
			}

			continue
		}

		switch st {
		case NAME:
			fmt.Sscanf(ln, "%s", &m.Name)

			st = TIMESTAMP

		case TIMESTAMP:
			fmt.Sscanf(ln, "Timestamp is %s", &m.Checksum)
			sp := strings.IndexRune(ln, '(')
			ep := strings.IndexRune(ln, ')')
			if sp >= 0 && ep >= 0 {
				m.Timestamp = ln[sp+1 : ep]
			}

			st = LOAD_ADDRESS

		case LOAD_ADDRESS:
			fmt.Sscanf(ln, "Preferred load address is %x", &m.BaseAddr)
			st = MATCH

		case MATCH:
			n, _ := fmt.Sscanf(ln, "%s %s %s %s", &t[0], &t[1], &t[2], &t[3])
			if n == 4 && streq(t[:4], []string{"Start", "Length", "Name", "Class"}) {
				st = SECTION
				continue
			}

			n, _ = fmt.Sscanf(ln, "%s %s %s %s %s %s", &t[0], &t[1], &t[2], &t[3], &t[4], &t[5])

			if n == 6 && streq(t[:6], []string{"Address", "Publics", "by", "Value", "Rva+Base", "Lib:Object"}) {
				st = SYMBOL
				continue
			}

		case SECTION:
			var p Section
			n, _ := fmt.Sscanf(ln, "%x:%x %xH %s %s", &p.Segment, &p.Offset, &p.Size, &p.Name, &p.Class)
			if n == 5 {
				m.Sections = append(m.Sections, p)
			}

		case SYMBOL:
			var p Symbol
			fmt.Sscanf(ln, "%x:%x %s %x", &p.Segment, &p.Offset, &p.Name, &p.Addr)

			t := strings.Split(ln, " ")
			if len(t) == 0 {
				continue
			}

			p.File = t[len(t)-1]
			for i := len(t) - 2; i > 2; i-- {
				t[i] = strings.TrimSpace(t[i])
				switch t[i] {
				case "f":
					p.Flags |= 0x1
				case "i":
					p.Flags |= 0x2
				}
			}
			m.Symbols = append(m.Symbols, p)
		}
	}

	sort.Slice(m.Sections, func(i, j int) bool {
		x := &m.Sections[i]
		y := &m.Sections[j]
		if x.Segment == y.Segment {
			return x.Offset < y.Offset
		}
		return x.Segment < y.Segment
	})

	sort.Slice(m.Symbols, func(i, j int) bool {
		x := &m.Symbols[i]
		y := &m.Symbols[j]
		if x.Segment == y.Segment {
			return x.Offset < y.Offset
		}
		return x.Segment < y.Segment
	})

	if l := len(m.Sections); l > 0 {
		m.Segments = make([]Segment, m.Sections[l-1].Segment+1)

		seg := uint64(math.MaxUint64)
		for i := l - 1; i >= 0; i-- {
			p := &m.Sections[i]
			if seg != p.Segment {
				seg = p.Segment
				m.Segments[seg].Size = p.Offset + p.Size
			}
		}

		var s []Segment
		for _, p := range m.Segments {
			if p.Size != 0 {
				s = append(s, p)
			}
		}
		m.Segments = s
	}

	l := len(m.Symbols)
	for i := 0; i < l-1; i++ {
		x := &m.Symbols[i]
		y := &m.Symbols[i+1]
		if x.Segment >= uint64(len(m.Segments)) || y.Segment >= uint64(len(m.Segments)) {
			continue
		}

		if x.Segment != y.Segment {
			x.Size = m.Segments[x.Segment].Size - x.Offset
		} else {
			x.Size = y.Offset - x.Offset
		}
	}
	if l > 0 {
		x := &m.Symbols[l-1]
		if x.Segment < uint64(len(m.Segments)) {
			x.Size = m.Segments[x.Segment].Size - x.Offset
		}
	}

	for i := 0; i < len(m.Symbols); i++ {
		x := &m.Symbols[i]
		for i, y := range m.Sections {
			if x.Segment == y.Segment && y.Offset <= x.Offset && x.Offset <= y.Offset+y.Size {
				x.Section = i
				break
			}
		}
	}

	return m, nil
}
