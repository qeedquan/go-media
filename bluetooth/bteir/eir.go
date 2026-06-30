package bteir

// Extended Inquiry Response
const (
	FLAGS               = 0x01 // flags
	UUID16_SOME         = 0x02 // 16-bit UUID, more available
	UUID16_ALL          = 0x03 // 16-bit UUID, all listed
	UUID32_SOME         = 0x04 // 32-bit UUID, more available
	UUID32_ALL          = 0x05 // 32-bit UUID, all listed
	UUID128_SOME        = 0x06 // 128-bit UUID, more available
	UUID128_ALL         = 0x07 // 128-bit UUID, all listed
	NAME_SHORT          = 0x08 // shortened local name
	NAME_COMPLETE       = 0x09 // complete local name
	TX_POWER            = 0x0A // transmit power level
	CLASS_OF_DEV        = 0x0D // Class of Device
	SSP_HASH            = 0x0E // SSP Hash
	SSP_RANDOMIZER      = 0x0F // SSP Randomizer
	DEVICE_ID           = 0x10 // device ID
	SOLICIT16           = 0x14 // LE: Solicit UUIDs, 16-bit
	SOLICIT128          = 0x15 // LE: Solicit UUIDs, 128-bit
	SVC_DATA16          = 0x16 // LE: Service data, 16-bit UUID
	PUB_TRGT_ADDR       = 0x17 // LE: Public Target Address
	RND_TRGT_ADDR       = 0x18 // LE: Random Target Address
	GAP_APPEARANCE      = 0x19 // GAP appearance
	SOLICIT32           = 0x1F // LE: Solicit UUIDs, 32-bit
	SVC_DATA32          = 0x20 // LE: Service data, 32-bit UUID
	SVC_DATA128         = 0x21 // LE: Service data, 128-bit UUID
	TRANSPORT_DISCOVERY = 0x26 // Transport Discovery Service
	MANUFACTURER_DATA   = 0xFF // Manufacturer Specific Data
)

// Flags Descriptions
const (
	LIM_DISC    = 0x01 // LE Limited Discoverable Mode
	GEN_DISC    = 0x02 // LE General Discoverable Mode
	BREDR_UNSUP = 0x04 // BR/EDR Not Supported
	CONTROLLER  = 0x08 // Simultaneous LE and BR/EDR to Same Device Capable (Controller)
	SIM_HOST    = 0x10 // Simultaneous LE and BR/EDR to Same Device Capable (Host)

	SD_MAX_LEN  = 238 // 240 (EIR) - 2 (len)
	MSD_MAX_LEN = 236 // 240 (EIR) - 2 (len & type) - 2
)
