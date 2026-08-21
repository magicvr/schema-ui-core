package composition

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/internal/config"
	"github.com/magicvr/schema-ui-core/apps/api/internal/objectstore"
)

// A-002 R-003 (GOAL-003): lock the driver=s3 -> HeadBucket probe wiring and
// the local-default no-op without standing up a full composition.
func TestNewObjectStoreWiring(t *testing.T) {
	t.Run("local default yields store without probe (root = db.path dir)", func(t *testing.T) {
		store, probe, err := newObjectStore(&config.Config{ObjectsDriver: "local", DBPath: "./data/schema-ui.db"})
		if err != nil || store == nil || probe != nil {
			t.Fatalf("local wiring store-nil=%t probe-non-nil=%t err=%v", store == nil, probe != nil, err)
		}
		local, ok := store.(*objectstore.LocalStore)
		if !ok {
			t.Fatalf("local driver must produce *objectstore.LocalStore, got %T", store)
		}
		if want := filepath.Dir("./data/schema-ui.db"); local.Root() != want {
			t.Fatalf("derived root = %q, want %q (filepath.Dir(db.path))", local.Root(), want)
		}
	})

	t.Run("ObjectsLocalRoot overrides the derived root", func(t *testing.T) {
		store, _, err := newObjectStore(&config.Config{ObjectsDriver: "local", ObjectsLocalRoot: "/tmp/override-root"})
		if err != nil {
			t.Fatalf("newObjectStore: %v", err)
		}
		local, ok := store.(*objectstore.LocalStore)
		if !ok || local.Root() != "/tmp/override-root" {
			t.Fatalf("override root not honored: %T %q", store, local.Root())
		}
	})

	t.Run("empty driver behaves like local", func(t *testing.T) {
		store, probe, err := newObjectStore(&config.Config{})
		if err != nil || store == nil || probe != nil {
			t.Fatalf("zero-value wiring store-nil=%t probe-non-nil=%t err=%v", store == nil, probe != nil, err)
		}
	})

	t.Run("unknown driver fails closed (GOAL-005 A-002 R-002)", func(t *testing.T) {
		for _, driver := range []string{"gcs", "S3"} { // unknown + case-sensitive guard
			store, probe, err := newObjectStore(&config.Config{ObjectsDriver: driver})
			if err == nil || store != nil || probe != nil {
				t.Fatalf("driver %q: store-nil=%t probe-non-nil=%t err=%v; want nil/nil/error", driver, store == nil, probe != nil, err)
			}
		}
	})

	t.Run("s3 driver returns store and failing probe for unreachable endpoint", func(t *testing.T) {
		cfg := &config.Config{
			ObjectsDriver:            "s3",
			ObjectsS3Endpoint:        "http://127.0.0.1:1",
			ObjectsS3Bucket:          "probe-bucket",
			ObjectsS3AccessKeyID:     "k",
			ObjectsS3SecretAccessKey: "s",
		}
		store, probe, err := newObjectStore(cfg)
		if err != nil || probe == nil {
			t.Fatalf("s3 wiring store-nil=%t probe-non-nil=%t err=%v; want store/probe/nil-err", store == nil, probe != nil, err)
		}
		if _, ok := store.(*objectstore.S3Store); !ok {
			t.Fatalf("s3 driver must produce *objectstore.S3Store, got %T", store)
		}
		// Nothing listens on port 1: HeadBucket must fail fast, proving the
		// probe is really wired to the configured endpoint.
		if err := probe(context.Background()); err == nil {
			t.Fatal("probe against unreachable endpoint must fail")
		} else if !strings.Contains(err.Error(), "head bucket") {
			t.Fatalf("probe error should name the operation: %v", err)
		}
	})
}