package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"

	pb "github.com/shipengqi/example.v1/grpc/01_simple/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const PORT = "9001"

// Auth 实现了 PerRPCCredentials 接口
type Auth struct {
	AppKey    string
	AppSecret string
}

func (a *Auth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{"app_key": a.AppKey, "app_secret": a.AppSecret}, nil
}

func (a *Auth) RequireTransportSecurity() bool {
	return true
}

func main() {
	cert, err := tls.LoadX509KeyPair("../../ssl/client/client.crt", "../../ssl/client/client.key")
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
		ServerName:   "grpc-example",
		RootCAs:      certPool,
	})

	auth := Auth{
		AppKey:    "example",
		AppSecret: "123456",
	}

	// 创建与 server 的连接
	// 用到了 grpc.WithPerRPCCredentials
	conn, err := grpc.Dial(fmt.Sprintf(":%s", PORT), grpc.WithTransportCredentials(c), grpc.WithPerRPCCredentials(&auth))
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

// AppSecret: "1234565",
// $ go run client.go
// Output:
// 2021/02/25 16:03:15 client.Search err: rpc error: code = Unauthenticated desc = 自定义认证 Token 无效

// AppSecret: "123456",
// $ go run client.go
// Output:
// 2021/02/25 16:04:46 resp: Hello, gRPC Server
