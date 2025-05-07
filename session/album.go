package session

import (
	"github.com/twoscott/go-fm/api"
	"github.com/twoscott/go-fm/lastfm"
)

type Album struct {
	*api.Album
	session *Session
}

// NewAlbum creates and returns a new Album API route.
func NewAlbum(session *Session) *Album {
	return &Album{Album: api.NewAlbum(session.API), session: session}
}

// AddTags adds tags to an album.
func (a Album) AddTags(artist, album string, tags []string) error {
	p := lastfm.AlbumAddTagsParams{Artist: artist, Album: album, Tags: tags}
	return a.session.Post(nil, api.AlbumAddTagsMethod, p)
}
