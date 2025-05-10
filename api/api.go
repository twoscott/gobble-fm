// Package api provides a client for interacting with the Last.fm API.
// It includes methods for constructing API URLs, generating signatures,
// and making unauthenticated requests to various Last.fm API endpoints.
//
// The package supports multiple levels of authorization, including API key,
// secret, and session-based authentication. It also provides utilities for
// parsing parameters, handling retries, and unmarshaling XML responses.
//
// The main entry point for interacting with the API is the `Client` struct,
// which provides access to specific service modules such as Album, Artist,
// User, and more.
//
// Key Features:
//   - Build API URLs with query parameters.
//   - Generate signatures for API requests using the API secret.
//   - Make GET and POST requests to the API with automatic XML unmarshaling.
//
// Usage:
//   - Create a new Client with your API key and optionally, your secret.
//   - Use the Client to access specific API methods
//   - Handle responses and errors using the provided types and utilities.
//   - Customize the client with user agent and timeout settings.
//   - Use the `Request` method for general-purpose API requests.
//
// For more information about the Last.fm API, visit:
// https://www.last.fm/api
package api

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/twoscott/gobble-fm/lastfm"
)

var (
	BaseEndpoint = "https://ws.audioscrobbler.com"
	Version      = "2.0"
	Endpoint     = BaseEndpoint + "/" + Version + "/"

	DefaultUserAgent      = "LastFM (https://github.com/twoscott/gobble-fm)"
	DefaultRetries   uint = 5
	DefaultTimeout        = 30
)

// BuildAPIURL constructs a Last.fm API URL with the specified parameters.
func BuildAPIURL(params url.Values) string {
	return Endpoint + "?" + params.Encode()
}

// ParseParameters takes a parameter of any type and converts it into a
// url.Values type. This is useful for converting structs into query parameters
// for API requests. If the parameter is nil, an empty url.Values is returned.
// If the parameter cannot be converted into a url.Values, an error is returned.
func ParseParameters(params any) (url.Values, error) {
	var p url.Values
	var err error

	if params == nil {
		p = url.Values{}
	} else {
		p, err = query.Values(params)
	}
	if err != nil {
		return nil, err
	}

	return p, nil
}

// Signature generates a Last.fm API signature for the given parameters and
// secret. The signature is created by concatenating the sorted parameter keys
// and their values, followed by the secret. The resulting string is then
// hashed using MD5 to produce a hexadecimal representation of the hash.
//
// Parameters:
//   - params: The parameters to include in the signature.
//   - secret: The secret key to use for signing the request.
//
// Returns:
//   - A hexadecimal string representing the signature.
//
// https://www.last.fm/api/authspec
func Signature(params url.Values, secret string) string {
	keys := slices.Sorted(maps.Keys(params))

	var sig string
	for _, k := range keys {
		// exclude format and callback params from signature
		if k == "format" || k == "callback" {
			continue
		}

		sig += k + params.Get(k)
	}

	sig += secret
	hash := md5.Sum([]byte(sig))

	return hex.EncodeToString(hash[:])
}

// RequestLevel specifies the level of authorisation and authentication required
// for an API request.
type RequestLevel int

const (
	RequestLevelNone RequestLevel = iota
	RequestLevelAPIKey
	RequestLevelSecret
	RequestLevelSession
)

// APIMethod represents a Last.fm API method parameter.
type APIMethod string

