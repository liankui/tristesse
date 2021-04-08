package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
	pb "github.com/liankui/tristesse/examples/grpc_example/proto"
)

type StreamService struct{}

const (
	PORT = "9002"
)

func main() {
	server := grpc.NewServer()
	pb.RegisterStreamServiceServer(server, &StreamService{})

	lis, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}

	server.Serve(lis)
}

func (s *StreamService) List(r *pb.StreamRequest, stream pb.StreamService_ListServer) error {
	return nil
}

func (s *StreamService) Record(stream pb.StreamService_RecordServer) error {
	return nil
}

func (s *StreamService) Route(stream pb.StreamService_RouteServer) error {
	return nil
}
