package api

import (
	"github.com/twoscott/gobble-fm/lastfm"
)

type AuthURLParams struct {
	APIKey   string `url:"api_key"`
	Callback string `url:"cb,omitempty"`
	Token    string `url:"token,omitempty"`
}

type Auth struct {
	api *API
}

// NewAuth creates and returns a new Auth API route.
func NewAuth(api *API) *Auth {
	return &Auth{api: api}
}

// Token returns a session token for the Last.fm API. This token can be used to
// authenticate requests to the API. The token is obtained by calling the
// AuthGetToken method of the Last.fm API. The token is cached in the session,
// so subsequent calls to this method will return the same cached token.
func (a Auth) Token() (string, error) {
	var token string
	err := a.api.GetSigned(&token, AuthGetTokenMethod, nil)
	return token, err
}

// Session returns a session for the given token. This session can be used to
// authenticate requests to the Last.fm API. The session is obtained by calling
// the AuthGetSession method of the Last.fm API.
func (a Auth) Session(token string) (*lastfm.Session, error) {
	var res lastfm.Session
	p := lastfm.SessionParams{Token: token}
	return &res, a.api.PostSigned(&res, AuthGetSessionMethod, p)
}

// MobileSession returns a session for the given mobile user credentials. This
// session can be used to authenticate requests to the Last.fm API. The session
// is obtained by calling the AuthGetMobileSession method of the Last.fm API.
func (a Auth) MobileSession(username, password string) (*lastfm.Session, error) {
	var res lastfm.Session
	p := lastfm.MobileSessionParams{Username: username, Password: password}
	return &res, a.api.PostSigned(&res, AuthGetMobileSessionMethod, p)
}
