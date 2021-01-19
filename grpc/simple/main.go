package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main()  {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	grpcServer.Serve(lis)
}
