package kernel

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"testing"
)

// testProvider implements Provider for the C2 contract tests. It declares a
// module with contributions and registers them through the Registrar unless
// mutate is set.
type testProvider struct {
	desc     Module
	mutate   func(Registrar) error
	register func(context.Context, Registrar) error
	migrate  func() ([]MigrationContribution, error)
}

func (p *testProvider) Descriptor() Module { return p.desc }
func (p *testProvider) CompiledPersistence() ([]MigrationContribution, error) {
	if p.migrate == nil {
		return nil, nil
	}
	return p.migrate()
}
func (p *testProvider) Register(ctx context.Context, r Registrar) error {
	if p.register != nil {
		return p.register(ctx, r)
	}
	return p.mutate(r)
}

func minimalModule(desc Module) Module {
	return Module{
		ID:             desc.ID,
		Version:        desc.Version,
		KernelAPIRange: ">=2.0 <3.0",
		DependsOn:      desc.DependsOn,
		Provides:       desc.Provides,
		Requires:       desc.Requires,
		Contributions:  desc.Contributions,
	}
}

func catalogModule(desc Module) *testProvider {
	return &testProvider{desc: minimalModule(desc)}
}

// sampleModule builds a minimal one-party module declaring all six registrar
// kinds plus a compiled migration. overrides merges ID / Version / Contributions
// so tests can exercise conflicts without repeating the full fixture.
func sampleModule(overrides Module) (Module, *testProvider) {
	base := Module{
		ID:             "test.sample",
		Version:        "2.0.0",
		KernelAPIRange: ">=2.0 <3.0",
		Provides:       []Capability{CapabilityHTTP, CapabilitySchema, CapabilityAuthorization, CapabilityNavigation, CapabilityManifest, CapabilityPersistence},
		Contributions: ContributionKeys{
			Routes:           []string{"GET /api/sample"},
			Pages:            []string{"sample"},
			Navigation:       []string{"menu_sample"},
			Permissions:      []string{"sample.read"},
			ConfigNamespaces: []string{"sample.runtime"},
			Fragments:        []string{"sample-fragment"},
		},
	}
	if overrides.ID != "" {
		base.ID = overrides.ID
	}
	if overrides.Version != "" {
		base.Version = overrides.Version
	}
	if len(overrides.Contributions.Routes) > 0 {
		base.Contributions.Routes = overrides.Contributions.Routes
	}
	if len(overrides.Contributions.Pages) > 0 {
		base.Contributions.Pages = overrides.Contributions.Pages
	}
	if len(overrides.Contributions.Navigation) > 0 {
		base.Contributions.Navigation = overrides.Contributions.Navigation
	}
	if len(overrides.Contributions.Permissions) > 0 {
		base.Contributions.Permissions = overrides.Contributions.Permissions
	}
	if len(overrides.Contributions.ConfigNamespaces) > 0 {
		base.Contributions.ConfigNamespaces = overrides.Contributions.ConfigNamespaces
	}
	if len(overrides.Contributions.Fragments) > 0 {
		base.Contributions.Fragments = overrides.Contributions.Fragments
	}
	module := minimalModule(base)
	provider := &testProvider{desc: module}
	provider.mutate = func(r Registrar) error {
		if err := r.HTTP(RouteContribution{
			ContributionIdentity: ContributionIdentity{ModuleID: module.ID, Key: base.Contributions.Routes[0]},
			Method:               methodOfKey(base.Contributions.Routes[0]), Pattern: patternOfKey(base.Contributions.Routes[0]),
			Handler: http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}),
		}); err != nil {
			return err
		}
		if err := r.Schema(PageContribution{
			ContributionIdentity: ContributionIdentity{ModuleID: module.ID, Key: base.Contributions.Pages[0]},
			PageID:               base.Contributions.Pages[0], Owner: module.ID,
			Document: testPageDocument(base.Contributions.Pages[0]),
		}); err != nil {
			return err
		}
		if err := r.Authorization(PermissionContribution{
			ContributionIdentity: ContributionIdentity{ModuleID: module.ID, Key: base.Contributions.Permissions[0]},
			Permission:           base.Contributions.Permissions[0], Resource: "sample", Action: "read",
			PolicyID: "system.admin", SystemDataVersion: 1,
		}); err != nil {
			return err
		}
		if err := r.Navigation(NavigationContribution{
			ContributionIdentity: ContributionIdentity{ModuleID: module.ID, Key: base.Contributions.Navigation[0]},
			NodeID:               base.Contributions.Navigation[0], PageID: "sample", Order: 1, Label: "Sample",
			Visibility: "system.admin", Permission: base.Contributions.Permissions[0], SystemDataVersion: 1,
		}); err != nil {
			return err
		}
		if err := r.Manifest(FragmentContribution{
			ContributionIdentity: ContributionIdentity{ModuleID: module.ID, Key: base.Contributions.Fragments[0]},
			FragmentID:           base.Contributions.Fragments[0], ProtocolVersion: "1.0",
			RequiredCapabilities: []string{"http", "schema"},
			JSON:                 []byte(`{"page":"sample"}`),
		}); err != nil {
			return err
		}
		return r.Configuration(ConfigurationContribution{
			ContributionIdentity: ContributionIdentity{ModuleID: module.ID, Key: base.Contributions.ConfigNamespaces[0]},
			Namespace:            base.Contributions.ConfigNamespaces[0],
			Defaults:             []byte(`{"enabled":true}`),
			Validate:             func(json.RawMessage) error { return nil },
		})
	}
	provider.migrate = func() ([]MigrationContribution, error) {
		return []MigrationContribution{{
			ContributionIdentity: ContributionIdentity{ModuleID: module.ID, Key: "sample_migration"},
			Version:              1, Name: "sample_migration", Checksum: "c1",
			Apply: func(*sql.Tx) error { return nil },
		}}, nil
	}
	return module, provider
}

