package session

import "github.com/twoscott/gobble-fm/api"

type Tag struct {
	*api.Tag
}

// NewTag creates and returns a new Tag API route.
func NewTag(session *Session) *Tag {
	return &Tag{Tag: api.NewTag(session.API)}
}
