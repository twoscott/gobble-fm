package session

import (
	"github.com/twoscott/gobble-fm/api"
	"github.com/twoscott/gobble-fm/lastfm"
)

type Album struct {
	*api.Album
	session *Session
}

// NewAlbum creates and returns a new Album API route.
func NewAlbum(session *Session) *Album {
	return &Album{Album: api.NewAlbum(session.API), session: session}
}

// AddTags adds tags to an album for the authenticated user.
func (a Album) AddTags(artist, album string, tags []string) error {
	p := lastfm.AlbumAddTagsParams{Artist: artist, Album: album, Tags: tags}
	return a.session.Post(nil, api.AlbumAddTagsMethod, p)
}

// SelfTags returns the tags of an album for the authenticated user by
// artist and album name.
func (a Album) SelfTags(params lastfm.AlbumSelfTagsParams) (*lastfm.AlbumTags, error) {
	var res lastfm.AlbumTags
	return &res, a.session.Get(&res, api.AlbumGetTagsMethod, params)
}

// SelfTagsByMBID returns the tags of an album for the authenticated user by MBID.
func (a Album) SelfTagsByMBID(params lastfm.AlbumSelfTagsMBIDParams) (*lastfm.AlbumTags, error) {
	var res lastfm.AlbumTags
	return &res, a.session.Get(&res, api.AlbumGetTagsMethod, params)
}

// RemoveTag removes a tag from an album for the authenticated user.
func (a Album) RemoveTag(artist, album, tag string) error {
	p := lastfm.AlbumRemoveTagParams{Artist: artist, Album: album, Tag: tag}
	return a.session.Post(nil, api.AlbumRemoveTagMethod, p)
}
