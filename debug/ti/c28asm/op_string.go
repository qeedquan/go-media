// Code generated by "stringer -type Op ."; DO NOT EDIT.

package c28asm

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[UNKNOWN-0]
	_ = x[ABORTI-1]
	_ = x[ABS-2]
	_ = x[ABSTC-3]
	_ = x[ADD-4]
	_ = x[ADDB-5]
	_ = x[ADDCL-6]
	_ = x[ADDCU-7]
	_ = x[ADDL-8]
	_ = x[ADDU-9]
	_ = x[ADDUL-10]
	_ = x[ADRK-11]
	_ = x[AND-12]
	_ = x[ANDB-13]
	_ = x[ASP-14]
	_ = x[ASR-15]
	_ = x[ASR64-16]
	_ = x[ASRL-17]
	_ = x[B-18]
	_ = x[BANZ-19]
	_ = x[BAR-20]
	_ = x[BF-21]
	_ = x[C27MAP-22]
	_ = x[C27OBJ-23]
	_ = x[C28ADDR-24]
	_ = x[C28MAP-25]
	_ = x[C28OBJ-26]
	_ = x[CLRC-27]
	_ = x[CMP-28]
	_ = x[CMP64-29]
	_ = x[CMPB-30]
	_ = x[CMPL-31]
	_ = x[CMPR-32]
	_ = x[CSB-33]
	_ = x[DEC-34]
	_ = x[DINT-35]
	_ = x[DMAC-36]
	_ = x[DMOV-37]
	_ = x[EALLOW-38]
	_ = x[EDIS-39]
	_ = x[EINT-40]
	_ = x[ESTOP0-41]
	_ = x[ESTOP1-42]
	_ = x[FFC-43]
	_ = x[FLIP-44]
	_ = x[IACK-45]
	_ = x[IDLE-46]
	_ = x[IMACL-47]
	_ = x[IMPYAL-48]
	_ = x[IMPYL-49]
	_ = x[IMPYSL-50]
	_ = x[IMPYXUL-51]
	_ = x[IN-52]
	_ = x[INC-53]
	_ = x[INTR-54]
	_ = x[IRET-55]
	_ = x[LB-56]
	_ = x[LC-57]
	_ = x[LCR-58]
	_ = x[LOOPNZ-59]
	_ = x[LOOPZ-60]
	_ = x[LPADDR-61]
	_ = x[LRET-62]
	_ = x[LRETE-63]
	_ = x[LRETR-64]
	_ = x[LSL-65]
	_ = x[LSL64-66]
	_ = x[LSLL-67]
	_ = x[LSR-68]
	_ = x[LSR64-69]
	_ = x[LSRL-70]
	_ = x[MAC-71]
	_ = x[MAX-72]
	_ = x[MAXCUL-73]
	_ = x[MAXL-74]
	_ = x[MIN-75]
	_ = x[MINCUL-76]
	_ = x[MINL-77]
	_ = x[MOV-78]
	_ = x[MOVA-79]
	_ = x[MOVAD-80]
	_ = x[MOVB-81]
	_ = x[MOVDL-82]
	_ = x[MOVH-83]
	_ = x[MOVL-84]
	_ = x[MOVP-85]
	_ = x[MOVS-86]
	_ = x[MOVU-87]
	_ = x[MOVW-88]
	_ = x[MOVX-89]
	_ = x[MOVZ-90]
	_ = x[MPY-91]
	_ = x[MPYA-92]
	_ = x[MPYB-93]
	_ = x[MPYS-94]
	_ = x[MPYU-95]
	_ = x[MPYXU-96]
	_ = x[NASP-97]
	_ = x[NEG-98]
	_ = x[NEG64-99]
	_ = x[NEGTC-100]
	_ = x[NOP-101]
	_ = x[NORM-102]
	_ = x[NOT-103]
	_ = x[OR-104]
	_ = x[ORB-105]
	_ = x[OUT-106]
	_ = x[POP-107]
	_ = x[PREAD-108]
	_ = x[PUSH-109]
	_ = x[PWRITE-110]
	_ = x[QMACL-111]
	_ = x[QMPYAL-112]
	_ = x[QMPYL-113]
	_ = x[QMPYSL-114]
	_ = x[QMPYUL-115]
	_ = x[QMPYXUL-116]
	_ = x[ROL-117]
	_ = x[ROR-118]
	_ = x[RPT-119]
	_ = x[SAT-120]
	_ = x[SAT64-121]
	_ = x[SB-122]
	_ = x[SBBU-123]
	_ = x[SBF-124]
	_ = x[SBRK-125]
	_ = x[SETC-126]
	_ = x[SFR-127]
	_ = x[SPM-128]
	_ = x[SQRA-129]
	_ = x[SUB-130]
	_ = x[SUBB-131]
	_ = x[SUBBL-132]
	_ = x[SUBCU-133]
	_ = x[SUBCUL-134]
	_ = x[SUBLSUBR-135]
	_ = x[SUBRL-136]
	_ = x[SUBU-137]
	_ = x[SUBUL-138]
	_ = x[TBIT-139]
	_ = x[TCLR-140]
	_ = x[TEST-141]
	_ = x[TRAP-142]
	_ = x[TSET-143]
	_ = x[UOUT-144]
	_ = x[XB-145]
	_ = x[XBANZ-146]
	_ = x[XCALL-147]
	_ = x[XMAC-148]
	_ = x[XMACD-149]
	_ = x[XOR-150]
	_ = x[XPREAD-151]
	_ = x[XPWRITE-152]
	_ = x[XRET-153]
	_ = x[XRETC-154]
	_ = x[ZALR-155]
	_ = x[ZAP-156]
	_ = x[ZAPA-157]
}

