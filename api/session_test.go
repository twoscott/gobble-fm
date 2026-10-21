package api

import (
	"errors"
	"net/http"
	"strings"
	"testing"
)

func TestSession_CheckCredentials(t *testing.T) {
	cases := []struct {
		name         string
		sessionKey   string
		requestLevel RequestLevel
		wantErr      error
	}{
		{
			name:         "Valid session key for session level",
			sessionKey:   "valid_session_key",
			requestLevel: RequestLevelSession,
			wantErr:      nil,
		},
		{
			name:         "Missing session key for session level",
			sessionKey:   "",
			requestLevel: RequestLevelSession,
			wantErr:      NewLastFMError(ErrSessionRequired, SessionRequiredMessage),
		},
		{
			name:         "No error for API level",
			sessionKey:   "",
			requestLevel: RequestLevelAPIKey,
			wantErr:      nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := Session{
				key:     c.sessionKey,
				AuthAPI: newTestAuthAPI("testapikey", "testsecret"),
			}

			err := s.CheckCredentials(c.requestLevel)
			if err != nil && c.wantErr == nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if err == nil && c.wantErr != nil {
				t.Errorf("Expected error: %v, got nil", c.wantErr)
			}
			if err != nil && !errors.Is(err, c.wantErr) {
				t.Errorf("Expected error: %v, got: %v", c.wantErr, err)
			}
		})
	}
}

func TestNewAPISession(t *testing.T) {
	cases := []struct {
		name       string
		api        *AuthAPI
		sessionKey string
		wantErr    error
		wantErrMsg string
	}{
		{
			name:       "Valid credentials",
			api:        newTestAuthAPI("testapikey", "testsecret"),
			sessionKey: "testsessionkey",
		},
		{
			name:       "Missing session key",
			api:        newTestAuthAPI("testapikey", "testsecret"),
			sessionKey: "",
			wantErr:    NewLastFMError(ErrSessionRequired, SessionRequiredMessage),
		},
		{
			name:       "Missing secret",
			api:        newTestAuthAPI("testapikey", ""),
			sessionKey: "testsessionkey",
			wantErr:    NewLastFMError(ErrSecretRequired, SecretRequiredMessage),
		},
		{
			name:       "Missing API key",
			api:        newTestAuthAPI("", "testsecret"),
			sessionKey: "testsessionkey",
			wantErr:    NewLastFMError(ErrAPIKeyMissing, APIKeyMissingMessage),
		},
		{
			name:       "Nil AuthAPI",
			api:        nil,
			sessionKey: "testsessionkey",
			wantErrMsg: "AuthAPI is nil",
		},
		{
			name:       "Nil embedded API",
			api:        &AuthAPI{secret: "testsecret"},
			sessionKey: "testsessionkey",
			wantErrMsg: "AuthAPI is nil",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s, err := NewAPISessionFromKey(c.api, c.sessionKey)

			if c.wantErr == nil && c.wantErrMsg == "" {
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				if s == nil {
					t.Fatal("Expected session, got nil")
				}
				if s.Key() != c.sessionKey {
					t.Errorf("Expected key %q, got %q", c.sessionKey, s.Key())
				}
				return
			}

			if s != nil {
				t.Errorf("Expected nil session, got %+v", s)
			}
			if err == nil {
				t.Fatal("Expected error, got nil")
			}
			if c.wantErr != nil && !errors.Is(err, c.wantErr) {
				t.Errorf("Expected error: %v, got: %v", c.wantErr, err)
			}
			if c.wantErrMsg != "" && err.Error() != c.wantErrMsg {
				t.Errorf("Expected error message: %q, got: %q", c.wantErrMsg, err.Error())
			}
		})
	}
}

// newTestSession returns a Session built directly, bypassing the credential
// checks in NewAPISession so tests can construct Sessions with missing
// credentials.
func newTestSession(apiKey, secret, sessionKey string) *Session {
	return &Session{AuthAPI: newTestAuthAPI(apiKey, secret), key: sessionKey}
}

