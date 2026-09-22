package xed

/*
#include <xed-interface.h>

xed_uint32_t get_flag_set_flat(xed_flag_set_t *c) {
	return c->flat;
}

xed_uint32_t set_flag_set_flat(xed_flag_set_t *c, xed_uint32_t flat) {
	c->flat = flat;
}
*/
import "C"

type (
	FlagAction C.xed_flag_action_t
	FlagSet    C.xed_flag_set_t
)

const (
	FLAGSET_ZF = (1 << 6)
)

func (c *SimpleFlag) ReadFlagSet() *FlagSet {
	return (*FlagSet)(C.xed_simple_flag_get_read_flag_set((*C.xed_simple_flag_t)(c)))
}

func (c *SimpleFlag) WrittenFlagSet() *FlagSet {
	return (*FlagSet)(C.xed_simple_flag_get_written_flag_set((*C.xed_simple_flag_t)(c)))
}

func (c *SimpleFlag) UndefinedFlagSet() *FlagSet {
	return (*FlagSet)(C.xed_simple_flag_get_undefined_flag_set((*C.xed_simple_flag_t)(c)))
}

func (c *SimpleFlag) ReadsFlags() bool {
	return xedbool(C.xed_simple_flag_reads_flags((*C.xed_simple_flag_t)(c)))
}

func (c *SimpleFlag) WritesFlags() bool {
	return xedbool(C.xed_simple_flag_writes_flags((*C.xed_simple_flag_t)(c)))
}

func (c *SimpleFlag) MayWrite() bool {
	return xedbool(C.xed_simple_flag_get_may_write((*C.xed_simple_flag_t)(c)))
}

func (c *SimpleFlag) MustWrite() bool {
	return xedbool(C.xed_simple_flag_get_must_write((*C.xed_simple_flag_t)(c)))
}

func (c *SimpleFlag) NumFlags() uint {
	return uint(C.xed_simple_flag_get_nflags((*C.xed_simple_flag_t)(c)))
}

func (c *SimpleFlag) FlagAction(i uint) *FlagAction {
	return (*FlagAction)(C.xed_simple_flag_get_flag_action((*C.xed_simple_flag_t)(c), C.uint(i)))
}

func (c *FlagSet) SetFlat(flat uint32) {
	C.set_flag_set_flat((*C.xed_flag_set_t)(c), C.xed_uint32_t(flat))
}

func (c *FlagSet) Flat() uint32 {
	return uint32(C.get_flag_set_flat((*C.xed_flag_set_t)(c)))
}

func (c *FlagSet) Mask() int {
	return int(C.xed_flag_set_mask((*C.xed_flag_set_t)(c)))
}

func (c *FlagSet) String() string {
	var buf [4096]C.char
	C.xed_flag_set_print((*C.xed_flag_set_t)(c), &buf[0], C.int(len(buf)))
	return C.GoString(&buf[0])
}

func (c *FlagAction) String() string {
	var buf [4096]C.char
	C.xed_flag_action_print((*C.xed_flag_action_t)(c), &buf[0], C.int(len(buf)))
	return C.GoString(&buf[0])
}