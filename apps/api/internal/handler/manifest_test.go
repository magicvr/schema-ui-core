package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestManifestIsPublicStableAndSupportsETag(t *testing.T) {
	mux := http.NewServeMux()
	if err := RegisterManifest(mux); err != nil {
		t.Fatal(err)
	}

	first := httptest.NewRecorder()
	mux.ServeHTTP(first, httptest.NewRequest(http.MethodGet, "/.well-known/schema-ui/app-manifest.json", nil))
	if first.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", first.Code)
	}
	if first.Header().Get("ETag") == "" || first.Header().Get("Content-Type") == "" {
		t.Fatalf("headers = %v", first.Header())
	}
	if got := first.Header().Get("X-Schema-UI-Manifest-Source"); got != "api" {
		t.Fatalf("manifest source = %q, want api", got)
	}

	second := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/.well-known/schema-ui/app-manifest.json", nil)
	request.Header.Set("If-None-Match", first.Header().Get("ETag"))
	mux.ServeHTTP(second, request)
	if second.Code != http.StatusNotModified {
		t.Fatalf("conditional status = %d, want 304", second.Code)
	}
}
