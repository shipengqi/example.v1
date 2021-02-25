package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	stdopentracing "github.com/opentracing/opentracing-go"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	zipkinreporter "github.com/openzipkin/zipkin-go/reporter"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	pb "github.com/shipengqi/example.v1/grpc/01_simple/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	PORT = "9001"

	SERVICE_NAME              = "zipkin_server"
	ZIPKIN_HTTP_ENDPOINT      = "http://shccdfrh75vm8.hpeswlab.net:9411/api/v2/spans"
	ZIPKIN_RECORDER_HOST_PORT = "localhost:8081"
)

func main() {
	tracer, repoter := getTracer()
	defer repoter.Close()

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
	conn, err := grpc.Dial(fmt.Sprintf(":%s", PORT),
		grpc.WithTransportCredentials(c),
		grpc.WithUnaryInterceptor(
			otgrpc.OpenTracingClientInterceptor(tracer, otgrpc.LogPayloads()),
		),
	)
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
	}
	defer conn.Close()
	// 创建 SearchService 的 Client
	client := pb.NewSearchServiceClient(conn)
	// 发送 RPC 请求
	resp, err := client.Search(context.Background(), &pb.SearchRequest{Request: "gRPC Zipkin"})
	if err != nil {
		log.Fatalf("client.Search err: %v", err)
	}
	log.Printf("resp: %s", resp.GetResponse())
}

func getTracer() (stdopentracing.Tracer, zipkinreporter.Reporter) {
	// set up a span reporter
	reporter := zipkinhttp.NewReporter(ZIPKIN_HTTP_ENDPOINT)
	// 在这里的 defer 在函数返回时错误的关闭 repoter
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
