package main

import (
	"context"
	rc "github.com/johannesrohwer/ringchat/grpc/ringchat"
	"log"
)

// RingchatServer implements required methods of the proto interface
type RingchatServer struct{}

func (s *RingchatServer) Ping(ctx context.Context, request *rc.EchoRequest) (*rc.EchoResponse, error) {
	log.Println("received ping")
	return &rc.EchoResponse{Message: "pong"}, nil
}
