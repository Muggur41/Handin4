package server

import (
	"Handin4/node" //node logic
	"Handin4/proto" //generated proto files
	"context" //required for grpc
	"log" //message logging
)

type Server struct {
	proto.UnimplementedMutualExclusionServer //defautt implementation of the server
	Node *node.Node	//pointer to node
}

func (s *Server) RequestAccess(ctx context.Context, req *proto.Request) (*proto.Reply, error) {
	log.Printf("Received request from node %d with timestamp %d", req.NodeId, req.Timestamp)
	
	s.Node.UpdateTimestamp(req.Timestamp) 

	if s.Node.ShouldDefer(req.NodeId, req.Timestamp) {
		log.Printf("Deffering request from node %d", req.NodeId)
		ch := s.Node.DeferReply(req.NodeId) //will reply later when I release
		<-ch //wait for the signal to reply
		log.Printf("now replying to deffered node %d", req.NodeId)
	} else {
		log.Printf("Sending reply to node %d", req.NodeId)
	}
	return &proto.Reply{NodeId: s.Node.ID}, nil
}

