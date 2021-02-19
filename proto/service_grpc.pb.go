// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ProductServiceClient is the client API for ProductService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ProductServiceClient interface {
	// read products from csv file
	ReadProducts(ctx context.Context, in *ReadFromCsvRequest, opts ...grpc.CallOption) (*Nothing, error)
	// list products
	ListProducts(ctx context.Context, in *Nothing, opts ...grpc.CallOption) (*ListProductsResponse, error)
}

type productServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewProductServiceClient(cc grpc.ClientConnInterface) ProductServiceClient {
	return &productServiceClient{cc}
}

func (c *productServiceClient) ReadProducts(ctx context.Context, in *ReadFromCsvRequest, opts ...grpc.CallOption) (*Nothing, error) {
	out := new(Nothing)
	err := c.cc.Invoke(ctx, "/product.service.ProductService/ReadProducts", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *productServiceClient) ListProducts(ctx context.Context, in *Nothing, opts ...grpc.CallOption) (*ListProductsResponse, error) {
	out := new(ListProductsResponse)
	err := c.cc.Invoke(ctx, "/product.service.ProductService/ListProducts", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProductServiceServer is the server API for ProductService service.
// All implementations must embed UnimplementedProductServiceServer
// for forward compatibility
type ProductServiceServer interface {
	// read products from csv file
	ReadProducts(context.Context, *ReadFromCsvRequest) (*Nothing, error)
	// list products
	ListProducts(context.Context, *Nothing) (*ListProductsResponse, error)
	mustEmbedUnimplementedProductServiceServer()
}

// UnimplementedProductServiceServer must be embedded to have forward compatible implementations.
type UnimplementedProductServiceServer struct {
}

func (UnimplementedProductServiceServer) ReadProducts(context.Context, *ReadFromCsvRequest) (*Nothing, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadProducts not implemented")
}
func (UnimplementedProductServiceServer) ListProducts(context.Context, *Nothing) (*ListProductsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListProducts not implemented")
}
func (UnimplementedProductServiceServer) mustEmbedUnimplementedProductServiceServer() {}

// UnsafeProductServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ProductServiceServer will
// result in compilation errors.
type UnsafeProductServiceServer interface {
	mustEmbedUnimplementedProductServiceServer()
}

func RegisterProductServiceServer(s grpc.ServiceRegistrar, srv ProductServiceServer) {
	s.RegisterService(&ProductService_ServiceDesc, srv)
}

func _ProductService_ReadProducts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadFromCsvRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProductServiceServer).ReadProducts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/product.service.ProductService/ReadProducts",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProductServiceServer).ReadProducts(ctx, req.(*ReadFromCsvRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProductService_ListProducts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Nothing)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProductServiceServer).ListProducts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/product.service.ProductService/ListProducts",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProductServiceServer).ListProducts(ctx, req.(*Nothing))
	}
	return interceptor(ctx, in, info, handler)
}

// ProductService_ServiceDesc is the grpc.ServiceDesc for ProductService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ProductService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "product.service.ProductService",
	HandlerType: (*ProductServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ReadProducts",
			Handler:    _ProductService_ReadProducts_Handler,
		},
		{
			MethodName: "ListProducts",
			Handler:    _ProductService_ListProducts_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service.proto",
}