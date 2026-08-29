package kernel

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
)

type lifecycleFailures struct {
	start map[string]error
	ready map[string]error
	stop  map[string]error
}

func resolvedLifecyclePlan(t *testing.T, profile ProfileName) Plan {
	t.Helper()
	resolution, err := ResolveProfile(string(profile), nil)
	if err != nil {
		t.Fatal(err)
	}
	registry, err := NewRegistry(BuiltinModules())
	if err != nil {
		t.Fatal(err)
	}
	plan, err := registry.Resolve(resolution.Modules)
	if err != nil {
		t.Fatal(err)
	}
	return plan
}

func instrumentLifecycle(plan Plan, events *[]string, failures lifecycleFailures) Plan {
	for i := range plan.Modules {
		moduleID := plan.Modules[i].ID
		plan.Modules[i].Hooks = Hooks{
			Start: func(context.Context) error {
				*events = append(*events, "start:"+moduleID)
				return failures.start[moduleID]
			},
			Ready: func(context.Context) error {
				*events = append(*events, "ready:"+moduleID)
				return failures.ready[moduleID]
			},
			Stop: func(context.Context) error {
				*events = append(*events, "stop:"+moduleID)
				return failures.stop[moduleID]
			},
		}
	}
	return plan
}

func phaseEvents(phase string, ids []string) []string {
	events := make([]string, 0, len(ids))
	for _, id := range ids {
		events = append(events, phase+":"+id)
	}
	return events
}

func reversePhaseEvents(phase string, ids []string) []string {
	events := make([]string, 0, len(ids))
	for i := len(ids) - 1; i >= 0; i-- {
		events = append(events, phase+":"+ids[i])
	}
	return events
}

func requireLifecycleError(t *testing.T, err error, code ErrorCode, moduleID string) *Error {
	t.Helper()
	var kernelErr *Error
	if !errors.As(err, &kernelErr) || kernelErr.Code != code || kernelErr.ModuleID != moduleID {
		t.Fatalf("error = %v, want %s [%s]", err, code, moduleID)
	}
	return kernelErr
}

func TestDualProfileLifecycleMatrix(t *testing.T) {
	for _, profile := range []ProfileName{ProfileMVP, ProfileAdmin} {
		t.Run(string(profile), func(t *testing.T) {
			base := resolvedLifecyclePlan(t, profile)
			ids := base.IDs()
			failureIndex := len(ids) / 2
			failureID := ids[failureIndex]
			ctx := context.Background()

			t.Run("success", func(t *testing.T) {
				var events []string
				runtime := NewRuntime(instrumentLifecycle(base, &events, lifecycleFailures{}))
				if err := runtime.Start(ctx); err != nil {
					t.Fatal(err)
				}
				if err := runtime.Ready(ctx); err != nil {
					t.Fatal(err)
				}
				if err := runtime.Stop(ctx); err != nil {
					t.Fatal(err)
				}
				want := append(phaseEvents("start", ids), phaseEvents("ready", ids)...)
				want = append(want, reversePhaseEvents("stop", ids)...)
				if !reflect.DeepEqual(events, want) {
					t.Fatalf("events = %v, want %v", events, want)
				}
			})

			t.Run("start failure cleans prior modules", func(t *testing.T) {
				var events []string
				failures := lifecycleFailures{
					start: map[string]error{failureID: errors.New("start sentinel")},
					stop:  map[string]error{ids[0]: errors.New("cleanup sentinel")},
				}
				runtime := NewRuntime(instrumentLifecycle(base, &events, failures))
				kernelErr := requireLifecycleError(t, runtime.Start(ctx), CodeLifecycleStartFailed, failureID)
				if !strings.Contains(kernelErr.Detail, "cleanup failed") {
					t.Fatalf("detail = %q", kernelErr.Detail)
				}
				want := phaseEvents("start", ids[:failureIndex+1])
				want = append(want, reversePhaseEvents("stop", ids[:failureIndex])...)
				if !reflect.DeepEqual(events, want) {
					t.Fatalf("events = %v, want %v", events, want)
				}
				if err := runtime.Stop(ctx); err != nil {
					t.Fatalf("repeat stop after cleanup: %v", err)
				}
			})

			t.Run("ready failure cleans all started modules", func(t *testing.T) {
				var events []string
				readyErr := &Error{Code: CodeModuleInvalid, ModuleID: "structured-ready", Detail: "ready sentinel"}
				failures := lifecycleFailures{
					ready: map[string]error{failureID: readyErr},
					stop:  map[string]error{ids[len(ids)-1]: errors.New("cleanup sentinel")},
				}
				runtime := NewRuntime(instrumentLifecycle(base, &events, failures))
				if err := runtime.Start(ctx); err != nil {
					t.Fatal(err)
				}
				kernelErr := requireLifecycleError(t, runtime.Ready(ctx), CodeModuleInvalid, "structured-ready")
				if !strings.Contains(kernelErr.Detail, "ready sentinel") || !strings.Contains(kernelErr.Detail, "cleanup failed") {
					t.Fatalf("detail = %q", kernelErr.Detail)
				}
				want := append(phaseEvents("start", ids), phaseEvents("ready", ids[:failureIndex+1])...)
				want = append(want, reversePhaseEvents("stop", ids)...)
				if !reflect.DeepEqual(events, want) {
					t.Fatalf("events = %v, want %v", events, want)
				}
				if err := runtime.Stop(ctx); err != nil {
					t.Fatalf("repeat stop after ready cleanup: %v", err)
				}
			})

			t.Run("stop continues and returns first reverse error", func(t *testing.T) {
				var events []string
				failures := lifecycleFailures{stop: map[string]error{
					ids[len(ids)-1]: errors.New("first reverse error"),
					ids[0]:          errors.New("later reverse error"),
				}}
				runtime := NewRuntime(instrumentLifecycle(base, &events, failures))
				if err := runtime.Start(ctx); err != nil {
					t.Fatal(err)
				}
				if err := runtime.Ready(ctx); err != nil {
					t.Fatal(err)
				}
				requireLifecycleError(t, runtime.Stop(ctx), CodeLifecycleStopFailed, ids[len(ids)-1])
				want := append(phaseEvents("start", ids), phaseEvents("ready", ids)...)
				want = append(want, reversePhaseEvents("stop", ids)...)
				if !reflect.DeepEqual(events, want) {
					t.Fatalf("events = %v, want %v", events, want)
				}
				if err := runtime.Stop(ctx); err != nil {
					t.Fatalf("repeat stop: %v", err)
				}
			})
		})
	}
}
