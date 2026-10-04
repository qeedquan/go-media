package xed

/*
#include <xed-interface.h>

*/
import "C"

type (
	DecodedInst C.xed_decoded_inst_t
)

func (c *DecodedInst) NumOperands() uint {
	return uint(C.xed_decoded_inst_noperands((*C.xed_decoded_inst_t)(c)))
}

func (c *DecodedInst) OperandAction(operand_index uint) OperandAction {
	return OperandAction(C.xed_decoded_inst_operand_action((*C.xed_decoded_inst_t)(c), C.uint(operand_index)))
}

func (c *DecodedInst) OperandLength(operand_index uint) uint {
	return uint(C.xed_decoded_inst_operand_length((*C.xed_decoded_inst_t)(c), C.uint(operand_index)))
}

func (c *DecodedInst) OperandLengthBits(operand_index uint) uint {
	return uint(C.xed_decoded_inst_operand_length_bits((*C.xed_decoded_inst_t)(c), C.uint(operand_index)))
}

func (c *DecodedInst) OperandElements(operand_index uint) uint {
	return uint(C.xed_decoded_inst_operand_elements((*C.xed_decoded_inst_t)(c), C.uint(operand_index)))
}

func (c *DecodedInst) OperandElementSizeBits(operand_index uint) uint {
	return uint(C.xed_decoded_inst_operand_element_size_bits((*C.xed_decoded_inst_t)(c), C.uint(operand_index)))
}

func (c *DecodedInst) Zero() {
	C.xed_decoded_inst_zero((*C.xed_decoded_inst_t)(c))
}

func (c *DecodedInst) ISASet() ISASet {
	return ISASet(C.xed_decoded_inst_get_isa_set((*C.xed_decoded_inst_t)(c)))
}

func (c *DecodedInst) SetMode(mmode MachineMode, stack_addr_width AddressWidth) {
	C.xed_decoded_inst_set_mode((*C.xed_decoded_inst_t)(c), C.xed_machine_mode_enum_t(mmode), C.xed_address_width_enum_t(stack_addr_width))
}

func (c *DecodedInst) OperandsConst() *OperandValues {
	return (*OperandValues)(C.xed_decoded_inst_operands_const((*C.xed_decoded_inst_t)(c)))
}

func (c *DecodedInst) Operands() *OperandValues {
	return (*OperandValues)(C.xed_decoded_inst_operands((*C.xed_decoded_inst_t)(c)))
}

func (c *DecodedInst) SetScale(scale int) {
	C.xed_decoded_inst_set_scale((*C.xed_decoded_inst_t)(c), C.xed_uint_t(scale))
}

func (c *DecodedInst) Valid() bool {
	return xedbool(C.xed_decoded_inst_valid((*C.xed_decoded_inst_t)(c)))
}

func (c *DecodedInst) Decode(itext []byte) error {
	if len(itext) == 0 {
		return xederror(C.xed_decode((*C.xed_decoded_inst_t)(c), nil, 0))
	}
	return xederror(C.xed_decode((*C.xed_decoded_inst_t)(c), (*C.xed_uint8_t)(&itext[0]), C.uint(len(itext))))
}

func (c *DecodedInst) Length() uint {
	return uint(C.xed_decoded_inst_get_length((*C.xed_decoded_inst_t)(c)))
}

func (c *DecodedInst) Extension() Extension {
	return Extension(C.xed_decoded_inst_get_extension((*C.xed_decoded_inst_t)(c)))
}

func (c *DecodedInst) MachineModeBits() int {
	return int(C.xed_decoded_inst_get_machine_mode_bits((*C.xed_decoded_inst_t)(c)))
}

func (c *DecodedInst) MemopAddressWidth(memop_idx int) int {
	return int(C.xed_decoded_inst_get_memop_address_width((*C.xed_decoded_inst_t)(c), C.xed_uint_t(memop_idx)))
}

