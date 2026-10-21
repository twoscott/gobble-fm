package client

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/twoscott/gobble-fm/v2/api"
	"github.com/twoscott/gobble-fm/v2/lastfm"
	"github.com/twoscott/gobble-fm/v2/util/option"
)

const (
	testAPIKey     = "testapikey"
	testSecret     = "testsecret"
	testSessionKey = "testsessionkey"

	testAlbum   = "testalbum"
	testArtist  = "testartist"
	testCountry = "testcountry"
	testMBID    = "testmbid"
	testTag     = "testtag"
	testTrack   = "testtrack"
	testUser    = "testuser"
)

var (
	testTags = []string{"tag1", "tag2"}
	testTime = time.Unix(1700000000, 0)
)

type recordedRequest struct {
	httpMethod string
	params     url.Values
}

// mockHTTPClient records the parameters of each request it receives and
// replies to all of them with the same response.
type mockHTTPClient struct {
	t        *testing.T
	status   int
	body     string
	requests []recordedRequest
}

func newMockClient(t *testing.T, status int, body string) *mockHTTPClient {
	return &mockHTTPClient{t: t, status: status, body: body}
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	m.t.Helper()

	params := req.URL.Query()
	if req.Body != nil {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			m.t.Fatal("failed to read request body:", err)
		}
		if params, err = url.ParseQuery(string(body)); err != nil {
			m.t.Fatal("failed to parse request body:", err)
		}
	}

	m.requests = append(m.requests, recordedRequest{httpMethod: req.Method, params: params})

	return &http.Response{
		StatusCode: m.status,
		Body:       io.NopCloser(strings.NewReader(m.body)),
	}, nil
}

// testClients holds one of each client type, all sending requests through the
// same HTTP client.
type testClients struct {
	client  *Client
	key     *KeyClient
	session *SessionClient
}

func newTestClients(t *testing.T, httpClient api.HTTPClient) testClients {
	t.Helper()

	c, err := New(testAPIKey, testSecret)
	if err != nil {
		t.Fatal("failed to create Client:", err)
	}
	c.API.Client = httpClient

	k, err := NewKeyOnly(testAPIKey)
	if err != nil {
		t.Fatal("failed to create KeyClient:", err)
	}
	k.API.Client = httpClient

	// shares c's AuthAPI, so it already uses httpClient
	s, err := c.NewSessionFromKey(testSessionKey)
	if err != nil {
		t.Fatal("failed to create SessionClient:", err)
	}

	return testClients{client: c, key: k, session: s}
}

// errOf discards a route's result so route calls fit in a table.
func errOf[T any](_ T, err error) error {
	return err
}

type authLevel int

const (
	unsigned authLevel = iota
	signed
	session
)

