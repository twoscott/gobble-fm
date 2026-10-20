package client

import (
	"github.com/twoscott/gobble-fm/v2/api"
	"github.com/twoscott/gobble-fm/v2/lastfm"
)

type Geo struct {
	api *api.API
}

// NewGeo creates and returns a new Geo API route.
func NewGeo(api *api.API) *Geo {
	return &Geo{api: api}
}

// TopArtists returns the top artists of a country.
func (g Geo) TopArtists(params lastfm.GeoTopArtistsParams) (*lastfm.GeoTopArtists, error) {
	var res lastfm.GeoTopArtists
	return &res, g.api.Get(&res, api.GeoGetTopArtistsMethod, params)
}

// TopTracks returns the top tracks of a country.
func (g Geo) TopTracks(params lastfm.GeoTopTracksParams) (*lastfm.GeoTopTracks, error) {
	var res lastfm.GeoTopTracks
	return &res, g.api.Get(&res, api.GeoGetTopTracksMethod, params)
}
