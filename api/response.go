package api

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
)

type ErrorCode int

// ErrorCode represents an error code returned by the Last.fm API.
//   - https://www.last.fm/api/errorcodes
//   - https://lastfm-docs.github.io/api-docs/codes/
const (
	NoError ErrorCode = iota // 0 (No error)

	_                         // 1 - "This error does not exist"
	ErrInvalidService         // 2
	ErrInvalidMethod          // 3
	ErrAuthenticationFailed   // 4
	ErrInvalidFormat          // 5
	ErrInvalidParameters      // 6
	ErrInvalidResource        // 7
	ErrOperationFailed        // 8
	ErrInvalidSessionKey      // 9
	ErrInvalidAPIKey          // 10
	ErrServiceOffline         // 11
	ErrSubscribersOnly        // 12
	ErrInvalidMethodSignature // 13
	ErrUnauthorizedToken      // 14
	ErrItemNotStreamable      // 15
	ErrServiceUnavailable     // 16
	ErrUserNotLoggedIn        // 17
	ErrTrialExpired           // 18
	_                         // 19 - "This error does not exist"
	ErrNotEnoughContent       // 20
	ErrNotEnoughMembers       // 21
	ErrNotEnoughFans          // 22
	ErrNotEnoughNeighbours    // 23
	ErrNoPeakRadio            // 24
	ErrRadioNotFound          // 25
	ErrAPIKeySuspended        // 26
	ErrDeprecated             // 27
	_                         // 28 - Missing
	ErrRateLimitExceeded      // 29
)

// Custom error codes.
const (
	ErrAPIKeyMissing ErrorCode = iota + 100
	ErrSecretRequired
	ErrSessionRequired
)

// Custom error messages.
const (
	APIKeyMissingMessage   = "API Key is missing"
	SecretRequiredMessage  = "Method requires API secret"
	SessionRequiredMessage = "Method requires user authentication (session key)"
)

type LFMWrapper struct {
	XMLName  xml.Name `xml:"lfm"`
	Status   string   `xml:"status,attr"`
	InnerXML []byte   `xml:",innerxml"`
}

// Empty checks if the InnerXML is empty.
func (lf *LFMWrapper) Empty() bool {
	return len(lf.InnerXML) == 0
}

// StatusOK checks if the status is "ok".
func (lf *LFMWrapper) StatusOK() bool {
	return lf.Status == "ok"
}

// StatusFailed checks if the status is "failed".
func (lf *LFMWrapper) StatusFailed() bool {
	return lf.Status == "failed"
}

// UnwrapError attempts to extract and return a LastFMError from the LFMWrapper
// instance. If the status of the LFMWrapper indicates success, it returns nil
// for both the error and LastFMError. Otherwise, it unmarshals the InnerXML
// field into a LastFMError structure. If the unmarshaling succeeds and the
// LastFMError contains an error code, it returns the LastFMError. If no error
// code is present or unmarshaling fails, it returns an appropriate error.
func (lf *LFMWrapper) UnwrapError() (*LastFMError, error) {
	if lf.StatusOK() {
		return nil, nil
	}

	var lferr LastFMError
	if err := lf.UnmarshalInnerXML(&lferr); err != nil {
		return nil, err
	}
	if lferr.HasErrorCode() {
		return &lferr, nil
	}

	return nil, errors.New("no error code in response")
}

// UnmarshalInnerXML unmarshals the InnerXML field into the provided destination
// variable.
func (lf *LFMWrapper) UnmarshalInnerXML(dest any) error {
	return xml.Unmarshal(lf.InnerXML, dest)
}

// LastFMError represents an error returned by Last.fm.
type LastFMError struct {
	Code      ErrorCode `xml:"code,attr"`
	Message   string    `xml:",chardata"`
	httpError *HTTPError
}

// NewLastFMError creates a new LastFMError instance from an error code and
// message.
func NewLastFMError(code ErrorCode, message string) *LastFMError {
	return &LastFMError{
		Code:    code,
		Message: message,
	}
}

// Error implements the error interface.
func (e *LastFMError) Error() string {
	return fmt.Sprintf("Last.fm Error: %d - %s", e.Code, e.Message)
}

// Is checks if the error matches the target error.
func (e *LastFMError) Is(target error) bool {
	if t, ok := target.(*LastFMError); ok {
		return e.IsCode(t.Code)
	}
	return false
}

// Unwrap implements the error interface to return the underlying HTTP error.
func (e *LastFMError) Unwrap() error {
	return e.httpError
}

// WrapHTTPError wraps an HTTP error into the LastFMError.
func (e *LastFMError) WrapHTTPError(httpError *HTTPError) *LastFMError {
	e.httpError = httpError
	return e
}

// WrapNewHTTPError wraps a new HTTP error from a HTTP response into the
// LastFMError.
func (e *LastFMError) WrapResponse(res *http.Response) *LastFMError {
	return e.WrapHTTPError(NewHTTPError(res))
}

// IsCode checks if the error code matches the given code.
func (e LastFMError) IsCode(code ErrorCode) bool {
	return e.Code == code
}

// HasErrorCode checks if the error code indicated an error (not 0).
func (e LastFMError) HasErrorCode() bool {
	return e.Code != NoError
}

// ShouldRetry returns true if the error code indicates that the request can be
// retried.
func (e LastFMError) ShouldRetry() bool {
	return e.Code == ErrOperationFailed ||
		e.Code == ErrServiceUnavailable ||
		e.Code == ErrRateLimitExceeded
}

// HTTPError represents an error that occurred during an HTTP request.
type HTTPError struct {
	StatusCode int
	Message    string
}

// NewHTTPError creates a new HTTPError instance from an HTTP response.
func NewHTTPError(res *http.Response) *HTTPError {
	if res == nil {
		return &HTTPError{
			StatusCode: http.StatusInternalServerError,
			Message:    "nil response",
		}
	}

	return &HTTPError{
		StatusCode: res.StatusCode,
		Message:    http.StatusText(res.StatusCode),
	}
}

// Error implements the error interface.
func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.StatusCode, e.Message)
}

// Is checks if the error matches the target error.
func (e *HTTPError) Is(target error) bool {
	if t, ok := target.(*HTTPError); ok {
		return e.StatusCode == t.StatusCode
	}
	return false
}