// https://www.last.fm/api
const (
	AlbumAddTagsMethod    APIMethod = "album.addTags"
	AlbumGetInfoMethod    APIMethod = "album.getInfo"
	AlbumGetTagsMethod    APIMethod = "album.getTags"
	AlbumGetTopTagsMethod APIMethod = "album.getTopTags"
	AlbumRemoveTagMethod  APIMethod = "album.removeTag"
	AlbumSearchMethod     APIMethod = "album.search"

	ArtistAddTagsMethod       APIMethod = "artist.addTags"
	ArtistGetCorrectionMethod APIMethod = "artist.getCorrection"
	ArtistGetInfoMethod       APIMethod = "artist.getInfo"
	ArtistGetSimilarMethod    APIMethod = "artist.getSimilar"
	ArtistGetTagsMethod       APIMethod = "artist.getTags"
	ArtistGetTopAlbumsMethod  APIMethod = "artist.getTopAlbums"
	ArtistGetTopTagsMethod    APIMethod = "artist.getTopTags"
	ArtistGetTopTracksMethod  APIMethod = "artist.getTopTracks"
	ArtistRemoveTagMethod     APIMethod = "artist.removeTag"
	ArtistSearchMethod        APIMethod = "artist.search"

	AuthGetMobileSessionMethod APIMethod = "auth.getMobileSession"
	AuthGetSessionMethod       APIMethod = "auth.getSession"
	AuthGetTokenMethod         APIMethod = "auth.getToken"

	ChartGetTopArtistsMethod APIMethod = "chart.getTopArtists"
	ChartGetTopTagsMethod    APIMethod = "chart.getTopTags"
	ChartGetTopTracksMethod  APIMethod = "chart.getTopTracks"

	GeoGetTopArtistsMethod APIMethod = "geo.getTopArtists"
	GeoGetTopTracksMethod  APIMethod = "geo.getTopTracks"

	LibraryGetArtistsMethod APIMethod = "library.getArtists"

	TagGetInfoMethod            APIMethod = "tag.getInfo"
	TagGetSimilarMethod         APIMethod = "tag.getSimilar"
	TagGetTopAlbumsMethod       APIMethod = "tag.getTopAlbums"
	TagGetTopArtistsMethod      APIMethod = "tag.getTopArtists"
	TagGetTopTagsMethod         APIMethod = "tag.getTopTags"
	TagGetTopTracksMethod       APIMethod = "tag.getTopTracks"
	TagGetWeeklyChartListMethod APIMethod = "tag.getWeeklyChartList"

	TrackAddTagsMethod          APIMethod = "track.addTags"
	TrackGetCorrectionMethod    APIMethod = "track.getCorrection"
	TrackGetInfoMethod          APIMethod = "track.getInfo"
	TrackGetSimilarMethod       APIMethod = "track.getSimilar"
	TrackGetTagsMethod          APIMethod = "track.getTags"
	TrackGetTopTagsMethod       APIMethod = "track.getTopTags"
	TrackLoveMethod             APIMethod = "track.love"
	TrackRemoveTagMethod        APIMethod = "track.removeTag"
	TrackScrobbleMethod         APIMethod = "track.scrobble"
	TrackSearchMethod           APIMethod = "track.search"
	TrackUnloveMethod           APIMethod = "track.unlove"
	TrackUpdateNowPlayingMethod APIMethod = "track.updateNowPlaying"

	UserGetFriendsMethod           APIMethod = "user.getFriends"
	UserGetInfoMethod              APIMethod = "user.getInfo"
	UserGetLovedTracksMethod       APIMethod = "user.getLovedTracks"
	UserGetPersonalTagsMethod      APIMethod = "user.getPersonalTags"
	UserGetRecentTracksMethod      APIMethod = "user.getRecentTracks"
	UserGetTopAlbumsMethod         APIMethod = "user.getTopAlbums"
	UserGetTopArtistsMethod        APIMethod = "user.getTopArtists"
	UserGetTopTagsMethod           APIMethod = "user.getTopTags"
	UserGetTopTracksMethod         APIMethod = "user.getTopTracks"
	UserGetWeeklyAlbumChartMethod  APIMethod = "user.getWeeklyAlbumChart"
	UserGetWeeklyArtistChartMethod APIMethod = "user.getWeeklyArtistChart"
	UserGetWeeklyChartListMethod   APIMethod = "user.getWeeklyChartList"
	UserGetWeeklyTrackChartMethod  APIMethod = "user.getWeeklyTrackChart"
)

// String returns the string representation of the APIMethod.
func (m APIMethod) String() string {
	return string(m)
}

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
	// Secret is the Last.fm API secret used to sign requests.
	Secret string
	// UserAgent is the user agent string sent with each request to the API.
	UserAgent string
	// Retries is the number of times to retry failed requests.
	Retries uint
	Client  HTTPClient
}

// New returns a new instance of API with the given API key.
func New(apiKey, secret string) *API {
	return NewWithTimeout(apiKey, secret, DefaultTimeout)
}

// NewWithTimeout returns a new instance of API with the given API key and
// timeout settings. The timeout is specified in seconds.
func NewWithTimeout(apiKey, secret string, timeout int) *API {
	t := time.Duration(timeout) * time.Second

	return &API{
		apiKey:    apiKey,
		Secret:    secret,
		UserAgent: DefaultUserAgent,
		Retries:   DefaultRetries,
		Client:    &http.Client{Timeout: t},
	}
}