func testPageDocument(pageID string) []byte {
	raw, _ := json.Marshal(map[string]any{
		"meta": map[string]any{"pageId": pageID},
		"body": map[string]any{"type": "section"},
	})
	return raw
}

func methodOfKey(key string) string {
	for i := 0; i < len(key); i++ {
		if key[i] == ' ' {
			return key[:i]
		}
	}
	return "GET"
}

func patternOfKey(key string) string {
	for i := 0; i < len(key); i++ {
		if key[i] == ' ' {
			return key[i+1:]
		}
	}
	return key
}

func sixCapabilities() []Capability {
	return []Capability{CapabilityHTTP, CapabilitySchema, CapabilityAuthorization, CapabilityNavigation, CapabilityManifest, CapabilityPersistence}
}

func TestRegisterContributionsHappyPath(t *testing.T) {
	module, provider := sampleModule(Module{})
	plan := Plan{Modules: []Module{module}, Capabilities: []Capability{CapabilityHTTP, CapabilitySchema, CapabilityAuthorization, CapabilityNavigation, CapabilityManifest, CapabilityPersistence}}
	set, err := RegisterContributions(context.Background(), plan, []Provider{provider})
	if err != nil {
		t.Fatalf("RegisterContributions: %v", err)
	}
	if len(set.Routes) != 1 || set.Routes[0].Key != "GET /api/sample" {
		t.Fatalf("routes = %+v, want single GET /api/sample", set.Routes)
	}
	if len(set.Pages) != 1 || set.Pages[0].PageID != "sample" {
		t.Fatalf("pages = %+v, want single sample", set.Pages)
	}
	if len(set.Permissions) != 1 || set.Permissions[0].Permission != "sample.read" {
		t.Fatalf("permissions = %+v, want sample.read", set.Permissions)
	}
	if len(set.Navigation) != 1 || set.Navigation[0].NodeID != "menu_sample" {
		t.Fatalf("navigation = %+v, want menu_sample", set.Navigation)
	}
	if len(set.Fragments) != 1 || set.Fragments[0].FragmentID != "sample-fragment" {
		t.Fatalf("fragments = %+v, want sample-fragment", set.Fragments)
	}
	if len(set.Configurations) != 1 || set.Configurations[0].Namespace != "sample.runtime" {
		t.Fatalf("configurations = %+v, want sample.runtime", set.Configurations)
	}
}

