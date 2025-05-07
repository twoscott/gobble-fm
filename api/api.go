package api

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/twoscott/go-fm/lastfm"
)

var (
	BaseEndpoint = "https://ws.audioscrobbler.com"
	Version      = "2.0"
	Endpoint     = BaseEndpoint + "/" + Version + "/"

	DefaultUserAgent = "LastFM (https://github.com/twoscott/go-fm)"
	DefaultRetries   = 5
	DefaultTimeout   = 30
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
	// APIKey is the Last.fm API key used to authenticate requests.
	APIKey string
	// UserAgent is the user agent string sent with each request to the API.
	UserAgent string
	Retries   uint
	Client    HTTPClient
}

// New returns a new instance of API with the given API key.
func New(apiKey string) *API {
	return NewWithTimeout(apiKey, DefaultTimeout)
}

// NewWithTimeout returns a new instance of API with the given API key and
// timeout settings. The timeout is specified in seconds.
func NewWithTimeout(apiKey string, timeout int) *API {
	return &API{
		APIKey:    apiKey,
		UserAgent: DefaultUserAgent,
		Client:    &http.Client{Timeout: time.Duration(timeout) * time.Second},
	}
}

// SetUserAgent sets the user agent for the API client.
func (a *API) SetUserAgent(userAgent string) {
	a.UserAgent = userAgent
}

// SetRetries sets the number of retries for failed requests.
func (a *API) SetRetries(retries uint) {
	a.Retries = retries
}

// AuthURL returns the authentication URL for the Last.fm API. This URL can be
// used to redirect users to Last.fm for authentication. The URL includes the
// API key.
func (a API) AuthURL() string {
	return a.AuthCallbackURL("")
}

// AuthCallbackURL returns the authentication URL for the Last.fm API with a
// callback URL. This URL can be used to redirect users to Last.fm for
// authentication. The URL includes the API key and the callback URL.
func (a API) AuthCallbackURL(callbackURL string) string {
	return a.authURL(callbackURL, "")
}

// AuthTokenURL returns the authentication URL for the Last.fm API with a
// token. This URL can be used to redirect users to Last.fm for authentication.
// The URL includes the API key and the token.
func (a API) AuthTokenURL(token string) string {
	return a.authURL("", token)
}

func (a API) authURL(cb, token string) string {
	p := url.Values{}
	p.Set("api_key", a.APIKey)

	if cb != "" {
		p.Set("cb", cb)
	}
	if token != "" {
		p.Set("token", token)
	}

	return lastfm.AuthURL + "?" + p.Encode()
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
func (a API) Get(dest any, method APIMethod, params any) error {
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
func (a API) Post(dest any, method APIMethod, params any) error {
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
func (a API) Request(dest any, httpMethod string, method APIMethod, params any) error {
	var p url.Values
	var err error

	if params != nil {
		p, err = query.Values(params)
	}
	if err != nil {
		return err
	}

	p.Set("api_key", a.APIKey)
	p.Set("method", method.String())
	url := BuildAPIURL(p)

	req, err := http.NewRequest(httpMethod, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", a.UserAgent)
	req.Header.Set("Accept", "application/xml")

	res, err := a.DoRequest(req)
	if err != nil {
		return err
	}

	if dest == nil {
		return nil
	}

	err = res.UnmarshalInnerXML(dest)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
}

// DoRequest sends the HTTP request and handles the response. It retries the
// request if the response indicates a server error or if the Last.fm API rate
// limit is exceeded. The function decodes the XML response into a LFMWrapper
// instance and checks for errors. If the response is successful, it returns
// the LFMWrapper instance. If an error occurs, it returns the error.
func (a API) DoRequest(req *http.Request) (*LFMWrapper, error) {
	var (
		res   *http.Response
		lfm   LFMWrapper
		lferr *LastFMError
		err   error
	)

	for i := uint(0); i <= a.Retries; i++ {
		res, err = a.Client.Do(req)
		if err != nil {
			return nil, err
		}

		err = xml.NewDecoder(res.Body).Decode(&lfm)
		res.Body.Close()
		if err == nil {
			lferr, _ = lfm.UnwrapError()
		}

		if res.StatusCode >= 500 || res.StatusCode == http.StatusTooManyRequests {
			continue
		}
		if lferr != nil && lferr.IsRateLimit() {
			continue
		}

		break
	}

	if lferr != nil {
		return nil, lferr
	}
	if res.StatusCode < http.StatusOK || res.StatusCode > http.StatusIMUsed {
		return nil, NewHTTPError(res)
	}
	if errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("invalid XML response: %w", err)
	}
	if err != nil {
		return nil, err
	}

	return &lfm, nil
}