func (c *DecodedInst) Modrm() int {
	return int(C.xed_decoded_inst_get_modrm((*C.xed_decoded_inst_t)(c)))
}

func (c *DecodedInst) Reg(reg_operand OperandMode) Reg {
	return Reg(C.xed_decoded_inst_get_reg((*C.xed_decoded_inst_t)(c), C.xed_operand_enum_t(reg_operand)))
}

func (c *DecodedInst) OperandWidth() uint32 {
	return uint32(C.xed_decoded_inst_get_operand_width((*C.xed_decoded_inst_t)(c)))
}

func (c *DecodedInst) IClass() IClass {
	return IClass(C.xed_decoded_inst_get_iclass((*C.xed_decoded_inst_t)(c)))
}

func (c *DecodedInst) MemRead(mem_idx uint) bool {
	return xedbool(C.xed_decoded_inst_mem_read((*C.xed_decoded_inst_t)(c), C.uint(mem_idx)))
}

func (c *DecodedInst) MemWritten(mem_idx uint) bool {
	return xedbool(C.xed_decoded_inst_mem_written((*C.xed_decoded_inst_t)(c), C.uint(mem_idx)))
}

func (c *DecodedInst) MemWrittenOnly(mem_idx uint) bool {
	return xedbool(C.xed_decoded_inst_mem_written_only((*C.xed_decoded_inst_t)(c), C.uint(mem_idx)))
}

func (c *DecodedInst) Merging() bool {
	return xedbool(C.xed_decoded_inst_merging((*C.xed_decoded_inst_t)(c)))
}

func (c *DecodedInst) UsesRflags() bool {
	return xedbool(C.xed_decoded_inst_uses_rflags((*C.xed_decoded_inst_t)(c)))
}

func (c *DecodedInst) BaseReg(mem_idx uint) Reg {
	return Reg(C.xed_decoded_inst_get_base_reg((*C.xed_decoded_inst_t)(c), C.uint(mem_idx)))
}

func (c *DecodedInst) BranchDisplacement() int32 {
	return int32(C.xed_decoded_inst_get_branch_displacement((*C.xed_decoded_inst_t)(c)))
}

func (c *DecodedInst) BranchDisplacementWidth() uint {
	return uint(C.xed_decoded_inst_get_branch_displacement_width((*C.xed_decoded_inst_t)(c)))
}

func (c *DecodedInst) BranchDisplacementWidthBits() uint {
	return uint(C.xed_decoded_inst_get_branch_displacement_width_bits((*C.xed_decoded_inst_t)(c)))
}

func (c *DecodedInst) ImmediateWidthBits() uint {
	return uint(C.xed_decoded_inst_get_immediate_width_bits((*C.xed_decoded_inst_t)(c)))
}

func (c *DecodedInst) ImmediateIsSigned() uint {
	return uint(C.xed_decoded_inst_get_immediate_is_signed((*C.xed_decoded_inst_t)(c)))
}

func (c *DecodedInst) OperandElementType(operand_index uint) OperandElementType {
	return OperandElementType(C.xed_decoded_inst_operand_element_type((*C.xed_decoded_inst_t)(c), C.uint(operand_index)))
}
func (c *DecodedInst) SignedImmediate() int32 {
	return int32(C.xed_decoded_inst_get_signed_immediate((*C.xed_decoded_inst_t)(c)))
}

func (c *DecodedInst) UnsignedImmediate() uint64 {
	return uint64(C.xed_decoded_inst_get_unsigned_immediate((*C.xed_decoded_inst_t)(c)))
}

func (c *DecodedInst) SecondImmediate() uint8 {
	return uint8(C.xed_decoded_inst_get_second_immediate((*C.xed_decoded_inst_t)(c)))
}

func (c *DecodedInst) PatchDisp(itext []byte, disp EncDisplacement) bool {
	if itext == nil {
		return xedbool(C.xed_patch_disp((*C.xed_decoded_inst_t)(c), nil, C.xed_enc_displacement_t(disp)))

	}
	return xedbool(C.xed_patch_disp((*C.xed_decoded_inst_t)(c), (*C.xed_uint8_t)(&itext[0]), C.xed_enc_displacement_t(disp)))
}

