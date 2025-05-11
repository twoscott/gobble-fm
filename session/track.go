package session

import (
	"github.com/twoscott/gobble-fm/api"
	"github.com/twoscott/gobble-fm/lastfm"
)

type multiScrobbleContainer struct {
	lastfm.ScrobbleMultiParams
}

type Track struct {
	*api.Track
	session *Session
}

// NewTrack creates and returns a new Track API route.
func NewTrack(session *Session) *Track {
	return &Track{Track: api.NewTrack(session.API), session: session}
}

// AddTags adds tags to a track for the authenticated user.
func (t Track) AddTags(artist, track string, tags []string) error {
	p := lastfm.TrackAddTagsParams{Artist: artist, Track: track, Tags: tags}
	return t.session.Post(nil, api.TrackAddTagsMethod, p)
}

// SelfTags returns the tags of a track for the authenticated user by artist and
// track name.
func (t Track) SelfTags(params lastfm.TrackSelfTagsParams) (*lastfm.TrackTags, error) {
	var res lastfm.TrackTags
	return &res, t.session.Get(&res, api.TrackGetTagsMethod, params)
}

// SelfTagsByMBID returns the tags of a track for the authenticated user by
// MBID.
func (t Track) SelfTagsByMBID(params lastfm.TrackSelfTagsMBIDParams) (*lastfm.TrackTags, error) {
	var res lastfm.TrackTags
	return &res, t.session.Get(&res, api.TrackGetTagsMethod, params)
}

// Love marks a track as loved for the authenticated user.
func (t Track) Love(artist, track string) error {
	p := lastfm.TrackLoveParams{Artist: artist, Track: track}
	return t.session.Post(nil, api.TrackLoveMethod, p)
}

// RemoveTag removes a tag from a track for the authenticated user.
func (t Track) RemoveTag(artist, track, tag string) error {
	p := lastfm.TrackRemoveTagParams{Artist: artist, Track: track, Tag: tag}
	return t.session.Post(nil, api.TrackRemoveTagMethod, p)
}

// Scrobble scrobbles a track for the authenticated user.
func (t Track) Scrobble(params lastfm.ScrobbleParams) (*lastfm.ScrobbleResult, error) {
	var res lastfm.ScrobbleResult
	return &res, t.session.Post(&res, api.TrackScrobbleMethod, params)
}

// ScrobbleMulti scrobbles multiple tracks for the authenticated user.
func (t Track) ScrobbleMulti(
	params lastfm.ScrobbleMultiParams) (*lastfm.ScrobbleMultiResult, error) {

	var res lastfm.ScrobbleMultiResult
	p := multiScrobbleContainer{ScrobbleMultiParams: params}
	return &res, t.session.Post(&res, api.TrackScrobbleMethod, p)
}

// Unlove unmarks a track as loved for the authenticated user.
func (t Track) Unlove(artist, track string) error {
	p := lastfm.TrackUnloveParams{Artist: artist, Track: track}
	return t.session.Post(nil, api.TrackUnloveMethod, p)
}

// UpdateNowPlaying updates the now playing track for the authenticated user.
func (t Track) UpdateNowPlaying(
	params lastfm.UpdateNowPlayingParams) (*lastfm.NowPlayingUpdate, error) {

	var res lastfm.NowPlayingUpdate
	return &res, t.session.Post(&res, api.TrackUpdateNowPlayingMethod, params)
}
