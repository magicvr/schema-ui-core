package kernel

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestBuiltinProfilesResolveDeterministically(t *testing.T) {
	registry, err := NewRegistry(BuiltinModules())
	if err != nil {
		t.Fatal(err)
	}

	mvp, err := ResolveProfile("mvp", nil)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := mvp.Modules, []string{"core.server-registration", "core.auth-session", "core.manifest-route", "core.navigation-capability", "core.schema-render", "core.operationlog", "admin.users", "admin.roles", "admin.account", "admin.dashboard", "admin.notifications"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("mvp modules = %v, want %v", got, want)
	}
	plan, err := registry.Resolve(mvp.Modules)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := plan.IDs(), []string{"core.server-registration", "core.auth-session", "core.schema-render", "core.manifest-route", "core.navigation-capability", "core.operationlog", "admin.account", "admin.dashboard", "admin.notifications", "admin.roles", "admin.users"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("mvp plan = %v, want %v", got, want)
	}

	admin, err := ResolveProfile("admin", nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := registry.Resolve(admin.Modules); err != nil {
		t.Fatal(err)
	}

	// W2 (GOAL-003 / workspace-010): demo = mvp capability set + dev.examples.
	demo, err := ResolveProfile("demo", nil)
	if err != nil {
		t.Fatalf("resolve demo: %v", err)
	}
	if got, want := demo.Modules, []string{"core.server-registration", "core.auth-session", "core.manifest-route", "core.navigation-capability", "core.schema-render", "core.operationlog", "admin.users", "admin.roles", "dev.examples", "admin.account", "admin.dashboard", "admin.notifications"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("demo modules = %v, want %v", got, want)
	}
	if _, err := registry.Resolve(demo.Modules); err != nil {
		t.Fatalf("resolve demo plan: %v", err)
	}
}

// TestDemoProfileIsNonProduction verifies demo is never used as a production
// default and mvp/admin keep excluding dev.examples (W1 hygiene, GOAL-003 S3).
func TestDemoProfileIsNonProduction(t *testing.T) {
	for _, name := range []ProfileName{ProfileMVP, ProfileAdmin} {
		resolution, err := ResolveProfile(string(name), nil)
		if err != nil {
			t.Fatalf("%s resolve: %v", name, err)
		}
		for _, id := range resolution.Modules {
			if id == "dev.examples" {
				t.Fatalf("%s default profile must not include dev.examples (W1 S5)", name)
			}
		}
	}
	demo, err := ResolveProfile(string(ProfileDemo), nil)
	if err != nil {
		t.Fatal(err)
	}
	if demo.Source != "profile.default" {
		t.Fatalf("demo source = %q, want profile.default", demo.Source)
	}
}

func TestProfileOverrideAndCustomRequireExplicitModules(t *testing.T) {
	resolved, err := ResolveProfile("admin", []string{"core.server-registration"})
	if err != nil {
		t.Fatal(err)
	}
	if resolved.Source != "modules.enabled" || !reflect.DeepEqual(resolved.Modules, []string{"core.server-registration"}) {
		t.Fatalf("unexpected override: %+v", resolved)
	}
	if _, err := ResolveProfile("custom", nil); err == nil {
		t.Fatal("custom profile must require explicit modules")
	}
}

func TestRegistryRejectsUnknownMissingCycleConflictAndCapability(t *testing.T) {
	base := Module{ID: "base", Version: "1", KernelAPIRange: ">=2 <3"}
	registry, err := NewRegistry([]Module{base})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := registry.Resolve([]string{"unknown"}); !hasCode(err, CodeModuleUnknown) {
		t.Fatalf("unknown error = %v", err)
	}

	missing, err := NewRegistry([]Module{{ID: "child", Version: "1", KernelAPIRange: ">=2 <3", DependsOn: []string{"base"}}})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := missing.Resolve([]string{"child"}); !hasCode(err, CodeModuleDependencyUnknown) {
		t.Fatalf("unknown dependency error = %v", err)
	}

	cycle, err := NewRegistry([]Module{
		{ID: "a", Version: "1", KernelAPIRange: ">=2 <3", DependsOn: []string{"b"}},
		{ID: "b", Version: "1", KernelAPIRange: ">=2 <3", DependsOn: []string{"a"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := cycle.Resolve([]string{"a", "b"}); !hasCode(err, CodeModuleDependencyCycle) {
		t.Fatalf("cycle error = %v", err)
	}

	conflict, err := NewRegistry([]Module{
		{ID: "a", Version: "1", KernelAPIRange: ">=2 <3", Contributions: ContributionKeys{Routes: []string{"/same"}}},
		{ID: "b", Version: "1", KernelAPIRange: ">=2 <3", Contributions: ContributionKeys{Routes: []string{"/same"}}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := conflict.Resolve([]string{"a", "b"}); !hasCode(err, CodeModuleContributionConflict) {
		t.Fatalf("conflict error = %v", err)
	}

	capability, err := NewRegistry([]Module{{ID: "needs", Version: "1", KernelAPIRange: ">=2 <3", Requires: []Capability{"missing"}}})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := capability.Resolve([]string{"needs"}); !hasCode(err, CodeModuleCapabilityMissing) {
		t.Fatalf("capability error = %v", err)
	}
}

func TestRegistryValidatesKernelAPIRanges(t *testing.T) {
	if _, err := NewRegistry([]Module{{ID: "invalid", Version: "1", KernelAPIRange: ">=bad"}}); !hasCode(err, CodeModuleAPIRangeInvalid) {
		t.Fatalf("invalid range error = %v", err)
	}
	if _, err := NewRegistry([]Module{{ID: "old", Version: "1", KernelAPIRange: "<2"}}); !hasCode(err, CodeModuleAPIMismatch) {
		t.Fatalf("incompatible range error = %v", err)
	}
}

func TestRuntimeStartsTopologicallyAndCleansUpOnFailure(t *testing.T) {
	var events []string
	registry, err := NewRegistry([]Module{
		{ID: "a", Version: "1", KernelAPIRange: ">=2 <3", Hooks: Hooks{Start: func(context.Context) error { events = append(events, "a+"); return nil }, Stop: func(context.Context) error { events = append(events, "a-"); return nil }}},
		{ID: "b", Version: "1", KernelAPIRange: ">=2 <3", DependsOn: []string{"a"}, Hooks: Hooks{Start: func(context.Context) error { events = append(events, "b+"); return errors.New("boom") }, Stop: func(context.Context) error { events = append(events, "b-"); return nil }}},
	})
	if err != nil {
		t.Fatal(err)
	}
	plan, err := registry.Resolve([]string{"b", "a"})
	if err != nil {
		t.Fatal(err)
	}
	runtime := NewRuntime(plan)
	if err := runtime.Start(context.Background()); !hasCode(err, CodeLifecycleStartFailed) {
		t.Fatalf("start error = %v", err)
	}
	if got, want := strings.Join(events, ","), "a+,b+,a-"; got != want {
		t.Fatalf("events = %q, want %q", got, want)
	}
}

func hasCode(err error, code ErrorCode) bool {
	var kernelErr *Error
	return errors.As(err, &kernelErr) && kernelErr.Code == code
}