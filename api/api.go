// Package api provides a client for interacting with the Last.fm API.
// It includes methods for constructing API URLs, generating signatures,
// and making unauthenticated requests to various Last.fm API endpoints.
//
// The package supports multiple levels of authorization, including API key,
// secret, and session-based authentication. It also provides utilities for
// parsing parameters, handling retries, and unmarshaling XML responses.
//
// The API struct makes unsigned requests using only an API key. AuthAPI adds
// the API secret for signed requests, and Session adds a user's session key
// for authenticated requests. For higher-level access to specific service
// modules such as Album, Artist, User, and more, use the client package.
//
// Key Features:
//   - Build API URLs with query parameters.
//   - Generate signatures for API requests using the API secret.
//   - Make GET and POST requests to the API with automatic XML unmarshaling.
//
// Usage:
//   - Create a new API with your API key, or an AuthAPI with your API key and
//     secret.
//   - Use Get, Post, GetSigned, or PostSigned to call specific API methods.
//   - Handle responses and errors using the provided types and utilities.
//   - Customize the client with user agent and timeout settings.
//   - Use the `Request` method for general-purpose API requests.
//
// For more information about the Last.fm API, visit:
// https://www.last.fm/api
package api

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
)

var (
	BaseEndpoint = "https://ws.audioscrobbler.com"
	Version      = "2.0"
	Endpoint     = BaseEndpoint + "/" + Version + "/"
)

const (
	DefaultUserAgent      = "LastFM (https://github.com/twoscott/gobble-fm)"
	DefaultRetries   uint = 5
	DefaultTimeout        = 30
)

// RequestLevel specifies the level of authorisation and authentication required
// for an API request.
type RequestLevel int

const (
	RequestLevelNone RequestLevel = iota
	RequestLevelAPIKey
	RequestLevelSecret
	RequestLevelSession
)

// HTTPClient is an interface that defines the Do method for making HTTP
// requests. This allows for easier testing and mocking of HTTP requests.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// API represents the Last.fm API client. It provides methods for making
// requests to the Last.fm API and handling responses. The API client is
// initialized with an API key and can be configured with a user agent and
// timeout settings. It also supports retries for failed requests.
type API struct {
	// apiKey is the Last.fm API key used to authenticate requests.
	apiKey string
	// UserAgent is the user agent string sent with each request to the API.
	UserAgent string
	// Retries is the number of times to retry failed requests.
	Retries uint
	Client  HTTPClient
}

// NewKeyOnly returns a new instance of API with the given API key but without
// a Last.fm API secret. This is useful if you don't plan to use the API secret
// to sign requests to the API such as auth methods. Returns an error if the API
// key is missing.
func NewKeyOnly(apiKey string) (*API, error) {
	t := time.Duration(DefaultTimeout) * time.Second

	a := &API{
		apiKey:    apiKey,
		UserAgent: DefaultUserAgent,
		Retries:   DefaultRetries,
		Client:    &http.Client{Timeout: t},
	}

	if err := a.CheckCredentials(RequestLevelAPIKey); err != nil {
		return nil, err
	}

	return a, nil
}

// APIKey returns the Last.fm API key used to authenticate requests.
func (a *API) APIKey() string {
	return a.apiKey
}

// SetUserAgent sets the user agent for the API client.
func (a *API) SetUserAgent(userAgent string) {
	a.UserAgent = userAgent
}

// SetRetries sets the number of retries for failed requests.
func (a *API) SetRetries(retries uint) {
	a.Retries = retries
}

// AuthURL returns the authorization URL for the Last.fm API. This method should
// be used for web authentication if you set a callback URL when creating your
// API account. Otherwise, use AuthCallbackURL and provide a custom callback
// URL.
//
// https://www.last.fm/api/webauth
func (a *API) AuthURL() string {
	return a.AuthCallbackURL("")
}

// AuthCallbackURL returns the authorization URL for the Last.fm API with a
// callback URL. This URL can be used for web authentication to send users to
// Last.fm for authentication. The user will be redirected to the callback URL
// after and an authorized token query parameter will be appended to the URL.
//
// https://www.last.fm/api/webauth
func (a *API) AuthCallbackURL(callbackURL string) string {
	return AuthURL(AuthURLParams{APIKey: a.apiKey, Callback: callbackURL})
}

