package main

import (
	"context"
	"github.com/google/uuid"
	rc "github.com/johannesrohwer/ringchat/grpc/ringchat"
	"log"
	"sync"
)

// Implement RingSlave Interface
type RingSlaveServer struct {
	ID   string
	self *rc.Node
	next *rc.Node
	sync.Mutex
}

func NewRingSlaveServer(selfHost string, selfPort int) *RingSlaveServer {
	id := uuid.New().String()
	self := &rc.Node{Name: selfHost, Port: int32(selfPort)}
	rss := RingSlaveServer{self: self, ID: id}
	return &rss
}

func (rss *RingSlaveServer) Broadcast(message string) {
	token := &rc.Token{Message: message, Id: rss.ID, Payload: "reserved"}
	rss.ForwardWithoutCheck(context.Background(), token)
}

func (rss *RingSlaveServer) Forward(ctx context.Context, token *rc.Token) (*rc.Empty, error) {
	// If message was sent by this node, ignore it
	if token.Id == rss.ID {
		log.Printf("[%s] destroy token with message %s\n", token.Id, token.Message)
		return &rc.Empty{}, nil
	}

	return rss.ForwardWithoutCheck(ctx, token)
}

func (rss *RingSlaveServer) ForwardWithoutCheck(ctx context.Context, token *rc.Token) (*rc.Empty, error) {
	log.Printf("[%s]\t%s", token.Id, token.Message)

	conn := dial(rss.next.Name, int(rss.next.Port))
	defer conn.Close()
	client := rc.NewRingSlaveClient(conn)
	_, err := client.Forward(context.Background(), token)
	if err != nil {
		log.Fatalf("could not forward: %v", err)
	}

	return &rc.Empty{}, nil
}

func (rss *RingSlaveServer) SetNext(ctx context.Context, next *rc.Node) (*rc.Empty, error) {
	rss.Lock()
	defer rss.Unlock()
	rss.next = next
	log.Printf("next for %v was set to %v", rss.self, next)
	return &rc.Empty{}, nil
}

func (rss *RingSlaveServer) JoinRing(host string, port int) {
	conn := dial(host, port)
	defer conn.Close()
	client := rc.NewRingMasterClient(conn)
	res, err := client.AddToRing(context.Background(), rss.self)
	if err != nil {
		log.Fatalf("could not add self to ring: %v", err)
	}

	rss.next = res
	log.Printf("sucessfully joined ring. next: %v", rss.next)
}
