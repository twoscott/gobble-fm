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

	fm := session.NewClientWithTimeout(apiKey, secret, 30)

	// If you set up a callback URL when applying for your Last.fm API key,
	// then you can use `Session.AuthURL()` instead, and it will redirect the
	// user to the callback URL after they authorize your app.
	url := fm.AuthCallbackURL("https://example.com/callback")

	// Redirect the user to the auth URL. In a real application, you would
	// redirect the user to this URL in their web browser.
	redirectUserToAuthURL(url)

	// After authorising your app, the user will be redirected to the callback
	// URL with a `token` query parameter. You need to extract this token and
	// use it to create a session key.
	// For example: https://example.com/callback?token=TOKEN

	token := "AUTHORISED_TOKEN" // obtained from the callback URL

	// After the user authorizes your app, you can use the token to log in.
	err := fm.TokenLogin(token)
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

	// Now you can use the session key to make requests on behalf of the user.

	// For example, to scrobble a track:
	res, err := fm.Track.Scrobble(lastfm.ScrobbleParams{
		Track:       "Move Me",
		Artist:      "Kohta Takahashi",
		Album:       "R4 / RIDGE RACER TYPE 4 / DIRECT AUDIO",
		AlbumArtist: "Various Artists",
		Duration:    lastfm.DurationMinSec(4, 33),
		Time:        time.Now(),
	})
	if err != nil {
		fmt.Println("Error scrobbling track:", err)
		return
	}

	s := res.Scrobble
	fmt.Println("Scrobbled track:", s.Track.Title, "by", s.Artist.Name)
}

func redirectUserToAuthURL(url string) {
	// This is a placeholder function. In a real application, you would
	// redirect the user to the auth URL.
	fmt.Println("Redirecting user to:", url)
}
