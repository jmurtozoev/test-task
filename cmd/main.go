package main

import (
	"github.com/jmurtozoev/test-task/proto"
	"github.com/jmurtozoev/test-task/server"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	// Create a listener on TCP port
	lis, err := net.Listen("tcp", ":8090")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	// Create a gRPC server object
	grpcServer := grpc.NewServer()

	// Initialize server
	s := server.New(&server.ServerOptions{})

	// Attach the product service to the grpc server
	proto.RegisterProductServiceServer(grpcServer, s)

	// Serve gRPC server
	log.Println("Serving gRPC on 0.0.0.0:8090")

	log.Fatalln(grpcServer.Serve(lis))
}
