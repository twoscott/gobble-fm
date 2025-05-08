package session

import (
	"github.com/twoscott/gobble-fm/api"
	"github.com/twoscott/gobble-fm/lastfm"
)

type Auth struct {
	session *Session
}

// NewAuth creates and returns a new Auth API route.
func NewAuth(session *Session) *Auth {
	return &Auth{session: session}
}

// Token returns a session token for the Last.fm API. This token can be used to
// authenticate requests to the API. The token is obtained by calling the
// AuthGetToken method of the Last.fm API. The token is cached in the session,
// so subsequent calls to this method will return the same cached token.
func (a Auth) Token() (string, error) {
	var token string
	err := a.session.Get(&token, api.AuthGetTokenMethod, nil)
	return token, err
}

// Session returns a session for the given token. This session can be used to
// authenticate requests to the Last.fm API. The session is obtained by calling
// the AuthGetSession method of the Last.fm API.
func (a Auth) Session(token string) (*lastfm.Session, error) {
	var res lastfm.Session
	return &res, a.session.Post(&res, api.AuthGetSessionMethod, lastfm.SessionParams{Token: token})
}

// MobileSession returns a session for the given mobile user credentials. This
// session can be used to authenticate requests to the Last.fm API. The session
// is obtained by calling the AuthGetMobileSession method of the Last.fm API.
func (a Auth) MobileSession(username, password string) (*lastfm.Session, error) {
	var res lastfm.Session
	return &res, a.session.Post(&res, api.AuthGetMobileSessionMethod, lastfm.MobileSessionParams{
		Username: username,
		Password: password,
	})
}
