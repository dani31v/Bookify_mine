package grpc

import (
	"context"
	"errors"

	reviewsgrpc "bookify/proto/reviews_proto/transport/grpc"

	"bookify/services/reviews/internal/controller"
	"bookify/services/reviews/internal/model"
)

type ReviewsServer struct {
	reviewsgrpc.UnimplementedReviewServiceServer
	ctrl *controller.Reviews
}

func NewReviewsServer(ctrl *controller.Reviews) *ReviewsServer {
	return &ReviewsServer{ctrl: ctrl}
}

func reviewModelToProto(r model.Review) *reviewsgrpc.Review {
	return &reviewsgrpc.Review{
		Id:     r.ID,
		BookId: r.BookID,
		UserId: r.UserID,
		Rating: int32(r.Rating),
		Text:   r.Text,
	}
}

func reviewProtoToModel(r *reviewsgrpc.Review) model.Review {
	return model.Review{
		ID:     r.GetId(),
		BookID: r.GetBookId(),
		UserID: r.GetUserId(),
		Rating: int(r.GetRating()),
		Text:   r.GetText(),
	}
}

func (s *ReviewsServer) GetReviews(ctx context.Context, req *reviewsgrpc.GetReviewsRequest) (*reviewsgrpc.GetReviewsResponse, error) {

	reviews := s.ctrl.ByBook(req.GetBookId())

	protoReviews := make([]*reviewsgrpc.Review, 0, len(reviews))
	for _, r := range reviews {
		protoReviews = append(protoReviews, &reviewsgrpc.Review{
			Id:     r.ID,
			BookId: r.BookID,
			UserId: r.UserID,
			Rating: int32(r.Rating),
			Text:   r.Text,
		})
	}

	return &reviewsgrpc.GetReviewsResponse{
		Reviews: protoReviews,
	}, nil
}

func (s *ReviewsServer) CreateReview(ctx context.Context, req *reviewsgrpc.CreateReviewRequest) (*reviewsgrpc.CreateReviewResponse, error) {

	if req.GetReview() == nil {
		return nil, errors.New("review is required")
	}

	newReview := reviewProtoToModel(req.GetReview())

	created, err := s.ctrl.Create(newReview)
	if err != nil {
		return nil, err
	}

	return &reviewsgrpc.CreateReviewResponse{
		Review: reviewModelToProto(created),
	}, nil
}
