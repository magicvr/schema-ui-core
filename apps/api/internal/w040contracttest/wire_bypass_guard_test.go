package w040contracttest

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/internal/temporalcontract"
)

// workspace-040 R3-A guard (independent audit A-002 F-I-001/F-I-002).
//
// The R3-A sweep found two public wire time fields that bypassed the shared
// fixed-6 formatter, because they were not inline Format(...) layouts:
//
//   - a *time.Time field of a struct handed straight to writeJSON
//     (PUT /api/mail/config), and
//   - a bare time.Time value placed in a response map whose key was outside the
//     hand-written scan list (GET /api/mfa/status "enrolledAt").
//
// These guards turn both classes into an executable, fail-closed check over the
// production sources, so the next such field cannot ship silently.

func wireRepoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	// test file lives at apps/api/internal/w040contracttest → root is ../../../..
	root := filepath.Clean(filepath.Join(wd, "..", "..", "..", ".."))
	if _, err := os.Stat(filepath.Join(root, "apps", "api", "go.mod")); err != nil {
		t.Fatalf("cannot locate the repository root from %s: %v", wd, err)
	}
	return root
}

// productionGoFiles lists every non-test Go source file of the API module.
func productionGoFiles(t *testing.T, root string) []string {
	t.Helper()
	apiRoot := filepath.Join(root, "apps", "api")
	var files []string
	err := filepath.Walk(apiRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if info.Name() == "testdata" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		t.Fatalf("walk %s: %v", apiRoot, err)
	}
	sort.Strings(files)
	if len(files) < 100 {
		t.Fatalf("only %d production Go files found; the walk is not seeing the module", len(files))
	}
	return files
}

func relPath(t *testing.T, root, path string) string {
	t.Helper()
	rel, err := filepath.Rel(root, path)
	if err != nil {
		t.Fatal(err)
	}
	return filepath.ToSlash(rel)
}

func readLines(t *testing.T, path string) []string {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	var lines []string
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return lines
}

// jsonTaggedTimeField matches a struct field declaration whose type is a time
// type and whose tag carries a JSON name.
var jsonTaggedTimeField = regexp.MustCompile("^\\s*[A-Z][A-Za-z0-9_]*\\s+\\*?(?:time\\.Time|sql\\.NullTime)\\s+`[^`]*json:\"[^\"]+\"[^`]*`")

// jsonTaggedTimeFieldAllowlist is the ratchet: adding a json-tagged time field to
// a production struct now requires a deliberate entry here with its reason. An
// entry does not mean "safe on the wire" — it means the file's wire path is
// projected explicitly (or the type never reaches HTTP at all).
var jsonTaggedTimeFieldAllowlist = map[string]string{
	"apps/api/internal/store/recovery.go":   "internal recovery-point marker (sidecar JSON / PG COMMENT), never an HTTP body; shape frozen by R2",
	"apps/api/internal/mail/runtime.go":     "mail.PublicView is the source of truth for the mailConfigResponse projection; the only wire path is mailConfigWire (GET+PUT covered by tests)",
	"apps/api/internal/mail/outbox.go":      "mail.OutboxRecord tags retained for the domain type; both wire paths go through outboxWireRecord",
	"apps/api/modules/wallet/voucher/voucher.go": "wallet domain model; the HTTP projection is voucherJSON (FormatWireTime)",
	"apps/api/modules/wallet/subject/subject.go": "wallet domain model; not referenced by any HTTP projection",
}

// TestNoJSONTaggedWireTimeFieldOutsideProjections is guard 1.
func TestNoJSONTaggedWireTimeFieldOutsideProjections(t *testing.T) {
	root := wireRepoRoot(t)
	var offenders []string
	seen := map[string]bool{}
	for _, path := range productionGoFiles(t, root) {
		rel := relPath(t, root, path)
		for _, line := range readLines(t, path) {
			if !jsonTaggedTimeField.MatchString(line) {
				continue
			}
			seen[rel] = true
			if _, ok := jsonTaggedTimeFieldAllowlist[rel]; !ok {
				offenders = append(offenders, rel+": "+strings.TrimSpace(line))
			}
		}
	}
	if len(offenders) > 0 {
		t.Fatalf("json-tagged time.Time field(s) outside the reviewed allowlist "+
			"(a struct handed to encoding/json renders variable-width RFC3339, not the D-003 fixed-6 shape):\n  %s",
			strings.Join(offenders, "\n  "))
	}
	// The allowlist must not rot: every entry has to still match a real field.
	for rel, reason := range jsonTaggedTimeFieldAllowlist {
		if !seen[rel] {
			t.Errorf("allowlist entry %s (%s) no longer matches any json-tagged time field; remove it", rel, reason)
		}
	}
}

// responseTimeKey matches a JSON/Go map key that names an instant. It is
// deliberately broader than the hand-written list that missed "enrolledAt":
// snake_case *_at, camelCase *At, and "timestamp".
var responseTimeKey = regexp.MustCompile("`?\"([A-Za-z_]*_at|[a-z][A-Za-z0-9]*At|timestamp)\"`?\\s*:")

// formattedValue reports whether the value expression routes the instant through
// a formatter (or is obviously not an instant at all).
func formattedValue(value string) bool {
	value = strings.TrimSpace(value)
	value = strings.TrimRight(value, ",})] \t")
	switch value {
	case "", "nil", "true", "false":
		return true
	}
	if strings.Contains(value, "Format") { // FormatWireTime / temporal.FormatWire
		return true
	}
	// String and numeric literals cannot be a time.Time value.
	if strings.HasPrefix(value, "\"") || strings.HasPrefix(value, "`") {
		return true
	}
	if value[0] >= '0' && value[0] <= '9' {
		return true
	}
	return false
}

