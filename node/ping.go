package node

import (
	"bufio"
	"context"
	"fmt"
	"time"

	"github.com/ellcrys/druid/wire"

	"github.com/ellcrys/druid/util"
	net "github.com/libp2p/go-libp2p-net"
	pc "github.com/multiformats/go-multicodec/protobuf"
)

func (pt *Inception) sendPing(remotePeer *Node) error {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	remotePeerIDShort := remotePeer.ShortID()
	s, err := pt.LocalPeer().addToPeerStore(remotePeer).newStream(ctx, remotePeer.ID(), util.PingVersion)
	if err != nil {
		pt.log.Debug("Ping failed. failed to connect to peer", "Err", err, "PeerID", remotePeerIDShort)
		return fmt.Errorf("ping failed. failed to connect to peer. %s", err)
	}
	defer s.Close()

	w := bufio.NewWriter(s)
	msg := &wire.Ping{}
	if err := pc.Multicodec(nil).Encoder(w).Encode(msg); err != nil {
		s.Reset()
		pt.log.Debug("ping failed. failed to write to stream", "Err", err, "PeerID", remotePeerIDShort)
		return fmt.Errorf("ping failed. failed to write to stream")
	}
	w.Flush()

	pt.log.Info("Sent ping to peer", "PeerID", remotePeerIDShort)

	// receive pong response
	pongMsg := &wire.Pong{}
	decoder := pc.Multicodec(nil).Decoder(bufio.NewReader(s))
	if err := decoder.Decode(pongMsg); err != nil {
		s.Reset()
		pt.log.Debug("Failed to read pong response", "Err", err, "PeerID", remotePeerIDShort)
		return fmt.Errorf("failed to read pong response")
	}

	remotePeer.Timestamp = time.Now()
	pt.PM().AddOrUpdatePeer(remotePeer)

	pt.log.Info("Received pong response from peer", "PeerID", remotePeerIDShort)

	return nil
}

// SendPing sends a ping message
func (pt *Inception) SendPing(remotePeers []*Node) {
	pt.log.Info("Sending ping to peer(s)", "NumPeers", len(remotePeers))
	for _, remotePeer := range remotePeers {
		_remotePeer := remotePeer
		go func() {
			if err := pt.sendPing(_remotePeer); err != nil {
				pt.PM().onFailedConnection(_remotePeer)
			}
		}()
	}
}

// OnPing handles incoming ping message
func (pt *Inception) OnPing(s net.Stream) {

	defer s.Close()
	remotePeer := NewRemoteNode(util.FullRemoteAddressFromStream(s), pt.LocalPeer())
	remotePeerIDShort := remotePeer.ShortID()

	pt.log.Info("Received ping message", "PeerID", remotePeerIDShort)

	msg := &wire.Ping{}
	if err := pc.Multicodec(nil).Decoder(bufio.NewReader(s)).Decode(msg); err != nil {
		s.Reset()
		pt.log.Error("failed to read ping message", "Err", err, "PeerID", remotePeerIDShort)
		return
	}

	// send pong message
	pongMsg := &wire.Pong{}
	w := bufio.NewWriter(s)
	enc := pc.Multicodec(nil).Encoder(w)
	if err := enc.Encode(pongMsg); err != nil {
		pt.log.Error("failed to send pong response", "Err", err)
		return
	}

	remotePeer.Timestamp = time.Now()
	pt.PM().AddOrUpdatePeer(remotePeer)
	pt.log.Info("Sent pong response to peer", "PeerID", remotePeerIDShort)

	w.Flush()
}