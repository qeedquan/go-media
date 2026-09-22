package pdb

type DbiStreamHeader struct {
	VersionSignature        int32
	VersionHeader           uint32
	Age                     uint32
	GlobalStreamIndex       uint16
	BuildNumber             uint16
	PublicStreamIndex       uint16
	PdbDllVersion           uint16
	SymRecordStream         uint16
	PdbDllRbld              uint16
	ModInfoSize             int32
	SectionContributionSize int32
	SectionMapSize          int32
	SourceInfoSize          int32
	TypeServerSize          int32
	MFCTypeServerIndex      uint32
	OptionalDbgHeaderSize   int32
	ECSubstreamSize         int32
	Flags                   uint16
	Machine                 uint16
	_                       uint32
}

type SectionContribEntry struct {
	Section         uint16
	_               [2]byte
	Offset          int32
	Size            int32
	Characteristics uint32
	ModuleIndex     uint16
	_               [2]byte
	DataCrc         uint32
	RelocCrc        uint32
}

type ModInfo struct {
	_                    uint32
	SectionContr         SectionContribEntry
	Flags                uint32
	ModuleSymStream      uint16
	SymByteSize          uint32
	C11ByteSize          uint32
	C13ByteSize          uint32
	SourceFileCount      uint16
	_                    [2]byte
	_                    uint32
	SourceFileNameIndex  uint32
	PdbFilePathNameIndex uint32
	ModuleName           string
	ObjFileName          string
}

const (
	DBI_VC41 = 930803
	DBI_V50  = 19960307
	DBI_V60  = 19970606
	DBI_V70  = 19990903
	DBI_V110 = 20091201
)
