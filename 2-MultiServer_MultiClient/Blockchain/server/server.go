// Defines the bahavior of a client's call to a server
// holding a blockchain
package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/google/uuid"
	"github.com/tubar2/go_MultiServer_Blockchain/Blockchain/blockchain"
	"github.com/tubar2/go_MultiServer_Blockchain/Blockchain/blockchainpb"
	"github.com/tubar2/go_MultiServer_Blockchain/Blockchain/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
// port string
// myBlockchain *blockchain.Blockchain
// adminUUID string
)

//
// MARK: Server Struct
//
type Server struct {
	blockchainpb.UnimplementedBlockchainServiceServer

	myBlockchain blockchain.Blockchain
	thisUser     blockchain.User

	transactions blockchain.Transactions
	nodes        map[string]blockchain.User // maps [uuid] -> User{.node}
}

//
// note: GRPC server services
//

func (server *Server) GetBlockchain(ctx context.Context, _ *blockchainpb.GetBlockchainRequest) (*blockchainpb.GetBlockchainResponse, error) {
	log.Println("Getting blockchain")
	return &blockchainpb.GetBlockchainResponse{
		Blockchain: server.myBlockchain.ToBlockchainpbfmt(),
	}, nil
}

func (server *Server) UpdateUserBalance(ctx context.Context, req *blockchainpb.UpdateUserBalanceRequest) (*blockchainpb.UpdateUserBalanceResponse, error) {
	log.Println("Updating user balance")
	amount := req.GetAmount()
	server.updateUserBalance(amount)

	return &blockchainpb.UpdateUserBalanceResponse{}, nil
}

func (server *Server) MineBlock(ctx context.Context, req *blockchainpb.MineBlockRequest) (*blockchainpb.MineBlockResponse, error) {
	log.Println("Mining block")
	block := server.MineServerBlock()

	// Propagate to network for the new chain.
	// If new chain is larger than 50% + 1 than other chains,
	// it will be replaced on the network
	chainSuccesfullyPropagated := true
	for _, user := range server.nodes {
		if user.Node == server.thisUser.Node {
			continue
		}
		log.Printf("Trying to replace chain @%s\n", user.Node)
		newChain := client.TryToReplaceChainAt(user.Node, &server.myBlockchain)
		if newChain != nil {
			log.Println("Chain on connection is larger. Replacing chain @host and repeating replacement in network")
			server.myBlockchain = *newChain
			chainSuccesfullyPropagated = false
			break
		}
	}

	// TODO: Use goroutines to speed this up
	// If the mined block was succesfully propafated to network, clear stored transaction
	if chainSuccesfullyPropagated {
		// Update balances and clear propagated transactions stored at other nodes in network
		for _, transaction := range server.getTransactions() {
			sender := transaction.Sender
			receiver := transaction.Receiver
			amount := transaction.Amount
			client.UpdateNodeBalanceBy(server.nodeAt(sender), (-1)*amount)
			client.UpdateNodeBalanceBy(server.nodeAt(receiver), amount)
			for _, user := range server.nodes {
				client.RemoveTransactionAt(transaction.UUID, user.Node)
			}
		}
		server.transactions = blockchain.NewTransactions()
	}

	return &blockchainpb.MineBlockResponse{
		Block: block.ToBlockpbfmt(),
	}, nil
}

func (server *Server) RemoveTransaction(ctx context.Context, req *blockchainpb.RemoveTransactionRequest) (*blockchainpb.RemoveTransactionResponse, error) {
	log.Println("Removing transaction")
	tId := req.GetTransactionId()
	removed := server.removeTransaction(tId)

	return &blockchainpb.RemoveTransactionResponse{
		Removed: removed,
	}, nil
}

