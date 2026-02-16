package middleware

import (
	"encoding/json"
	"net/http"
	"strings"
)

// JSONErrors intercepts http.Error() calls (which output text/plain) and
// automatically converts them to JSON format: {"error": "message"}.
// Handlers don't need to change — this middleware handles the conversion.
func JSONErrors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jw := &jsonErrorWriter{ResponseWriter: w}
		next.ServeHTTP(jw, r)
	})
}

type jsonErrorWriter struct {
	http.ResponseWriter
	code        int
	intercepted bool
}

func (w *jsonErrorWriter) WriteHeader(code int) {
	if code >= 400 {
		ct := w.Header().Get("Content-Type")
		if ct == "text/plain; charset=utf-8" {
			w.code = code
			w.intercepted = true
			return
		}
	}
	w.ResponseWriter.WriteHeader(code)
}

func (w *jsonErrorWriter) Write(b []byte) (int, error) {
	if w.intercepted {
		msg := strings.TrimSpace(string(b))
		escaped, _ := json.Marshal(msg)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Del("X-Content-Type-Options")
		w.ResponseWriter.WriteHeader(w.code)
		w.intercepted = false
		return w.ResponseWriter.Write([]byte(`{"error":` + string(escaped) + `}`))
	}
	return w.ResponseWriter.Write(b)
}
