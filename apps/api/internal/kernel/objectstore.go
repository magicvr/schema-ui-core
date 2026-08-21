package kernel

// Kernel object-storage port (VP-014 / workspace-014 GOAL-002 D-001, R1).
//
// The port is the only storage contract for the three first-party object
// families (avatars, brand-assets, uploads). Public types carry neither local
// paths nor os.File: callers address objects by (namespace, id) pairs and the
// active adapter (local disk default, S3-compatible explicit config) resolves
// the physical location. This mirrors the kernel.Store persistence port
// precedent: domain code consumes the port, never a backend handle.
//
// Bucket model (GOAL-001 D-002): one bucket / one local root; the namespace
// is the isolation unit. Local layout <root>/<ns>/<id> therefore matches the
// historical per-family directories byte for byte (zero migration), and the
// S3 key space is <ns>/<id> inside the single configured bucket.
//
// Enumeration (GOAL-001 D-003): List + Stat are port primitives so GC and
// per-user quota scans keep identical semantics on both adapters. GC/quota
// policies themselves stay at the call sites.

import (
	"context"
	"errors"
	"regexp"
)

// ObjectNamespace isolates the three first-party object families. Anything
// outside the frozen set is rejected fail-closed by adapters and helpers.
type ObjectNamespace string

const (
	// ObjectNamespaceAvatars is the account avatar family (server re-encoded
	// rasters; startup GC against users.avatar_url).
	ObjectNamespaceAvatars ObjectNamespace = "avatars"
	// ObjectNamespaceBrandAssets is the branding family (logo/favicon;
	// startup GC against site_settings).
	ObjectNamespaceBrandAssets ObjectNamespace = "brand-assets"
	// ObjectNamespaceUploads is the shared upload family (generic uploads,
	// file-library and data-transfer import source directory).
	ObjectNamespaceUploads ObjectNamespace = "uploads"
)

// ValidObjectNamespace reports whether ns belongs to the frozen namespace
// set. Unknown namespaces are a programming error, not a runtime condition:
// adapters must reject them before touching storage.
func ValidObjectNamespace(ns ObjectNamespace) bool {
	switch ns {
	case ObjectNamespaceAvatars, ObjectNamespaceBrandAssets, ObjectNamespaceUploads:
		return true
	default:
		return false
	}
}

// objectIDPattern is the only id shape the port accepts: 16 random bytes as
// lowercase hex — identical to the historical uploadFileIDPattern every
// current writer uses. Enforcing it at the port keeps crafted ids from
// turning into path escapes (local) or key injection (S3), fail-closed.
var objectIDPattern = regexp.MustCompile("^[0-9a-f]{32}$")

// ValidObjectID reports whether id matches the frozen object-id shape.
func ValidObjectID(id string) bool { return objectIDPattern.MatchString(id) }

// ErrObjectNotFound is the kernel "no such object" sentinel. Adapters wrap
// their backend miss with it so callers can errors.Is across dialects, the
// same way kernel.ErrNoRows works for the SQL port.
var ErrObjectNotFound = errors.New("kernel: object not found")

// ObjectMeta is the dialect-neutral metadata record stored alongside every
// object. It unifies the two historical sidecar shapes (uploads:
// name/type/owner; rasters: type/kind/owner). All fields are optional at the
// port level — required-field policy (e.g. owner-only uploads) belongs to the
// HTTP surface, not to storage.
type ObjectMeta struct {
	// Name is the client-supplied filename (uploads family only).
	Name string
	// Type is the server-detected content type.
	Type string
	// Kind is the raster processing kind (logo/favicon/avatar) when set.
	Kind string
	// Owner is the authenticated uploader id (quota + owner-only reads).
	Owner string
}

// ObjectInfo is Stat's result: identity, byte size and metadata without the
// body, so quota scans stay O(1) in bytes read per object.
type ObjectInfo struct {
	ID   string
	Size int64
	Meta ObjectMeta
}

// ObjectStore is the kernel object-storage port (R1). Implementations must
// validate namespace and id (ValidObjectNamespace / ValidObjectID) and fail
// closed on violations. Semantics frozen by GOAL-002 D-001:
//
//   - Put is an idempotent upsert (current writers always mint fresh random
//     ids, so overwrite never occurs in practice).
//   - Delete is idempotent: a missing object returns nil.
//   - Get/Stat map a missing object to ErrObjectNotFound (errors.Is-able).
//   - List returns object ids in ascending order (deterministic on every
//     adapter) and treats a not-yet-created namespace as empty.
type ObjectStore interface {
	Put(ctx context.Context, ns ObjectNamespace, id string, body []byte, meta ObjectMeta) error
	Get(ctx context.Context, ns ObjectNamespace, id string) ([]byte, ObjectMeta, error)
	Stat(ctx context.Context, ns ObjectNamespace, id string) (ObjectInfo, error)
	Delete(ctx context.Context, ns ObjectNamespace, id string) error
	Exists(ctx context.Context, ns ObjectNamespace, id string) (bool, error)
	List(ctx context.Context, ns ObjectNamespace) ([]string, error)
}
