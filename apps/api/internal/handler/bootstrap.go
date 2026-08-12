package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
)

// RegisterBootstrap mounts the optional public bootstrap document
// (ADR-0035 / spec 10 §2.1). It is assembled from the same manifest bytes the
// manifest route serves, so the declared manifest.sha256 always matches the
// real response — the Host's integrity stage verifies production bytes.
//
// The first version always serves availability mode "normal"; maintenance /
// upgrade-required / degraded documents are exercised through the Host's
// browser-level tests (route interception), not through this endpoint.

type bootstrapDocument struct {
	BootstrapVersion     string            `json:"bootstrapVersion"`
	RequiredCapabilities []string          `json:"requiredCapabilities"`
	Manifest             bootstrapManifest `json:"manifest"`
	Availability         bootstrapMode     `json:"availability"`
}

type bootstrapManifest struct {
	URL    string `json:"url"`
	Sha256 string `json:"sha256"`
}

type bootstrapMode struct {
	Mode string `json:"mode"`
}

func RegisterBootstrap(mux *http.ServeMux, manifestBytes []byte) error {
	sum := sha256.Sum256(manifestBytes)
	document := bootstrapDocument{
		BootstrapVersion:     "1.0",
		RequiredCapabilities: []string{"host.bootstrap"},
		Manifest: bootstrapManifest{
			URL:    "/.well-known/schema-ui/app-manifest.json",
			Sha256: hex.EncodeToString(sum[:]),
		},
		Availability: bootstrapMode{Mode: "normal"},
	}
	data, err := json.Marshal(document)
	if err != nil {
		return fmt.Errorf("bootstrap: encode document: %w", err)
	}
	// The response is assembled before registration so an invalid fragment
	// fails app construction rather than becoming a silent runtime fallback.
	mux.Handle("GET /.well-known/schema-ui/host-bootstrap.json", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	}))
	return nil
}
