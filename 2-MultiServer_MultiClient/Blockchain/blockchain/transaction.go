package blockchain

import (
	"github.com/tubar2/go_MultiServer_Blockchain/Blockchain/blockchainpb"
)

//
// MARK: Transaction struct
//
type Transaction struct {
	UUID     string  `json:"uuid"`
	Sender   string  `json:"sender"`
	Receiver string  `json:"receiver"`
	Amount   float64 `json:"amount"`
}

//
// note: Formatting From and To transactionpb
//

// Returns a new blockchainpb_transaction from a transaction object
func (transaction *Transaction) ToTransactionpbfmt() *blockchainpb.Transaction {
	return &blockchainpb.Transaction{
		Uuid:     transaction.UUID,
		Receiver: transaction.Receiver,
		Sender:   transaction.Sender,
		Amount:   transaction.Amount,
	}
}

// Returns a new transaction object from a blockcainpb_transaction
func FromTransactionpbfmt(transaction *blockchainpb.Transaction) Transaction {
	return Transaction{
		UUID:     transaction.GetUuid(),
		Sender:   transaction.GetSender(),
		Receiver: transaction.GetReceiver(),
		Amount:   transaction.GetAmount(),
	}
}

//
// MARK: TransactionArray struct
//
type Transactions struct {
	// maps uuid -> Transaction
	TransactionMap map[string]Transaction `json:"transactions"`
}

//
// note: Methods
//

func (transactions *Transactions) AddTransaction(transaction Transaction) bool {
	_, ok := transactions.TransactionMap[transaction.UUID]
	if !ok {
		transactions.TransactionMap[transaction.UUID] = transaction
	}

	return ok
}

func (transactions *Transactions) RemoveTransaction(tId string) bool {
	_, ok := transactions.TransactionMap[tId]
	if ok {
		delete(transactions.TransactionMap, tId)
	}

	return ok
}

func (transactions *Transactions) GetTransactions() map[string]Transaction {
	return transactions.TransactionMap
}

//
// note: Formatting From and To transactionArraypb
//

// Returns a new []blockchainpb_transaction from a transactionArray object
func (transactions *Transactions) ToTransactionpbArrayfmt() []*blockchainpb.Transaction {
	var transactionspb []*blockchainpb.Transaction

	for _, transaction := range transactions.TransactionMap {
		transactionspb = append(transactionspb, transaction.ToTransactionpbfmt())
	}

	return transactionspb
}

// Returns a new TransactionArray object from a []blockcainpb_transaction
func FromTransactionpbArray(transactionspb []*blockchainpb.Transaction) []Transaction {
	var transactions []Transaction
	for _, transactionpb := range transactionspb {
		transaction := FromTransactionpbfmt(transactionpb)
		transactions = append(transactions, transaction)

	}

	return transactions
}

func (transactions *Transactions) ToTransactionArray() []Transaction {
	var transactionsArr []Transaction
	for _, transaction := range transactions.TransactionMap {
		transactionsArr = append(transactionsArr, transaction)
	}

	return transactionsArr
}

//
// note: Functions
//
func NewTransactions() Transactions {
	return Transactions{
		TransactionMap: make(map[string]Transaction),
	}
}