func (c *DecodedInst) PatchRelBr(itext []byte, disp EncoderOperand) bool {
	if itext == nil {
		return xedbool(C.xed_patch_relbr((*C.xed_decoded_inst_t)(c), nil, C.xed_encoder_operand_t(disp)))

	}
	return xedbool(C.xed_patch_relbr((*C.xed_decoded_inst_t)(c), (*C.xed_uint8_t)(&itext[0]), C.xed_encoder_operand_t(disp)))
}

func (c *DecodedInst) PatchImm0(itext []byte, disp EncoderOperand) bool {
	if itext == nil {
		return xedbool(C.xed_patch_imm0((*C.xed_decoded_inst_t)(c), nil, C.xed_encoder_operand_t(disp)))

	}
	return xedbool(C.xed_patch_imm0((*C.xed_decoded_inst_t)(c), (*C.xed_uint8_t)(&itext[0]), C.xed_encoder_operand_t(disp)))
}

func (c *DecodedInst) ZeroSetMode(state *State) {
	C.xed_decoded_inst_zero_set_mode((*C.xed_decoded_inst_t)(c), (*C.xed_state_t)(state))
}

func (c *DecodedInst) Zeroing() bool {
	return xedbool(C.xed_decoded_inst_zeroing((*C.xed_decoded_inst_t)(c)))
}

func (c *DecodedInst) SetInputChip(chip Chip) {
	C.xed_decoded_inst_set_input_chip((*C.xed_decoded_inst_t)(c), C.xed_chip_enum_t(chip))
}

func (c *DecodedInst) RflagsInfo() *SimpleFlag {
	return (*SimpleFlag)(C.xed_decoded_inst_get_rflags_info((*C.xed_decoded_inst_t)(c)))
}

func (c *DecodedInst) SetMpxMode(opval Bits) {
	C.xed3_operand_set_mpxmode((*C.xed_decoded_inst_t)(c), C.xed_bits_t(opval))
}

func (c *DecodedInst) SetCetMode(opval Bits) {
	C.xed3_operand_set_cet((*C.xed_decoded_inst_t)(c), C.xed_bits_t(opval))
}

func (c *DecodedInst) Category() Category {
	return Category(C.xed_decoded_inst_get_category((*C.xed_decoded_inst_t)(c)))
}

func (c *DecodedInst) Iform() Iform {
	return Iform(C.xed_decoded_inst_get_iform_enum((*C.xed_decoded_inst_t)(c)))
}

func (c *DecodedInst) Inst() *Inst {
	return (*Inst)(C.xed_decoded_inst_inst((*C.xed_decoded_inst_t)(c)))
}

func (c *DecodedInst) IformDispatch() uint {
	return uint(C.xed_decoded_inst_get_iform_enum_dispatch((*C.xed_decoded_inst_t)(c)))
}

func (c *DecodedInst) NumPrefixes() uint {
	return uint(C.xed_decoded_inst_get_nprefixes((*C.xed_decoded_inst_t)(c)))
}

func (c *DecodedInst) Attribute(attr Attribute) uint32 {
	return uint32(C.xed_decoded_inst_get_attribute((*C.xed_decoded_inst_t)(c), C.xed_attribute_enum_t(attr)))
}

func (c *DecodedInst) Attributes() Attributes {
	return Attributes(C.xed_decoded_inst_get_attributes((*C.xed_decoded_inst_t)(c)))
}

func (c *DecodedInst) IsBroadcast() bool {
	return xedbool(C.xed_decoded_inst_is_broadcast((*C.xed_decoded_inst_t)(c)))
}

func (c *DecodedInst) ClassifySSE() bool {
	return xedbool(C.xed_classify_sse((*C.xed_decoded_inst_t)(c)))
}

