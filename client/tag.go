package client

import (
	"github.com/twoscott/gobble-fm/v2/api"
	"github.com/twoscott/gobble-fm/v2/lastfm"
)

type Tag struct {
	api *api.API
}

// NewTag creates and returns a new Tag API route.
func NewTag(api *api.API) *Tag {
	return &Tag{api: api}
}

// Info returns the information of a tag.
func (t Tag) Info(params lastfm.TagInfoParams) (*lastfm.TagInfo, error) {
	var res lastfm.TagInfo
	return &res, t.api.Get(&res, api.TagGetInfoMethod, params)
}

// Similar returns tags similar to the given tag.
func (t Tag) Similar(tag string) (*lastfm.SimilarTags, error) {
	var res lastfm.SimilarTags
	p := lastfm.TagSimilarParams{Tag: tag}
	return &res, t.api.Get(&res, api.TagGetSimilarMethod, p)
}

// TopAlbums returns the top albums tagged with the given tag.
func (t Tag) TopAlbums(params lastfm.TagTopAlbumsParams) (*lastfm.TagTopAlbums, error) {
	var res lastfm.TagTopAlbums
	return &res, t.api.Get(&res, api.TagGetTopAlbumsMethod, params)
}

// TopArtists returns the top artists tagged with the given tag.
func (t Tag) TopArtists(params lastfm.TagTopArtistsParams) (*lastfm.TagTopArtists, error) {
	var res lastfm.TagTopArtists
	return &res, t.api.Get(&res, api.TagGetTopArtistsMethod, params)
}

// TopTags returns the top tags on Last.fm.
func (t Tag) TopTags() (*lastfm.TagTopTags, error) {
	return t.TopTagsLimit(nil)
}

// TopTagsLimit returns the top tags on Last.fm, with optional limit and offset.
func (t Tag) TopTagsLimit(params *lastfm.TagTopTagsParams) (*lastfm.TagTopTags, error) {
	var res lastfm.TagTopTags
	return &res, t.api.Get(&res, api.TagGetTopTagsMethod, params)
}

// TopTracks returns the top tracks tagged with the given tag.
func (t Tag) TopTracks(params lastfm.TagTopTracksParams) (*lastfm.TagTopTracks, error) {
	var res lastfm.TagTopTracks
	return &res, t.api.Get(&res, api.TagGetTopTracksMethod, params)
}

// WeeklyChartList returns the weekly chart list of a tag.
func (t Tag) WeeklyChartList(tag string) (*lastfm.TagWeeklyChartList, error) {
	var res lastfm.TagWeeklyChartList
	p := lastfm.TagWeeklyChartListParams{Tag: tag}
	return &res, t.api.Get(&res, api.TagGetWeeklyChartListMethod, p)
}
