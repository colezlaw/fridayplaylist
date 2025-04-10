package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	fridayplaylist "github.com/colezlaw/fridayPlaylist"
)

type PlaylistClient interface {
	GetPlaylistsForUser(string) ([]fridayplaylist.Playlist, error)
	GetTracksForPlaylist(string) ([]fridayplaylist.Track, error)
}

func run(args []string, stdout io.Writer, c PlaylistClient) error {
	log.SetOutput(stdout)
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	var (
		fn   = flags.String("f", "output.csv", "output filename")
		user = flags.String("user", "jshields14", "username for playlists to collect")
	)
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	playlists, err := c.GetPlaylistsForUser(*user)
	if err != nil {
		return fmt.Errorf("getplaylistsforuser: %w", err)
	}

	of, err := os.Create(*fn)
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}
	defer of.Close()

	w := csv.NewWriter(of)
	defer w.Flush()

	w.Write([]string{"PLAYLIST", "SONG", "ARTIST", "RELEASED"})

	for _, playlist := range playlists {
		fmt.Printf("%#v\n", playlist)
		fmt.Println(playlist.Name)
		tracks, err := c.GetTracksForPlaylist(playlist.ID)
		if err != nil {
			return fmt.Errorf("gettracks: %w", err)
		}
		for _, track := range tracks {
			w.Write([]string{playlist.Name, track.Name, track.Artist, track.Album.ReleaseDate})
		}
	}

	return nil
}

func main() {
	c := &fridayplaylist.Client{}
	if err := c.GetToken(os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET")); err != nil {
		fmt.Fprintf(os.Stderr, "unable to get token: %v", err)
	}
	if err := run(os.Args, os.Stdout, c); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
