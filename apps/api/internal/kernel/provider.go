package kernel

import (
	"context"
	"sort"
	"strings"
)

// This file implements the R4 C2 Provider/Registrar surface frozen in
// GOAL-005/attachments/r4-c1-freeze-package-draft.md §2.1/§2.3. Composition
// constructs providers (framework-agnostic dependencies only) and calls
// RegisterContributions; Fx stays out of module and kernel packages.

// Provider is a compiled one-party module surface. Descriptor() returns the
// registration-time metadata; CompiledPersistence() returns the module's
// compiled-global migration descriptors; Register contributes HTTP/Schema/
// Authorization/Navigation/Manifest surfaces for an enabled Plan.
type Provider interface {
	Descriptor() Module
	CompiledPersistence() ([]MigrationContribution, error)
	Register(context.Context, Registrar) error
}

// Registrar accepts structured surface contributions. There is deliberately no
// Persistence method: migrations flow through Provider.CompiledPersistence()
// into the compiled-global catalog (freeze package §4).
type Registrar interface {
	HTTP(RouteContribution) error
	Schema(PageContribution) error
	Authorization(PermissionContribution) error
	Navigation(NavigationContribution) error
	Manifest(FragmentContribution) error
}

// ContributionSet is the validated, deterministic surface set produced by
// RegisterContributions. It is never partially published: any validation error
// discards the whole set (freeze package §3 step 6).
type ContributionSet struct {
	Routes      []RouteContribution
	Pages       []PageContribution
	Permissions []PermissionContribution
	Navigation  []NavigationContribution
	Fragments   []FragmentContribution
}

// contributionConflictError is the stable fail-closed error for a duplicate
// contribution identity (freeze package §3: conflicts map to a structured code).
func contributionConflictError(kind ContributionKind, key, first, second string) error {
	return kernelError(CodeModuleContributionConflict, second, "%s contribution key %q conflicts with %s", kind, key, first)
}

// RegisterContributions builds the validated surface set from the compiled
// providers that are enabled in the Plan, following freeze package §2.3:
//
//  1. only enabled providers are registered (a plan module without a provider
//     remains centrally registered until C3 migration — freeze §7 step 1/2);
//  2. each enabled provider's Descriptor must exactly match the Plan module
//     (id and version);
//  3. Register may only write declared (Kind, Key) contributions with a matching
//     identity field; undeclared / duplicate / field-mismatch fails immediately;
//  4. before finalize, global conflict, reference-integrity, capability and
//     deterministic-ordering checks run; any failure discards the whole set.
//
// Providers that are compiled but not enabled are skipped for surface
// registration (their persistence still enters the compiled-global catalog via
// CollectPersistence).
func RegisterContributions(ctx context.Context, plan Plan, providers []Provider) (ContributionSet, error) {
	planModules := make(map[string]Module, len(plan.Modules))
	for _, module := range plan.Modules {
		planModules[module.ID] = module
	}

	// Register in deterministic order by descriptor id.
	sorted := append([]Provider(nil), providers...)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Descriptor().ID < sorted[j].Descriptor().ID
	})

	set := ContributionSet{}
	for _, provider := range sorted {
		desc := provider.Descriptor()
		if strings.TrimSpace(desc.ID) == "" {
			return ContributionSet{}, kernelError(CodeModuleInvalid, desc.ID, "provider descriptor has empty module id")
		}
		module, enabled := planModules[desc.ID]
		if !enabled {
			continue // compiled but not enabled: no surface registration
		}
		if err := descriptorsMatch(desc, module); err != nil {
			return ContributionSet{}, err
		}
		reg := &validatingRegistrar{module: module, set: &set, seen: map[ContributionKind]map[string]string{}}
		if err := provider.Register(ctx, reg); err != nil {
			return ContributionSet{}, err
		}
	}

	if err := set.finalize(plan); err != nil {
		return ContributionSet{}, err
	}
	return set, nil
}

// descriptorsMatch implements freeze package §2.3 step 2: the provider
// Descriptor must exactly match the Plan module. All normative fields are
// compared (slice fields as order-insensitive sets).
func descriptorsMatch(desc, plan Module) error {
	mismatch := func(what string) error {
		return kernelError(CodeModuleAPIMismatch, plan.ID, "provider descriptor %s differs from plan module: %s", desc.ID, what)
	}
	if desc.ID != plan.ID {
		return mismatch("id")
	}
	if desc.Version != plan.Version {
		return mismatch("version")
	}
	if desc.KernelAPIRange != plan.KernelAPIRange {
		return mismatch("kernel API range")
	}
	if !stringSetEqual(desc.DependsOn, plan.DependsOn) {
		return mismatch("depends-on set")
	}
	if !capabilitySetEqual(desc.Provides, plan.Provides) {
		return mismatch("provides set")
	}
	if !capabilitySetEqual(desc.Requires, plan.Requires) {
		return mismatch("requires set")
	}
	if !contributionsEqual(desc.Contributions, plan.Contributions) {
		return mismatch("contribution declaration keys")
	}
	return nil
}

