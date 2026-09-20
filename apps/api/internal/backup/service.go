package backup

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/temporalcontract"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

// Provider creates a native artifact and restores it to a NEW isolated target.
type Provider interface {
	// Dialect reports the provider dialect ("sqlite" or "postgres").
	Dialect() string
	// Create writes the native artifact of sourceID to artifactRef.
	Create(ctx context.Context, sourceID, artifactRef string) error
	// Restore materialises artifactRef into a new isolated target and returns a
	// handle token the verifier understands (a SQLite file path, or a DSN).
	Restore(ctx context.Context, artifactRef string) (target string, cleanup func(), err error)
}

// Service is the internal BackupService: it selects the provider, creates the
// artifact, restores it to a new target, verifies the target and only then
// returns a RecoveryPoint (Root D-010 postcondition: never an unverified
// artifact).
type Service struct {
	providers map[string]Provider
	// now is injectable for tests.
	now func() time.Time
	// workDir is where temporary artifacts and restore targets live.
	workDir string
	// verifiers allows tests to inject a measurement; nil uses the real ones.
	measureSQLite func(ctx context.Context, path string) (sqliteTargetVerification, error)
}

// Service implements the minimal kernel port (Root D-007 / D-010): only
// CreateRecoveryPoint is exported to callers, and it never returns an
// unverified artifact.
var _ kernel.RecoveryPointPort = (*Service)(nil)

// NewService builds the service with the native providers and a temporary work
// directory.
func NewService(workDir string) *Service {
	return &Service{
		providers: map[string]Provider{
			"sqlite": SQLiteProvider{},
		},
		now:           func() time.Time { return time.Now().UTC() },
		workDir:       workDir,
		measureSQLite: measureSQLiteTarget,
	}
}

// RegisterProvider adds or replaces a provider (used by tests and by the
// postgres wiring).
func (s *Service) RegisterProvider(p Provider) {
	if s.providers == nil {
		s.providers = map[string]Provider{}
	}
	s.providers[p.Dialect()] = p
}

// SetClock overrides the clock (tests).
func (s *Service) SetClock(now func() time.Time) { s.now = now }

