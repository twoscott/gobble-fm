package api

import (
	"errors"
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
