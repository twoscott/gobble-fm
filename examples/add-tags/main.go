package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/twoscott/gobble-fm/api"
	"github.com/twoscott/gobble-fm/lastfm"
	"github.com/twoscott/gobble-fm/session"
	"github.com/twoscott/gobble-fm/util/option"
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

	tags := []string{"electronic", "idm", "ambient", "experimental", "techno"}

	err = fm.Artist.AddTags("Aphex Twin", tags)
	if err != nil {
		fmt.Println("Error adding tags:", err)
		return
	}

	fmt.Println("Tags added successfully.")

	res, err := fm.Artist.SelfTags(lastfm.ArtistSelfTagsParams{
		Artist:      "Aphex Twin",
		AutoCorrect: option.False,
	})
	if err != nil {
		fmt.Println("Error getting tags:", err)
		return
	}

	var tagNames []string
	for _, tag := range res.Tags {
		tagNames = append(tagNames, tag.Name)
	}

	fmt.Println("Tags:", strings.Join(tagNames, ", "))
}
