package xed

/*
#include <xed-interface.h>
*/
import "C"
import "fmt"

type (
	FormatOptions       C.xed_format_options_t
	PrintInfo           C.xed_print_info_t
	Syntax              C.xed_syntax_enum_t
	DisassemblyCallback func(address uint64, symbol_buffer []byte, offset *uint64)
)

const (
	SYNTAX_INVALID Syntax = C.XED_SYNTAX_INVALID
	SYNTAX_XED     Syntax = C.XED_SYNTAX_XED
	SYNTAX_ATT     Syntax = C.XED_SYNTAX_ATT
	SYNTAX_INTEL   Syntax = C.XED_SYNTAX_INTEL
	SYNTAX_LAST    Syntax = C.XED_SYNTAX_LAST
)

func (c Syntax) String() string {
	return C.GoString(C.xed_syntax_enum_t2str(C.xed_syntax_enum_t(c)))
}

func (c *PrintInfo) Init() {
	C.xed_init_print_info((*C.xed_print_info_t)(c))
}

func (c *FormatOptions) SetLowerCaseHex(lowercase_hex bool) {
	c.lowercase_hex = boolxed(lowercase_hex)
}

func (c *FormatOptions) LowerCase() bool {
	return xedbool(c.lowercase_hex)
}

func FormatContext(syntax Syntax, xedd *DecodedInst, runtime_instruction_address uint64, symbolic_callback DisassemblyCallback) (string, error) {
	var buf [256]C.char
	rc := C.xed_format_context(C.xed_syntax_enum_t(syntax), (*C.xed_decoded_inst_t)(xedd), &buf[0], C.int(len(buf)), C.xed_uint64_t(runtime_instruction_address), nil, nil)
	if rc == 0 {
		return "", fmt.Errorf("failed to format instruction")
	}
	return C.GoString(&buf[0]), nil
}