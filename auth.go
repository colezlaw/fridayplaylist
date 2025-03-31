package fridayplaylist

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	base_url     string
	auth_url     string
	access_token string
}

func (c *Client) GetToken(clientID, clientSecret string) error {
	if c.auth_url == "" {
		c.auth_url = "https://accounts.spotify.com/api/token"
	}

	param := url.Values{
		"grant_type": []string{"client_credentials"},
	}
	req, err := http.NewRequest(http.MethodPost, c.auth_url, strings.NewReader(param.Encode()))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(clientID, clientSecret)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("unable to send request: %w", err)
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("got unexpected return code: %d", res.StatusCode)
	}
	result := struct {
		AccessToken string `json:"access_token"`
	}{}
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return fmt.Errorf("unable to decode: %w", err)
	}

	c.access_token = result.AccessToken
	return nil
}
