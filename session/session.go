package session

import (
	"github.com/twoscott/go-fm/api"
)

type Session struct {
	*api.API
	Secret     string
	SessionKey string
}

// New returns a new instance of Session with the given API key and secret.
func New(apiKey, secret string) *Session {
	return NewWithTimeout(apiKey, secret, api.DefaultTimeout)
}

// NewWithTimeout returns a new instance of Session with the given API key,
// secret, and timeout settings. The timeout is specified in seconds.
func NewWithTimeout(apiKey, secret string, timeout int) *Session {
	return &Session{
		API:    api.NewWithTimeout(apiKey, timeout),
		Secret: secret,
	}
}
