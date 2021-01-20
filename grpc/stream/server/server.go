package main

import (
	"fmt"
	"log"
	"net"

	pb "github.com/shipengqi/example.v1/grpc/stream/proto"
	"google.golang.org/grpc"
)

type StreamService struct {

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

const (
	PORT = "9002"
)

func main()  {
	server := grpc.NewServer()
	pb.RegisterStreamServiceServer(server, &StreamService{})

	// 监听 TCP 端口
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", PORT))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// grpc Server 开始 lis.Accept，直到 Stop 或 GracefulStop
	if err := server.Serve(lis); err != nil {
		log.Fatalf("gRPC.Serve: %v", err)
	}
}

