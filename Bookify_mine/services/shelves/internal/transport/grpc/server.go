package grpc

import (
	"context"
	"errors"

	shelvesgrpc "bookify/proto/shelves_proto/transport/grpc"

	"bookify/services/shelves/internal/controller"
	"bookify/services/shelves/internal/model"
)

type ShelvesServer struct {
	shelvesgrpc.UnimplementedShelvesServiceServer
	ctrl *controller.Shelves
}

func NewShelvesServer(ctrl *controller.Shelves) *ShelvesServer {
	return &ShelvesServer{ctrl: ctrl}
}

func shelfModelToProto(s model.Shelf) shelvesgrpc.Shelf {
	switch s {
	case model.ToRead:
		return shelvesgrpc.Shelf_SHELF_TO_READ
	case model.Reading:
		return shelvesgrpc.Shelf_SHELF_READING
	case model.Done:
		return shelvesgrpc.Shelf_SHELF_FINISHED
	default:
		return shelvesgrpc.Shelf_SHELF_UNSPECIFIED
	}
}

func shelfProtoToModel(s shelvesgrpc.Shelf) model.Shelf {
	switch s {
	case shelvesgrpc.Shelf_SHELF_TO_READ:
		return model.ToRead
	case shelvesgrpc.Shelf_SHELF_READING:
		return model.Reading
	case shelvesgrpc.Shelf_SHELF_FINISHED:
		return model.Done
	default:
		return model.ToRead
	}
}

func (s *ShelvesServer) GetShelfItem(ctx context.Context, req *shelvesgrpc.GetShelfItemRequest) (*shelvesgrpc.GetShelfItemResponse, error) {
	if req.GetUserId() == "" || req.GetBookId() == "" {
		return nil, errors.New("user_id and book_id are required")
	}

	data, ok := s.ctrl.Get(req.GetUserId(), req.GetBookId())
	if !ok {
		return nil, errors.New("shelf item not found")
	}

	item := data

	return &shelvesgrpc.GetShelfItemResponse{
		Item: &shelvesgrpc.ShelfItem{
			UserId: item.UserID,
			BookId: item.BookID,
			Shelf:  shelfModelToProto(item.Shelf),
		},
	}, nil
}

func (s *ShelvesServer) CreateShelfItem(ctx context.Context, req *shelvesgrpc.CreateShelfItemRequest) (*shelvesgrpc.CreateShelfItemResponse, error) {
	if req.GetItem() == nil {
		return nil, errors.New("item is required")
	}

	in := req.GetItem()

	newItem := model.ShelfItem{
		UserID: in.GetUserId(),
		BookID: in.GetBookId(),
		Shelf:  shelfProtoToModel(in.GetShelf()),
	}

	created, err := s.ctrl.Create(newItem)
	if err != nil {
		return nil, err
	}

	return &shelvesgrpc.CreateShelfItemResponse{
		Item: &shelvesgrpc.ShelfItem{
			UserId: created.UserID,
			BookId: created.BookID,
			Shelf:  shelfModelToProto(created.Shelf),
		},
	}, nil
}
