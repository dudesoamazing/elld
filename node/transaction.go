package node

import (
	"context"
	"fmt"

	"github.com/ellcrys/elld/blockchain"
	"github.com/ellcrys/elld/config"
	"github.com/ellcrys/elld/node/histcache"
	"github.com/ellcrys/elld/types"
	"github.com/ellcrys/elld/types/core"
	"github.com/ellcrys/elld/types/core/objects"
	"github.com/ellcrys/elld/util"
	net "github.com/libp2p/go-libp2p-net"
)

func MakeTxHistoryKey(tx core.Transaction, peer types.Engine) histcache.MultiKey {
	return []interface{}{tx.ID(), peer.StringID()}
}

// addTransaction adds a transaction to the transaction pool.
func (n *Node) addTransaction(tx core.Transaction) error {

	// Validate the transactions
	txValidator := blockchain.NewTxValidator(tx, n.GetTxPool(), n.GetBlockchain())
	if errs := txValidator.Validate(); len(errs) > 0 {
		return errs[0]
	}

	return n.GetTxPool().Put(tx)
}

// OnTx handles incoming transaction message
func (g *Gossip) OnTx(s net.Stream) {
	defer s.Close()

	remotePeer := NewRemoteNode(util.FullRemoteAddressFromStream(s), g.engine)
	remotePeerIDShort := remotePeer.ShortID()

	// read the message
	msg := &objects.Transaction{}
	if err := ReadStream(s, msg); err != nil {
		s.Reset()
		g.log.Error("failed to read tx message", "Err", err, "PeerID", remotePeerIDShort)
		return
	}

	g.log.Info("Received new transaction", "PeerID", remotePeerIDShort)

	// make a key for this transaction to be added to the history cache so we always know
	// when we have processed this transaction in case we see it again.
	historyKey := MakeTxHistoryKey(msg, remotePeer)

	// check if we have an history about this transaction with the remote peer,
	// if no, add the transaction.
	if g.engine.history().Has(historyKey) {
		return
	}

	// AllocCoin transactions are meant to be added by a miner
	// to a block and not relayed like regular transactions.
	if msg.Type == objects.TxTypeAlloc {
		s.Reset()
		g.log.Debug("Refusing to process allocation transaction")
		g.engine.event.Emit(EventTransactionProcessed, fmt.Errorf("unexpected allocation transaction received"))
		return
	}

	// Ignore the transaction if already in our transaction pool
	if g.engine.transactionsPool.Has(msg) {
		return
	}

	// Validate the transaction and check whether
	// it already exists in the transaction pool,
	// main chain and side chains. If so, reject it
	errs := blockchain.NewTxValidator(msg, g.engine.GetTxPool(), g.engine.bchain).Validate()
	if len(errs) > 0 {
		s.Reset()
		g.log.Debug("Transaction is not valid", "Err", errs[0])
		g.engine.event.Emit(EventTransactionProcessed, errs[0])
		return
	}

	// Add the transaction to the transaction pool and wait for error response
	if err := g.engine.addTransaction(msg); err != nil {
		g.log.Error("failed to add transaction to pool", "Err", msg)
		g.engine.event.Emit(EventTransactionProcessed, err)
		return
	}

	// add transaction to the history cache using the key we created earlier
	g.engine.history().Add(historyKey)

	g.engine.event.Emit(EventTransactionProcessed)

	g.log.Info("Added new transaction to pool", "TxID", msg.ID())
}

// RelayTx relays transactions to peers
func (g *Gossip) RelayTx(tx core.Transaction, remotePeers []types.Engine) error {

	txID := tx.ID()
	sent := 0
	g.log.Debug("Relaying transaction to peers", "TxID", txID, "NumPeers", len(remotePeers))
	for _, peer := range remotePeers {

		historyKey := MakeTxHistoryKey(tx, peer)

		// check if we have an history of sending or receiving this transaction
		// from this remote peer. If yes, do not relay
		if g.engine.history().Has(historyKey) {
			continue
		}

		// create a stream to the remote peer
		s, err := g.NewStream(context.Background(), peer, config.TxVersion)
		if err != nil {
			g.log.Debug("Tx message failed. failed to connect to peer", "Err", err, "PeerID", peer.ShortID())
			continue
		}
		defer s.Close()

		// write to the stream
		if err := WriteStream(s, tx); err != nil {
			s.Reset()
			g.log.Debug("Tx message failed. failed to write to stream", "Err", err, "PeerID", peer.ShortID())
			continue
		}

		// add new history
		g.engine.history().Add(historyKey)

		sent++
	}

	g.log.Info("Finished relaying transaction", "TxID", txID, "NumPeersSentTo", sent)
	return nil
}
