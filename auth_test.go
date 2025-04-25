package fridayplaylist

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetToken(t *testing.T) {
	// Arrange
	gotAuth := ""
	wantAuth := "VFJPR0RPUjpCVVJOSU5BVE9S"
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		io.Copy(w, strings.NewReader(`{"access_token":"SOME_ACCESS_TOKEN","token_type":"Bearer","expires_in":3600}`))
	}))
	defer s.Close()

	c := &Client{
		AuthURL: s.URL,
	}

	// Act
	err := c.GetToken(context.TODO(), "TROGDOR", "BURNINATOR")

	// Assert
	if err != nil {
		t.Fatalf("unable to get token: %v", err)
	}
	if c.access_token != "SOME_ACCESS_TOKEN" {
		t.Errorf("got: %s want: SOME_ACCESS_TOKEN", c.access_token)
	}
	if gotAuth != "Basic "+wantAuth {
		t.Errorf("got auth %q, want %q", gotAuth, "Basic "+wantAuth)
	}
}
