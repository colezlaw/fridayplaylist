package fridayplaylist

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGetPlaylistForUser(t *testing.T) {
	b, err := os.ReadFile("testdata/playlist_response.json")
	if err != nil {
		t.Fatalf("unable to read response: %v", err)
	}
	responses := make([]json.RawMessage, 0)
	if err := json.Unmarshal(b, &responses); err != nil {
		t.Fatalf("unable to unmarshal responses")
	}
	if len(responses) != 2 {
		t.Errorf("expected 2 responses, got %d", len(responses))
	}

	rc := 0
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(responses[rc])
		rc++
	}))
	defer s.Close()
	responses[0] = json.RawMessage(bytes.ReplaceAll(responses[0], []byte("@@TEST_SERVER_URL@@"), []byte(s.URL)))

	c := &Client{
		access_token: "SOME_ACCESS_TOKEN",
		base_url:     s.URL + "/",
	}
	playlists, err := c.GetPlaylistsForUser("someuser")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(playlists) != 2 {
		t.Errorf("expected 2 playlists, got %d", len(playlists))
	}
	if playlists[0].Name != "Condiments" {
		t.Errorf(`expected "Condiments", got %q`, playlists[0].Name)
	}
}

func TestGetTracks(t *testing.T) {
	b, err := os.ReadFile("testdata/tracks_response.json")
	if err != nil {
		t.Fatalf("unable to read response: %v", err)
	}
	responses := make([]json.RawMessage, 0)
	if err := json.Unmarshal(b, &responses); err != nil {
		t.Fatalf("unable to unmarshal responses")
	}
	if len(responses) != 2 {
		t.Errorf("expected 2 responses, got %d", len(responses))
	}

	rc := 0
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(responses[rc])
		rc++
	}))
	defer s.Close()
	responses[0] = json.RawMessage(bytes.ReplaceAll(responses[0], []byte("@@TEST_SERVER_URL@@"), []byte(s.URL)))

	c := &Client{
		access_token: "SOME_ACCESS_TOKEN",
		base_url:     s.URL + "/",
	}
	tracks, err := c.GetTracksForPlaylist("someplaylist")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tracks) != 2 {
		t.Errorf("expected 2 playlists, got %d", len(tracks))
	}
	if tracks[0].Name != "The Ketchup Song (Aserej├⌐) - Motown Club Single Edit" {
		t.Errorf(`expected "The Ketchup Song (Aserej├⌐) - Motown Club Single Edit", got %q`, tracks[0].Name)
	}
}
