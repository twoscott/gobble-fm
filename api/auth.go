package api

import (
	"net/url"

	"github.com/google/go-querystring/query"
	"github.com/twoscott/gobble-fm/lastfm"
)

type AuthURLParams struct {
	APIKey   string `url:"api_key"`
	Callback string `url:"cb,omitempty"`
	Token    string `url:"token,omitempty"`
}

// AuthURL returns the authentication URL for the Last.fm API with the specified
// parameters. This URL can be used to send users to Last.fm for authentication.
// The URL includes the API key and a callback URL if specified.
//
// Parameters:
//   - params.APIKey: The Last.fm API key.
//   - params.Callback: The URL to redirect the user to and provided the
//     authenticated token to as a query parameter.
//   - params.Token: The token for the user to authenticate.
//
// Returns:
//   - The authentication URL as a string.
func AuthURL(params AuthURLParams) string {
	p, err := query.Values(params)

	// Shouldn't happen
	if err != nil {
		p = url.Values{}
		p.Set("api_key", params.APIKey)
		if params.Callback != "" {
			p.Set("cb", params.Callback)
		}
		if params.Token != "" {
			p.Set("token", params.Token)
		}
	}

	return lastfm.AuthURL + "?" + p.Encode()
}

// AuthURL returns the authentication URL for the Last.fm API.
func (a API) AuthURL() string {
	return a.AuthCallbackURL("")
}

// AuthCallbackURL returns the authentication URL for the Last.fm API with a
// callback URL. This URL can be used to send users to Last.fm for
// authentication. The user will be redirected to the callback URL after and
// an authenticated token parameter will be appended to the URL.
//
// https://www.last.fm/api/webauth
func (a API) AuthCallbackURL(callbackURL string) string {
	return AuthURL(AuthURLParams{APIKey: a.APIKey, Callback: callbackURL})
}

// AuthTokenURL returns the authentication URL for the Last.fm API with a
// token. This URL can be used to send users to Last.fm for authentication where
// the provided token will be authenticated.
//
// https://www.last.fm/api/desktopauth
func (a API) AuthTokenURL(token string) string {
	return AuthURL(AuthURLParams{APIKey: a.APIKey, Token: token})
}
