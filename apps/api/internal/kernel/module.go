package kernel

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

const KernelAPIVersion = "2.0.0"

// Capability is a stable, framework-agnostic capability identifier.
type Capability string

const (
	CapabilityHTTP          Capability = "http"
	CapabilitySchema        Capability = "schema"
	CapabilityAuthorization Capability = "authorization"
	CapabilityNavigation    Capability = "navigation"
	CapabilityManifest      Capability = "manifest"
	CapabilityPersistence   Capability = "persistence"
	CapabilityExpressions   Capability = "expressions"
	CapabilityValidation    Capability = "validation"
	CapabilityOperationLog  Capability = "operationlog"
)

// ContributionKeys are checked globally after Profile expansion. Modules do
// not mutate a central application registry; they declare ownership here.
type ContributionKeys struct {
	Routes           []string
	Pages            []string
	Navigation       []string
	Permissions      []string
	ConfigNamespaces []string
}

// Hooks are deliberately independent of Fx. The composition root may adapt
// them to a framework-specific lifecycle, while module APIs remain portable.
type Hooks struct {
	Start func(context.Context) error
	Ready func(context.Context) error
	Stop  func(context.Context) error
}

// Module is the public semantic contract used by the kernel registry.
type Module struct {
	ID             string
	Version        string
	KernelAPIRange string
	DependsOn      []string
	Provides       []Capability
	Requires       []Capability
	Contributions  ContributionKeys
	Hooks          Hooks
}

type ErrorCode string

const (
	CodeModuleInvalid              ErrorCode = "MODULE_INVALID"
	CodeModuleDuplicate            ErrorCode = "MODULE_DUPLICATE"
	CodeModuleUnknown              ErrorCode = "MODULE_UNKNOWN"
	CodeModuleDuplicateEnabled     ErrorCode = "MODULE_DUPLICATE_ENABLED"
	CodeModuleDependencyUnknown    ErrorCode = "MODULE_DEPENDENCY_UNKNOWN"
	CodeModuleDependencyMissing    ErrorCode = "MODULE_DEPENDENCY_MISSING"
	CodeModuleDependencyCycle      ErrorCode = "MODULE_DEPENDENCY_CYCLE"
	CodeModuleContributionConflict ErrorCode = "MODULE_CONTRIBUTION_CONFLICT"
	CodeModuleCapabilityMissing    ErrorCode = "MODULE_CAPABILITY_MISSING"
	CodeModuleAPIRangeInvalid      ErrorCode = "MODULE_API_RANGE_INVALID"
	CodeModuleAPIMismatch          ErrorCode = "MODULE_API_MISMATCH"
	CodeProfileUnknown             ErrorCode = "PROFILE_UNKNOWN"
	CodeProfileModulesRequired     ErrorCode = "PROFILE_MODULES_REQUIRED"
	CodeLifecycleStartFailed       ErrorCode = "LIFECYCLE_START_FAILED"
	CodeLifecycleReadyFailed       ErrorCode = "LIFECYCLE_READY_FAILED"
	CodeLifecycleStopFailed        ErrorCode = "LIFECYCLE_STOP_FAILED"
)

// Error is intentionally structured so startup diagnostics can retain a
// stable code and module identifier without exposing an Fx error type.
type Error struct {
	Code     ErrorCode
	ModuleID string
	Detail   string
}

func (e *Error) Error() string {
	if e.ModuleID == "" {
		return fmt.Sprintf("%s: %s", e.Code, e.Detail)
	}
	return fmt.Sprintf("%s [%s]: %s", e.Code, e.ModuleID, e.Detail)
}

func kernelError(code ErrorCode, moduleID, format string, args ...any) error {
	return &Error{Code: code, ModuleID: moduleID, Detail: fmt.Sprintf(format, args...)}
}

type Registry struct {
	modules map[string]Module
}

func NewRegistry(modules []Module) (*Registry, error) {
	registry := &Registry{modules: make(map[string]Module, len(modules))}
	for _, module := range modules {
		module.ID = strings.TrimSpace(module.ID)
		module.Version = strings.TrimSpace(module.Version)
		module.KernelAPIRange = strings.TrimSpace(module.KernelAPIRange)
		if module.ID == "" || module.Version == "" || module.KernelAPIRange == "" {
			return nil, kernelError(CodeModuleInvalid, module.ID, "id, version and kernel API range are required")
		}
		compatible, err := kernelAPIRangeAccepts(module.KernelAPIRange, KernelAPIVersion)
		if err != nil {
			return nil, kernelError(CodeModuleAPIRangeInvalid, module.ID, "%v", err)
		}
		if !compatible {
			return nil, kernelError(CodeModuleAPIMismatch, module.ID, "declared range %q does not accept kernel API %s", module.KernelAPIRange, KernelAPIVersion)
		}
		if _, exists := registry.modules[module.ID]; exists {
			return nil, kernelError(CodeModuleDuplicate, module.ID, "module id is registered more than once")
		}
		registry.modules[module.ID] = cloneModule(module)
	}
	return registry, nil
}

