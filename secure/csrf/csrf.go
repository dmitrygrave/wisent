package csrf

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/dmitrygrave/wisent/secure/helpers"
	"github.com/dmitrygrave/wisent/utils/logging"
	"github.com/gorilla/securecookie"
)

// Handler wraps the default handler and provides CSRF protection
type Handler struct {
	// Handlers
	success http.Handler
	failure http.Handler
	// The base cookie
	baseCookie *cookie
	// A list of paths which will are exempt from CSRF protection
	exemptPaths []string
	// A list of methods which are exempt from CSRF protection
	exemptMethods []string
}

const (
	tokenLength = 32
)

const (
	// CookieName is the name of the CSRF cookie
	CookieName = "csrf_token"
	// FormFieldName is the name of the hidden CSRF form field
	FormFieldName = "csrf_token"
	// HeaderName is the name of the CSRF header
	HeaderName = "X-CSRF-Token"
	// MaxAge is the default MaxAge cookie value used
	MaxAge = 365 * 24 * 60 * 60
)

var (
	// ExemptMethods are the default methods which are exempt from CSRF protection
	ExemptMethods = []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodConnect,
		http.MethodTrace,
	}

	// Errors

	// ErrServerIssue is an error which is returned whenever there is a
	// server error processing the request
	ErrServerIssue = errors.New("There was an error handling the CSRF protected request")

	// ErrNoReferer is an error which is returned whenever a secure request
	// comes through without a valid referer header.
	ErrNoReferer = errors.New("A secure request was received with an empty or malformed Referer")

	// ErrBadReferer is an error which is returned whenever a secure request
	// comes through with a mismatched referer header / request URL
	ErrBadReferer = errors.New("A secure request's Referer comes from a different Origin")

	// ErrNoToken is an error which is returned whenever a request comes in
	// without a valid token in either the header, form, or multipart
	ErrNoToken = errors.New("No CSRF token found in the request")

	// ErrBadToken is an error which is returned whenever a token comes in
	// which does not match the CSRF cookie
	ErrBadToken = errors.New("Invalid CSRF token received")
)

func defaultFailureHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(400), 400)
}

// New initializes a new CSRFHandler
func New(handler http.Handler) *Handler {
	csrf := &Handler{
		success:       handler,
		failure:       http.HandlerFunc(defaultFailureHandler),
		baseCookie:    nil,
		exemptMethods: ExemptMethods,
	}

	return csrf
}

// SetCookieOptions sets the cookie
func (handler *Handler) SetCookieOptions(name string, maxAge int, secure bool, httpOnly bool, path string, domain string, authKey []byte) {
	handler.baseCookie = &cookie{
		name:         name,
		maxAge:       maxAge,
		secure:       secure,
		httpOnly:     httpOnly,
		path:         path,
		domain:       domain,
		secureCookie: securecookie.New(authKey, nil),
	}
}

// SetFailureHandler sets the handler to call when there is an error in the csrf
// check
func (handler *Handler) SetFailureHandler(failureHandler http.Handler) {
	handler.failure = failureHandler
}

// ServeHTTP wraps around the http.Handler and validates that the CSRF token
// is valid. If it is then Handler.success is called, if noot Handler.failure
// is called
func (handler *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r = addContext(r)

	w.Header().Add("Vary", "Cookie")

	realToken, err := handler.baseCookie.Get(r)

	// If there is an error (cookie not set or HMAC failure) we need to
	// generate a new token. This will fail further down when equality checks
	// are performed.
	if err != nil || len(realToken) != tokenLength {
		realToken, err = helpers.GenerateSecureBytes(tokenLength)

		if err != nil {
			logging.Error("Error generating random bytes for CSRF token")
			setReqErrDesc(r, ErrServerIssue)
			handler.handleFailure(w, r)
			return
		}

		// Set the cookie with the new token
		err = handler.baseCookie.Set(realToken, w)

		if err != nil {
			logging.Error("Error setting CSRF cookie")
			setReqErrDesc(r, ErrServerIssue)
			handler.handleFailure(w, r)
			return
		}
	}

	setReqToken(r, maskToken(realToken, r))

	// TODO implement exemptions
	if stringInSlice(handler.exemptMethods, r.Method) {
		// If it's a safe method according to RFC7231#section-4.2.1
		// then we can stop here and pass to the success handler
		handler.handleSuccess(w, r)
		return
	}

	// Check if the request is secure, then enforce origin checks
	if r.URL.Scheme == "https" {
		referer, err := url.Parse(r.Header.Get("Referer"))

		// If the referer is empty then it's probably unspecified
		if err != nil {
			setReqErrDesc(r, ErrNoReferer)
			handler.handleFailure(w, r)
			return
		}

		if !sameOrigin(referer, r.URL) {
			setReqErrDesc(r, ErrNoReferer)
			handler.handleFailure(w, r)
			return
		}
	}

	// Get token from request
	sentToken, err := getTokenFromRequest(r)

	if err != nil {
		setReqErrDesc(r, ErrNoToken)
		handler.handleFailure(w, r)
		return
	}

	unmaskedToken := unmaskToken(sentToken)

	if !compareTokens(unmaskedToken, realToken) {
		setReqErrDesc(r, ErrBadToken)
		handler.handleFailure(w, r)
		return
	}

	// Prevent clients from caching the response
	w.Header().Add("Vary", "Cookie")

	handler.handleSuccess(w, r)
}

// handleSuccess calls the sucess http.Handler and should be called after
// setting the request context
func (handler *Handler) handleSuccess(w http.ResponseWriter, r *http.Request) {
	handler.success.ServeHTTP(w, r)
}

// handleFailure calls the failure http.Handler and should be called after
// settings the request error description and setting the request context
func (handler *Handler) handleFailure(w http.ResponseWriter, r *http.Request) {
	handler.failure.ServeHTTP(w, r)
}
