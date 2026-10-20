package client

import (
	"github.com/twoscott/gobble-fm/v2/api"
	"github.com/twoscott/gobble-fm/v2/lastfm"
)

type Album struct {
	api *api.API
}

// NewAlbum creates and returns a new Album API route.
func NewAlbum(api *api.API) *Album {
	return &Album{api: api}
}

// Info returns the information of an album by artist and album name.
func (a Album) Info(params lastfm.AlbumInfoParams) (*lastfm.AlbumInfo, error) {
	var res lastfm.AlbumInfo
	return &res, a.api.Get(&res, api.AlbumGetInfoMethod, params)
}

// InfoByMBID returns the information of an album by MBID.
func (a Album) InfoByMBID(params lastfm.AlbumInfoMBIDParams) (*lastfm.AlbumInfo, error) {
	var res lastfm.AlbumInfo
	return &res, a.api.Get(&res, api.AlbumGetInfoMethod, params)
}

// UserInfo returns the information of an album for user by artist and album
// name.
func (a Album) UserInfo(params lastfm.AlbumUserInfoParams) (*lastfm.AlbumUserInfo, error) {
	var res lastfm.AlbumUserInfo
	return &res, a.api.Get(&res, api.AlbumGetInfoMethod, params)
}

// UserInfoByMBID returns the information of an album for user by MBID.
func (a Album) UserInfoByMBID(
	params lastfm.AlbumUserInfoMBIDParams) (*lastfm.AlbumUserInfo, error) {

	var res lastfm.AlbumUserInfo
	return &res, a.api.Get(&res, api.AlbumGetInfoMethod, params)
}

// UserTags returns the tags of an album for user by artist and album name.
func (a Album) UserTags(params lastfm.AlbumTagsParams) (*lastfm.AlbumTags, error) {
	var res lastfm.AlbumTags
	return &res, a.api.Get(&res, api.AlbumGetTagsMethod, params)
}

// UserTagsByMBID returns the tags of an album for user by MBID.
func (a Album) UserTagsByMBID(params lastfm.AlbumTagsMBIDParams) (*lastfm.AlbumTags, error) {
	var res lastfm.AlbumTags
	return &res, a.api.Get(&res, api.AlbumGetTagsMethod, params)
}

// TopTags returns the top tags of an album by artist and album name.
func (a Album) TopTags(params lastfm.AlbumTopTagsParams) (*lastfm.AlbumTopTags, error) {
	var res lastfm.AlbumTopTags
	return &res, a.api.Get(&res, api.AlbumGetTopTagsMethod, params)
}

// TopTagsByMBID returns the top tags of an album by MBID.
//
// Deprecated: Fetching top tags by MBID doesn't seem to work. Use TopTags
// instead.
func (a Album) TopTagsByMBID(params lastfm.AlbumTopTagsMBIDParams) (*lastfm.AlbumTopTags, error) {
	var res lastfm.AlbumTopTags
	return &res, a.api.Get(&res, api.AlbumGetTopTagsMethod, params)
}

// Search returns the results of an album search.
func (a Album) Search(params lastfm.AlbumSearchParams) (*lastfm.AlbumSearchResult, error) {
	var res lastfm.AlbumSearchResult
	return &res, a.api.Get(&res, api.AlbumSearchMethod, params)
}

type AlbumSession struct {
	*Album
	session *api.Session
}

// NewAlbumSession creates and returns a new Album API route.
func NewAlbumSession(session *api.Session) *AlbumSession {
	return &AlbumSession{Album: NewAlbum(session.API), session: session}
}

// AddTags adds tags to an album for the authenticated user.
func (a AlbumSession) AddTags(artist, album string, tags []string) error {
	p := lastfm.AlbumAddTagsParams{Artist: artist, Album: album, Tags: tags}
	return a.session.Post(nil, api.AlbumAddTagsMethod, p)
}

// SelfTags returns the tags of an album for the authenticated user by
// artist and album name.
func (a AlbumSession) SelfTags(params lastfm.AlbumSelfTagsParams) (*lastfm.AlbumTags, error) {
	var res lastfm.AlbumTags
	return &res, a.session.Get(&res, api.AlbumGetTagsMethod, params)
}

// SelfTagsByMBID returns the tags of an album for the authenticated user by MBID.
func (a AlbumSession) SelfTagsByMBID(params lastfm.AlbumSelfTagsMBIDParams) (*lastfm.AlbumTags, error) {
	var res lastfm.AlbumTags
	return &res, a.session.Get(&res, api.AlbumGetTagsMethod, params)
}

// RemoveTag removes a tag from an album for the authenticated user.
func (a AlbumSession) RemoveTag(artist, album, tag string) error {
	p := lastfm.AlbumRemoveTagParams{Artist: artist, Album: album, Tag: tag}
	return a.session.Post(nil, api.AlbumRemoveTagMethod, p)
}