func TestSchemaDocumentValidationAndCopy(t *testing.T) {
	module := minimalModule(Module{
		ID: "test.schema", Version: "2.0.0",
		Contributions: ContributionKeys{Pages: []string{"sample"}},
	})
	plan := Plan{Modules: []Module{module}, Capabilities: []Capability{CapabilitySchema}}

	for name, document := range map[string][]byte{
		"invalid JSON":  []byte(`{`),
		"not an object": []byte(`[]`),
		"missing meta":  []byte(`{"body":{}}`),
		"wrong page id": testPageDocument("other"),
	} {
		t.Run(name, func(t *testing.T) {
			provider := &testProvider{desc: module, mutate: func(r Registrar) error {
				return r.Schema(PageContribution{
					ContributionIdentity: ContributionIdentity{ModuleID: module.ID, Key: "sample"},
					PageID:               "sample", Owner: module.ID, Document: document,
				})
			}}
			if _, err := RegisterContributions(context.Background(), plan, []Provider{provider}); err == nil {
				t.Fatal("invalid page document should fail closed")
			}
		})
	}

	document := testPageDocument("sample")
	provider := &testProvider{desc: module, mutate: func(r Registrar) error {
		return r.Schema(PageContribution{
			ContributionIdentity: ContributionIdentity{ModuleID: module.ID, Key: "sample"},
			PageID:               "sample", Owner: module.ID, Document: document,
		})
	}}
	set, err := RegisterContributions(context.Background(), plan, []Provider{provider})
	if err != nil {
		t.Fatal(err)
	}
	document[0] = 'x'
	if !json.Valid(set.Pages[0].Document) {
		t.Fatal("registered document aliases provider-owned bytes")
	}
}

func TestConfigurationValidationCopyConflictAndOrdering(t *testing.T) {
	module := func(id, namespace string) Module {
		return minimalModule(Module{
			ID: id, Version: "2.0.0",
			Contributions: ContributionKeys{ConfigNamespaces: []string{namespace}},
		})
	}
	provider := func(desc Module, namespace string, defaults []byte, validate func(json.RawMessage) error) Provider {
		return &testProvider{desc: desc, mutate: func(r Registrar) error {
			return r.Configuration(ConfigurationContribution{
				ContributionIdentity: ContributionIdentity{ModuleID: desc.ID, Key: namespace},
				Namespace:            namespace, Defaults: defaults, Validate: validate,
			})
		}}
	}

	valid := func(json.RawMessage) error { return nil }
	for _, tc := range []struct {
		name      string
		namespace string
		defaults  []byte
		validator func(json.RawMessage) error
	}{
		{"invalid namespace", "Sample.Runtime", []byte(`{}`), valid},
		{"non-object defaults", "sample.runtime", []byte(`[]`), valid},
		{"missing validator", "sample.runtime", []byte(`{}`), nil},
		{"defaults rejected", "sample.runtime", []byte(`{}`), func(json.RawMessage) error { return errors.New("rejected") }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			desc := module("test.config", tc.namespace)
			if _, err := RegisterContributions(context.Background(), Plan{Modules: []Module{desc}}, []Provider{provider(desc, tc.namespace, tc.defaults, tc.validator)}); err == nil {
				t.Fatal("invalid configuration should fail closed")
			}
		})
	}

	defaults := []byte(`{"enabled":true}`)
	a := module("test.a", "a.runtime")
	z := module("test.z", "z.runtime")
	set, err := RegisterContributions(context.Background(), Plan{Modules: []Module{z, a}}, []Provider{
		provider(z, "z.runtime", []byte(`{}`), valid),
		provider(a, "a.runtime", defaults, valid),
	})
	if err != nil {
		t.Fatal(err)
	}
	defaults[0] = 'x'
	if got := []string{set.Configurations[0].Namespace, set.Configurations[1].Namespace}; !reflect.DeepEqual(got, []string{"a.runtime", "z.runtime"}) {
		t.Fatalf("configuration order = %v", got)
	}
	if !json.Valid(set.Configurations[0].Defaults) {
		t.Fatal("registered configuration aliases provider-owned defaults")
	}

	dupA := module("test.dup-a", "dup.runtime")
	dupB := module("test.dup-b", "dup.runtime")
	_, err = RegisterContributions(context.Background(), Plan{Modules: []Module{dupA, dupB}}, []Provider{
		provider(dupA, "dup.runtime", []byte(`{}`), valid),
		provider(dupB, "dup.runtime", []byte(`{}`), valid),
	})
	var kerr *Error
	if !errors.As(err, &kerr) || kerr.Code != CodeModuleContributionConflict {
		t.Fatalf("duplicate configuration err = %v", err)
	}
}

