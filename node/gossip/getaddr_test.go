package gossip_test

import (
	"time"

	"github.com/ellcrys/elld/crypto"
	"github.com/ellcrys/elld/node"
	"github.com/ellcrys/elld/types/core"
	"github.com/ellcrys/elld/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetAddr", func() {

	var lp, rp *node.Node
	var sender, _ = crypto.NewKey(nil)
	var lpPort, rpPort int

	BeforeEach(func() {
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

	AfterEach(func() {
		closeNode(lp)
		closeNode(rp)
	})

	Describe(".SendGetAddrToPeer", func() {

		It("should return error if connection fail", func() {
			_, err := rp.Gossip().SendGetAddrToPeer(rp)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("dial to self attempted"))
		})

		Context("when a remote peer knowns an address that is the same as the requesting peer", func() {

			BeforeEach(func() {
				rp.PM().AddOrUpdateNode(lp)
			})

			It("should not be returned", func() {
				addrs, err := lp.Gossip().SendGetAddrToPeer(rp)
				Expect(err).To(BeNil())
				Expect(addrs).To(HaveLen(0))
			})
		})

		Context("when a remote peer knowns an address that is valid", func() {
			var remoteAddr *node.Node

			BeforeEach(func() {
				remoteAddr = makeTestNode(getPort())
				remoteAddr.SetLastSeen(time.Now())
				rp.PM().AddOrUpdateNode(remoteAddr)
			})

			AfterEach(func() {
				closeNode(remoteAddr)
			})

			It("should return the address", func() {
				addrs, err := lp.Gossip().SendGetAddrToPeer(rp)
				Expect(err).To(BeNil())
				Expect(addrs).To(HaveLen(1))
				Expect(addrs[0].Address).To(Equal(remoteAddr.GetAddress()))
			})
		})

		Context("when a remote peer knowns an address that is hardcoded", func() {
			var remoteAddr *node.Node

			BeforeEach(func() {
				remoteAddr = makeTestNode(getPort())
				remoteAddr.SetHardcoded(true)
				remoteAddr.SetLastSeen(time.Now())
				rp.PM().AddOrUpdateNode(remoteAddr)
			})

			AfterEach(func() {
				closeNode(remoteAddr)
			})

			It("should not return the address", func() {
				addrs, err := lp.Gossip().SendGetAddrToPeer(rp)
				Expect(err).To(BeNil())
				Expect(addrs).To(HaveLen(0))
			})
		})

		Context("when a remote peer returns addresses greater than MaxAddrsExpected", func() {
			var remoteAddr *node.Node

			BeforeEach(func() {
				lp.GetCfg().Node.MaxAddrsExpected = 0
				remoteAddr = makeTestNode(getPort())
				remoteAddr.SetLastSeen(time.Now())
				rp.PM().AddOrUpdateNode(remoteAddr)
			})

			AfterEach(func() {
				closeNode(remoteAddr)
			})

			It("should not return the address", func() {
				addrs, err := lp.Gossip().SendGetAddrToPeer(rp)
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("too many addresses received. Ignoring addresses"))
				Expect(addrs).To(HaveLen(0))
			})
		})
	})

})
