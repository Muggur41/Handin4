package main

  import (
  	"Handin4/client"
  	"Handin4/node"
  	"Handin4/proto"
  	"Handin4/server"
  	"flag"
  	"log"
  	"net"
  	"strings"
  	"time"
	"fmt"
	"sync"

  	"google.golang.org/grpc"
  )

func main() {
	id := flag.Int("id", 1, "Node ID") //Which node are we
	port := flag.Int("port", 5001, "Port to listen on") //what port to listen on
	peers := flag.String("peers", "", "Comma-separated list of peer addresses e.g., localhost:5002,localhost:5003") //addresses of other nodes
	
	flag.Parse()

	// Create a new node with the given ID and number of peers
	var peerList []string
	if *peers != "" {
		peerList = strings.Split(*peers, ",")
	}

	log.Printf("Starting node %d on port %d with peers: %v", *id, *port, peerList)

	n := node.NewNode(int32(*id), len(peerList))

  	go startServer(n, *port)
  	
	time.Sleep(2 * time.Second)

	// Request critical section
  	log.Printf("Node %d requesting critical section", *id)
  	timestamp := n.RequestCS()

  	// Send request to all peers
  	  	// Send request to all peers and wait for replies
  	var wg sync.WaitGroup
  	for _, peer := range peerList {
  		wg.Add(1)
  		go func(addr string) {
  			defer wg.Done()
  			client.SendRequest(addr, int32(*id), timestamp)
  			n.ReceiveReply()
  		}(peer)
  	}
  	wg.Wait()

  	// Now we have all replies, enter critical section
  	log.Printf("Node %d received all replies, entering critical section", *id)
  	n.EnterCS()

  	// Critical section, just print something
  	log.Printf(">>> Node %d is IN the critical section <<<", *id)
  	time.Sleep(3 * time.Second)
  	log.Printf(">>> Node %d is LEAVING the critical section <<<", *id)

  	n.ReleaseCS()

  	// Keep the program running
  	select {}

  	// Keep the program running
  	select {}
}

func startServer(n *node.Node, port int) {
	// Set up a TCP listener on the specified port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create a new gRPC server and register the MutualExclusion service
	s := grpc.NewServer()
	proto.RegisterMutualExclusionServer(s, &server.Server{Node: n})

	// Start serving incoming connections
	log.Printf("Server listening on port %d", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

