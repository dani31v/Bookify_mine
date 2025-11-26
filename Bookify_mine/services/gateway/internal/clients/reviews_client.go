package clients

import (
	"context"
	"fmt"

	reviewspb "bookify/proto/reviews_proto/transport/grpc"

	"google.golang.org/grpc"
)

type ReviewsClient struct {
	client reviewspb.ReviewServiceClient
}

func NewReviewsClient(addr string) (*ReviewsClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("cannot dial reviews service: %w", err)
	}

	return &ReviewsClient{
		client: reviewspb.NewReviewServiceClient(conn),
	}, nil
}

func (c *ReviewsClient) GetReviews(ctx context.Context, bookID string) (*reviewspb.GetReviewsResponse, error) {
	req := &reviewspb.GetReviewsRequest{
		BookId: bookID,
	}
	return c.client.GetReviews(ctx, req)
}
