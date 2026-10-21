package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/twoscott/gobble-fm/v2/api"
	"github.com/twoscott/gobble-fm/v2/client"
	"github.com/twoscott/gobble-fm/v2/lastfm"
)

func main() {
	apiKey := os.Getenv("LASTFM_API_KEY")
	secret := os.Getenv("LASTFM_API_SECRET")

	fm, err := client.NewWithTimeout(apiKey, secret, 30)
	if err != nil {
		fmt.Println("Error creating client:", err)
		return
	}

	// Fetch an unauthorised token from the API.
	token, err := fm.Auth.Token()
	if err != nil {
		fmt.Println("Error getting token:", err)
		return
	}

	// Now the token must be authorized by the user. This is done by redirecting
	// the user to the authorization URL, where they can log in and authorize
	// the application.
	url := fm.AuthTokenURL(token)

	// In a real application, you would redirect the user to this URL in their
	// web browser. For this example, we'll just print it to the console.
	fmt.Printf("Visit the following URL to authorize the application: %s\n", url)

	// The token can't be used until the user authorizes it.
	fmt.Println("Once you've authorized the application, press Enter... ")
	fmt.Scanln()

	// After the user authorizes your app, you can use the token to log in.
	session, err := fm.TokenLogin(token)
	if err != nil {
		var fmerr *api.LastFMError
		if errors.As(err, &fmerr) {
			switch fmerr.Code {
			case api.ErrUnauthorizedToken:
				fmt.Println("You must authorize the token before using it.")
			default:
				fmt.Println("Authorization failed:", err)
				// ...
			}
		} else {
			fmt.Println(err)
		}

		return
	}

	// Now you can use the session to make requests on behalf of the user.

	// For example, to scrobble a track:
	res, err := session.Track.Scrobble(lastfm.ScrobbleParams{
		Track:    "Xtal",
		Artist:   "Aphex Twin",
		Album:    "Selected Ambient Works 85-92",
		Duration: lastfm.DurationMinSec(4, 54),
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
