package api

import "github.com/twoscott/gobble-fm/lastfm"

type Geo struct {
	api *API
}

// NewGeo creates and returns a new Geo API route.
func NewGeo(api *API) *Geo {
	return &Geo{api: api}
}

// TopArtists returns the top artists of a country.
func (g Geo) TopArtists(params lastfm.GeoTopArtistsParams) (*lastfm.GeoTopArtists, error) {
	var res lastfm.GeoTopArtists
	return &res, g.api.Get(&res, GeoGetTopArtistsMethod, params)
}

// TopTracks returns the top tracks of a country.
func (g Geo) TopTracks(params lastfm.GeoTopTracksParams) (*lastfm.GeoTopTracks, error) {
	var res lastfm.GeoTopTracks
	return &res, g.api.Get(&res, GeoGetTopTracksMethod, params)
}
