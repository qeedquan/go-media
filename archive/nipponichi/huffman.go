package nipponichi

import (
	"encoding/binary"
	"fmt"
)

type Huffman struct {
	freq []uint32
}

func (c *Huffman) Encode(src []byte) ([]byte, error) {
	return nil, fmt.Errorf("huffman: wip")
}

func (c *Huffman) Decode(src []byte) ([]byte, error) {
	c.freq = make([]uint32, 256)
	if len(src) < len(c.freq)*4+0x14 {
		return nil, fmt.Errorf("huffman: size too small")
	}
	if src[0x8] != 3 {
		return nil, fmt.Errorf("huffman: compression mismatch: expected 3, got %d", src[0x8])
	}

	typ := src[0x9]

	// this can either be total bits to represent the decoded file or the total file size of the encoded file
	// depending on the type
	total := binary.LittleEndian.Uint32(src[0x0c:])

	// size of the data output
	size := binary.LittleEndian.Uint32(src[0x10:])

	fmt.Println(typ, total, size)

	for i := range c.freq {
		c.freq[i] = binary.LittleEndian.Uint32(src[0x14+4*i:])
	}

	return nil, fmt.Errorf("huffman: wip")
}
