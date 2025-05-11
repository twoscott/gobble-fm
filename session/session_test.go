package session

import (
	"errors"
	"testing"

	"github.com/twoscott/gobble-fm/api"
)

func TestSession_CheckCredentials(t *testing.T) {
	cases := []struct {
		name         string
		sessionKey   string
		requestLevel api.RequestLevel
		wantErr      error
	}{
		{
			name:         "Valid session key for session level",
			sessionKey:   "valid_session_key",
			requestLevel: api.RequestLevelSession,
			wantErr:      nil,
		},
		{
			name:         "Missing session key for session level",
			sessionKey:   "",
			requestLevel: api.RequestLevelSession,
			wantErr:      api.NewLastFMError(api.ErrSessionRequired, api.SessionRequiredMessage),
		},
		{
			name:         "No error for API level",
			sessionKey:   "",
			requestLevel: api.RequestLevelAPIKey,
			wantErr:      nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := Session{
				SessionKey: c.sessionKey,
				API:        api.New("testapikey", "testsecret"),
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
