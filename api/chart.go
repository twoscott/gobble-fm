package api

import "github.com/twoscott/gobble-fm/lastfm"

type Chart struct {
	api *API
}

// NewChart creates and returns a new Chart API route.
func NewChart(api *API) *Chart {
	return &Chart{api: api}
}

// TopArtistsLimit returns the top artists of the chart.
func (c Chart) TopArtistsLimit(params *lastfm.ChartTopArtistsParams) (*lastfm.ChartTopArtists, error) {
	var res lastfm.ChartTopArtists
	return &res, c.api.Get(&res, ChartGetTopArtistsMethod, params)
}

// TopArtists returns all the top artists of the chart. Same as
// TopArtistsLimit(nil).
func (c Chart) TopArtists() (*lastfm.ChartTopArtists, error) {
	return c.TopArtistsLimit(nil)
}

// TopTagsLimit returns the top tags of the chart.
func (c Chart) TopTagsLimit(params *lastfm.ChartTopTagsParams) (*lastfm.ChartTopTags, error) {
	var res lastfm.ChartTopTags
	return &res, c.api.Get(&res, ChartGetTopTagsMethod, params)
}

// TopTags returns the top tags of the chart. Same as TopTagsLimit(nil).
func (c Chart) TopTags() (*lastfm.ChartTopTags, error) {
	return c.TopTagsLimit(nil)
}

// TopTracksLimit returns the top tracks of the chart.
func (c Chart) TopTracksLimit(params *lastfm.ChartTopTracksParams) (*lastfm.ChartTopTracks, error) {
	var res lastfm.ChartTopTracks
	return &res, c.api.Get(&res, ChartGetTopTracksMethod, params)
}

// TopTracks returns all the top tracks of the chart. Same as
// TopTracksLimit(nil).
func (c Chart) TopTracks() (*lastfm.ChartTopTracks, error) {
	return c.TopTracksLimit(nil)
}
