package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	pb "github.com/shipengqi/example.v1/grpc/01_simple/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

const PORT = "9001"

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

	// 创建与 server 的连接
	conn, err := grpc.Dial(fmt.Sprintf(":%s", PORT), grpc.WithTransportCredentials(c))
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
	}
	defer conn.Close()

	// context.WithDeadline：会返回最终上下文截止时间。第一个形参为父上下文，第二个形参为调整的截止时间。
	// 若父级时间早于子级时间，则以父级时间为准，否则以子级时间为最终截止时间
	// ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Duration(5 * time.Second)))
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	// 创建 SearchService 的 Client
	client := pb.NewSearchServiceClient(conn)
	// 发送 RPC 请求
	resp, err := client.Search(ctx, &pb.SearchRequest{Request: "gRPC"})
	if err != nil {
		// FromError：返回 GRPCStatus 的具体错误码，若为非法，则直接返回 codes.Unknown
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				log.Fatalln("client.Search err: deadline")
			}
		}

		log.Fatalf("client.Search err: %v", err)
	}
	log.Printf("resp: %s", resp.GetResponse())
}

// $ go run client.go
// 2021/02/25 16:38:48 client.Search err: deadline
