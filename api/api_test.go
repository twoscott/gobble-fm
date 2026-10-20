package api

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

var (
	maxTries             = DefaultRetries + 1
	errNetwork           = errors.New("network error")
	errUnexpectedRequest = errors.New("unexpected request")
)

type mockHTTPClient struct {
	tries       uint
	capturedReq *http.Request
	// bodies holds the body of each request in the order they were sent.
	bodies []string
	doFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	m.capturedReq = req
	m.tries++

	var body []byte
	if req.Body != nil {
		var err error
		if body, err = io.ReadAll(req.Body); err != nil {
			return nil, err
		}
		req.Body = io.NopCloser(bytes.NewReader(body))
	}
	m.bodies = append(m.bodies, string(body))

	return m.doFunc(req)
}

func TestNewKeyOnly(t *testing.T) {
	cases := []struct {
		name    string
		apiKey  string
		wantErr error
	}{
		{
			name:   "Valid API key",
			apiKey: "testapikey",
		},
		{
			name:    "Missing API key",
			apiKey:  "",
			wantErr: NewLastFMError(ErrAPIKeyMissing, APIKeyMissingMessage),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			a, err := NewKeyOnly(c.apiKey)

			if c.wantErr == nil {
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				if a == nil {
					t.Fatal("Expected API, got nil")
				}
				if a.APIKey() != c.apiKey {
					t.Errorf("Expected key %q, got %q", c.apiKey, a.APIKey())
				}
				return
			}

			if a != nil {
				t.Errorf("Expected nil API, got %+v", a)
			}
			if !errors.Is(err, c.wantErr) {
				t.Errorf("Expected error: %v, got: %v", c.wantErr, err)
			}
		})
	}
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
			name:   "API operation failed response",
			method: UserGetInfoMethod,
			params: struct {
				User string `url:"user"`
			}{User: "testuser"},
			mockStatusCode: http.StatusBadRequest,
			mockBody:       `<lfm status="failed"><error code="8">Operation failed</error></lfm>`,
			wantURL:        "https://ws.audioscrobbler.com/2.0/?api_key=testapikey&method=user.getInfo&user=testuser",
			wantError:      &LastFMError{Code: ErrOperationFailed, Message: "Operation failed"},
			wantTries:      maxTries,
		},
		{
			name:   "API service unavailable response",
			method: UserGetInfoMethod,
			params: struct {
				User string `url:"user"`
			}{User: "testuser"},
			mockStatusCode: http.StatusBadRequest,
			mockBody:       `<lfm status="failed"><error code="16">Service unavailable</error></lfm>`,
			wantURL:        "https://ws.audioscrobbler.com/2.0/?api_key=testapikey&method=user.getInfo&user=testuser",
			wantError:      &LastFMError{Code: ErrServiceUnavailable, Message: "Service unavailable"},
			wantTries:      maxTries,
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
		{
			name:           "nil params",
			method:         UserGetInfoMethod,
			params:         nil,
			mockStatusCode: http.StatusOK,
			mockBody:       `<lfm status="ok"><user><name>testuser</name></user></lfm>`,
			wantResult:     "testuser",
			wantURL:        "https://ws.audioscrobbler.com/2.0/?api_key=testapikey&method=user.getInfo",
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
				api.apiKey = "testapikey"
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
				api.apiKey = "testapikey"
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

func TestAPI_Retries(t *testing.T) {
	const (
		okBody               = `<lfm status="ok"><user><name>testuser</name></user></lfm>`
		operationFailedBody  = `<lfm status="failed"><error code="8">Operation failed</error></lfm>`
		invalidParamsBody    = `<lfm status="failed"><error code="6">Invalid parameters</error></lfm>`
		wantPostBody         = "api_key=testapikey&method=user.getInfo&user=testuser"
		internalServerErrMsg = "Internal Server Error"
	)

	var (
		ok              = mockResponse{status: http.StatusOK, body: okBody}
		serverError     = mockResponse{status: http.StatusInternalServerError}
		operationFailed = mockResponse{status: http.StatusBadRequest, body: operationFailedBody}
	)

	cases := []struct {
		name string

		retries   uint
		post      bool
		responses []mockResponse

		wantResult string
		wantError  error
		wantTries  uint
	}{
		{
			name:       "Recovers after HTTP 500",
			retries:    DefaultRetries,
			responses:  []mockResponse{serverError, ok},
			wantResult: "testuser",
			wantTries:  2,
		},
		{
			name:       "Recovers after HTTP 429",
			retries:    DefaultRetries,
			responses:  []mockResponse{{status: http.StatusTooManyRequests}, ok},
			wantResult: "testuser",
			wantTries:  2,
		},
		{
			name:       "Recovers after retryable Last.fm error",
			retries:    DefaultRetries,
			responses:  []mockResponse{operationFailed, ok},
			wantResult: "testuser",
			wantTries:  2,
		},
		{
			name:       "Resends POST body on retry",
			retries:    DefaultRetries,
			post:       true,
			responses:  []mockResponse{serverError, serverError, ok},
			wantResult: "testuser",
			wantTries:  3,
		},
		{
			name:      "Retries disabled",
			retries:   0,
			responses: []mockResponse{serverError},
			wantError: &HTTPError{StatusCode: http.StatusInternalServerError, Message: internalServerErrMsg},
			wantTries: 1,
		},
		{
			name:      "Custom retry count",
			retries:   2,
			responses: []mockResponse{serverError, serverError, serverError},
			wantError: &HTTPError{StatusCode: http.StatusInternalServerError, Message: internalServerErrMsg},
			wantTries: 3,
		},
		{
			name:    "Stops at non-retryable Last.fm error",
			retries: DefaultRetries,
			responses: []mockResponse{
				serverError,
				{status: http.StatusBadRequest, body: invalidParamsBody},
			},
			wantError: &LastFMError{Code: ErrInvalidParameters, Message: "Invalid parameters"},
			wantTries: 2,
		},
		{
			name:      "Returns network error without retrying it",
			retries:   DefaultRetries,
			responses: []mockResponse{serverError, {err: errNetwork}},
			wantError: errNetwork,
			wantTries: 2,
		},
		{
			name:      "Doesn't reuse Last.fm error from earlier attempt",
			retries:   DefaultRetries,
			responses: []mockResponse{operationFailed, {status: http.StatusOK, body: "invalid xml"}},
			wantError: io.EOF,
			wantTries: 2,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockClient := &mockHTTPClient{doFunc: respondWith(t, c.responses...)}

			api := &API{
				APIKey:    "testapikey",
				UserAgent: DefaultUserAgent,
				Retries:   c.retries,
				Client:    mockClient,
			}

			params := struct {
				User string `url:"user"`
			}{User: "testuser"}

			var user struct {
				Name string `xml:"name"`
			}

			var err error
			if c.post {
				err = api.Post(&user, UserGetInfoMethod, params)
			} else {
				err = api.Get(&user, UserGetInfoMethod, params)
			}

			if c.wantError == nil && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if c.wantError != nil && !errors.Is(err, c.wantError) {
				t.Errorf("expected error %v, got %v", c.wantError, err)
			}
			if user.Name != c.wantResult {
				t.Errorf("expected result %q, got %q", c.wantResult, user.Name)
			}
			if c.wantTries != mockClient.tries {
				t.Errorf("expected %d tries, got %d", c.wantTries, mockClient.tries)
			}

			if c.post {
				for i, body := range mockClient.bodies {
					if body != wantPostBody {
						t.Errorf("expected attempt %d body %q, got %q", i+1, wantPostBody, body)
					}
				}
			}
		})
	}
}