// stringSetEqual compares two string slices as sets (order-insensitive,
// duplicates ignored). It is used for DependsOn where the normative declaration
// is a set, not an ordered list.
func stringSetEqual(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	seen := make(map[string]int, len(left))
	for _, item := range left {
		seen[item]++
	}
	for _, item := range right {
		seen[item]--
		if seen[item] < 0 {
			return false
		}
	}
	return true
}

// capabilitySetEqual compares two Capability slices as sets.
func capabilitySetEqual(left, right []Capability) bool {
	if len(left) != len(right) {
		return false
	}
	seen := make(map[Capability]int, len(left))
	for _, item := range left {
		seen[item]++
	}
	for _, item := range right {
		seen[item]--
		if seen[item] < 0 {
			return false
		}
	}
	return true
}

func contributionsEqual(left, right ContributionKeys) bool {
	return stringSetEqual(left.Routes, right.Routes) &&
		stringSetEqual(left.Pages, right.Pages) &&
		stringSetEqual(left.Navigation, right.Navigation) &&
		stringSetEqual(left.Permissions, right.Permissions) &&
		stringSetEqual(left.ConfigNamespaces, right.ConfigNamespaces) &&
		stringSetEqual(left.Fragments, right.Fragments)
}

// validatingRegistrar collects one module's contributions, failing closed on
// undeclared keys, duplicate keys and identity-field mismatches.
type validatingRegistrar struct {
	module Module
	set    *ContributionSet
	seen   map[ContributionKind]map[string]string
}

func (r *validatingRegistrar) declare(kind ContributionKind, key string) error {
	if r.seen[kind] == nil {
		r.seen[kind] = map[string]string{}
	}
	if previous, exists := r.seen[kind][key]; exists {
		return contributionConflictError(kind, key, previous, r.module.ID)
	}
	declared, err := contributionDeclared(r.module, kind, key)
	if err != nil {
		return kernelError(CodeModuleInvalid, r.module.ID, "%v", err)
	}
	if !declared {
		return kernelError(CodeModuleInvalid, r.module.ID, "%s contribution key %q is not declared in the module contributions", kind, key)
	}
	r.seen[kind][key] = r.module.ID
	return nil
}

func (r *validatingRegistrar) HTTP(c RouteContribution) error {
	if c.ModuleID != r.module.ID {
		return kernelError(CodeModuleInvalid, r.module.ID, "route contribution module id %q does not match provider", c.ModuleID)
	}
	if err := validateRoute(r.module.ID, c); err != nil {
		return err
	}
	if err := r.declare(KindHTTP, c.Key); err != nil {
		return err
	}
	r.set.Routes = append(r.set.Routes, c)
	return nil
}

func (r *validatingRegistrar) Schema(c PageContribution) error {
	if c.ModuleID != r.module.ID {
		return kernelError(CodeModuleInvalid, r.module.ID, "page contribution module id %q does not match provider", c.ModuleID)
	}
	if err := validatePage(r.module.ID, c); err != nil {
		return err
	}
	if err := r.declare(KindSchema, c.Key); err != nil {
		return err
	}
	r.set.Pages = append(r.set.Pages, c)
	return nil
}

func (r *validatingRegistrar) Authorization(c PermissionContribution) error {
	if c.ModuleID != r.module.ID {
		return kernelError(CodeModuleInvalid, r.module.ID, "permission contribution module id %q does not match provider", c.ModuleID)
	}
	if err := validatePermission(r.module.ID, c); err != nil {
		return err
	}
	if err := r.declare(KindAuthorization, c.Key); err != nil {
		return err
	}
	r.set.Permissions = append(r.set.Permissions, c)
	return nil
}

func (r *validatingRegistrar) Navigation(c NavigationContribution) error {
	if c.ModuleID != r.module.ID {
		return kernelError(CodeModuleInvalid, r.module.ID, "navigation contribution module id %q does not match provider", c.ModuleID)
	}
	if err := validateNavigation(r.module.ID, c); err != nil {
		return err
	}
	if err := r.declare(KindNavigation, c.Key); err != nil {
		return err
	}
	r.set.Navigation = append(r.set.Navigation, c)
	return nil
}

