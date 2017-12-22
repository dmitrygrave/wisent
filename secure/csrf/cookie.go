package csrf

import (
	"net/http"
	"time"

	"github.com/gorilla/securecookie"
)

// cookie is a secure cookie where the CSRF cookie will be stored
type cookie struct {
	name         string
	maxAge       int
	secure       bool
	httpOnly     bool
	path         string
	domain       string
	secureCookie *securecookie.SecureCookie
}

// Get returns the CSRF token from the cookie. If no token can be found it
// returns nil
func (c *cookie) Get(r *http.Request) ([]byte, error) {
	// Get the cookie
	cookie, err := r.Cookie(c.name)

	if err != nil {
		return nil, err
	}

	token := make([]byte, tokenLength)

	// Decode the cookie
	err = c.secureCookie.Decode(c.name, cookie.Value, &token)

	if err != nil {
		return nil, err
	}

	return token, nil
}

func (c *cookie) Set(token []byte, w http.ResponseWriter) error {
	// Encode the CSRF token
	encToken, err := c.secureCookie.Encode(c.name, token)

	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:     c.name,
		Value:    encToken,
		MaxAge:   c.maxAge,
		HttpOnly: c.httpOnly,
		Secure:   c.secure,
		Path:     c.path,
		Domain:   c.domain,
	}

	// Set the max age value on the cookie
	if c.maxAge > 0 {
		cookie.Expires = time.Now().Add(time.Duration(c.maxAge) * time.Second)
	}

	http.SetCookie(w, cookie)

	return nil
}
