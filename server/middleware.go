package server

import (
	"net/http"
)

// middleware to perform authentication for a request
func AuthMiddleware(authchecker AuthChecker) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// check if API key provided
			if r.URL.Query().Has("apikey") {
				if apikey := r.URL.Query().Get("apikey"); authchecker.CheckApiKey(apikey) {
					next.ServeHTTP(w, r)
					return
				}
				http.Error(w, "invalid api key", http.StatusUnauthorized)
				return
			}

			if username, password, ok := r.BasicAuth(); ok {
				if authchecker.CheckBasic(username, password) {
					next.ServeHTTP(w, r)
					return
				} else {
					http.Error(w, "invalid credentials", http.StatusUnauthorized)
					return
				}
			}
			w.Header().Set("WWW-Authenticate", `Basic realm="restricted"`)
			w.WriteHeader(http.StatusUnauthorized)
		})
	}
}