func (c *DecodedInst) ClassifyAVX() bool {
	return xedbool(C.xed_classify_avx((*C.xed_decoded_inst_t)(c)))
}

func (c *DecodedInst) ClassifyAVX512() bool {
	return xedbool(C.xed_classify_avx512((*C.xed_decoded_inst_t)(c)))
}

func (c *DecodedInst) ClassifyAVX512MaskOp() bool {
	return xedbool(C.xed_classify_avx512_maskop((*C.xed_decoded_inst_t)(c)))
}

func (c *DecodedInst) MaskedVectorOperation() bool {
	return xedbool(C.xed_decoded_inst_masked_vector_operation((*C.xed_decoded_inst_t)(c)))
}

func (c *DecodedInst) Byte(byte_index uint) uint8 {
	return uint8(C.xed_decoded_inst_get_byte((*C.xed_decoded_inst_t)(c), C.uint(byte_index)))
}

func (c *DecodedInst) VectorLengthBits() uint {
	return uint(C.xed_decoded_inst_vector_length_bits((*C.xed_decoded_inst_t)(c)))
}

func (c *DecodedInst) NumberOfMemoryOperands() uint {
	return uint(C.xed_decoded_inst_number_of_memory_operands((*C.xed_decoded_inst_t)(c)))
}

func (c *DecodedInst) SegReg(mem_idx uint) Reg {
	return Reg(C.xed_decoded_inst_get_seg_reg((*C.xed_decoded_inst_t)(c), C.uint(mem_idx)))
}

func (c *DecodedInst) IndexReg(mem_idx uint) Reg {
	return Reg(C.xed_decoded_inst_get_index_reg((*C.xed_decoded_inst_t)(c), C.uint(mem_idx)))
}

func (c *DecodedInst) Scale(mem_idx uint) uint {
	return uint(C.xed_decoded_inst_get_scale((*C.xed_decoded_inst_t)(c), C.uint(mem_idx)))
}

func (c *DecodedInst) MemoryDisplacement(mem_idx uint) uint {
	return uint(C.xed_decoded_inst_get_memory_displacement((*C.xed_decoded_inst_t)(c), C.uint(mem_idx)))
}

func (c *DecodedInst) MemoryDisplacementWidth(mem_idx uint) uint {
	return uint(C.xed_decoded_inst_get_memory_displacement_width((*C.xed_decoded_inst_t)(c), C.uint(mem_idx)))
}

func (c *DecodedInst) MemoryDisplacementWidthBits(mem_idx uint) uint {
	return uint(C.xed_decoded_inst_get_memory_displacement_width_bits((*C.xed_decoded_inst_t)(c), C.uint(mem_idx)))
}

func (c *DecodedInst) MemOpAddressWidth(mem_idx uint) uint {
	return uint(C.xed_decoded_inst_get_memop_address_width((*C.xed_decoded_inst_t)(c), C.uint(mem_idx)))
}

func (c *DecodedInst) MemoryOperandLength(mem_idx uint) uint {
	return uint(C.xed_decoded_inst_get_memory_operand_length((*C.xed_decoded_inst_t)(c), C.uint(mem_idx)))
}

func (c *DecodedInst) Dump() string {
	var buf [512]C.char
	C.xed_decoded_inst_dump((*C.xed_decoded_inst_t)(c), &buf[0], C.int(len(buf)))
	return C.GoString(&buf[0])
}

func (c *DecodedInst) DumpXEDFormat(runtime_address uint64) (string, bool) {
	var buf [512]C.char
	b := C.xed_decoded_inst_dump_xed_format((*C.xed_decoded_inst_t)(c), &buf[0], C.int(len(buf)), C.xed_uint64_t(runtime_address))
	return C.GoString(&buf[0]), xedbool(b)
}

func (c *DecodedInst) String() string {
	str, _ := FormatContext(C.XED_SYNTAX_XED, c, 0, nil)
	return str
}