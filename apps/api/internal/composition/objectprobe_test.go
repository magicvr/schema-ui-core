package composition

import (
	"context"
	"strings"
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/internal/config"
)

// A-002 R-003 (GOAL-003): lock the driver=s3 -> HeadBucket probe wiring and
// the local-default no-op without standing up a full composition.
func TestNewObjectProbe(t *testing.T) {
	t.Run("local default yields no probe", func(t *testing.T) {
		probe, err := newObjectProbe(&config.Config{ObjectsDriver: "local"})
		if err != nil || probe != nil {
			t.Fatalf("local probe non-nil=%t err=%v; want false/nil", probe != nil, err)
		}
	})

	t.Run("empty driver behaves like local", func(t *testing.T) {
		probe, err := newObjectProbe(&config.Config{})
		if err != nil || probe != nil {
			t.Fatalf("zero-value probe non-nil=%t err=%v; want false/nil", probe != nil, err)
		}
	})

	t.Run("s3 driver returns failing probe for unreachable endpoint", func(t *testing.T) {
		cfg := &config.Config{
			ObjectsDriver:            "s3",
			ObjectsS3Endpoint:        "http://127.0.0.1:1",
			ObjectsS3Bucket:          "probe-bucket",
			ObjectsS3AccessKeyID:     "k",
			ObjectsS3SecretAccessKey: "s",
		}
		probe, err := newObjectProbe(cfg)
		if err != nil || probe == nil {
			t.Fatalf("s3 probe non-nil=%t err=%v; want true/nil", probe != nil, err)
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