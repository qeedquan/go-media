package pdb

type Header struct {
	Sig      [4]byte
	GUID     [16]byte
	Age      uint32
	Filename [18]byte
}

type StreamHeader struct {
	Version uint32
	Sig     uint32
	Age     uint32
	GUID    [16]byte
}

const (
	VC2     = 19941610
	VC4     = 19950623
	VC41    = 19950814
	VC50    = 19960307
	VC98    = 19970604
	VC70Dep = 19990604
	VC70    = 20000404
	VC80    = 20030901
	VC110   = 20091201
	VC140   = 20140508
)
