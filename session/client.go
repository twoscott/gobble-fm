package session

// Client provides low-level functionality for making calls to the Last.fm Client.
type Client struct {
	*Session
	User *User
}

// New returns a new instance of Session Client with the given API key and secret.
func NewClient(apiKey, secret string) *Client {
	s := New(apiKey, secret)

	return &Client{
		Session: s,
		User:    NewUser(s),
	}
}