func (server *Server) ReplaceChain(ctx context.Context, req *blockchainpb.ReplaceChainRequest) (*blockchainpb.ReplaceChainResponse, error) {
	log.Println("Replacing blockchain")
	tempBlockchain := req.GetBlockchain()
	size := len(tempBlockchain.GetChain())

	if size > server.myBlockchain.Len() {
		log.Println("New chain is bigger, replacing")
		newChain := blockchain.FromBlockchainpbfmt(tempBlockchain)
		server.myBlockchain = newChain
		return &blockchainpb.ReplaceChainResponse{
			Replaced: true,
		}, nil
	} else if size < server.myBlockchain.Len() {
		log.Println("Current node's chain is bigger, replacing on request")
		return &blockchainpb.ReplaceChainResponse{
			Replaced:   false,
			Blockchain: server.myBlockchain.ToBlockchainpbfmt(),
		}, nil
	} else {
		// TODO: What happens if chains are the same size but different
		log.Println("Nodes have the same size")
		return &blockchainpb.ReplaceChainResponse{
			Replaced: false,
		}, nil
	}
}

// Returns current node's uuid
func (server *Server) GetUserId(ctx context.Context, req *blockchainpb.GetUserIDRequest) (*blockchainpb.GetUserIDResponse, error) {
	log.Println("Getting userId")
	return &blockchainpb.GetUserIDResponse{
		Uuid: server.thisUser.Uuid,
	}, nil
}

func (server *Server) ConnectToNetwork(req *blockchainpb.ConnectToNetworkRequest, stream blockchainpb.BlockchainService_ConnectToNetworkServer) error {
	log.Printf("ConnectToNetwork call from @%s", server.thisUser.Node)

	// Connects to provided node's network
	_node := req.GetNode()
	log.Printf("Connecting node to network @%s", _node)
	if _node == server.thisUser.Node {
		if err := stream.Send(&blockchainpb.ConnectToNetworkResponse{
			Result: "Can't connect node to itself",
		}); err != nil {
			return err
		}
		return nil
	}
	connections := client.ConnectUserToNet(server.thisUser, _node)

	log.Println("Appending new connections")
	for _, connection := range connections {
		server.addUser(connection.Node, connection.Uuid)
		if err := stream.Send(&blockchainpb.ConnectToNetworkResponse{
			Result: fmt.Sprintf("connection to @%s : true", connection.GetNode()),
		}); err != nil {
			return err
		}
	}
	return nil
}
func (server *Server) ConnectUser(req *blockchainpb.ConnectUserRequest, stream blockchainpb.BlockchainService_ConnectUserServer) error {
	log.Println("Connecting user to current node's network")

	// {"50052"}
	for _, _user := range server.nodes {
		log.Println("\t>\033[32m", _user.Node, "\033[0m")
		if _user.Node == req.GetUser().GetNode() {
			continue
		}
		// _user = "50052"
		connected := client.ConnectUserToNode(req.GetUser(), _user.Node)
		if err := stream.Send(&blockchainpb.ConnectUserResponse{
			Result: connected,
			User:   _user.ToUserpbfmt(),
		}); err != nil {
			return err
		}
	}
	return nil
}
func (server *Server) ConnectToNode(ctx context.Context, req *blockchainpb.ConnectToNodeRequest) (*blockchainpb.ConnectToNodeResponse, error) {
	log.Printf("Adding node @%s to network", req.GetUser().GetNode())
	return &blockchainpb.ConnectToNodeResponse{
		Result: !server.addUser(req.User.Node, req.User.Uuid),
	}, nil
}

// TODO: Use server.addTransaction return to see if transaction already existed
func (server *Server) AddTransaction(ctx context.Context, req *blockchainpb.AddTransactionRequest) (*blockchainpb.AddTransactionResponse, error) {
	log.Println("Adding transaction to network")
	amount := req.GetAmount()
	sender := req.GetSender()
	receiver := req.GetReceiver()

	tId := uuid.NewString()
	transaction := blockchain.Transaction{
		UUID:     tId,
		Sender:   sender,
		Receiver: receiver,
		Amount:   amount,
	}

	// Propagate transaction to other nodes
	for _, user := range server.nodes {
		if user.Node == server.thisUser.Node {
			continue
		}
		client.PropagateTransactionTo(user.Node, transaction)
	}

	server.addTransaction(transaction)
	return &blockchainpb.AddTransactionResponse{
		TransactionId: tId,
	}, nil
}

