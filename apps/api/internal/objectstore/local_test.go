package objectstore

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
)

const (
	idA = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" // 32 lowercase hex
	idB = "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
)

func newTestStore(t *testing.T) (*LocalStore, string) {
	t.Helper()
	root := t.TempDir()
	return NewLocal(root), root
}

// R1 checkpoint 2a: full round-trip over the frozen port surface.
func TestLocalRoundTrip(t *testing.T) {
	s, _ := newTestStore(t)
	ctx := context.Background()

	if err := s.Put(ctx, kernel.ObjectNamespaceUploads, idA, []byte("hello"), kernel.ObjectMeta{
		Name: "notes.txt", Type: "text/plain", Owner: "user-1",
	}); err != nil {
		t.Fatalf("put: %v", err)
	}
	body, meta, err := s.Get(ctx, kernel.ObjectNamespaceUploads, idA)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if string(body) != "hello" {
		t.Fatalf("body = %q, want hello", body)
	}
	if meta.Name != "notes.txt" || meta.Type != "text/plain" || meta.Owner != "user-1" || meta.Kind != "" {
		t.Fatalf("meta = %+v", meta)
	}

	info, err := s.Stat(ctx, kernel.ObjectNamespaceUploads, idA)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.ID != idA || info.Size != 5 || info.Meta.Owner != "user-1" {
		t.Fatalf("info = %+v", info)
	}

	ok, err := s.Exists(ctx, kernel.ObjectNamespaceUploads, idA)
	if err != nil || !ok {
		t.Fatalf("exists = %v, %v; want true, nil", ok, err)
	}

	// Upsert semantics: Put replaces body and meta.
	if err := s.Put(ctx, kernel.ObjectNamespaceUploads, idA, []byte("replacement!"), kernel.ObjectMeta{Type: "application/json"}); err != nil {
		t.Fatalf("re-put: %v", err)
	}
	body, _, _ = s.Get(ctx, kernel.ObjectNamespaceUploads, idA)
	if string(body) != "replacement!" {
		t.Fatalf("body after overwrite = %q", body)
	}
	info, _ = s.Stat(ctx, kernel.ObjectNamespaceUploads, idA)
	if info.Size != int64(len("replacement!")) {
		t.Fatalf("size after overwrite = %d", info.Size)
	}

	// Delete is idempotent.
	for i := 0; i < 2; i++ {
		if err := s.Delete(ctx, kernel.ObjectNamespaceUploads, idA); err != nil {
			t.Fatalf("delete #%d: %v", i+1, err)
		}
	}
	if _, _, err := s.Get(ctx, kernel.ObjectNamespaceUploads, idA); !errors.Is(err, kernel.ErrObjectNotFound) {
		t.Fatalf("get after delete: err = %v, want ErrObjectNotFound", err)
	}
	if ok, err := s.Exists(ctx, kernel.ObjectNamespaceUploads, idA); err != nil || ok {
		t.Fatalf("exists after delete = %v, %v; want false, nil", ok, err)
	}
}

// Get/Stat misses map to the kernel sentinel (errors.Is across adapters).
func TestLocalNotFoundSentinel(t *testing.T) {
	s, _ := newTestStore(t)
	ctx := context.Background()
	if _, _, err := s.Get(ctx, kernel.ObjectNamespaceAvatars, idB); !errors.Is(err, kernel.ErrObjectNotFound) {
		t.Fatalf("get miss err = %v, want errors.Is ErrObjectNotFound", err)
	}
	if _, err := s.Stat(ctx, kernel.ObjectNamespaceAvatars, idB); !errors.Is(err, kernel.ErrObjectNotFound) {
		t.Fatalf("stat miss err = %v, want errors.Is ErrObjectNotFound", err)
	}
}

// R1 checkpoint 2b: fail-closed validation - unknown namespaces and crafted
// ids are rejected before touching the filesystem (nothing is created).
func TestLocalValidationFailClosed(t *testing.T) {
	s, root := newTestStore(t)
	ctx := context.Background()
	badNS := kernel.ObjectNamespace("../../escape")
	badIDs := []string{"", "../" + idA, strings.ToUpper(idA), idA[:31], idA + "0", "gfffffffffffffffffffffffffffffff"}

	if err := s.Put(ctx, badNS, idA, []byte("x"), kernel.ObjectMeta{}); err == nil {
		t.Fatal("put with unknown namespace must fail")
	}
	for _, bad := range badIDs {
		if err := s.Put(ctx, kernel.ObjectNamespaceUploads, bad, []byte("x"), kernel.ObjectMeta{}); err == nil {
			t.Fatalf("put accepted invalid id %q", bad)
		}
		if _, _, err := s.Get(ctx, kernel.ObjectNamespaceUploads, bad); err == nil {
			t.Fatalf("get accepted invalid id %q", bad)
		}
		if err := s.Delete(ctx, kernel.ObjectNamespaceUploads, bad); err == nil {
			t.Fatalf("delete accepted invalid id %q", bad)
		}
	}
	if _, err := s.List(ctx, badNS); err == nil {
		t.Fatal("list with unknown namespace must fail")
	}
	// Nothing was created on disk.
	if entries, err := os.ReadDir(root); err == nil && len(entries) > 0 {
		t.Fatalf("rejected calls created filesystem entries: %v", entries)
	}
}

