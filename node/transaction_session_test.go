package node

import (
	evbus "github.com/asaskevich/EventBus"
	"github.com/ellcrys/elld/crypto"
	"github.com/ellcrys/elld/testutil"
	"github.com/ellcrys/elld/util/logger"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TransactionSession", func() {

	var err error
	var log = logger.NewLogrusNoOp()
	var n *Node

	BeforeEach(func() {
		var err error
		cfg, err = testutil.SetTestCfg()
		Expect(err).To(BeNil())
	})

	BeforeEach(func() {
		n, err = NewNode(cfg, "127.0.0.1:40001", crypto.NewKeyFromIntSeed(1), log)
		Expect(err).To(BeNil())
		n.SetLogicBus(evbus.New())
	})

	AfterEach(func() {
		Expect(testutil.RemoveTestCfgDir()).To(BeNil())
	})

	Describe(".HasTxSession", func() {
		It("should return false when txId is not in the session map", func() {
			Expect(n.HasTxSession("some_id")).To(BeFalse())
		})
	})

	Describe(".AddTxSession", func() {
		It("should successfully add txId to the session map", func() {
			n.AddTxSession("my_id")
			Expect(n.openTransactionsSession).To(HaveKey("my_id"))
			Expect(n.HasTxSession("my_id")).To(BeTrue())
		})
	})

	Describe(".RemoveTxSession", func() {
		It("should successfully remove txId from the session map", func() {
			n.AddTxSession("my_id")
			Expect(n.openTransactionsSession).To(HaveKey("my_id"))
			n.RemoveTxSession("my_id")
			Expect(n.openTransactionsSession).ToNot(HaveKey("my_id"))
			Expect(n.HasTxSession("my_id")).To(BeFalse())
		})
	})

	Describe(".CountTxSession", func() {
		It("should return 2", func() {
			n.AddTxSession("my_id")
			n.AddTxSession("my_id_2")
			Expect(n.CountTxSession()).To(Equal(2))
		})
	})
})
