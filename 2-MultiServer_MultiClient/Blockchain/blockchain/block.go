package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/tubar2/go_MultiServer_Blockchain/Blockchain/blockchainpb"
)

//
// MARK: Block struct
//
type Block struct {
	Index         int64         `json:"index"`
	Proof         int64         `json:"proof"`
	Previous_hash string        `json:"previous_hash"`
	Timestamp     string        `json:"timestamp"`
	Transactions  []Transaction `json:"transactions"`
}

//
// note: Methods
//

// Returns a sha256 hex string representation of a block in json format
func (block Block) Hash() string {
	obj, err := json.Marshal(block)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", sha256.Sum256([]byte(obj)))
}

//
// note: Formatting From and To blockpb
//

// Returns a new blockchainpb_block from a block object
func (block *Block) ToBlockpbfmt() *blockchainpb.Block {
	var _transactionspb []*blockchainpb.Transaction
	for _, transaction := range block.Transactions {
		_transactionspb = append(_transactionspb, transaction.ToTransactionpbfmt())
	}

	return &blockchainpb.Block{
		Index:        block.Index,
		Proof:        block.Proof,
		PreviousHash: block.Previous_hash,
		Timestamp:    block.Timestamp,
		Transactions: _transactionspb,
	}
}

// Returns a new block object from a blockcainpb_block
func FromBlockpbfmt(block *blockchainpb.Block) Block {
	return Block{
		Index:         block.GetIndex(),
		Proof:         block.GetProof(),
		Previous_hash: block.GetPreviousHash(),
		Timestamp:     block.GetTimestamp(),
		Transactions:  FromTransactionpbArray(block.GetTransactions()),
	}
}

//
// MARK: BlockArray struct
//
type BlockArray struct {
	Blocks []Block
}

//
// note: Methods
//

// Returns block at position specified
func (blockArr *BlockArray) At(position int) *Block {
	return &blockArr.Blocks[position]
}

// Returns block at position specified
func (blockArr *BlockArray) Len() int {
	return len(blockArr.Blocks)
}

func (blockArr *BlockArray) Append(block Block) {
	blockArr.Blocks = append(blockArr.Blocks, block)
}

//
// note: Formatting From and To blockarraypb
//

func (blockArray *BlockArray) ToBlockpbArrayfmt() []*blockchainpb.Block {
	var blockspb []*blockchainpb.Block

	for _, block := range blockArray.Blocks {
		blockspb = append(blockspb, block.ToBlockpbfmt())
	}

	return blockspb
}

func FromBlockpbArrayfmt(blockspb []*blockchainpb.Block) BlockArray {
	var blocks BlockArray
	for _, blockpb := range blockspb {
		block := FromBlockpbfmt(blockpb)
		blocks.Blocks = append(blocks.Blocks, block)
	}

	return blocks
}