func TestPolicyReferenceGrammar(t *testing.T) {
	for _, value := range []string{"system.admin", "system.admin-editor", "system.admin-editor-viewer", "a1.b-2"} {
		if !validDottedIdentifier(value) {
			t.Fatalf("valid reference rejected: %q", value)
		}
	}
	for _, value := range []string{"", " system.admin", "System.admin", "system..admin", "system.-admin", "system.admin-", "system.admin--editor", "system.admin_editor", "系统.admin"} {
		if validDottedIdentifier(value) {
			t.Fatalf("invalid reference accepted: %q", value)
		}
	}

	base := PermissionContribution{
		ContributionIdentity: ContributionIdentity{ModuleID: "test.policy", Key: "sample.read"},
		Permission:           "sample.read", Resource: "sample", Action: "read", PolicyID: "system.custom", SystemDataVersion: 1,
	}
	if err := validatePermission("test.policy", base); err != nil {
		t.Fatalf("well-formed owner-specific policy rejected by kernel: %v", err)
	}
	base.PolicyID = "system.admin || true"
	if err := validatePermission("test.policy", base); err == nil {
		t.Fatal("unversioned policy expression should fail kernel grammar")
	}
	if err := validateNavigation("test.policy", NavigationContribution{
		ContributionIdentity: ContributionIdentity{ModuleID: "test.policy", Key: "menu_sample"},
		NodeID:               "menu_sample", PageID: "sample", Label: "Sample", Visibility: "system.admin && system.editor", SystemDataVersion: 1,
	}); err == nil {
		t.Fatal("unversioned visibility expression should fail kernel grammar")
	}
}

func TestRegisterContributionsSkipsPlanModulesWithoutProvider(t *testing.T) {
	// A plan module without a compiled provider stays centrally registered
	// until C3 migration; the contract function does not fail on it (freeze §7).
	module, provider := sampleModule(Module{})
	extra := Module{ID: "core.not-migrated", Version: "2.0.0", KernelAPIRange: ">=2.0 <3.0"}
	plan := Plan{Modules: []Module{extra, module}, Capabilities: sixCapabilities()}
	set, err := RegisterContributions(context.Background(), plan, []Provider{provider})
	if err != nil {
		t.Fatalf("RegisterContributions: %v", err)
	}
	if len(set.Routes) != 1 {
		t.Fatalf("routes = %d, want 1 (only the migrated module)", len(set.Routes))
	}
}

func TestRegisterContributionsSkipsProviderNotInPlan(t *testing.T) {
	moduleA, providerA := sampleModule(Module{ID: "test.a"})
	_, providerB := sampleModule(Module{ID: "test.b"})
	// providerB is compiled but not enabled in the plan → skipped.
	plan := Plan{Modules: []Module{moduleA}, Capabilities: sixCapabilities()}
	set, err := RegisterContributions(context.Background(), plan, []Provider{providerA, providerB})
	if err != nil {
		t.Fatalf("RegisterContributions: %v", err)
	}
	if len(set.Routes) != 1 || set.Routes[0].ModuleID != "test.a" {
		t.Fatalf("routes = %+v, want only test.a route", set.Routes)
	}
}

