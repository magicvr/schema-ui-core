// Package objectstore holds the adapters for the kernel object-storage port
// (VP-014 / workspace-014 GOAL-002 D-001). LocalStore is the default
// implementation: a single root directory whose per-namespace subdirectories
// reproduce the historical layout (<db dir>/avatars, /brand-assets,
// /uploads) byte for byte, so adopting the port requires zero migration.
//
// On-disk contract (frozen): each object is <root>/<ns>/<id> plus a JSON
// sidecar <root>/<ns>/<id>.meta.json carrying only the non-empty ObjectMeta
// keys. New writes therefore stay compatible with both legacy shapes
// (uploads name/type/owner; rasters type/kind/owner), and reads tolerate a
// missing or corrupt sidecar exactly like the historical loaders did.
package objectstore

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
)

// LocalStore implements kernel.ObjectStore on the local filesystem.
// It carries no mutexes: every method touches only path-derived state and
// the OS provides the synchronization, matching the historical stores.
type LocalStore struct {
	root string
}

// compile-time proof the default adapter satisfies the frozen port.
var _ kernel.ObjectStore = (*LocalStore)(nil)

// NewLocal constructs the local disk adapter rooted at root (the namespace
// subdirectories live directly beneath it). The caller decides the root:
// composition derives it from filepath.Dir(cfg.DBPath) unless an explicit
// storage.objects.local.root override is configured.
func NewLocal(root string) *LocalStore {
	return &LocalStore{root: root}
}

// Root reports the storage root (diagnostics/tests).
func (s *LocalStore) Root() string { return s.root }

// nsDir returns the namespace directory after fail-closed validation.
func (s *LocalStore) nsDir(ns kernel.ObjectNamespace) (string, error) {
	if !kernel.ValidObjectNamespace(ns) {
		return "", fmt.Errorf("objectstore: unknown namespace %q", string(ns))
	}
	return filepath.Join(s.root, string(ns)), nil
}

// objectPath validates ns and id, then returns the body and sidecar paths.
// The id shape check happens before any join so a crafted id can never
// escape the namespace directory (same guard as uploadFileIDPattern).
func (s *LocalStore) objectPaths(ns kernel.ObjectNamespace, id string) (body, side string, err error) {
	if !kernel.ValidObjectNamespace(ns) {
		return "", "", fmt.Errorf("objectstore: unknown namespace %q", string(ns))
	}
	if !kernel.ValidObjectID(id) {
		return "", "", fmt.Errorf("objectstore: invalid object id %q", id)
	}
	dir := filepath.Join(s.root, string(ns))
	return filepath.Join(dir, id), filepath.Join(dir, id+".meta.json"), nil
}

// metaFile renders the sidecar record with empty keys omitted, so new files
// are byte-compatible with the historical upload {name,type,owner} and
// raster {type,kind,owner} shapes.
func metaFile(meta kernel.ObjectMeta) ([]byte, error) {
	m := map[string]string{}
	for k, v := range map[string]string{
		"name":  meta.Name,
		"type":  meta.Type,
		"kind":  meta.Kind,
		"owner": meta.Owner,
	} {
		if v != "" {
			m[k] = v
		}
	}
	if len(m) == 0 {
		return []byte("{}"), nil
	}
	return json.Marshal(m)
}

// parseMeta reads a sidecar tolerantly: missing or corrupt files yield an
// empty meta instead of failing the read (historical load() behavior — the
// owner-only policy at the HTTP surface already fails closed on empty owner).
func parseMeta(path string) kernel.ObjectMeta {
	meta := kernel.ObjectMeta{}
	raw, err := os.ReadFile(path)
	if err != nil {
		return meta
	}
	var m map[string]string
	if err := json.Unmarshal(raw, &m); err != nil {
		return meta
	}
	meta.Name, meta.Type, meta.Kind, meta.Owner = m["name"], m["type"], m["kind"], m["owner"]
	return meta
}

