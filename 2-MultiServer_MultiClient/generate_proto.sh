#!/bin/bash

protoc blockchain/blockchainpb/blockchain.proto --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative