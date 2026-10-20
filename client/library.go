package client

import (
	"github.com/twoscott/gobble-fm/v2/api"
	"github.com/twoscott/gobble-fm/v2/lastfm"
)

type Library struct {
	api *api.API
}

// NewLibrary creates and returns a new Library API route.
func NewLibrary(api *api.API) *Library {
	return &Library{api: api}
}

func (l Library) Artists(params lastfm.LibraryArtistsParams) (*lastfm.LibraryArtists, error) {
	var res lastfm.LibraryArtists
	return &res, l.api.Get(&res, api.LibraryGetArtistsMethod, params)
}
