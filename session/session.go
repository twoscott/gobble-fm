package session

import (
	"crypto/md5"
	"fmt"
	"maps"
	"net/http"
	"net/url"
	"slices"

	"github.com/google/go-querystring/query"
	"github.com/twoscott/gobble-fm/api"
)

type Session struct {
	*api.API
	// Secret is the Last.fm API secret used to sign requests.
	Secret string
	// SessionKey is the session key for the Last.fm API session.
	// This key is used to authenticate requests made to the API.
	// Last.fm session keys have infinite lifetime, so you can store it and
	// reuse it for future requests without needing to re-authenticate the user.
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

// SetSessionKey sets the session key for the Last.fm API session. This key is
// used to authenticate requests made to the API. The session key is typically
// obtained after a user has logged in and authorized the application.
//
// Use this method to set the session key manually if you have obtained it
// through other means, such as a login process or an authentication flow, or
// a stored session key from a previous session.
func (s *Session) SetSessionKey(sessionKey string) {
	s.SessionKey = sessionKey
}

// BuildAPIURL constructs a Last.fm API URL with the specified parameters and
// signs the call with an API signature. The signature is generated using the
// session secret and the provided parameters.
func (s Session) BuildAPIURL(params url.Values) string {
	params.Set("api_sig", s.Signature(params))
	return api.BuildAPIURL(params)
}

// Get sends an HTTP GET request to the API using the specified method and
// parameters, and decodes the response into the provided destination.
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

// Post sends an HTTP POST request to the API with the specified method and
// parameters. The response is unmarshaled into the provided destination.
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

// Request sends an HTTP request to the API with the specified parameters and
// unmarshals the response into the provided destination.
//
// Parameters:
//   - dest: A pointer to the variable where the unmarshaled XML response will
//     be stored.
//   - httpMethod: The HTTP method to use for the request (e.g., "GET", "POST").
//   - method: The API method to call, represented as an APIMethod type.
//   - params: The parameters to include in the API request, typically a struct
//     that can be serialized into query parameters.
//
// Returns:
//   - An error if the request fails, the response cannot be unmarshaled,
//     or any other issue occurs.
//
// The function constructs the request URL with the provided parameters, sets
// the necessary headers, sends the request, and processes the response. It uses
// the `query.Values` package to encode parameters and the
// `xml.Unmarshal` function to parse the XML response.
func (s Session) Request(dest any, httpMethod string, method api.APIMethod, params any) error {
	var p url.Values
	var err error

	if params == nil {
		p = url.Values{}
	} else {
		p, err = query.Values(params)
	}
	if err != nil {
		return err
	}

	p.Set("api_key", s.APIKey)
	p.Set("sk", s.SessionKey)
	p.Set("method", method.String())
	url := s.BuildAPIURL(p)

	req, err := http.NewRequest(httpMethod, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", s.UserAgent)
	req.Header.Set("Accept", "application/xml")

	lfm, err := s.DoRequest(req)
	if err != nil {
		return err
	}

	if dest == nil {
		return nil
	}

	err = lfm.UnmarshalInnerXML(dest)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
}

// Signature generates a signature for the given parameters using the session
// secret. The signature is created by concatenating the sorted parameter keys
// and their values, followed by the session secret. The resulting string is
// then hashed using MD5 to produce a hexadecimal representation of the hash.
func (s Session) Signature(params url.Values) string {
	keys := slices.Sorted(maps.Keys(params))

	var sig string
	for _, k := range keys {
		sig += k + params.Get(k)
	}

	sig += s.Secret

	return fmt.Sprintf("%x", md5.Sum([]byte(sig)))
}

// Client is a struct that serves as a central point for making authenticated
// API calls. It embeds a Session and provides fields for interacting with
// different API routes such as Album, Artist, User, etc.
type Client struct {
	*Session
	Album   *Album
	Artist  *Artist
	Auth    *Auth
	Chart   *Chart
	Geo     *Geo
	Library *Library
	Tag     *Tag
	User    *User
}

// New returns a new instance of Session Client with the given API key and secret.
func NewClient(apiKey, secret string) *Client {
	s := New(apiKey, secret)

	return &Client{
		Session: s,
		Album:   NewAlbum(s),
		Artist:  NewArtist(s),
		Auth:    NewAuth(s),
		Chart:   NewChart(s),
		Geo:     NewGeo(s),
		Library: NewLibrary(s),
		Tag:     NewTag(s),
		User:    NewUser(s),
	}
}

// FetchLoginURL fetches a token for the user and returns the URL for the user
// to authorize the application. The token is obtained by calling the
// AuthGetToken method of the Last.fm API. The URL is constructed using the
// API key and the token. If the token cannot be fetched, an error is returned.
func (c Client) FetchLoginURL() (string, error) {
	token, err := c.Auth.Token()
	if err != nil {
		return "", err
	}

	return c.AuthTokenURL(token), nil
}

// Login authenticates a user using their username and password credentials.
// Calls the AuthGetMobileSession method of the Last.fm API and sets the session
// key in the Client.
func (c Client) Login(username, password string) error {
	s, err := c.Auth.MobileSession(username, password)
	if err != nil {
		return err
	}

	c.SessionKey = s.Key
	return nil
}
