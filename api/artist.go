package api

import "github.com/twoscott/gobble-fm/lastfm"

type Artist struct {
	api *API
}

// NewArtist creates and returns a new Artist API route.
func NewArtist(api *API) *Artist {
	return &Artist{api: api}
}

// Correction returns the artist name corrections of an artist.
func (a Artist) Correction(artist string) (*lastfm.ArtistCorrection, error) {
	var res lastfm.ArtistCorrection
	p := lastfm.ArtistCorrectionParams{Artist: artist}
	return &res, a.api.Get(&res, ArtistGetCorrectionMethod, p)
}

// Info returns the information of an artist by artist name.
func (a Artist) Info(params lastfm.ArtistInfoParams) (*lastfm.ArtistInfo, error) {
	var res lastfm.ArtistInfo
	return &res, a.api.Get(&res, ArtistGetInfoMethod, params)
}

// InfoByMBID returns the information of an artist by MBID.
func (a Artist) InfoByMBID(params lastfm.ArtistInfoMBIDParams) (*lastfm.ArtistInfo, error) {
	var res lastfm.ArtistInfo
	return &res, a.api.Get(&res, ArtistGetInfoMethod, params)
}

// Similar returns the similar artists of an artist by artist name.
func (a Artist) Similar(params lastfm.ArtistSimilarParams) (*lastfm.SimilarArtists, error) {
	var res lastfm.SimilarArtists
	return &res, a.api.Get(&res, ArtistGetSimilarMethod, params)
}

// SimilarByMBID returns the similar artists of an artist by MBID.
func (a Artist) SimilarByMBID(
	params lastfm.ArtistSimilarMBIDParams) (*lastfm.SimilarArtists, error) {

	var res lastfm.SimilarArtists
	return &res, a.api.Get(&res, ArtistGetSimilarMethod, params)
}

// UserTags returns the tags of an artist for user by artist name.
func (a Artist) UserTags(params lastfm.ArtistTagsParams) (*lastfm.ArtistTags, error) {
	var res lastfm.ArtistTags
	return &res, a.api.Get(&res, ArtistGetTagsMethod, params)
}

// UserTagsByMBID returns the tags of an artist for user by MBID.
func (a Artist) UserTagsByMBID(params lastfm.ArtistTagsMBIDParams) (*lastfm.ArtistTags, error) {
	var res lastfm.ArtistTags
	return &res, a.api.Get(&res, ArtistGetTagsMethod, params)
}

// TopAlbums returns the top albums of an artist by artist name.
func (a Artist) TopAlbums(params lastfm.ArtistTopAlbumsParams) (*lastfm.ArtistTopAlbums, error) {
	var res lastfm.ArtistTopAlbums
	return &res, a.api.Get(&res, ArtistGetTopAlbumsMethod, params)
}

// TopAlbumsByMBID returns the top albums of an artist by MBID.
func (a Artist) TopAlbumsByMBID(
	params lastfm.ArtistTopAlbumsMBIDParams) (*lastfm.ArtistTopAlbums, error) {

	var res lastfm.ArtistTopAlbums
	return &res, a.api.Get(&res, ArtistGetTopAlbumsMethod, params)
}

// TopTracks returns the top tracks of an artist by artist name.
func (a Artist) TopTags(params lastfm.ArtistTopTagsParams) (*lastfm.ArtistTopTags, error) {
	var res lastfm.ArtistTopTags
	return &res, a.api.Get(&res, ArtistGetTopTagsMethod, params)
}

// TopTagsByMBID returns the top tracks of an artist by MBID.
func (a Artist) TopTagsByMBID(
	params lastfm.ArtistTopTagsMBIDParams) (*lastfm.ArtistTopTags, error) {

	var res lastfm.ArtistTopTags
	return &res, a.api.Get(&res, ArtistGetTopTagsMethod, params)
}

// TopTracks returns the top tracks of an artist by artist name.
func (a Artist) TopTracks(params lastfm.ArtistTopTracksParams) (*lastfm.ArtistTopTracks, error) {
	var res lastfm.ArtistTopTracks
	return &res, a.api.Get(&res, ArtistGetTopTracksMethod, params)
}

// TopTracksByMBID returns the top tracks of an artist by MBID.
func (a Artist) TopTracksByMBID(
	params lastfm.ArtistTopTracksMBIDParams) (*lastfm.ArtistTopTracks, error) {

	var res lastfm.ArtistTopTracks
	return &res, a.api.Get(&res, ArtistGetTopTracksMethod, params)
}

// Search returns the results of an album search.
func (a Artist) Search(params lastfm.ArtistSearchParams) (*lastfm.ArtistSearchResult, error) {
	var res lastfm.ArtistSearchResult
	return &res, a.api.Get(&res, ArtistSearchMethod, params)
}
