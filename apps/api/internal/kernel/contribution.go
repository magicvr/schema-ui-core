package kernel

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// This file implements the R4 C2 structured contribution contract frozen in
// GOAL-005/attachments/r4-c1-freeze-package-draft.md §2. The types are
// framework-agnostic: only standard library (net/http, database/sql) types are
// referenced; no go.uber.org/fx import is permitted in this package.

// ContributionIdentity is the immutable identity every structured contribution
// carries. Key is the canonical semantic id for its kind (see keyFor*), not an
// arbitrary alias; implementations must validate Key equals the corresponding
// field.
type ContributionIdentity struct {
	ModuleID string
	Key      string
}

// RouteContribution registers one HTTP route (freeze package §2.2).
type RouteContribution struct {
	ContributionIdentity
	Method     string
	Pattern    string
	Handler    http.Handler
	Middleware []string
	Public     bool
}

// PageContribution registers one Schema page (freeze package §2.2).
type PageContribution struct {
	ContributionIdentity
	PageID     string
	Resources  []string
	Actions    []string
	DataSource string
	Owner      string
}

// PermissionContribution registers one authorization permission (freeze §2.2).
type PermissionContribution struct {
	ContributionIdentity
	Permission        string
	Resource          string
	Action            string
	PolicyID          string
	SecretSensitivity string
	// SystemDataVersion identifies the versioned persistence contract for this
	// permission's system-data reconcile entry.
	SystemDataVersion int
}

// NavigationContribution registers one navigation node (freeze package §2.2).
type NavigationContribution struct {
	ContributionIdentity
	NodeID     string
	PageID     string
	Parent     string
	Order      int
	Label      string
	Visibility string
	Permission string
	// SystemDataVersion identifies the versioned persistence contract for this
	// navigation node's system-data reconcile entry.
	SystemDataVersion int
}

// FragmentContribution registers one Manifest fragment (freeze package §2.2).
type FragmentContribution struct {
	ContributionIdentity
	FragmentID           string
	ProtocolVersion      string
	ModuleAPIVersion     string
	RequiredCapabilities []string
	JSON                 []byte
}

// MigrationContribution is a compiled-global migration descriptor collected via
// Provider.CompiledPersistence(); it never enters the enablement-gated Registrar
// (freeze package §4).
type MigrationContribution struct {
	ContributionIdentity
	Version           int
	Name              string
	Checksum          string
	Apply             func(*sql.Tx) error
	Tombstone         bool
	ReconcileVersion  int
	ReconcileChecksum string
	Reconcile         func(*sql.Tx) error
}

// ContributionKind enumerates the structured contribution kinds.
type ContributionKind string

const (
	KindHTTP          ContributionKind = "http"
	KindSchema        ContributionKind = "schema"
	KindAuthorization ContributionKind = "authorization"
	KindNavigation    ContributionKind = "navigation"
	KindManifest      ContributionKind = "manifest"
	KindPersistence   ContributionKind = "persistence"
)

// keyForRoute derives the canonical contribution key for an HTTP route:
// uppercase method, single space, pattern (freeze package §2.2).
func keyForRoute(method, pattern string) string {
	return strings.ToUpper(strings.TrimSpace(method)) + " " + strings.TrimSpace(pattern)
}

// RouteKey is the exported form of keyForRoute so providers and tests build a
// contribution key that matches the kernel's canonical rule.
func RouteKey(method, pattern string) string {
	return keyForRoute(method, pattern)
}

// validateIdentity ensures the contribution's canonical Key equals its semantic
// field, and that the identity is complete.
func validateIdentity(moduleID string, kind ContributionKind, key, canonical string) error {
	if strings.TrimSpace(moduleID) == "" {
		return kernelError(CodeModuleInvalid, moduleID, "%s contribution has empty module id", kind)
	}
	key = strings.TrimSpace(key)
	canonical = strings.TrimSpace(canonical)
	if key == "" {
		return kernelError(CodeModuleInvalid, moduleID, "%s contribution has empty key", kind)
	}
	if key != canonical {
		return kernelError(CodeModuleInvalid, moduleID, "%s contribution key %q does not match canonical key %q", kind, key, canonical)
	}
	return nil
}

func validateRoute(moduleID string, r RouteContribution) error {
	if err := validateIdentity(moduleID, KindHTTP, r.Key, keyForRoute(r.Method, r.Pattern)); err != nil {
		return err
	}
	if r.Handler == nil {
		return kernelError(CodeModuleInvalid, moduleID, "route %q has nil handler", r.Key)
	}
	return nil
}

func validatePage(moduleID string, p PageContribution) error {
	if err := validateIdentity(moduleID, KindSchema, p.Key, p.PageID); err != nil {
		return err
	}
	if strings.TrimSpace(p.Owner) == "" || p.Owner != moduleID {
		return kernelError(CodeModuleInvalid, moduleID, "page %q owner must equal the contributing module", p.PageID)
	}
	return nil
}

