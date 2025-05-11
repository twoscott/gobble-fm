package api

import (
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

var (
	maxTries   = DefaultRetries + 1
	errNetwork = errors.New("network error")
)

type mockHTTPClient struct {
	tries       uint
	capturedReq *http.Request
	doFunc      func(req *http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	m.capturedReq = req
	m.tries++
	return m.doFunc(req)
}

func TestAPI_Get(t *testing.T) {
	cases := []struct {
		name string

		noAPIkey bool
		method   APIMethod
		params   any

		mockStatusCode int
		mockBody       string
		mockError      error

		wantResult    string
		wantURL       string
		wantError     error
		wantErrorType error
		wantTries     uint
	}{
		{
			name:   "Successful request",
			method: UserGetInfoMethod,
			params: struct {
				User string `url:"user"`
			}{User: "testuser"},
			mockStatusCode: http.StatusOK,
			mockBody:       `<lfm status="ok"><user><name>testuser</name></user></lfm>`,
			wantResult:     "testuser",
			wantURL:        "https://ws.audioscrobbler.com/2.0/?api_key=testapikey&method=user.getInfo&user=testuser",
			wantTries:      1,
		},
		{
			name:   "Invalid XML response",
			method: UserGetInfoMethod,
			params: struct {
				User string `url:"user"`
			}{User: "testuser"},
			mockStatusCode: http.StatusOK,
			mockBody:       "invalid xml",
			wantURL:        "https://ws.audioscrobbler.com/2.0/?api_key=testapikey&method=user.getInfo&user=testuser",
			wantError:      io.EOF,
			wantTries:      1,
		},
		{
			name:   "Malformed XML response",
			method: UserGetInfoMethod,
			params: struct {
				User string `url:"user"`
			}{User: "testuser"},
			mockStatusCode: http.StatusOK,
			mockBody:       `<lfm status="ok"><user><name>testuser<name><user><lfm>`,
			wantURL:        "https://ws.audioscrobbler.com/2.0/?api_key=testapikey&method=user.getInfo&user=testuser",
			wantErrorType:  &xml.SyntaxError{},
			wantTries:      1,
		},
		{
			name:   "Wrong XML element",
			method: UserGetInfoMethod,
			params: struct {
				User string `url:"user"`
			}{User: "testuser"},
			mockStatusCode: http.StatusOK,
			mockBody:       `<xml></xml>`,
			wantURL:        "https://ws.audioscrobbler.com/2.0/?api_key=testapikey&method=user.getInfo&user=testuser",
			wantErrorType:  xml.UnmarshalError(""),
			wantTries:      1,
		},
		{
			name:   "API error response",
			method: UserGetInfoMethod,
			params: struct {
				User string `url:"user"`
			}{User: "testuser"},
			mockStatusCode: http.StatusBadRequest,
			mockBody:       `<lfm status="failed"><error code="6">Invalid parameters</error></lfm>`,
			wantURL:        "https://ws.audioscrobbler.com/2.0/?api_key=testapikey&method=user.getInfo&user=testuser",
			wantError:      &LastFMError{Code: ErrInvalidParameters, Message: "Invalid parameters"},
			wantTries:      1,
		},
		{
			name:   "HTTP Error wrapped in LastFMError",
			method: UserGetInfoMethod,
			params: struct {
				User string `url:"invalidparam"`
			}{User: "testuser"},
			mockStatusCode: http.StatusBadRequest,
			mockBody:       `<lfm status="failed"><error code="6">Invalid parameters</error></lfm>`,
			wantURL:        "https://ws.audioscrobbler.com/2.0/?api_key=testapikey&invalidparam=testuser&method=user.getInfo",
			wantError:      &HTTPError{StatusCode: http.StatusBadRequest, Message: "Bad Request"},
			wantTries:      1,
		},
		{
			name:   "API rate limit error response",
			method: UserGetInfoMethod,
			params: struct {
				User string `url:"user"`
			}{User: "testuser"},
			mockStatusCode: http.StatusBadRequest,
			mockBody:       `<lfm status="failed"><error code="29">Rate limit exceeded</error></lfm>`,
			wantURL:        "https://ws.audioscrobbler.com/2.0/?api_key=testapikey&method=user.getInfo&user=testuser",
			wantError:      &LastFMError{Code: ErrRateLimitExceeded, Message: "Rate limit exceeded"},
			wantTries:      maxTries,
		},
		{
			name:   "HTTP 500 error response",
			method: UserGetInfoMethod,
			params: struct {
				User string `url:"user"`
			}{User: "testuser"},
			mockStatusCode: http.StatusInternalServerError,
			wantURL:        "https://ws.audioscrobbler.com/2.0/?api_key=testapikey&method=user.getInfo&user=testuser",
			wantError: &HTTPError{
				StatusCode: http.StatusInternalServerError,
				Message:    "Internal Server Error",
			},
			wantTries: maxTries,
		},
		{
			name:   "HTTP 429 error response",
			method: UserGetInfoMethod,
			params: struct {
				User string `url:"user"`
			}{User: "testuser"},
			mockStatusCode: http.StatusTooManyRequests,
			wantURL:        "https://ws.audioscrobbler.com/2.0/?api_key=testapikey&method=user.getInfo&user=testuser",
			wantError: &HTTPError{
				StatusCode: http.StatusTooManyRequests,
				Message:    "Too Many Requests",
			},
			wantTries: maxTries,
		},
		{
			name:   "HTTP 400 error response",
			method: "nil.invalidmethod",
			params: struct {
				User string `url:"user"`
			}{User: ""},
			mockStatusCode: http.StatusBadRequest,
			wantURL:        "https://ws.audioscrobbler.com/2.0/?api_key=testapikey&method=nil.invalidmethod&user=",
			wantError:      &HTTPError{StatusCode: http.StatusBadRequest, Message: "Bad Request"},
			wantTries:      1,
		},
		{
			name:   "HTTP 200 empty response",
			method: UserGetInfoMethod,
			params: struct {
				User string `url:"user"`
			}{User: ""},
			mockStatusCode: http.StatusOK,
			wantURL:        "https://ws.audioscrobbler.com/2.0/?api_key=testapikey&method=user.getInfo&user=",
			wantError:      io.EOF,
			wantTries:      1,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockClient := &mockHTTPClient{
				doFunc: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: c.mockStatusCode,
						Body:       io.NopCloser(strings.NewReader(c.mockBody)),
					}, c.mockError
				},
			}

			api := &API{
				UserAgent: DefaultUserAgent,
				Retries:   DefaultRetries,
				Client:    mockClient,
			}

			if !c.noAPIkey {
				api.APIKey = "testapikey"
			}

			var user struct {
				Name string `xml:"name"`
			}

			err := api.Get(&user, c.method, c.params)

			if c.wantError == nil && c.wantErrorType == nil && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if c.wantError != nil && !errors.Is(err, c.wantError) {
				t.Errorf("expected error %v, got %v", c.wantError, err)
			}

			if c.wantErrorType != nil {
				switch c.wantErrorType.(type) {
				case *xml.SyntaxError:
					if _, ok := err.(*xml.SyntaxError); ok {
						return
					}
				case xml.UnmarshalError:
					if _, ok := err.(xml.UnmarshalError); ok {
						return
					}
				default:
					t.Errorf("expected error type %T, got %T", c.wantErrorType, err)
				}
			}

			if c.wantResult != "" && user.Name != c.wantResult {
				t.Errorf("expected result %s, got %s", c.wantResult, user.Name)
			}

			if c.wantURL != "" {
				req := mockClient.capturedReq
				if req == nil {
					t.Fatalf("expected request to be captured, but it was nil")
				}
				if req.URL == nil || req.URL.String() != c.wantURL {
					t.Errorf("expected URL %s, got %s", c.wantURL, mockClient.capturedReq.URL.String())
				}
			}

			if c.wantTries != mockClient.tries {
				t.Errorf("expected %d tries, got %d", c.wantTries, mockClient.tries)
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
			mockClient := &mockHTTPClient{
				doFunc: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: c.mockStatusCode,
						Body:       io.NopCloser(strings.NewReader(c.mockBody)),
					}, c.mockError
				},
			}

			api := &API{
				UserAgent: DefaultUserAgent,
				Retries:   DefaultRetries,
				Client:    mockClient,
			}

			if !c.noAPIkey {
				api.APIKey = "testapikey"
			}
			if !c.noSecret {
				api.Secret = "testsecret"
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

				if c.wantURL != "" && req.URL == nil || req.URL.String() != c.wantURL {
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

func TestAPI_Request(t *testing.T) {
	cases := []struct {
		name string

		noAPIkey bool
		noClient bool

		method     APIMethod
		httpMethod string
		params     any

		mockError error

		wantResult   string
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
			mockClient := &mockHTTPClient{
				doFunc: func(req *http.Request) (*http.Response, error) {
					return nil, c.mockError
				},
			}

			api := &API{
				UserAgent: DefaultUserAgent,
				Retries:   DefaultRetries,
				Client:    mockClient,
			}

			if !c.noAPIkey {
				api.APIKey = "testapikey"
			}

			var user struct {
				Name string `xml:"name"`
			}

			err := api.Request(&user, c.httpMethod, c.method, c.params)

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
			mockClient := &mockHTTPClient{
				doFunc: func(req *http.Request) (*http.Response, error) {
					return nil, c.mockError
				},
			}

			api := &API{
				UserAgent: DefaultUserAgent,
				Retries:   DefaultRetries,
				Client:    mockClient,
			}

			if !c.noAPIkey {
				api.APIKey = "testapikey"
			}
			if !c.noSecret {
				api.Secret = "testsecret"
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

func TestSignature(t *testing.T) {
	cases := []struct {
		name       string
		params     url.Values
		secret     string
		wantResult string
	}{
		{
			name: "Valid parameters",
			params: url.Values{
				"api_key": {"testapikey"},
				"artist":  {"testartist"},
				"method":  {"track.scrobble"},
				"track":   {"testtrack"},
			},
			secret:     "testsecret",
			wantResult: "fd8d12c778a57123318097b2dcc5911b",
		},
		{
			name: "Exclude format and callback",
			params: url.Values{
				"api_key":  {"testapikey"},
				"artist":   {"testartist"},
				"method":   {"track.scrobble"},
				"track":    {"testtrack"},
				"format":   {"xml"},
				"callback": {"http://example.com"},
			},
			secret:     "testsecret",
			wantResult: "fd8d12c778a57123318097b2dcc5911b",
		},
		{
			name:       "Empty parameters",
			params:     url.Values{},
			secret:     "testsecret",
			wantResult: "217df19d942a4a990ebeed63a983292f",
		},
		{
			name: "Empty secret",
			params: url.Values{
				"api_key": {"testapikey"},
				"method":  {"track.scrobble"},
				"track":   {"testtrack"},
			},
			secret:     "",
			wantResult: "00caa126c50bc88fc0afbbcc78961334",
		},
		{
			name:   "Empty parameters and secret",
			params: url.Values{},
			secret: "",
			// d41d8cd98f00b204e9800998ecf8427e is the MD5 hash of an empty string
			wantResult: "d41d8cd98f00b204e9800998ecf8427e",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := Signature(c.params, c.secret)
			if got != c.wantResult {
				t.Errorf("expected %s, got %s", c.wantResult, got)
			}
		})
	}
}
