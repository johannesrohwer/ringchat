package main

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
)

func dial(host string, port int) *grpc.ClientConn {
	// Establish connection to server
	target := fmt.Sprintf("%s:%d", host, port)
	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %s", err)
	}

	return conn
}