// responseTimeKeyAllowlist holds instants that are deliberately NOT rendered by
// the public wire formatter, each with the decision that puts it outside the
// contract (Root D-009: human prose, audit detail JSON, arbitrary payload).
var responseTimeKeyAllowlist = map[string]string{
	// D-009: recycle_items.payload is arbitrary payload text with its own
	// parser/compatibility rules. The matches here are the six payload→domain
	// struct field initializers (timeField(...)), not a response map.
	"apps/api/modules/recyclebin/service.go": "D-009: recycle_items.payload reconstruction keeps its own parser rules",
}

// TestNoUnformattedResponseTimeValue is guard 2: every response map key that
// names an instant must either pass through the shared formatter or be listed
// with the decision that excludes it.
func TestNoUnformattedResponseTimeValue(t *testing.T) {
	root := wireRepoRoot(t)
	var offenders []string
	for _, path := range productionGoFiles(t, root) {
		rel := relPath(t, root, path)
		if _, ok := responseTimeKeyAllowlist[rel]; ok {
			continue
		}
		for i, line := range readLines(t, path) {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "*") {
				continue // documentation, not code
			}
			loc := responseTimeKey.FindStringSubmatchIndex(line)
			if loc == nil {
				continue
			}
			key := strings.Trim(line[loc[2]:loc[3]], "`\"")
			// Skip SQL/DDL and sort/column mapping tables: they never encode a
			// response body. Detected by the absence of a Go map/call context.
			if !strings.Contains(line, ":") || strings.Contains(line, "SELECT") ||
				strings.Contains(line, "ALTER TABLE") || strings.Contains(line, "json:\"") {
				continue
			}
			if formattedValue(line[loc[1]:]) {
				continue
			}
			offenders = append(offenders, rel+":"+itoa(i+1)+" ("+key+"): "+trimmed)
		}
	}
	if len(offenders) > 0 {
		t.Fatalf("response value(s) naming an instant bypass the shared formatter "+
			"(encoding/json renders time.Time variable-width, violating D-003):\n  %s",
			strings.Join(offenders, "\n  "))
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}

// TestUnitFamilyProvenanceIsMachineChecked answers the independent audit's
// F-I-004: temporalcontract.Column only carries Unit, so "this endpoint covers
// the nullable family" and "this endpoint covers the sentinel-converted family"
// were prose, not a checkable provenance. This test derives both from the frozen
// artifacts themselves:
//
//   - nullable: the table DDL declares the column without NOT NULL;
//   - sentinel-converted (D0): the generated VP-040 descriptor's verify entry
//     marks the column NonNull: false, i.e. the conversion makes it nullable
//     because a legacy 0 becomes NULL.
func TestUnitFamilyProvenanceIsMachineChecked(t *testing.T) {
	root := wireRepoRoot(t)

	// 1. The frozen denominator owns the unit of every column the matrix uses.
	wanted := []struct{ table, column, unit string }{
		{"digital_offers", "created_at", "sec"},
		{"digital_offers", "updated_at", "sec"},
		{"jobs", "created_at", "ms"},
		{"jobs", "updated_at", "ms"},
		{"jobs", "finished_at", "ms"},
		{"service_credentials", "revoked_at", "sec"},
		{"service_credentials", "last_used_at", "sec"},
		{"mail_config", "updated_at", "ms"},
	}
	index := map[string]string{}
	for _, c := range temporalcontract.Columns() {
		index[c.Table+"."+c.Column] = c.Unit
	}
	for _, w := range wanted {
		unit, ok := index[w.table+"."+w.column]
		if !ok {
			t.Errorf("%s.%s is not in the frozen VP-040 denominator", w.table, w.column)
			continue
		}
		if unit != w.unit {
			t.Errorf("denominator unit of %s.%s = %q, want %q", w.table, w.column, unit, w.unit)
		}
	}

	// 2. The generated descriptor marks the D0/sentinel column as becoming
	// nullable — that marker, not the unit, is what makes it the sentinel family.
	generated := readFileString(t, filepath.Join(root, "apps", "api", "modules", "corepersistence", "migration", "vp040_temporal.go"))
	const sentinelMarker = `{Table: "mail_config", Column: "updated_at", NonNull: false}`
	if !strings.Contains(generated, sentinelMarker) {
		t.Errorf("the generated VP-040 descriptor does not mark %s; the sentinel-family claim would be unverifiable", sentinelMarker)
	}

	// 3. The nullable family is nullable in the table DDL itself.
	ddl := readFileString(t, filepath.Join(root, "apps", "api", "modules", "authsession", "migration", "migration.go"))
	for _, column := range []string{"revoked_at", "last_used_at"} {
		line, ok := ddlLine(ddl, column)
		if !ok {
			t.Errorf("service_credentials DDL has no %s column", column)
			continue
		}
		if strings.Contains(line, "NOT NULL") {
			t.Errorf("service_credentials.%s is declared NOT NULL (%q); it cannot be the nullable family", column, strings.TrimSpace(line))
		}
	}
}

func readFileString(t *testing.T, path string) string {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(raw)
}

// ddlLine returns the DDL line declaring a column, ignoring index/SELECT lines.
func ddlLine(source, column string) (string, bool) {
	for _, line := range strings.Split(source, "\n") {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, column+" ") {
			continue
		}
		if strings.Contains(trimmed, "INDEX") || strings.Contains(trimmed, "SELECT") {
			continue
		}
		return line, true
	}
	return "", false
}
