# 1-SingleServer_MultiClient

This is a implementation of blockchain using a single server (localized blockchain) that can be accessed by multiple clients using gRPC buffers.
If you want to see a implementation where each client is also a server, check 2-MultiServer.

## To run:
From current directory run `go run Blockchain/server/server.go` to start the server containing the blockchain

Next, from a different terminal window, connect a client to it by running `go run Blockchain/client/client.go`.
You'll become a new user connected to the network that is able to make a transaction and ask for the mine to be blocked
