package composition

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/config"
	"github.com/magicvr/schema-ui-core/apps/api/internal/ratelimit"
	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/modules/authsession"
	"github.com/magicvr/schema-ui-core/apps/api/modules/operationlog"
	settingsrepository "github.com/magicvr/schema-ui-core/apps/api/modules/settings/repository"
)

// GOAL-003 R2 (D-001 §4, R-1): the durable Job runtime used to be enabled only
// when admin.wallet was in the plan, so a profile carrying admin.jobs without
// admin.wallet got a runner that silently never started. The flag now flips for
// either module. This test pins the jobs-without-wallet case, which the admin
// profile (it contains wallet) cannot distinguish — a revert of that branch
// would otherwise leave every other test green.
func TestJobRuntimeEnablesForJobsWithoutWallet(t *testing.T) {
	plan := jobsOnlyPlan(t)
	if plan.HasModule("admin.wallet") {
		t.Fatalf("fixture is wrong: plan unexpectedly contains admin.wallet")
	}
	if !plan.HasModule("admin.jobs") {
		t.Fatalf("fixture is wrong: plan is missing admin.jobs")
	}

	st, err := testsupport.OpenStore(":memory:", "admin", "hash", false)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })

	runtime, err := newJobRuntime(st)
	if err != nil {
		t.Fatal(err)
	}
	if runtime.enabled.Load() {
		t.Fatalf("fixture is wrong: jobRuntime starts enabled")
	}

	cachePort, err := newCache(&config.Config{DBPath: "test.db"})
	if err != nil {
		t.Fatal(err)
	}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	if _, err := newMux(
		&config.Config{DBPath: "test.db"},
		auth.New([]byte("test-secret"), 0, 0, st, true),
		st,
		authsession.NewRepository(st),
		operationlog.NewRepository(st),
		settingsrepository.New(st),
		plan,
		&readinessGate{},
		jwtSecret("test-secret"),
		runtime,
		nil,
		logger,
		cachePort,
		newEventBus(&config.Config{DBPath: "test.db"}, logger),
		ratelimit.NewProvider(),
		nil,
	); err != nil {
		t.Fatal(err)
	}

	if !runtime.enabled.Load() {
		t.Fatalf("jobRuntime.enabled = false for a plan with admin.jobs but no admin.wallet; " +
			"the R-1 fix regressed and the runner would never start")
	}
	// The flag is what Start() gates on, so prove the runner really starts
	// rather than only that the flag was set.
	if err := runtime.Start(); err != nil {
		t.Fatalf("jobRuntime.Start() = %v, want a started runner", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := runtime.Stop(ctx); err != nil {
		t.Fatalf("jobRuntime.Stop() = %v", err)
	}
}

// jobsOnlyPlan resolves a custom profile containing admin.jobs and the
// dependencies its Descriptor declares, and nothing else.
func jobsOnlyPlan(t *testing.T) kernel.Plan {
	t.Helper()
	registry, err := kernel.NewRegistry(kernel.BuiltinModules())
	if err != nil {
		t.Fatal(err)
	}
	plan, err := registry.Resolve([]string{
		"core.server-registration",
		"core.auth-session",
		"core.navigation-capability",
		"core.schema-render",
		"core.manifest-route",
		"core.operationlog",
		"admin.jobs",
	})
	if err != nil {
		t.Fatalf("resolve jobs-only plan: %v", err)
	}
	return plan
}
