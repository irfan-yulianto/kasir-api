package middleware

import "net/http"

func APIKey(validAPIKey string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get("X-API-Key")

			if apiKey == "" {
				http.Error(w, "API key required", http.StatusUnauthorized)
				return
			}

			if apiKey != validAPIKey {
				http.Error(w, "Invalid API key", http.StatusUnauthorized)
				return
			}

			next(w, r)
		}
	}
}
