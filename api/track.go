package api

import "github.com/twoscott/gobble-fm/lastfm"

type Track struct {
	api *API
}

// NewTrack creates and returns a new Track API route.
func NewTrack(api *API) *Track {
	return &Track{api: api}
}

// Correction returns the track and artist name corrections of a track.
func (t Track) Correction(artist, track string) (*lastfm.TrackCorrection, error) {
	var res lastfm.TrackCorrection
	p := lastfm.TrackCorrectionParams{Artist: artist, Track: track}
	return &res, t.api.Get(&res, TrackGetCorrectionMethod, p)
}

// Info returns the information of a track by artist and track name.
func (t Track) Info(params lastfm.TrackInfoParams) (*lastfm.TrackInfo, error) {
	var res lastfm.TrackInfo
	return &res, t.api.Get(&res, TrackGetInfoMethod, params)
}

// InfoByMBID returns the information of a track by MBID.
func (t Track) InfoByMBID(params lastfm.TrackInfoMBIDParams) (*lastfm.TrackInfo, error) {
	var res lastfm.TrackInfo
	return &res, t.api.Get(&res, TrackGetInfoMethod, params)
}

// UserInfo returns the information of a track for user by artist and track
// name.
func (t Track) UserInfo(params lastfm.TrackUserInfoParams) (*lastfm.TrackUserInfo, error) {
	var res lastfm.TrackUserInfo
	return &res, t.api.Get(&res, TrackGetInfoMethod, params)
}

// UserInfoByMBID returns the information of a track for user by MBID.
func (t Track) UserInfoByMBID(
	params lastfm.TrackUserInfoMBIDParams) (*lastfm.TrackUserInfo, error) {

	var res lastfm.TrackUserInfo
	return &res, t.api.Get(&res, TrackGetInfoMethod, params)
}

// Similar returns the similar tracks of a track by artist and track name.
func (t Track) Similar(params lastfm.TrackSimilarParams) (*lastfm.SimilarTracks, error) {
	var res lastfm.SimilarTracks
	return &res, t.api.Get(&res, TrackGetSimilarMethod, params)
}

// SimilarByMBID returns the similar tracks of a track by MBID.
func (t Track) SimilarByMBID(params lastfm.TrackSimilarMBIDParams) (*lastfm.SimilarTracks, error) {
	var res lastfm.SimilarTracks
	return &res, t.api.Get(&res, TrackGetSimilarMethod, params)
}

// Tags returns the tags of a track by artist and track name.
func (t Track) Tags(params lastfm.TrackTagsParams) (*lastfm.TrackTags, error) {
	var res lastfm.TrackTags
	return &res, t.api.Get(&res, TrackGetTagsMethod, params)
}

// TagsByMBID returns the tags of a track by MBID.
func (t Track) TagsByMBID(params lastfm.TrackTagsMBIDParams) (*lastfm.TrackTags, error) {
	var res lastfm.TrackTags
	return &res, t.api.Get(&res, TrackGetTagsMethod, params)
}

// TopTags returns the top tags of a track by artist and track name.
func (t Track) TopTags(params lastfm.TrackTopTagsParams) (*lastfm.TrackTopTags, error) {
	var res lastfm.TrackTopTags
	return &res, t.api.Get(&res, TrackGetTopTagsMethod, params)
}

// TopTagsByMBID returns the top tags of a track by MBID.
func (t Track) TopTagsByMBID(params lastfm.TrackTopTagsMBIDParams) (*lastfm.TrackTopTags, error) {
	var res lastfm.TrackTopTags
	return &res, t.api.Get(&res, TrackGetTopTagsMethod, params)
}

// Search searches for tracks by track name, and optionally artist name.
func (t Track) Search(params lastfm.TrackSearchParams) (*lastfm.TrackSearchResult, error) {
	var res lastfm.TrackSearchResult
	return &res, t.api.Get(&res, TrackSearchMethod, params)
}
