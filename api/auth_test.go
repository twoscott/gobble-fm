package api

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

// newTestAuthAPI returns an AuthAPI built directly, bypassing the credential
// checks in New so tests can construct AuthAPIs with missing credentials.
func newTestAuthAPI(apiKey, secret string) *AuthAPI {
	return &AuthAPI{
		API: &API{
			apiKey:    apiKey,
			UserAgent: DefaultUserAgent,
			Retries:   DefaultRetries,
			Client:    &http.Client{},
		},
		secret: secret,
	}
}

func TestNew(t *testing.T) {
	cases := []struct {
		name    string
		apiKey  string
		secret  string
		wantErr error
	}{
		{
			name:   "Valid credentials",
			apiKey: "testapikey",
			secret: "testsecret",
		},
		{
			name:    "Missing secret",
			apiKey:  "testapikey",
			secret:  "",
			wantErr: NewLastFMError(ErrSecretRequired, SecretRequiredMessage),
		},
		{
			name:    "Missing API key",
			apiKey:  "",
			secret:  "testsecret",
			wantErr: NewLastFMError(ErrAPIKeyMissing, APIKeyMissingMessage),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			a, err := New(c.apiKey, c.secret)

			if c.wantErr == nil {
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				if a == nil {
					t.Fatal("Expected AuthAPI, got nil")
				}
				if a.APIKey() != c.apiKey || a.secret != c.secret {
					t.Errorf("Expected key %q and secret %q, got %q and %q",
						c.apiKey, c.secret, a.APIKey(), a.secret)
				}
				return
			}

			if a != nil {
				t.Errorf("Expected nil AuthAPI, got %+v", a)
			}
			if !errors.Is(err, c.wantErr) {
				t.Errorf("Expected error: %v, got: %v", c.wantErr, err)
			}
		})
	}
}

func TestAPI_PostSigned(t *testing.T) {
	cases := []struct {
		name string

		noAPIkey bool
		noSecret bool
		noDest   bool

		method APIMethod
		params any

		mockStatusCode int
		mockBody       string
		mockError      error

		wantResult    string
		wantURL       string
		wantPostBody  string
		wantError     error
		wantErrorType error
		wantTries     uint
	}{
		{
			name:   "Successful request",
			method: TrackScrobbleMethod,
			params: struct {
				Track string `url:"track"`
			}{Track: "testtrack"},
			mockStatusCode: http.StatusOK,
			mockBody:       `<lfm status="ok"><scrobble><track>testtrack</track></scrobble></lfm>`,
			wantResult:     "testtrack",
			wantURL:        "https://ws.audioscrobbler.com/2.0/",
			wantPostBody:   "api_key=testapikey&api_sig=80699f4494e92c2f065f9e42eb26f5c0&method=track.scrobble&track=testtrack",
			wantTries:      1,
		},
		{
			name:   "Invalid XML response",
			method: TrackScrobbleMethod,
			params: struct {
				Track string `url:"track"`
			}{Track: "testtrack"},
			mockStatusCode: http.StatusOK,
			mockBody:       "invalid xml",
			wantError:      io.EOF,
			wantURL:        "https://ws.audioscrobbler.com/2.0/",
			wantPostBody:   "api_key=testapikey&api_sig=80699f4494e92c2f065f9e42eb26f5c0&method=track.scrobble&track=testtrack",
			wantTries:      1,
		},
		{
			name:   "API error response",
			method: TrackScrobbleMethod,
			params: struct {
				Track string `url:"track"`
			}{Track: "testtrack"},
			mockStatusCode: http.StatusBadRequest,
			mockBody:       `<lfm status="failed"><error code="6">Invalid parameters</error></lfm>`,
			wantError:      &LastFMError{Code: ErrInvalidParameters, Message: "Invalid parameters"},
			wantURL:        "https://ws.audioscrobbler.com/2.0/",
			wantPostBody:   "api_key=testapikey&api_sig=80699f4494e92c2f065f9e42eb26f5c0&method=track.scrobble&track=testtrack",
			wantTries:      1,
		},
		{
			name:   "No destination provided",
			noDest: true,
			method: TrackScrobbleMethod,
			params: struct {
				Track string `url:"track"`
			}{Track: "testtrack"},
			mockStatusCode: http.StatusOK,
			mockBody:       `<lfm status="ok"><scrobble><track>testtrack</track></scrobble></lfm>`,
			wantTries:      1,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockClient := newMockClient(t, func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: c.mockStatusCode,
					Body:       io.NopCloser(strings.NewReader(c.mockBody)),
				}, c.mockError
			})

			api := &AuthAPI{
				API: &API{
					UserAgent: DefaultUserAgent,
					Retries:   DefaultRetries,
					Client:    mockClient,
				},
			}

			if !c.noAPIkey {
				api.apiKey = "testapikey"
			}
			if !c.noSecret {
				api.secret = "testsecret"
			}

			var result struct {
				Track string `xml:"track"`
			}
			var err error
			if c.noDest {
				err = api.PostSigned(nil, c.method, c.params)
			} else {
				err = api.PostSigned(&result, c.method, c.params)
			}

			if c.wantError == nil && c.wantErrorType == nil && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if c.wantError != nil && !errors.Is(err, c.wantError) {
				t.Errorf("expected error %v, got %v", c.wantError, err)
			}

			if c.wantResult != "" && result.Track != c.wantResult {
				t.Errorf("expected result %s, got %s", c.wantResult, result.Track)
			}

			if c.wantURL != "" || c.wantPostBody != "" {
				req := mockClient.capturedReq
				if req == nil {
					t.Fatalf("expected request to be captured, but it was nil")
				}

				if c.wantURL != "" && (req.URL == nil || req.URL.String() != c.wantURL) {
					t.Errorf("expected URL %s, got %s", c.wantURL, mockClient.capturedReq.URL.String())
				}
				if c.wantPostBody != "" {
					if req.Body == nil {
						t.Fatalf("expected request body to be captured, but it was nil")
					}
					body, err := io.ReadAll(req.Body)
					if err != nil {
						t.Fatal("failed to read request body:", err)
					}
					if string(body) != c.wantPostBody {
						t.Errorf("expected post body %s, got %s", c.wantPostBody, string(body))
					}
				}
			}

			if c.wantTries != mockClient.tries {
				t.Errorf("expected %d tries, got %d", c.wantTries, mockClient.tries)
			}
		})
	}
}

