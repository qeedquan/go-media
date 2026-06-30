package xed

/*
#include <xed-interface.h>
*/
import "C"

func (r Reg) GPRRegClass() RegClass {
	return RegClass(C.xed_gpr_reg_class(C.xed_reg_enum_t(r)))
}

func (r Reg) RegClass() RegClass {
	return RegClass(C.xed_reg_class(C.xed_reg_enum_t(r)))
}