package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/twoscott/gobble-fm/api"
	"github.com/twoscott/gobble-fm/lastfm"
)

func main() {
	apiKey := os.Getenv("LASTFM_API_KEY")

	fm := api.NewClientKeyOnly(apiKey)

	params := lastfm.UserTopAlbumsParams{
		User:   "Username",
		Limit:  5,
		Period: lastfm.PeriodMonth,
	}

	res, err := fm.User.TopAlbums(params)
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

	for _, t := range res.Albums {
		fmt.Printf("%d.\t%s by %s\n", t.Rank, t.Title, t.Artist.Name)
		fmt.Printf("\tPlaycount: %d\n", t.Playcount)
		fmt.Printf("\tURL: %s\n", t.URL)
		fmt.Printf("\tCover: %s\n", t.Cover.OriginalURL())
	}
}
