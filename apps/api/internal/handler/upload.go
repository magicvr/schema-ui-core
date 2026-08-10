// Upload endpoint family (I-PROTO-FULL-001 · D-UPLOAD include).
//
// Server-side contract for docs/07-actions-contract.md §7.2: POST /api/upload
// accepts one multipart file part (any field name; the Renderer always uses
// "file" by default) and responds with `{url, id, name, size}` — the client
// takes url (priority) or id as the field value. Errors use the semantic
// codes suggested by §7.3 (FILE_TOO_LARGE / UNSUPPORTED_FILE_TYPE /
// STORAGE_UNAVAILABLE) inside the frozen {error, message} envelope.
//
// GOAL-003: every stored object records the authenticated uploader as owner;
// GET /api/files/{id} is owner-only (fail-closed when meta.owner is missing).
package handler

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
)

// maxUploadBytes bounds a single multipart upload (8 MiB); the server rejects
// oversize files with FILE_TOO_LARGE regardless of client-side constraints
// (ADR-0012 D2: server must independently validate).
const maxUploadBytes = 8 << 20

// uploadAllowedTypes, when non-empty (comma-separated MIME types), makes the
// server independently reject other types with UNSUPPORTED_FILE_TYPE. The
// allow-list is checked against the server-detected content type, never the
// client-declared header.
var uploadAllowedTypes = strings.TrimSpace(os.Getenv("UPLOAD_ALLOWED_TYPES"))

// dangerousInlineTypes are MIME types browsers may render inline as active
// (script-bearing) content. They are rejected regardless of any allow-list:
// serving them from the API's same origin would let any authenticated user
// store and deliver script that runs in the context of every logged-in admin
// (stored XSS → refresh-token theft).
var dangerousInlineTypes = map[string]bool{
	"text/html":             true,
	"application/xhtml+xml": true,
	"image/svg+xml":         true,
}

// activeContentMarkers are byte sequences that indicate script-bearing or
// SVG markup regardless of what http.DetectContentType reports (A-002 F-001):
// DetectContentType sniffs SVG as text/plain and HTML can be smuggled past the
// first 512 bytes or behind a GIF89a header. Matching is case-insensitive
// because HTML/SVG tag names are case-insensitive (A-003 N-001).
var activeContentMarkers = []string{
	"<svg",
	"<script",
	"<?xml",
}

// containsActiveContent reports whether body contains any active-content
// marker, case-insensitively. A hard rejection gate on top of the sniffed MIME.
func containsActiveContent(body []byte) bool {
	lower := bytes.ToLower(body)
	for _, marker := range activeContentMarkers {
		if bytes.Contains(lower, []byte(marker)) {
			return true
		}
	}
	return false
}

type uploadStore struct {
	dir string
}

func (s *uploadStore) save(name string, contentType string, ownerID string, body []byte) (string, error) {
	idBytes := make([]byte, 16)
	if _, err := rand.Read(idBytes); err != nil {
		return "", err
	}
	id := hex.EncodeToString(idBytes)
	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return "", err
	}
	if err := os.WriteFile(filepath.Join(s.dir, id), body, 0o644); err != nil {
		return "", err
	}
	meta := map[string]string{"name": name, "type": contentType, "owner": ownerID}
	raw, err := json.Marshal(meta)
	if err == nil {
		_ = os.WriteFile(filepath.Join(s.dir, id+".meta.json"), raw, 0o644)
	}
	return id, nil
}

