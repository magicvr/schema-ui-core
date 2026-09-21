package jobs

import (
	"context"
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/internal/handler"
	"github.com/magicvr/schema-ui-core/apps/api/internal/jobs"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

// GOAL-005 R4: the provider's descriptor must declare exactly what it serves.
// The kernel rejects a mismatch between the plan descriptor and the provider
// (MODULE_API_MISMATCH), but that check can only compare the two DECLARATIONS.
// This test pins the declaration against the optional surfaces actually bound,
// including the conditional route sets a read-only composition produces.

type stubSubmitter struct{}

func (stubSubmitter) Supported(string) bool { return true }

func (stubSubmitter) SubmitBatchExport(context.Context, string, []string, string, string) (*jobs.Job, error) {
	return nil, nil
}

type stubActions struct{}

func (stubActions) CancelAny(context.Context, string) (*jobs.Job, error) { return nil, nil }

func (stubActions) RetryAny(context.Context, string) (*jobs.Job, error) { return nil, nil }

func TestDescriptorDeclaresEveryBoundSurface(t *testing.T) {
	cases := []struct {
		name      string
		submitter handler.JobBatchSubmitter
		actions   handler.JobActions
		want      []string
	}{
		{
			name: "read only",
			want: []string{"GET /api/jobs", "GET /api/jobs/{id}", "GET /api/jobs/{id}/result"},
		},
		{
			name:      "batch export only",
			submitter: stubSubmitter{},
			want:      []string{"GET /api/jobs", "GET /api/jobs/{id}", "GET /api/jobs/{id}/result", "POST /api/jobs/batch-export"},
		},
		{
			name:    "actions only",
			actions: stubActions{},
			want:    []string{"GET /api/jobs", "GET /api/jobs/{id}", "GET /api/jobs/{id}/result", "POST /api/jobs/{id}/cancel", "POST /api/jobs/{id}/retry"},
		},
		{
			name:      "full result center",
			submitter: stubSubmitter{},
			actions:   stubActions{},
			want:      []string{"GET /api/jobs", "GET /api/jobs/{id}", "GET /api/jobs/{id}/result", "POST /api/jobs/batch-export", "POST /api/jobs/{id}/cancel", "POST /api/jobs/{id}/retry"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			module := New(nil, nil, tc.submitter, tc.actions).Descriptor()
			got := module.Contributions.Routes
			if len(got) != len(tc.want) {
				t.Fatalf("routes = %v, want %v", got, tc.want)
			}
			for i, route := range tc.want {
				if got[i] != route {
					t.Fatalf("routes[%d] = %q, want %q", i, got[i], route)
				}
			}
			// jobs.write is declared exactly when a mutating surface is bound: a
			// read-only composition must not advertise a write key.
			hasWrite := false
			for _, permission := range module.Contributions.Permissions {
				if permission == "jobs.write" {
					hasWrite = true
				}
			}
			wantWrite := tc.submitter != nil || tc.actions != nil
			if hasWrite != wantWrite {
				t.Fatalf("jobs.write declared = %v, want %v (%v)", hasWrite, wantWrite, module.Contributions.Permissions)
			}
		})
	}
}

// The plan descriptor in kernel.BuiltinModules must carry the same route set as
// a fully composed provider: the plan is what a profile negotiates against, so a
// stale plan makes the module unreachable (or over-declared) at composition.
func TestPlanDescriptorMatchesTheFullProvider(t *testing.T) {
	var plan *kernel.Module
	for _, module := range kernel.BuiltinModules() {
		if module.ID == ModuleID {
			candidate := module
			plan = &candidate
			break
		}
	}
	if plan == nil {
		t.Fatalf("%s is not in the built-in module plan", ModuleID)
	}
	provider := New(nil, nil, stubSubmitter{}, stubActions{}).Descriptor()
	if len(plan.Contributions.Routes) != len(provider.Contributions.Routes) {
		t.Fatalf("plan routes = %v, provider routes = %v", plan.Contributions.Routes, provider.Contributions.Routes)
	}
	for i, route := range provider.Contributions.Routes {
		if plan.Contributions.Routes[i] != route {
			t.Fatalf("plan routes[%d] = %q, provider = %q", i, plan.Contributions.Routes[i], route)
		}
	}
	if len(plan.Contributions.Permissions) != 2 {
		t.Fatalf("plan permissions = %v, want jobs.read + jobs.write", plan.Contributions.Permissions)
	}
}
