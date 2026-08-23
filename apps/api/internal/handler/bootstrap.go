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
// Runtime modes are projected onto the existing Host availability enum; the
// read-only mode deliberately uses the existing degraded representation.

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
	Mode                 string   `json:"mode"`
	MessageKey           string   `json:"messageKey,omitempty"`
	RetryAfterSeconds    int      `json:"retryAfterSeconds,omitempty"`
	DisabledCapabilities []string `json:"disabledCapabilities,omitempty"`
}

func RegisterBootstrap(mux *http.ServeMux, manifestBytes []byte) error {
	return RegisterBootstrapWithAvailability(mux, manifestBytes, "normal")
}

// RegisterBootstrapWithAvailability projects the backend runtime mode onto
// the existing Host availability enum. read-only intentionally uses the
// existing degraded mode; its precise distinction is carried by the status
// endpoint, not by a new protocol capability or mode.
func RegisterBootstrapWithAvailability(mux routeRegistrar, manifestBytes []byte, runtimeMode string) error {
	availability, err := bootstrapAvailability(runtimeMode)
	if err != nil {
		return err
	}
	sum := sha256.Sum256(manifestBytes)
	document := bootstrapDocument{
		BootstrapVersion:     "1.0",
		RequiredCapabilities: []string{"host.bootstrap"},
		Manifest: bootstrapManifest{
			URL:    "/.well-known/schema-ui/app-manifest.json",
			Sha256: hex.EncodeToString(sum[:]),
		},
		Availability: availability,
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

func bootstrapAvailability(runtimeMode string) (bootstrapMode, error) {
	switch runtimeMode {
	case "", "normal":
		return bootstrapMode{Mode: "normal"}, nil
	case "maintenance":
		return bootstrapMode{Mode: "maintenance"}, nil
	case "degraded", "read-only":
		return bootstrapMode{Mode: "degraded"}, nil
	default:
		return bootstrapMode{}, fmt.Errorf("bootstrap: invalid runtime mode %q", runtimeMode)
	}
}
