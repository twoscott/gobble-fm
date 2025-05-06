package api

import (
	"github.com/twoscott/go-fm/lastfm"
)

type User struct {
	api *API
}

// NewUser creates and returns a new User API route.
func NewUser(api *API) *User {
	return &User{api: api}
}

// TODO: User can be left empty if call is authenticated; fetches logged in user. Method can be added to Session User object.
//
// Info returns the information of a user.
func (u User) Info(user string) (*lastfm.UserInfo, error) {
	var res lastfm.UserInfo
	p := lastfm.UserInfoParams{User: user}
	return &res, u.api.Get(&res, UserGetInfoMethod, p)
}

// Friends returns the friends of a user.
func (u User) Friends(params lastfm.FriendsParams) (*lastfm.Friends, error) {
	var res lastfm.Friends
	return &res, u.api.Get(&res, UserGetFriendsMethod, params)
}

// RecentTracks returns the recent tracks of a user.
func (u User) RecentTracks(params lastfm.RecentTracksParams) (*lastfm.RecentTracks, error) {
	var res lastfm.RecentTracks
	return &res, u.api.Get(&res, UserGetRecentTracksMethod, params)
}

// RecentTracksExtended returns the recent tracks of a user with extended
// information.
func (u User) RecentTracksExtended(
	params lastfm.RecentTracksParams) (*lastfm.RecentTracksExtended, error) {

	var res lastfm.RecentTracksExtended
	return &res, u.api.Get(&res, UserGetRecentTracksMethod, params)
}
