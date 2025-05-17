package api

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

// NewClientKeyOnly returns a new instance of Client with the given API key but
// without a Last.fm API secret. This is useful if you don't plan to use the API
// secret to sign requests to the API such as auth methods.
func NewClientKeyOnly(apiKey string) *Client {
	return newClient(NewKeyOnly(apiKey))
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
