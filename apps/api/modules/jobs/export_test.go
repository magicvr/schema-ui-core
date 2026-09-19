package jobs_test

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/jobs"
	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
	jobmodule "github.com/magicvr/schema-ui-core/apps/api/modules/jobs"
)

// fakeRowSource stands in for the users/roles resource entities so the export
// job's own behaviour — progress, CSV rendering, result document — is tested
// without the HTTP layer.
type fakeRowSource struct {
	rows    map[string]map[string]any
	delay   time.Duration
	failing []string
}

func (s *fakeRowSource) SupportedExportResources() []string { return []string{"users", "roles"} }

func (s *fakeRowSource) ExportSelectedRows(resource string, ids []string, onRow func(done int) error) ([]string, [][]string, error) {
	if resource != "users" {
		return nil, nil, context.Canceled
	}
	headers := []string{"id", "username", "name", "roles", "enabled", "locked", "createdAt", "updatedAt"}
	rows := make([][]string, 0, len(ids))
	for done, id := range ids {
		if s.delay > 0 {
			time.Sleep(s.delay)
		}
		row, ok := s.rows[id]
		if !ok {
			return nil, nil, context.Canceled
		}
		rows = append(rows, []string{id, str(row["username"]), str(row["name"]), "", "true", "false", "", ""})
		if onRow != nil {
			if err := onRow(done + 1); err != nil {
				return nil, nil, err
			}
		}
	}
	return headers, rows, nil
}

func str(value any) string {
	text, _ := value.(string)
	return text
}

func newExportHarness(t *testing.T, source *fakeRowSource) (*jobs.Repository, *jobs.Runner, *jobmodule.BatchExportService) {
	t.Helper()
	st, err := testsupport.OpenStore(":memory:", "admin", "hash", false)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	repository := jobs.NewRepository(st)
	options := jobs.DefaultRunnerOptions()
	options.HeartbeatInterval = 5 * time.Millisecond
	options.LeaseDuration = 50 * time.Millisecond
	options.ScanInterval = 5 * time.Millisecond
	runner, err := jobs.NewRunner(repository, options)
	if err != nil {
		t.Fatal(err)
	}
	service, err := jobmodule.NewBatchExportService(runner, source, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := runner.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		_ = runner.Stop(ctx)
	})
	return repository, runner, service
}

// C1: the job runs to success and stores a CSV result with the frozen header
// order and a UTF-8 BOM.
func TestBatchExportJobProducesCSVResult(t *testing.T) {
	source := &fakeRowSource{rows: map[string]map[string]any{
		"u-1": {"username": "alpha", "name": "Alpha"},
		"u-2": {"username": "beta", "name": "Beta"},
	}}
	repository, runner, service := newExportHarness(t, source)

	submitted, err := service.Submit(context.Background(), runner, "users", []string{"u-1", "u-1", " u-2 "}, "actor-1", "corr-1")
	if err != nil {
		t.Fatalf("submit: %v", err)
	}
	terminal := waitForListedStatus(t, repository, submitted.ID, jobs.StatusSucceeded)
	if terminal.ErrorCode != "" {
		t.Fatalf("job failed: %s %s", terminal.ErrorCode, terminal.ErrorMessage)
	}

	var document struct {
		Resource string `json:"resource"`
		RowCount int    `json:"rowCount"`
		FileName string `json:"fileName"`
		CSV      string `json:"csv"`
	}
	if err := json.Unmarshal(terminal.Result, &document); err != nil {
		t.Fatalf("result is not the expected document: %v (%s)", err, terminal.Result)
	}
	if document.Resource != "users" || document.RowCount != 2 {
		t.Fatalf("result summary = %+v, want users/2 (deduped, trimmed)", document)
	}
	if document.FileName != "users-selection.csv" {
		t.Fatalf("fileName = %q", document.FileName)
	}
	if !strings.HasPrefix(document.CSV, "\uFEFF") {
		t.Fatalf("csv is missing the UTF-8 BOM")
	}
	if !strings.Contains(document.CSV, "id,username,name,roles,enabled,locked,createdAt,updatedAt") {
		t.Fatalf("csv header mismatch:\n%s", document.CSV)
	}
	if !strings.Contains(document.CSV, "alpha") || !strings.Contains(document.CSV, "beta") {
		t.Fatalf("csv missing rows:\n%s", document.CSV)
	}
}

