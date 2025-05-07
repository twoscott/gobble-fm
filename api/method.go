package api

import "net/url"

// APIMethod represents a Last.fm API method parameter.
type APIMethod string

// String returns the string representation of the APIMethod.
func (m APIMethod) String() string {
	return string(m)
}

// https://www.last.fm/api
const (
	AlbumAddTagsMethod    APIMethod = "album.addTags"
	AlbumGetInfoMethod    APIMethod = "album.getInfo"
	AlbumGetTagsMethod    APIMethod = "album.getTags"
	AlbumGetTopTagsMethod APIMethod = "album.getTopTags"
	AlbumRemoveTagMethod  APIMethod = "album.removeTag"
	AlbumSearchMethod     APIMethod = "album.search"

	ArtistAddTagsMethod       APIMethod = "artist.addTags"
	ArtistGetCorrectionMethod APIMethod = "artist.getCorrection"
	ArtistGetInfoMethod       APIMethod = "artist.getInfo"
	ArtistGetSimilarMethod    APIMethod = "artist.getSimilar"
	ArtistGetTagsMethod       APIMethod = "artist.getTags"
	ArtistGetTopAlbumsMethod  APIMethod = "artist.getTopAlbums"
	ArtistGetTopTagsMethod    APIMethod = "artist.getTopTags"
	ArtistGetTopTracksMethod  APIMethod = "artist.getTopTracks"
	ArtistRemoveTagMethod     APIMethod = "artist.removeTag"
	ArtistSearchMethod        APIMethod = "artist.search"

	AuthGetMobileSessionMethod APIMethod = "auth.getMobileSession"
	AuthGetSessionMethod       APIMethod = "auth.getSession"
	AuthGetTokenMethod         APIMethod = "auth.getToken"

	ChartGetTopArtistsMethod APIMethod = "chart.getTopArtists"
	ChartGetTopTagsMethod    APIMethod = "chart.getTopTags"
	ChartGetTopTracksMethod  APIMethod = "chart.getTopTracks"

	GeoGetTopArtistsMethod APIMethod = "geo.getTopArtists"
	GeoGetTopTracksMethod  APIMethod = "geo.getTopTracks"

	LibraryGetArtistsMethod APIMethod = "library.getArtists"

	TagGetInfoMethod            APIMethod = "tag.getInfo"
	TagGetSimilarMethod         APIMethod = "tag.getSimilar"
	TagGetTopAlbumsMethod       APIMethod = "tag.getTopAlbums"
	TagGetTopArtistsMethod      APIMethod = "tag.getTopArtists"
	TagGetTopTagsMethod         APIMethod = "tag.getTopTags"
	TagGetTopTracksMethod       APIMethod = "tag.getTopTracks"
	TagGetWeeklyChartListMethod APIMethod = "tag.getWeeklyChartList"

	TrackAddTagsMethod          APIMethod = "track.addTags"
	TrackGetCorrectionMethod    APIMethod = "track.getCorrection"
	TrackGetInfoMethod          APIMethod = "track.getInfo"
	TrackGetSimilarMethod       APIMethod = "track.getSimilar"
	TrackGetTagsMethod          APIMethod = "track.getTags"
	TrackGetTopTagsMethod       APIMethod = "track.getTopTags"
	TrackLoveMethod             APIMethod = "track.love"
	TrackRemoveTagMethod        APIMethod = "track.removeTag"
	TrackScrobbleMethod         APIMethod = "track.scrobble"
	TrackSearchMethod           APIMethod = "track.search"
	TrackUnloveMethod           APIMethod = "track.unlove"
	TrackUpdateNowPlayingMethod APIMethod = "track.updateNowPlaying"

	UserGetFriendsMethod           APIMethod = "user.getFriends"
	UserGetInfoMethod              APIMethod = "user.getInfo"
	UserGetLovedTracksMethod       APIMethod = "user.getLovedTracks"
	UserGetPersonalTagsMethod      APIMethod = "user.getPersonalTags"
	UserGetRecentTracksMethod      APIMethod = "user.getRecentTracks"
	UserGetTopAlbumsMethod         APIMethod = "user.getTopAlbums"
	UserGetTopArtistsMethod        APIMethod = "user.getTopArtists"
	UserGetTopTagsMethod           APIMethod = "user.getTopTags"
	UserGetTopTracksMethod         APIMethod = "user.getTopTracks"
	UserGetWeeklyAlbumChartMethod  APIMethod = "user.getWeeklyAlbumChart"
	UserGetWeeklyArtistChartMethod APIMethod = "user.getWeeklyArtistChart"
	UserGetWeeklyChartListMethod   APIMethod = "user.getWeeklyChartList"
	UserGetWeeklyTrackChartMethod  APIMethod = "user.getWeeklyTrackChart"
)

// BuildAPIURL constructs a Last.fm API URL with the specified parameters.
func BuildAPIURL(params url.Values) string {
	return buildURL(Endpoint, params)
}

func buildURL(url string, params url.Values) string {
	return url + "?" + params.Encode()
}
