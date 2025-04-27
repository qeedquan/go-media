package xed

/*
#include <xed-interface.h>
*/
import "C"

type (
	EncoderRequest C.xed_encoder_request_t
)

func (c *EncoderRequest) SetEffectiveOperandWidth(width_bits uint) {
	C.xed_encoder_request_set_effective_operand_width((*C.xed_encoder_request_t)(c), C.xed_uint_t(width_bits))
}

func (c *EncoderRequest) SetIclass(iclass IClass) {
	C.xed_encoder_request_set_iclass((*C.xed_encoder_request_t)(c), C.xed_iclass_enum_t(iclass))
}

func (c *EncoderRequest) SetIndex(index_reg Reg) {
	C.xed_encoder_request_set_index((*C.xed_encoder_request_t)(c), C.xed_reg_enum_t(index_reg))
}

func (c *EncoderRequest) SetMem0() {
	C.xed_encoder_request_set_mem0((*C.xed_encoder_request_t)(c))
}

func (c *EncoderRequest) SetMem1() {
	C.xed_encoder_request_set_mem1((*C.xed_encoder_request_t)(c))
}

func (c *EncoderRequest) SetMemoryDisplacement(memdisp int64, nbytes uint) {
	C.xed_encoder_request_set_memory_displacement((*C.xed_encoder_request_t)(c), C.xed_int64_t(memdisp), C.xed_uint_t(nbytes))
}

func (c *EncoderRequest) SetMemoryOperandLength(nbytes uint) {
	C.xed_encoder_request_set_memory_operand_length((*C.xed_encoder_request_t)(c), C.xed_uint_t(nbytes))
}

func (c *EncoderRequest) SetOperandOrder(operand_index uint, name OperandMode) {
	C.xed_encoder_request_set_operand_order((*C.xed_encoder_request_t)(c), C.xed_uint_t(operand_index), C.xed_operand_enum_t(name))
}

func (c *EncoderRequest) SetPtr() {
	C.xed_encoder_request_set_ptr((*C.xed_encoder_request_t)(c))
}

func (c *EncoderRequest) SetRelbr() {
	C.xed_encoder_request_set_relbr((*C.xed_encoder_request_t)(c))
}

func (c *EncoderRequest) SetScale(scale uint) {
	C.xed_encoder_request_set_scale((*C.xed_encoder_request_t)(c), C.xed_uint_t(scale))
}

func (c *EncoderRequest) SetUimm1(uimm uint8) {
	C.xed_encoder_request_set_uimm1((*C.xed_encoder_request_t)(c), C.xed_uint8_t(uimm))
}

func (c *EncoderRequest) Zero() {
	C.xed_encoder_request_zero((*C.xed_encoder_request_t)(c))
}

func (c *EncoderRequest) ZeroOperandOrder() {
	C.xed_encoder_request_zero_operand_order((*C.xed_encoder_request_t)(c))
}

func (c *EncoderRequest) String() string {
	var buf [512]C.char
	C.xed_encode_request_print((*C.xed_encoder_request_t)(c), &buf[0], C.xed_uint_t(len(buf)))
	return C.GoString(&buf[0])
}

func EncodeNop(array []uint8) error {
	if array == nil {
		return nil
	}
	return xederror(C.xed_encode_nop((*C.uint8_t)(&array[0]), C.uint(len(array))))
}
