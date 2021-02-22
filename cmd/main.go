package main

import (
	"flag"
	"github.com/jmurtozoev/test-task/db"
	"github.com/jmurtozoev/test-task/proto"
	"github.com/jmurtozoev/test-task/server"
	"github.com/jmurtozoev/test-task/storage"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"
)

var (
	grpcPort = flag.Int("port", 8090, "server address")
)

func main() {
	flag.Parse()
	port := strconv.Itoa(*grpcPort)

	// connect to db
	database := db.Connect()
	store := storage.New(database)

	// Create a listener on TCP port
	lis, err := net.Listen("tcp", ":" + port)

	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	// Create a gRPC server object
	grpcServer := grpc.NewServer()

	// Initialize server
	s := server.New(&server.Options{
		Storage: store,
	})

	// Attach the product service to the grpc server
	proto.RegisterProductServiceServer(grpcServer, s)

	// Serve gRPC server
	log.Println("Serving gRPC on 0.0.0.0:8090")
	log.Fatalln(grpcServer.Serve(lis))
}
