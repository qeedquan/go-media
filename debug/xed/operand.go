package xed

/*
#include <xed-interface.h>
*/
import "C"

type (
	Operand       C.xed_operand_t
	OperandValues C.xed_operand_values_t
)

func (c *Operand) Width() OperandWidth {
	return OperandWidth(C.xed_operand_width((*C.xed_operand_t)(c)))
}

func (c *Operand) Visibility() OperandVisibility {
	return OperandVisibility(C.xed_operand_operand_visibility((*C.xed_operand_t)(c)))
}

func (c *OperandValues) EffectiveOperandWidth() uint32 {
	return uint32(C.xed_operand_values_get_effective_operand_width((*C.xed_operand_values_t)(c)))
}

func (c *OperandValues) EffectiveAddressWidth() uint32 {
	return uint32(C.xed_operand_values_get_effective_address_width((*C.xed_operand_values_t)(c)))
}

func (c *OperandValues) StackAddressWidth() uint32 {
	return uint32(C.xed_operand_values_get_stack_address_width((*C.xed_operand_values_t)(c)))
}

func (c *OperandValues) IClass() IClass {
	return IClass(C.xed_operand_values_get_iclass((*C.xed_operand_values_t)(c)))
}

func (c *OperandValues) Init() {
	C.xed_operand_values_init((*C.xed_operand_values_t)(c))
}

func (c *OperandValues) HasMemoryDisplacement() bool {
	return xedbool(C.xed_operand_values_has_memory_displacement((*C.xed_operand_values_t)(c)))
}

func (c *OperandValues) HasRepPrefix() bool {
	return xedbool(C.xed_operand_values_has_rep_prefix((*C.xed_operand_values_t)(c)))
}

func (c *OperandValues) HasRepnePrefix() bool {
	return xedbool(C.xed_operand_values_has_repne_prefix((*C.xed_operand_values_t)(c)))
}

func (c *OperandValues) HasRexwPrefix() bool {
	return xedbool(C.xed_operand_values_has_rexw_prefix((*C.xed_operand_values_t)(c)))
}

func (c *OperandValues) HasSegmentPrefix() bool {
	return xedbool(C.xed_operand_values_has_segment_prefix((*C.xed_operand_values_t)(c)))
}

func (c *OperandValues) HasSibByte() bool {
	return xedbool(C.xed_operand_values_has_sib_byte((*C.xed_operand_values_t)(c)))
}

func (c *OperandValues) IsNop() bool {
	return xedbool(C.xed_operand_values_is_nop((*C.xed_operand_values_t)(c)))
}

func (c *OperandValues) Lockable() bool {
	return xedbool(C.xed_operand_values_lockable((*C.xed_operand_values_t)(c)))
}

func (c *OperandValues) MemopWithoutModrm() bool {
	return xedbool(C.xed_operand_values_memop_without_modrm((*C.xed_operand_values_t)(c)))
}

func (c *OperandValues) HasRealRep() bool {
	return xedbool(C.xed_operand_values_has_real_rep((*C.xed_operand_values_t)(c)))
}

func (c *OperandValues) HasAddressSizePrefix() bool {
	return xedbool(C.xed_operand_values_has_address_size_prefix((*C.xed_operand_values_t)(c)))
}

func (c *OperandValues) HasOperandSizePrefix() bool {
	return xedbool(C.xed_operand_values_has_operand_size_prefix((*C.xed_operand_values_t)(c)))
}

func (c *OperandValues) Has66Prefix() bool {
	return xedbool(C.xed_operand_values_has_66_prefix((*C.xed_operand_values_t)(c)))
}

func (c *OperandValues) HasBranchDisplacement() bool {
	return xedbool(C.xed_operand_values_has_branch_displacement((*C.xed_operand_values_t)(c)))
}

func (c *OperandValues) ImmediateInt64() int64 {
	return int64(C.xed_operand_values_get_immediate_int64((*C.xed_operand_values_t)(c)))
}

func (c *OperandValues) ImmediateUint64() uint64 {
	return uint64(C.xed_operand_values_get_immediate_uint64((*C.xed_operand_values_t)(c)))
}

func (c *OperandValues) ImmediateIsSigned() bool {
	return C.xed_operand_values_get_immediate_is_signed((*C.xed_operand_values_t)(c)) != 0
}

func (c *OperandValues) ImmediateByte(i uint) uint8 {
	return uint8(C.xed_operand_values_get_immediate_byte((*C.xed_operand_values_t)(c), C.uint(i)))
}

func (c *OperandValues) SecondImmediate() uint8 {
	return uint8(C.xed_operand_values_get_second_immediate((*C.xed_operand_values_t)(c)))
}

func (c *OperandValues) BranchDisplacementLength() uint32 {
	return uint32(C.xed_operand_values_get_branch_displacement_length((*C.xed_operand_values_t)(c)))
}

func (c *OperandValues) BranchDisplacementLengthBits() uint32 {
	return uint32(C.xed_operand_values_get_branch_displacement_length_bits((*C.xed_operand_values_t)(c)))
}

func (c *OperandValues) BranchDisplacementInt32() int32 {
	return int32(C.xed_operand_values_get_branch_displacement_int32((*C.xed_operand_values_t)(c)))
}

func (c *OperandValues) BranchDisplacementByte(i uint) uint8 {
	return uint8(C.xed_operand_values_get_branch_displacement_byte((*C.xed_operand_values_t)(c), C.uint(i)))
}

func (c *OperandValues) MemoryDisplacementInt64() int64 {
	return int64(C.xed_operand_values_get_memory_displacement_int64((*C.xed_operand_values_t)(c)))
}

func (c *OperandValues) MemoryDisplacementInt64Raw() int64 {
	return int64(C.xed_operand_values_get_memory_displacement_int64_raw((*C.xed_operand_values_t)(c)))
}

func (c *OperandValues) MemoryDisplacementByte(i uint) uint8 {
	return uint8(C.xed_operand_values_get_memory_displacement_byte((*C.xed_operand_values_t)(c), C.uint(i)))
}