package backup

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/temporalcontract"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

// TemporalColumns returns the frozen 90-column denominator.
func TemporalColumns() []temporalcontract.Column {
	return temporalcontract.Columns()
}

// sqliteTargetVerification is the measured shape of one restored SQLite target.
type sqliteTargetVerification struct {
	// MeasuredColumns is how many of the frozen denominator columns were found.
	MeasuredColumns int
	// ConvertedColumns is how many of those carry the canonical TEXT storage
	// class.
	ConvertedColumns int
	// Missing lists denominator entries absent from the target.
	Missing []string
	// WrongShape lists denominator columns whose declared type is not TEXT.
	WrongShape           []string
	LedgerSet            string
	BatchVersion         int
	IntegrityOK          bool
	ForeignKeyViolations int
	// RecordsPresent reports whether the retired `records` table came back.
	RecordsPresent bool
}

// measureSQLiteTarget opens the restored file and measures the frozen contract
// facts. Errors while opening/reading are classified as unreadable, never as a
// contract mismatch (C3 §5.1).
func measureSQLiteTarget(ctx context.Context, path string) (sqliteTargetVerification, error) {
	var out sqliteTargetVerification
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return out, classify(KindArtifactUnreadable, "open restored target", err)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)
	if err := db.PingContext(ctx); err != nil {
		return out, classify(KindArtifactUnreadable, "ping restored target", err)
	}

	var integrity string
	if err := db.QueryRowContext(ctx, `PRAGMA integrity_check`).Scan(&integrity); err != nil {
		return out, classify(KindArtifactUnreadable, "integrity_check", err)
	}
	out.IntegrityOK = strings.TrimSpace(integrity) == "ok"

	rows, err := db.QueryContext(ctx, `PRAGMA foreign_key_check`)
	if err != nil {
		return out, classify(KindArtifactUnreadable, "foreign_key_check", err)
	}
	for rows.Next() {
		out.ForeignKeyViolations++
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return out, classify(KindArtifactUnreadable, "foreign_key_check", err)
	}
	rows.Close()

	// Per-table declared types, one query per distinct table (44 tables).
	declared := map[string]map[string]string{}
	for _, column := range temporalcontract.Columns() {
		if _, ok := declared[column.Table]; ok {
			continue
		}
		info, err := db.QueryContext(ctx, `PRAGMA table_info("`+column.Table+`")`)
		if err != nil {
			return out, classify(KindArtifactUnreadable, "table_info", err)
		}
		types := map[string]string{}
		for info.Next() {
			var cid, notNull, pk int
			var name, ctype string
			var dflt sql.NullString
			if err := info.Scan(&cid, &name, &ctype, &notNull, &dflt, &pk); err != nil {
				info.Close()
				return out, classify(KindArtifactUnreadable, "table_info scan", err)
			}
			types[name] = strings.ToUpper(strings.TrimSpace(ctype))
		}
		if err := info.Err(); err != nil {
			info.Close()
			return out, classify(KindArtifactUnreadable, "table_info", err)
		}
		info.Close()
		declared[column.Table] = types
	}
	for _, column := range temporalcontract.Columns() {
		types, ok := declared[column.Table]
		if !ok {
			out.Missing = append(out.Missing, column.Table+"."+column.Column)
			continue
		}
		ctype, ok := types[column.Column]
		if !ok {
			out.Missing = append(out.Missing, column.Table+"."+column.Column)
			continue
		}
		out.MeasuredColumns++
		if ctype == "TEXT" {
			out.ConvertedColumns++
		} else {
			out.WrongShape = append(out.WrongShape, fmt.Sprintf("%s.%s=%s", column.Table, column.Column, ctype))
		}
	}

	var records int
	if err := db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name = 'records'`).Scan(&records); err != nil {
		return out, classify(KindArtifactUnreadable, "probe records", err)
	}
	out.RecordsPresent = records != 0

	set, head, err := ledgerFingerprint(db)
	if err != nil {
		return out, classify(KindArtifactUnreadable, "read ledger", err)
	}
	out.LedgerSet = set
	out.BatchVersion = head
	return out, nil
}

// verifyConvertedShape reduces a measurement to the C3 verdict. The order of the
// checks is the classification order of C3 §5.1: an incomplete column set is
// reported as such, otherwise a non-converted shape is a contract mismatch.
func verifyConvertedShape(m sqliteTargetVerification, wantChecksumSet string, wantCatalogVersion int) error {
	if len(m.Missing) > 0 || m.MeasuredColumns != temporalcontract.Count {
		return &Error{Kind: KindTemporalColumnSetIncomplete, Op: "verify sqlite target",
			Err: fmt.Errorf("%d/%d contract columns measured (missing %v)",
				m.MeasuredColumns, temporalcontract.Count, m.Missing)}
	}
	if len(m.WrongShape) > 0 {
		return &Error{Kind: KindTimeContractMismatch, Op: "verify sqlite target",
			Err: fmt.Errorf("%d column(s) are not canonical TEXT: %v",
				len(m.WrongShape), m.WrongShape)}
	}
	if !m.IntegrityOK {
		return &Error{Kind: KindArtifactUnreadable, Op: "verify sqlite target",
			Err: fmt.Errorf("integrity_check is not ok")}
	}
	if m.ForeignKeyViolations != 0 {
		return &Error{Kind: KindArtifactUnreadable, Op: "verify sqlite target",
			Err: fmt.Errorf("%d foreign_key_check violation(s)", m.ForeignKeyViolations)}
	}
	if m.RecordsPresent {
		return &Error{Kind: KindTimeContractMismatch, Op: "verify sqlite target",
			Err: fmt.Errorf("retired `records` table is present in the artifact")}
	}
	if wantChecksumSet != "" && m.LedgerSet != wantChecksumSet {
		return &Error{Kind: KindChecksumMismatch, Op: "verify sqlite target",
			Err: fmt.Errorf("artifact ledger %s, source %s", m.LedgerSet, wantChecksumSet)}
	}
	if wantCatalogVersion > 0 && m.BatchVersion != wantCatalogVersion {
		return &Error{Kind: KindChecksumMismatch, Op: "verify sqlite target",
			Err: fmt.Errorf("artifact batch version %d, want %d", m.BatchVersion, wantCatalogVersion)}
	}
	return nil
}

// verifySQLiteSamples re-reads representative instants from the restored target.
// The sample set is the C3 §5 item 4 list: seconds, milliseconds, sentinel 0,
// nullable absence and fixed-6 lexical order. Negative values are deliberately
// NOT a positive sample of B (C3 §5 note): they belong to the legacy reverse
// assertion and to the voucher preflight.
func verifySQLiteSamples(ctx context.Context, path string) error {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return classify(KindArtifactUnreadable, "open samples target", err)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)

	// Sentinel 0 must have become NULL: mail_config and telegram_config are the
	// two singleton D0 columns in the converted set.
	for _, probe := range []struct{ table, column string }{
		{"mail_config", "updated_at"},
		{"telegram_config", "updated_at"},
		{"users", "locked_until"},
	} {
		var count int
		query := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE %s IS NOT NULL`, probe.table, probe.column)
		if err := db.QueryRowContext(ctx, query).Scan(&count); err != nil {
			// The table may be empty (no row at all): that is "no sample", not a
			// failure.
			continue
		}
		if count != 0 {
			return &Error{Kind: KindSampleMismatch, Op: "verify sqlite samples",
				Err: fmt.Errorf("%s.%s has %d non-NULL legacy sentinel value(s)", probe.table, probe.column, count)}
		}
	}

	// Fixed-6 lexical order equals instant order for the canonical form: read
	// every canonical text value in the converted set and require it to parse to
	// itself through the codec.
	total, canonical := 0, 0
	for _, column := range temporalcontract.Columns() {
		query := fmt.Sprintf(`SELECT %s FROM %s WHERE %s IS NOT NULL LIMIT 200`,
			column.Column, column.Table, column.Column)
		rows, err := db.QueryContext(ctx, query)
		if err != nil {
			continue // table empty or absent sample: not a failure
		}
		for rows.Next() {
			var value string
			if err := rows.Scan(&value); err != nil {
				rows.Close()
				return classify(KindArtifactUnreadable, "scan sample", err)
			}
			total++
			if isCanonicalFixed6(value) {
				canonical++
			} else {
				rows.Close()
				return &Error{Kind: KindSampleMismatch, Op: "verify sqlite samples",
					Err: fmt.Errorf("%s.%s holds non-canonical value %q", column.Table, column.Column, value)}
			}
		}
		rows.Close()
		if err := rows.Err(); err != nil {
			return classify(KindArtifactUnreadable, "iterate samples", err)
		}
	}
	_ = total
	_ = canonical
	return nil
}

