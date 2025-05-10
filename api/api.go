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
)

var (
	BaseEndpoint = "https://ws.audioscrobbler.com"
	Version      = "2.0"
	Endpoint     = BaseEndpoint + "/" + Version + "/"

	DefaultUserAgent = "LastFM (https://github.com/twoscott/gobble-fm)"
	DefaultRetries   = 5
	DefaultTimeout   = 30
)

// APIMethod represents a Last.fm API method parameter.
type APIMethod string

// String returns the string representation of the APIMethod.
func (m APIMethod) String() string {
	return string(m)
}

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

// BuildAPIURL constructs a Last.fm API URL with the specified parameters.
func BuildAPIURL(params url.Values) string {
	return Endpoint + "?" + params.Encode()
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
	// Retries is the number of times to retry failed requests.
	Retries uint
	Client  HTTPClient
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

	if params == nil {
		p = url.Values{}
	} else {
		p, err = query.Values(params)
	}
	if err != nil {
		return err
	}

	p.Set("api_key", a.APIKey)
	p.Set("method", method.String())
	url := BuildAPIURL(p)

	return a.RequestURL(dest, httpMethod, url)
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
	Chart   *Chart
	Geo     *Geo
	Library *Library
	Tag     *Tag
	Track   *Track
	User    *User
}

// New returns a new instance of API Client with the given API key.
func NewClient(apiKey string) *Client {
	a := New(apiKey)

	return &Client{
		API:     a,
		Album:   NewAlbum(a),
		Artist:  NewArtist(a),
		Chart:   NewChart(a),
		Geo:     NewGeo(a),
		Library: NewLibrary(a),
		Tag:     NewTag(a),
		Track:   NewTrack(a),
		User:    NewUser(a),
	}
}
