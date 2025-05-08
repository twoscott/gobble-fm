package api

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
)

type ErrorCode int

// https://www.last.fm/api/errorcodes
const (
	NoError ErrorCode = iota // 0 (No error)

	_                      // 1 - "This error does not exist"
	InvalidService         // 2
	InvalidMethod          // 3
	AuthenticationFailed   // 4
	InvalidFormat          // 5
	InvalidParameters      // 6
	InvalidResource        // 7
	OperationFailed        // 8
	InvalidSessionKey      // 9
	InvalidAPIKey          // 10
	ServiceOffline         // 11
	SubscribersOnly        // 12
	InvalidMethodSignature // 13
	UnauthorizedToken      // 14
	ItemNotStreamable      // 15
	ServiceUnavailable     // 16
	UserNotLoggedIn        // 17
	TrialExpired           // 18
	_                      // 19 - "This error does not exist"
	NotEnoughContent       // 20
	NotEnoughMembers       // 21
	NotEnoughFans          // 22
	NotEnoughNeighbours    // 23
	NoPeakRadio            // 24
	RadioNotFound          // 25
	APIKeySuspended        // 26
	Deprecated             // 27
	_                      // 28 - Missing
	RateLimitExceeded      // 29
)

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

type LastFMError struct {
	Code      ErrorCode `xml:"code,attr"`
	Message   string    `xml:",chardata"`
	httpError *HTTPError
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
func (e *LastFMError) WrapHTTPError(httpError HTTPError) *LastFMError {
	e.httpError = &httpError
	return e
}

// WrapNewHTTPError wraps a new HTTP error from a HTTP response into the
// LastFMError.
func (e *LastFMError) WrapResponse(res *http.Response) *LastFMError {
	return e.WrapHTTPError(*NewHTTPError(res))
}

// IsCode checks if the error code matches the given code.
func (e LastFMError) IsCode(code ErrorCode) bool {
	return e.Code == code
}

// HasErrorCode checks if the error code indicated an error (not 0).
func (e LastFMError) HasErrorCode() bool {
	return e.Code != NoError
}

// IsRateLimit checks if the error is a rate limit exceeded error.
func (e LastFMError) IsRateLimit() bool {
	return e.Code == RateLimitExceeded
}

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