// R1 checkpoint 2c: legacy sidecar shapes stay readable - zero migration.
func TestLocalLegacyMetaCompatibility(t *testing.T) {
	s, root := newTestStore(t)
	ctx := context.Background()
	upDir := filepath.Join(root, "uploads")
	avDir := filepath.Join(root, "avatars")
	for _, dir := range []string{upDir, avDir} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	// Historical upload object: name/type/owner sidecar.
	if err := os.WriteFile(filepath.Join(upDir, idA), []byte("csv-bytes"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(upDir, idA+".meta.json"), []byte(`{"name":"data.csv","type":"text/csv","owner":"admin"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	// Historical raster object: type/kind/owner sidecar (owner may be empty).
	if err := os.WriteFile(filepath.Join(avDir, idB), []byte("png-bytes"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(avDir, idB+".meta.json"), []byte(`{"type":"image/png","kind":"avatar","owner":""}`), 0o644); err != nil {
		t.Fatal(err)
	}

	_, meta, err := s.Get(ctx, kernel.ObjectNamespaceUploads, idA)
	if err != nil || meta.Name != "data.csv" || meta.Type != "text/csv" || meta.Owner != "admin" {
		t.Fatalf("legacy upload meta = %+v, err %v", meta, err)
	}
	_, rmeta, err := s.Get(ctx, kernel.ObjectNamespaceAvatars, idB)
	if err != nil || rmeta.Kind != "avatar" || rmeta.Type != "image/png" {
		t.Fatalf("legacy raster meta = %+v, err %v", rmeta, err)
	}

	// Corrupt sidecar is tolerated with empty meta (historical behavior).
	if err := os.WriteFile(filepath.Join(upDir, idB), []byte("orphan"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(upDir, idB+".meta.json"), []byte("{not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	body, cmeta, err := s.Get(ctx, kernel.ObjectNamespaceUploads, idB)
	if err != nil || string(body) != "orphan" || cmeta != (kernel.ObjectMeta{}) {
		t.Fatalf("corrupt sidecar: body=%q meta=%+v err=%v", body, cmeta, err)
	}
}

// Namespaces are isolated; List is ascending and treats a missing namespace
// as empty; foreign files never masquerade as objects.
func TestLocalListAndIsolation(t *testing.T) {
	s, root := newTestStore(t)
	ctx := context.Background()

	if err := s.Put(ctx, kernel.ObjectNamespaceBrandAssets, idB, []byte("b"), kernel.ObjectMeta{}); err != nil {
		t.Fatal(err)
	}
	if err := s.Put(ctx, kernel.ObjectNamespaceBrandAssets, idA, []byte("a"), kernel.ObjectMeta{}); err != nil {
		t.Fatal(err)
	}
	// Same id in another namespace stays independent.
	if err := s.Put(ctx, kernel.ObjectNamespaceUploads, idB, []byte("other"), kernel.ObjectMeta{}); err != nil {
		t.Fatal(err)
	}

	got, err := s.List(ctx, kernel.ObjectNamespaceBrandAssets)
	if err != nil || len(got) != 2 || got[0] != idA || got[1] != idB {
		t.Fatalf("list = %v, err %v; want [idA idB] ascending", got, err)
	}

	// Foreign / leftover files are invisible to List but harmless.
	foreign := filepath.Join(root, "brand-assets", "README.txt")
	if err := os.WriteFile(foreign, []byte("hi"), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err = s.List(ctx, kernel.ObjectNamespaceBrandAssets)
	if err != nil || len(got) != 2 {
		t.Fatalf("list with foreign file = %v, err %v", got, err)
	}

	// Deleting one namespace must not touch the other.
	if err := s.Delete(ctx, kernel.ObjectNamespaceBrandAssets, idB); err != nil {
		t.Fatal(err)
	}
	if ok, _ := s.Exists(ctx, kernel.ObjectNamespaceUploads, idB); !ok {
		t.Fatal("upload namespace object was deleted cross-namespace")
	}

	// Missing namespace = empty list, no error (GC/quota startup semantics).
	got, err = s.List(ctx, kernel.ObjectNamespaceAvatars)
	if err != nil || got == nil || len(got) != 0 {
		t.Fatalf("missing namespace list = %v, err %v; want empty non-nil", got, err)
	}
}