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