func validatePermission(moduleID string, p PermissionContribution) error {
	if err := validateIdentity(moduleID, KindAuthorization, p.Key, p.Permission); err != nil {
		return err
	}
	if strings.TrimSpace(p.Resource) == "" || strings.TrimSpace(p.Action) == "" {
		return kernelError(CodeModuleInvalid, moduleID, "permission %q requires resource and action", p.Permission)
	}
	if strings.TrimSpace(p.PolicyID) == "" || strings.TrimSpace(p.PolicyID) != p.PolicyID {
		return kernelError(CodeModuleInvalid, moduleID, "permission %q requires a trimmed policy id", p.Permission)
	}
	if p.SystemDataVersion <= 0 {
		return kernelError(CodeModuleInvalid, moduleID, "permission %q requires a positive system-data version", p.Permission)
	}
	return nil
}

func validateNavigation(moduleID string, n NavigationContribution) error {
	if err := validateIdentity(moduleID, KindNavigation, n.Key, n.NodeID); err != nil {
		return err
	}
	if strings.TrimSpace(n.Label) == "" {
		return kernelError(CodeModuleInvalid, moduleID, "navigation node %q requires a label", n.NodeID)
	}
	if strings.TrimSpace(n.PageID) == "" || strings.TrimSpace(n.PageID) != n.PageID {
		return kernelError(CodeModuleInvalid, moduleID, "navigation node %q requires a trimmed page id", n.NodeID)
	}
	if strings.TrimSpace(n.Visibility) == "" || strings.TrimSpace(n.Visibility) != n.Visibility {
		return kernelError(CodeModuleInvalid, moduleID, "navigation node %q requires a trimmed visibility expression", n.NodeID)
	}
	if n.SystemDataVersion <= 0 {
		return kernelError(CodeModuleInvalid, moduleID, "navigation node %q requires a positive system-data version", n.NodeID)
	}
	return nil
}

func validateFragment(moduleID string, f FragmentContribution) error {
	if err := validateIdentity(moduleID, KindManifest, f.Key, f.FragmentID); err != nil {
		return err
	}
	if strings.TrimSpace(f.ProtocolVersion) == "" {
		return kernelError(CodeModuleInvalid, moduleID, "fragment %q requires protocol version", f.FragmentID)
	}
	// JSON must be valid JSON and deterministically encodable (freeze §2.2).
	// The authored file may be pretty-printed; a canonical re-encode must
	// succeed so the published bytes are deterministic.
	var canonical any
	if err := json.Unmarshal(f.JSON, &canonical); err != nil {
		return kernelError(CodeModuleInvalid, moduleID, "fragment %q JSON is not valid JSON: %v", f.FragmentID, err)
	}
	if _, err := json.Marshal(canonical); err != nil {
		return kernelError(CodeModuleInvalid, moduleID, "fragment %q JSON cannot be canonically encoded: %v", f.FragmentID, err)
	}
	return nil
}

func validateMigration(moduleID string, m MigrationContribution) error {
	if err := validateIdentity(moduleID, KindPersistence, m.Key, m.Name); err != nil {
		return err
	}
	if m.Version <= 0 || strings.TrimSpace(m.Checksum) == "" {
		return kernelError(CodeModuleInvalid, moduleID, "migration %q requires a positive version and a checksum", m.Name)
	}
	if m.Tombstone && m.Apply != nil {
		return kernelError(CodeModuleInvalid, moduleID, "tombstone migration %q must not carry an Apply function", m.Name)
	}
	if !m.Tombstone && m.Apply == nil {
		return kernelError(CodeModuleInvalid, moduleID, "migration %q requires an Apply function unless it is a tombstone", m.Name)
	}
	return nil
}

// contributionDeclared reports whether a (kind, key) pair is declared in the
// module's pre-registration ContributionKeys (freeze package §2.3 step 3:
// "Register 只能写入 descriptor 已声明的 Kind + Key").
func contributionDeclared(module Module, kind ContributionKind, key string) (bool, error) {
	key = strings.TrimSpace(key)
	contains := func(items []string) bool {
		for _, item := range items {
			if strings.TrimSpace(item) == key {
				return true
			}
		}
		return false
	}
	switch kind {
	case KindHTTP:
		return contains(module.Contributions.Routes), nil
	case KindSchema:
		return contains(module.Contributions.Pages), nil
	case KindAuthorization:
		return contains(module.Contributions.Permissions), nil
	case KindNavigation:
		return contains(module.Contributions.Navigation), nil
	case KindManifest:
		return contains(module.Contributions.Fragments), nil
	default:
		return false, fmt.Errorf("kind %s is not registrar-declared", kind)
	}
}