func TestRegisterContributionsDescriptorMismatch(t *testing.T) {
	module, provider := sampleModule(Module{Version: "2.0.0"})
	// Plan module version differs from provider descriptor.
	planModule := module
	planModule.Version = "3.0.0"
	plan := Plan{Modules: []Module{planModule}}
	_, err := RegisterContributions(context.Background(), plan, []Provider{provider})
	var kerr *Error
	if !errors.As(err, &kerr) || kerr.Code != CodeModuleAPIMismatch {
		t.Fatalf("err = %v, want MODULE_API_MISMATCH", err)
	}
}

func TestRegisterContributionsDescriptorKernelAPIMismatch(t *testing.T) {
	module, provider := sampleModule(Module{})
	// Provider descriptor KernelAPIRange differs from the plan module (full-match).
	planModule := module
	planModule.KernelAPIRange = ">=9.0 <10.0"
	plan := Plan{Modules: []Module{planModule}}
	_, err := RegisterContributions(context.Background(), plan, []Provider{provider})
	var kerr *Error
	if !errors.As(err, &kerr) || kerr.Code != CodeModuleAPIMismatch {
		t.Fatalf("err = %v, want MODULE_API_MISMATCH on kernel API range", err)
	}
}

func TestRegisterContributionsUndeclaredKey(t *testing.T) {
	module, provider := sampleModule(Module{})
	provider.mutate = func(r Registrar) error {
		return r.HTTP(RouteContribution{
			ContributionIdentity: ContributionIdentity{ModuleID: module.ID, Key: "GET /api/undeclared"},
			Method:               "GET", Pattern: "/api/undeclared",
			Handler: http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}),
		})
	}
	plan := Plan{Modules: []Module{module}}
	_, err := RegisterContributions(context.Background(), plan, []Provider{provider})
	if err == nil {
		t.Fatal("undeclared route key should fail")
	}
}

func TestRegisterContributionsKeyFieldMismatch(t *testing.T) {
	module, provider := sampleModule(Module{})
	provider.mutate = func(r Registrar) error {
		return r.Schema(PageContribution{
			ContributionIdentity: ContributionIdentity{ModuleID: module.ID, Key: "sample"},
			PageID:               "other", Owner: module.ID, // Key != PageID
		})
	}
	plan := Plan{Modules: []Module{module}}
	_, err := RegisterContributions(context.Background(), plan, []Provider{provider})
	if err == nil {
		t.Fatal("key/field mismatch should fail")
	}
}

func TestRegisterContributionsCrossModuleConflict(t *testing.T) {
	dupRoutes := ContributionKeys{Routes: []string{"GET /api/dup"}}
	moduleA, providerA := sampleModule(Module{ID: "test.a", Contributions: dupRoutes})
	moduleB, providerB := sampleModule(Module{ID: "test.b", Contributions: dupRoutes})
	plan := Plan{Modules: []Module{moduleA, moduleB}}
	_, err := RegisterContributions(context.Background(), plan, []Provider{providerA, providerB})
	var kerr *Error
	if !errors.As(err, &kerr) || kerr.Code != CodeModuleContributionConflict {
		t.Fatalf("err = %v, want MODULE_CONTRIBUTION_CONFLICT", err)
	}
}

