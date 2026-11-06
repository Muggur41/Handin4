package client

import (
	"Handin4/proto"
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)


func SendRequest(address string, nodeID int32, timestamp int64) bool {
	// Set up a connection to the server.
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("did not connect to: %s: %v",address, err)
		return false
	} 
	defer conn.Close()

	// Create a new client for the MutualExclusion service.
	client := proto.NewMutualExclusionClient(conn)

	// Set a timeout for the request.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// Create the request message.
	req := &proto.Request{
		NodeId:    nodeID,
		Timestamp: timestamp,
	}

	// Send the request to the server and receive a reply.
	_, err = client.RequestAccess(ctx, req)
	if err != nil {
		log.Printf("could not request access: %v", err)
		return false
	}

	// Log the successful receipt of the reply.
	log.Printf("received reply from node %s", address)
	return true
}