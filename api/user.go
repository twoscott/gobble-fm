package api

import (
	"github.com/twoscott/gobble-fm/lastfm"
)

type recentTracksExtendedParams struct {
	lastfm.RecentTracksParams
	Extended bool `url:"extended,int,omitempty"`
}

type personalTagsExtendedParams struct {
	lastfm.UserTagsParams
	Type lastfm.TagType `url:"taggingtype"`
}

type User struct {
	api *API
}

// NewUser creates and returns a new User API route.
func NewUser(api *API) *User {
	return &User{api: api}
}

// Friends returns the friends of a user.
func (u User) Friends(params lastfm.FriendsParams) (*lastfm.Friends, error) {
	var res lastfm.Friends
	return &res, u.api.Get(&res, UserGetFriendsMethod, params)
}

// Info returns the information of a user.
func (u User) Info(user string) (*lastfm.UserInfo, error) {
	var res lastfm.UserInfo
	p := lastfm.UserInfoParams{User: user}
	return &res, u.api.Get(&res, UserGetInfoMethod, p)
}

// LovedTracks returns the loved tracks of a user.
func (u User) LovedTracks(params lastfm.LovedTracksParams) (*lastfm.LovedTracks, error) {
	var res lastfm.LovedTracks
	return &res, u.api.Get(&res, UserGetLovedTracksMethod, params)
}

// TaggedAlbums returns the albums tagged by a user with the given tag.
func (u User) TaggedAlbums(params lastfm.UserTagsParams) (*lastfm.UserAlbumTags, error) {
	var res lastfm.UserAlbumTags
	p := personalTagsExtendedParams{UserTagsParams: params, Type: lastfm.TagTypeAlbum}
	return &res, u.api.Get(&res, UserGetPersonalTagsMethod, p)
}

// TaggedArtists returns the artists tagged by a user with the given tag.
func (u User) TaggedArtists(params lastfm.UserTagsParams) (*lastfm.UserArtistTags, error) {
	var res lastfm.UserArtistTags
	p := personalTagsExtendedParams{UserTagsParams: params, Type: lastfm.TagTypeArtist}
	return &res, u.api.Get(&res, UserGetPersonalTagsMethod, p)
}

// TaggedTracks returns the tracks tagged by a user with the given tag.
func (u User) TaggedTracks(params lastfm.UserTagsParams) (*lastfm.UserTrackTags, error) {
	var res lastfm.UserTrackTags
	p := personalTagsExtendedParams{UserTagsParams: params, Type: lastfm.TagTypeTrack}
	return &res, u.api.Get(&res, UserGetPersonalTagsMethod, p)
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
	p := recentTracksExtendedParams{RecentTracksParams: params, Extended: true}
	return &res, u.api.Get(&res, UserGetRecentTracksMethod, p)
}

// TopAlbums returns the top albums of a user.
func (u User) TopAlbums(params lastfm.UserTopAlbumsParams) (*lastfm.UserTopAlbums, error) {
	var res lastfm.UserTopAlbums
	return &res, u.api.Get(&res, UserGetTopAlbumsMethod, params)
}

// TopArtists returns the top artists of a user.
func (u User) TopArtists(params lastfm.UserTopArtistsParams) (*lastfm.UserTopArtists, error) {
	var res lastfm.UserTopArtists
	return &res, u.api.Get(&res, UserGetTopArtistsMethod, params)
}

// TopTags returns the top tags of a user.
func (u User) TopTags(params lastfm.UserTopTagsParams) (*lastfm.UserTopTags, error) {
	var res lastfm.UserTopTags
	return &res, u.api.Get(&res, UserGetTopTagsMethod, params)
}

// TopTracks returns the top tracks of a user.
func (u User) TopTracks(params lastfm.UserTopTracksParams) (*lastfm.UserTopTracks, error) {
	var res lastfm.UserTopTracks
	return &res, u.api.Get(&res, UserGetTopTracksMethod, params)
}

// WeeklyAlbumChart returns the weekly album chart of a user.
func (u User) WeeklyAlbumChart(
	params lastfm.WeeklyAlbumChartParams) (*lastfm.WeeklyAlbumChart, error) {

	var res lastfm.WeeklyAlbumChart
	return &res, u.api.Get(&res, UserGetWeeklyAlbumChartMethod, params)
}

// WeeklyArtistChart returns the weekly artist chart of a user.
func (u User) WeeklyArtistChart(
	params lastfm.WeeklyArtistChartParams) (*lastfm.WeeklyArtistChart, error) {

	var res lastfm.WeeklyArtistChart
	return &res, u.api.Get(&res, UserGetWeeklyArtistChartMethod, params)
}

// WeeklyChartList returns the weekly chart list of a user.
func (u User) WeeklyChartList(user string) (*lastfm.WeeklyChartList, error) {
	var res lastfm.WeeklyChartList
	p := lastfm.WeeklyChartListParams{User: user}
	return &res, u.api.Get(&res, UserGetWeeklyChartListMethod, p)
}

// WeeklyTrackChart returns the weekly track chart of a user.
func (u User) WeeklyTrackChart(
	params lastfm.WeeklyTrackChartParams) (*lastfm.WeeklyTrackChart, error) {

	var res lastfm.WeeklyTrackChart
	return &res, u.api.Get(&res, UserGetWeeklyTrackChartMethod, params)
}
