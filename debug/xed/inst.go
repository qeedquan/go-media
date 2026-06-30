package xed

/*
#include <xed-interface.h>
*/
import "C"

type (
	Inst C.xed_inst_t
)

func (c *Inst) Attribute(attr Attribute) uint32 {
	return uint32(C.xed_inst_get_attribute((*C.xed_inst_t)(c), C.xed_attribute_enum_t(attr)))
}

func (c *Inst) NumOperands() uint {
	return uint(C.xed_inst_noperands((*C.xed_inst_t)(c)))
}

func (c *Inst) Operand(i uint) *Operand {
	return (*Operand)(C.xed_inst_operand((*C.xed_inst_t)(c), C.uint(i)))
}

func (c *Inst) Exception() Exception {
	return Exception(C.xed_inst_exception((*C.xed_inst_t)(c)))
}