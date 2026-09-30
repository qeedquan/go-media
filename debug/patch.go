package debug

import (
	"bytes"
	"io"
	"sort"
)

type Patch struct {
	Start int64
	End   int64
	Data  []byte
}

type Patchset struct {
	List []Patch
	Gap  []Patch
	Data []byte
	Len  int64
	off  int64
}

func FindSignatures(data, sig []byte) []int {
	var r []int
	p := data
	for {
		i := bytes.Index(p, sig)
		if i < 0 {
			break
		}
		r = append(r, i)
		p = p[i+len(sig):]
	}

	for i := 1; i < len(r); i++ {
		r[i] += r[i-1]
	}

	return r
}

func FindPatches(data, start, end []byte) []Patch {
	x := FindSignatures(data, start)
	y := FindSignatures(data, end)

	var p []Patch
	for i, j := 0, 0; i < len(x) && j < len(y); {
		if y[j] < x[i] {
			j++
			continue
		}

		for i+1 < len(x) && x[i+1] < y[j] {
			i++
		}

		p = append(p, Patch{int64(x[i]), int64(y[j]), data[x[i] : y[j]+len(end)]})
		i++
		j++
	}
	return p
}

func (p *Patchset) WriteAt(b []byte, off int64) (n int, err error) {
	if len(b) == 0 {
		return
	}
	p.List = append(p.List, Patch{
		Start: off,
		End:   off + int64(len(b)) - 1,
		Data:  append([]byte{}, b...),
	})
	if off+int64(len(b)) > p.Len {
		p.Len = off + int64(len(b))
	}
	return len(b), nil
}

func (p *Patchset) Write(b []byte) (n int, err error) {
	p.WriteAt(b, p.off)
	p.off += int64(len(b))
	return len(b), nil
}

func (p *Patchset) Seek(off int64, whence int) {
	switch whence {
	case io.SeekStart:
		p.off = off
	case io.SeekCurrent:
		p.off += off
	case io.SeekEnd:
		p.off = p.Len + off
	}
	if p.Len < p.off {
		p.Len = p.off
	}
}

func (p *Patchset) Merge() {
	p.Data = make([]byte, p.Len)
	for _, l := range p.List {
		copy(p.Data[l.Start:l.End+1], l.Data)
	}

	sort.SliceStable(p.List, func(i, j int) bool {
		return p.List[i].Start < p.List[j].Start
	})

	n := 0
	for i := 0; i < len(p.List); i++ {
		if n == 0 || p.List[n-1].End < p.List[i].Start {
			p.List[n], n = p.List[i], n+1
		} else {
			p.List[n-1].End = max(p.List[n-1].End, p.List[i].End)
		}
	}
	p.List = p.List[:n]

	n = 0
	for i := 0; i < len(p.List); {
		l := p.List[i]
		j := i
		for ; j+1 < len(p.List) && p.List[j].End+1 == p.List[j+1].Start; j++ {
			l.End = p.List[j+1].End
		}
		i = j + 1
		p.List[n], n = l, n+1
	}
	p.List = p.List[:n]

	for i := range p.List {
		l := &p.List[i]
		l.Data = make([]byte, l.End-l.Start+1)
		copy(l.Data, p.Data[l.Start:])
	}

	p.Gap = p.Gap[:0]
	for i := 0; i < len(p.List); i++ {
		if i == 0 && p.List[i].Start > 0 {
			p.Gap = append(p.Gap, Patch{End: p.List[i].Start - 1})
		} else if i+1 < len(p.List) {
			p.Gap = append(p.Gap, Patch{Start: p.List[i].End + 1, End: p.List[i+1].Start - 1})
		}
	}
}