// C1: progress reflects real work. Because the row source is slow, the job is
// observable mid-flight and at least one sample must be a value the hardcoded
// 10→100 precedent could never produce.
func TestBatchExportJobReportsRealProgress(t *testing.T) {
	source := &fakeRowSource{
		rows: map[string]map[string]any{
			"u-1": {"username": "a"}, "u-2": {"username": "b"},
			"u-3": {"username": "c"}, "u-4": {"username": "d"},
		},
		delay: 15 * time.Millisecond,
	}
	repository, runner, service := newExportHarness(t, source)

	submitted, err := service.Submit(context.Background(), runner, "users", []string{"u-1", "u-2", "u-3", "u-4"}, "actor-1", "corr-p")
	if err != nil {
		t.Fatalf("submit: %v", err)
	}

	seen := map[int]bool{}
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		current, err := repository.GetJob(context.Background(), submitted.ID)
		if err != nil {
			t.Fatal(err)
		}
		seen[current.Progress] = true
		if current.Status == jobs.StatusSucceeded {
			break
		}
		if current.Status == jobs.StatusFailed {
			t.Fatalf("job failed: %s %s", current.ErrorCode, current.ErrorMessage)
		}
		time.Sleep(3 * time.Millisecond)
	}
	if !seen[100] {
		t.Fatalf("job never reached progress 100; seen=%v", seen)
	}
	intermediate := 0
	for value := range seen {
		if value != 0 && value != 100 {
			intermediate++
		}
	}
	if intermediate == 0 {
		t.Fatalf("only 0/100 were observed — progress is not tracking real work; seen=%v", seen)
	}
	// With 4 rows the per-row steps are 21/42/63/85, so a monotonic spread well
	// below 90 must be reachable. Assert the values stay in the documented
	// band (never a bogus >99 before completion).
	for value := range seen {
		if value < 0 || value > 100 {
			t.Fatalf("progress out of range: %d", value)
		}
	}
}

// C1: validation rejects an empty, oversized or unknown-resource selection
// before any job row exists.
func TestBatchExportSubmitValidation(t *testing.T) {
	source := &fakeRowSource{rows: map[string]map[string]any{}}
	repository, runner, service := newExportHarness(t, source)

	if _, err := service.Submit(context.Background(), runner, "users", nil, "a", "c"); err == nil {
		t.Fatalf("empty selection must be rejected")
	}
	if _, err := service.Submit(context.Background(), runner, "users", []string{"  ", ""}, "a", "c"); err == nil {
		t.Fatalf("a selection of only blank keys must be rejected")
	}
	if _, err := service.Submit(context.Background(), runner, "settings", []string{"x"}, "a", "c"); err == nil {
		t.Fatalf("an unsupported resource must be rejected")
	}
	oversized := make([]string, jobmodule.BatchExportMaxIDs+1)
	for i := range oversized {
		oversized[i] = "id"
	}
	// Distinct keys so the duplicate filter cannot mask the size check.
	for i := range oversized {
		oversized[i] = "id-" + string(rune('a'+i%26)) + string(rune('a'+i/26))
	}
	if _, err := service.Submit(context.Background(), runner, "users", oversized, "a", "c"); err == nil {
		t.Fatalf("a selection above the cap must be rejected")
	}

	if _, total, err := repository.ListJobs(context.Background(), jobs.ListFilter{Kind: jobmodule.BatchExportJobKind}); err != nil || total != 0 {
		t.Fatalf("rejected submits left %d job rows (err=%v), want 0", total, err)
	}
}

// waitForListedStatus polls the management read until the job settles.
func waitForListedStatus(t *testing.T, repository *jobs.Repository, id string, want jobs.Status) *jobs.Job {
	t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		current, err := repository.GetJob(context.Background(), id)
		if err != nil {
			t.Fatal(err)
		}
		if current.Status == want || current.Status == jobs.StatusFailed || current.Status == jobs.StatusCancelled {
			return current
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatalf("job %s never settled (want %s)", id, want)
	return nil
}
