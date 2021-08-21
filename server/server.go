package server

import (
	"google.golang.org/grpc"
	"net"
)

var s *grpc.Server

func Init() *grpc.Server {
	s = grpc.NewServer()
	return s
}

func Run(address string) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	return s.Serve(lis)
}
