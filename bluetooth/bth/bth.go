package bth

import "net"

type Addr struct {
	Net string
	HW  net.HardwareAddr
}

func (a Addr) Network() string { return a.Net }
func (a Addr) String() string  { return a.HW.String() }

var (
	ADDR_ANY   = net.HardwareAddr{0, 0, 0, 0, 0, 0}
	ADDR_ALL   = net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	ADDR_LOCAL = net.HardwareAddr{0, 0, 0, 0xff, 0xff, 0xff}
)

const (
	ADDR_BREDR     = 0x00
	ADDR_LE_PUBLIC = 0x01
	ADDR_LE_RANDOM = 0x02
)

type Security struct {
	Level   uint8
	KeySize uint8
}

type Voice struct {
	Setting uint8
}

const (
	SECURITY       = 4
	POWER          = 9
	CHANNEL_POLICY = 10
	VOICE          = 11
	SNDMTU         = 12
	RCVMTU         = 13
)

const (
	SECURITY_SDP    = 0
	SECURITY_LOW    = 1
	SECURITY_MEDIUM = 2
	SECURITY_HIGH   = 3
	SECURITY_FIPS   = 4
)

const (
	POWER_FORCE_ACTIVE_OFF = 0
	POWER_FORCE_ACTIVE_ON  = 1
)

const (
	CHANNEL_POLICY_BREDR_ONLY      = 0
	CHANNEL_POLICY_BREDR_PREFERRED = 1
	CHANNEL_POLICY_AMP_PREFERRED   = 2
)

const (
	VOICE_TRANSPARENT = 0x0003
	VOICE_CVSD_16BIT  = 0x0060
)
