package session

import (
	"github.com/twoscott/go-fm/api"
	"github.com/twoscott/go-fm/lastfm"
)

type User struct {
	*api.User
	session *Session
}

// NewUser creates and returns a new User API route.
func NewUser(session *Session) *User {
	return &User{User: api.NewUser(session.API), session: session}
}

// Info returns the information the authenticated user.
func (u User) SelfInfo() (*lastfm.UserInfo, error) {
	var res lastfm.UserInfo
	return &res, u.session.Get(&res, api.UserGetInfoMethod, nil)
}
