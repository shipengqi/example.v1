package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/shipengqi/example.v1/grpc/03_simple_tls/proto"
)

type Auth struct {
	appKey    string
	appSecret string
}

func (a *Auth) Check(ctx context.Context) error {
	// 调用 metadata.FromIncomingContext 从上下文中获取 metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "自定义认证 Token 失败")
	}

	var (
		appKey    string
		appSecret string
	)
	if value, ok := md["app_key"]; ok {
		appKey = value[0]
	}
	if value, ok := md["app_secret"]; ok {
		appSecret = value[0]
	}

	if appKey != "example" || appSecret != "123456" {
		return status.Errorf(codes.Unauthenticated, "自定义认证 Token 无效")
	}

	return nil
}

type SearchService struct{}

func (s *SearchService) Search(ctx context.Context, r *pb.SearchRequest) (*pb.SearchResponse, error) {
	log.Println("get request: ", r.GetRequest())
	return &pb.SearchResponse{Response: "Hello, " + r.GetRequest() + " Server"}, nil
}

const PORT = "9001"
var auth *Auth

func main()  {
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

	auth = &Auth{}
	// 创建 grpc Server 对象
	server := grpc.NewServer(grpc.Creds(c), grpc.UnaryInterceptor(AuthenticateInterceptor))
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

func AuthenticateInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("Auth method: %s, %v", info.FullMethod, req)
	if err := auth.Check(ctx); err != nil {
		return nil, err
	}
	resp, err := handler(ctx, req)
	log.Printf("gRPC method: %s, %v", info.FullMethod, resp)
	return resp, err
}
