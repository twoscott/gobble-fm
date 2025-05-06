package api

import (
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

const defaultRetries = 3

type mockHTTPClient struct {
	tries  uint
	doFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	m.tries++
	return m.doFunc(req)
}

func TestRequest(t *testing.T) {
	networkError := errors.New("network error")

	cases := []struct {
		name string

		httpMethod string
		method     APIMethod
		params     any

		mockStatusCode int
		mockBody       string
		wantResult     string

		mockError     error
		wantError     error
		wantErrorType error

		expectedRetries uint
	}{
		{
			name:       "Successful GET request",
			httpMethod: http.MethodGet,
			method:     UserGetInfoMethod,
			params: struct {
				User string `url:"user"`
			}{User: "testuser"},
			mockStatusCode:  http.StatusOK,
			mockBody:        `<lfm status="ok"><user><name>testuser</name></user></lfm>`,
			wantResult:      "testuser",
			expectedRetries: 0,
		},
		{
			name:            "HTTP client error",
			httpMethod:      http.MethodGet,
			method:          UserGetInfoMethod,
			params:          nil,
			mockError:       networkError,
			wantError:       networkError,
			expectedRetries: 0,
		},
		{
			name:       "Invalid XML response",
			httpMethod: http.MethodGet,
			method:     UserGetInfoMethod,
			params: struct {
				User string `url:"user"`
			}{User: "testuser"},
			mockStatusCode:  http.StatusOK,
			mockBody:        "invalid xml",
			wantError:       io.EOF,
			expectedRetries: 0,
		},
		{
			name:       "Malformed XML response",
			httpMethod: http.MethodGet,
			method:     UserGetInfoMethod,
			params: struct {
				User string `url:"user"`
			}{User: "testuser"},
			mockStatusCode:  http.StatusOK,
			mockBody:        `<lfm status="ok"><user><name>testuser<name><user><lfm>`,
			wantErrorType:   &xml.SyntaxError{},
			expectedRetries: 0,
		},
		{
			name:       "Wrong XML element response",
			httpMethod: http.MethodGet,
			method:     UserGetInfoMethod,
			params: struct {
				User string `url:"user"`
			}{User: "testuser"},
			mockStatusCode:  http.StatusOK,
			mockBody:        `<xml></xml>`,
			wantErrorType:   xml.UnmarshalError(""),
			expectedRetries: 0,
		},
		{
			name:       "API parameters error response",
			httpMethod: http.MethodGet,
			method:     UserGetInfoMethod,
			params: struct {
				User string `url:"user"`
			}{User: "testuser"},
			mockStatusCode:  http.StatusBadRequest,
			mockBody:        `<lfm status="failed"><error code="6">Invalid parameters</error></lfm>`,
			wantError:       &LastFMError{Code: InvalidParameters, Message: "Invalid parameters"},
			expectedRetries: 0,
		},
		{
			name:       "API rate limit error response",
			httpMethod: http.MethodGet,
			method:     UserGetInfoMethod,
			params: struct {
				User string `url:"user"`
			}{User: "testuser"},
			mockStatusCode:  http.StatusBadRequest,
			mockBody:        `<lfm status="failed"><error code="29">Rate limit exceeded</error></lfm>`,
			wantError:       &LastFMError{Code: RateLimitExceeded, Message: "Rate limit exceeded"},
			expectedRetries: defaultRetries,
		},
		{
			name:       "HTTP 500 error response",
			httpMethod: http.MethodGet,
			method:     UserGetInfoMethod,
			params: struct {
				User string `url:"user"`
			}{User: "testuser"},
			mockStatusCode: http.StatusInternalServerError,
			wantError: &HTTPError{
				StatusCode: http.StatusInternalServerError,
				Message:    "Internal Server Error",
			},
			expectedRetries: defaultRetries,
		},
		{
			name:       "HTTP 429 error response",
			httpMethod: http.MethodGet,
			method:     UserGetInfoMethod,
			params: struct {
				User string `url:"user"`
			}{User: "testuser"},
			mockStatusCode: http.StatusTooManyRequests,
			wantError: &HTTPError{
				StatusCode: http.StatusTooManyRequests,
				Message:    "Too Many Requests",
			},
			expectedRetries: defaultRetries,
		},
		{
			name:       "HTTP 400 error response",
			httpMethod: http.MethodGet,
			method:     "nil.invalidmethod",
			params: struct {
				User string `url:"user"`
			}{User: ""},
			mockStatusCode: http.StatusBadRequest,
			wantError: &HTTPError{
				StatusCode: http.StatusBadRequest,
				Message:    "Bad Request",
			},
			expectedRetries: 0,
		},
		{
			name:       "HTTP 200 empty response",
			httpMethod: http.MethodGet,
			method:     UserGetInfoMethod,
			params: struct {
				User string `url:"user"`
			}{User: ""},
			mockStatusCode:  http.StatusOK,
			wantError:       io.EOF,
			expectedRetries: 0,
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
				APIKey:    "testapikey",
				userAgent: DefaultUserAgent,
				retries:   defaultRetries,
				client:    mockClient,
			}

			var user struct {
				Name string `xml:"name"`
			}

			err := api.Request(&user, c.httpMethod, c.method, c.params)

			if c.wantError != nil {
				if !errors.Is(err, c.wantError) {
					t.Errorf("expected error %v, got %v", c.wantError, err)
				}
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

			if c.wantResult != "" {
				if user.Name != c.wantResult {
					t.Errorf("expected result %s, got %s", c.wantResult, user.Name)
				}
			}

			if c.expectedRetries != mockClient.tries-1 {
				t.Errorf("expected %d retries, got %d", c.expectedRetries, mockClient.tries-1)
			}
		})
	}
}
