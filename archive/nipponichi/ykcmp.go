package nipponichi

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"os/exec"
	"runtime"
	"strings"
)

type YKCMP struct{}

func (z *YKCMP) Encode(kind string, src []byte) ([]byte, error) {
	w := new(bytes.Buffer)

	var (
		enc []byte
		err error
	)
	switch kind {
	case "S2", "s2":
		var slide Slide
		enc, err = slide.Encode2(src)

	case "H", "h": // Huffman
		fallthrough

	case "R", "r": // RLE
		fallthrough

	case "Z", "z": // ZLIB
		fallthrough

	case "lzma": // LZMA
		fallthrough

	case "lz4": // LZ4
		fallthrough

	default:
		enc, err = encdecex("e "+kind, src)
	}

	w.Write(enc)
	buf := w.Bytes()
	copy(buf, []byte("YKCMP_V1"))
	return buf, err
}

func (z *YKCMP) Decode(src []byte) ([]byte, error) {
	if !bytes.HasPrefix(src, []byte("YKCMP_V1")) {
		return nil, fmt.Errorf("ykcmp: invalid header: %q", src[:8])
	}

	var (
		dst []byte
		err error
	)
	typ := binary.LittleEndian.Uint32(src[8:])
	switch typ {
	case 4:
		var slide Slide
		dst, err = slide.Decode2(src, math.MaxUint32)

	case 1: // RLE
		fallthrough
	case 2: // Slide V1
		fallthrough
	case 3: // Huffman
		fallthrough
	case 8, 9: // LZ4
		fallthrough
	case 10: // LZMA
		fallthrough

	default:
		dst, err = encdecex("d", src)
	}

	return dst, err
}

func encdecex(op string, data []byte) ([]byte, error) {
	prog := "nmpl.exe"
	argv := strings.Split(op, " ")
	if runtime.GOOS != "windows" {
		argv = append([]string{prog}, argv...)
		prog = "wine64"
	}
	stdout := new(bytes.Buffer)
	cmd := exec.Command(prog, argv...)
	cmd.Stdin = bytes.NewBuffer(data)
	cmd.Stdout = stdout
	err := cmd.Run()
	return stdout.Bytes(), err
}
