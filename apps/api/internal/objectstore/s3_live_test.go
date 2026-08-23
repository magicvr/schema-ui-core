package objectstore

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
)

// Live round-trip against a real S3-compatible endpoint (MinIO suffices).
// Skipped cleanly unless every S3_TEST_* variable is set, so a plain
// go test ./... stays offline - mirroring the pgtest precedent.
func TestS3LiveRoundTrip(t *testing.T) {
	endpoint := os.Getenv("S3_TEST_ENDPOINT")
	bucket := os.Getenv("S3_TEST_BUCKET")
	key := os.Getenv("S3_TEST_ACCESS_KEY")
	secret := os.Getenv("S3_TEST_SECRET")
	if endpoint == "" || bucket == "" || key == "" || secret == "" {
		t.Skip("S3_TEST_ENDPOINT/S3_TEST_BUCKET/S3_TEST_ACCESS_KEY/S3_TEST_SECRET not set; skipping live object-store test")
	}
	s, err := NewS3(endpoint, os.Getenv("S3_TEST_REGION"), bucket, key, secret, true)
	if err != nil {
		t.Fatalf("NewS3: %v", err)
	}
	ctx := context.Background()
	if err := s.Ping(ctx); err != nil {
		t.Fatalf("ping (bucket must exist): %v", err)
	}

	id := "0123456789abcdef0123456789abcdef"
	defer func() { _ = s.Delete(ctx, kernel.ObjectNamespaceUploads, id) }()
	if err := s.Put(ctx, kernel.ObjectNamespaceUploads, id, []byte("live-bytes"), kernel.ObjectMeta{Type: "text/plain", Owner: "it"}); err != nil {
		t.Fatalf("put: %v", err)
	}
	body, meta, err := s.Get(ctx, kernel.ObjectNamespaceUploads, id)
	if err != nil || string(body) != "live-bytes" || meta.Owner != "it" {
		t.Fatalf("get = %q/%+v/%v", body, meta, err)
	}
	info, err := s.Stat(ctx, kernel.ObjectNamespaceUploads, id)
	if err != nil || info.Size != int64(len("live-bytes")) {
		t.Fatalf("stat = %+v/%v", info, err)
	}
	ids, err := s.List(ctx, kernel.ObjectNamespaceUploads)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	found := false
	for _, v := range ids {
		if v == id {
			found = true
		}
	}
	if !found {
		t.Fatalf("list missing %s: %v", id, ids)
	}
	if err := s.Delete(ctx, kernel.ObjectNamespaceUploads, id); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if ok, err := s.Exists(ctx, kernel.ObjectNamespaceUploads, id); ok || err != nil {
		t.Fatalf("exists after delete = %v/%v", ok, err)
	}
	if _, _, err := s.Get(ctx, kernel.ObjectNamespaceUploads, id); !errors.Is(err, kernel.ErrObjectNotFound) {
		t.Fatalf("get after delete err = %v", err)
	}
}
