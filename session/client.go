package session

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
	Track   *Track
	User    *User
}

// New returns a new instance of Session Client with the given API key and
// secret.
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

func newClient(s *Session) *Client {
	return &Client{
		Session: s,
		Album:   NewAlbum(s),
		Artist:  NewArtist(s),
		Auth:    NewAuth(s),
		Chart:   NewChart(s),
		Geo:     NewGeo(s),
		Library: NewLibrary(s),
		Tag:     NewTag(s),
		Track:   NewTrack(s),
		User:    NewUser(s),
	}
}

// TokenLogin authenticates a user using an authorized token obtained from the
// AuthGetToken method of the Last.fm API. The user must have authorised the
// token via the URL returned by AuthTokenURL. The token is used to create a
// session key, which is then set in the Client. If the token cannot be used,
// an error is returned.
func (c *Client) TokenLogin(token string) error {
	s, err := c.Auth.Session(token)
	if err != nil {
		return err
	}

	c.SessionKey = s.Key
	return nil
}

// Login authenticates a user using their username and password credentials.
// Calls the AuthGetMobileSession method of the Last.fm API and sets the session
// key in the Client.
func (c *Client) Login(username, password string) error {
	s, err := c.Auth.MobileSession(username, password)
	if err != nil {
		return err
	}

	c.SessionKey = s.Key
	return nil
}
