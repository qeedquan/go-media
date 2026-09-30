// +build linux

package bthci

import (
	"errors"
	"fmt"
	"net"
	"unsafe"

	"golang.org/x/sys/unix"
)

const (
	DEV_NONE = 0xffff
)

type Version struct {
	Manufacturer uint16
	HCIVer       uint8
	HCIRev       uint16
	LMPVer       uint8
	LMPSubver    uint16
}

type DeviceInfo struct {
	DevID uint16
	Name  [8]byte

	BDAddr [6]uint8

	Flags uint32
	Type  uint8

	Features [8]uint8

	PktType    uint32
	LinkPolicy uint32
	LinkMode   uint32

	AclMtu  uint16
	AclPkts uint16
	ScoMtu  uint16
	ScoPkts uint16

	Stat DeviceStats
}

type DeviceStats struct {
	ErrRx  uint32
	ErrTx  uint32
	CmdTx  uint32
	EvtRx  uint32
	AclTx  uint32
	AclRx  uint32
	ScoTx  uint32
	ScoRx  uint32
	ByteRx uint32
	ByteTx uint32
}

const (
	LE_PUBLIC_ADDRESS = 0x00
	LE_RANDOM_ADDRESS = 0x01
)

const (
	DEV_REG     = 1
	DEV_UNREG   = 2
	DEV_UP      = 3
	DEV_DOWN    = 4
	DEV_SUSPEND = 5
	DEV_RESUME  = 6
)

const (
	VIRTUAL = 0
	USB     = 1
	PCCARD  = 2
	UART    = 3
	RS232   = 4
	PCI     = 5
	SDIO    = 6
	SPI     = 7
	I2C     = 8
	SMD     = 9
)

const (
	COMMAND_PKT = 0x01
	ACLDATA_PKT = 0x02
	SCODATA_PKT = 0x03
	EVENT_PKT   = 0x04
	VENDOR_PKT  = 0xff

	P2DH1 = 0x0003
	P3DH1 = 0x0004
	DM1   = 0x0008
	DH1   = 0x0010
	P2DH3 = 0x0100
	P3DH3 = 0x0200
	DM3   = 0x0400
	DH3   = 0x0800
	P2DH5 = 0x1000
	P3DH5 = 0x2000
	DM5   = 0x4000
	DH5   = 0x8000

	HV1 = 0x0020
	HV2 = 0x0040
	HV3 = 0x0080

	EV3   = 0x0008
	EV4   = 0x0010
	EV5   = 0x0020
	P2EV3 = 0x0040
	P3EV3 = 0x0080
	P2EV5 = 0x0100
	P3EV5 = 0x0200

	SCO_PTYPE_MASK = (HV1 | HV2 | HV3)
	ACL_PTYPE_MASK = (DM1 | DH1 | DM3 | DH3 | DM5 | DH5)
)

func ior(t, nr, size uintptr) uintptr {
	return (2 << 30) | (t << 8) | nr | (size << 16)
}

func iow(t, nr, size uintptr) uintptr {
	return (1 << 30) | (t << 8) | nr | (size << 16)
}

var (
	DEVUP     = iow('H', 201, 4)
	DEVDOWN   = iow('H', 202, 4)
	DEVRESET  = iow('H', 203, 4)
	DEVRESTAT = iow('H', 204, 4)

	GETDEVLIST  = ior('H', 210, 4)
	GETDEVINFO  = ior('H', 211, 4)
	GETCONNLIST = ior('H', 212, 4)
	GETCONNINFO = ior('H', 213, 4)
	GETAUTINFO  = ior('H', 215, 4)

	SETRAW      = iow('H', 220, 4)
	SETSCAN     = iow('H', 221, 4)
	SETAUTH     = iow('H', 222, 4)
	SETENCRYPT  = iow('H', 223, 4)
	SETTYPE     = iow('H', 224, 4)
	SETLINKPOL  = iow('H', 225, 4)
	SETLINKMODE = iow('H', 226, 4)
	SETACLMTU   = iow('H', 227, 4)
	SETSCOMTU   = iow('H', 228, 4)

	BLOCKADDR   = iow('H', 230, 4)
	UNBLOCKADDR = iow('H', 231, 4)

	INQUIRY = ior('H', 240, 4)
)

func LookupDeviceInfo(s string) (di DeviceInfo, err error) {
	var d int
	n, _ := fmt.Sscanf(s, "hci%d", &d)
	if n != 1 {
		err = fmt.Errorf("invalid hci device: %s", s)
		return
	}

	dd, err := unix.Socket(unix.AF_BLUETOOTH, unix.SOCK_RAW|unix.SOCK_CLOEXEC, unix.BTPROTO_HCI)
	if err != nil {
		return
	}
	defer unix.Close(dd)

	_, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(dd), uintptr(GETDEVINFO), uintptr(unsafe.Pointer(&di)))
	if errno != 0 {
		err = errno
	}

	return
}