func (s *uploadStore) load(id string) ([]byte, map[string]string, error) {
	if id == "" || strings.ContainsAny(id, `/\`) {
		return nil, nil, os.ErrNotExist
	}
	body, err := os.ReadFile(filepath.Join(s.dir, id))
	if err != nil {
		return nil, nil, err
	}
	meta := map[string]string{}
	raw, err := os.ReadFile(filepath.Join(s.dir, id+".meta.json"))
	if err == nil {
		_ = json.Unmarshal(raw, &meta)
	}
	return body, meta, nil
}

// RegisterUpload mounts the upload/file endpoints behind the auth middleware.
func RegisterUpload(mux *http.ServeMux, a authMiddleware, dir string) {
	store := &uploadStore{dir: dir}
	mux.Handle("POST /api/upload", a.Middleware(store.upload()))
	mux.Handle("GET /api/files/{id}", a.Middleware(store.file()))
}

// authMiddleware is the minimal surface the upload handler needs; the concrete
// *auth.Authenticator satisfies it.
type authMiddleware interface {
	Middleware(next http.Handler) http.Handler
}

func (s *uploadStore) upload() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := auth.IdentityFrom(r.Context())
		if !ok || strings.TrimSpace(user.ID) == "" {
			// Middleware should always inject identity; fail closed if it did not.
			writeLocalizedError(w, r, http.StatusUnauthorized, "UNAUTHENTICATED", "no active session")
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, maxUploadBytes)
		file, header, err := r.FormFile("file")
		if err != nil {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_UPLOAD", "expected a multipart file part named file")
			return
		}
		defer file.Close()
		if header.Size <= 0 {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_FILE", "empty files are rejected")
			return
		}
		if header.Size > maxUploadBytes {
			writeLocalizedError(w, r, http.StatusRequestEntityTooLarge, "FILE_TOO_LARGE", "file exceeds the server size limit")
			return
		}
		body, err := io.ReadAll(file)
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "STORAGE_UNAVAILABLE", "could not read upload")
			return
		}
		if len(body) == 0 {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_FILE", "empty files are rejected")
			return
		}
		// The server sniffs the actual content type; the client-declared MIME is
		// never trusted for storage or serving decisions (stored-XSS hardening).
		detected := http.DetectContentType(body)
		base := detected
		if i := strings.IndexByte(base, ';'); i >= 0 {
			base = strings.TrimSpace(base[:i])
		}
		// A-002 F-001: DetectContentType alone misses SVG (sniffed as text/plain)
		// and header-smuggled HTML; the active-content marker check is the hard
		// rejection gate for script/SVG-bearing bodies regardless of the sniffed
		// type. This is defense-in-depth on top of the download headers.
		if dangerousInlineTypes[base] || containsActiveContent(body) {
			writeLocalizedError(w, r, http.StatusUnsupportedMediaType, "UNSUPPORTED_FILE_TYPE", "file type is not allowed")
			return
		}
		if uploadAllowedTypes != "" {
			allowed := strings.Split(uploadAllowedTypes, ",")
			matched := false
			for _, entry := range allowed {
				if strings.EqualFold(strings.TrimSpace(entry), base) {
					matched = true
					break
				}
			}
			if !matched {
				writeLocalizedError(w, r, http.StatusUnsupportedMediaType, "UNSUPPORTED_FILE_TYPE", "file type is not allowed")
				return
			}
		}
		id, err := s.save(header.Filename, detected, user.ID, body)
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "STORAGE_UNAVAILABLE", "could not store upload")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"id":   id,
			"name": header.Filename,
			"size": len(body),
			"url":  fmt.Sprintf("/api/files/%s", id),
		})
	})
}

func (s *uploadStore) file() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := auth.IdentityFrom(r.Context())
		if !ok || strings.TrimSpace(user.ID) == "" {
			writeLocalizedError(w, r, http.StatusUnauthorized, "UNAUTHENTICATED", "no active session")
			return
		}
		id := r.PathValue("id")
		body, meta, err := s.load(id)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				writeLocalizedError(w, r, http.StatusNotFound, "FILE_NOT_FOUND", "no file with that id")
				return
			}
			writeLocalizedError(w, r, http.StatusInternalServerError, "STORAGE_UNAVAILABLE", "could not read file")
			return
		}
		// GOAL-003: owner-only download. Missing owner (legacy/corrupt meta) fails
		// closed — never treat an unowned object as world-readable among authed users.
		owner := strings.TrimSpace(meta["owner"])
		if owner == "" || owner != user.ID {
			writeLocalizedError(w, r, http.StatusForbidden, "FORBIDDEN", "not the owner of this file")
			return
		}
		contentType := meta["type"]
		if contentType == "" {
			contentType = "application/octet-stream"
		}
		w.Header().Set("Content-Type", contentType)
		w.Header().Set("X-Content-Type-Options", "nosniff")
		// Force a download instead of inline rendering: even a type that slipped
		// past detection cannot execute in the API's same-origin context. The
		// filename is a fixed literal — the stored name is attacker-controlled
		// and must never reach a header.
		w.Header().Set("Content-Disposition", `attachment; filename="download"`)
		w.Header().Set("Content-Security-Policy", "sandbox")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(body)))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(body)
	})
}
