// Package session provides functionality for managing authenticated sessions
// with the Last.fm API. It includes methods for setting session keys, verifying
// credentials, and making authenticated API requests.
//
// The Session struct is the core of this package, encapsulating the API client
// and session key required for authenticated requests. It provides methods for
// sending GET and POST requests, as well as a general-purpose Request method
// for handling different HTTP methods.
//
// The Client struct extends the Session functionality by embedding it and
// providing access to various API routes such as Album, Artist, User, and more.
// It serves as a central point for interacting with the Last.fm API.
//
// Key Features:
//   - Create and manage authenticated sessions with the Last.fm API.
//   - Manage Last.fm user authentication and session keys.
//   - Send authenticated HTTP GET and POST requests.
//   - Access different API routes through the Client struct.
//
// Usage:
//   - Create a new Session or Client instance using the provided constructors.
//   - Set the session key using SetSessionKey or through the Login method.
//   - Use the Get, Post, or Request methods to interact with the Last.fm API.
package session

import (
	"errors"
	"net/http"

	"github.com/google/go-querystring/query"
	"github.com/twoscott/gobble-fm/api"
)

type Session struct {
	*api.API
	// SessionKey is the session key for the Last.fm API session. This key is used to
	// authenticate requests made to the API. Last.fm session keys have infinite
	// lifetime, so you can store it and reuse it for future requests without
	// needing to re-authenticate the user.
	SessionKey string
}

// New returns a new instance of Session with the given API key and secret.
func New(apiKey, secret string) *Session {
	return NewWithTimeout(apiKey, secret, api.DefaultTimeout)
}

// NewWithTimeout returns a new instance of Session with the given API key,
// secret, and timeout settings. The timeout is specified in seconds.
func NewWithTimeout(apiKey, secret string, timeout int) *Session {
	return &Session{API: api.NewWithTimeout(apiKey, secret, timeout)}
}

// SetSessionKey sets the session key for the Last.fm API session. This key is
// used to authenticate requests made to the API. The session key is typically
// obtained after a user has logged in and authorized the application.
//
// Use this method to set the session key manually if you have obtained it
// through other means, such as a login process or an authentication flow, or
// a stored session key from a previous session.
func (s *Session) SetSessionKey(key string) {
	s.SessionKey = key
}

// CheckCredentials verifies the authentication level required for an API request and
// ensures the necessary credentials are present. It checks the presence of the
// Last.fm API session key for requests requiring a session, the secret for
// requests requiring a secret, and the API key for requests requiring an API
// key. Returns an error if required credentials are missing.
//
// Parameters:
//   - level: The RequestLevel indicating the level of authorization required.
//
// Returns:
//   - An error if the required authentication credentials are not present.
func (s Session) CheckCredentials(level api.RequestLevel) error {
	switch level {
	case api.RequestLevelSession:
		if s.SessionKey == "" {
			return api.NewLastFMError(api.ErrSessionRequired, api.SessionRequiredMessage)
		}
		fallthrough
	default:
		return s.API.CheckCredentials(level)
	}
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
func (s Session) Get(dest any, method api.APIMethod, params any) error {
	return s.Request(dest, http.MethodGet, method, params)
}

// Post sends an authenticated HTTP POST request to the API with the specified
// method and parameters. The response is unmarshaled into the provided destination.
//
// Parameters:
//   - dest: A pointer to the variable where the response will be unmarshaled.
//   - method: The APIMethod representing the API endpoint to call.
//   - params: The parameters to include in the POST request.
//
// Returns:
//   - An error if the request fails or the response cannot be unmarshaled.
func (s Session) Post(dest any, method api.APIMethod, params any) error {
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
func (s Session) Request(dest any, httpMethod string, method api.APIMethod, params any) error {
	err := s.CheckCredentials(api.RequestLevelSession)
	if err != nil {
		return err
	}

	p, err := query.Values(params)
	if err != nil {
		return err
	}

	p.Set("api_key", s.APIKey)
	p.Set("sk", s.SessionKey)
	p.Set("method", method.String())
	p.Set("api_sig", s.Signature(p))

	switch httpMethod {
	case http.MethodGet:
		return s.GetURL(dest, api.BuildAPIURL(p))
	case http.MethodPost:
		return s.PostBody(dest, api.Endpoint, p.Encode())
	default:
		return errors.New("unsupported HTTP method")
	}
}
