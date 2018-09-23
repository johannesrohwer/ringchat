package main

import (
	"context"
	rc "github.com/johannesrohwer/ringchat/grpc/ringchat"
	"log"
	"sync"
)

// Implement RingMaster Interface
type RingMasterServer struct {
	sync.Mutex
	ring []*rc.Node
}

func (rms *RingMasterServer) AddToRing(ctx context.Context, node *rc.Node) (*rc.Node, error) {
	rms.Lock()
	defer rms.Unlock()

	// TODO: Add support for insertion and removal

	if len(rms.ring) == 0 {
		// Handle ring of length 0
		// No previous node, no need to update anything
	} else {
		// Handle ring of length > 0
		// Update previous tail node
		tail := rms.ring[len(rms.ring)-1]
		if err := setNext(tail, node); err != nil {
			log.Fatalf("could not add to ring %v", err)
		}

	}

	// Add new node to the ring
	rms.ring = append(rms.ring, node)
	return rms.ring[0], nil
}

func setNext(node *rc.Node, next *rc.Node) error {
	conn := dial(node.Name, int(node.Port))
	defer conn.Close()

	client := rc.NewRingSlaveClient(conn)
	if _, err := client.SetNext(context.Background(), next); err != nil {
		log.Fatalf("error when calling setNext(): %s", err)
	}

	return nil
}
