package session

import (
	"github.com/twoscott/gobble-fm/api"
	"github.com/twoscott/gobble-fm/lastfm"
)

type Artist struct {
	*api.Artist
	session *Session
}

// NewArtist creates and returns a new Artist API route.
func NewArtist(session *Session) *Artist {
	return &Artist{Artist: api.NewArtist(session.API), session: session}
}

// AddTags adds tags to an artist for the authenticated user.
func (a Artist) AddTags(artist string, tags []string) error {
	p := lastfm.ArtistAddTagsParams{Artist: artist, Tags: tags}
	return a.session.Post(nil, api.ArtistAddTagsMethod, p)
}

// SelfTags returns the tags of an artist for the authenticated user by artist
// name.
func (a Artist) SelfTags(params lastfm.ArtistSelfTagsParams) (*lastfm.ArtistTags, error) {
	var res lastfm.ArtistTags
	return &res, a.session.Get(&res, api.ArtistGetTagsMethod, params)
}

// SelfTagsByMBID returns the tags of an artist for the authenticated user by
// MBID.
func (a Artist) SelfTagsByMBID(params lastfm.ArtistSelfTagsMBIDParams) (*lastfm.ArtistTags, error) {
	var res lastfm.ArtistTags
	return &res, a.session.Get(&res, api.ArtistGetTagsMethod, params)
}

// RemoveTag removes a tag from an artist for the authenticated user.
func (a Artist) RemoveTag(artist string, tag string) error {
	p := lastfm.ArtistRemoveTagParams{Artist: artist, Tag: tag}
	return a.session.Post(nil, api.ArtistRemoveTagMethod, p)
}
