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

	// Now you can use the session key to make requests on behalf of the user.

	// For example, to scrobble a track:
	res, err := fm.Track.Scrobble(lastfm.ScrobbleParams{
		Track:    "L'Amour Toujours",
		Artist:   "Gigi D'Agostino",
		Album:    "L'Amour Toujours",
		Duration: lastfm.DurationMinSec(4, 02),
		Time:     time.Now(),
	})
	if err != nil {
		fmt.Println("Error scrobbling track:", err)
		return
	}

	s := res.Scrobble
	if s.Ignored.Code == lastfm.ScrobbleNotIgnored {
		fmt.Println("Scrobbled track:", s.Track.Title, "by", s.Artist.Name)
	} else {
		fmt.Println("Track was ignored:", s.Ignored.Message())
	}
}