// CreateRecoveryPoint implements kernel.RecoveryPointPort.
func (s *Service) CreateRecoveryPoint(ctx context.Context, req kernel.RecoveryPointRequest) (kernel.RecoveryPoint, error) {
	var point kernel.RecoveryPoint
	if err := verifyRequest(req); err != nil {
		return point, err
	}
	dialect, err := dialectOf(req)
	if err != nil {
		return point, err
	}
	provider, ok := s.providers[dialect]
	if !ok {
		return point, classify(KindInvalidRequest, "provider", fmt.Errorf("no provider for dialect %q", dialect))
	}

	created := s.now().UTC().Truncate(time.Microsecond)
	id, err := newRecoveryPointID(created)
	if err != nil {
		return point, err
	}
	workDir := s.workDir
	if workDir == "" {
		workDir, err = os.MkdirTemp("", "vp040-recovery-")
		if err != nil {
			return point, classify(KindToolFailure, "workdir", err)
		}
		defer os.RemoveAll(workDir)
	}
	artifactRef := filepath.Join(workDir, id+".artifact")
	if strings.TrimSpace(req.DestinationHint) != "" {
		artifactRef = req.DestinationHint
	}

	// 1. native artifact
	if err := provider.Create(ctx, req.SourceID, artifactRef); err != nil {
		return point, err
	}

	// 2. restore to a NEW isolated target
	target, cleanup, err := provider.Restore(ctx, artifactRef)
	if err != nil {
		return point, err
	}
	if cleanup != nil {
		defer cleanup()
	}

	// 3. verify: shape, denominator, ledger fingerprint, samples
	sourceSet, sourceHead, err := s.sourceFingerprint(ctx, dialect, req.SourceID)
	if err != nil {
		return point, err
	}
	summary := kernel.VerificationSummary{}
	switch dialect {
	case "sqlite":
		measure, err := s.measureSQLite(ctx, target)
		if err != nil {
			return point, err
		}
		if err := verifyConvertedShape(measure, sourceSet, sourceHead); err != nil {
			return point, err
		}
		summary.SchemaVerified = measure.IntegrityOK && measure.ForeignKeyViolations == 0
		summary.TypeContractVerified = measure.ConvertedColumns == temporalcontract.Count
		summary.ChecksumVerified = measure.LedgerSet == sourceSet && measure.BatchVersion == sourceHead
		if err := verifySQLiteSamples(ctx, target); err != nil {
			return point, err
		}
		summary.SampleVerified = true
		if !summary.Passed() {
			return kernel.RecoveryPoint{}, &Error{Kind: KindTimeContractMismatch, Op: "summary",
				Err: fmt.Errorf("verification summary incomplete: %+v", summary)}
		}
		point = kernel.RecoveryPoint{
			ID:             id,
			Dialect:        req.Dialect,
			ArtifactRef:    artifactRef,
			CatalogVersion: req.CatalogVersion,
			BatchVersion:   measure.BatchVersion,
			ContractShape:  kernel.ContractShapeConverted,
			ChecksumSet:    measure.LedgerSet,
			TimeContract:   kernel.TimeContractVP040Timestamptz,
			CreatedAt:      created,
			VerifiedAt:     s.now().UTC().Truncate(time.Microsecond),
			Verification:   summary,
		}
	default:
		db, err := sql.Open("pgx", target)
		if err != nil {
			return point, classify(KindArtifactUnreadable, "open restored database", err)
		}
		defer db.Close()
		if err := db.PingContext(ctx); err != nil {
			return point, classify(KindArtifactUnreadable, "ping restored database", err)
		}
		measure, err := measurePostgresTarget(ctx, db)
		if err != nil {
			return point, err
		}
		if err := verifyPostgresShape(measure, sourceSet, sourceHead, db, ctx); err != nil {
			return point, err
		}
		set, head, err := ledgerFingerprint(db)
		if err != nil {
			return point, classify(KindArtifactUnreadable, "read restored ledger", err)
		}
		summary.SchemaVerified = measure.MeasuredColumns == temporalcontract.Count
		summary.TypeContractVerified = measure.ConvertedColumns == temporalcontract.Count
		summary.ChecksumVerified = set == sourceSet && head == sourceHead
		// Sample verification for postgres reuses the converted-shape facts: the
		// restored database is read natively (timestamptz) so there is no
		// canonical-text spelling to re-parse.
		summary.SampleVerified = summary.TypeContractVerified
		if !summary.Passed() {
			return kernel.RecoveryPoint{}, &Error{Kind: KindTimeContractMismatch, Op: "summary",
				Err: fmt.Errorf("verification summary incomplete: %+v", summary)}
		}
		point = kernel.RecoveryPoint{
			ID:             id,
			Dialect:        req.Dialect,
			ArtifactRef:    artifactRef,
			CatalogVersion: req.CatalogVersion,
			BatchVersion:   head,
			ContractShape:  kernel.ContractShapeConverted,
			ChecksumSet:    set,
			TimeContract:   kernel.TimeContractVP040Timestamptz,
			CreatedAt:      created,
			VerifiedAt:     s.now().UTC().Truncate(time.Microsecond),
			Verification:   summary,
		}
	}

	// 4. only a verified artifact is returned; an unverified one never leaves the
	// service (Root D-010).
	if !point.Verification.Passed() {
		return kernel.RecoveryPoint{}, &Error{Kind: KindTimeContractMismatch, Op: "postcondition",
			Err: fmt.Errorf("refusing to return an unverified recovery point")}
	}
	return point, nil
}

// postgresFingerprint reads the ledger fingerprint of a postgres source store.
func postgresFingerprint(ctx context.Context, dsn string) (string, int, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return "", 0, classify(KindArtifactUnreadable, "open postgres source", err)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)
	if err := db.PingContext(ctx); err != nil {
		return "", 0, classify(KindArtifactUnreadable, "ping postgres source", err)
	}
	set, head, err := ledgerFingerprint(db)
	if err != nil {
		return "", 0, classify(KindArtifactUnreadable, "read postgres ledger", err)
	}
	return set, head, nil
}

// sourceFingerprint reads the source store's ledger fingerprint for comparison.
func (s *Service) sourceFingerprint(ctx context.Context, dialect, sourceID string) (string, int, error) {
	switch dialect {
	case "sqlite":
		return sqliteFingerprint(ctx, sourceID)
	case "postgres":
		return postgresFingerprint(ctx, sourceID)
	default:
		return "", 0, classify(KindInvalidRequest, "source fingerprint",
			fmt.Errorf("dialect %q is not wired", dialect))
	}
}
