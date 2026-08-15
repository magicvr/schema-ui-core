// Package datapermission provides the admin.data-permission module surface
// as a kernel.Provider (S-09 · GOAL-016 D-002): scope policy + assignment
// management endpoints, the RowScopeProvider consumed by the resource factory,
// data-permission.read / data-permission.write permission keys, the
// data-permission page and data-permission.* audit events. v1 wires no
// production resource into the enforceable set (D-002 §2 — registration is
// left to future domain modules; PATCH rejects unwired resources).
package datapermission

import (
	"context"
	"errors"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/handler"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	authsessiondata "github.com/magicvr/schema-ui-core/apps/api/internal/modules/authsession/systemdata"
	datapermissionstore "github.com/magicvr/schema-ui-core/apps/api/internal/modules/datapermission/store"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/datapermission/manifest"
	datapermissionschema "github.com/magicvr/schema-ui-core/apps/api/internal/modules/datapermission/schema"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
)

// ModuleID is the stable admin.data-permission module identifier.
const ModuleID = "admin.data-permission"

// Service is the data-permission domain service: enforcement resolution
// (ScopeFor) plus policy/assignment management. It satisfies
// handler.RowScopeProvider and handler.DataPermissionService structurally.
type Service struct {
	repo        *datapermissionstore.Repository
	enforceable map[string]bool
}

// NewService constructs the service. enforceable lists resources whose entity
// consumes filter.Scope (ScopeAware) — wired at composition; v1 registers no
// production resource (GOAL-016 D-002 §2).
func NewService(repo *datapermissionstore.Repository, enforceable []string) *Service {
	set := make(map[string]bool, len(enforceable))
	for _, r := range enforceable {
		set[r] = true
	}
	return &Service{repo: repo, enforceable: set}
}

// ScopeFor implements handler.RowScopeProvider. A resource that is not wired
// as enforceable, has no policy, or is disabled yields nil (no constraint).
func (s *Service) ScopeFor(userID, resource string) (*handler.ScopeConstraint, error) {
	if !s.enforceable[resource] {
		return nil, nil
	}
	p, err := s.repo.GetPolicy(resource)
	if errors.Is(err, datapermissionstore.ErrNotFound) || err == nil && !p.Enabled {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	effective := p.DefaultScope
	if a, err := s.repo.GetAssignment(userID, resource); err == nil {
		effective = a.ScopeType
	} else if !errors.Is(err, datapermissionstore.ErrNotFound) {
		return nil, err
	}
	if effective != datapermissionstore.ScopeSelf {
		return nil, nil
	}
	return &handler.ScopeConstraint{
		Resource: resource, ScopeType: datapermissionstore.ScopeSelf,
		OwnerColumn: p.OwnerColumn, ActorID: userID,
	}, nil
}

// ListPolicies returns all registered policies.
func (s *Service) ListPolicies() ([]datapermissionstore.Policy, error) {
	return s.repo.ListPolicies()
}

// UpsertPolicy registers a resource policy. The enforceability gate lives on
// the policy write (A-005 recommended): unwired resources are rejected.
func (s *Service) UpsertPolicy(resource, ownerColumn, defaultScope string, enabled bool, now time.Time) error {
	if !s.enforceable[resource] {
		return datapermissionstore.ErrNotEnforceable
	}
	return s.repo.UpsertPolicy(resource, ownerColumn, defaultScope, enabled, now)
}

// ListAssignments returns one user's assignments.
func (s *Service) ListAssignments(userID string) ([]datapermissionstore.Assignment, error) {
	return s.repo.ListAssignments(userID)
}

// UpsertAssignments upserts one user's assignment map (resource → scope).
func (s *Service) UpsertAssignments(userID string, scopes map[string]string, now time.Time) error {
	for resource := range scopes {
		if !s.enforceable[resource] {
			return datapermissionstore.ErrNotEnforceable
		}
	}
	return s.repo.UpsertAssignments(userID, scopes, now)
}

// Provider implements kernel.Provider for admin.data-permission.
type Provider struct {
	a          *auth.Authenticator
	service    *Service
	operations operationlog.Recorder
}

// New constructs the data-permission provider.
func New(a *auth.Authenticator, service *Service, operations operationlog.Recorder) *Provider {
	return &Provider{a: a, service: service, operations: operations}
}

func (p *Provider) Descriptor() kernel.Module {
	return kernel.Module{
		ID:             ModuleID,
		Version:        "2.0.0",
		KernelAPIRange: ">=2.0 <3.0",
		DependsOn:      []string{"core.auth-session", "core.navigation-capability", "core.schema-render", "core.operationlog"},
		Requires:       kernel.StandardAdminCapabilities(),
		Contributions: kernel.ContributionKeys{
			Routes: []string{
				"GET /api/data-permission/policies", "PATCH /api/data-permission/policies/{resource}",
				"GET /api/data-permission/scopes", "PATCH /api/data-permission/scopes",
			},
			Pages:       []string{"data-permission"},
			Navigation:  []string{"menu_data_permission"},
			Permissions: []string{"data-permission.read", "data-permission.write"},
			Fragments:   []string{"data-permission"},
		},
	}
}

func (p *Provider) CompiledPersistence() ([]kernel.MigrationContribution, error) {
	return nil, nil // tables are owned by the datapermission/migration provider (0027)
}

func (p *Provider) Register(ctx context.Context, reg kernel.Registrar) error {
	for _, route := range handler.DataPermissionRoutes(p.a, p.service, p.operations, ModuleID) {
		if err := reg.HTTP(route); err != nil {
			return err
		}
	}
	for _, pageID := range []string{"data-permission"} {
		if err := reg.Schema(kernel.PageContribution{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: pageID},
			PageID:               pageID,
			Resources:            []string{"data-permission"},
			Actions:              []string{"list", "update"},
			DataSource:           "/api/data-permission/policies",
			Owner:                ModuleID,
			Document:             datapermissionschema.SchemaDocuments()[pageID],
		}); err != nil {
			return err
		}
	}
	for _, permission := range []kernel.PermissionContribution{
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "data-permission.read"}, Permission: "data-permission.read", Resource: "data-permission", Action: "read", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "data-permission.write"}, Permission: "data-permission.write", Resource: "data-permission", Action: "write", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
	} {
		if err := reg.Authorization(permission); err != nil {
			return err
		}
	}
	if err := reg.Navigation(kernel.NavigationContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "menu_data_permission"},
		NodeID:               "menu_data_permission",
		PageID:               "data-permission",
		Order:                9,
		Label:                "Data permission",
		Visibility:           authsessiondata.PolicyAdmin,
		Permission:           "data-permission.read",
		SystemDataVersion:    authsessiondata.SystemDataVersion,
	}); err != nil {
		return err
	}
	return reg.Manifest(kernel.FragmentContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "data-permission"},
		FragmentID:           "data-permission",
		ProtocolVersion:      "2.7",
		RequiredCapabilities: []string{"manifest", "navigation"},
		JSON:                 manifest.FragmentJSON,
	})
}
