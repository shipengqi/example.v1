package client

import (
	"fmt"
	"log"

	pb "github.com/shipengqi/example.v1/grpc/stream/proto"
	"google.golang.org/grpc"
)

const PORT = "9001"

func main() {
	// 创建与 server 的连接
	conn, err := grpc.Dial(fmt.Sprintf(":%s", PORT), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
	}
	defer conn.Close()
	// 创建 StreamService 的 Client
	client := pb.NewStreamServiceClient(conn)

	err = printLists(client, &pb.StreamRequest{Pt: &pb.StreamPoint{Name: "gRPC Stream Client: List", Value: 2018}})
	if err != nil {
		log.Fatalf("printLists.err: %v", err)
	}

	err = printRecord(client, &pb.StreamRequest{Pt: &pb.StreamPoint{Name: "gRPC Stream Client: Record", Value: 2018}})
	if err != nil {
		log.Fatalf("printRecord.err: %v", err)
	}

	err = printRoute(client, &pb.StreamRequest{Pt: &pb.StreamPoint{Name: "gRPC Stream Client: Route", Value: 2018}})
	if err != nil {
		log.Fatalf("printRoute.err: %v", err)
	}
}

func printLists(client pb.StreamServiceClient, r *pb.StreamRequest) error {
	return nil
}

func printRecord(client pb.StreamServiceClient, r *pb.StreamRequest) error {
	return nil
}

func printRoute(client pb.StreamServiceClient, r *pb.StreamRequest) error {
	return nil
}
