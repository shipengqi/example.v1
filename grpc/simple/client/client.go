package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/shipengqi/example.v1/grpc/simple/proto"
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
	// 创建 SearchService 的 Client
	client := pb.NewSearchServiceClient(conn)
	// 发送 RPC 请求
	resp, err := client.Search(context.Background(), &pb.SearchRequest{Request: "gRPC"})
	if err != nil {
		log.Fatalf("client.Search err: %v", err)
	}
	log.Printf("resp: %s", resp.GetResponse())
}
