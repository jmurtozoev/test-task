package server

import (
	"fmt"
	"github.com/jmurtozoev/test-task/proto"
	"context"
)

type ServerOptions struct {

}

type server struct{
	proto.UnimplementedProductServiceServer

}

func New(opts *ServerOptions) *server {
	return &server{

	}
}

func (s *server) ReadProducts(ctx context.Context, req *proto.ReadFromCsvRequest) (*proto.Nothing, error) {
	fmt.Println("Hey buddy, wazzup!")
	return &proto.Nothing{}, nil
}

func (s *server) ListProducts(ctx context.Context, req *proto.Nothing) (*proto.ListProductsResponse, error) {

	return &proto.ListProductsResponse{}, nil
}