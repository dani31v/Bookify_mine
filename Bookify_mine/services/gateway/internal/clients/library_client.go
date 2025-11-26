package clients

import (
	"context"
	"fmt"

	librarygrpc "bookify/proto/library_proto/transport/grpc"

	"google.golang.org/grpc"
)

type LibraryClient struct {
	client librarygrpc.LibraryServiceClient
}

func NewLibraryClient(addr string) (*LibraryClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to dial library gRPC: %w", err)
	}

	return &LibraryClient{
		client: librarygrpc.NewLibraryServiceClient(conn),
	}, nil
}

func (c *LibraryClient) GetBook(ctx context.Context, id string) (*librarygrpc.Book, error) {
	req := &librarygrpc.GetBookRequest{Id: id}
	resp, err := c.client.GetBook(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Book, nil
}
