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
	// Arrange
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
	gotToken := ""
	wantToken := "WANT_TOKEN"

	rc := 0
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotToken = r.Header.Get("Authorization")

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(responses[rc])
		rc++
	}))
	defer s.Close()
	responses[0] = json.RawMessage(bytes.ReplaceAll(responses[0], []byte("@@TEST_SERVER_URL@@"), []byte(s.URL)))

	c := &Client{
		access_token: wantToken,
		base_url:     s.URL + "/",
	}

	// Act
	playlists, err := c.GetPlaylistsForUser("someuser")

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(playlists) != 2 {
		t.Errorf("expected 2 playlists, got %d", len(playlists))
	}
	if playlists[0].Name != "Condiments" {
		t.Errorf(`expected "Condiments", got %q`, playlists[0].Name)
	}
	if gotToken != "Bearer "+wantToken {
		t.Errorf("expected authorization %q, got %q", gotToken, "Bearer "+wantToken)
	}
}

func TestGetTracks(t *testing.T) {
	// Arrange
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
	gotToken := ""
	wantToken := "WANT_TOKEN"

	rc := 0
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotToken = r.Header.Get("Authorization")

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(responses[rc])
		rc++
	}))
	defer s.Close()
	responses[0] = json.RawMessage(bytes.ReplaceAll(responses[0], []byte("@@TEST_SERVER_URL@@"), []byte(s.URL)))

	c := &Client{
		access_token: wantToken,
		base_url:     s.URL + "/",
	}

	// Act
	tracks, err := c.GetTracksForPlaylist("someplaylist")

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tracks) != 2 {
		t.Errorf("expected 2 playlists, got %d", len(tracks))
	}
	if tracks[0].Name != "The Ketchup Song (Aserej├⌐) - Motown Club Single Edit" {
		t.Errorf(`expected "The Ketchup Song (Aserej├⌐) - Motown Club Single Edit", got %q`, tracks[0].Name)
	}
	if gotToken != "Bearer "+wantToken {
		t.Errorf("expected authorization %q, got %q", gotToken, "Bearer "+wantToken)
	}
}

func TestUnmarshalTrack(t *testing.T) {
	// Arrange
	data := []byte(`{"album":{"name":"Siamese Dream (Deluxe Edition)","release_date":"1993"},
"artists":[{"name": "The Smashing Pumpkins"},{"name": "Andy Griffith"}],"name":"Mayonaise - 2011 Remaster"
}`)
	got := Track{}
	want := "The Smashing Pumpkins,Andy Griffith"

	// Act
	err := json.Unmarshal(data, &got)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error unmarshaling: %v", err)
	}
	if got.Artist != want {
		t.Errorf("got %q, want %q", got.Artist, want)
	}
}