func TestCheckCredentials(t *testing.T) {
	errNilClient := "API HTTP Client is nil"

	cases := []struct {
		name       string
		check      func() error
		wantErr    error
		wantErrMsg string
	}{
		{
			name:  "API key not required for no level",
			check: func() error { return (&API{Client: &http.Client{}}).CheckCredentials(RequestLevelNone) },
		},
		{
			name:       "Nil HTTP client for API",
			check:      func() error { return (&API{apiKey: "testapikey"}).CheckCredentials(RequestLevelAPIKey) },
			wantErrMsg: errNilClient,
		},
		{
			name:  "Secret not required for API key level",
			check: func() error { return newTestAuthAPI("testapikey", "").CheckCredentials(RequestLevelAPIKey) },
		},
		{
			name: "Nil HTTP client for AuthAPI",
			check: func() error {
				a := &AuthAPI{API: &API{apiKey: "testapikey"}, secret: "testsecret"}
				return a.CheckCredentials(RequestLevelSecret)
			},
			wantErrMsg: errNilClient,
		},
		{
			name: "Session key not required for secret level",
			check: func() error {
				return newTestSession("testapikey", "testsecret", "").CheckCredentials(RequestLevelSecret)
			},
		},
		{
			name: "Missing secret for session level",
			check: func() error {
				return newTestSession("testapikey", "", "testsessionkey").CheckCredentials(RequestLevelSession)
			},
			wantErr: NewLastFMError(ErrSecretRequired, SecretRequiredMessage),
		},
		{
			name: "Missing API key for session level",
			check: func() error {
				return newTestSession("", "testsecret", "testsessionkey").CheckCredentials(RequestLevelSession)
			},
			wantErr: NewLastFMError(ErrAPIKeyMissing, APIKeyMissingMessage),
		},
		{
			name: "Nil HTTP client for Session",
			check: func() error {
				s := newTestSession("testapikey", "testsecret", "testsessionkey")
				s.Client = nil
				return s.CheckCredentials(RequestLevelSession)
			},
			wantErrMsg: errNilClient,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.check()

			if c.wantErr == nil && c.wantErrMsg == "" {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				return
			}

			if err == nil {
				t.Fatal("Expected error, got nil")
			}
			if c.wantErr != nil && !errors.Is(err, c.wantErr) {
				t.Errorf("Expected error: %v, got: %v", c.wantErr, err)
			}
			if c.wantErrMsg != "" && err.Error() != c.wantErrMsg {
				t.Errorf("Expected error message: %q, got: %q", c.wantErrMsg, err.Error())
			}
		})
	}
}

func TestSession_Request(t *testing.T) {
	cases := []struct {
		name string

		noSessionKey bool
		httpMethod   string

		mockError error

		wantError    error
		wantErrorMsg string
		wantTries    uint
	}{
		{
			name:         "Missing session key",
			noSessionKey: true,
			httpMethod:   http.MethodPost,
			wantError:    NewLastFMError(ErrSessionRequired, SessionRequiredMessage),
			wantTries:    0,
		},
		{
			name:         "Unsupported HTTP method",
			httpMethod:   http.MethodPut,
			wantErrorMsg: "unsupported HTTP method",
			wantTries:    0,
		},
		{
			name:       "Network error",
			httpMethod: http.MethodPost,
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

			sessionKey := "testsessionkey"
			if c.noSessionKey {
				sessionKey = ""
			}
			s := newTestSession("testapikey", "testsecret", sessionKey)
			s.Client = mockClient

			params := struct {
				Track string `url:"track"`
			}{Track: "testtrack"}

			err := s.Request(nil, c.httpMethod, TrackLoveMethod, params)

			if c.wantError == nil && c.wantErrorMsg == "" && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if c.wantError != nil && !errors.Is(err, c.wantError) {
				t.Errorf("expected error %v, got %v", c.wantError, err)
			}
			if c.wantErrorMsg != "" && (err == nil || !strings.Contains(err.Error(), c.wantErrorMsg)) {
				t.Errorf("expected error message %s, got %v", c.wantErrorMsg, err)
			}
			if c.wantTries != mockClient.tries {
				t.Errorf("expected %d tries, got %d", c.wantTries, mockClient.tries)
			}
		})
	}
}
