# Gobble.fm

[![Go Reference](https://img.shields.io/badge/reference-009bc2?style=flat-round&logo=go&logoColor=ffffff)](https://pkg.go.dev/github.com/twoscott/gobble-fm)
[![Go Version 1.23+](https://img.shields.io/badge/go-1.23+-009bc2?style=flat-round)](https://golang.org/dl/)
[![Tag](https://img.shields.io/github/v/tag/twoscott/gobble-fm?style=flat-round&color=00b1b1)](https://github.com/twoscott/gobble-fm/tags)
[![Go Tests](https://img.shields.io/github/actions/workflow/status/twoscott/gobble-fm/test.yml?branch=master&style=flat-round&label=tests)](https://github.com/twoscott/gobble-fm/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/twoscott/gobble-fm?style=flat-round)](https://goreportcard.com/report/github.com/twoscott/gobble-fm)
[![Last Commit](https://img.shields.io/github/last-commit/twoscott/gobble-fm?logo=github&logoColor=ffffff&style=flat-round)](https://github.com/twoscott/gobble-fm/commits/master)

Gobble.fm is a Go (Golang) library for interacting with the Last.fm API.

## Why Gobble.fm?

- Comprehensive API coverage.
- Package separation between unauthenticated and authenticated API methods.
- Typed API parameter structs for URL encoding—no need to reference API docs or manually enter parameter names.
- Typed response struct fields—no need to convert from strings.
- Helper types and constants for easier API interaction.

## Installation

	go get github.com/twoscott/gobble-fm

## Documentation

- [Gobble.fm documentation](https://pkg.go.dev/github.com/twoscott/gobble-fm)
- [Last.fm API documentation](https://www.last.fm/api)

## Usage

First you need to instatiate the Last.fm API. You can choose the level of abstraction you'd like to use to interact with the API:
```go
import "github.com/twoscott/gobble-fm/api"
// Basic API client with only the API key. No access to auth methods.
fm := api.NewClientKeyOnly("API_KEY")
```
```go
// Make calls to auth.[getMobileSession|getSession|getToken] methods.
fm := api.NewClient("API_KEY", "SECRET")
```
```go
import "github.com/twoscott/gobble-fm/session"
// Authenticate API calls on behalf of a user.
fm := session.NewClient("API_KEY", "SECRET")
// Must authenticate a user first. e.g.,
fm.Login("USERNAME", "PASSWORD")
// or
fm.TokenLogin("AUTHORIZED_TOKEN")
```
#
Low-level abstractions:
```go
import "github.com/twoscott/gobble-fm/api"
// Provides methods for making API requests such as Get, Post, and Request.
fm := api.New("API_KEY", "SECRET")
```
```go
import "github.com/twoscott/gobble-fm/session"
// Provides methods for making authenticated API requests.
fm := session.New("API_KEY", "SECRET")
// Must authenticate a user first. e.g.,
// Obtain session key from one of the auth methods.
fm.SetSessionKey("SESSION_KEY")
```

## Simple Example
```go
package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/twoscott/gobble-fm/api"
	"github.com/twoscott/gobble-fm/lastfm"
)

func main() {
	fm := api.NewClientKeyOnly("API_KEY")

	params := lastfm.RecentTracksParams{
		User:  "Username",
		Limit: 5,
		From:  time.Now().Add(-24 * time.Hour),
	}

	res, err := fm.User.RecentTracks(params)
	if err != nil {
		var fmerr *api.LastFMError
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
```

## More Examples

- #### [Mobile Auth Example](https://github.com/twoscott/gobble-fm/blob/master/examples/auth/auth-flow-mobile/main.go)
- #### [Desktop Auth Example](https://github.com/twoscott/gobble-fm/blob/master/examples/auth/auth-flow-desktop/main.go)
- #### [Web Auth Example](https://github.com/twoscott/gobble-fm/blob/master/examples/auth/auth-flow-web/main.go)
- #### [Multi Scrobble Example](https://github.com/twoscott/gobble-fm/blob/master/examples/multi-scrobble/main.go)
- #### [Recent Tracks Example](https://github.com/twoscott/gobble-fm/blob/master/examples/recent-tracks/main.go)
- #### [Top Albums Example](https://github.com/twoscott/gobble-fm/blob/master/examples/top-albums/main.go)
- #### [Add Tags Example](https://github.com/twoscott/gobble-fm/blob/master/examples/add-tags/main.go)
