package blockchain

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/ellcrys/elld/blockchain/common"
	"github.com/ellcrys/elld/util"
	"github.com/ellcrys/elld/wire"
	"github.com/shopspring/decimal"
)

// validateBlock handles block validation. A block that successfully
// passes this validation is considered safe to add to the chain.
func (b *Blockchain) validateBlock(block *wire.Block) error {

	// validate the block
	if err := block.Validate(); err != nil {
		return fmt.Errorf("failed block validation: %s", err)
	}

	// check if the signature of the block is valid and signed by the block creator.
	if err := wire.BlockVerify(block); err != nil {
		return fmt.Errorf("failed block signature verification: %s", err)
	}

	// validate the transaction root
	if block.Header.TransactionsRoot == util.ToHex(ComputeTxsRoot(block.Transactions)) {
		return fmt.Errorf("failed transaction root check")
	}

	return nil
}

// addOp adds a transition operation object to the list of
// operations (ops). It a similar transition to op already exists,
// it will replaced by the new op.
// @ops 	The current list of operations to add to.
// @op 		The operation to be added
// @returns	A new slice of operations with op included
func addOp(ops []common.Transition, op common.Transition) []common.Transition {
	var newOps []common.Transition
	for _, _op := range ops {
		if !_op.Equal(op) {
			newOps = append(newOps, _op)
		}
	}
	return append(newOps, op)
}

// processBalanceTx processes a transaction. It returns series of transition
// operations that must be applied to effect the transaction.
// It accepts the following args:
//
// @tx: 	The transaction
// @ops: 	The list of current operations generated from other transactions of same block as tx.
//			We use ops to check the latest proposed operation of an account initiated by other transactions.
// @returns	A slice of transitions to be applied to the chain state or error if something bad happened.
func (b *Blockchain) processBalanceTx(tx *wire.Transaction, ops []common.Transition, chain *Chain) ([]common.Transition, error) {
	var err error
	var txOps []common.Transition
	var senderAcct, recipientAcct *wire.Account
	var senderAcctBalance = decimal.Zero
	var recipientAcctBalance = decimal.Zero

	// first, we check if we can determine the balances of the sender and recipient accounts
	// from OpNewAccountBalance operations by previous transactions. Because, if an account was
	// updated by an previous transaction, the new balance will be found in the ops list.
	for _, prevOp := range ops {
		// check for balance change for the sender
		if opNewBalance, yes := prevOp.(*common.OpNewAccountBalance); yes && opNewBalance.Address() == tx.From {
			senderAcctBalance, _ = util.StrToDecimal(opNewBalance.Account.Balance)
		}
		// check for balance change for the recipient
		if opNewBalance, yes := prevOp.(*common.OpNewAccountBalance); yes && opNewBalance.Address() == tx.To {
			recipientAcctBalance, _ = util.StrToDecimal(opNewBalance.Account.Balance)
		}
	}

	// find the sender account. Return error if sender account
	// does not exist. This should never happen here as the caller must
	// have validated all transactions in the containing block.
	senderAcct, err = b.GetAccount(chain, tx.From)
	if err != nil {
		return nil, fmt.Errorf("failed to get sender's account: %s", err)
	}

	// if we were unable to learn about the sender's latest balance from the ops list
	// as a result of previous transactions in same block, then we use the current account balance.
	if senderAcctBalance.Equals(decimal.Zero) {
		senderAcctBalance, _ = util.StrToDecimal(senderAcct.Balance)
	}

	// find the account of the recipient. If the recipient account does not
	// exists, then we must create a OpCreateAccount transition to instruct the creation of a new account.
	recipientAcct, err = b.GetAccount(chain, tx.To)
	if err != nil {
		if err != ErrAccountNotFound {
			return nil, fmt.Errorf("failed to retrieve recipient account: %s", err)
		}
		txOps = append(txOps, &common.OpCreateAccount{
			OpBase: &common.OpBase{Addr: tx.To},
			Account: &wire.Account{
				Type:    wire.AccountTypeBalance,
				Address: tx.To,
				Balance: "0",
			},
		})
		recipientAcct = &wire.Account{
			Type:    wire.AccountTypeBalance,
			Address: tx.To,
			Balance: "0",
		}
	}

	// if we are unable to learn about the recipient's latest balance from the ops list as
	// then we can use the balance of the recipient account
	if recipientAcctBalance.Equals(decimal.Zero) {
		recipientAcctBalance, _ = util.StrToDecimal(recipientAcct.Balance)
	}

	// convert the amount to be sent to decimal
	sendingAmount, err := decimal.NewFromString(tx.Value)
	if err != nil {
		return nil, fmt.Errorf("sending amount error: %s", err)
	}

	// ensure the sender's account balance is sufficient for this transaction
	if senderAcctBalance.LessThan(sendingAmount) {
		return nil, fmt.Errorf("insufficient sender account balance")
	}

	// add an operation to set a new account balance for the sender
	senderAcct.Balance = senderAcctBalance.Sub(sendingAmount).String()
	txOps = append(txOps, &common.OpNewAccountBalance{
		OpBase:  &common.OpBase{Addr: tx.From},
		Account: senderAcct,
	})

	// add an operation to set a new balance of the recipient
	recipientAcct.Balance = recipientAcctBalance.Add(sendingAmount).String()
	txOps = append(txOps, &common.OpNewAccountBalance{
		OpBase:  &common.OpBase{Addr: tx.To},
		Account: recipientAcct,
	})

	return txOps, nil
}

