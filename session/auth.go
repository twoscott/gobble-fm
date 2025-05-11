package session

import (
	"github.com/twoscott/gobble-fm/api"
)

type Auth struct {
	*api.Auth
}

// NewAuth creates and returns a new Auth API route.
func NewAuth(session *Session) *Auth {
	return &Auth{Auth: api.NewAuth(session.API)}
}
