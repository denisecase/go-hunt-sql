// Package middleware provides common funcationality for requests
package middleware

import (
	"net/http"

	"github.com/denisecase/go-hunt-sql/api/auth"
)

// SetMiddlewareJSON formats responses as JSON
func SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

// SetMiddlewareAuthentication checks authentication token
func SetMiddlewareAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := auth.TokenValid(r)
		if err != nil {
			return
		}
		next(w, r)
	}
}
