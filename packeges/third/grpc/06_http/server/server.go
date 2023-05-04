package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/pkg/errors"
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
	c, err := GetTLSCredentialsByCA()
	if err != nil {
		log.Fatalf("GetTLSCredentialsByCA err: %v", err)
	}

	mux := GetHTTPServeMux()
	// 创建 grpc Server 对象
	server := grpc.NewServer(grpc.Creds(c))
	// 注册 SearchService
	pb.RegisterSearchServiceServer(server, &SearchService{})

	certFile := "../../ssl/server/server.crt"
	keyFile := "../../ssl/server/server.key"

	http.ListenAndServeTLS(fmt.Sprintf(":%s", PORT), certFile, keyFile,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 检测请求协议是否为 HTTP/2
			// 判断 Content-Type 是否为 application/grpc（gRPC 的默认标识位）
			// 根据协议的不同转发到不同的服务处理
			if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
				server.ServeHTTP(w, r)
			} else {
				mux.ServeHTTP(w, r)
			}

			return
		}),
	)
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

func GetTLSCredentialsByCA() (credentials.TransportCredentials, error) {
	cert, err := tls.LoadX509KeyPair("../../ssl/server/server.crt", "../../ssl/server/server.key")
	if err != nil {
		return nil, errors.Wrap(err, "tls.LoadX509KeyPair")
	}

	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("../../ssl/ca/ca.crt")
	if err != nil {
		return nil, errors.Wrap(err, "ioutil.ReadFile")
	}

	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		return nil, errors.New("certPool.AppendCertsFromPEM")
	}

	c := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	})

	return c, nil
}

func GetHTTPServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("example.v1: grpc-example"))
	})

	return mux
}
