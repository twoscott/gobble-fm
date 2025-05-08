package api

import "github.com/twoscott/gobble-fm/lastfm"

type Chart struct {
	api *API
}

// NewChart creates and returns a new Chart API route.
func NewChart(api *API) *Chart {
	return &Chart{api: api}
}

// TopArtists returns the top artists of the chart.
func (c Chart) TopArtists(params *lastfm.ChartTopArtistsParams) (*lastfm.ChartTopArtists, error) {
	var res lastfm.ChartTopArtists
	return &res, c.api.Get(&res, ChartGetTopArtistsMethod, params)
}

// AllTopArtists returns all the top artists of the chart.
// Same as TopArtists(nil).
func (c Chart) AllTopArtists() (*lastfm.ChartTopArtists, error) {
	return c.TopArtists(nil)
}

// TopTags returns the top tags of the chart.
func (c Chart) TopTags(params *lastfm.ChartTopTagsParams) (*lastfm.ChartTopTags, error) {
	var res lastfm.ChartTopTags
	return &res, c.api.Get(&res, ChartGetTopTagsMethod, params)
}

// AllTopTags returns all the top tags of the chart.
// Same as TopTags(nil).
func (c Chart) AllTopTags() (*lastfm.ChartTopTags, error) {
	return c.TopTags(nil)
}

// TopTracks returns the top tracks of the chart.
func (c Chart) TopTracks(params *lastfm.ChartTopTracksParams) (*lastfm.ChartTopTracks, error) {
	var res lastfm.ChartTopTracks
	return &res, c.api.Get(&res, ChartGetTopTracksMethod, params)
}

// AllTopTracks returns all the top tracks of the chart.
// Same as TopTracks(nil).
func (c Chart) AllTopTracks() (*lastfm.ChartTopTracks, error) {
	return c.TopTracks(nil)
}
