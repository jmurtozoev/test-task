package server

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/jmurtozoev/test-task/models"
	"github.com/jmurtozoev/test-task/proto"
	"github.com/jmurtozoev/test-task/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"net/url"
	"strconv"
)

type Options struct {
	Storage storage.Storage
}

type server struct {
	proto.UnimplementedProductServiceServer
	store storage.Storage
}

func New(opts *Options) *server {
	return &server{
		store: opts.Storage,
	}
}

func (s *server) ReadProducts(ctx context.Context, req *proto.ReadFromCsvRequest) (*proto.Nothing, error) {
	u, err := url.ParseRequestURI(req.GetUrl())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid url")
	}

	// fetch data from url
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, status.Error(codes.Unknown, fmt.Sprintf("request error: %v", err))
	}
	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)
	reader.Comma = ',' // can be changed according to csv file data
	data, err := reader.ReadAll()
	if err != nil {
		return nil, status.Error(codes.Unknown, fmt.Sprintf("reading data error: %v", err))
	}

	for idx, row := range data {
		var product models.Product
		// skip header
		if idx == 0 {
			continue
		}

		product.Name = row[0]
		price, err := strconv.ParseFloat(row[1], 32)
		if err != nil {
			return nil, status.Error(codes.Unknown, fmt.Sprintf("parsing float error: %v", err))
		}

		product.Price = float32(price)
		err = s.store.Product().Create(product)
		if err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf("create product error: %v", err))
		}
	}

	return &proto.Nothing{}, nil
}

func (s *server) ListProducts(ctx context.Context, req *proto.ListProductsRequest) (*proto.ListProductsResponse, error) {
	var Products []*proto.Product
	filter := make(map[string]interface{})
	var limit, page int

	if limit = int(req.GetLimit()); limit <= 0 {
		return nil, status.Error(codes.InvalidArgument, "limit is invalid or missing")
	}

	if page = int(req.GetPage()); page <= 0 {
		return nil, status.Error(codes.InvalidArgument, "page is invalid or missing")
	}

	if req.GetName() != "" {
		filter["name"] = req.GetName()
	}

	if req.GetCostMin() > 0 {
		filter["cost_min"] = req.GetCostMin()
	}

	if req.GetCostMax() > 0 && req.GetCostMax() > req.GetCostMin() {
		filter["cost_max"] = req.GetCostMax()
	}

	products, count, err := s.store.Product().List(page, limit, filter)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("create product error: %v", err))
	}

	for _, p := range products {
		product := proto.Product{
			Id:          int32(p.ID),
			Name:        p.Name,
			Price:       p.Price,
			UpdateCount: int32(p.UpdateCount),
		}

		Products = append(Products, &product)
	}

	resp := &proto.ListProductsResponse{
		Count:    int32(count),
		Products: Products,
	}

	return resp, nil
}
