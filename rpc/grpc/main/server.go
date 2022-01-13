package main

import (
	"log"
	"net"

	grpc "google.golang.org/grpc"
)

func main() {
	grpcServer := grpc.NewServer()
	RegisterHelloServiceServer(grpcServer, new(HelloServiceImpl))

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer.Serve(listener)
}