func TestRegisterContributionsNavigationDanglingPermission(t *testing.T) {
	module, provider := sampleModule(Module{})
	provider.mutate = func(r Registrar) error {
		if err := r.HTTP(RouteContribution{
			ContributionIdentity: ContributionIdentity{ModuleID: module.ID, Key: "GET /api/sample"},
			Method:               "GET", Pattern: "/api/sample",
			Handler: http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}),
		}); err != nil {
			return err
		}
		if err := r.Schema(PageContribution{
			ContributionIdentity: ContributionIdentity{ModuleID: module.ID, Key: "sample"},
			PageID:               "sample", Owner: module.ID, Document: testPageDocument("sample"),
		}); err != nil {
			return err
		}
		if err := r.Authorization(PermissionContribution{
			ContributionIdentity: ContributionIdentity{ModuleID: module.ID, Key: "sample.read"},
			Permission:           "sample.read", Resource: "sample", Action: "read", PolicyID: "system.admin", SystemDataVersion: 1,
		}); err != nil {
			return err
		}
		// navigation references a permission that is never registered
		return r.Navigation(NavigationContribution{
			ContributionIdentity: ContributionIdentity{ModuleID: module.ID, Key: "menu_sample"},
			NodeID:               "menu_sample", Order: 1, Label: "Sample", Permission: "ghost.permission",
		})
	}
	plan := Plan{Modules: []Module{module}, Capabilities: []Capability{CapabilityHTTP, CapabilitySchema, CapabilityAuthorization, CapabilityNavigation, CapabilityManifest, CapabilityPersistence}}
	_, err := RegisterContributions(context.Background(), plan, []Provider{provider})
	if err == nil {
		t.Fatal("dangling navigation permission should fail")
	}
}

func TestRegisterContributionsFragmentCapabilityMissing(t *testing.T) {
	module, provider := sampleModule(Module{})
	provider.mutate = func(r Registrar) error {
		// only register the fragment whose required capability the plan lacks
		return r.Manifest(FragmentContribution{
			ContributionIdentity: ContributionIdentity{ModuleID: module.ID, Key: "sample-fragment"},
			FragmentID:           "sample-fragment", ProtocolVersion: "1.0",
			RequiredCapabilities: []string{"uploads"}, // not provided by plan
			JSON:                 []byte(`{}`),
		})
	}
	plan := Plan{Modules: []Module{module}, Capabilities: []Capability{CapabilityManifest}}
	_, err := RegisterContributions(context.Background(), plan, []Provider{provider})
	if err == nil {
		t.Fatal("fragment requiring missing capability should fail")
	}
}

func TestDeterministicOrdering(t *testing.T) {
	module, provider := sampleModule(Module{})
	plan := Plan{Modules: []Module{module}, Capabilities: []Capability{CapabilityHTTP, CapabilitySchema, CapabilityAuthorization, CapabilityNavigation, CapabilityManifest, CapabilityPersistence}}
	first, err := RegisterContributions(context.Background(), plan, []Provider{provider})
	if err != nil {
		t.Fatalf("RegisterContributions: %v", err)
	}
	second, err := RegisterContributions(context.Background(), plan, []Provider{provider})
	if err != nil {
		t.Fatalf("RegisterContributions (2nd): %v", err)
	}
	if first.Routes[0].Key != second.Routes[0].Key || first.Pages[0].PageID != second.Pages[0].PageID {
		t.Fatal("registration order is not deterministic")
	}
}

// TestDualProfileContractMatrix exercises the C2 contract mechanics under both
// compiled profile shapes (freeze package §3: each profile must run the matrix).
// The runtime register/conflict/Start/Ready failure matrix over real modules is
// a C3/C5 gate once business modules are wired as providers.
func TestDualProfileContractMatrix(t *testing.T) {
	for _, name := range []ProfileName{ProfileMVP, ProfileAdmin} {
		resolution, err := ResolveProfile(string(name), nil)
		if err != nil {
			t.Fatalf("resolve %s profile: %v", name, err)
		}
		module, provider := sampleModule(Module{})
		registry, err := NewRegistry(append(BuiltinModules(), module))
		if err != nil {
			t.Fatalf("registry: %v", err)
		}
		plan, err := registry.Resolve(append(resolution.Modules, module.ID))
		if err != nil {
			t.Fatalf("%s plan resolve: %v", name, err)
		}
		if !plan.HasModule(module.ID) {
			t.Fatalf("%s plan missing sample module", name)
		}
		set, err := RegisterContributions(context.Background(), plan, []Provider{provider})
		if err != nil {
			t.Fatalf("%s RegisterContributions: %v", name, err)
		}
		if len(set.Routes) != 1 || len(set.Pages) != 1 || len(set.Permissions) != 1 || len(set.Navigation) != 1 || len(set.Fragments) != 1 || len(set.Configurations) != 1 {
			t.Fatalf("%s set = %+v, want one of each surface kind", name, set)
		}
	}
}

