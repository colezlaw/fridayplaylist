package fridayplaylist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Playlist struct {
	Href string `json:"href"`
	Name string `json:"name"`
	URI  string `json:"uri"`
	ID   string `json:"ID"`
}

type UserPlaylistResponse struct {
	Next  string     `json:"next"`
	Items []Playlist `json:"items"`
}

type Album struct {
	Name        string `json:"name"`
	ReleaseDate string `json:"release_date"`
}

type Artist struct {
	Name string `json:"name"`
}

type Track struct {
	Name    string   `json:"name"`
	Artists []Artist `json:"artists"`
	Album   Album    `json:"album"`
}

type GetPlaylistTracksResult struct {
	Next  string `json:"next"`
	Items []struct {
		Track Track `json:"track"`
	}
}

func (c *Client) GetTracksForPlaylist(id string) ([]Track, error) {
	tracks := make([]Track, 0)

	if c.base_url == "" {
		c.base_url = "https://api.spotify.com/v1/"
	}

	if c.access_token == "" {
		return tracks, fmt.Errorf("you must get a token. Call GetToken")
	}

	next := c.base_url + "playlists/" + id + "/tracks?limit=50"

	for next != "" {
		req, err := http.NewRequest(http.MethodGet, next, nil)
		if err != nil {
			return tracks, fmt.Errorf("unable to create request: %w", err)
		}
		req.Header.Set("Authorization", "Bearer "+c.access_token)
		req.Header.Set("Accept", "application/json")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return tracks, fmt.Errorf("do: %w", err)
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			io.Copy(os.Stderr, res.Body)
			return tracks, fmt.Errorf("unexpected error code: %d", res.StatusCode)
		}
		buff := &bytes.Buffer{}
		io.Copy(buff, res.Body)
		body := GetPlaylistTracksResult{}
		if err := json.NewDecoder(buff).Decode(&body); err != nil {
			return tracks, fmt.Errorf("decode: %w", err)
		}
		for _, track := range body.Items {
			tracks = append(tracks, track.Track)
		}

		next = body.Next
	}

	return tracks, nil
}

func (c *Client) GetPlaylistsForUser(username string) ([]Playlist, error) {
	result := make([]Playlist, 0)

	if c.access_token == "" {
		return result, fmt.Errorf("you must authenticate first - call GetToken")
	}

	if c.base_url == "" {
		c.base_url = "https://api.spotify.com/v1/"
	}

	next := c.base_url + "users/" + username + "/playlists?limit=50"
	for next != "" {
		req, err := http.NewRequest(http.MethodGet, next, nil)
		if err != nil {
			return []Playlist{}, fmt.Errorf("unable to create request: %w", err)
		}
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Authorization", "Bearer "+c.access_token)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return result, fmt.Errorf("unable to request: %w", err)
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			io.Copy(os.Stderr, res.Body)
			return result, fmt.Errorf("unexpected return code: %d", res.StatusCode)
		}
		body := UserPlaylistResponse{}
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			return result, fmt.Errorf("error decoding: %w", err)
		}
		result = append(result, body.Items...)
		next = body.Next
	}

	return result, nil
}
