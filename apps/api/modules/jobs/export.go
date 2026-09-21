package jobs

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/jobs"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

// BatchExportJobKind is the first real batch operation carried by the durable
// Job runtime (GOAL-004 R3 · VP-038 exit criterion 3). It exports the SELECTED
// rows of a supported resource as CSV and stores the CSV as the job result.
const BatchExportJobKind = "jobs.batch-export"

// BatchExportMaxIDs bounds one selection, keeping a single job a bounded unit
// of work (the synchronous batch-delete body cap is 4 KiB, ≈ hundreds of ids).
const BatchExportMaxIDs = 500

// batchExportPayload is the durable job input. Only the resource and the ids
// are carried; the actor lives on the job row itself.
type batchExportPayload struct {
	Resource string   `json:"resource"`
	IDs      []string `json:"ids"`
}

// batchExportResult is the job result document. `csv` carries the export bytes
// (the Job result column must hold valid JSON, so the CSV rides as a string);
// the remaining fields let the result center render counts and a filename
// without re-parsing the CSV.
type batchExportResult struct {
	Resource string `json:"resource"`
	RowCount int    `json:"rowCount"`
	FileName string `json:"fileName"`
	CSV      string `json:"csv"`
}

// ExportRowSource builds the CSV rows for an explicit selection. Implemented by
// the handler package over the same column sets and formula neutralization the
// synchronous export uses, so the two surfaces cannot drift apart (GOAL-004
// D-001: the batch export reuses the frozen column order rather than restating
// it). onRow reports rows read so the job can publish real progress.
type ExportRowSource interface {
	ExportSelectedRows(resource string, ids []string, onRow func(done int) error) (headers []string, rows [][]string, err error)
	SupportedExportResources() []string
}

// BatchExportService binds the Job runtime to batch export. The kind is
// registered at construction; the composition root builds this before the
// runner starts (R1 D-001 §1.2 K-5).
type BatchExportService struct {
	rows       ExportRowSource
	operations OperationRecorder
	now        func() time.Time
}

// OperationRecorder records a batch-export audit event. nil disables recording.
type OperationRecorder interface {
	RecordBatchExport(ctx context.Context, jobID, resource string, rowCount int, actorName string)
}

// NewBatchExportService registers the batch-export handler on the runner.
// runner may be nil in read-only compositions: the service is then inert and
// Submit/Register are unavailable, which keeps the module composable without a
// Job runtime.
func NewBatchExportService(runner *jobs.Runner, rows ExportRowSource, operations OperationRecorder) (*BatchExportService, error) {
	if rows == nil {
		return nil, jobs.ErrInvalid
	}
	s := &BatchExportService{rows: rows, operations: operations, now: time.Now}
	if runner == nil {
		return s, nil
	}
	if err := runner.Register(BatchExportJobKind, s.runExport); err != nil {
		return nil, err
	}
	return s, nil
}

// SupportedResources mirrors the export denominator of the synchronous export.
func (s *BatchExportService) SupportedResources() []string { return s.rows.SupportedExportResources() }

// Supported reports whether resource can be batch-exported.
func (s *BatchExportService) Supported(resource string) bool {
	for _, supported := range s.rows.SupportedExportResources() {
		if supported == resource {
			return true
		}
	}
	return false
}

// NormalizeExportIDs trims, drops empties and de-duplicates while preserving
// order — the same normalization the synchronous batch action applies, so both
// paths agree on what "the selection" means. It rejects an empty selection and
// a selection larger than BatchExportMaxIDs.
func NormalizeExportIDs(ids []string) ([]string, error) {
	if len(ids) == 0 {
		return nil, jobs.ErrInvalid
	}
	seen := make(map[string]bool, len(ids))
	out := make([]string, 0, len(ids))
	for _, raw := range ids {
		id := strings.TrimSpace(raw)
		if id == "" || seen[id] {
			continue
		}
		seen[id] = true
		out = append(out, id)
	}
	if len(out) == 0 || len(out) > BatchExportMaxIDs {
		return nil, jobs.ErrInvalid
	}
	return out, nil
}

