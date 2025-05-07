package api

import "github.com/twoscott/go-fm/lastfm"

type Album struct {
	api *API
}

// NewAlbum creates and returns a new Album API route.
func NewAlbum(api *API) *Album {
	return &Album{api: api}
}

// Info returns the information of an album.
func (a Album) Info(params lastfm.AlbumInfoParams) (*lastfm.AlbumInfo, error) {
	var res lastfm.AlbumInfo
	return &res, a.api.Get(&res, AlbumGetInfoMethod, params)
}
