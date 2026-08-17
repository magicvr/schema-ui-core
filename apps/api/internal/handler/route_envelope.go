package handler

import (
	"net/http"
	"strings"
)

// WithJSONRouteErrors maps unregistered paths and method mismatches onto the
// cataloged JSON error envelope (W15-F04) instead of Go's text/plain 404/405.
func WithJSONRouteErrors(mux *http.ServeMux) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, pattern := mux.Handler(r)
		if pattern == "" || isMethodMismatch(mux, r, pattern) {
			if pathHasOtherMethod(mux, r) {
				writeLocalizedError(w, r, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
				return
			}
			if pattern == "" {
				writeLocalizedError(w, r, http.StatusNotFound, "NOT_FOUND", "not found")
				return
			}
		}
		mux.ServeHTTP(w, r)
	})
}

func isMethodMismatch(mux *http.ServeMux, r *http.Request, pattern string) bool {
	if method, _, ok := strings.Cut(pattern, " "); ok && method != r.Method {
		return pathHasOtherMethod(mux, r)
	}
	return false
}

func pathHasOtherMethod(mux *http.ServeMux, r *http.Request) bool {
	probe := r.Clone(r.Context())
	for _, method := range []string{
		http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut,
		http.MethodPatch, http.MethodDelete, http.MethodOptions,
	} {
		if method == r.Method {
			continue
		}
		probe.Method = method
		if _, pattern := mux.Handler(probe); pattern != "" {
			return true
		}
	}
	return false
}
