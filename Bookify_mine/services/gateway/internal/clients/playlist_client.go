package clients

import (
	"context"
	"fmt"

	playlistgrpc "bookify/proto/playlist_proto/transport/grpc"

	"google.golang.org/grpc"
)

type PlaylistClient struct {
	client playlistgrpc.PlaylistServiceClient
}

func NewPlaylistClient(addr string) (*PlaylistClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("cannot dial playlist grpc: %w", err)
	}

	return &PlaylistClient{
		client: playlistgrpc.NewPlaylistServiceClient(conn),
	}, nil
}

func (c *PlaylistClient) GetPlaylistForBook(ctx context.Context, bookID string) (*playlistgrpc.Playlist, error) {
	resp, err := c.client.GetPlaylistForBook(ctx, &playlistgrpc.GetPlaylistForBookRequest{
		BookId: bookID,
	})
	if err != nil {
		return nil, err
	}
	return resp.GetPlaylist(), nil
}

func (c *PlaylistClient) CreatePlaylist(ctx context.Context, p *playlistgrpc.Playlist) (*playlistgrpc.Playlist, error) {
	resp, err := c.client.CreatePlaylist(ctx, &playlistgrpc.CreatePlaylistRequest{
		Playlist: p,
	})
	if err != nil {
		return nil, err
	}
	return resp.GetPlaylist(), nil
}
