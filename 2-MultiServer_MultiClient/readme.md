# 2-MultiServer_MultiClient

## Before:
To test this program's functionality I recommend using the evans cli (https://github.com/ktr0731/evans)

## Running:
From the programs root directory run `go run main.go` \
This will start the program to run a server on default port localhost:50051 \
To access this server with evans use `evans -p 50051 -r`

## Running on different ports:
From the programs root directory run `go run main.go -port 50052` \
This will start the program to run a server on provided port localhost:50052 \
To access this server with evans use `evans -p 50052 -r`

## Running on different host:
From the programs root directory run `go run main.go -host 0.0.0.0` \
This will start the program to run a server on provided port 0.0.0.0:50051

## First Steps
- I recommend first, from evans, using `call ConnectToNetwork` and entering a port running another instance of the server such as `localhost:50052`
- Next, use functions `call AddTransaction` and `call MineBlock` to send blocks to each node on the networks
