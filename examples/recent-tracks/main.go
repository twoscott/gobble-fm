package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/twoscott/gobble-fm/api"
	"github.com/twoscott/gobble-fm/lastfm"
)

func main() {
	apiKey := os.Getenv("LASTFM_API_KEY")

	fm := api.NewClientKeyOnly(apiKey)

	params := lastfm.RecentTracksParams{
		User:  "Username",
		Limit: 5,
		From:  time.Now().Add(-24 * time.Hour),
	}

	res, err := fm.User.RecentTracks(params)
	if err != nil {
		fmerr := &api.LastFMError{}
		if errors.As(err, &fmerr) {
			switch fmerr.Code {
			case api.ErrInvalidParameters:
				fmt.Println("Invalid parameters")
			case api.ErrOperationFailed:
				fmt.Println("Operation failed")
			default:
				fmt.Println(err)
				// ...
			}
		} else {
			fmt.Println(err)
		}

		return
	}

	for i, t := range res.Tracks {
		fmt.Printf("%d.\t%s by %s\n", i+1, t.Title, t.Artist.Name)

		if t.NowPlaying {
			fmt.Println("\tNow playing...")
		} else {
			ago := time.Since(t.ScrobbledAt.Time()).Truncate(time.Second)
			fmt.Printf("\tScrobbled %s ago\n", ago)
		}

		fmt.Printf("\n\tArt: %s\n", t.Image.OriginalURL())
		fmt.Println()
	}
}
