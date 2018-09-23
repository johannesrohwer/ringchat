package main

import "fmt"
import (
	"bufio"
	"flag"
	rc "github.com/johannesrohwer/ringchat/grpc/ringchat"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"time"
)

func startRingMaster(masterPort int) {
	// Start listening
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", masterPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("ring master listening on port %d\n", masterPort)

	// Create a gRPC server object
	grpcServer := grpc.NewServer()
	rms := RingMasterServer{}
	rc.RegisterRingMasterServer(grpcServer, &rms)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("ring master failed to serve: %s", err)
	}
}

func startRingSlave(masterHost string, masterPort int, returnCh chan *RingSlaveServer) {
	// Start listening
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 0))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	selfHost := "" // FIXME: add current IP / hostname
	selfPort := listener.Addr().(*net.TCPAddr).Port
	rss := NewRingSlaveServer(selfHost, selfPort)
	rc.RegisterRingSlaveServer(grpcServer, rss)

	// FIXME: Add error handling
	go grpcServer.Serve(listener)

	rss.JoinRing(masterHost, masterPort)
	returnCh <- rss
}

func main() {
	// Load parameters & constants
	// FIXME: adjust description of parameters
	masterBoolPtr := flag.Bool("master", false, "set to true if this is a master node")
	masterPortIntPtr := flag.Int("master-port", 9999, "port of the master node")
	masterHostStrPtr := flag.String("master-host", "", "host of the master node")
	flag.Parse()

	isMaster := *masterBoolPtr
	masterPort := *masterPortIntPtr
	masterHost := *masterHostStrPtr

	// Start master server if required
	if isMaster {
		go startRingMaster(masterPort)
	}

	// FIXME: Do not just wait until its up but actually check
	log.Printf("sleep 3 seconds...\n")
	time.Sleep(3 * time.Second)

	// Start ring slave
	slaveCh := make(chan *RingSlaveServer, 1)
	go startRingSlave(masterHost, masterPort, slaveCh)
	rss := <-slaveCh

	// Read and evaluate input
	// FIXME: Refactor this into own method
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		if input == "!q" {
			break
		} else {
			rss.Broadcast(input)
		}

	}

}
