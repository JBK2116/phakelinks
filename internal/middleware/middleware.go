package middleware

import (
	"net/http"
	"strings"
)

// StripTrailingSlashMiddleware() removes the trailing `/` from all incoming request URLs
func StripTrailingSlashMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		url := *request.URL
		url.Path = strings.TrimSuffix(request.URL.Path, "/")
		// prevent "" paths to minimize bugs in production
		if url.Path == "" {
			url.Path = "/"
		}
		request.URL = &url
		handler.ServeHTTP(writer, request)
	})
}
