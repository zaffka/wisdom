package serve

import (
	"net"
)

//go:generate mockery --name ConnHandler
type ConnHandler interface {
	net.Conn
}