type semanticVersion struct {
	major int
	minor int
	patch int
}

type rangeConstraint struct {
	operator string
	version  semanticVersion
}

func kernelAPIRangeAccepts(rawRange, rawVersion string) (bool, error) {
	constraints, err := parseRange(rawRange)
	if err != nil {
		return false, err
	}
	current, err := parseSemanticVersion(rawVersion)
	if err != nil {
		return false, err
	}
	for _, constraint := range constraints {
		if !compareVersions(current, constraint.version, constraint.operator) {
			return false, nil
		}
	}
	return true, nil
}

func parseRange(raw string) ([]rangeConstraint, error) {
	fields := strings.Fields(strings.TrimSpace(raw))
	if len(fields) == 0 {
		return nil, fmt.Errorf("kernel API range is empty")
	}
	constraints := make([]rangeConstraint, 0, len(fields))
	for _, field := range fields {
		op := "="
		value := field
		for _, candidate := range []string{">=", "<=", "==", ">", "<", "="} {
			if strings.HasPrefix(field, candidate) {
				op = candidate
				value = strings.TrimPrefix(field, candidate)
				break
			}
		}
		if op == "==" {
			op = "="
		}
		version, err := parseSemanticVersion(value)
		if err != nil {
			return nil, fmt.Errorf("invalid kernel API range token %q: %w", field, err)
		}
		constraints = append(constraints, rangeConstraint{operator: op, version: version})
	}
	return constraints, nil
}

func parseSemanticVersion(raw string) (semanticVersion, error) {
	parts := strings.Split(strings.TrimSpace(raw), ".")
	if len(parts) < 1 || len(parts) > 3 {
		return semanticVersion{}, fmt.Errorf("version %q must have one to three numeric components", raw)
	}
	values := [3]int{}
	for i, part := range parts {
		if part == "" {
			return semanticVersion{}, fmt.Errorf("version %q contains an empty component", raw)
		}
		value, err := strconv.Atoi(part)
		if err != nil || value < 0 {
			return semanticVersion{}, fmt.Errorf("version %q contains only non-negative integer components", raw)
		}
		values[i] = value
	}
	return semanticVersion{major: values[0], minor: values[1], patch: values[2]}, nil
}

func compareVersions(left, right semanticVersion, operator string) bool {
	comparison := 0
	for _, pair := range [][2]int{{left.major, right.major}, {left.minor, right.minor}, {left.patch, right.patch}} {
		if pair[0] < pair[1] {
			comparison = -1
			break
		}
		if pair[0] > pair[1] {
			comparison = 1
			break
		}
	}
	switch operator {
	case ">":
		return comparison > 0
	case ">=":
		return comparison >= 0
	case "<":
		return comparison < 0
	case "<=":
		return comparison <= 0
	default:
		return comparison == 0
	}
}

func (r *Registry) Module(id string) (Module, bool) {
	module, ok := r.modules[id]
	return cloneModule(module), ok
}

type Plan struct {
	Modules      []Module
	Capabilities []Capability
}

// HasModule reports whether a module was selected in this resolved plan.
// Runtime registration must consume this plan rather than re-resolving a
// profile, so all protocol surfaces share one enablement decision.
func (p Plan) HasModule(id string) bool {
	for _, module := range p.Modules {
		if module.ID == id {
			return true
		}
	}
	return false
}

func (p Plan) IDs() []string {
	ids := make([]string, 0, len(p.Modules))
	for _, module := range p.Modules {
		ids = append(ids, module.ID)
	}
	return ids
}

