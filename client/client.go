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

	client := proto.NewMutualExclusionClient(conn)

  	req := &proto.Request{
  		NodeId:    nodeID,
  		Timestamp: timestamp,
  	}

  	// Retry the RPC call up to 5 times
  	for i := 0; i < 5; i++ {
  		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  		_, err = client.RequestAccess(ctx, req)
  		cancel()

  		if err == nil {
  			log.Printf("Received reply from %s", address)
  			return true
  		}
  		log.Printf("Attempt %d: Request to %s failed: %v", i+1, address, err)
  		time.Sleep(1 * time.Second)
  	}

  	log.Printf("Failed to get reply from %s after 5 attempts", address)
  	return false
  }
