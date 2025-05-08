package api

import "github.com/twoscott/go-fm/lastfm"

type Album struct {
	api *API
}

// NewAlbum creates and returns a new Album API route.
func NewAlbum(api *API) *Album {
	return &Album{api: api}
}

// Info returns the information of an album by artist and album name.
func (a Album) Info(params lastfm.AlbumInfoParams) (*lastfm.AlbumInfo, error) {
	var res lastfm.AlbumInfo
	return &res, a.api.Get(&res, AlbumGetInfoMethod, params)
}

// InfoByMBID returns the information of an album by MBID.
func (a Album) InfoByMBID(params lastfm.AlbumInfoMBIDParams) (*lastfm.AlbumInfo, error) {
	var res lastfm.AlbumInfo
	return &res, a.api.Get(&res, AlbumGetInfoMethod, params)
}

// UserTags returns the tags of an album for user by artist and album name.
func (a Album) UserTags(params lastfm.AlbumTagsParams) (*lastfm.AlbumTags, error) {
	var res lastfm.AlbumTags
	return &res, a.api.Get(&res, AlbumGetTagsMethod, params)
}

// UserTagsByMBID returns the tags of an album for user by MBID.
func (a Album) UserTagsByMBID(params lastfm.AlbumTagsMBIDParams) (*lastfm.AlbumTags, error) {
	var res lastfm.AlbumTags
	return &res, a.api.Get(&res, AlbumGetTagsMethod, params)
}

// TopTags returns the top tags of an album by artist and album name.
func (a Album) TopTags(params lastfm.AlbumTopTagsParams) (*lastfm.AlbumTopTags, error) {
	var res lastfm.AlbumTopTags
	return &res, a.api.Get(&res, AlbumGetTopTagsMethod, params)
}

// TopTagsByMBID returns the top tags of an album by MBID.
//
// Deprecated: Fetching top tags by MBID doesn't seem to work. Use TopTags
// instead.
func (a Album) TopTagsByMBID(params lastfm.AlbumTopTagsMBIDParams) (*lastfm.AlbumTopTags, error) {
	var res lastfm.AlbumTopTags
	return &res, a.api.Get(&res, AlbumGetTopTagsMethod, params)
}

// Search returns the results of an album search.
func (a Album) Search(params lastfm.AlbumSearchParams) (*lastfm.AlbumSearchResult, error) {
	var res lastfm.AlbumSearchResult
	return &res, a.api.Get(&res, AlbumSearchMethod, params)
}
