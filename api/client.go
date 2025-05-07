package api

// Client provides low-level functionality for making calls to the Last.fm Client.
type Client struct {
	*API
	Album *Album
	User  *User
}

// New returns a new instance of API Client with the given API key.
func NewClient(apiKey string) *Client {
	a := New(apiKey)

	return &Client{
		API:   a,
		Album: NewAlbum(a),
		User:  NewUser(a),
	}
}
