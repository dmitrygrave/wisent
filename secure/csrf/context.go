package csrf

import (
	"context"
	"encoding/base64"
	"net/http"

	"github.com/dmitrygrave/wisent/utils/logging"
)

// contextKey defines a string type to use as the context key since it is
// recommended not to use any default types
type contextKey string

const csrfContextKey = contextKey("csrf-key")

// csrfContext contains either the csrf token or an error describing why
// the csrf validation failed
type csrfContext struct {
	token   string
	errDesc error
}

// addContext adds the CSRF context to the request
func addContext(r *http.Request) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), csrfContextKey, &csrfContext{}))
}

// GetReqToken returns the CSRF token for the provided request or an empty
// string if none is found
func GetReqToken(r *http.Request) string {
	context := r.Context().Value(csrfContextKey).(*csrfContext)

	return context.token
}

// GetReqErrDesc returns the error describing why the CSRF validation failed
func GetReqErrDesc(r *http.Request) error {
	context := r.Context().Value(csrfContextKey).(*csrfContext)

	return context.errDesc
}

// setReqToken sets the CSRF token on the request
func setReqToken(r *http.Request, token string) {
	context := r.Context().Value(csrfContextKey).(*csrfContext)
	context.token = base64.URLEncoding.EncodeToString([]byte(token))
}

// setReqErrDesc sets the error describing why the CSRF validation failed
func setReqErrDesc(r *http.Request, reason error) {
	context := r.Context().Value(csrfContextKey).(*csrfContext)

	if context.token == "" {
		logging.Panic("An error description should never be set when there is no token in the context")
	}

	context.errDesc = reason
}
