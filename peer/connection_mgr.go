package peer

import (
	"sync"
	"time"

	"github.com/ellcrys/druid/util"
	net "github.com/libp2p/go-libp2p-net"
	ma "github.com/multiformats/go-multiaddr"
)

// ConnectionEstTickerTime is the interval between attempts
// to establish connections with peers
var ConnectionEstTickerTime = 10 * time.Second

// ConnectionManager manages the active connections
// ensuring the required number of connections at any given
// time is maintained
type ConnectionManager struct {
	gmx        *sync.Mutex
	pm         *Manager
	activeConn int
	connEstInt *time.Ticker
}

// NewConnMrg creates a new connection manager
func NewConnMrg(m *Manager) *ConnectionManager {
	return &ConnectionManager{
		pm:  m,
		gmx: &sync.Mutex{},
	}
}

// Manage starts connection management
func (m *ConnectionManager) Manage() {
	go m.establishConnections()
}

// connectionCount returns the number of active connections
func (m *ConnectionManager) connectionCount() int {
	m.gmx.Lock()
	defer m.gmx.Unlock()
	return m.activeConn
}

// establishConnections will attempt to send a handshake to
// addresses that have not been connected to as long as the max
// connection limit has not been reached
func (m *ConnectionManager) establishConnections() {
	m.connEstInt = time.NewTicker(ConnectionEstTickerTime)
	for {
		select {
		case <-m.connEstInt.C:
			if m.pm.NeedMorePeers() {
				for _, p := range m.pm.getUnconnectedPeers() {
					go m.pm.establishConnection(p.StringID())
				}
			}
		}
	}
}

// Listen is called when network starts listening on an address
func (m *ConnectionManager) Listen(net.Network, ma.Multiaddr) {

}

// ListenClose is called when network stops listening on an address
func (m *ConnectionManager) ListenClose(net.Network, ma.Multiaddr) {

}

// Connected is called when a connection is opened
func (m *ConnectionManager) Connected(net net.Network, conn net.Conn) {
	m.gmx.Lock()
	defer m.gmx.Unlock()
	m.activeConn++
}

// Disconnected is called when a connection is closed
func (m *ConnectionManager) Disconnected(net net.Network, conn net.Conn) {
	fullAddr := util.FullRemoteAddressFromConn(conn)
	go m.pm.onPeerDisconnect(fullAddr)

	m.gmx.Lock()
	defer m.gmx.Unlock()
	m.activeConn--
}

// OpenedStream is called when a stream is openned
func (m *ConnectionManager) OpenedStream(net.Network, net.Stream) {

}

// ClosedStream is called when a stream is openned
func (m *ConnectionManager) ClosedStream(nt net.Network, s net.Stream) {
}
