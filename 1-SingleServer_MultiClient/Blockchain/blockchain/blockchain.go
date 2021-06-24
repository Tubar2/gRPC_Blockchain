package blockchain

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type Blockchain struct {
	chain        []Block
	transactions []Transaction
	nodes        map[string]*User
}

func (blockchain *Blockchain) GetUser(uuid string) *User {
	return blockchain.nodes[uuid]
}

func (blockchain *Blockchain) ClearTransactions() {
	blockchain.transactions = []Transaction{}
}

func (blockchain *Blockchain) GetChain() []Block {
	return blockchain.chain
}

func (blockchain *Blockchain) GetNodes() map[string]*User {
	return blockchain.nodes
}

func (blockchain *Blockchain) UserInNodes(user string) bool {
	if _, ok := blockchain.nodes[user]; ok {
		return true
	}
	return false
}

func (blockchain *Blockchain) IsEmptyNodes() bool {
	return blockchain.nodes == nil
}

func (blockchain *Blockchain) IsEmptyblockchain() bool {
	return blockchain == nil
}

// Createas a new block and appends it to the given blockchain
func (blockchain *Blockchain) CreateBlock(proof int, previous_hash string) Block {
	block := Block{
		Index:         int64(len(blockchain.chain)),
		Proof:         int64(proof),
		Previous_hash: previous_hash,
		Timestamp:     time.Now().Format("2006-01-02 15:04:05.058898"),
		Transactions:  blockchain.transactions,
	}
	blockchain.chain = append(blockchain.chain, block)
	for _, transaction := range blockchain.transactions {
		amount := transaction.Amount
		blockchain.nodes[transaction.Sender].Ballance -= amount
		blockchain.nodes[transaction.Receiver].Ballance += amount
	}
	blockchain.ClearTransactions()
	return block
}

// Returns last element from given blockchain
func (blockchain *Blockchain) GetPreviousBlock() *Block {
	return &blockchain.chain[len(blockchain.chain)-1]
}

// Generates a proof for the blockchain
func (blockchain *Blockchain) ProofOfWork(previous_proof int) int {
	new_proof := 1

	for {
		hashOperation := getHashFromProof(new_proof, previous_proof)
		if hashOperation[:4] == "0000" {
			break
		} else {
			new_proof += 1
		}
	}
	return new_proof
}

func (blockchain *Blockchain) IsChainValid() bool {
	blockIndex := 1
	previousBlock := blockchain.chain[0]

	for blockIndex < len(blockchain.chain) {
		block := blockchain.chain[blockIndex]
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

func (blockchain *Blockchain) AddTransaction(sender, receiver string, amount float64) (int64, string) {
	userId := uuid.NewString()
	blockchain.transactions = append(blockchain.transactions, Transaction{
		UUID:     userId,
		Amount:   amount,
		Sender:   sender,
		Receiver: receiver,
	})
	previousBlock := blockchain.GetPreviousBlock()
	return (previousBlock.Index + 1), userId
}

func (blockchain *Blockchain) AddNode(uuid string) {
	if !blockchain.UserInNodes(uuid) {
		blockchain.nodes[uuid] = &User{
			Ballance: 0,
		}
	}
}

// Retuns a completely new blockchani
// with only it's genesis block
func NewBlockChain() *Blockchain {
	newBlockchain := &Blockchain{
		transactions: []Transaction{},
		nodes:        make(map[string]*User),
		chain: []Block{
			{
				Index:         0,
				Proof:         1,
				Previous_hash: "0",
				Timestamp:     time.Now().Format("2006-01-02 15:04:05.058898"),
				Transactions:  []Transaction{},
			},
		},
	}
	return newBlockchain
}

func getHashFromProof(proof, previous_proof int) string {
	sum := (proof * proof) - (previous_proof * previous_proof)
	hashSum := sha256.Sum256([]byte(strconv.Itoa(sum)))
	return fmt.Sprintf("%x", hashSum)
}
