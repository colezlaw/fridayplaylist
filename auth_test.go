package fridayplaylist

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestGetToken(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		io.Copy(w, strings.NewReader(`{"access_token":"SOME_ACCESS_TOKEN","token_type":"Bearer","expires_in":3600}`))
	}))
	defer s.Close()

	c := &Client{
		auth_url: s.URL,
	}
	err := c.GetToken(os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"))
	if err != nil {
		t.Fatalf("unable to get token: %v", err)
	}
	if c.access_token != "SOME_ACCESS_TOKEN" {
		t.Errorf("got: %s want: SOME_ACCESS_TOKEN", c.access_token)
	}
}
