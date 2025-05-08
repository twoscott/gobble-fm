package session

import "github.com/twoscott/gobble-fm/api"

type Chart struct {
	*api.Chart
}

// NewChart creates and returns a new Chart API route.
func NewChart(session *Session) *Chart {
	return &Chart{Chart: api.NewChart(session.API)}
}
