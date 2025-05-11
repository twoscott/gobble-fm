package session

import "github.com/twoscott/gobble-fm/api"

type Library struct {
	*api.Library
}

// NewLibrary creates and returns a new Library API route.
func NewLibrary(session *Session) *Library {
	return &Library{Library: api.NewLibrary(session.API)}
}
