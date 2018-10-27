package gossip_test

import (
	"testing"

	. "github.com/ellcrys/elld/blockchain/testutil"
	"github.com/ellcrys/elld/crypto"
	"github.com/ellcrys/elld/node"
	"github.com/ellcrys/elld/node/gossip"
	"github.com/ellcrys/elld/params"
	"github.com/ellcrys/elld/types"
	"github.com/ellcrys/elld/types/core"
	"github.com/ellcrys/elld/util"
	. "github.com/ncodes/goblin"
	"github.com/olebedev/emitter"
	. "github.com/onsi/gomega"
	"github.com/shopspring/decimal"
)

func TestBlock(t *testing.T) {
	params.FeePerByte = decimal.NewFromFloat(0.01)
	g := Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("Block", func() {

		var lp, rp *node.Node
		var sender, _ = crypto.NewKey(nil)
		var receiver, _ = crypto.NewKey(nil)
		var lpPort, rpPort int

		g.BeforeEach(func() {
			lpPort = getPort()
			rpPort = getPort()

			lp = makeTestNode(lpPort)
			Expect(lp.GetBlockchain().Up()).To(BeNil())

			rp = makeTestNode(rpPort)
			Expect(rp.GetBlockchain().Up()).To(BeNil())

			// Create sender account on the remote peer
			Expect(rp.GetBlockchain().CreateAccount(1, rp.GetBlockchain().GetBestChain(), &core.Account{
				Type:    core.AccountTypeBalance,
				Address: util.String(sender.Addr()),
				Balance: "100",
			})).To(BeNil())

			// Create sender account on the local peer
			Expect(lp.GetBlockchain().CreateAccount(1, lp.GetBlockchain().GetBestChain(), &core.Account{
				Type:    core.AccountTypeBalance,
				Address: util.String(sender.Addr()),
				Balance: "100",
			})).To(BeNil())
		})

		g.AfterEach(func() {
			closeNode(lp)
			closeNode(rp)
		})

		g.Describe(".RelayBlock", func() {
			g.Context("when a block has a transaction that is not in the txs pool of the remote peer", func() {
				g.Context("when block is successfully relayed to a remote peer", func() {

					var block types.Block

					g.BeforeEach(func() {
						block = MakeBlockWithSingleTx(lp.GetBlockchain(), lp.GetBlockchain().GetBestChain(), sender, sender, 1)
					})

					g.Specify("relayed block must be processed by the remote peer", func(done Done) {

						err := lp.Gossip().RelayBlock(block, []core.Engine{rp})
						Expect(err).To(BeNil())

						go func() {
							evtArgs := <-rp.GetEventEmitter().Once(gossip.EventBlockProcessed)
							Expect(evtArgs.Args).To(HaveLen(2))
							_, err := evtArgs.Args[0], evtArgs.Args[1]
							Expect(err).ToNot(BeNil())
							Expect(err.(error).Error()).To(Equal("tx:0, error:transaction does not exist in the transactions pool"))
							done()
						}()
					})
				})
			})

			g.Context("when a block is successfully relayed to a remote peer", func() {

				var evtArgs emitter.Event
				var block types.Block

				g.BeforeEach(func() {
					block = MakeBlockWithSingleTx(lp.GetBlockchain(), lp.GetBlockchain().GetBestChain(), sender, sender, 1)

					// Add the transaction to the remote node's
					// transaction pool to prevent block rejection
					// because of the tx being unknown.
					err := rp.GetTxPool().Put(block.GetTransactions()[0])
					Expect(err).To(BeNil())
				})

				g.Specify("relayed block must be processed by the remote peer", func(done Done) {
					err := lp.Gossip().RelayBlock(block, []core.Engine{rp})
					Expect(err).To(BeNil())

					go func() {
						evtArgs = <-rp.GetEventEmitter().Once(gossip.EventBlockProcessed)
						Expect(evtArgs.Args).To(HaveLen(2))
						processedBlock, err := evtArgs.Args[0], evtArgs.Args[1]
						Expect(err).To(BeNil())
						Expect(processedBlock).ToNot(BeNil())
						Expect(block.GetHashAsHex()).To(Equal(processedBlock.(types.Block).GetHashAsHex()))

						rpCurBlock, err := rp.GetBlockchain().ChainReader().Current()
						Expect(err).To(BeNil())
						Expect(rpCurBlock.GetNumber()).To(Equal(block.GetNumber()))

						done()
					}()
				})
			})

			g.Context("when relayed block is considered an orphan by the remote peer", func() {

				var block2, block3 types.Block
				var evtOrphanBlock = make(chan emitter.Event)

				g.BeforeEach(func() {
					block2 = MakeBlockWithSingleTx(lp.GetBlockchain(), lp.GetBlockchain().GetBestChain(), sender, sender, 1)
					_, err := lp.GetBlockchain().ProcessBlock(block2)
					Expect(err).To(BeNil())

					block3 = MakeBlockWithSingleTx(lp.GetBlockchain(), lp.GetBlockchain().GetBestChain(), sender, sender, 2)
					_, err = lp.GetBlockchain().ProcessBlock(block3)
					Expect(err).To(BeNil())
				})

				g.It("should emit core.EventOrphanBlock", func(done Done) {
					err := lp.Gossip().RelayBlock(block3, []core.Engine{rp})
					Expect(err).To(BeNil())

					go func() {
						evt := <-rp.GetBlockchain().GetEventEmitter().Once(core.EventOrphanBlock)
						evtOrphanBlock <- evt
					}()

					go func() {
						<-rp.GetEventEmitter().Once(gossip.EventBlockProcessed)

						evt := <-evtOrphanBlock
						Expect(evt).ToNot(BeNil())
						Expect(evt.Args).To(HaveLen(1))
						orphanBlock := evt.Args[0].(*core.Block)

						Expect(orphanBlock.GetNumber()).To(Equal(block3.GetNumber()))
						Expect(orphanBlock.GetBroadcaster().StringID()).To(Equal(lp.StringID()))
						Expect(rp.GetBlockchain().OrphanBlocks().Len()).To(Equal(1))
						done()
					}()
				})
			})
		})

		g.Describe(".RequestBlock", func() {

			var block2 types.Block

			g.BeforeEach(func() {
				block2 = MakeBlockWithSingleTx(lp.GetBlockchain(), lp.GetBlockchain().GetBestChain(), sender, sender, 1)
				_, err := lp.GetBlockchain().ProcessBlock(block2)
				Expect(err).To(BeNil())
			})

			g.It("Should request and process block from remote peer", func() {
				err := rp.Gossip().RequestBlock(lp, block2.GetHash())
				Expect(err).To(BeNil())
				curBlock, err := rp.GetBlockchain().ChainReader().Current()
				Expect(err).To(BeNil())
				Expect(curBlock.GetHashAsHex()).To(Equal(block2.GetHashAsHex()))
			})
		})

		g.Describe(".SendGetBlockHashes", func() {

			var block2, block3 types.Block

			// Target shape:
			// Remote Peer
			// [1]-[2]-[3]
			//
			// Local Peer
			// [1]
			g.Context("when remote blockchain shape is [1]-[2]-[3] and local blockchain shape: [1]", func() {

				g.BeforeEach(func() {
					block2 = MakeBlockWithSingleTx(rp.GetBlockchain(), rp.GetBlockchain().GetBestChain(), sender, receiver, 1)
					_, err := rp.GetBlockchain().ProcessBlock(block2)
					Expect(err).To(BeNil())

					block3 = MakeBlockWithSingleTx(rp.GetBlockchain(), rp.GetBlockchain().GetBestChain(), sender, receiver, 2)
					_, err = rp.GetBlockchain().ProcessBlock(block3)
					Expect(err).To(BeNil())

					err = lp.Gossip().SendGetBlockHashes(rp, nil)
					Expect(err).To(BeNil())
				})

				g.It("should get 2 block hashes from remote peer", func() {
					Expect(lp.GetBlockHashQueue().Size()).To(Equal(2))
				})

				g.Specify("first header number = 2 and second header number = 3", func() {
					h2 := lp.GetBlockHashQueue().Shift()
					h3 := lp.GetBlockHashQueue().Shift()
					Expect(h2).To(BeAssignableToTypeOf(&core.BlockHash{}))
					Expect(h2.(*core.BlockHash).Hash).To(Equal(block2.GetHash()))
					Expect(h3).To(BeAssignableToTypeOf(&core.BlockHash{}))
					Expect(h3.(*core.BlockHash).Hash).To(Equal(block3.GetHash()))
				})
			})

			// Target shape:
			// Remote Peer
			// [1]-[2]
			//
			// Local Peer
			// [1]
			g.Context("when remote blockchain shape is [1]-[2] and local peer blockchain shape is [1]", func() {

				var block2 types.Block

				g.BeforeEach(func() {
					block2 = MakeBlockWithSingleTx(rp.GetBlockchain(), rp.GetBlockchain().GetBestChain(), sender, receiver, 1)
					_, err := rp.GetBlockchain().ProcessBlock(block2)
					Expect(err).To(BeNil())

					err = lp.Gossip().SendGetBlockHashes(rp, nil)
					Expect(err).To(BeNil())
				})

				g.It("should get 1 block hash from remote peer", func() {
					Expect(lp.GetBlockHashQueue().Size()).To(Equal(1))
				})

				g.Specify("the block hash must match hash of block [2]", func() {
					hash := lp.GetBlockHashQueue().Shift()
					Expect(hash).To(BeAssignableToTypeOf(&core.BlockHash{}))
					Expect(hash.(*core.BlockHash).Hash).To(Equal(block2.GetHash()))
				})
			})

			// Target shape:
			// Remote Peer
			// [1]
			//
			// Local Peer
			// [1]
			g.Context("when remote peer's blockchain shape is [1] and local peer's blockchain shape is [1]", func() {

				g.BeforeEach(func() {
					err := lp.Gossip().SendGetBlockHashes(rp, nil)
					Expect(err).To(BeNil())
				})

				g.Specify("local peer block hash queue should be empty", func() {
					Expect(lp.GetBlockHashQueue().Size()).To(Equal(0))
				})
			})

			// Target shape:
			// Remote Peer
			// [1]-[2]-[3] 	ChainA
			//  |__[2]		ChainB
			g.Context("when remote peer's blockchain shape is: [1]-[2]-[3] and local peer's blockchain shape is: [2]", func() {

				var block2, block3, chainBBlock2 types.Block

				g.BeforeEach(func() {
					block2 = MakeBlockWithSingleTx(rp.GetBlockchain(), rp.GetBlockchain().GetBestChain(), sender, receiver, 1)
					chainBBlock2 = MakeBlockWithSingleTx(rp.GetBlockchain(), rp.GetBlockchain().GetBestChain(), sender, receiver, 1)

					_, err := rp.GetBlockchain().ProcessBlock(block2)
					Expect(err).To(BeNil())

					// processing chainBBlock2 will create a fork
					_, err = rp.GetBlockchain().ProcessBlock(chainBBlock2)
					Expect(err).To(BeNil())
					Expect(rp.GetBlockchain().GetChainsReader()).To(HaveLen(2))

					block3 = MakeBlockWithSingleTx(rp.GetBlockchain(), rp.GetBlockchain().GetBestChain(), sender, receiver, 2)
					_, err = rp.GetBlockchain().ProcessBlock(block3)
					Expect(err).To(BeNil())
				})

				g.When("locator hash is block [2] of Chain B", func() {

					g.BeforeEach(func() {
						err := lp.Gossip().SendGetBlockHashes(rp, []util.Hash{chainBBlock2.GetHash()})
						Expect(err).To(BeNil())
					})

					g.Specify("the block hash queue should contain 2 hashes from ChainA", func() {
						Expect(lp.GetBlockHashQueue().Size()).To(Equal(2))
						g.Context("first header number = [2] and second header number = [3]", func() {
							h2 := lp.GetBlockHashQueue().Shift()
							h3 := lp.GetBlockHashQueue().Shift()
							Expect(h2).To(BeAssignableToTypeOf(&core.BlockHash{}))
							Expect(h2.(*core.BlockHash).Hash).To(Equal(block2.GetHash()))
							Expect(h3).To(BeAssignableToTypeOf(&core.BlockHash{}))
							Expect(h3.(*core.BlockHash).Hash).To(Equal(block3.GetHash()))
						})
					})
				})
			})

			g.Context("when no locator/block hash is shared with the remote peer", func() {

				g.BeforeEach(func() {
					err := lp.Gossip().SendGetBlockHashes(rp, []util.Hash{util.StrToHash("unknown")})
					Expect(err).To(BeNil())
				})

				g.It("local peer's block hash queue should be empty", func() {
					Expect(lp.GetBlockHashQueue().Size()).To(Equal(0))
				})
			})
		})

		g.Describe(".SendGetBlockBodies", func() {

			var block2, block3 types.Block

			// Target shape:
			// Remote Peer
			// [1]-[2]-[3]
			//
			// Local Peer
			// [1]
			g.Context("when remote peer's blockchain shape is [1]-[2]-[3] and local peer's blockchain shape is: [1]", func() {
				g.Context("at least one hash is requested", func() {
					g.BeforeEach(func() {
						block2 = MakeBlockWithSingleTx(rp.GetBlockchain(), rp.GetBlockchain().GetBestChain(), sender, sender, 1)
						_, err := rp.GetBlockchain().ProcessBlock(block2)
						Expect(err).To(BeNil())

						block3 = MakeBlockWithSingleTx(rp.GetBlockchain(), rp.GetBlockchain().GetBestChain(), sender, sender, 2)
						_, err = rp.GetBlockchain().ProcessBlock(block3)
						Expect(err).To(BeNil())
					})

					g.It("should successfully fetch block [2] and [3] and append to local peer's chain", func() {
						hashes := []util.Hash{block2.GetHash(), block3.GetHash()}
						err := lp.Gossip().SendGetBlockBodies(rp, hashes)
						Expect(err).To(BeNil())

						lpTip, err := lp.GetBlockchain().ChainReader().Current()
						Expect(err).To(BeNil())
						Expect(lpTip.GetNumber()).To(Equal(uint64(3)))
					})
				})

				g.Context("no hash is requested", func() {
					g.BeforeEach(func() {
						err := lp.Gossip().SendGetBlockBodies(rp, []util.Hash{})
						Expect(err).To(BeNil())
					})

					g.Specify("local chain tip should remain unchanged", func() {
						lpTip, err := lp.GetBlockchain().ChainReader().Current()
						Expect(err).To(BeNil())
						Expect(lpTip.GetNumber()).To(Equal(uint64(1)))
					})
				})
			})
		})
	})
}