// Error codes
const (
	UNKNOWN_COMMAND                 = 0x01
	NO_CONNECTION                   = 0x02
	HARDWARE_FAILURE                = 0x03
	PAGE_TIMEOUT                    = 0x04
	AUTHENTICATION_FAILURE          = 0x05
	PIN_OR_KEY_MISSING              = 0x06
	MEMORY_FULL                     = 0x07
	CONNECTION_TIMEOUT              = 0x08
	MAX_NUMBER_OF_CONNECTIONS       = 0x09
	MAX_NUMBER_OF_SCO_CONNECTIONS   = 0x0a
	ACL_CONNECTION_EXISTS           = 0x0b
	COMMAND_DISALLOWED              = 0x0c
	REJECTED_LIMITED_RESOURCES      = 0x0d
	REJECTED_SECURITY               = 0x0e
	REJECTED_PERSONAL               = 0x0f
	HOST_TIMEOUT                    = 0x10
	UNSUPPORTED_FEATURE             = 0x11
	INVALID_PARAMETERS              = 0x12
	OE_USER_ENDED_CONNECTION        = 0x13
	OE_LOW_RESOURCES                = 0x14
	OE_POWER_OFF                    = 0x15
	CONNECTION_TERMINATED           = 0x16
	REPEATED_ATTEMPTS               = 0x17
	PAIRING_NOT_ALLOWED             = 0x18
	UNKNOWN_LMP_PDU                 = 0x19
	UNSUPPORTED_REMOTE_FEATURE      = 0x1a
	SCO_OFFSET_REJECTED             = 0x1b
	SCO_INTERVAL_REJECTED           = 0x1c
	AIR_MODE_REJECTED               = 0x1d
	INVALID_LMP_PARAMETERS          = 0x1e
	UNSPECIFIED_ERROR               = 0x1f
	UNSUPPORTED_LMP_PARAMETER_VALUE = 0x20
	ROLE_CHANGE_NOT_ALLOWED         = 0x21
	LMP_RESPONSE_TIMEOUT            = 0x22
	LMP_ERROR_TRANSACTION_COLLISION = 0x23
	LMP_PDU_NOT_ALLOWED             = 0x24
	ENCRYPTION_MODE_NOT_ACCEPTED    = 0x25
	UNIT_LINK_KEY_USED              = 0x26
	QOS_NOT_SUPPORTED               = 0x27
	INSTANT_PASSED                  = 0x28
	PAIRING_NOT_SUPPORTED           = 0x29
	TRANSACTION_COLLISION           = 0x2a
	QOS_UNACCEPTABLE_PARAMETER      = 0x2c
	QOS_REJECTED                    = 0x2d
	CLASSIFICATION_NOT_SUPPORTED    = 0x2e
	INSUFFICIENT_SECURITY           = 0x2f
	PARAMETER_OUT_OF_RANGE          = 0x30
	ROLE_SWITCH_PENDING             = 0x32
	SLOT_VIOLATION                  = 0x34
	ROLE_SWITCH_FAILED              = 0x35
	EIR_TOO_LARGE                   = 0x36
	SIMPLE_PAIRING_NOT_SUPPORTED    = 0x37
	HOST_BUSY_PAIRING               = 0x38
)

// ACL flags
const (
	ACL_START_NO_FLUSH = 0x00
	ACL_CONT           = 0x01
	ACL_START          = 0x02
	ACL_ACTIVE_BCAST   = 0x04
	ACL_PICO_BCAST     = 0x08
)

// Baseband links
const (
	SCO_LINK  = 0x00
	ACL_LINK  = 0x01
	ESCO_LINK = 0x02
)

type Conn struct {
	fd int
}

func Dial(network, address string) (*Conn, error) {
	var dev, channel uint16
	n, _ := fmt.Sscanf(address, "%d:%d", &dev, &channel)
	if n != 2 {
		return nil, &net.OpError{Op: "dial", Err: errors.New("invalid address")}
	}

	fd, err := unix.Socket(unix.AF_BLUETOOTH, unix.SOCK_RAW, unix.BTPROTO_HCI)
	if err != nil {
		return nil, &net.OpError{Op: "dial", Err: err}
	}

	err = unix.Bind(fd, &unix.SockaddrHCI{
		Dev:     dev,
		Channel: channel,
	})
	if err != nil {
		unix.Close(fd)
		return nil, err
	}

	return &Conn{
		fd: fd,
	}, nil
}

func (c *Conn) Fd() uintptr {
	return uintptr(c.fd)
}

func (c *Conn) Close() error {
	return unix.Close(c.fd)
}