// TestFragmentsDeclarationConflict verifies ContributionKeys.Fragments join the
// Plan-level global uniqueness check (F-IND-C2-002).
func TestFragmentsDeclarationConflict(t *testing.T) {
	dupFragments := ContributionKeys{Fragments: []string{"shared-fragment"}}
	moduleA := Module{ID: "test.a", Version: "2.0.0", KernelAPIRange: ">=2.0 <3.0", Contributions: dupFragments}
	moduleB := Module{ID: "test.b", Version: "2.0.0", KernelAPIRange: ">=2.0 <3.0", Contributions: dupFragments}
	registry, err := NewRegistry([]Module{moduleA, moduleB})
	if err != nil {
		t.Fatalf("registry: %v", err)
	}
	_, err = registry.Resolve([]string{"test.a", "test.b"})
	var kerr *Error
	if !errors.As(err, &kerr) || kerr.Code != CodeModuleContributionConflict {
		t.Fatalf("err = %v, want MODULE_CONTRIBUTION_CONFLICT on fragment declaration", err)
	}
}

// TestNavigationParentOrderIndependent verifies parent references are validated
// after all NodeIDs are collected, so a child appearing before its parent never
// fails (F-IND-C2-007).
func TestNavigationParentOrderIndependent(t *testing.T) {
	navKeys := ContributionKeys{Navigation: []string{"menu_sample", "menu_sample_child"}}
	module, provider := sampleModule(Module{Contributions: navKeys})
	provider.mutate = func(r Registrar) error {
		if err := r.HTTP(RouteContribution{
			ContributionIdentity: ContributionIdentity{ModuleID: module.ID, Key: "GET /api/sample"},
			Method:               "GET", Pattern: "/api/sample",
			Handler: http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}),
		}); err != nil {
			return err
		}
		if err := r.Schema(PageContribution{
			ContributionIdentity: ContributionIdentity{ModuleID: module.ID, Key: "sample"},
			PageID:               "sample", Owner: module.ID, Document: testPageDocument("sample"),
		}); err != nil {
			return err
		}
		if err := r.Authorization(PermissionContribution{
			ContributionIdentity: ContributionIdentity{ModuleID: module.ID, Key: "sample.read"},
			Permission:           "sample.read", Resource: "sample", Action: "read", PolicyID: "system.admin", SystemDataVersion: 1,
		}); err != nil {
			return err
		}
		// child registered before its parent on purpose
		if err := r.Navigation(NavigationContribution{
			ContributionIdentity: ContributionIdentity{ModuleID: module.ID, Key: "menu_sample_child"},
			NodeID:               "menu_sample_child", PageID: "sample", Parent: "menu_sample", Order: 2, Label: "Child",
			Visibility: "system.admin", SystemDataVersion: 1,
		}); err != nil {
			return err
		}
		return r.Navigation(NavigationContribution{
			ContributionIdentity: ContributionIdentity{ModuleID: module.ID, Key: "menu_sample"},
			NodeID:               "menu_sample", PageID: "sample", Order: 1, Label: "Sample",
			Visibility: "system.admin", SystemDataVersion: 1,
		})
	}
	plan := Plan{Modules: []Module{module}, Capabilities: sixCapabilities()}
	set, err := RegisterContributions(context.Background(), plan, []Provider{provider})
	if err != nil {
		t.Fatalf("RegisterContributions with parent-after-child: %v", err)
	}
	if len(set.Navigation) != 2 {
		t.Fatalf("navigation = %d, want 2", len(set.Navigation))
	}
}
