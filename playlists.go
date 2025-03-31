package fridayplaylist

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Playlist struct {
	Href string `json:"href"`
	Name string `json:"name"`
	URL  string `json:"URL"`
}

type UserPlaylistResponse struct {
	Next  string     `json:"next"`
	Items []Playlist `json:"items"`
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
