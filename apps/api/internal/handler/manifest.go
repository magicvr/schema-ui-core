package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/magicvr/schema-ui-core/apps/api/internal/manifest"
)

// RegisterManifest mounts the public, same-origin protocol endpoint. The
// response is assembled before registration so an invalid fragment fails app
// construction rather than becoming a silent runtime fallback.

func RegisterManifest(mux *http.ServeMux, supplied ...[]byte) error {
	if len(supplied) > 1 {
		return fmt.Errorf("manifest: at most one prebuilt document is allowed")
	}
	var data []byte
	if len(supplied) == 1 {
		data = append([]byte(nil), supplied[0]...)
	} else {
		var err error
		data, err = manifest.Default()
		if err != nil {
			return err
		}
	}
	sum := sha256.Sum256(data)
	etag := `"` + hex.EncodeToString(sum[:]) + `"`
	mux.Handle("GET /.well-known/schema-ui/app-manifest.json", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("ETag", etag)
		if r.Header.Get("If-None-Match") == etag {
			w.WriteHeader(http.StatusNotModified)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	}))
	return nil
}