func TestAPI_RequestSigned(t *testing.T) {
	cases := []struct {
		name string

		noAPIkey   bool
		noSecret   bool
		method     APIMethod
		httpMethod string
		params     any

		mockError error

		wantError    error
		wantErrorMsg string
		wantTries    uint
	}{
		{
			name:   "Unsupported HTTP method",
			method: UserGetInfoMethod,
			params: struct {
				User string `url:"user"`
			}{User: "testuser"},
			httpMethod:   http.MethodPut,
			wantErrorMsg: "unsupported HTTP method",
			wantTries:    0,
		},
		{
			name:     "Missing API key",
			noAPIkey: true,
			method:   UserGetInfoMethod,
			params: struct {
				User string `url:"user"`
			}{User: "testuser"},
			httpMethod: http.MethodGet,
			wantError:  &LastFMError{Code: ErrAPIKeyMissing, Message: APIKeyMissingMessage},
			wantTries:  0,
		},
		{
			name:     "Missing secret",
			noSecret: true,
			method:   TrackScrobbleMethod,
			params: struct {
				Track string `url:"track"`
			}{Track: "testtrack"},
			wantError: &LastFMError{Code: ErrSecretRequired, Message: SecretRequiredMessage},
			wantTries: 0,
		},
		{
			name:   "Network error",
			method: UserGetInfoMethod,
			params: struct {
				User string `url:"user"`
			}{User: "testuser"},
			httpMethod: http.MethodGet,
			mockError:  errNetwork,
			wantError:  errNetwork,
			wantTries:  1,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var doFunc func(req *http.Request) (*http.Response, error)
			if c.mockError != nil {
				doFunc = respondWith(t, mockResponse{err: c.mockError})
			}
			mockClient := newMockClient(t, doFunc)

			api := &AuthAPI{
				API: &API{
					UserAgent: DefaultUserAgent,
					Retries:   DefaultRetries,
					Client:    mockClient,
				},
			}

			if !c.noAPIkey {
				api.apiKey = "testapikey"
			}
			if !c.noSecret {
				api.secret = "testsecret"
			}

			var user struct {
				Name string `xml:"name"`
			}

			err := api.RequestSigned(&user, c.httpMethod, c.method, c.params)

			if c.wantError == nil && c.wantErrorMsg == "" && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if c.wantError != nil && !errors.Is(err, c.wantError) {
				t.Errorf("expected error %v, got %v", c.wantError, err)
			}

			if c.wantErrorMsg != "" && err != nil && !strings.Contains(err.Error(), c.wantErrorMsg) {
				t.Errorf("expected error message %s, got %s", c.wantErrorMsg, err.Error())
			}

			if c.wantTries != mockClient.tries {
				t.Errorf("expected %d tries, got %d", c.wantTries, mockClient.tries)
			}
		})
	}
}
