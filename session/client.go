package session

// Client provides low-level functionality for making calls to the Last.fm Client.
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

// MobileLogin authenticates a user using their mobile credentials. Calls the
// AuthGetMobileSession method of the Last.fm API and sets the session key
// in the Client.
func (c Client) MobileLogin(username, password string) error {
	s, err := c.Auth.MobileSession(username, password)
	if err != nil {
		return err
	}

	c.SessionKey = s.Key
	return nil
}
