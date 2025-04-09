package main

import (
	"bytes"
	"testing"

	fridayplaylist "github.com/colezlaw/fridayPlaylist"
)

type MockPlaylistClient struct {
	getPlaylists func(string) ([]fridayplaylist.Playlist, error)
	getTracks    func(string) ([]fridayplaylist.Track, error)
}

func (m *MockPlaylistClient) GetPlaylistsForUser(user string) ([]fridayplaylist.Playlist, error) {
	if m.getPlaylists == nil {
		return []fridayplaylist.Playlist{}, nil
	}

	return m.getPlaylists(user)
}

func (m *MockPlaylistClient) GetTracksForPlaylist(playlistID string) ([]fridayplaylist.Track, error) {
	if m.getTracks == nil {
		return []fridayplaylist.Track{}, nil
	}

	return m.getTracks(playlistID)
}

func TestRun(t *testing.T) {
	// Arrange
	mock := &MockPlaylistClient{}
	mock.getPlaylists = func(u string) ([]fridayplaylist.Playlist, error) {
		return []fridayplaylist.Playlist{
			{Name: "As the World Turns", ID: "1234"},
		}, nil
	}
	mock.getTracks = func(p string) ([]fridayplaylist.Track, error) {
		return []fridayplaylist.Track{
			{
				Name:   "As the Worm Turns",
				Artist: "Faith No More",
				Album: fridayplaylist.Album{
					Name:        "Faith No More",
					ReleaseDate: "1985",
				},
			}}, nil
	}
	var buff *bytes.Buffer

	// Act
	if err := run([]string{"-f", "test_out.csv", "-user", "c0leslaw"}, buff, mock); err != nil {
		t.Fatalf("unexpected error running run: %v", err)
	}
	t.Logf("%s", buff)
}
