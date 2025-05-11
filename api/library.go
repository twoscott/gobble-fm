package api

import "github.com/twoscott/gobble-fm/lastfm"

type Library struct {
	api *API
}

// NewLibrary creates and returns a new Library API route.
func NewLibrary(api *API) *Library {
	return &Library{api: api}
}

func (l *Library) Artists(params lastfm.LibraryArtistsParams) (*lastfm.LibraryArtists, error) {
	var res lastfm.LibraryArtists
	return &res, l.api.Get(&res, LibraryGetArtistsMethod, params)
}
