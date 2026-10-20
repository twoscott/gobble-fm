package client

import (
	"github.com/twoscott/gobble-fm/v2/api"
	"github.com/twoscott/gobble-fm/v2/lastfm"
)

type recentTracksExtendedParams struct {
	lastfm.RecentTracksParams
	Extended bool `url:"extended,int,omitempty"`
}

type userTagsExtendedParams struct {
	lastfm.UserTagsParams
	Type lastfm.TagType `url:"taggingtype"`
}

type UserAPI struct {
	api *api.API
}

// NewUserAPI creates and returns a new User API route.
func NewUserAPI(api *api.API) *UserAPI {
	return &UserAPI{api: api}
}

// Friends returns the friends of a user.
func (u UserAPI) Friends(params lastfm.FriendsParams) (*lastfm.Friends, error) {
	var res lastfm.Friends
	return &res, u.api.Get(&res, api.UserGetFriendsMethod, params)
}

// Info returns the information of a user.
func (u UserAPI) Info(user string) (*lastfm.UserInfo, error) {
	var res lastfm.UserInfo
	p := lastfm.UserInfoParams{User: user}
	return &res, u.api.Get(&res, api.UserGetInfoMethod, p)
}

// LovedTracks returns the loved tracks of a user.
func (u UserAPI) LovedTracks(params lastfm.LovedTracksParams) (*lastfm.LovedTracks, error) {
	var res lastfm.LovedTracks
	return &res, u.api.Get(&res, api.UserGetLovedTracksMethod, params)
}

// RecentTrack returns the most recent track of a user. This is a convenience
// method that calls RecentTracks with limit=1.
func (u UserAPI) RecentTrack(user string) (*lastfm.RecentTrack, error) {
	var res lastfm.RecentTrack
	p := lastfm.RecentTracksParams{User: user, Limit: 1}
	return &res, u.api.Get(&res, api.UserGetRecentTracksMethod, p)
}

// RecentTracks returns the recent tracks of a user.
func (u UserAPI) RecentTracks(params lastfm.RecentTracksParams) (*lastfm.RecentTracks, error) {
	var res lastfm.RecentTracks
	return &res, u.api.Get(&res, api.UserGetRecentTracksMethod, params)
}

// RecentTrackExtended returns the most recent track of a user with extended
// information. This is a convenience method that calls RecentTracksExtended
// with limit=1.
func (u UserAPI) RecentTrackExtended(user string) (*lastfm.RecentTrackExtended, error) {
	var res lastfm.RecentTrackExtended
	p := lastfm.RecentTracksParams{User: user, Limit: 1}
	exp := recentTracksExtendedParams{RecentTracksParams: p, Extended: true}
	return &res, u.api.Get(&res, api.UserGetRecentTracksMethod, exp)
}

// RecentTracksExtended returns the recent tracks of a user with extended
// information.
func (u UserAPI) RecentTracksExtended(
	params lastfm.RecentTracksParams) (*lastfm.RecentTracksExtended, error) {

	var res lastfm.RecentTracksExtended
	exp := recentTracksExtendedParams{RecentTracksParams: params, Extended: true}
	return &res, u.api.Get(&res, api.UserGetRecentTracksMethod, exp)
}

// TaggedAlbums returns the albums tagged by a user with the given tag.
func (u UserAPI) TaggedAlbums(params lastfm.UserTagsParams) (*lastfm.UserAlbumTags, error) {
	var res lastfm.UserAlbumTags
	p := userTagsExtendedParams{UserTagsParams: params, Type: lastfm.TagTypeAlbum}
	return &res, u.api.Get(&res, api.UserGetPersonalTagsMethod, p)
}

// TaggedArtists returns the artists tagged by a user with the given tag.
func (u UserAPI) TaggedArtists(params lastfm.UserTagsParams) (*lastfm.UserArtistTags, error) {
	var res lastfm.UserArtistTags
	p := userTagsExtendedParams{UserTagsParams: params, Type: lastfm.TagTypeArtist}
	return &res, u.api.Get(&res, api.UserGetPersonalTagsMethod, p)
}

// TaggedTracks returns the tracks tagged by a user with the given tag.
func (u UserAPI) TaggedTracks(params lastfm.UserTagsParams) (*lastfm.UserTrackTags, error) {
	var res lastfm.UserTrackTags
	p := userTagsExtendedParams{UserTagsParams: params, Type: lastfm.TagTypeTrack}
	return &res, u.api.Get(&res, api.UserGetPersonalTagsMethod, p)
}

// TopAlbums returns the top albums of a user.
func (u UserAPI) TopAlbums(params lastfm.UserTopAlbumsParams) (*lastfm.UserTopAlbums, error) {
	var res lastfm.UserTopAlbums
	return &res, u.api.Get(&res, api.UserGetTopAlbumsMethod, params)
}

// TopArtists returns the top artists of a user.
func (u UserAPI) TopArtists(params lastfm.UserTopArtistsParams) (*lastfm.UserTopArtists, error) {
	var res lastfm.UserTopArtists
	return &res, u.api.Get(&res, api.UserGetTopArtistsMethod, params)
}

// TopTags returns the top tags of a user.
func (u UserAPI) TopTags(params lastfm.UserTopTagsParams) (*lastfm.UserTopTags, error) {
	var res lastfm.UserTopTags
	return &res, u.api.Get(&res, api.UserGetTopTagsMethod, params)
}

// TopTracks returns the top tracks of a user.
func (u UserAPI) TopTracks(params lastfm.UserTopTracksParams) (*lastfm.UserTopTracks, error) {
	var res lastfm.UserTopTracks
	return &res, u.api.Get(&res, api.UserGetTopTracksMethod, params)
}

// WeeklyAlbumChart returns the weekly album chart of a user.
func (u UserAPI) WeeklyAlbumChart(
	params lastfm.WeeklyAlbumChartParams) (*lastfm.WeeklyAlbumChart, error) {

	var res lastfm.WeeklyAlbumChart
	return &res, u.api.Get(&res, api.UserGetWeeklyAlbumChartMethod, params)
}

// WeeklyArtistChart returns the weekly artist chart of a user.
func (u UserAPI) WeeklyArtistChart(
	params lastfm.WeeklyArtistChartParams) (*lastfm.WeeklyArtistChart, error) {

	var res lastfm.WeeklyArtistChart
	return &res, u.api.Get(&res, api.UserGetWeeklyArtistChartMethod, params)
}

// WeeklyChartList returns the weekly chart list of a user.
func (u UserAPI) WeeklyChartList(user string) (*lastfm.WeeklyChartList, error) {
	var res lastfm.WeeklyChartList
	p := lastfm.WeeklyChartListParams{User: user}
	return &res, u.api.Get(&res, api.UserGetWeeklyChartListMethod, p)
}

// WeeklyTrackChart returns the weekly track chart of a user.
func (u UserAPI) WeeklyTrackChart(
	params lastfm.WeeklyTrackChartParams) (*lastfm.WeeklyTrackChart, error) {

	var res lastfm.WeeklyTrackChart
	return &res, u.api.Get(&res, api.UserGetWeeklyTrackChartMethod, params)
}

type UserSession struct {
	*UserAPI
	session *api.Session
}

// NewUserSession creates and returns a new User API route.
func NewUserSession(session *api.Session) *UserSession {
	return &UserSession{UserAPI: NewUserAPI(session.API), session: session}
}

// SelfInfo returns the information of the authenticated user.
func (u UserSession) SelfInfo() (*lastfm.UserInfo, error) {
	var res lastfm.UserInfo
	return &res, u.session.Get(&res, api.UserGetInfoMethod, nil)
}
