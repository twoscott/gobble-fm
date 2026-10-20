package api

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"maps"
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/twoscott/gobble-fm/v2/lastfm"
)

// AuthAPI embeds API and provides methods for making signed requests which are required for
// fetching and authenticating session tokens for users.
type AuthAPI struct {
	*API
	// secret is the Last.fm API secret used to sign requests. It is stored as a
	// pointer because fmt can't call String on unexported fields, so a Secret
	// value here would be printed in full. fmt prints nested pointers as
	// addresses, so the secret isn't leaked when an AuthAPI, or anything
	// embedding it, is printed or logged.
	secret string
}

// New returns a new instance of AuthAPI with the given API key and secret.
// Returns an error if the API key or secret is missing.
func New(apiKey, secret string) (*AuthAPI, error) {
	return NewWithTimeout(apiKey, secret, DefaultTimeout)
}

// NewWithTimeout returns a new instance of AuthAPI with the given API key,
// secret, and timeout settings. The timeout is specified in seconds. Returns an
// error if the API key or secret is missing.
func NewWithTimeout(apiKey, secret string, timeout int) (*AuthAPI, error) {
	t := time.Duration(timeout) * time.Second

	a := &AuthAPI{
		API: &API{
			apiKey:    apiKey,
			UserAgent: DefaultUserAgent,
			Retries:   DefaultRetries,
			Client:    &http.Client{Timeout: t},
		},
		secret: secret,
	}

	if err := a.CheckCredentials(RequestLevelSecret); err != nil {
		return nil, err
	}

	return a, nil
}

// GetSigned sends an HTTP GET request to the API with the specified method and
// parameters, signed with the API secret. The response is unmarshaled into the
// provided destination.
func (a *AuthAPI) GetSigned(dest any, method APIMethod, params any) error {
	return a.RequestSigned(dest, http.MethodGet, method, params)
}

// PostSigned sends an HTTP POST request to the API with the specified method
// and parameters, signed with the API secret. The response is unmarshaled into
// the provided destination.
func (a *AuthAPI) PostSigned(dest any, method APIMethod, params any) error {
	return a.RequestSigned(dest, http.MethodPost, method, params)
}

// RequestSigned sends an HTTP request to the API with the specified method and
// parameters, signed with the API secret. The response is unmarshaled into the
// provided destination.
//
// Parameters:
//   - dest: A pointer to the variable where the unmarshaled response will be
//     stored.
//   - httpMethod: The HTTP method to use for the request (e.g., "GET", "POST").
//   - method: The API method to call, represented as an APIMethod type.
//   - params: The parameters to include in the API request, typically a struct
//     that can be serialized into query parameters.
//
// Returns:
//   - An error if the request fails, the response cannot be unmarshaled, or any
//     other issue occurs.
func (a *AuthAPI) RequestSigned(dest any, httpMethod string, method APIMethod, params any) error {
	err := a.CheckCredentials(RequestLevelSecret)
	if err != nil {
		return err
	}

	p, err := query.Values(params)
	if err != nil {
		return err
	}

	p.Set("api_key", a.apiKey)
	p.Set("method", method.String())
	p.Set("api_sig", a.Signature(p))

	switch httpMethod {
	case http.MethodGet:
		return a.GetURL(dest, BuildAPIURL(p))
	case http.MethodPost:
		return a.PostBody(dest, Endpoint, p.Encode())
	default:
		return errors.New("unsupported HTTP method")
	}
}

type AuthURLParams struct {
	APIKey   string `url:"api_key"`
	Callback string `url:"cb,omitempty"`
	Token    string `url:"token,omitempty"`
}

// AuthURL returns the authentication URL for the Last.fm API with the specified
// parameters. This URL can be used to send users to Last.fm for authentication.
// The URL includes the API key and a callback URL if specified.
//
// Parameters:
//   - params.APIKey: The Last.fm API key.
//   - params.Callback: The URL to redirect the user to and provided the
//     authenticated token to as a query parameter.
//   - params.Token: The token for the user to authenticate.
//
// Returns:
//   - The authentication URL as a string.
func AuthURL(params AuthURLParams) string {
	p := url.Values{}

	p.Set("api_key", params.APIKey)

	if params.Callback != "" {
		p.Set("cb", params.Callback)
	}
	if params.Token != "" {
		p.Set("token", params.Token)
	}

	return lastfm.AuthURL + "?" + p.Encode()
}

// Signature generates a signature for the given parameters using the API
// secret. The signature is created by concatenating the sorted parameter keys
// and their values, followed by the session secret. The resulting string is
// then hashed using MD5 to produce a hexadecimal representation of the hash.
func (a *AuthAPI) Signature(params url.Values) string {
	return Signature(params, a.secret)
}

// Signature generates a Last.fm API signature for the given parameters and
// secret. The signature is created by concatenating the sorted parameter keys
// and their values, followed by the secret. The resulting string is then
// hashed using MD5 to produce a hexadecimal representation of the hash.
//
// Parameters:
//   - params: The parameters to include in the signature.
//   - secret: The secret key to use for signing the request.
//
// Returns:
//   - A hexadecimal string representing the signature.
//
// https://www.last.fm/api/authspec
func Signature(params url.Values, secret string) string {
	keys := slices.Sorted(maps.Keys(params))

	var sig string
	for _, k := range keys {
		// exclude format and callback params from signature
		if k == "format" || k == "callback" {
			continue
		}

		sig += k + params.Get(k)
	}

	sig += secret
	hash := md5.Sum([]byte(sig))

	return hex.EncodeToString(hash[:])
}

// CheckCredentials verifies the authorization level required for an API
// request and ensures the necessary credentials are present. It checks the
// presence of the API secret for requests requiring a session or secret, and
// the API key for requests requiring an API key. Returns an error if required
// credentials are missing.
//
// Parameters:
//   - level: The RequestLevel indicating the level of authorization required.
//
// Returns:
//   - An error if the required authentication credentials are not present.
func (a *AuthAPI) CheckCredentials(level RequestLevel) error {
	switch level {
	case RequestLevelSession, RequestLevelSecret:
		if a.secret == "" {
			return NewLastFMError(ErrSecretRequired, SecretRequiredMessage)
		}
		fallthrough
	default:
		return a.API.CheckCredentials(level)
	}
}
