package client

import (
	"github.com/twoscott/gobble-fm/v2/api"
	"github.com/twoscott/gobble-fm/v2/lastfm"
)

type Artist struct {
	api *api.API
}

// NewArtist creates and returns a new Artist API route.
func NewArtist(api *api.API) *Artist {
	return &Artist{api: api}
}

// Correction returns the artist name corrections of an artist.
func (a Artist) Correction(artist string) (*lastfm.ArtistCorrection, error) {
	var res lastfm.ArtistCorrection
	p := lastfm.ArtistCorrectionParams{Artist: artist}
	return &res, a.api.Get(&res, api.ArtistGetCorrectionMethod, p)
}

// Info returns the information of an artist by artist name.
func (a Artist) Info(params lastfm.ArtistInfoParams) (*lastfm.ArtistInfo, error) {
	var res lastfm.ArtistInfo
	return &res, a.api.Get(&res, api.ArtistGetInfoMethod, params)
}

// InfoByMBID returns the information of an artist by MBID.
func (a Artist) InfoByMBID(params lastfm.ArtistInfoMBIDParams) (*lastfm.ArtistInfo, error) {
	var res lastfm.ArtistInfo
	return &res, a.api.Get(&res, api.ArtistGetInfoMethod, params)
}

// UserInfo returns the information of an artist for user by artist name.
func (a Artist) UserInfo(params lastfm.ArtistUserInfoParams) (*lastfm.ArtistUserInfo, error) {
	var res lastfm.ArtistUserInfo
	return &res, a.api.Get(&res, api.ArtistGetInfoMethod, params)
}

// UserInfoByMBID returns the information of an artist for user by MBID.
func (a Artist) UserInfoByMBID(
	params lastfm.ArtistUserInfoMBIDParams) (*lastfm.ArtistUserInfo, error) {

	var res lastfm.ArtistUserInfo
	return &res, a.api.Get(&res, api.ArtistGetInfoMethod, params)
}

// Similar returns the similar artists of an artist by artist name.
func (a Artist) Similar(params lastfm.ArtistSimilarParams) (*lastfm.SimilarArtists, error) {
	var res lastfm.SimilarArtists
	return &res, a.api.Get(&res, api.ArtistGetSimilarMethod, params)
}

// SimilarByMBID returns the similar artists of an artist by MBID.
func (a Artist) SimilarByMBID(
	params lastfm.ArtistSimilarMBIDParams) (*lastfm.SimilarArtists, error) {

	var res lastfm.SimilarArtists
	return &res, a.api.Get(&res, api.ArtistGetSimilarMethod, params)
}

// UserTags returns the tags of an artist for user by artist name.
func (a Artist) UserTags(params lastfm.ArtistTagsParams) (*lastfm.ArtistTags, error) {
	var res lastfm.ArtistTags
	return &res, a.api.Get(&res, api.ArtistGetTagsMethod, params)
}

// UserTagsByMBID returns the tags of an artist for user by MBID.
func (a Artist) UserTagsByMBID(params lastfm.ArtistTagsMBIDParams) (*lastfm.ArtistTags, error) {
	var res lastfm.ArtistTags
	return &res, a.api.Get(&res, api.ArtistGetTagsMethod, params)
}

// TopAlbums returns the top albums of an artist by artist name.
func (a Artist) TopAlbums(params lastfm.ArtistTopAlbumsParams) (*lastfm.ArtistTopAlbums, error) {
	var res lastfm.ArtistTopAlbums
	return &res, a.api.Get(&res, api.ArtistGetTopAlbumsMethod, params)
}

// TopAlbumsByMBID returns the top albums of an artist by MBID.
func (a Artist) TopAlbumsByMBID(
	params lastfm.ArtistTopAlbumsMBIDParams) (*lastfm.ArtistTopAlbums, error) {

	var res lastfm.ArtistTopAlbums
	return &res, a.api.Get(&res, api.ArtistGetTopAlbumsMethod, params)
}

// TopTags returns the top tags of an artist by artist name.
func (a Artist) TopTags(params lastfm.ArtistTopTagsParams) (*lastfm.ArtistTopTags, error) {
	var res lastfm.ArtistTopTags
	return &res, a.api.Get(&res, api.ArtistGetTopTagsMethod, params)
}

// TopTagsByMBID returns the top tags of an artist by MBID.
func (a Artist) TopTagsByMBID(
	params lastfm.ArtistTopTagsMBIDParams) (*lastfm.ArtistTopTags, error) {

	var res lastfm.ArtistTopTags
	return &res, a.api.Get(&res, api.ArtistGetTopTagsMethod, params)
}

// TopTracks returns the top tracks of an artist by artist name.
func (a Artist) TopTracks(params lastfm.ArtistTopTracksParams) (*lastfm.ArtistTopTracks, error) {
	var res lastfm.ArtistTopTracks
	return &res, a.api.Get(&res, api.ArtistGetTopTracksMethod, params)
}

// TopTracksByMBID returns the top tracks of an artist by MBID.
func (a Artist) TopTracksByMBID(
	params lastfm.ArtistTopTracksMBIDParams) (*lastfm.ArtistTopTracks, error) {

	var res lastfm.ArtistTopTracks
	return &res, a.api.Get(&res, api.ArtistGetTopTracksMethod, params)
}

// Search returns the results of an artist search.
func (a Artist) Search(params lastfm.ArtistSearchParams) (*lastfm.ArtistSearchResult, error) {
	var res lastfm.ArtistSearchResult
	return &res, a.api.Get(&res, api.ArtistSearchMethod, params)
}

type ArtistSession struct {
	*Artist
	session *api.Session
}

// NewArtistSession creates and returns a new Artist API route.
func NewArtistSession(session *api.Session) *ArtistSession {
	return &ArtistSession{Artist: NewArtist(session.API), session: session}
}

// AddTags adds tags to an artist for the authenticated user.
func (a ArtistSession) AddTags(artist string, tags []string) error {
	p := lastfm.ArtistAddTagsParams{Artist: artist, Tags: tags}
	return a.session.Post(nil, api.ArtistAddTagsMethod, p)
}

// SelfTags returns the tags of an artist for the authenticated user by artist
// name.
func (a ArtistSession) SelfTags(params lastfm.ArtistSelfTagsParams) (*lastfm.ArtistTags, error) {
	var res lastfm.ArtistTags
	return &res, a.session.Get(&res, api.ArtistGetTagsMethod, params)
}

// SelfTagsByMBID returns the tags of an artist for the authenticated user by
// MBID.
func (a ArtistSession) SelfTagsByMBID(params lastfm.ArtistSelfTagsMBIDParams) (*lastfm.ArtistTags, error) {
	var res lastfm.ArtistTags
	return &res, a.session.Get(&res, api.ArtistGetTagsMethod, params)
}

// RemoveTag removes a tag from an artist for the authenticated user.
func (a ArtistSession) RemoveTag(artist string, tag string) error {
	p := lastfm.ArtistRemoveTagParams{Artist: artist, Tag: tag}
	return a.session.Post(nil, api.ArtistRemoveTagMethod, p)
}
