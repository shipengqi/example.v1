package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"

	pb "github.com/shipengqi/example.v1/packeges/third/grpc/03_simple_tls/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type SearchService struct{}

func (s *SearchService) Search(ctx context.Context, r *pb.SearchRequest) (*pb.SearchResponse, error) {
	log.Println("get request: ", r.GetRequest())
	return &pb.SearchResponse{Response: "Hello, " + r.GetRequest() + " Server"}, nil
}

const PORT = "9001"

func main() {
	cert, err := tls.LoadX509KeyPair("../../ssl/server/server.crt", "../../ssl/server/server.key")
	if err != nil {
		log.Fatalf("tls.LoadX509KeyPair err: %v", err)
	}

	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("../../ssl/ca/ca.crt")
	if err != nil {
		log.Fatalf("ioutil.ReadFile err: %v", err)
	}

	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatalf("certPool.AppendCertsFromPEM err")
	}

	c := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	})
	// 创建 grpc Server 对象
	server := grpc.NewServer(grpc.Creds(c))
	// 注册 SearchService
	pb.RegisterSearchServiceServer(server, &SearchService{})

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
