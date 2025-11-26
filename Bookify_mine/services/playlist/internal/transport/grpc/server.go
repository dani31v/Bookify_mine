package grpc

import (
	"context"
	"errors"

	playlistgrpc "bookify/proto/playlist_proto/transport/grpc"

	"bookify/services/playlist/internal/controller"
	"bookify/services/playlist/internal/model"
)

type PlaylistServer struct {
	playlistgrpc.UnimplementedPlaylistServiceServer
	ctrl *controller.Playlists
}

func NewPlaylistServer(ctrl *controller.Playlists) *PlaylistServer {
	return &PlaylistServer{ctrl: ctrl}
}

// ------------ Helpers (MODEL <-> PROTO) ------------

func songModelToProto(s model.Song) *playlistgrpc.Song {
	return &playlistgrpc.Song{
		Title:  s.Title,
		Artist: s.Artist,
	}
}

func songProtoToModel(s *playlistgrpc.Song) model.Song {
	return model.Song{
		Title:  s.GetTitle(),
		Artist: s.GetArtist(),
	}
}

func playlistModelToProto(p model.Playlist) *playlistgrpc.Playlist {
	protoPlaylist := &playlistgrpc.Playlist{
		BookId: p.BookID,
	}

	// Convert tracks
	for _, t := range p.Tracks {
		protoPlaylist.Tracks = append(protoPlaylist.Tracks, songModelToProto(t))
	}

	return protoPlaylist
}

func playlistProtoToModel(p *playlistgrpc.Playlist) model.Playlist {
	var tracks []model.Song

	for _, t := range p.GetTracks() {
		tracks = append(tracks, songProtoToModel(t))
	}

	return model.Playlist{
		BookID: p.GetBookId(),
		Tracks: tracks,
	}
}

// ------------ RPC: GetPlaylistForBook ------------

func (s *PlaylistServer) GetPlaylistForBook(ctx context.Context, req *playlistgrpc.GetPlaylistForBookRequest) (*playlistgrpc.GetPlaylistForBookResponse, error) {

	if req.GetBookId() == "" {
		return nil, errors.New("book_id is required")
	}

	// controller ahora devuelve directamente model.Playlist
	playlist := s.ctrl.ForBook(req.GetBookId())

	return &playlistgrpc.GetPlaylistForBookResponse{
		Playlist: playlistModelToProto(playlist),
	}, nil
}

// ------------ RPC: CreatePlaylist ------------

func (s *PlaylistServer) CreatePlaylist(ctx context.Context, req *playlistgrpc.CreatePlaylistRequest) (*playlistgrpc.CreatePlaylistResponse, error) {

	if req.GetPlaylist() == nil {
		return nil, errors.New("playlist is required")
	}

	newPlaylist := playlistProtoToModel(req.GetPlaylist())

	created, err := s.ctrl.Create(newPlaylist)
	if err != nil {
		return nil, err
	}

	return &playlistgrpc.CreatePlaylistResponse{
		Playlist: playlistModelToProto(created),
	}, nil
}