func (r *validatingRegistrar) Manifest(c FragmentContribution) error {
	if c.ModuleID != r.module.ID {
		return kernelError(CodeModuleInvalid, r.module.ID, "fragment contribution module id %q does not match provider", c.ModuleID)
	}
	if err := validateFragment(r.module.ID, c); err != nil {
		return err
	}
	if err := r.declare(KindManifest, c.Key); err != nil {
		return err
	}
	r.set.Fragments = append(r.set.Fragments, c)
	return nil
}

// finalize runs global conflict, reference-integrity, capability and
// deterministic-ordering checks. Any failure discards the whole set.
func (s *ContributionSet) finalize(plan Plan) error {
	seen := map[string]string{}
	checkUnique := func(kind ContributionKind, key, moduleID string) error {
		identity := string(kind) + "\x00" + key
		if previous, exists := seen[identity]; exists {
			return contributionConflictError(kind, key, previous, moduleID)
		}
		seen[identity] = moduleID
		return nil
	}

	permissions := map[string]string{}
	for _, p := range s.Permissions {
		if err := checkUnique(KindAuthorization, p.Key, p.ModuleID); err != nil {
			return err
		}
		permissions[p.Permission] = p.ModuleID
	}
	for _, r := range s.Routes {
		if err := checkUnique(KindHTTP, r.Key, r.ModuleID); err != nil {
			return err
		}
	}
	pages := map[string]string{}
	for _, p := range s.Pages {
		if err := checkUnique(KindSchema, p.Key, p.ModuleID); err != nil {
			return err
		}
		pages[p.PageID] = p.ModuleID
	}
	// Two-pass navigation validation: collect every NodeID first, then check
	// parent references, so registration order never affects the result.
	nodes := map[string]string{}
	for _, n := range s.Navigation {
		if err := checkUnique(KindNavigation, n.Key, n.ModuleID); err != nil {
			return err
		}
		nodes[n.NodeID] = n.ModuleID
	}
	for _, n := range s.Navigation {
		if owner, ok := pages[n.PageID]; !ok || owner != n.ModuleID {
			return kernelError(CodeModuleInvalid, n.ModuleID, "navigation node %q references page %q outside its module", n.NodeID, n.PageID)
		}
		if n.Permission != "" {
			if _, ok := permissions[n.Permission]; !ok {
				return kernelError(CodeModuleInvalid, n.ModuleID, "navigation node %q references undeclared permission %q", n.NodeID, n.Permission)
			}
		}
		if n.Parent != "" {
			if _, ok := nodes[n.Parent]; !ok {
				return kernelError(CodeModuleInvalid, n.ModuleID, "navigation node %q references undeclared parent %q", n.NodeID, n.Parent)
			}
		}
	}
	capabilities := map[Capability]bool{}
	for _, c := range plan.Capabilities {
		capabilities[c] = true
	}
	for _, f := range s.Fragments {
		if err := checkUnique(KindManifest, f.Key, f.ModuleID); err != nil {
			return err
		}
		for _, required := range f.RequiredCapabilities {
			if !capabilities[Capability(required)] {
				return kernelError(CodeModuleInvalid, f.ModuleID, "fragment %q requires capability %q that the plan does not provide", f.FragmentID, required)
			}
		}
	}

	sortRoutes(s.Routes)
	sortPages(s.Pages)
	sortPermissions(s.Permissions)
	sortNavigation(s.Navigation)
	sortFragments(s.Fragments)
	return nil
}

func sortRoutes(routes []RouteContribution) {
	sort.Slice(routes, func(i, j int) bool { return routes[i].Key < routes[j].Key })
}

func sortPages(pages []PageContribution) {
	sort.Slice(pages, func(i, j int) bool { return pages[i].PageID < pages[j].PageID })
}

func sortPermissions(perms []PermissionContribution) {
	sort.Slice(perms, func(i, j int) bool { return perms[i].Permission < perms[j].Permission })
}

func sortNavigation(nodes []NavigationContribution) {
	sort.Slice(nodes, func(i, j int) bool {
		if nodes[i].Parent != nodes[j].Parent {
			return nodes[i].Parent < nodes[j].Parent
		}
		if nodes[i].Order != nodes[j].Order {
			return nodes[i].Order < nodes[j].Order
		}
		return nodes[i].NodeID < nodes[j].NodeID
	})
}

func sortFragments(fragments []FragmentContribution) {
	sort.Slice(fragments, func(i, j int) bool { return fragments[i].FragmentID < fragments[j].FragmentID })
}