// Put stores body under (ns, id), replacing any previous object (upsert),
// then writes the metadata sidecar. The body lands via a temp file + rename
// inside the same directory (atomic on all supported platforms); if the
// sidecar cannot be written the body is removed again so a successful-looking
// return never leaves an invisible object (W7 F-013 precedent).
func (s *LocalStore) Put(_ context.Context, ns kernel.ObjectNamespace, id string, body []byte, meta kernel.ObjectMeta) error {
	bodyPath, sidePath, err := s.objectPaths(ns, id)
	if err != nil {
		return err
	}
	dir := filepath.Dir(bodyPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(dir, ".put-"+id+"-*.tmp")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer func() {
		_ = os.Remove(tmpName) // no-op after a successful rename
	}()
	if _, err := tmp.Write(body); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Chmod(tmpName, 0o644); err != nil {
		return err
	}
	if err := os.Rename(tmpName, bodyPath); err != nil {
		return err
	}
	raw, err := metaFile(meta)
	if err == nil {
		err = os.WriteFile(sidePath, raw, 0o644)
	}
	if err != nil {
		// Roll the object back: a stored body without its meta would be
		// invisible to quota scans and owner checks (fail closed).
		_ = os.Remove(bodyPath)
		return err
	}
	return nil
}

// Get returns the object body and metadata, or kernel.ErrObjectNotFound.
func (s *LocalStore) Get(_ context.Context, ns kernel.ObjectNamespace, id string) ([]byte, kernel.ObjectMeta, error) {
	bodyPath, sidePath, err := s.objectPaths(ns, id)
	if err != nil {
		return nil, kernel.ObjectMeta{}, err
	}
	body, err := os.ReadFile(bodyPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, kernel.ObjectMeta{}, fmt.Errorf("objectstore: get %s/%s: %w", string(ns), id, kernel.ErrObjectNotFound)
		}
		return nil, kernel.ObjectMeta{}, err
	}
	return body, parseMeta(sidePath), nil
}

// Stat returns identity, size and metadata without reading the body.
func (s *LocalStore) Stat(_ context.Context, ns kernel.ObjectNamespace, id string) (kernel.ObjectInfo, error) {
	bodyPath, sidePath, err := s.objectPaths(ns, id)
	if err != nil {
		return kernel.ObjectInfo{}, err
	}
	info, err := os.Stat(bodyPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return kernel.ObjectInfo{}, fmt.Errorf("objectstore: stat %s/%s: %w", string(ns), id, kernel.ErrObjectNotFound)
		}
		return kernel.ObjectInfo{}, err
	}
	return kernel.ObjectInfo{ID: id, Size: info.Size(), Meta: parseMeta(sidePath), ModTime: info.ModTime()}, nil
}

// Delete removes the object and its sidecar. A missing object is a no-op
// (idempotent), mirroring RasterAssetStore.Delete.
func (s *LocalStore) Delete(_ context.Context, ns kernel.ObjectNamespace, id string) error {
	bodyPath, sidePath, err := s.objectPaths(ns, id)
	if err != nil {
		return err
	}
	if err := os.Remove(bodyPath); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	_ = os.Remove(sidePath)
	return nil
}

// Exists reports whether the object body is present.
func (s *LocalStore) Exists(_ context.Context, ns kernel.ObjectNamespace, id string) (bool, error) {
	bodyPath, _, err := s.objectPaths(ns, id)
	if err != nil {
		return false, err
	}
	switch _, err := os.Stat(bodyPath); {
	case err == nil:
		return true, nil
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	default:
		return false, err
	}
}

// List returns every stored id in the namespace, ascending. A not-yet-created
// namespace directory is an empty list (the historical GC/quota IsNotExist
// semantics), while real read errors propagate. Sidecar files map to their
// object ids; names outside the frozen id shape are ignored so foreign or
// leftover files can neither crash callers nor masquerade as objects.
func (s *LocalStore) List(_ context.Context, ns kernel.ObjectNamespace) ([]string, error) {
	dir, err := s.nsDir(ns)
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []string{}, nil
		}
		return nil, err
	}
	seen := map[string]bool{}
	ids := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := strings.TrimSuffix(entry.Name(), ".meta.json")
		if !kernel.ValidObjectID(name) || seen[name] {
			continue
		}
		seen[name] = true
		ids = append(ids, name)
	}
	sort.Strings(ids)
	return ids, nil
}
