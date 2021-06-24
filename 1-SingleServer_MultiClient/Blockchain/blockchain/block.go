package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/tubar2/go_Block-Chain/Blockchain/blockchainpb"
)

type Block struct {
	Index         int64         `json:"index"`
	Proof         int64         `json:"proof"`
	Previous_hash string        `json:"previous_hash"`
	Timestamp     string        `json:"timestamp"`
	Transactions  []Transaction `json:"transactions"`
}

// Return a block understandable by the grpc client
func (block *Block) ToBlockpbfmt() *blockchainpb.Block {
	var transactions []*blockchainpb.Transaction

	for _, transaction := range block.Transactions {
		transactions = append(transactions, transaction.ToTransactionpbfmt())
	}

	return &blockchainpb.Block{
		Index:        block.Index,
		Proof:        block.Proof,
		PreviousHash: block.Previous_hash,
		Timestamp:    block.Timestamp,
		Transactions: transactions,
	}
}

// Returns a sha256 representation of a block in json format
func (block Block) Hash() string {
	obj, err := json.Marshal(block)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", sha256.Sum256([]byte(obj)))
}

func FromBlockpbfmt(block *blockchainpb.Block) *Block {
	return &Block{
		Index:         block.GetIndex(),
		Proof:         block.GetProof(),
		Previous_hash: block.GetPreviousHash(),
		Timestamp:     block.GetTimestamp(),
		Transactions:  FromTransactionpbArray(block.GetTransactions()),
	}
}