// NewAPIOnly returns a new instance of API with the given API key but without
// a Last.fm API secret. This is useful if you don't plan to use the API secret
// to sign requests to the API such as auth methods.
func NewAPIOnly(apiKey string) *API {
	t := time.Duration(DefaultTimeout) * time.Second

	return &API{
		apiKey:    apiKey,
		UserAgent: DefaultUserAgent,
		Retries:   DefaultRetries,
		Client:    &http.Client{Timeout: t},
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

// APIKey returns the API key used to authorise requests to the Last.fm API.
func (a API) APIKey() string {
	return a.apiKey
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
	p := url.Values{}

	p.Set("api_key", params.APIKey)

	if params.Callback != "" {
		p.Set("cb", params.Callback)
	}
	if params.Token != "" {
		p.Set("token", params.Token)
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
	return AuthURL(AuthURLParams{APIKey: a.apiKey, Callback: callbackURL})
}

// AuthTokenURL returns the authentication URL for the Last.fm API with a
// token. This URL can be used to send users to Last.fm for authentication where
// the provided token will be authenticated.
//
// https://www.last.fm/api/desktopauth
func (a API) AuthTokenURL(token string) string {
	return AuthURL(AuthURLParams{APIKey: a.apiKey, Token: token})
}

// Signature generates a signature for the given parameters using the API
// secret. The signature is created by concatenating the sorted parameter keys
// and their values, followed by the session secret. The resulting string is
// then hashed using MD5 to produce a hexadecimal representation of the hash.
func (a API) Signature(params url.Values) string {
	return Signature(params, a.Secret)
}

// CheckCredentials verifies the authentication level required for an API request and
// ensures the necessary credentials are present. It checks the presence of the
// API secret for requests requiring a session or secret, and the API key for
// requests requiring an API key. Returns an error if required credentials are
// missing.
//
// Parameters:
//   - level: The RequestLevel indicating the level of authorization required.
//
// Returns:
//   - An error if the required authentication credentials are not present.
func (a API) CheckCredentials(level RequestLevel) error {
	switch level {
	case RequestLevelSession, RequestLevelSecret:
		if a.Secret == "" {
			return NewLastFMError(ErrSecretRequired, SecretRequiredMessage)
		}
		fallthrough
	case RequestLevelAPIKey:
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
	err := a.CheckCredentials(RequestLevelAPIKey)
	if err != nil {
		return err
	}

	p, err := ParseParameters(params)
	if err != nil {
		return err
	}

	p.Set("api_key", a.apiKey)
	p.Set("method", method.String())

	switch httpMethod {
	case http.MethodGet:
		return a.RequestURL(dest, httpMethod, BuildAPIURL(p))
	case http.MethodPost:
		return a.RequestBody(dest, httpMethod, Endpoint, p.Encode())
	default:
		return errors.New("unsupported HTTP method")
	}
}

// GetSigned sends an HTTP GET request to the API with the specified method and
// parameters, signed with the API secret. The response is unmarshaled into the
// provided destination.
//
// Parameters:
//   - dest: A pointer to the variable where the response will be unmarshaled.
//   - method: The APIMethod representing the API endpoint to call.
//   - params: The parameters to include in the request.
//
// Returns:
//   - An error if the request fails or the response cannot be decoded.
func (a API) GetSigned(dest any, method APIMethod, params any) error {
	return a.RequestSigned(dest, http.MethodGet, method, params)
}

// PostSigned sends an HTTP POST request to the API with the specified method
// and parameters, signed with the API secret. The response is unmarshaled into
// the provided destination.
//
// Parameters:
//   - dest: A pointer to the variable where the response will be unmarshaled.
//   - method: The APIMethod representing the API endpoint to call.
//   - params: The parameters to include in the request.
//
// Returns:
//   - An error if the request fails or the response cannot be decoded.
func (a API) PostSigned(dest any, method APIMethod, params any) error {
	return a.RequestSigned(dest, http.MethodPost, method, params)
}

// RequestSigned sends an HTTP request to the API with the specified parameters
// and method, signed with the API secret. The request is made with the
// specified HTTP method, and the response is unmarshaled into the provided
// destination.
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
func (a API) RequestSigned(dest any, httpMethod string, method APIMethod, params any) error {
	err := a.CheckCredentials(RequestLevelSecret)
	if err != nil {
		return err
	}

	p, err := ParseParameters(params)
	if err != nil {
		return err
	}

	p.Set("api_key", a.apiKey)
	p.Set("method", method.String())
	p.Set("api_sig", a.Signature(p))

	switch httpMethod {
	case http.MethodGet:
		return a.RequestURL(dest, httpMethod, BuildAPIURL(p))
	case http.MethodPost:
		return a.RequestBody(dest, httpMethod, Endpoint, p.Encode())
	default:
		return errors.New("unsupported HTTP method")
	}
}

// RequestURL sends an HTTP request to the API with the specified URL and
// method, and unmarshals the response into the provided destination.
//
// Parameters:
//   - dest: A pointer to the variable where the unmarshaled XML response will
//     be stored.
//   - method: The HTTP method to use for the request (e.g., "GET", "POST").
//   - url: The URL to send the request to.
//
// Returns:
//   - An error if the request fails, the response cannot be unmarshaled,
//     or any other issue occurs.
func (a API) RequestURL(dest any, method, url string) error {
	return a.tryRequest(dest, method, url, "")
}

// RequestBody sends an HTTP request to the API with the specified URL and
// method, with the given request body, and unmarshals the response into the
// provided destination.
//
// Parameters:
//   - dest: A pointer to the variable where the unmarshaled XML response will
//     be stored.
//   - method: The HTTP method to use for the request (e.g., "GET", "POST").
//   - url: The URL to send the request to.
//   - body: The request body to send with the request.
//
// Returns:
//   - An error if the request fails, the response cannot be unmarshaled,
//     or any other issue occurs.
func (a API) RequestBody(dest any, method, url, body string) error {
	return a.tryRequest(dest, method, url, body)
}

func (a API) tryRequest(dest any, method, url, body string) error {
	var (
		res   *http.Response
		lfm   LFMWrapper
		lferr *LastFMError
		err   error
	)

	for i := uint(0); i <= a.Retries; i++ {
		var req *http.Request

		switch method {
		case http.MethodGet:
			req, err = a.createGetRequest(method, url)
		case http.MethodPost:
			req, err = a.createPostRequest(method, url, body)
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

		if res.StatusCode >= 500 || res.StatusCode == http.StatusTooManyRequests {
			continue
		}
		if lferr != nil && lferr.IsRateLimit() {
			continue
		}

		break
	}

	if lferr != nil {
		return lferr.WrapResponse(res)
	}
	if res.StatusCode < http.StatusOK || res.StatusCode > http.StatusIMUsed {
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

func (a API) createGetRequest(method, url string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", a.UserAgent)
	req.Header.Set("Accept", "application/xml")

	return req, nil
}

func (a API) createPostRequest(method, url, body string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", a.UserAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/xml")

	return req, nil
}

// Client is the main struct that provides access to various API services.
// It embeds the API struct and includes fields for accessing specific
// service modules such as Album, Artist, User, etc.
type Client struct {
	*API
	Album   *Album
	Artist  *Artist
	Auth    *Auth
	Chart   *Chart
	Geo     *Geo
	Library *Library
	Tag     *Tag
	Track   *Track
	User    *User
}

// New returns a new instance of API Client with the given API key and secret.
func NewClient(apiKey, secret string) *Client {
	return newClient(New(apiKey, secret))
}

// NewClientWithTimeout returns a new instance of Client with the given API key,
// secret, and timeout settings. The timeout is specified in seconds and is used
// to configure the HTTP client for making API requests. This allows for better
// control over network timeouts when interacting with the API.
func NewClientWithTimeout(apiKey, secret string, timeout int) *Client {
	return newClient(NewWithTimeout(apiKey, secret, timeout))
}

// NewClientAPIOnly returns a new instance of Client with the given API key but
// without a Last.fm API secret. This is useful if you don't plan to use the API
// secret to sign requests to the API such as auth methods.
func NewClientAPIOnly(apiKey string) *Client {
	return newClient(NewAPIOnly(apiKey))
}

func newClient(a *API) *Client {
	return &Client{
		API:     a,
		Album:   NewAlbum(a),
		Artist:  NewArtist(a),
		Auth:    NewAuth(a),
		Chart:   NewChart(a),
		Geo:     NewGeo(a),
		Library: NewLibrary(a),
		Tag:     NewTag(a),
		Track:   NewTrack(a),
		User:    NewUser(a),
	}
}
