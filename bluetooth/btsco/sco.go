package btsco

const (
	DEFAULT_MTU      = 500
	DEFAULT_FLUSH_TO = 0xFFFF
)

const (
	OPTIONS  = 0x01
	CONNINFO = 0x02
)

type Options struct {
	MTU uint16
}

type ConnInfo struct {
	HCIHandle uint16
	DevClass  [3]uint8
}