// opsToKVObjects takes a slice of operations and apply them to the provided chain
func (b *Blockchain) opsToStateObjects(block *wire.Block, chain *Chain, ops []common.Transition) ([]*common.StateObject, error) {

	stateObjs := []*common.StateObject{}

	for _, op := range ops {
		switch _op := op.(type) {

		case *common.OpCreateAccount:
			stateObjs = append(stateObjs, &common.StateObject{
				Key:   common.MakeAccountKey(block.GetNumber(), chain.id, _op.Address()),
				Value: util.ObjectToBytes(_op.Account),
			})

		case *common.OpNewAccountBalance:
			stateObjs = append(stateObjs, &common.StateObject{
				Key:   common.MakeAccountKey(block.GetNumber(), chain.id, _op.Address()),
				Value: util.ObjectToBytes(_op.Account),
			})

		default:
			return nil, fmt.Errorf("unknown transition sub-type")
		}
	}

	return stateObjs, nil
}

// processTransactions computes the operations that must be applied to the
// hash tree and world state.
func (b *Blockchain) processTransactions(txs []*wire.Transaction, chain *Chain) ([]common.Transition, error) {

	var ops []common.Transition

	// here we will process each transaction and attempt
	// to decide what should happen to the chain state by
	// producing transition objects.
	for _, tx := range txs {
		switch tx.Type {

		case wire.TxTypeBalance:
			newOps, err := b.processBalanceTx(tx, ops, chain)
			if err != nil {
				return nil, err
			}
			for _, op := range newOps {
				ops = addOp(ops, op)
			}
		}
	}

	return ops, nil
}

// ComputeTxsRoot computes the merkle root of a set of transactions.
// Transactions are first lexicographically sorted and added to a
// brand new tree. Returns the tree root.
func ComputeTxsRoot(txs []*wire.Transaction) []byte {

	sort.Slice(txs, func(i, j int) bool {
		iBytes, _ := util.FromHex(txs[i].Hash)
		jBytes, _ := util.FromHex(txs[j].Hash)
		return bytes.Compare(iBytes, jBytes) == -1
	})

	tree := NewMemHashTree(nil, nil)
	for _, tx := range txs {
		tree.Upsert([]byte(tx.GetHash()), []byte(""))
	}

	root, _, _ := tree.Root()
	return root
}

func (b *Blockchain) maybeAcceptBlock(block *wire.Block) error {

	// find the chain where the parent of the block exists on. If a chain is not found,
	// then the block is considered an orphan. If the chain is found but the block at the tip
	// is has the same or a greater block number compared to the new block, it is considered a stale block.
	parentBlock, chain, chainTip, err := b.findBlockChainByHash(block.Header.ParentHash)
	if err != nil {
		b.log.Debug("failed to find chain", "Err", err.Error())
		return err

	} else if chain == nil {
		b.addOrphanBlock(block)
		return nil

	} else if block.Header.Number < chainTip.Number {
		// This is a much older stale block. We only support stale blocks of same height
		// as the current block on the chain.
		b.addRejectedBlock(block)
		return common.ErrVeryStaleBlock

	} else if block.GetNumber() == chainTip.Number {
		// create the new chain, set its root to the parent of the stale block
		// and also add the stale block to it.
		if _, err := b.newChain(block, parentBlock); err != nil {
			return fmt.Errorf("failed to create subtree out of stale block")
		}

		return nil
	} else if block.GetNumber()-chainTip.Number != 1 {
		b.addRejectedBlock(block)
		return common.ErrBlockFailedValidation
	}

	// Mock execute block to derive the state objects and the resulting
	// state root should the state object be applied to the blockchain state tree.
	mockRoot, stateObjs, err := b.mockExecBlock(chain, block)
	if err != nil {
		return err
	}

	// Compare the state root in the block header with the root obtained
	// from the mock execution of the block.
	if block.Header.StateRoot != util.ToHex(mockRoot) {
		return common.ErrBlockStateRootInvalid
	}

	tx := chain.store.NewTx()

	// Next we need to update the blockchain objects in the store
	// as described by the state objects
	for _, so := range stateObjs {
		if err := chain.store.PutWithTx(tx, so.Key, so.Value); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to add state object to store: %s", err)
		}
	}

	// At this point, the block is good to go. We add it to the chain
	if err := chain.appendBlockWithTx(tx, block); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to add block: %s", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return fmt.Errorf("commit error: %s", err)
	}

	return nil
}

