package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	fridayplaylist "github.com/colezlaw/fridayPlaylist"
)

func main() {
	c := fridayplaylist.Client{}
	if err := c.GetToken(os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET")); err != nil {
		log.Fatalf("gettoken: %v", err)
	}

	playlists, err := c.GetPlaylistsForUser("jshields14")
	if err != nil {
		log.Fatalf("getplaylistsforuser: %v", err)
	}

	of, err := os.Create("output.csv")
	if err != nil {
		log.Fatalf("create: %v", err)
	}
	defer of.Close()

	w := csv.NewWriter(of)
	defer w.Flush()

	w.Write([]string{"PLAYLIST", "SONG", "ARTIST", "RELEASED"})

	for _, playlist := range playlists {
		fmt.Println(playlist.Name)
		tracks, err := c.GetTracksForPlaylist(playlist.ID)
		if err != nil {
			log.Fatalf("gettracks: %v", err)
		}
		for _, track := range tracks {
			w.Write([]string{playlist.Name, track.Name, track.Artist, track.Album.ReleaseDate})
		}
	}
}
