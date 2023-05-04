package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	stdopentracing "github.com/opentracing/opentracing-go"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	zipkinreporter "github.com/openzipkin/zipkin-go/reporter"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	pb "github.com/shipengqi/example.v1/packeges/third/grpc/03_simple_tls/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type SearchService struct{}

func (s *SearchService) Search(ctx context.Context, r *pb.SearchRequest) (*pb.SearchResponse, error) {
	log.Println("get request: ", r.GetRequest())
	return &pb.SearchResponse{Response: "Hello, " + r.GetRequest() + " Server"}, nil
}

const (
	PORT = "9001"

	SERVICE_NAME              = "zipkin_server"
	ZIPKIN_HTTP_ENDPOINT      = "http://shccdfrh75vm8.hpeswlab.net:9411/api/v2/spans"
	ZIPKIN_RECORDER_HOST_PORT = "localhost:8081"
)

func main() {
	tracer, repoter := getTracer()
	defer repoter.Close()

	// do other bootstrapping stuff...
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
	server := grpc.NewServer(
		grpc.Creds(c),
		grpc.UnaryInterceptor(
			otgrpc.OpenTracingServerInterceptor(tracer, otgrpc.LogPayloads()),
		),
	)
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

func getTracer() (stdopentracing.Tracer, zipkinreporter.Reporter) {
	// set up a span reporter
	reporter := zipkinhttp.NewReporter(ZIPKIN_HTTP_ENDPOINT)
	// 注意，在这里的 defer ，会在函数返回时错误的关闭 repoter
	// defer reporter.Close()

	// create our local service endpoint
	endpoint, err := zipkin.NewEndpoint(SERVICE_NAME, ZIPKIN_RECORDER_HOST_PORT)
	if err != nil {
		log.Fatalf("unable to create local endpoint: %+v\n", err)
	}

	// set-up our sampling strategy
	// sampler := zipkin.NewModuloSampler(1)
	// initialize the tracer
	zipkinTracer, err := zipkin.NewTracer(
		reporter,
		zipkin.WithLocalEndpoint(endpoint),
		// zipkin.WithSampler(sampler),
	)

	if err != nil {
		log.Fatalf("unable to create tracer: %+v\n", err)
	}

	// use zipkin-go-opentracing to wrap our tracer
	tracer := zipkinot.Wrap(zipkinTracer)

	// optionally set as Global OpenTracing tracer instance
	// opentracing.SetGlobalTracer(tracer)

	return tracer, reporter
}