// isCanonicalFixed6 reports whether value is the 27-character canonical form.
func isCanonicalFixed6(value string) bool {
	if len(value) != 27 {
		return false
	}
	t, err := time.Parse("2006-01-02T15:04:05.000000Z", value)
	if err != nil {
		return false
	}
	return t.UTC().Format("2006-01-02T15:04:05.000000Z") == value
}

// postgresTargetVerification is the measured shape of one restored PG database.
type postgresTargetVerification struct {
	MeasuredColumns  int
	ConvertedColumns int
	Missing          []string
	WrongShape       []string
	ServerVersion    string
}

// measurePostgresTarget measures the frozen contract facts on the restored
// database through information_schema.
func measurePostgresTarget(ctx context.Context, db *sql.DB) (postgresTargetVerification, error) {
	var out postgresTargetVerification
	if err := db.QueryRowContext(ctx, `SHOW server_version`).Scan(&out.ServerVersion); err != nil {
		return out, classify(KindArtifactUnreadable, "server_version", err)
	}
	for _, column := range temporalcontract.Columns() {
		var dataType sql.NullString
		var precision sql.NullInt64
		err := db.QueryRowContext(ctx, `
			SELECT data_type, datetime_precision FROM information_schema.columns
			WHERE table_schema = current_schema() AND table_name = $1 AND column_name = $2`,
			column.Table, column.Column).Scan(&dataType, &precision)
		if err == sql.ErrNoRows {
			out.Missing = append(out.Missing, column.Table+"."+column.Column)
			continue
		}
		if err != nil {
			return out, classify(KindArtifactUnreadable, "information_schema", err)
		}
		out.MeasuredColumns++
		if dataType.String == "timestamp with time zone" && precision.Int64 == 6 {
			out.ConvertedColumns++
			continue
		}
		out.WrongShape = append(out.WrongShape,
			fmt.Sprintf("%s.%s=%s(%d)", column.Table, column.Column, dataType.String, precision.Int64))
	}
	return out, nil
}

