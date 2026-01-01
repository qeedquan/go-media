// +build linux

package btl2cap

import (
	"fmt"
	"net"
	"os"
	"unsafe"

	"github.com/qeedquan/go-media/bluetooth/bth"
	"golang.org/x/sys/unix"
)

const (
	DEFAULT_MTU      = 672
	DEFAULT_FLUSH_TO = 0xFFFF
)

const (
	LM          = 0x03
	LM_MASTER   = 0x0001
	LM_AUTH     = 0x0002
	LM_ENCRYPT  = 0x0004
	LM_TRUSTED  = 0x0008
	LM_RELIABLE = 0x0010
	LM_SECURE   = 0x0020
)

const (
	COMMAND_REJ  = 0x01
	CONN_REQ     = 0x02
	CONN_RSP     = 0x03
	CONF_REQ     = 0x04
	CONF_RSP     = 0x05
	DISCONN_REQ  = 0x06
	DISCONN_RSP  = 0x07
	ECHO_REQ     = 0x08
	ECHO_RSP     = 0x09
	INFO_REQ     = 0x0a
	INFO_RSP     = 0x0b
	CREATE_REQ   = 0x0c
	CREATE_RSP   = 0x0d
	MOVE_REQ     = 0x0e
	MOVE_RSP     = 0x0f
	MOVE_CFM     = 0x10
	MOVE_CFM_RSP = 0x11
)

const (
	FEAT_FLOWCTL    = 0x00000001
	FEAT_RETRANS    = 0x00000002
	FEAT_BIDIR_QOS  = 0x00000004
	FEAT_ERTM       = 0x00000008
	FEAT_STREAMING  = 0x00000010
	FEAT_FCS        = 0x00000020
	FEAT_EXT_FLOW   = 0x00000040
	FEAT_FIXED_CHAN = 0x00000080
	FEAT_EXT_WINDOW = 0x00000100
	FEAT_UCD        = 0x00000200
)

const (
	FC_L2CAP    = 0x02
	FC_CONNLESS = 0x04
	FC_A2MP     = 0x08
)

const (
	OPTIONS  = 0x01
	CONNINFO = 0x02
)

type Options struct {
	OMTU      uint16
	IMTU      uint16
	FlushTo   uint16
	Mode      uint8
	Fcs       uint8
	MaxTx     uint8
	TxWinSize uint16
}

type ConnInfo struct {
	HCIHandle uint16
	DevClass  [3]uint8
}

type Conn struct {
	*os.File
	laddr bth.Addr
	raddr bth.Addr
}

func (c *Conn) LocalAddr() net.Addr {
	return c.laddr
}

func (c *Conn) RemoteAddr() net.Addr {
	return c.raddr
}

type Listener struct {
	fd    int
	laddr bth.Addr
	opt   *ListenerOption
}

type ListenerOption struct {
	Sockaddr unix.SockaddrL2
	Security bth.Security
}

func Listen(network, address string, option *ListenerOption) (*Listener, error) {
	fd, err := unix.Socket(unix.AF_BLUETOOTH, unix.SOCK_SEQPACKET, unix.BTPROTO_L2CAP)
	if err != nil {
		return nil, &net.OpError{Op: "listen", Err: err}
	}

	bdaddr, err := net.ParseMAC(address)
	if err != nil || len(bdaddr) != 6 {
		return nil, &net.AddrError{Err: "invalid bluetooth address", Addr: address}
	}

	sa := option.Sockaddr
	copy(sa.Addr[:], bdaddr)

	err = unix.Bind(fd, &sa)
	if err != nil {
		return nil, &net.OpError{Op: "listen", Err: err}
	}

	if option.Security.Level != 0 {
		_, _, errno := unix.Syscall6(unix.SYS_SETSOCKOPT, uintptr(fd), unix.SOL_BLUETOOTH, bth.SECURITY,
			uintptr(unsafe.Pointer(&option.Security)), unsafe.Sizeof(option.Security), 0)
		if errno != 0 {
			return nil, &net.OpError{Op: "listen", Err: fmt.Errorf("failed to set security level: %v", errno)}
		}
	}

	err = unix.Listen(fd, 1024)
	if err != nil {
		return nil, &net.OpError{Op: "listen", Err: err}
	}

	return &Listener{
		fd:    fd,
		opt:   option,
		laddr: bth.Addr{"l2cap", bdaddr},
	}, nil
}

func (l *Listener) Accept() (net.Conn, error) {
	return l.AcceptL2CAP()
}

func (l *Listener) AcceptL2CAP() (*Conn, error) {
	cfd, sa, err := unix.Accept(l.fd)
	if err != nil {
		return nil, err
	}

	err = unix.SetNonblock(cfd, true)
	if err != nil {
		unix.Close(cfd)
		return nil, err
	}

	l2sa := sa.(*unix.SockaddrL2)
	return &Conn{
		File:  os.NewFile(uintptr(cfd), fmt.Sprintf("l2cap:%d", cfd)),
		laddr: l.laddr,
		raddr: bth.Addr{"l2cap", l2sa.Addr[:]},
	}, nil
}

func (l *Listener) Addr() net.Addr {
	return l.laddr
}

func (l *Listener) Close() error {
	return unix.Close(l.fd)
}
