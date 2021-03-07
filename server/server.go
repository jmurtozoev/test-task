package server

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"github.com/jmurtozoev/test-task/models"
	"github.com/jmurtozoev/test-task/proto"
	"github.com/jmurtozoev/test-task/storage"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
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

func (s *server) ReadProducts(ctx context.Context, req *proto.ReadFromCsvRequest) (*proto.ReadFromCsvResponse, error) {
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

	productCount := 0
	for idx, row := range data {
		var product models.Product
		// skip header
		if idx == 0 {
			continue
		}

		// initialize product name
		if product.Name = row[0]; product.Name == "" {
			continue
		}

		// initialize product price
		price, err := strconv.ParseFloat(row[1], 32)
		if err != nil {
			log.Printf("parsing float error: %v", err)
			continue
		}

		product.Price = float32(price)
		err = s.store.Product().Create(product)
		if err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf("create product error: %v", err))
		}

		productCount++
	}

	return &proto.ReadFromCsvResponse{Inserted: int32(productCount)}, nil
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
			UpdatedAt:   p.UpdatedAt,
		}

		Products = append(Products, &product)
	}

	resp := &proto.ListProductsResponse{
		Count:    int32(count),
		Products: Products,
	}

	return resp, nil
}

func (s *server) UpdateProduct(ctx context.Context, req *proto.UpdateProductRequest) (resp *proto.Product, err error) {
	var productId int
	if productId = int(req.GetId()); productId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "id field is invalid or missing")
	}

	// get product by id
	var product *models.Product
	product, err = s.store.Product().Get(productId)
	switch err {
	case sql.ErrNoRows:
		return nil, status.Error(codes.NotFound, fmt.Sprintf("product with id=%d not found", productId))
	case nil:
		break
	default:
		return nil, status.Error(codes.Internal, fmt.Sprintf("error getting product with id=%d", productId))
	}

	if req.GetName() != "" {
		product.Name = req.GetName()
	}

	if req.GetPrice() > 0 {
		product.Price = req.GetPrice()
	}

	product.UpdateCount++

	err = s.store.Product().Update(product)
	if err != nil {
		if dbErr, ok := err.(*pq.Error); ok && dbErr.Constraint == "products_name_key" {
			return nil, status.Error(codes.AlreadyExists, fmt.Sprintf("product with name=%s already exists", product.Name))
		}

		return nil, status.Error(codes.Internal, fmt.Sprintf("update product error: %v", err))
	}

	resp = &proto.Product{
		Id:          int32(product.ID),
		Name:        product.Name,
		Price:       product.Price,
		UpdateCount: int32(product.UpdateCount),
	}

	return
}
