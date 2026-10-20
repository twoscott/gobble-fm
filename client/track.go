package client

import (
	"github.com/twoscott/gobble-fm/v2/api"
	"github.com/twoscott/gobble-fm/v2/lastfm"
)

type Track struct {
	api *api.API
}

// NewTrack creates and returns a new Track API route.
func NewTrack(api *api.API) *Track {
	return &Track{api: api}
}

// Correction returns the track and artist name corrections of a track.
func (t Track) Correction(artist, track string) (*lastfm.TrackCorrection, error) {
	var res lastfm.TrackCorrection
	p := lastfm.TrackCorrectionParams{Artist: artist, Track: track}
	return &res, t.api.Get(&res, api.TrackGetCorrectionMethod, p)
}

// Info returns the information of a track by artist and track name.
func (t Track) Info(params lastfm.TrackInfoParams) (*lastfm.TrackInfo, error) {
	var res lastfm.TrackInfo
	return &res, t.api.Get(&res, api.TrackGetInfoMethod, params)
}

// InfoByMBID returns the information of a track by MBID.
func (t Track) InfoByMBID(params lastfm.TrackInfoMBIDParams) (*lastfm.TrackInfo, error) {
	var res lastfm.TrackInfo
	return &res, t.api.Get(&res, api.TrackGetInfoMethod, params)
}

// UserInfo returns the information of a track for user by artist and track
// name.
func (t Track) UserInfo(params lastfm.TrackUserInfoParams) (*lastfm.TrackUserInfo, error) {
	var res lastfm.TrackUserInfo
	return &res, t.api.Get(&res, api.TrackGetInfoMethod, params)
}

// UserInfoByMBID returns the information of a track for user by MBID.
func (t Track) UserInfoByMBID(
	params lastfm.TrackUserInfoMBIDParams) (*lastfm.TrackUserInfo, error) {

	var res lastfm.TrackUserInfo
	return &res, t.api.Get(&res, api.TrackGetInfoMethod, params)
}

// Similar returns the similar tracks of a track by artist and track name.
func (t Track) Similar(params lastfm.TrackSimilarParams) (*lastfm.SimilarTracks, error) {
	var res lastfm.SimilarTracks
	return &res, t.api.Get(&res, api.TrackGetSimilarMethod, params)
}

// SimilarByMBID returns the similar tracks of a track by MBID.
func (t Track) SimilarByMBID(params lastfm.TrackSimilarMBIDParams) (*lastfm.SimilarTracks, error) {
	var res lastfm.SimilarTracks
	return &res, t.api.Get(&res, api.TrackGetSimilarMethod, params)
}

// Tags returns the tags of a track by artist and track name.
func (t Track) Tags(params lastfm.TrackTagsParams) (*lastfm.TrackTags, error) {
	var res lastfm.TrackTags
	return &res, t.api.Get(&res, api.TrackGetTagsMethod, params)
}

// TagsByMBID returns the tags of a track by MBID.
func (t Track) TagsByMBID(params lastfm.TrackTagsMBIDParams) (*lastfm.TrackTags, error) {
	var res lastfm.TrackTags
	return &res, t.api.Get(&res, api.TrackGetTagsMethod, params)
}

// TopTags returns the top tags of a track by artist and track name.
func (t Track) TopTags(params lastfm.TrackTopTagsParams) (*lastfm.TrackTopTags, error) {
	var res lastfm.TrackTopTags
	return &res, t.api.Get(&res, api.TrackGetTopTagsMethod, params)
}

// TopTagsByMBID returns the top tags of a track by MBID.
func (t Track) TopTagsByMBID(params lastfm.TrackTopTagsMBIDParams) (*lastfm.TrackTopTags, error) {
	var res lastfm.TrackTopTags
	return &res, t.api.Get(&res, api.TrackGetTopTagsMethod, params)
}

// Search searches for tracks by track name, and optionally artist name.
func (t Track) Search(params lastfm.TrackSearchParams) (*lastfm.TrackSearchResult, error) {
	var res lastfm.TrackSearchResult
	return &res, t.api.Get(&res, api.TrackSearchMethod, params)
}

type TrackSession struct {
	*Track
	session *api.Session
}

// NewTrackSession creates and returns a new Track API route.
func NewTrackSession(session *api.Session) *TrackSession {
	return &TrackSession{Track: NewTrack(session.API), session: session}
}

// AddTags adds tags to a track for the authenticated user.
func (t TrackSession) AddTags(artist, track string, tags []string) error {
	p := lastfm.TrackAddTagsParams{Artist: artist, Track: track, Tags: tags}
	return t.session.Post(nil, api.TrackAddTagsMethod, p)
}

// SelfTags returns the tags of a track for the authenticated user by artist and
// track name.
func (t TrackSession) SelfTags(params lastfm.TrackSelfTagsParams) (*lastfm.TrackTags, error) {
	var res lastfm.TrackTags
	return &res, t.session.Get(&res, api.TrackGetTagsMethod, params)
}

// SelfTagsByMBID returns the tags of a track for the authenticated user by
// MBID.
func (t TrackSession) SelfTagsByMBID(params lastfm.TrackSelfTagsMBIDParams) (*lastfm.TrackTags, error) {
	var res lastfm.TrackTags
	return &res, t.session.Get(&res, api.TrackGetTagsMethod, params)
}

// Love marks a track as loved for the authenticated user.
func (t TrackSession) Love(artist, track string) error {
	p := lastfm.TrackLoveParams{Artist: artist, Track: track}
	return t.session.Post(nil, api.TrackLoveMethod, p)
}

// RemoveTag removes a tag from a track for the authenticated user.
func (t TrackSession) RemoveTag(artist, track, tag string) error {
	p := lastfm.TrackRemoveTagParams{Artist: artist, Track: track, Tag: tag}
	return t.session.Post(nil, api.TrackRemoveTagMethod, p)
}

// Scrobble scrobbles a track for the authenticated user.
func (t TrackSession) Scrobble(params lastfm.ScrobbleParams) (*lastfm.ScrobbleResult, error) {
	var res lastfm.ScrobbleResult
	return &res, t.session.Post(&res, api.TrackScrobbleMethod, params)
}

// ScrobbleMulti scrobbles multiple tracks for the authenticated user.
func (t TrackSession) ScrobbleMulti(
	params lastfm.ScrobbleMultiParams) (*lastfm.ScrobbleMultiResult, error) {
	type multiScrobbleContainer struct{ lastfm.ScrobbleMultiParams }

	var res lastfm.ScrobbleMultiResult
	p := multiScrobbleContainer{ScrobbleMultiParams: params}
	return &res, t.session.Post(&res, api.TrackScrobbleMethod, p)
}

// Unlove unmarks a track as loved for the authenticated user.
func (t TrackSession) Unlove(artist, track string) error {
	p := lastfm.TrackUnloveParams{Artist: artist, Track: track}
	return t.session.Post(nil, api.TrackUnloveMethod, p)
}

// UpdateNowPlaying updates the now playing track for the authenticated user.
func (t TrackSession) UpdateNowPlaying(
	params lastfm.UpdateNowPlayingParams) (*lastfm.NowPlayingUpdate, error) {

	var res lastfm.NowPlayingUpdate
	return &res, t.session.Post(&res, api.TrackUpdateNowPlayingMethod, params)
}
