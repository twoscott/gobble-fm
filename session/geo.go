package session

import "github.com/twoscott/go-fm/api"

type Geo struct {
	*api.Geo
}

// NewGeo creates and returns a new Geo API route.
func NewGeo(session *Session) *Geo {
	return &Geo{Geo: api.NewGeo(session.API)}
}
