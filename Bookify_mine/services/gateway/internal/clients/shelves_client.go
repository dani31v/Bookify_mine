package clients

import (
	"context"
	"fmt"

	shelvesgrpc "bookify/proto/shelves_proto/transport/grpc"

	"google.golang.org/grpc"
)

type ShelvesClient struct {
	client shelvesgrpc.ShelvesServiceClient
}

func NewShelvesClient(addr string) (*ShelvesClient, error) {
	conn, err := grpc.Dial(addr,
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(10*1024*1024)),
	)
	if err != nil {
		return nil, fmt.Errorf("cannot dial shelves grpc: %w", err)
	}

	return &ShelvesClient{
		client: shelvesgrpc.NewShelvesServiceClient(conn),
	}, nil
}

func (c *ShelvesClient) GetShelfItem(ctx context.Context, userID, bookID string) (*shelvesgrpc.ShelfItem, error) {
	resp, err := c.client.GetShelfItem(ctx, &shelvesgrpc.GetShelfItemRequest{
		UserId: userID,
		BookId: bookID,
	})
	if err != nil {
		return nil, err
	}
	return resp.GetItem(), nil
}

func (c *ShelvesClient) CreateShelfItem(ctx context.Context, item *shelvesgrpc.ShelfItem) (*shelvesgrpc.ShelfItem, error) {
	resp, err := c.client.CreateShelfItem(ctx, &shelvesgrpc.CreateShelfItemRequest{
		Item: item,
	})
	if err != nil {
		return nil, err
	}
	return resp.GetItem(), nil
}
