package types

import (
	"github.com/ellcrys/elld/wire"
	net "github.com/libp2p/go-libp2p-net"
)

// Gossip defines an interface for a gossip protocol
type Gossip interface {
	SendHandshake(Engine) error
	OnHandshake(net.Stream)
	SendPing([]Engine)
	OnPing(net.Stream)
	SendGetAddr([]Engine) error
	OnGetAddr(net.Stream)
	OnAddr(net.Stream)
	RelayAddr([]*wire.Address) []error
	SelfAdvertise([]Engine) int
	OnTx(net.Stream)
	RelayTx(*wire.Transaction, []Engine) error
}
