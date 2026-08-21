package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterBootstrapServesDeterministicDocument(t *testing.T) {
	manifestBytes := []byte(`{"protocolVersion":"2.7","requiredCapabilities":["app.manifest"]}`)
	mux := http.NewServeMux()
	if err := RegisterManifest(mux, manifestBytes); err != nil {
		t.Fatalf("RegisterManifest: %v", err)
	}
	if err := RegisterBootstrap(mux, manifestBytes); err != nil {
		t.Fatalf("RegisterBootstrap: %v", err)
	}

	response := httptest.NewRecorder()
	mux.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/.well-known/schema-ui/host-bootstrap.json", nil))
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", response.Code)
	}
	if contentType := response.Header().Get("Content-Type"); contentType != "application/json; charset=utf-8" {
		t.Fatalf("Content-Type = %q", contentType)
	}

	var document struct {
		BootstrapVersion     string   `json:"bootstrapVersion"`
		RequiredCapabilities []string `json:"requiredCapabilities"`
		Manifest             struct {
			URL    string `json:"url"`
			Sha256 string `json:"sha256"`
		} `json:"manifest"`
		Availability struct {
			Mode string `json:"mode"`
		} `json:"availability"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &document); err != nil {
		t.Fatalf("parse document: %v", err)
	}
	if document.BootstrapVersion != "1.0" {
		t.Fatalf("bootstrapVersion = %q", document.BootstrapVersion)
	}
	if len(document.RequiredCapabilities) != 1 || document.RequiredCapabilities[0] != "host.bootstrap" {
		t.Fatalf("requiredCapabilities = %v", document.RequiredCapabilities)
	}
	if document.Manifest.URL != "/.well-known/schema-ui/app-manifest.json" {
		t.Fatalf("manifest.url = %q", document.Manifest.URL)
	}
	sum := sha256.Sum256(manifestBytes)
	if document.Manifest.Sha256 != hex.EncodeToString(sum[:]) {
		t.Fatalf("manifest.sha256 = %q, want %q", document.Manifest.Sha256, hex.EncodeToString(sum[:]))
	}
	if document.Availability.Mode != "normal" {
		t.Fatalf("availability.mode = %q", document.Availability.Mode)
	}
}

func TestRegisterBootstrapOnlyGET(t *testing.T) {
	mux := http.NewServeMux()
	if err := RegisterBootstrap(mux, []byte(`{}`)); err != nil {
		t.Fatalf("RegisterBootstrap: %v", err)
	}
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/.well-known/schema-ui/host-bootstrap.json", nil))
	if response.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want 405", response.Code)
	}
}

func TestRegisterBootstrapProjectsRuntimeAvailability(t *testing.T) {
	for _, tc := range []struct {
		mode string
		want string
	}{
		{mode: "normal", want: "normal"},
		{mode: "maintenance", want: "maintenance"},
		{mode: "degraded", want: "degraded"},
		{mode: "read-only", want: "degraded"},
	} {
		t.Run(tc.mode, func(t *testing.T) {
			mux := http.NewServeMux()
			if err := RegisterBootstrapWithAvailability(mux, []byte(`{}`), tc.mode); err != nil {
				t.Fatalf("RegisterBootstrapWithAvailability: %v", err)
			}
			recorder := httptest.NewRecorder()
			mux.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/.well-known/schema-ui/host-bootstrap.json", nil))
			var document struct {
				Availability struct {
					Mode string `json:"mode"`
				} `json:"availability"`
			}
			if err := json.Unmarshal(recorder.Body.Bytes(), &document); err != nil {
				t.Fatal(err)
			}
			if document.Availability.Mode != tc.want {
				t.Fatalf("availability.mode = %q, want %q", document.Availability.Mode, tc.want)
			}
		})
	}
	if err := RegisterBootstrapWithAvailability(http.NewServeMux(), []byte(`{}`), "invalid"); err == nil {
		t.Fatal("invalid runtime mode must fail closed")
	}
}
