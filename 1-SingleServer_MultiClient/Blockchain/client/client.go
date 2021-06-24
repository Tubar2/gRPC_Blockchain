package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/tubar2/go_Block-Chain/Blockchain/blockchain"
	"github.com/tubar2/go_Block-Chain/Blockchain/blockchainpb"
	"google.golang.org/grpc"
)

var (
	clientUserId string
	myBlockchain []blockchain.Block
	network      []string
	scanner      = bufio.NewScanner(os.Stdin)
)

func main() {
	log.Println("Initiating new blockchain client")

	opts := grpc.WithInsecure()

	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalln("Client couldn't connect", err)
	}
	defer cc.Close()

	c := blockchainpb.NewBlockchainClient(cc)

	// Adds current client, to server's network
	addClient(c)
	log.Println("New client appended to network with id:", clientUserId)

	// Adds server's blockchain to current client
	addBlockchain(c)
	// b, err := json.MarshalIndent(myBlockchain, "", " ")
	// if err != nil {
	// 	panic(err)
	// }
	// log.Println("Network at current time:", string(b))

	// Adds server's clients, to user's
	// addUsers(c)
	// log.Println("Created network", network)

	var command int
	run := true
	for run {
		if command, err = shell(getBalance(c)); err != nil {
			fmt.Printf("Invalid command (%d)", command)
			continue
		}
		switch command {
		case 0:
			run = false
		case 1:
			buy(c)
		case 2:
			mineBlockchain(c)
		case 9:
			continue
		}

	}
}

func mineBlockchain(c blockchainpb.BlockchainClient) {
	res, err := c.MineBlock(context.Background(), &blockchainpb.MineBlockRequest{
		Miner: clientUserId,
	})
	if err != nil {
		log.Fatalln("Couldn't mine block", err)
	}
	res.GetBlock()
}

func getBalance(c blockchainpb.BlockchainClient) float64 {
	res, err := c.GetBallance(context.Background(), &blockchainpb.GetBallanceRequest{
		Uuid: clientUserId,
	})
	if err != nil {
		log.Fatalln("Couldn't add user to network", err)
	}
	return res.GetBallance()
}

func addClient(c blockchainpb.BlockchainClient) {
	res, err := c.AddUser(context.Background(), &blockchainpb.AddUserRequest{})
	if err != nil {
		log.Fatalln("Couldn't add user to network", err)
	}
	clientUserId = res.GetUuid()
}

func addBlockchain(c blockchainpb.BlockchainClient) {
	stream, err := c.GetBlockchain(context.Background(), &blockchainpb.GetBlocksRequest{})
	if err != nil {
		log.Fatalln("Couldn't add user to network", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln("Couldn't add blockchain to client", err)
		}
		data := res.GetBlock()
		var transactions []blockchain.Transaction
		for _, transactionpb := range data.GetTransactions() {
			transactions = append(transactions, blockchain.Transaction{
				UUID:     transactionpb.GetUuid(),
				Sender:   transactionpb.GetSender(),
				Receiver: transactionpb.GetReceiver(),
				Amount:   transactionpb.GetAmount(),
			})
		}
		myBlockchain = append(myBlockchain, blockchain.Block{
			Index:         data.GetIndex(),
			Proof:         data.GetProof(),
			Previous_hash: data.GetPreviousHash(),
			Timestamp:     data.GetTimestamp(),
			Transactions:  transactions,
		})
	}
}

func addUsers(c blockchainpb.BlockchainClient) {
	stream, err := c.ListUsers(context.Background(), &blockchainpb.ListUserRequest{})
	if err != nil {
		log.Fatalln("Couldn't add user to network", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln("Couldn't add blockchain to client", err)
		}
		userId := res.GetUuid()
		network = append(network, userId)
	}
}

func buy(c blockchainpb.BlockchainClient) {
	fmt.Println("Receiver:")
	scanner.Scan()
	receiver := scanner.Text()
	fmt.Println("Amount:")
	scanner.Scan()
	amountStr := scanner.Text()
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		panic(err)
	}
	buyTransaction(c, receiver, amount)

}

func buyTransaction(c blockchainpb.BlockchainClient, receiver string, amount float64) {
	c.AddTransaction(context.Background(), &blockchainpb.AddTransactionRequest{
		Sender: &blockchainpb.User{
			Uuid: clientUserId,
		},
		Receiver: &blockchainpb.User{
			Uuid: receiver,
		},
		Amount: amount,
	})
}

func shell(balance float64) (int, error) {
	fmt.Print("\033[H\033[2J")
	fmt.Println("Simple Shell")
	fmt.Println("0: exit")
	fmt.Println("1: buy")
	fmt.Println("2: mine block")
	fmt.Printf("Client: %s\tBalanace: %f \n", clientUserId, balance)
	fmt.Println("Command:")
	scanner.Scan()
	cmd := scanner.Text()
	command, err := strconv.Atoi(cmd)
	if err != nil {
		return 0, err
	}

	return command, nil
}