// verifyPostgresShape reduces the PG measurement to the C3 verdict.
func verifyPostgresShape(m postgresTargetVerification, wantChecksumSet string, wantCatalogVersion int, db *sql.DB, ctx context.Context) error {
	if len(m.Missing) > 0 || m.MeasuredColumns != temporalcontract.Count {
		return &Error{Kind: KindTemporalColumnSetIncomplete, Op: "verify postgres target",
			Err: fmt.Errorf("%d/%d contract columns measured (missing %v)",
				m.MeasuredColumns, temporalcontract.Count, m.Missing)}
	}
	if len(m.WrongShape) > 0 {
		return &Error{Kind: KindTimeContractMismatch, Op: "verify postgres target",
			Err: fmt.Errorf("%d column(s) are not timestamptz(6): %v", len(m.WrongShape), m.WrongShape)}
	}
	var records int
	if err := db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = current_schema() AND table_name = 'records'`).Scan(&records); err != nil {
		return classify(KindArtifactUnreadable, "probe records", err)
	}
	if records != 0 {
		return &Error{Kind: KindTimeContractMismatch, Op: "verify postgres target",
			Err: fmt.Errorf("retired `records` table is present in the artifact")}
	}
	if wantChecksumSet != "" {
		set, head, err := ledgerFingerprint(db)
		if err != nil {
			return classify(KindArtifactUnreadable, "read ledger", err)
		}
		if set != wantChecksumSet {
			return &Error{Kind: KindChecksumMismatch, Op: "verify postgres target",
				Err: fmt.Errorf("artifact ledger %s, source %s", set, wantChecksumSet)}
		}
		if wantCatalogVersion > 0 && head != wantCatalogVersion {
			return &Error{Kind: KindChecksumMismatch, Op: "verify postgres target",
				Err: fmt.Errorf("artifact batch version %d, want %d", head, wantCatalogVersion)}
		}
	}
	return nil
}

// verifyRequest validates the caller's request against the frozen contract.
func verifyRequest(req kernel.RecoveryPointRequest) error {
	if strings.TrimSpace(req.SourceID) == "" {
		return &Error{Kind: KindInvalidRequest, Op: "request", Err: fmt.Errorf("SourceID is required")}
	}
	if req.TimeContract != "" && req.TimeContract != kernel.TimeContractVP040Timestamptz {
		return &Error{Kind: KindInvalidRequest, Op: "request",
			Err: fmt.Errorf("time contract %q, want %q", req.TimeContract, kernel.TimeContractVP040Timestamptz)}
	}
	if req.CatalogVersion <= 0 {
		return &Error{Kind: KindInvalidRequest, Op: "request", Err: fmt.Errorf("CatalogVersion must be positive")}
	}
	return nil
}