func (r *Registry) Resolve(enabled []string) (Plan, error) {
	if len(enabled) == 0 {
		return Plan{}, kernelError(CodeModuleInvalid, "", "enabled module set must not be empty")
	}
	selected := make(map[string]Module, len(enabled))
	for _, rawID := range enabled {
		id := strings.TrimSpace(rawID)
		if id == "" {
			return Plan{}, kernelError(CodeModuleInvalid, "", "enabled module id must not be empty")
		}
		if _, exists := selected[id]; exists {
			return Plan{}, kernelError(CodeModuleDuplicateEnabled, id, "module is enabled more than once")
		}
		module, exists := r.modules[id]
		if !exists {
			return Plan{}, kernelError(CodeModuleUnknown, id, "module is not in the compiled candidate set")
		}
		selected[id] = cloneModule(module)
	}

	ids := make([]string, 0, len(selected))
	for id := range selected {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		for _, dep := range selected[id].DependsOn {
			dep = strings.TrimSpace(dep)
			if _, compiled := r.modules[dep]; !compiled {
				return Plan{}, kernelError(CodeModuleDependencyUnknown, id, "dependency %q is not compiled", dep)
			}
			if _, enabledDep := selected[dep]; !enabledDep {
				return Plan{}, kernelError(CodeModuleDependencyMissing, id, "dependency %q is not enabled", dep)
			}
		}
	}

	state := make(map[string]uint8, len(selected))
	ordered := make([]Module, 0, len(selected))
	var visit func(string) error
	visit = func(id string) error {
		switch state[id] {
		case 1:
			return kernelError(CodeModuleDependencyCycle, id, "dependency graph contains a cycle")
		case 2:
			return nil
		}
		state[id] = 1
		deps := append([]string(nil), selected[id].DependsOn...)
		sort.Strings(deps)
		for _, dep := range deps {
			if err := visit(dep); err != nil {
				return err
			}
		}
		state[id] = 2
		ordered = append(ordered, cloneModule(selected[id]))
		return nil
	}
	for _, id := range ids {
		if err := visit(id); err != nil {
			return Plan{}, err
		}
	}

	provided := make(map[Capability]string)
	for _, module := range ordered {
		for _, capability := range module.Provides {
			if _, exists := provided[capability]; !exists {
				provided[capability] = module.ID
			}
		}
	}
	for _, module := range ordered {
		for _, capability := range module.Requires {
			if _, exists := provided[capability]; !exists {
				return Plan{}, kernelError(CodeModuleCapabilityMissing, module.ID, "required capability %q is not provided", capability)
			}
		}
	}
	if err := validateContributions(ordered); err != nil {
		return Plan{}, err
	}

	capabilities := make([]Capability, 0, len(provided))
	for capability := range provided {
		capabilities = append(capabilities, capability)
	}
	sort.Slice(capabilities, func(i, j int) bool { return capabilities[i] < capabilities[j] })
	return Plan{Modules: ordered, Capabilities: capabilities}, nil
}

func validateContributions(modules []Module) error {
	type ownership struct{ kind, key, module string }
	seen := map[string]ownership{}
	add := func(module, kind, key string) error {
		key = strings.TrimSpace(key)
		if key == "" {
			return kernelError(CodeModuleInvalid, module, "%s contribution key must not be empty", kind)
		}
		identity := kind + "\x00" + key
		if previous, exists := seen[identity]; exists {
			return kernelError(CodeModuleContributionConflict, module, "%s %q conflicts with %s", kind, key, previous.module)
		}
		seen[identity] = ownership{kind: kind, key: key, module: module}
		return nil
	}
	for _, module := range modules {
		for _, key := range module.Contributions.Routes {
			if err := add(module.ID, "route", key); err != nil {
				return err
			}
		}
		for _, key := range module.Contributions.Pages {
			if err := add(module.ID, "page", key); err != nil {
				return err
			}
		}
		for _, key := range module.Contributions.Navigation {
			if err := add(module.ID, "navigation", key); err != nil {
				return err
			}
		}
		for _, key := range module.Contributions.Permissions {
			if err := add(module.ID, "permission", key); err != nil {
				return err
			}
		}
		for _, key := range module.Contributions.ConfigNamespaces {
			if err := add(module.ID, "config", key); err != nil {
				return err
			}
		}
	}
	return nil
}

func cloneModule(module Module) Module {
	module.DependsOn = append([]string(nil), module.DependsOn...)
	module.Provides = append([]Capability(nil), module.Provides...)
	module.Requires = append([]Capability(nil), module.Requires...)
	module.Contributions.Routes = append([]string(nil), module.Contributions.Routes...)
	module.Contributions.Pages = append([]string(nil), module.Contributions.Pages...)
	module.Contributions.Navigation = append([]string(nil), module.Contributions.Navigation...)
	module.Contributions.Permissions = append([]string(nil), module.Contributions.Permissions...)
	module.Contributions.ConfigNamespaces = append([]string(nil), module.Contributions.ConfigNamespaces...)
	return module
}
