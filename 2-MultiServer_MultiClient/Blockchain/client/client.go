//
// Defines the capabilities of a user in a blockchain
// -> Mine a block
// -> Create a transaction
//
package client

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/tubar2/go_MultiServer_Blockchain/Blockchain/blockchain"
	"github.com/tubar2/go_MultiServer_Blockchain/Blockchain/blockchainpb"
	"google.golang.org/grpc"
)

// Returns the uuuid from a given node. (eg. "0.0.0.0:50051")
func GetUserIdFrom(node string) (string, error) {
	log.Println("Getting uuid from node:", node)
	opts := grpc.WithInsecure()
	cc, err := grpc.Dial(node, opts)
	if err != nil {
		log.Fatalln("Client couldn't connect", err)
	}
	defer cc.Close()

	c := blockchainpb.NewBlockchainServiceClient(cc)
	res, err := c.GetUserId(context.Background(), &blockchainpb.GetUserIDRequest{})
	if err != nil {
		return "", err
	}
	return res.GetUuid(), nil
}

// Updates the balance of a given node
func UpdateNodeBalanceBy(node string, amount float64) error {
	log.Println("Updating balance for node:", node)
	opts := grpc.WithInsecure()
	cc, err := grpc.Dial(node, opts)
	if err != nil {
		log.Fatalln("Client couldn't connect", err)
	}
	defer cc.Close()

	c := blockchainpb.NewBlockchainServiceClient(cc)
	_, err = c.UpdateUserBalance(context.Background(), &blockchainpb.UpdateUserBalanceRequest{
		Amount: amount,
	})

	return err
}

// TODO: Implement
func RemoveTransactionAt(transactionId, node string) {
	log.Printf("Removing transaction %s at node:%s\n", transactionId, node)
	opts := grpc.WithInsecure()
	cc, err := grpc.Dial(node, opts)
	if err != nil {
		log.Fatalln("Client couldn't connect", err)
	}
	defer cc.Close()

	c := blockchainpb.NewBlockchainServiceClient(cc)

	res, err := c.RemoveTransaction(context.Background(), &blockchainpb.RemoveTransactionRequest{
		TransactionId: transactionId,
	})
	if err != nil {
		log.Fatalln("Error removing transaction")
	}
	if !res.GetRemoved() {
		log.Println("Didn't remove transaction. Perhaps it doesn't exist in node", node)
	}
}

// Returns blokchain from given node
func GetBlockchainFrom(node string) *blockchainpb.Blockchain {
	log.Println("Getting blockchain at node:", node)
	opts := grpc.WithInsecure()
	cc, err := grpc.Dial(node, opts)
	if err != nil {
		log.Fatalln("Client couldn't connect", err)
	}
	defer cc.Close()

	c := blockchainpb.NewBlockchainServiceClient(cc)

	res, err := c.GetBlockchain(context.Background(), &blockchainpb.GetBlockchainRequest{})
	if err != nil {
		log.Fatalln("Error getting blockchain", err)
	}
	return res.GetBlockchain()
}

func PropagateTransactionTo(node string, transaction blockchain.Transaction) {
	log.Println("Propagating transaction to node:", node)
	opts := grpc.WithInsecure()
	cc, err := grpc.Dial(node, opts)
	if err != nil {
		log.Fatalln("Client couldn't connect", err)
	}
	defer cc.Close()

	c := blockchainpb.NewBlockchainServiceClient(cc)

	res, err := c.AppendTransactions(context.Background(), &blockchainpb.AppendTransactionsRequest{
		Transaction: transaction.ToTransactionpbfmt(),
	})
	if err != nil {
		log.Fatalln("Error appending transaction", err)
	}
	log.Println(res.GetResult())
}

func ConnectUserToNet(newUser blockchain.User, otherNode string) []*blockchainpb.User {
	log.Printf("Joining network from node @%s", otherNode)
	opts := grpc.WithInsecure()
	cc, err := grpc.Dial(otherNode, opts)
	if err != nil {
		log.Fatalln("Client couldn't connect", err)
	}
	defer cc.Close()
	c := blockchainpb.NewBlockchainServiceClient(cc)

	// Request connection from newNode to otherNode's netowrk
	stream, err := c.ConnectUser(context.Background(), &blockchainpb.ConnectUserRequest{
		User: newUser.ToUserpbfmt(),
	})
	if err != nil {
		log.Fatalln("Error appending transaction", err)
	}
	log.Println("Receiving new conections")
	var connections []*blockchainpb.User
	log.Println("Adding to connections nodes:")
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			log.Println("Received all connections")
			break
		}
		if err != nil {
			log.Fatalln("Error receiving", err)
		}
		if res.GetResult() {
			fmt.Println("\t\t\t>\033[32m", res.GetUser().Node, "\033[0m")
			connections = append(connections, res.GetUser())
		}
	}
	return connections
}

func ConnectUserToNode(newUser *blockchainpb.User, node string) bool {
	log.Printf("Connecting user @%s to node @%s", newUser.Node, node)
	opts := grpc.WithInsecure()
	cc, err := grpc.Dial(node, opts)
	if err != nil {
		log.Fatalln("Client couldn't connect", err)
	}
	defer cc.Close()
	c := blockchainpb.NewBlockchainServiceClient(cc)

	res, err := c.ConnectToNode(context.Background(), &blockchainpb.ConnectToNodeRequest{
		User: newUser,
	})
	if err != nil {
		log.Println("Error:", err)
		return false
	}
	return res.GetResult()
}

func TryToReplaceChainAt(node string, _blockchain *blockchain.Blockchain) *blockchain.Blockchain {
	log.Println("Propagating blockchain to node:", node)
	opts := grpc.WithInsecure()
	cc, err := grpc.Dial(node, opts)
	if err != nil {
		log.Fatalln("Client couldn't connect", err)
	}
	defer cc.Close()

	c := blockchainpb.NewBlockchainServiceClient(cc)

	res, err := c.ReplaceChain(context.Background(), &blockchainpb.ReplaceChainRequest{
		Blockchain: _blockchain.ToBlockchainpbfmt(),
	})
	if err != nil {
		log.Fatalln("Error appending transaction", err)
	}
	if res.GetReplaced() {
		log.Println("Replaced chain at node", node)
	} else {
		log.Println("Chain not replaced")
		if res.GetBlockchain() != nil {
			log.Printf("Replacing this node's chain with connection's @%s\n", node)
			newChain := blockchain.FromBlockchainpbfmt(res.GetBlockchain())
			return &newChain
		} else {
			log.Println("Chains of same size")
		}
	}
	return nil
}
