package blockchain

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"

	"github.com/tubar2/go_MultiServer_Blockchain/Blockchain/blockchain/utils"
	"github.com/tubar2/go_MultiServer_Blockchain/Blockchain/blockchainpb"
)

// MARK: Blockchain
type Blockchain struct {
	chain BlockArray
}

//
// note: Methods
//

// Returns last block from the blockchain
func (blockchain *Blockchain) GetLastBlock() *Block {
	return blockchain.chain.At(blockchain.chain.Len() - 1)
}

// Generates a proof for the blockchain
func (blockchain *Blockchain) ProofOfWork(previousProof int) int {
	newProof := 1

	for {
		hashOperation := getHashFromProof(newProof, previousProof)
		if hashOperation[:4] == "0000" {
			break
		} else {
			newProof += 1
		}
	}
	return newProof
}

// Createas a new block and appends it to the given blockchain
func (blockchain *Blockchain) CreateBlock(proof int, previous_hash string, transactions Transactions) Block {
	block := Block{
		Index:         int64(blockchain.chain.Len()),
		Proof:         int64(proof),
		Previous_hash: previous_hash,
		Timestamp:     time.Now().Format("2006-01-02 15:04:05.058898"),
		Transactions:  transactions.ToTransactionArray(),
	}
	blockchain.chain.Append(block)

	return block
}

func (blockchain *Blockchain) Len() int {
	return blockchain.chain.Len()
}

// Checks if stored blockchain is valid
func (blockchain *Blockchain) IsChainValid() bool {
	blockIndex := 1
	previousBlock := blockchain.chain.At(0)

	for blockIndex < blockchain.chain.Len() {
		block := blockchain.chain.At(blockIndex)
		if block.Previous_hash != previousBlock.Hash() {
			return false
		}

		previousProof := previousBlock.Proof
		proof := block.Proof
		hash_operation := getHashFromProof(int(proof), int(previousProof))
		if hash_operation[:4] != "0000" {
			return false
		}
		previousBlock = block
		blockIndex += 1
	}
	return true
}

//
// note: Formatting From and To blockchainpb
//

// Retuns Blockchainpb
func (blockchain *Blockchain) ToBlockchainpbfmt() *blockchainpb.Blockchain {
	return &blockchainpb.Blockchain{
		Chain: blockchain.chain.ToBlockpbArrayfmt(),
	}
}

// Retuns Blockchain
func FromBlockchainpbfmt(myBlockchainpb *blockchainpb.Blockchain) Blockchain {
	return Blockchain{
		chain: FromBlockpbArrayfmt(myBlockchainpb.GetChain()),
	}
}

//
// MARK: Functions
//

func getHashFromProof(proof, previousProof int) string {
	sum := (proof * proof) - (previousProof * previousProof)
	hashSum := sha256.Sum256([]byte(strconv.Itoa(sum)))
	return fmt.Sprintf("%x", hashSum)
}

// Retuns a completely new blockchain with only it's genesis block
func NewBlockChain() Blockchain {
	newBlockchain := Blockchain{
		chain: BlockArray{
			Blocks: []Block{
				{
					Index:         0,
					Proof:         1,
					Previous_hash: "0",
					Timestamp:     utils.GetCurrentTime(),
					Transactions:  []Transaction{},
				},
			},
		},
	}
	return newBlockchain
}
