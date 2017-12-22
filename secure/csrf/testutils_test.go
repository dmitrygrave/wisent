package csrf

import (
	"net/http"
	"time"
)

// GetTestHandler generates a dummy http.Handler for testing the middleware
func GetTestHandler() http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("Success"))
	}

	return http.HandlerFunc(fn)
}

func SetDummyCookie(handler *Handler) {
	name := CookieName
	maxAge := int(30 * time.Minute)
	secure := true
	httpOnly := true
	path := "/"
	domain := ".dummy.url"
	authKey := "supersecure"

	handler.SetCookieOptions(name, maxAge, secure, httpOnly, path, domain, authKey)
}
