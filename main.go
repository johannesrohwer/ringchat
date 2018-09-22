package main

import "fmt"
import (
	"context"
	rc "github.com/johannesrohwer/ringchat/grpc/ringchat"
	"google.golang.org/grpc"
	"log"
	"net"
)

func printHelp() {
	fmt.Println("Help")
	fmt.Println("!q \t quit")
	fmt.Println("!c \t connect")
	fmt.Println("!p \t ping")
}

func startServer(port int) {
	// Start listening
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("listening on port %d\n", port)

	// Create a gRPC server object
	grpcServer := grpc.NewServer()
	rs := RingchatServer{}
	rc.RegisterRingchatServer(grpcServer, &rs)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

func connect(host string, port int) (*rc.RingchatClient, *grpc.ClientConn) {
	// FIXME: Conn needs to be close manually

	// Establish connection to server
	target := fmt.Sprintf("%s:%d", host, port)
	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %s", err)
	}

	client := rc.NewRingchatClient(conn)
	return &client, conn
}

func ping(client *rc.RingchatClient) {
	response, err := (*client).Ping(context.Background(), &rc.EchoRequest{Message: "hi! :)"})
	if err != nil {
		log.Fatalf("Error when calling Ping: %s", err)
	}
	log.Printf("Response from server: %s", response.Message)

}

func main() {
	// Load parameters & constants
	port := 7777

	go startServer(port)
	printHelp()

	var input = ""
	var client *rc.RingchatClient
	var conn *grpc.ClientConn
inputLoop:
	for {
		switch input {
		case "!q":
			break inputLoop
		case "!c":
			client, conn = connect("", port)
		case "!p":
			ping(client)
		}

		fmt.Print("$ ")
		fmt.Scanln(&input)
	}

	conn.Close()
}