func TestConstructors(t *testing.T) {
	base, err := New(testAPIKey, testSecret)
	if err != nil {
		t.Fatal("failed to create Client:", err)
	}

	var (
		errAPIKeyMissing = api.NewLastFMError(api.ErrAPIKeyMissing, api.APIKeyMissingMessage)
		errSecret        = api.NewLastFMError(api.ErrSecretRequired, api.SecretRequiredMessage)
		errSession       = api.NewLastFMError(api.ErrSessionRequired, api.SessionRequiredMessage)
	)

	cases := []struct {
		name    string
		new     func() (any, error)
		wantErr error
	}{
		{
			name: "New",
			new:  func() (any, error) { return New(testAPIKey, testSecret) },
		},
		{
			name:    "New missing API key",
			new:     func() (any, error) { return New("", testSecret) },
			wantErr: errAPIKeyMissing,
		},
		{
			name:    "New missing secret",
			new:     func() (any, error) { return New(testAPIKey, "") },
			wantErr: errSecret,
		},
		{
			name: "NewWithTimeout",
			new:  func() (any, error) { return NewWithTimeout(testAPIKey, testSecret, 5) },
		},
		{
			name:    "NewWithTimeout missing secret",
			new:     func() (any, error) { return NewWithTimeout(testAPIKey, "", 5) },
			wantErr: errSecret,
		},
		{
			name: "NewKeyOnly",
			new:  func() (any, error) { return NewKeyOnly(testAPIKey) },
		},
		{
			name:    "NewKeyOnly missing API key",
			new:     func() (any, error) { return NewKeyOnly("") },
			wantErr: errAPIKeyMissing,
		},
		{
			name: "NewSession",
			new: func() (any, error) {
				return NewSession(testAPIKey, testSecret, lastfm.Session{Name: testUser, Key: testSessionKey})
			},
		},
		{
			name: "NewSession missing session key",
			new: func() (any, error) {
				return NewSession(testAPIKey, testSecret, lastfm.Session{Name: testUser})
			},
			wantErr: errSession,
		},
		{
			name: "NewSessionFromKey",
			new:  func() (any, error) { return NewSessionFromKey(testAPIKey, testSecret, testSessionKey) },
		},
		{
			name:    "NewSessionFromKey missing secret",
			new:     func() (any, error) { return NewSessionFromKey(testAPIKey, "", testSessionKey) },
			wantErr: errSecret,
		},
		{
			name: "Client.NewSession",
			new:  func() (any, error) { return base.NewSession(lastfm.Session{Key: testSessionKey}) },
		},
		{
			name: "Client.NewSessionFromKey",
			new:  func() (any, error) { return base.NewSessionFromKey(testSessionKey) },
		},
		{
			name:    "Client.NewSessionFromKey missing session key",
			new:     func() (any, error) { return base.NewSessionFromKey("") },
			wantErr: errSession,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := c.new()
			v := reflect.ValueOf(got)

			if c.wantErr != nil {
				if !v.IsNil() {
					t.Errorf("expected nil client, got %+v", got)
				}
				if !errors.Is(err, c.wantErr) {
					t.Errorf("expected error %v, got %v", c.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if v.IsNil() {
				t.Fatal("expected client, got nil")
			}

			// every embedded API and route must be set, or calling it panics
			elem := v.Elem()
			for i := range elem.NumField() {
				if f := elem.Field(i); f.Kind() == reflect.Pointer && f.IsNil() {
					t.Errorf("%s.%s is nil", elem.Type().Name(), elem.Type().Field(i).Name)
				}
			}
		})
	}
}

func TestRoutes(t *testing.T) {
	const (
		get  = http.MethodGet
		post = http.MethodPost
	)

	cases := []struct {
		name string
		call func(c testClients) error

		httpMethod string
		method     api.APIMethod
		level      authLevel
		// params are the route-specific parameters. The request must contain
		// exactly these, plus the method, API key, and any signature or session
		// key its auth level requires.
		params map[string]string
	}{
		// Album
		{
			name: "Client.Album.Info",
			call: func(c testClients) error {
				return errOf(c.client.Album.Info(lastfm.AlbumInfoParams{Artist: testArtist, Album: testAlbum}))
			},
			httpMethod: get, method: api.AlbumGetInfoMethod, level: unsigned,
			params: map[string]string{"artist": testArtist, "album": testAlbum},
		},
		{
			name: "Client.Album.InfoByMBID",
			call: func(c testClients) error {
				return errOf(c.client.Album.InfoByMBID(lastfm.AlbumInfoMBIDParams{MBID: testMBID}))
			},
			httpMethod: get, method: api.AlbumGetInfoMethod, level: unsigned,
			params: map[string]string{"mbid": testMBID},
		},
		{
			name: "Client.Album.UserInfo",
			call: func(c testClients) error {
				return errOf(c.client.Album.UserInfo(lastfm.AlbumUserInfoParams{
					Artist: testArtist, Album: testAlbum, User: testUser,
				}))
			},
			httpMethod: get, method: api.AlbumGetInfoMethod, level: unsigned,
			params: map[string]string{"artist": testArtist, "album": testAlbum, "username": testUser},
		},
		{
			name: "Client.Album.UserInfoByMBID",
			call: func(c testClients) error {
				return errOf(c.client.Album.UserInfoByMBID(lastfm.AlbumUserInfoMBIDParams{
					MBID: testMBID, User: testUser,
				}))
			},
			httpMethod: get, method: api.AlbumGetInfoMethod, level: unsigned,
			params: map[string]string{"mbid": testMBID, "username": testUser},
		},
		{
			name: "Client.Album.UserTags",
			call: func(c testClients) error {
				return errOf(c.client.Album.UserTags(lastfm.AlbumTagsParams{
					Artist: testArtist, Album: testAlbum, User: testUser,
				}))
			},
			httpMethod: get, method: api.AlbumGetTagsMethod, level: unsigned,
			params: map[string]string{"artist": testArtist, "album": testAlbum, "username": testUser},
		},
		{
			name: "Client.Album.UserTagsByMBID",
			call: func(c testClients) error {
				return errOf(c.client.Album.UserTagsByMBID(lastfm.AlbumTagsMBIDParams{
					MBID: testMBID, User: testUser,
				}))
			},
			httpMethod: get, method: api.AlbumGetTagsMethod, level: unsigned,
			params: map[string]string{"mbid": testMBID, "username": testUser},
		},
		{
			name: "Client.Album.TopTags",
			call: func(c testClients) error {
				return errOf(c.client.Album.TopTags(lastfm.AlbumTopTagsParams{Artist: testArtist, Album: testAlbum}))
			},
			httpMethod: get, method: api.AlbumGetTopTagsMethod, level: unsigned,
			params: map[string]string{"artist": testArtist, "album": testAlbum},
		},
		{
			name: "Client.Album.TopTagsByMBID",
			call: func(c testClients) error {
				return errOf(c.client.Album.TopTagsByMBID(lastfm.AlbumTopTagsMBIDParams{MBID: testMBID}))
			},
			httpMethod: get, method: api.AlbumGetTopTagsMethod, level: unsigned,
			params: map[string]string{"mbid": testMBID},
		},
		{
			name: "Client.Album.Search",
			call: func(c testClients) error {
				return errOf(c.client.Album.Search(lastfm.AlbumSearchParams{Album: testAlbum, Limit: 5, Page: 2}))
			},
			httpMethod: get, method: api.AlbumSearchMethod, level: unsigned,
			params: map[string]string{"album": testAlbum, "limit": "5", "page": "2"},
		},
		{
			name: "SessionClient.Album.AddTags",
			call: func(c testClients) error {
				return c.session.Album.AddTags(testArtist, testAlbum, testTags)
			},
			httpMethod: post, method: api.AlbumAddTagsMethod, level: session,
			params: map[string]string{"artist": testArtist, "album": testAlbum, "tags": "tag1,tag2"},
		},
		{
			name: "SessionClient.Album.SelfTags",
			call: func(c testClients) error {
				return errOf(c.session.Album.SelfTags(lastfm.AlbumSelfTagsParams{Artist: testArtist, Album: testAlbum}))
			},
			httpMethod: get, method: api.AlbumGetTagsMethod, level: session,
			params: map[string]string{"artist": testArtist, "album": testAlbum},
		},
		{
			name: "SessionClient.Album.SelfTagsByMBID",
			call: func(c testClients) error {
				return errOf(c.session.Album.SelfTagsByMBID(lastfm.AlbumSelfTagsMBIDParams{MBID: testMBID}))
			},
			httpMethod: get, method: api.AlbumGetTagsMethod, level: session,
			params: map[string]string{"mbid": testMBID},
		},
		{
			name: "SessionClient.Album.RemoveTag",
			call: func(c testClients) error {
				return c.session.Album.RemoveTag(testArtist, testAlbum, testTag)
			},
			httpMethod: post, method: api.AlbumRemoveTagMethod, level: session,
			params: map[string]string{"artist": testArtist, "album": testAlbum, "tag": testTag},
		},

		// Artist
		{
			name: "Client.Artist.Correction",
			call: func(c testClients) error {
				return errOf(c.client.Artist.Correction(testArtist))
			},
			httpMethod: get, method: api.ArtistGetCorrectionMethod, level: unsigned,
			params: map[string]string{"artist": testArtist},
		},
		{
			name: "Client.Artist.Info",
			call: func(c testClients) error {
				return errOf(c.client.Artist.Info(lastfm.ArtistInfoParams{
					Artist: testArtist, AutoCorrect: option.True, Language: "en",
				}))
			},
			httpMethod: get, method: api.ArtistGetInfoMethod, level: unsigned,
			params: map[string]string{"artist": testArtist, "autocorrect": "1", "lang": "en"},
		},
		{
			name: "Client.Artist.InfoByMBID",
			call: func(c testClients) error {
				return errOf(c.client.Artist.InfoByMBID(lastfm.ArtistInfoMBIDParams{MBID: testMBID}))
			},
			httpMethod: get, method: api.ArtistGetInfoMethod, level: unsigned,
			params: map[string]string{"mbid": testMBID},
		},
		{
			name: "Client.Artist.UserInfo",
			call: func(c testClients) error {
				return errOf(c.client.Artist.UserInfo(lastfm.ArtistUserInfoParams{Artist: testArtist, User: testUser}))
			},
			httpMethod: get, method: api.ArtistGetInfoMethod, level: unsigned,
			params: map[string]string{"artist": testArtist, "username": testUser},
		},
		{
			name: "Client.Artist.UserInfoByMBID",
			call: func(c testClients) error {
				return errOf(c.client.Artist.UserInfoByMBID(lastfm.ArtistUserInfoMBIDParams{
					MBID: testMBID, User: testUser,
				}))
			},
			httpMethod: get, method: api.ArtistGetInfoMethod, level: unsigned,
			params: map[string]string{"mbid": testMBID, "username": testUser},
		},
		{
			name: "Client.Artist.Similar",
			call: func(c testClients) error {
				return errOf(c.client.Artist.Similar(lastfm.ArtistSimilarParams{Artist: testArtist}))
			},
			httpMethod: get, method: api.ArtistGetSimilarMethod, level: unsigned,
			params: map[string]string{"artist": testArtist},
		},
		{
			name: "Client.Artist.SimilarByMBID",
			call: func(c testClients) error {
				return errOf(c.client.Artist.SimilarByMBID(lastfm.ArtistSimilarMBIDParams{MBID: testMBID}))
			},
			httpMethod: get, method: api.ArtistGetSimilarMethod, level: unsigned,
			params: map[string]string{"mbid": testMBID},
		},
		{
			name: "Client.Artist.UserTags",
			call: func(c testClients) error {
				return errOf(c.client.Artist.UserTags(lastfm.ArtistTagsParams{Artist: testArtist, User: testUser}))
			},
			httpMethod: get, method: api.ArtistGetTagsMethod, level: unsigned,
			params: map[string]string{"artist": testArtist, "username": testUser},
		},
		{
			name: "Client.Artist.UserTagsByMBID",
			call: func(c testClients) error {
				return errOf(c.client.Artist.UserTagsByMBID(lastfm.ArtistTagsMBIDParams{MBID: testMBID, User: testUser}))
			},
			httpMethod: get, method: api.ArtistGetTagsMethod, level: unsigned,
			params: map[string]string{"mbid": testMBID, "username": testUser},
		},
		{
			name: "Client.Artist.TopAlbums",
			call: func(c testClients) error {
				return errOf(c.client.Artist.TopAlbums(lastfm.ArtistTopAlbumsParams{Artist: testArtist}))
			},
			httpMethod: get, method: api.ArtistGetTopAlbumsMethod, level: unsigned,
			params: map[string]string{"artist": testArtist},
		},
		{
			name: "Client.Artist.TopAlbumsByMBID",
			call: func(c testClients) error {
				return errOf(c.client.Artist.TopAlbumsByMBID(lastfm.ArtistTopAlbumsMBIDParams{MBID: testMBID}))
			},
			httpMethod: get, method: api.ArtistGetTopAlbumsMethod, level: unsigned,
			params: map[string]string{"mbid": testMBID},
		},
		{
			name: "Client.Artist.TopTags",
			call: func(c testClients) error {
				return errOf(c.client.Artist.TopTags(lastfm.ArtistTopTagsParams{Artist: testArtist}))
			},
			httpMethod: get, method: api.ArtistGetTopTagsMethod, level: unsigned,
			params: map[string]string{"artist": testArtist},
		},
		{
			name: "Client.Artist.TopTagsByMBID",
			call: func(c testClients) error {
				return errOf(c.client.Artist.TopTagsByMBID(lastfm.ArtistTopTagsMBIDParams{MBID: testMBID}))
			},
			httpMethod: get, method: api.ArtistGetTopTagsMethod, level: unsigned,
			params: map[string]string{"mbid": testMBID},
		},
		{
			name: "Client.Artist.TopTracks",
			call: func(c testClients) error {
				return errOf(c.client.Artist.TopTracks(lastfm.ArtistTopTracksParams{Artist: testArtist}))
			},
			httpMethod: get, method: api.ArtistGetTopTracksMethod, level: unsigned,
			params: map[string]string{"artist": testArtist},
		},
		{
			name: "Client.Artist.TopTracksByMBID",
			call: func(c testClients) error {
				return errOf(c.client.Artist.TopTracksByMBID(lastfm.ArtistTopTracksMBIDParams{MBID: testMBID}))
			},
			httpMethod: get, method: api.ArtistGetTopTracksMethod, level: unsigned,
			params: map[string]string{"mbid": testMBID},
		},
		{
			name: "Client.Artist.Search",
			call: func(c testClients) error {
				return errOf(c.client.Artist.Search(lastfm.ArtistSearchParams{Artist: testArtist}))
			},
			httpMethod: get, method: api.ArtistSearchMethod, level: unsigned,
			params: map[string]string{"artist": testArtist},
		},
		{
			name: "SessionClient.Artist.AddTags",
			call: func(c testClients) error {
				return c.session.Artist.AddTags(testArtist, testTags)
			},
			httpMethod: post, method: api.ArtistAddTagsMethod, level: session,
			params: map[string]string{"artist": testArtist, "tags": "tag1,tag2"},
		},
		{
			name: "SessionClient.Artist.SelfTags",
			call: func(c testClients) error {
				return errOf(c.session.Artist.SelfTags(lastfm.ArtistSelfTagsParams{Artist: testArtist}))
			},
			httpMethod: get, method: api.ArtistGetTagsMethod, level: session,
			params: map[string]string{"artist": testArtist},
		},
		{
			name: "SessionClient.Artist.SelfTagsByMBID",
			call: func(c testClients) error {
				return errOf(c.session.Artist.SelfTagsByMBID(lastfm.ArtistSelfTagsMBIDParams{MBID: testMBID}))
			},
			httpMethod: get, method: api.ArtistGetTagsMethod, level: session,
			params: map[string]string{"mbid": testMBID},
		},
		{
			name: "SessionClient.Artist.RemoveTag",
			call: func(c testClients) error {
				return c.session.Artist.RemoveTag(testArtist, testTag)
			},
			httpMethod: post, method: api.ArtistRemoveTagMethod, level: session,
			params: map[string]string{"artist": testArtist, "tag": testTag},
		},

		// Auth
		{
			name: "Client.Auth.Token",
			call: func(c testClients) error {
				return errOf(c.client.Auth.Token())
			},
			httpMethod: get, method: api.AuthGetTokenMethod, level: signed,
			params: map[string]string{},
		},
		{
			name: "Client.Auth.Session",
			call: func(c testClients) error {
				return errOf(c.client.Auth.Session("testtoken"))
			},
			httpMethod: post, method: api.AuthGetSessionMethod, level: signed,
			params: map[string]string{"token": "testtoken"},
		},
		{
			name: "Client.Auth.MobileSession",
			call: func(c testClients) error {
				return errOf(c.client.Auth.MobileSession(testUser, "testpassword"))
			},
			httpMethod: post, method: api.AuthGetMobileSessionMethod, level: signed,
			params: map[string]string{"username": testUser, "password": "testpassword"},
		},

		// Chart
		{
			name: "Client.Chart.TopArtists",
			call: func(c testClients) error {
				return errOf(c.client.Chart.TopArtists())
			},
			httpMethod: get, method: api.ChartGetTopArtistsMethod, level: unsigned,
			params: map[string]string{},
		},
		{
			name: "Client.Chart.TopArtistsLimit",
			call: func(c testClients) error {
				return errOf(c.client.Chart.TopArtistsLimit(&lastfm.ChartTopArtistsParams{Limit: 5, Page: 2}))
			},
			httpMethod: get, method: api.ChartGetTopArtistsMethod, level: unsigned,
			params: map[string]string{"limit": "5", "page": "2"},
		},
		{
			name: "Client.Chart.TopTags",
			call: func(c testClients) error {
				return errOf(c.client.Chart.TopTags())
			},
			httpMethod: get, method: api.ChartGetTopTagsMethod, level: unsigned,
			params: map[string]string{},
		},
		{
			name: "Client.Chart.TopTagsLimit",
			call: func(c testClients) error {
				return errOf(c.client.Chart.TopTagsLimit(&lastfm.ChartTopTagsParams{Limit: 5, Page: 2}))
			},
			httpMethod: get, method: api.ChartGetTopTagsMethod, level: unsigned,
			params: map[string]string{"limit": "5", "page": "2"},
		},
		{
			name: "Client.Chart.TopTracks",
			call: func(c testClients) error {
				return errOf(c.client.Chart.TopTracks())
			},
			httpMethod: get, method: api.ChartGetTopTracksMethod, level: unsigned,
			params: map[string]string{},
		},
		{
			name: "Client.Chart.TopTracksLimit",
			call: func(c testClients) error {
				return errOf(c.client.Chart.TopTracksLimit(&lastfm.ChartTopTracksParams{Limit: 5, Page: 2}))
			},
			httpMethod: get, method: api.ChartGetTopTracksMethod, level: unsigned,
			params: map[string]string{"limit": "5", "page": "2"},
		},

		// Geo
		{
			name: "Client.Geo.TopArtists",
			call: func(c testClients) error {
				return errOf(c.client.Geo.TopArtists(lastfm.GeoTopArtistsParams{Country: testCountry}))
			},
			httpMethod: get, method: api.GeoGetTopArtistsMethod, level: unsigned,
			params: map[string]string{"country": testCountry},
		},
		{
			name: "Client.Geo.TopTracks",
			call: func(c testClients) error {
				return errOf(c.client.Geo.TopTracks(lastfm.GeoTopTracksParams{
					Country: testCountry, Location: "testlocation",
				}))
			},
			httpMethod: get, method: api.GeoGetTopTracksMethod, level: unsigned,
			params: map[string]string{"country": testCountry, "location": "testlocation"},
		},

		// Library
		{
			name: "Client.Library.Artists",
			call: func(c testClients) error {
				return errOf(c.client.Library.Artists(lastfm.LibraryArtistsParams{User: testUser}))
			},
			httpMethod: get, method: api.LibraryGetArtistsMethod, level: unsigned,
			params: map[string]string{"user": testUser},
		},

		// Tag
		{
			name: "Client.Tag.Info",
			call: func(c testClients) error {
				return errOf(c.client.Tag.Info(lastfm.TagInfoParams{Tag: testTag}))
			},
			httpMethod: get, method: api.TagGetInfoMethod, level: unsigned,
			params: map[string]string{"tag": testTag},
		},
		{
			name: "Client.Tag.Similar",
			call: func(c testClients) error {
				return errOf(c.client.Tag.Similar(testTag))
			},
			httpMethod: get, method: api.TagGetSimilarMethod, level: unsigned,
			params: map[string]string{"tag": testTag},
		},
		{
			name: "Client.Tag.TopAlbums",
			call: func(c testClients) error {
				return errOf(c.client.Tag.TopAlbums(lastfm.TagTopAlbumsParams{Tag: testTag}))
			},
			httpMethod: get, method: api.TagGetTopAlbumsMethod, level: unsigned,
			params: map[string]string{"tag": testTag},
		},
		{
			name: "Client.Tag.TopArtists",
			call: func(c testClients) error {
				return errOf(c.client.Tag.TopArtists(lastfm.TagTopArtistsParams{Tag: testTag}))
			},
			httpMethod: get, method: api.TagGetTopArtistsMethod, level: unsigned,
			params: map[string]string{"tag": testTag},
		},
		{
			name: "Client.Tag.TopTags",
			call: func(c testClients) error {
				return errOf(c.client.Tag.TopTags())
			},
			httpMethod: get, method: api.TagGetTopTagsMethod, level: unsigned,
			params: map[string]string{},
		},
		{
			name: "Client.Tag.TopTagsLimit",
			call: func(c testClients) error {
				return errOf(c.client.Tag.TopTagsLimit(&lastfm.TagTopTagsParams{Limit: 5, Offset: 10}))
			},
			httpMethod: get, method: api.TagGetTopTagsMethod, level: unsigned,
			params: map[string]string{"num_res": "5", "offset": "10"},
		},
		{
			name: "Client.Tag.TopTracks",
			call: func(c testClients) error {
				return errOf(c.client.Tag.TopTracks(lastfm.TagTopTracksParams{Tag: testTag}))
			},
			httpMethod: get, method: api.TagGetTopTracksMethod, level: unsigned,
			params: map[string]string{"tag": testTag},
		},
		{
			name: "Client.Tag.WeeklyChartList",
			call: func(c testClients) error {
				return errOf(c.client.Tag.WeeklyChartList(testTag))
			},
			httpMethod: get, method: api.TagGetWeeklyChartListMethod, level: unsigned,
			params: map[string]string{"tag": testTag},
		},

		// Track
		{
			name: "Client.Track.Correction",
			call: func(c testClients) error {
				return errOf(c.client.Track.Correction(testArtist, testTrack))
			},
			httpMethod: get, method: api.TrackGetCorrectionMethod, level: unsigned,
			params: map[string]string{"artist": testArtist, "track": testTrack},
		},
		{
			name: "Client.Track.Info",
			call: func(c testClients) error {
				return errOf(c.client.Track.Info(lastfm.TrackInfoParams{Artist: testArtist, Track: testTrack}))
			},
			httpMethod: get, method: api.TrackGetInfoMethod, level: unsigned,
			params: map[string]string{"artist": testArtist, "track": testTrack},
		},
		{
			name: "Client.Track.InfoByMBID",
			call: func(c testClients) error {
				return errOf(c.client.Track.InfoByMBID(lastfm.TrackInfoMBIDParams{MBID: testMBID}))
			},
			httpMethod: get, method: api.TrackGetInfoMethod, level: unsigned,
			params: map[string]string{"mbid": testMBID},
		},
		{
			name: "Client.Track.UserInfo",
			call: func(c testClients) error {
				return errOf(c.client.Track.UserInfo(lastfm.TrackUserInfoParams{
					Artist: testArtist, Track: testTrack, User: testUser,
				}))
			},
			httpMethod: get, method: api.TrackGetInfoMethod, level: unsigned,
			params: map[string]string{"artist": testArtist, "track": testTrack, "username": testUser},
		},
		{
			name: "Client.Track.UserInfoByMBID",
			call: func(c testClients) error {
				return errOf(c.client.Track.UserInfoByMBID(lastfm.TrackUserInfoMBIDParams{
					MBID: testMBID, User: testUser,
				}))
			},
			httpMethod: get, method: api.TrackGetInfoMethod, level: unsigned,
			params: map[string]string{"mbid": testMBID, "username": testUser},
		},
		{
			name: "Client.Track.Similar",
			call: func(c testClients) error {
				return errOf(c.client.Track.Similar(lastfm.TrackSimilarParams{Artist: testArtist, Track: testTrack}))
			},
			httpMethod: get, method: api.TrackGetSimilarMethod, level: unsigned,
			params: map[string]string{"artist": testArtist, "track": testTrack},
		},
		{
			name: "Client.Track.SimilarByMBID",
			call: func(c testClients) error {
				return errOf(c.client.Track.SimilarByMBID(lastfm.TrackSimilarMBIDParams{MBID: testMBID}))
			},
			httpMethod: get, method: api.TrackGetSimilarMethod, level: unsigned,
			params: map[string]string{"mbid": testMBID},
		},
		{
			name: "Client.Track.Tags",
			call: func(c testClients) error {
				return errOf(c.client.Track.Tags(lastfm.TrackTagsParams{
					Artist: testArtist, Track: testTrack, User: testUser,
				}))
			},
			httpMethod: get, method: api.TrackGetTagsMethod, level: unsigned,
			params: map[string]string{"artist": testArtist, "track": testTrack, "username": testUser},
		},
		{
			name: "Client.Track.TagsByMBID",
			call: func(c testClients) error {
				return errOf(c.client.Track.TagsByMBID(lastfm.TrackTagsMBIDParams{MBID: testMBID, User: testUser}))
			},
			httpMethod: get, method: api.TrackGetTagsMethod, level: unsigned,
			params: map[string]string{"mbid": testMBID, "username": testUser},
		},
		{
			name: "Client.Track.TopTags",
			call: func(c testClients) error {
				return errOf(c.client.Track.TopTags(lastfm.TrackTopTagsParams{Artist: testArtist, Track: testTrack}))
			},
			httpMethod: get, method: api.TrackGetTopTagsMethod, level: unsigned,
			params: map[string]string{"artist": testArtist, "track": testTrack},
		},
		{
			name: "Client.Track.TopTagsByMBID",
			call: func(c testClients) error {
				return errOf(c.client.Track.TopTagsByMBID(lastfm.TrackTopTagsMBIDParams{MBID: testMBID}))
			},
			httpMethod: get, method: api.TrackGetTopTagsMethod, level: unsigned,
			params: map[string]string{"mbid": testMBID},
		},
		{
			name: "Client.Track.Search",
			call: func(c testClients) error {
				return errOf(c.client.Track.Search(lastfm.TrackSearchParams{Track: testTrack, Artist: testArtist}))
			},
			httpMethod: get, method: api.TrackSearchMethod, level: unsigned,
			params: map[string]string{"track": testTrack, "artist": testArtist},
		},
		{
			name: "SessionClient.Track.AddTags",
			call: func(c testClients) error {
				return c.session.Track.AddTags(testArtist, testTrack, testTags)
			},
			httpMethod: post, method: api.TrackAddTagsMethod, level: session,
			params: map[string]string{"artist": testArtist, "track": testTrack, "tags": "tag1,tag2"},
		},
		{
			name: "SessionClient.Track.SelfTags",
			call: func(c testClients) error {
				return errOf(c.session.Track.SelfTags(lastfm.TrackSelfTagsParams{Artist: testArtist, Track: testTrack}))
			},
			httpMethod: get, method: api.TrackGetTagsMethod, level: session,
			params: map[string]string{"artist": testArtist, "track": testTrack},
		},
		{
			name: "SessionClient.Track.SelfTagsByMBID",
			call: func(c testClients) error {
				return errOf(c.session.Track.SelfTagsByMBID(lastfm.TrackSelfTagsMBIDParams{MBID: testMBID}))
			},
			httpMethod: get, method: api.TrackGetTagsMethod, level: session,
			params: map[string]string{"mbid": testMBID},
		},
		{
			name: "SessionClient.Track.Love",
			call: func(c testClients) error {
				return c.session.Track.Love(testArtist, testTrack)
			},
			httpMethod: post, method: api.TrackLoveMethod, level: session,
			params: map[string]string{"artist": testArtist, "track": testTrack},
		},
		{
			name: "SessionClient.Track.RemoveTag",
			call: func(c testClients) error {
				return c.session.Track.RemoveTag(testArtist, testTrack, testTag)
			},
			httpMethod: post, method: api.TrackRemoveTagMethod, level: session,
			params: map[string]string{"artist": testArtist, "track": testTrack, "tag": testTag},
		},
		{
			name: "SessionClient.Track.Scrobble",
			call: func(c testClients) error {
				return errOf(c.session.Track.Scrobble(lastfm.ScrobbleParams{
					Artist:   testArtist,
					Track:    testTrack,
					Time:     testTime,
					Duration: lastfm.DurationMinSec(4, 2),
				}))
			},
			httpMethod: post, method: api.TrackScrobbleMethod, level: session,
			params: map[string]string{
				"artist":    testArtist,
				"track":     testTrack,
				"timestamp": "1700000000",
				"duration":  "242",
			},
		},
		{
			name: "SessionClient.Track.ScrobbleMulti",
			call: func(c testClients) error {
				return errOf(c.session.Track.ScrobbleMulti(lastfm.ScrobbleMultiParams{
					{Artist: testArtist, Track: testTrack, Time: testTime},
					{Artist: "testartist2", Track: "testtrack2", Time: testTime.Add(-5 * time.Minute)},
				}))
			},
			httpMethod: post, method: api.TrackScrobbleMethod, level: session,
			params: map[string]string{
				"artist[0]":    testArtist,
				"track[0]":     testTrack,
				"timestamp[0]": "1700000000",
				"artist[1]":    "testartist2",
				"track[1]":     "testtrack2",
				"timestamp[1]": "1699999700",
			},
		},
		{
			name: "SessionClient.Track.Unlove",
			call: func(c testClients) error {
				return c.session.Track.Unlove(testArtist, testTrack)
			},
			httpMethod: post, method: api.TrackUnloveMethod, level: session,
			params: map[string]string{"artist": testArtist, "track": testTrack},
		},
		{
			name: "SessionClient.Track.UpdateNowPlaying",
			call: func(c testClients) error {
				return errOf(c.session.Track.UpdateNowPlaying(lastfm.UpdateNowPlayingParams{
					Artist: testArtist, Track: testTrack,
				}))
			},
			httpMethod: post, method: api.TrackUpdateNowPlayingMethod, level: session,
			params: map[string]string{"artist": testArtist, "track": testTrack},
		},

		// User
		{
			name: "Client.User.Friends",
			call: func(c testClients) error {
				return errOf(c.client.User.Friends(lastfm.FriendsParams{User: testUser}))
			},
			httpMethod: get, method: api.UserGetFriendsMethod, level: unsigned,
			params: map[string]string{"user": testUser},
		},
		{
			name: "Client.User.Info",
			call: func(c testClients) error {
				return errOf(c.client.User.Info(testUser))
			},
			httpMethod: get, method: api.UserGetInfoMethod, level: unsigned,
			params: map[string]string{"user": testUser},
		},
		{
			name: "Client.User.LovedTracks",
			call: func(c testClients) error {
				return errOf(c.client.User.LovedTracks(lastfm.LovedTracksParams{User: testUser}))
			},
			httpMethod: get, method: api.UserGetLovedTracksMethod, level: unsigned,
			params: map[string]string{"user": testUser},
		},
		{
			name: "Client.User.RecentTrack",
			call: func(c testClients) error {
				return errOf(c.client.User.RecentTrack(testUser))
			},
			httpMethod: get, method: api.UserGetRecentTracksMethod, level: unsigned,
			params: map[string]string{"user": testUser, "limit": "1"},
		},
		{
			name: "Client.User.RecentTracks",
			call: func(c testClients) error {
				return errOf(c.client.User.RecentTracks(lastfm.RecentTracksParams{User: testUser, From: testTime}))
			},
			httpMethod: get, method: api.UserGetRecentTracksMethod, level: unsigned,
			params: map[string]string{"user": testUser, "from": "1700000000"},
		},
		{
			name: "Client.User.RecentTrackExtended",
			call: func(c testClients) error {
				return errOf(c.client.User.RecentTrackExtended(testUser))
			},
			httpMethod: get, method: api.UserGetRecentTracksMethod, level: unsigned,
			params: map[string]string{"user": testUser, "limit": "1", "extended": "1"},
		},
		{
			name: "Client.User.RecentTracksExtended",
			call: func(c testClients) error {
				return errOf(c.client.User.RecentTracksExtended(lastfm.RecentTracksParams{User: testUser}))
			},
			httpMethod: get, method: api.UserGetRecentTracksMethod, level: unsigned,
			params: map[string]string{"user": testUser, "extended": "1"},
		},
		{
			name: "Client.User.TaggedAlbums",
			call: func(c testClients) error {
				return errOf(c.client.User.TaggedAlbums(lastfm.UserTagsParams{User: testUser, Tag: testTag}))
			},
			httpMethod: get, method: api.UserGetPersonalTagsMethod, level: unsigned,
			params: map[string]string{"user": testUser, "tag": testTag, "taggingtype": "album"},
		},
		{
			name: "Client.User.TaggedArtists",
			call: func(c testClients) error {
				return errOf(c.client.User.TaggedArtists(lastfm.UserTagsParams{User: testUser, Tag: testTag}))
			},
			httpMethod: get, method: api.UserGetPersonalTagsMethod, level: unsigned,
			params: map[string]string{"user": testUser, "tag": testTag, "taggingtype": "artist"},
		},
		{
			name: "Client.User.TaggedTracks",
			call: func(c testClients) error {
				return errOf(c.client.User.TaggedTracks(lastfm.UserTagsParams{User: testUser, Tag: testTag}))
			},
			httpMethod: get, method: api.UserGetPersonalTagsMethod, level: unsigned,
			params: map[string]string{"user": testUser, "tag": testTag, "taggingtype": "track"},
		},
		{
			name: "Client.User.TopAlbums",
			call: func(c testClients) error {
				return errOf(c.client.User.TopAlbums(lastfm.UserTopAlbumsParams{User: testUser, Period: lastfm.PeriodMonth}))
			},
			httpMethod: get, method: api.UserGetTopAlbumsMethod, level: unsigned,
			params: map[string]string{"user": testUser, "period": "1month"},
		},
		{
			name: "Client.User.TopArtists",
			call: func(c testClients) error {
				return errOf(c.client.User.TopArtists(lastfm.UserTopArtistsParams{User: testUser}))
			},
			httpMethod: get, method: api.UserGetTopArtistsMethod, level: unsigned,
			params: map[string]string{"user": testUser},
		},
		{
			name: "Client.User.TopTags",
			call: func(c testClients) error {
				return errOf(c.client.User.TopTags(lastfm.UserTopTagsParams{User: testUser}))
			},
			httpMethod: get, method: api.UserGetTopTagsMethod, level: unsigned,
			params: map[string]string{"user": testUser},
		},
		{
			name: "Client.User.TopTracks",
			call: func(c testClients) error {
				return errOf(c.client.User.TopTracks(lastfm.UserTopTracksParams{User: testUser}))
			},
			httpMethod: get, method: api.UserGetTopTracksMethod, level: unsigned,
			params: map[string]string{"user": testUser},
		},
		{
			name: "Client.User.WeeklyAlbumChart",
			call: func(c testClients) error {
				return errOf(c.client.User.WeeklyAlbumChart(lastfm.WeeklyAlbumChartParams{User: testUser}))
			},
			httpMethod: get, method: api.UserGetWeeklyAlbumChartMethod, level: unsigned,
			params: map[string]string{"user": testUser},
		},
		{
			name: "Client.User.WeeklyArtistChart",
			call: func(c testClients) error {
				return errOf(c.client.User.WeeklyArtistChart(lastfm.WeeklyArtistChartParams{User: testUser}))
			},
			httpMethod: get, method: api.UserGetWeeklyArtistChartMethod, level: unsigned,
			params: map[string]string{"user": testUser},
		},
		{
			name: "Client.User.WeeklyChartList",
			call: func(c testClients) error {
				return errOf(c.client.User.WeeklyChartList(testUser))
			},
			httpMethod: get, method: api.UserGetWeeklyChartListMethod, level: unsigned,
			params: map[string]string{"user": testUser},
		},
		{
			name: "Client.User.WeeklyTrackChart",
			call: func(c testClients) error {
				return errOf(c.client.User.WeeklyTrackChart(lastfm.WeeklyTrackChartParams{User: testUser}))
			},
			httpMethod: get, method: api.UserGetWeeklyTrackChartMethod, level: unsigned,
			params: map[string]string{"user": testUser},
		},
		{
			name: "SessionClient.User.SelfInfo",
			call: func(c testClients) error {
				return errOf(c.session.User.SelfInfo())
			},
			httpMethod: get, method: api.UserGetInfoMethod, level: session,
			params: map[string]string{},
		},

		// Unsigned routes reached through the other client types
		{
			name: "KeyClient.Album.Info",
			call: func(c testClients) error {
				return errOf(c.key.Album.Info(lastfm.AlbumInfoParams{Artist: testArtist, Album: testAlbum}))
			},
			httpMethod: get, method: api.AlbumGetInfoMethod, level: unsigned,
			params: map[string]string{"artist": testArtist, "album": testAlbum},
		},
		{
			name: "KeyClient.Artist.Info",
			call: func(c testClients) error {
				return errOf(c.key.Artist.Info(lastfm.ArtistInfoParams{Artist: testArtist}))
			},
			httpMethod: get, method: api.ArtistGetInfoMethod, level: unsigned,
			params: map[string]string{"artist": testArtist},
		},
		{
			name: "KeyClient.Chart.TopTags",
			call: func(c testClients) error {
				return errOf(c.key.Chart.TopTags())
			},
			httpMethod: get, method: api.ChartGetTopTagsMethod, level: unsigned,
			params: map[string]string{},
		},
		{
			name: "KeyClient.Geo.TopArtists",
			call: func(c testClients) error {
				return errOf(c.key.Geo.TopArtists(lastfm.GeoTopArtistsParams{Country: testCountry}))
			},
			httpMethod: get, method: api.GeoGetTopArtistsMethod, level: unsigned,
			params: map[string]string{"country": testCountry},
		},
		{
			name: "KeyClient.Library.Artists",
			call: func(c testClients) error {
				return errOf(c.key.Library.Artists(lastfm.LibraryArtistsParams{User: testUser}))
			},
			httpMethod: get, method: api.LibraryGetArtistsMethod, level: unsigned,
			params: map[string]string{"user": testUser},
		},
		{
			name: "KeyClient.Tag.Info",
			call: func(c testClients) error {
				return errOf(c.key.Tag.Info(lastfm.TagInfoParams{Tag: testTag}))
			},
			httpMethod: get, method: api.TagGetInfoMethod, level: unsigned,
			params: map[string]string{"tag": testTag},
		},
		{
			name: "KeyClient.Track.Info",
			call: func(c testClients) error {
				return errOf(c.key.Track.Info(lastfm.TrackInfoParams{Artist: testArtist, Track: testTrack}))
			},
			httpMethod: get, method: api.TrackGetInfoMethod, level: unsigned,
			params: map[string]string{"artist": testArtist, "track": testTrack},
		},
		{
			name: "KeyClient.User.Info",
			call: func(c testClients) error {
				return errOf(c.key.User.Info(testUser))
			},
			httpMethod: get, method: api.UserGetInfoMethod, level: unsigned,
			params: map[string]string{"user": testUser},
		},
		{
			name: "SessionClient.Album.Info",
			call: func(c testClients) error {
				return errOf(c.session.Album.Info(lastfm.AlbumInfoParams{Artist: testArtist, Album: testAlbum}))
			},
			httpMethod: get, method: api.AlbumGetInfoMethod, level: unsigned,
			params: map[string]string{"artist": testArtist, "album": testAlbum},
		},
		{
			name: "SessionClient.Artist.Info",
			call: func(c testClients) error {
				return errOf(c.session.Artist.Info(lastfm.ArtistInfoParams{Artist: testArtist}))
			},
			httpMethod: get, method: api.ArtistGetInfoMethod, level: unsigned,
			params: map[string]string{"artist": testArtist},
		},
		{
			name: "SessionClient.Chart.TopTags",
			call: func(c testClients) error {
				return errOf(c.session.Chart.TopTags())
			},
			httpMethod: get, method: api.ChartGetTopTagsMethod, level: unsigned,
			params: map[string]string{},
		},
		{
			name: "SessionClient.Geo.TopArtists",
			call: func(c testClients) error {
				return errOf(c.session.Geo.TopArtists(lastfm.GeoTopArtistsParams{Country: testCountry}))
			},
			httpMethod: get, method: api.GeoGetTopArtistsMethod, level: unsigned,
			params: map[string]string{"country": testCountry},
		},
		{
			name: "SessionClient.Library.Artists",
			call: func(c testClients) error {
				return errOf(c.session.Library.Artists(lastfm.LibraryArtistsParams{User: testUser}))
			},
			httpMethod: get, method: api.LibraryGetArtistsMethod, level: unsigned,
			params: map[string]string{"user": testUser},
		},
		{
			name: "SessionClient.Tag.Info",
			call: func(c testClients) error {
				return errOf(c.session.Tag.Info(lastfm.TagInfoParams{Tag: testTag}))
			},
			httpMethod: get, method: api.TagGetInfoMethod, level: unsigned,
			params: map[string]string{"tag": testTag},
		},
		{
			name: "SessionClient.Track.Info",
			call: func(c testClients) error {
				return errOf(c.session.Track.Info(lastfm.TrackInfoParams{Artist: testArtist, Track: testTrack}))
			},
			httpMethod: get, method: api.TrackGetInfoMethod, level: unsigned,
			params: map[string]string{"artist": testArtist, "track": testTrack},
		},
		{
			name: "SessionClient.User.Info",
			call: func(c testClients) error {
				return errOf(c.session.User.Info(testUser))
			},
			httpMethod: get, method: api.UserGetInfoMethod, level: unsigned,
			params: map[string]string{"user": testUser},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockClient := newMockClient(t, http.StatusOK, `<lfm status="ok"></lfm>`)

			err := c.call(newTestClients(t, mockClient))

			// The empty mock response can't be unmarshaled into most result
			// types. Only the request matters here, so that error is expected.
			if err != nil && !strings.Contains(err.Error(), "failed to unmarshal response") {
				t.Errorf("unexpected error: %v", err)
			}

			if len(mockClient.requests) != 1 {
				t.Fatalf("expected 1 request, got %d", len(mockClient.requests))
			}
			req := mockClient.requests[0]

			if req.httpMethod != c.httpMethod {
				t.Errorf("expected HTTP method %s, got %s", c.httpMethod, req.httpMethod)
			}

			want := map[string]string{"method": c.method.String(), "api_key": testAPIKey}
			for k, v := range c.params {
				want[k] = v
			}

			switch c.level {
			case signed, session:
				if !req.params.Has("api_sig") {
					t.Error("expected request to be signed, but api_sig is missing")
					break
				}

				unsignedParams := url.Values{}
				for k, v := range req.params {
					if k != "api_sig" {
						unsignedParams[k] = v
					}
				}
				want["api_sig"] = api.Signature(unsignedParams, testSecret)
			}
			if c.level == session {
				want["sk"] = testSessionKey
			}

			for k, v := range want {
				if !req.params.Has(k) {
					t.Errorf("missing param %s", k)
				} else if got := req.params.Get(k); got != v {
					t.Errorf("expected param %s=%q, got %q", k, v, got)
				}
			}
			for k := range req.params {
				if _, ok := want[k]; !ok {
					t.Errorf("unexpected param %s=%q", k, req.params.Get(k))
				}
			}
		})
	}
}

func TestClient_Login(t *testing.T) {
	const sessionBody = `<lfm status="ok"><session>` +
		`<name>testuser</name><key>testsessionkey</key><subscriber>1</subscriber>` +
		`</session></lfm>`

	login := func(c *Client) (*SessionClient, error) {
		return c.Login(testUser, "testpassword")
	}
	tokenLogin := func(c *Client) (*SessionClient, error) {
		return c.TokenLogin("testtoken")
	}

	cases := []struct {
		name  string
		login func(c *Client) (*SessionClient, error)

		mockStatus int
		mockBody   string

		wantMethod api.APIMethod
		wantErr    error
	}{
		{
			name:       "Login",
			login:      login,
			mockStatus: http.StatusOK,
			mockBody:   sessionBody,
			wantMethod: api.AuthGetMobileSessionMethod,
		},
		{
			name:       "TokenLogin",
			login:      tokenLogin,
			mockStatus: http.StatusOK,
			mockBody:   sessionBody,
			wantMethod: api.AuthGetSessionMethod,
		},
		{
			name:       "Login with wrong credentials",
			login:      login,
			mockStatus: http.StatusForbidden,
			mockBody:   `<lfm status="failed"><error code="4">Authentication Failed</error></lfm>`,
			wantMethod: api.AuthGetMobileSessionMethod,
			wantErr:    api.NewLastFMError(api.ErrAuthenticationFailed, ""),
		},
		{
			name:       "TokenLogin with unauthorized token",
			login:      tokenLogin,
			mockStatus: http.StatusForbidden,
			mockBody:   `<lfm status="failed"><error code="14">Unauthorized Token</error></lfm>`,
			wantMethod: api.AuthGetSessionMethod,
			wantErr:    api.NewLastFMError(api.ErrUnauthorizedToken, ""),
		},
		{
			name:       "Login response without session key",
			login:      login,
			mockStatus: http.StatusOK,
			mockBody:   `<lfm status="ok"><session><name>testuser</name><key></key><subscriber>0</subscriber></session></lfm>`,
			wantMethod: api.AuthGetMobileSessionMethod,
			wantErr:    api.NewLastFMError(api.ErrSessionRequired, api.SessionRequiredMessage),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockClient := newMockClient(t, c.mockStatus, c.mockBody)
			client := newTestClients(t, mockClient).client

			s, err := c.login(client)

			methods := make([]string, len(mockClient.requests))
			for i, req := range mockClient.requests {
				methods[i] = req.params.Get("method")
			}
			if !slices.Equal(methods, []string{c.wantMethod.String()}) {
				t.Errorf("expected one %s request, got %v", c.wantMethod, methods)
			}

			if c.wantErr != nil {
				if s != nil {
					t.Errorf("expected nil SessionClient, got %+v", s)
				}
				if !errors.Is(err, c.wantErr) {
					t.Errorf("expected error %v, got %v", c.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if s.Key() != testSessionKey || s.Name() != testUser || !s.Subscriber() {
				t.Errorf("expected session %q for %q (subscriber), got %q for %q (subscriber: %t)",
					testSessionKey, testUser, s.Key(), s.Name(), s.Subscriber())
			}
			if s.AuthAPI != client.AuthAPI {
				t.Error("expected SessionClient to share the Client's AuthAPI")
			}
		})
	}
}
