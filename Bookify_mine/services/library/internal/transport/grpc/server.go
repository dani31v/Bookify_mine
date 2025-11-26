package grpc

import (
	"context"
	"errors"

	librarygrpc "bookify/proto/library_proto/transport/grpc"

	"bookify/services/library/internal/controller"
	"bookify/services/library/internal/model"
)

type LibraryServer struct {
	librarygrpc.UnimplementedLibraryServiceServer
	ctrl *controller.Books
}

func NewLibraryServer(ctrl *controller.Books) *LibraryServer {
	return &LibraryServer{ctrl: ctrl}
}

func (s *LibraryServer) GetBook(ctx context.Context, req *librarygrpc.GetBookRequest) (*librarygrpc.GetBookResponse, error) {
	data, ok := s.ctrl.Get(req.Id)
	if !ok {
		return nil, errors.New("book not found")
	}

	book, ok := data.(model.Book)
	if !ok {
		return nil, errors.New("invalid book data")
	}

	return &librarygrpc.GetBookResponse{
		Book: &librarygrpc.Book{
			Id:      book.ID,
			Title:   book.Title,
			Author:  book.Author,
			Pages:   int32(book.Pages),
			Edition: book.Edition,
		},
	}, nil
}

func (s *LibraryServer) CreateBook(ctx context.Context, req *librarygrpc.CreateBookRequest) (*librarygrpc.CreateBookResponse, error) {
	b := req.Book

	newBook := model.Book{
		ID:      b.Id,
		Title:   b.Title,
		Author:  b.Author,
		Pages:   int(b.Pages),
		Edition: b.Edition,
	}

	created, err := s.ctrl.Create(newBook)
	if err != nil {
		return nil, err
	}

	return &librarygrpc.CreateBookResponse{
		Book: &librarygrpc.Book{
			Id:      created.ID,
			Title:   created.Title,
			Author:  created.Author,
			Pages:   int32(created.Pages),
			Edition: created.Edition,
		},
	}, nil
}