const _Op_name = "UNKNOWNABORTIABSABSTCADDADDBADDCLADDCUADDLADDUADDULADRKANDANDBASPASRASR64ASRLBBANZBARBFC27MAPC27OBJC28ADDRC28MAPC28OBJCLRCCMPCMP64CMPBCMPLCMPRCSBDECDINTDMACDMOVEALLOWEDISEINTESTOP0ESTOP1FFCFLIPIACKIDLEIMACLIMPYALIMPYLIMPYSLIMPYXULININCINTRIRETLBLCLCRLOOPNZLOOPZLPADDRLRETLRETELRETRLSLLSL64LSLLLSRLSR64LSRLMACMAXMAXCULMAXLMINMINCULMINLMOVMOVAMOVADMOVBMOVDLMOVHMOVLMOVPMOVSMOVUMOVWMOVXMOVZMPYMPYAMPYBMPYSMPYUMPYXUNASPNEGNEG64NEGTCNOPNORMNOTORORBOUTPOPPREADPUSHPWRITEQMACLQMPYALQMPYLQMPYSLQMPYULQMPYXULROLRORRPTSATSAT64SBSBBUSBFSBRKSETCSFRSPMSQRASUBSUBBSUBBLSUBCUSUBCULSUBLSUBRSUBRLSUBUSUBULTBITTCLRTESTTRAPTSETUOUTXBXBANZXCALLXMACXMACDXORXPREADXPWRITEXRETXRETCZALRZAPZAPA"

var _Op_index = [...]uint16{0, 7, 13, 16, 21, 24, 28, 33, 38, 42, 46, 51, 55, 58, 62, 65, 68, 73, 77, 78, 82, 85, 87, 93, 99, 106, 112, 118, 122, 125, 130, 134, 138, 142, 145, 148, 152, 156, 160, 166, 170, 174, 180, 186, 189, 193, 197, 201, 206, 212, 217, 223, 230, 232, 235, 239, 243, 245, 247, 250, 256, 261, 267, 271, 276, 281, 284, 289, 293, 296, 301, 305, 308, 311, 317, 321, 324, 330, 334, 337, 341, 346, 350, 355, 359, 363, 367, 371, 375, 379, 383, 387, 390, 394, 398, 402, 406, 411, 415, 418, 423, 428, 431, 435, 438, 440, 443, 446, 449, 454, 458, 464, 469, 475, 480, 486, 492, 499, 502, 505, 508, 511, 516, 518, 522, 525, 529, 533, 536, 539, 543, 546, 550, 555, 560, 566, 574, 579, 583, 588, 592, 596, 600, 604, 608, 612, 614, 619, 624, 628, 633, 636, 642, 649, 653, 658, 662, 665, 669}

func (i Op) String() string {
	if i >= Op(len(_Op_index)-1) {
		return "Op(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Op_name[_Op_index[i]:_Op_index[i+1]]
}