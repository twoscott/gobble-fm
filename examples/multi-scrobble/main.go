package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/twoscott/gobble-fm/api"
	"github.com/twoscott/gobble-fm/lastfm"
	"github.com/twoscott/gobble-fm/session"
)

func main() {
	apiKey := os.Getenv("LASTFM_API_KEY")
	secret := os.Getenv("LASTFM_API_SECRET")

	fm := session.NewClientWithTimeout(apiKey, secret, 5)

	username := os.Getenv("LASTFM_USERNAME")
	password := os.Getenv("LASTFM_PASSWORD")

	err := fm.Login(username, password)
	if err != nil {
		fmerr := &api.LastFMError{}
		if errors.As(err, &fmerr) {
			switch fmerr.Code {
			case api.ErrAuthenticationFailed:
				fmt.Println("Login failed. Check your username and password.")
			default:
				fmt.Println("Login failed:", err)
				// ...
			}
		} else {
			fmt.Println(err)
		}

		return
	}

	scrobbles := lastfm.ScrobbleMultiParams{
		{
			Track:    "Xtal",
			Artist:   "Aphex Twin",
			Album:    "Selected Ambient Works 85-92",
			Duration: lastfm.DurationMinSec(4, 54),
			Time:     time.Now(),
		},
		{
			Track:    "L'Amour Toujours",
			Artist:   "Gigi D'Agostino",
			Album:    "L'Amour Toujours",
			Duration: lastfm.DurationMinSec(4, 02),
			Time:     time.Now(),
		},
		{
			Track:       "Move Me",
			Artist:      "Kohta Takahashi",
			Album:       "R4 / RIDGE RACER TYPE 4 / DIRECT AUDIO",
			AlbumArtist: "Various Artists",
			Duration:    lastfm.DurationMinSec(4, 33),
			Time:        time.Now(),
		},
	}

	res, err := fm.Track.ScrobbleMulti(scrobbles)
	if err != nil {
		fmt.Println("Error scrobbling tracks:", err)
		return
	}

	for _, s := range res.Scrobbles {
		if s.Ignored.Code == lastfm.ScrobbleNotIgnored {
			fmt.Printf("Scrobbled track: %s by %s\n", s.Track.Title, s.Artist.Name)
		} else {
			fmt.Printf("Ignored track: %s by %s - ", s.Track.Title, s.Artist.Name)
			fmt.Println("Reason:", s.Ignored.Message())
		}
	}
}
