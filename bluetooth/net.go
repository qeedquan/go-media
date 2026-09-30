package bluetooth

import (
	"errors"
	"net"

	"github.com/qeedquan/go-media/bluetooth/btl2cap"
)

func Listen(network, address string, option interface{}) (net.Listener, error) {
	switch network {
	case "l2cap":
		return btl2cap.Listen(network, address, option.(*btl2cap.ListenerOption))
	}
	return nil, &net.OpError{Op: "listen", Err: errors.New("unsupported protocol")}
}
