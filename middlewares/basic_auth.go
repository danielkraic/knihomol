package middlewares

import (
	"encoding/base64"
	"net/http"
	"strings"
)

// AuthenticationMiddleware basic auth middleware
type AuthenticationMiddleware struct {
	username string
	password string
}

// NewAuthenticationMiddleware creates new AuthenticationMiddleware
func NewAuthenticationMiddleware(username, password string) AuthenticationMiddleware {
	return AuthenticationMiddleware{
		username: username,
		password: password,
	}
}

// Middleware middleware to provide basic auth
func (amw *AuthenticationMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if amw.username == "" && amw.password == "" {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)

		auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(auth) != 2 || auth[0] != "Basic" {
			http.Error(w, "Not authorized", 401)
			return
		}

		payload, err := base64.StdEncoding.DecodeString(auth[1])
		if err != nil {
			http.Error(w, "authorization failed", http.StatusBadRequest)
			return
		}
		pair := strings.SplitN(string(payload), ":", 2)

		if len(pair) != 2 || amw.username != pair[0] || amw.password != pair[1] {
			http.Error(w, "Not authorized", 401)
			return
		}

		next.ServeHTTP(w, r)
	})
}
