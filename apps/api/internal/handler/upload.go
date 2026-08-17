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
	"regexp"
	"strings"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
)

// maxUploadBytes bounds a single multipart upload (8 MiB); the server rejects
// oversize files with FILE_TOO_LARGE regardless of client-side constraints
// (ADR-0012 D2: server must independently validate).
const maxUploadBytes = 8 << 20

// uploadPolicy carries the W7-configurable upload limits. The old package-level
// UPLOAD_* environment reads moved into config.Config (single configuration
// authority); RegisterUpload accepts optional UploadOption overrides and keeps
// the historical defaults when none are given, so legacy callers are
// unaffected.
type uploadPolicy struct {
	// allowedTypes, when non-empty (comma-separated MIME types), makes the
	// server independently reject other types with UNSUPPORTED_FILE_TYPE. The
	// allow-list is checked against the server-detected content type, never the
	// client-declared header.
	allowedTypes string
	// Per-user upload quotas (W4 P0-2): the permission gate stops a
	// low-privilege viewer from filling the disk, but a compromised
	// files.write holder could still pump 8 MiB objects indefinitely. These
	// best-effort limits bound each user's stored file count and total bytes.
	// They are enforced by scanning the upload directory's owner meta
	// (owner-only contract already guarantees every stored object records its
	// uploader); a corrupt/unreadable meta entry counts toward the total
	// conservatively so a failed read cannot bypass the quota.
	maxUserFiles int
	maxUserBytes int
}

// UploadOption configures the upload policy; zero or more may be passed to
// RegisterUpload.
type UploadOption func(*uploadPolicy)

// WithAllowedTypes sets the comma-separated MIME allow-list.
func WithAllowedTypes(types string) UploadOption {
	return func(p *uploadPolicy) { p.allowedTypes = strings.TrimSpace(types) }
}

// WithUserLimits sets the per-user file count and total byte quotas. A
// non-positive value keeps the caller-supplied default (never disables the
// server-side bound unintentionally).
func WithUserLimits(maxFiles, maxBytes int) UploadOption {
	return func(p *uploadPolicy) {
		if maxFiles > 0 {
			p.maxUserFiles = maxFiles
		}
		if maxBytes > 0 {
			p.maxUserBytes = maxBytes
		}
	}
}

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
	dir    string
	policy uploadPolicy
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

// uploadFileIDPattern is the only shape save() ever writes (16 random bytes as
// lowercase hex). load() must reject everything else so a crafted PathValue
// cannot turn filepath.Join into a volume-relative path (Windows `C:name`)
// or `..` / reserved names.
var uploadFileIDPattern = regexp.MustCompile(`^[0-9a-f]{32}$`)

func (s *uploadStore) load(id string) ([]byte, map[string]string, error) {
	if !uploadFileIDPattern.MatchString(id) {
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

// quotaReached reports whether storing one more object (with the given size)
// for ownerID would exceed the per-user file count or total byte quota (W4
// P0-2). It scans the upload directory's owner meta; a corrupt or unreadable
// meta entry still counts toward the total conservatively, so a failed read
// cannot bypass the limit. Scanning is O(files) per upload — acceptable for an
// admin tool where the permission gate already bounds who can upload.
func (s *uploadStore) quotaReached(ownerID string, nextSize int) (reason string, reached bool) {
	if s.policy.maxUserFiles <= 0 && s.policy.maxUserBytes <= 0 {
		return "", false // quotas disabled
	}
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		// Dir absent = no stored files; a real read error fails closed (count
		// toward quota rather than silently allowing an unbounded store).
		if os.IsNotExist(err) {
			return "", false
		}
		return "quota unavailable", true
	}
	files := 0
	bytes := 0
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".meta.json") {
			continue
		}
		raw, err := os.ReadFile(filepath.Join(s.dir, entry.Name()))
		if err != nil {
			files++
			bytes += maxUploadBytes // conservative: unreadable meta counts max
			continue
		}
		meta := map[string]string{}
		if err := json.Unmarshal(raw, &meta); err != nil {
			files++
			bytes += maxUploadBytes
			continue
		}
		if meta["owner"] != ownerID {
			continue
		}
		// The stored object is the meta name minus the ".meta.json" suffix;
		// stat it for the real byte count so the byte quota is accurate.
		fileSize := 1 << 20 // nominal minimum when stat fails
		if info, err := os.Stat(filepath.Join(s.dir, strings.TrimSuffix(entry.Name(), ".meta.json"))); err == nil {
			fileSize = int(info.Size())
		}
		files++
		bytes += fileSize
	}
	if s.policy.maxUserFiles > 0 && files+1 > s.policy.maxUserFiles {
		return "per-user file count quota exceeded", true
	}
	if s.policy.maxUserBytes > 0 && bytes+nextSize > s.policy.maxUserBytes {
		return "per-user byte quota exceeded", true
	}
	return "", false
}

// RegisterUpload mounts the upload/file endpoints behind the auth middleware
// and the files.write permission gate (W4 P0-2): any authenticated user can
// read their own downloads (owner-only), but storing a new object requires the
// files.write key — default-held by admin only, so a low-privilege viewer
// account cannot fill the disk. The gate mirrors the resource-factory
// permission model (requirePermission, fail-closed 403 for authenticated users
// without the key).
func RegisterUpload(mux *http.ServeMux, a authMiddleware, dir string, opts ...UploadOption) {
	store := &uploadStore{
		dir: dir,
		policy: uploadPolicy{
			// Historical defaults (pre-W7 package-level env values).
			maxUserFiles: 1000,
			maxUserBytes: 256 << 20,
		},
	}
	for _, o := range opts {
		o(&store.policy)
	}
	mux.Handle("POST /api/upload", a.Middleware(uploadPermissionGate(store.upload())))
	mux.Handle("GET /api/files/{id}", a.Middleware(store.file()))
}

// uploadPermissionGate wraps the upload handler with the files.write
// permission check. It is an http.Handler so it composes with the auth
// middleware exactly like the resource factory's permission gates.
func uploadPermissionGate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := requirePermission(w, r, "files.write"); !ok {
			return
		}
		next.ServeHTTP(w, r)
	})
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
		if s.policy.allowedTypes != "" {
			allowed := strings.Split(s.policy.allowedTypes, ",")
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
		// W4 P0-2: per-user quota gate before storage. Combined with the
		// files.write permission gate, a single account cannot fill the disk.
		if reason, reached := s.quotaReached(user.ID, len(body)); reached {
			writeLocalizedError(w, r, http.StatusRequestEntityTooLarge, "UPLOAD_QUOTA_EXCEEDED", "upload rejected: "+reason)
			return
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