func (server *Server) AppendTransactions(ctx context.Context, req *blockchainpb.AppendTransactionsRequest) (*blockchainpb.AppendTransactionsResponse, error) {
	log.Println("Appending new transactions to this node")
	mTransaction := blockchain.FromTransactionpbfmt(req.GetTransaction())
	server.addTransaction(mTransaction)
	return &blockchainpb.AppendTransactionsResponse{
		Result: fmt.Sprintln("Succesfully added transaction", mTransaction.UUID, "to node", server.thisUser.Node),
	}, nil
}

func (server *Server) GetTransactions(_ *blockchainpb.GetTransactionsRequest, stream blockchainpb.BlockchainService_GetTransactionsServer) error {
	log.Println("Getting transactions")

	for _, transaction := range server.getTransactions() {
		if err := stream.Send(&blockchainpb.GetTransactionsResponse{
			Transaction: transaction.ToTransactionpbfmt(),
		}); err != nil {
			return err
		}
	}
	return nil
}

func (server *Server) GetUsers(req *blockchainpb.GetUsersRequest, stream blockchainpb.BlockchainService_GetUsersServer) error {
	log.Println("Getting users")

	for _, user := range server.nodes {
		if err := stream.Send(&blockchainpb.User{
			Uuid: user.Uuid,
			Node: user.Node,
		}); err != nil {
			return err
		}
	}
	return nil
}

//
// note: Server Functions
//
func Start(port string) (net.Listener, *grpc.Server) {
	log.Println("Creating Listener on port:", port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalln("Failed to listen", err)
	}

	log.Println("Creating new user for current node")
	id := uuid.NewString()
	newUser := blockchain.User{
		Uuid:     id,
		Ballance: 0,
		Node:     port,
	}
	log.Printf("Created user: @%s {%s}\n", port, id)

	log.Println("Creating and registering new server grpc service")
	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)
	blockchainpb.RegisterBlockchainServiceServer(s, &Server{
		thisUser:     newUser,
		myBlockchain: blockchain.NewBlockChain(),
		transactions: blockchain.NewTransactions(),
		nodes: map[string]blockchain.User{
			id: newUser,
		},
	})

	log.Println("Registering server reflection service")
	reflection.Register(s)

	return lis, s
}

//
// note: Server Methods
//

// Mines the a block for blockchain in server
func (server *Server) MineServerBlock() blockchain.Block {
	previousBlock := server.myBlockchain.GetLastBlock()
	previousProof := previousBlock.Proof
	proof := server.myBlockchain.ProofOfWork(int(previousProof))
	previousHash := previousBlock.Hash()

	// TODO: Propagate information for balance change
	return server.myBlockchain.CreateBlock(proof, previousHash, server.transactions)
}

func (server *Server) addUser(node, uuid string) bool {
	_, ok := server.nodes[uuid]
	if !ok {
		server.nodes[uuid] = blockchain.User{
			Uuid: uuid,
			Node: node,
		}
	}
	return ok
}

func (server *Server) nodeAt(uuid string) string {
	return server.nodes[uuid].Node
}

func (server *Server) removeTransaction(tId string) bool {
	return server.transactions.RemoveTransaction(tId)
}

func (server *Server) getTransactions() map[string]blockchain.Transaction {
	return server.transactions.GetTransactions()
}

func (server *Server) updateUserBalance(amount float64) {
	server.thisUser.UpdateUserBalance(amount)
}

func (server *Server) addTransaction(transaction blockchain.Transaction) bool {
	return server.transactions.AddTransaction(transaction)
}