// Submit validates the request and enqueues the job. Validation runs before the
// row is created so a rejected request never leaves a queued job behind.
func (s *BatchExportService) Submit(ctx context.Context, runner *jobs.Runner, resource string, ids []string, actorID, correlationID string) (*jobs.Job, error) {
	if runner == nil || !s.Supported(resource) {
		return nil, jobs.ErrInvalid
	}
	normalized, err := NormalizeExportIDs(ids)
	if err != nil {
		return nil, err
	}
	now := s.now().UTC()
	id, err := jobs.NewID(now)
	if err != nil {
		return nil, err
	}
	payload, err := json.Marshal(batchExportPayload{Resource: resource, IDs: normalized})
	if err != nil {
		return nil, err
	}
	return runner.Submit(ctx, jobs.CreateInput{
		ID: id, Kind: BatchExportJobKind, Payload: payload,
		ActorID: actorID, CorrelationID: correlationID,
		MaxAttempts: jobs.DefaultMaxAttempts, Now: now,
	})
}

// runExport reads the selected rows, renders CSV and returns the commit that
// stores the bytes as the job result.
//
// Progress is reported across the actual selection (R1 D-001 §1.2 K-6: the
// existing wallet precedent only reports a hardcoded 10 → 100, so observable
// progress has to be implemented here). reporter.Progress accepts 0..99; the
// terminal 100 is set by the atomic completion.
func (s *BatchExportService) runExport(_ context.Context, job jobs.Job, reporter jobs.Reporter) (jobs.CommitFunc, error) {
	var payload batchExportPayload
	if err := json.Unmarshal(job.Payload, &payload); err != nil {
		return nil, fmt.Errorf("decode batch export payload: %w", err)
	}
	if !s.Supported(payload.Resource) {
		return nil, fmt.Errorf("unsupported batch export resource %q", payload.Resource)
	}
	if len(payload.IDs) == 0 {
		return nil, fmt.Errorf("batch export selection is empty")
	}

	total := len(payload.IDs)
	// Reserve 85..99 for rendering and completion so every reported value stays
	// strictly below 100 (the completion transition owns 100).
	stepToPercent := func(done int) int { return done * 85 / total }

	headers, rows, err := s.rows.ExportSelectedRows(payload.Resource, payload.IDs, func(done int) error {
		if reporter.Cancelled() {
			return jobs.ErrInvalid
		}
		return reporter.Progress(stepToPercent(done))
	})
	if err != nil {
		return nil, err
	}
	if reporter.Cancelled() {
		return nil, jobs.ErrInvalid
	}
	if err := reporter.Progress(90); err != nil {
		return nil, err
	}
	rendered, err := RenderExportCSV(headers, rows)
	if err != nil {
		return nil, err
	}
	if err := reporter.Progress(99); err != nil {
		return nil, err
	}

	result := batchExportResult{
		Resource: payload.Resource,
		RowCount: len(rows),
		FileName: payload.Resource + "-selection.csv",
		CSV:      rendered,
	}
	if s.operations != nil {
		s.operations.RecordBatchExport(context.Background(), job.ID, payload.Resource, len(rows), job.ActorID)
	}
	return func(kernel.Tx) (json.RawMessage, error) {
		return json.Marshal(result)
	}, nil
}

// RenderExportCSV writes the export with the same UTF-8 BOM and RFC 4180
// escaping the synchronous export uses, so a batch download opens identically
// in Excel.
func RenderExportCSV(headers []string, rows [][]string) (string, error) {
	var out strings.Builder
	out.WriteString("\uFEFF")
	writer := csv.NewWriter(&out)
	if err := writer.Write(headers); err != nil {
		return "", err
	}
	for _, row := range rows {
		if err := writer.Write(row); err != nil {
			return "", err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", err
	}
	return out.String(), nil
}
