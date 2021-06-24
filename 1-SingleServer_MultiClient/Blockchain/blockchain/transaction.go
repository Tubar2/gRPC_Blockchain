package blockchain

import "github.com/tubar2/go_Block-Chain/Blockchain/blockchainpb"

type Transaction struct {
	UUID     string  `json:"uuid"`
	Sender   string  `json:"sender"`
	Receiver string  `json:"receiver"`
	Amount   float64 `json:"amount"`
}

func (transaction *Transaction) ToTransactionpbfmt() *blockchainpb.Transaction {
	return &blockchainpb.Transaction{
		Uuid:     transaction.UUID,
		Receiver: transaction.Receiver,
		Sender:   transaction.Sender,
		Amount:   transaction.Amount,
	}
}

func FromTransactionpbfmt(transaction *blockchainpb.Transaction) Transaction {
	return Transaction{
		UUID:     transaction.GetUuid(),
		Sender:   transaction.GetSender(),
		Receiver: transaction.GetReceiver(),
		Amount:   transaction.GetAmount(),
	}
}

func FromTransactionpbArray(transactionspb []*blockchainpb.Transaction) []Transaction {
	var transactions []Transaction
	for _, transaction := range transactionspb {
		transactions = append(transactions, FromTransactionpbfmt(transaction))
	}

	return transactions
}
