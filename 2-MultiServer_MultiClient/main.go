package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/tubar2/go_MultiServer_Blockchain/Blockchain/server"
)

var (
	host = "localhost"
	port = 50051
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	hostFlag := flag.String("host", host, fmt.Sprintf("Host to listen on. Defaults to %s", host))
	portFlag := flag.Int("port", port, fmt.Sprintf("Port to listen on. Defaults to %d", port))
	flag.Parse()

	lis, s := server.Start(fmt.Sprintf("%s:%d", *hostFlag, *portFlag))

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