// AuthTokenURL returns the authorization URL for the Last.fm API with a
// token. This URL can be used to send users to Last.fm for authentication where
// the provided token will be authorized.
//
// https://www.last.fm/api/desktopauth
func (a *API) AuthTokenURL(token string) string {
	return AuthURL(AuthURLParams{APIKey: a.apiKey, Token: token})
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
func (a *API) Get(dest any, method APIMethod, params any) error {
	return a.Request(dest, http.MethodGet, method, params)
}

// Post sends an HTTP POST request to the API using the specified method and
// parameters, and decodes the response into the provided destination.
//
// Parameters:
//   - dest: A pointer to the variable where the response will be unmarshaled.
//   - method: The APIMethod representing the endpoint to call.
//   - params: The parameters to include in the request.
//
// Returns:
//   - An error if the request fails or the response cannot be decoded.
func (a *API) Post(dest any, method APIMethod, params any) error {
	return a.Request(dest, http.MethodPost, method, params)
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
func (a *API) Request(dest any, httpMethod string, method APIMethod, params any) error {
	err := a.CheckCredentials(RequestLevelAPIKey)
	if err != nil {
		return err
	}

	p, err := query.Values(params)
	if err != nil {
		return err
	}

	p.Set("api_key", a.apiKey)
	p.Set("method", method.String())

	switch httpMethod {
	case http.MethodGet:
		return a.GetURL(dest, BuildAPIURL(p))
	case http.MethodPost:
		return a.PostBody(dest, Endpoint, p.Encode())
	default:
		return errors.New("unsupported HTTP method")
	}
}

// GetURL sends an HTTP GET request to the specified URL and unmarshals the
// response into the provided destination.
func (a *API) GetURL(dest any, url string) error {
	return a.tryRequest(dest, http.MethodGet, url, "")
}

// PostBody sends an HTTP POST request to the specified URL with the given
// request body and unmarshals the response into the provided destination.
func (a *API) PostBody(dest any, url, body string) error {
	return a.tryRequest(dest, http.MethodPost, url, body)
}

func (a *API) tryRequest(dest any, method, url, body string) error {
	var (
		res   *http.Response
		err   error
		lfm   LFMWrapper
		lferr *LastFMError
	)

	for i := uint(0); i <= a.Retries; i++ {
		var req *http.Request

		// reset so a failed decode doesn't leave the previous attempt's
		// response or error in place
		lfm, lferr = LFMWrapper{}, nil

		switch method {
		case http.MethodGet:
			req, err = a.createGetRequest(url)
		case http.MethodPost:
			req, err = a.createPostRequest(url, body)
		default:
			req, err = a.createRequest(method, url, body)
		}
		if err != nil {
			return err
		}

		res, err = a.Client.Do(req)
		if err != nil {
			return err
		}

		err = xml.NewDecoder(res.Body).Decode(&lfm)
		res.Body.Close()
		if err == nil {
			lferr, _ = lfm.UnwrapError()
		}

		if res.StatusCode == http.StatusTooManyRequests || res.StatusCode >= 500 {
			continue
		}
		if lferr != nil && lferr.ShouldRetry() {
			continue
		}

		break
	}

	if lferr != nil {
		return lferr.WrapResponse(res)
	}
	if res.StatusCode < http.StatusOK || res.StatusCode >= 300 {
		return NewHTTPError(res)
	}
	if errors.Is(err, io.EOF) {
		return fmt.Errorf("invalid XML response: %w", err)
	}
	if err != nil {
		return err
	}

	if dest == nil {
		return nil
	}
	if err = lfm.UnmarshalInnerXML(dest); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
}

func (a *API) createGetRequest(url string) (*http.Request, error) {
	return a.createRequest(http.MethodGet, url, "")
}

func (a *API) createPostRequest(url, body string) (*http.Request, error) {
	req, err := a.createRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return req, nil
}

func (a *API) createRequest(method, url, body string) (*http.Request, error) {
	var r io.Reader
	if body != "" {
		r = strings.NewReader(body)
	}

	req, err := http.NewRequest(method, url, r)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", a.UserAgent)
	req.Header.Set("Accept", "application/xml")

	return req, nil
}

// BuildAPIURL constructs a Last.fm API URL with the specified parameters.
func BuildAPIURL(params url.Values) string {
	return Endpoint + "?" + params.Encode()
}

// CheckCredentials verifies the authorization level required for an API
// request and ensures the necessary credentials are present. It checks the
// presence of the API secret for requests requiring a session or secret, and
// the API key for requests requiring an API key. Returns an error if required
// credentials are missing.
//
// Parameters:
//   - level: The RequestLevel indicating the level of authorization required.
//
// Returns:
//   - An error if the required authentication credentials are not present.
func (a *API) CheckCredentials(level RequestLevel) error {
	switch level {
	case RequestLevelSession, RequestLevelSecret, RequestLevelAPIKey:
		if a.apiKey == "" {
			return NewLastFMError(ErrAPIKeyMissing, APIKeyMissingMessage)
		}
		fallthrough
	default:
		if a.Client == nil {
			return errors.New("API HTTP Client is nil")
		}
	}

	return nil
}
