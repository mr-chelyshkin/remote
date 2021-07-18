package server

import (
	"github.com/mr-chelyshkin/remote/pb"
)

type TransportServer struct {
	pb.UnimplementedTransportServer
}

// NewTransportServer return TransportServer object as pointer.
func NewTransportServer() *TransportServer {
	return &TransportServer{}
}