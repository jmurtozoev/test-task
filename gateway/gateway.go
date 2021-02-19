package main

import (
	"context"
	"flag"
	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jmurtozoev/test-task/proto"
	"google.golang.org/grpc"
	"net/http"
	"log"
)

var (
	// command-line options:
	// gRPC server endpoint
	serverEndpoint = flag.String("grpc-server-endpoint", "localhost:8090", "gRPC server endpoint")
)

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Register gRPC server endpoint
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	err := proto.RegisterProductServiceHandlerFromEndpoint(ctx, mux, *serverEndpoint, opts)
	if err != nil {
		return err
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	return http.ListenAndServe(":8080", mux)
}

func main() {
	flag.Parse()
	defer glog.Flush()

	log.Println("grpc-gateway is running on port 8080")
	if err := run(); err != nil {
		glog.Fatal(err)
	}
}