// ProcessBlock takes a block and attempts to add it to the
// tip of the blockchain.
func (b *Blockchain) ProcessBlock(block *wire.Block) error {
	b.mLock.Lock()
	defer b.mLock.Unlock()

	b.log.Debug("Processing block", "Hash", block.Hash)

	// validate the content and format of the block as well as the signature.
	// if err := b.validateBlock(block); err != nil {
	// 	return nil
	// }

	// if the block has been previously rejected, return err
	if b.isRejected(block) {
		return common.ErrBlockRejected
	}

	// check if the block has previously been detected as an orphan.
	// We do not need to go re-process this block if it is an orphan.
	if b.isOrphanBlock(block.Hash) {
		return common.ErrOrphanBlock
	}

	// check if the block exists in any known chain
	exists, err := b.HaveBlock(block.Hash)
	if err != nil {
		return fmt.Errorf("failed to check block existence: %s", err)
	}
	if exists {
		b.log.Debug("Block already exists", "Hash", block.Hash)
		return common.ErrBlockExists
	}

	// attempt to add the block to a chain
	if err := b.maybeAcceptBlock(block); err != nil {
		return err
	}

	// process any remaining orphan blocks
	b.processOrphanBlocks(block.Hash)

	return nil
}

// mockExecBlock performs a mock execution of the blocks to output
// the resulting state objects and state root without making permanent
// commits to the current blockchain state. chain is the specific chain
// to perform this execution against.
func (b *Blockchain) mockExecBlock(chain *Chain, block *wire.Block) (root []byte, stateObjs []*common.StateObject, err error) {

	// Process the transactions to produce a series of transitions
	// that must be applied to the blockchain state.
	ops, err := b.processTransactions(block.Transactions, chain)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to process transactions: rejected: %s", err)
	}

	// Create state objects from the transition objects. State objects when written
	// to the blockchain state (store and tree) change the values of data.
	stateObjs, err = b.opsToStateObjects(block, chain, ops)
	if err != nil {
		return nil, nil, err
	}

	// Get the current root hash and node of the chain.
	rootHash, rootNode, err := chain.stateTree.Root()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get root of chain")
	}

	// Create a new memory-backed state tree and set its initial root to the
	// last root information of the chain.
	stateTree := NewMemHashTree(rootHash, rootNode)
	for _, so := range stateObjs {
		if err := stateTree.Upsert(so.Value, so.Value); err != nil {
			return nil, nil, fmt.Errorf("failed to add state object to state hash tree: %s", err)
		}
	}

	if root, _, err = stateTree.Root(); err != nil {
		return nil, nil, err
	}

	return
}

// processOrphanBlocks finds orphan blocks in the cache and attempts
// to re-process the blocks that are parented by the latestBlockHash.
//
// This method is not protected by any lock. It must be called with
// the chain lock held.
func (b *Blockchain) processOrphanBlocks(latestBlockHash string) error {

	// Add the passed block hash to this internal slice. This is
	// the slice we will use to perform repetitive orphan processing
	// without needing to recursively call this method at the end.
	var parentHashes = []string{latestBlockHash}

	// As long as we have parent hashes, We will continuously pick a
	// parent hash and try to find an orphan block that references the parent hash.
	for len(parentHashes) > 0 {
		// pick the next parent hash and remove it from the slice
		curParentHash := parentHashes[0]
		parentHashes[0] = ""
		parentHashes = parentHashes[1:]

		// Retrieve the keys of blocks in the orphan cache.
		// Go through them and attempt to append them to a chain
		// using maybeAcceptBlock.
		orphansKey := b.orphanBlocks.Keys()
		for i := 0; i < len(orphansKey); i++ {

			oBKey := orphansKey[i]

			// find an orphan block with a parent hash that
			// is same has the latestBlockHash
			orphanBlock := b.orphanBlocks.Peek(oBKey).(*wire.Block)
			if orphanBlock.Header.ParentHash != curParentHash {
				continue
			}

			// remove from the orphan from the cache
			b.orphanBlocks.Remove(orphanBlock.GetHash())

			// re-attempt to process the block
			if err := b.maybeAcceptBlock(orphanBlock); err != nil {
				return err
			}

			parentHashes = append(parentHashes, orphanBlock.Hash)
		}
	}

	return nil
}