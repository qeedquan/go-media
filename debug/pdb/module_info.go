package pdb

type ModiStream struct {
	Signature      uint32
	Symbols        []uint8
	C11LineInfo    []uint8
	C13LineInfo    []uint8
	GlobalRefsSize uint32
	GlobalRefs     []uint8
}
