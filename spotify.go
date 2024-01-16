package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/tabwriter"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

const redirectURI = "http://localhost:8080/callback"

var (
	auth = spotifyauth.New(
		spotifyauth.WithRedirectURL(redirectURI),
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserReadPrivate,
			spotifyauth.ScopeUserReadCurrentlyPlaying,
			spotifyauth.ScopeUserReadPlaybackState,
			spotifyauth.ScopeUserModifyPlaybackState,
		),
	)
	ch    = make(chan *spotify.Client)
	state = "abc123"
)

func main() {
	ctx := context.Background()

	// Start HTTP server
	http.HandleFunc("/callback", completeAuth)
	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Log in user
	url := auth.AuthURL(state)
	fmt.Printf("Login:\n%s\n", url)
	client := <-ch

	user, err := client.CurrentUser(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Logged In As:", user.ID)

	// search for albums with the name Sempiternal
	results, err := client.Search(ctx, "Petrodragonic Apocolypse", spotify.SearchTypeAlbum)
	if err != nil {
		log.Fatal(err)
	}

	// select the top album
	item := results.Albums.Albums[0]

	// get tracks from album
	res, err := client.GetAlbumTracks(ctx, item.ID, spotify.Market("US"))

	if err != nil {
		log.Fatal("error getting tracks ....", err.Error())
		return
	}

	// *display in tabular form using TabWriter
	w := tabwriter.NewWriter(os.Stdout, 10, 2, 5, ' ', 0)
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t\n\n", "Songs", "Energy", "Danceability", "Valence")

	// loop through tracks
	for _, track := range res.Tracks {

		// retrieve features
		features, err := client.GetAudioFeatures(ctx, track.ID)
		if err != nil {
			log.Fatal("error getting audio features...", err.Error())
			return
		}
		fmt.Fprintf(w, "%s\t%v\t%v\t%v\t\n", track.Name, features[0].Energy, features[0].Danceability, features[0].Valence)
		w.Flush()
	}
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}

	client := spotify.New(auth.Client(r.Context(), tok))
	fmt.Fprintf(w, "Login Completed!")
	ch <- client
}
