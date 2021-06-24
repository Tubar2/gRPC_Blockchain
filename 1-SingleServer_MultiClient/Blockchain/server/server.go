package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/google/uuid"
	"github.com/tubar2/go_Block-Chain/Blockchain/blockchain"
	"github.com/tubar2/go_Block-Chain/Blockchain/blockchainpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var (
	port         = "0.0.0.0:50051"
	myBlockchain *blockchain.Blockchain
	adminUUID    = "SuAdm"
)

type server struct {
	blockchainpb.UnimplementedBlockchainServer
}

func (server *server) CreateBlockchain(ctx context.Context, _ *blockchainpb.CreateBlockRequest) (*blockchainpb.CreateBlockResponse, error) {
	log.Println("Creating New blockchain")
	myBlockchain = blockchain.NewBlockChain()
	myBlockchain.AddNode(adminUUID)
	log.Println("Blockchain created")

	genBlock := myBlockchain.GetPreviousBlock()
	// Returns genesis block information
	return &blockchainpb.CreateBlockResponse{
		Block: genBlock.ToBlockpbfmt(),
	}, nil
}

func (server *server) MineBlock(ctx context.Context, req *blockchainpb.MineBlockRequest) (*blockchainpb.MineBlockResponse, error) {
	log.Println("Mining blockchain")
	if myBlockchain.IsEmptyblockchain() {
		return nil, status.Errorf(
			codes.FailedPrecondition,
			fmt.Sprintln("Blockchain wasn't created"),
		)
	}

	log.Println("Getting previous block and proof")
	previousBlock := myBlockchain.GetPreviousBlock()
	previousProof := previousBlock.Proof

	log.Println("Proofing block and hashing")
	proof := myBlockchain.ProofOfWork(int(previousProof))
	previousHash := previousBlock.Hash()

	log.Println("Getting Miner, and adding transaction")
	miner := req.GetMiner()
	myBlockchain.AddTransaction(adminUUID, miner, 0.1)

	log.Println("Creating block")
	block := myBlockchain.CreateBlock(proof, previousHash)

	log.Println("Success")
	return &blockchainpb.MineBlockResponse{
		Block: block.ToBlockpbfmt(),
	}, nil
}

func (server *server) GetBlockchain(_ *blockchainpb.GetBlocksRequest, stream blockchainpb.Blockchain_GetBlockchainServer) error {
	log.Println("Listing blocks in blockchain")
	for _, block := range myBlockchain.GetChain() {
		if err := stream.Send(&blockchainpb.GetBlocksResponse{
			Block: block.ToBlockpbfmt(),
		}); err != nil {
			return status.Errorf(
				codes.Unknown,
				fmt.Sprintln("Error sending block data", err),
			)
		}
	}
	return nil
}

func (server *server) AddUser(ctx context.Context, _ *blockchainpb.AddUserRequest) (*blockchainpb.AddUserResponse, error) {
	node := uuid.NewString()
	myBlockchain.AddNode(node)
	log.Println("Adding user to blokchain network", node)

	return &blockchainpb.AddUserResponse{
		Uuid: node,
	}, nil
}

func (server *server) ListUsers(req *blockchainpb.ListUserRequest, stream blockchainpb.Blockchain_ListUsersServer) error {
	log.Println("Listing users in blockchain")
	if myBlockchain.IsEmptyblockchain() {
		return status.Errorf(
			codes.FailedPrecondition,
			fmt.Sprintln("Blockchain wasn't created"),
		)
	}
	if myBlockchain.IsEmptyNodes() {
		return status.Errorf(
			codes.FailedPrecondition,
			fmt.Sprintln("User map is empty"),
		)
	}

	for userId, _ := range myBlockchain.GetNodes() {
		if err := stream.Send(&blockchainpb.ListUserResponse{
			Uuid: userId,
		}); err != nil {
			return status.Errorf(
				codes.Internal,
				fmt.Sprintln("Error reading user map", err),
			)
		}
	}

	return nil
}

func (server *server) AddTransaction(ctx context.Context, req *blockchainpb.AddTransactionRequest) (*blockchainpb.AddTransactionResponse, error) {
	sender := req.GetSender().GetUuid()
	receiver := req.GetReceiver().GetUuid()
	amount := req.GetAmount()

	log.Println("Creating transaction. Amount:", amount)
	log.Printf("\tSender: %s\n", sender)
	log.Printf("\tReceiver: %s\n", receiver)

	if !myBlockchain.UserInNodes(sender) {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintln("Sender user ->", sender, "<- not part of chain"),
		)
	}
	if !myBlockchain.UserInNodes(receiver) {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintln("Receiver user ->", receiver, "<- not part of chain"),
		)
	}

	_, userId := myBlockchain.AddTransaction(sender, receiver, float64(amount))

	return &blockchainpb.AddTransactionResponse{
		Transaction: &blockchainpb.Transaction{
			Uuid: userId,
		},
	}, nil
}

func (server *server) GetBallance(ctx context.Context, req *blockchainpb.GetBallanceRequest) (*blockchainpb.GetBallanceResponse, error) {
	userId := req.GetUuid()
	log.Println("Getting balance fo user:", userId)
	if !myBlockchain.UserInNodes(userId) {
		return nil, status.Errorf(
			codes.NotFound, fmt.Sprintln("User not in node"),
		)
	}
	user := myBlockchain.GetUser(userId)

	return &blockchainpb.GetBallanceResponse{
		Ballance: user.GetBallance(),
	}, nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Println("Creating Listener on port:", port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalln("Failed to listen", err)
	}

	log.Println("Registering new server")

	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)
	blockchainpb.RegisterBlockchainServer(s, &server{})

	// Register reflection on server
	reflection.Register(s)

	log.Println("Creating New blockchain")
	myBlockchain = blockchain.NewBlockChain()
	myBlockchain.AddNode(adminUUID)
	log.Println("Blockchain created")

	go func() {
		log.Println("Starting Server")
		if err := s.Serve(lis); err != nil {
			log.Fatalln("Failed to server", err)
		}
	}()

	// Wait for control C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	<-ch
	log.Println("Stopping Server")
	s.Stop()
	lis.Close()
}
