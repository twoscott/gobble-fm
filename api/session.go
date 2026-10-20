package api

import (
	"errors"
	"net/http"

	"github.com/google/go-querystring/query"
	"github.com/twoscott/gobble-fm/v2/lastfm"
)

// Session embeds AuthAPI and provides methods for making authenticated requests
// on behalf of a user, signed with the user's session key.
type Session struct {
	*AuthAPI
	// name is the name of the authenticated user.
	name string
	// key is the session key for the Last.fm API session. This key is
	// used to authenticate requests made to the API. Last.fm session keys have
	// infinite lifetime, so you can store it and reuse it for future requests
	// without needing to re-authenticate the user.
	key string
	// subscriber shows whether the authenticated user is a subscriber or not.
	subscriber bool
}

// NewSession returns a new API session from a user's session details.
func NewSession(apiKey, secret string, session lastfm.Session) *Session {
	return NewAPISession(New(apiKey, secret), session)
}

// NewSessionFromKey returns a new API session from a user's session key.
func NewSessionFromKey(apiKey, secret, sessionKey string) *Session {
	return NewAPISession(New(apiKey, secret), lastfm.Session{Key: sessionKey})
}

// NewAPISession returns a new API session from an AuthAPI instance.
func NewAPISession(api *AuthAPI, session lastfm.Session) *Session {
	return &Session{
		AuthAPI:    api,
		name:       session.Name,
		key:        session.Key,
		subscriber: session.Subscriber.Bool(),
	}
}

// // NewAPISession returns a new API session from an AuthAPI instance.
func NewAPISessionFromKey(api *AuthAPI, sessionKey string) *Session {
	return NewAPISession(api, lastfm.Session{Key: sessionKey})
}

func (s *Session) Key() string {
	return s.key
}

func (s *Session) Name() string {
	return s.name
}

func (s *Session) Subscriber() bool {
	return s.subscriber
}

// Get sends an authenticated HTTP GET request to the API using the specified
// method and parameters, and decodes the response into the provided destination.
//
// Parameters:
//   - dest: A pointer to the variable where the response will be unmarshaled.
//   - method: The APIMethod representing the endpoint to call.
//   - params: The parameters to include in the request.
//
// Returns:
//   - An error if the request fails or the response cannot be decoded.
func (s *Session) Get(dest any, method APIMethod, params any) error {
	return s.Request(dest, http.MethodGet, method, params)
}

// Post sends an authenticated HTTP POST request to the API with the specified
// method and parameters. The response is unmarshaled into the provided
// destination.
//
// Parameters:
//   - dest: A pointer to the variable where the response will be unmarshaled.
//   - method: The APIMethod representing the API endpoint to call.
//   - params: The parameters to include in the POST request.
//
// Returns:
//   - An error if the request fails or the response cannot be unmarshaled.
func (s *Session) Post(dest any, method APIMethod, params any) error {
	return s.Request(dest, http.MethodPost, method, params)
}

// Request sends an authenticated HTTP request to the API using the specified
// parameters and unmarshals the response into the provided destination.
//
// Parameters:
//   - dest: A pointer to the variable where the response will be unmarshaled.
//   - httpMethod: The HTTP method to use for the request (e.g., "GET", "POST").
//   - method: The APIMethod representing the endpoint to call.
//   - params: The parameters to include in the request.
//
// Returns:
//   - An error if the request fails or the response cannot be unmarshaled.
func (s *Session) Request(dest any, httpMethod string, method APIMethod, params any) error {
	err := s.CheckCredentials(RequestLevelSession)
	if err != nil {
		return err
	}

	p, err := query.Values(params)
	if err != nil {
		return err
	}

	p.Set("api_key", s.APIKey)
	p.Set("sk", s.key)
	p.Set("method", method.String())
	p.Set("api_sig", s.Signature(p))

	switch httpMethod {
	case http.MethodGet:
		return s.GetURL(dest, BuildAPIURL(p))
	case http.MethodPost:
		return s.PostBody(dest, Endpoint, p.Encode())
	default:
		return errors.New("unsupported HTTP method")
	}
}

// CheckCredentials verifies the authentication level required for an API
// request and ensures the necessary credentials are present. It checks the
// presence of the Last.fm API session key for requests requiring a session,
// the secret for requests requiring a secret, and the API key for requests
// requiring an API key. Returns an error if required credentials are missing.
//
// Parameters:
//   - level: The RequestLevel indicating the level of authorization required.
//
// Returns:
//   - An error if the required authentication credentials are not present.
func (s *Session) CheckCredentials(level RequestLevel) error {
	switch level {
	case RequestLevelSession:
		if s.key == "" {
			return NewLastFMError(ErrSessionRequired, SessionRequiredMessage)
		}
		fallthrough
	default:
		return s.AuthAPI.CheckCredentials(level)
	}
}
