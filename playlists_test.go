package fridayplaylist

import (
	"os"
	"testing"
)

func TestGetPlaylistForUser(t *testing.T) {
	c := &Client{}
	if err := c.GetToken(os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET")); err != nil {
		t.Fatalf("%v", err)
	}

	result, err := c.GetPlaylistsForUser("jshields14")
	if err != nil {
		t.Fatalf("unable to get playlists: %v", err)
	}
	for _, p := range result {
		t.Logf("%+v", p)
	}
	t.Logf("%d", len(result))
}
