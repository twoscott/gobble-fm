package client

import (
	"github.com/twoscott/gobble-fm/v2/api"
	"github.com/twoscott/gobble-fm/v2/lastfm"
)

// Client is the main struct that provides access to various API services.
// It embeds the AuthAPI struct and includes fields for accessing specific
// service modules such as Album, Artist, Auth, User, etc.
type Client struct {
	*api.AuthAPI
	Album   *Album
	Artist  *Artist
	Auth    *Auth
	Chart   *Chart
	Geo     *Geo
	Library *Library
	Tag     *Tag
	Track   *Track
	User    *UserAPI
}

// New returns a new instance of Client with the given API key and secret.
// Returns an error if the API key or secret is missing.
func New(apiKey, secret string) (*Client, error) {
	return newClient(api.New(apiKey, secret))
}

// NewWithTimeout returns a new instance of Client with the given API key,
// secret, and timeout settings. The timeout is specified in seconds and is used
// to configure the HTTP client for making API requests. This allows for better
// control over network timeouts when interacting with the API. Returns an
// error if the API key or secret is missing.
func NewWithTimeout(apiKey, secret string, timeout int) (*Client, error) {
	return newClient(api.NewWithTimeout(apiKey, secret, timeout))
}

func newClient(a *api.AuthAPI, err error) (*Client, error) {
	if err != nil {
		return nil, err
	}

	return &Client{
		AuthAPI: a,
		Album:   NewAlbum(a.API),
		Artist:  NewArtist(a.API),
		Auth:    NewAuth(a),
		Chart:   NewChart(a.API),
		Geo:     NewGeo(a.API),
		Library: NewLibrary(a.API),
		Tag:     NewTag(a.API),
		Track:   NewTrack(a.API),
		User:    NewUserAPI(a.API),
	}, nil
}

// NewSession creates a new SessionClient for the provided Last.fm session based
// on Client. Returns an error if the session key is missing.
func (c *Client) NewSession(s lastfm.Session) (*SessionClient, error) {
	return newSession(api.NewAPISession(c.AuthAPI, s))
}

// NewSessionFromKey creates a new SessionClient for the provided Last.fm
// session key based on Client. Returns an error if the session key is missing.
func (c *Client) NewSessionFromKey(sessionKey string) (*SessionClient, error) {
	return newSession(api.NewAPISessionFromKey(c.AuthAPI, sessionKey))
}

// KeyClient allows you to interact with the Last.fm API with only an API key. Namely this means
// you can't access Auth API methods that require an API secret.
type KeyClient struct {
	*api.API
	Album   *Album
	Artist  *Artist
	Chart   *Chart
	Geo     *Geo
	Library *Library
	Tag     *Tag
	Track   *Track
	User    *UserAPI
}

// NewKeyOnly returns a new instance of KeyClient with the given API key but
// without a Last.fm API secret. This is useful if you don't plan to use the API
// secret to sign requests to the API.
//
// If you want to use the API secret to sign requests, you should use New or
// NewWithTimeout instead, which will create a Client that gives you access to
// Auth API methods.
//
// Returns an error if the API key is missing.
func NewKeyOnly(apiKey string) (*KeyClient, error) {
	return newKeyOnly(api.NewKeyOnly(apiKey))
}

func newKeyOnly(a *api.API, err error) (*KeyClient, error) {
	if err != nil {
		return nil, err
	}

	return &KeyClient{
		API:     a,
		Album:   NewAlbum(a),
		Artist:  NewArtist(a),
		Chart:   NewChart(a),
		Geo:     NewGeo(a),
		Library: NewLibrary(a),
		Tag:     NewTag(a),
		Track:   NewTrack(a),
		User:    NewUserAPI(a),
	}, nil
}

// SessionClient is a struct that serves as a central point for making
// authenticated API calls. It embeds a Session and provides fields for
// interacting with different API routes such as Album, Artist, User, etc.
type SessionClient struct {
	*api.Session
	Album   *AlbumSession
	Artist  *ArtistSession
	Chart   *Chart
	Geo     *Geo
	Library *Library
	Tag     *Tag
	Track   *TrackSession
	User    *UserSession
}

// NewSession returns a new session client from a user's session details.
// Returns an error if the API key, secret, or session key is missing.
func NewSession(apiKey, secret string, session lastfm.Session) (*SessionClient, error) {
	return newSession(api.NewSession(apiKey, secret, session))
}

// NewSessionFromKey returns a new session client from a user's session key.
// Returns an error if the API key, secret, or session key is missing.
func NewSessionFromKey(apiKey, secret, sessionKey string) (*SessionClient, error) {
	return newSession(api.NewSessionFromKey(apiKey, secret, sessionKey))
}

func newSession(s *api.Session, err error) (*SessionClient, error) {
	if err != nil {
		return nil, err
	}

	return &SessionClient{
		Session: s,
		Album:   NewAlbumSession(s),
		Artist:  NewArtistSession(s),
		Chart:   NewChart(s.API),
		Geo:     NewGeo(s.API),
		Library: NewLibrary(s.API),
		Tag:     NewTag(s.API),
		Track:   NewTrackSession(s),
		User:    NewUserSession(s),
	}, nil
}

// TokenLogin authenticates a user using an authorized token obtained from the
// AuthGetToken method of the Last.fm API. The user must have authorised the
// token via the URL returned by AuthTokenURL. The token is used to create a
// session, which is returned as a SessionClient. If the token cannot be used,
// an error is returned.
func (c *Client) TokenLogin(token string) (*SessionClient, error) {
	s, err := c.Auth.Session(token)
	if err != nil {
		return nil, err
	}

	return c.NewSession(*s)
}

// Login authenticates a user using their username and password credentials.
// Calls the AuthGetMobileSession method of the Last.fm API and returns the
// session as a SessionClient.
func (c *Client) Login(username, password string) (*SessionClient, error) {
	s, err := c.Auth.MobileSession(username, password)
	if err != nil {
		return nil, err
	}

	return c.NewSession(*s)
}
