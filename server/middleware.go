package server

import (
	"net/http"
)

// middleware to perform authentication for a request
func MiddlewareAPIKeys(apikeyCheck ApiKeyChecker) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// check if API key provided
			if r.URL.Query().Has("apikey") {
				if apikey := r.URL.Query().Get("apikey"); apikeyCheck.IsValidApiKey(apikey) {
					next.ServeHTTP(w, r)
					return
				}
				http.Error(w, "invalid api key", http.StatusUnauthorized)
				return
			} else {
				http.Error(w, "no api key provided", http.StatusUnauthorized)
			}
		})
	}
}
