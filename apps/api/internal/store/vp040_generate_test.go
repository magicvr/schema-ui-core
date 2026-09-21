package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/internal/temporalmigrate"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

// TestGenerateVP040Descriptors regenerates the workspace-040 R2 conversion
// descriptors (GOAL-003 checkpoints B/C) and their statement manifest.
//
// It is a guarded tool, not a normal test: it is skipped unless
// VP040_GENERATE=1, and when enabled it materialises the frozen v1–v72 schema,
// derives the v73–v87 conversion descriptors mechanically (live sqlite_master
// text plus the frozen per-column edits of
// r1-c2-per-column-conversion-contract-v1.0-fc.md), rewrites the per-module
// migration files, and rewrites the statement/checksum manifest attachment.
//
// Keeping it in-tree makes the derivation reproducible: re-running it on the
// same v1–v72 history must reproduce the committed descriptors byte for byte
// (any difference is transcription drift), and the manifest records the exact
// canonical statements and the real MigrationChecksum of every descriptor.
func TestGenerateVP040Descriptors(t *testing.T) {
	if os.Getenv("VP040_GENERATE") != "1" {
		t.Skip("set VP040_GENERATE=1 to regenerate the v73–v87 descriptors")
	}
	path := filepath.Join(t.TempDir(), "v72.db")
	// The conversion descriptors are written against the frozen v1–v72 shape, so
	// the derivation DB stops at v72 (the compiled catalog is now v1–v87).
	catalog := MigrationCatalog()
	if len(catalog) < 72 {
		t.Fatalf("compiled catalog has %d entries, want at least 72", len(catalog))
	}
	st, err := OpenWithCatalog(path, catalog[:72])
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if err := st.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	defer db.Close()

	live := readLiveSchema(t, db)
	children := childMap(live)

	apiRoot, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	repoRoot := filepath.Dir(filepath.Dir(apiRoot))

	var manifest strings.Builder
	manifest.WriteString(manifestHeader(live, children))

	for _, plan := range vp040Plans() {
		src, manifestSection := generateDescriptor(t, live, children, plan)
		target := filepath.Join(apiRoot, "modules", plan.dir, "vp040_temporal.go")
		if err := os.WriteFile(target, []byte(src), 0o644); err != nil {
			t.Fatalf("write %s: %v", target, err)
		}
		registerDescriptor(t, filepath.Join(apiRoot, "modules", plan.dir, "migration.go"))
		manifest.WriteString(manifestSection)
		t.Logf("wrote %s", target)
	}

	out := filepath.Join(repoRoot, "docs", "workspaces", "workspace-040-timestamptz-persistence-contract",
		"GOAL-003-r2-codec-and-descriptor-m1-m2", "attachments", "r2-v73-v87-generated-statements-v0.1.md")
	if err := os.WriteFile(out, []byte(manifest.String()), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}
	t.Logf("wrote %s", out)
}

// registerDescriptor inserts the generated descriptor into the module's
// Descriptors() catalog. The kernel orders the compiled catalog by Version, so
// the entry is inserted at the head of the returned literal with that noted.
func registerDescriptor(t *testing.T, path string) {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if strings.Contains(string(raw), "VP040TemporalDescriptor(),") {
		return
	}
	sep := "\n"
	if strings.Contains(string(raw), "\r\n") {
		sep = "\r\n"
	}
	lines := strings.Split(string(raw), sep)
	fn := -1
	for i, line := range lines {
		if strings.HasPrefix(line, "func Descriptors(") {
			fn = i
			break
		}
	}
	if fn < 0 {
		t.Fatalf("%s: func Descriptors not found", path)
	}
	anchor := -1
	for i := fn; i < len(lines); i++ {
		if strings.Contains(lines[i], "return []kernel.MigrationContribution{") {
			if anchor >= 0 {
				t.Fatalf("%s: more than one catalog literal in Descriptors", path)
			}
			anchor = i
		}
		if lines[i] == "}" && anchor >= 0 {
			break
		}
	}
	if anchor < 0 {
		t.Fatalf("%s: catalog literal not found", path)
	}
	entry := []string{
		"\t\t// workspace-040 R2 (GOAL-003 M2): v73–v87 timestamp conversions.",
		"\t\t// The kernel orders the compiled catalog by Version; source order",
		"\t\t// is not significant.",
		"\t\tVP040TemporalDescriptor(),",
	}
	out := append([]string{}, lines[:anchor+1]...)
	out = append(out, entry...)
	out = append(out, lines[anchor+1:]...)
	if err := os.WriteFile(path, []byte(strings.Join(out, sep)), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

type liveColumn struct {
	name    string
	ctype   string
	notNull bool
	dflt    string
	pk      int
}

type liveTable struct {
	name    string
	ddl     string
	columns []liveColumn
	indexes []string
}

func readLiveSchema(t *testing.T, db *sql.DB) map[string]liveTable {
	t.Helper()
	out := map[string]liveTable{}
	rows, err := db.Query(`SELECT type, name, tbl_name, sql FROM sqlite_master WHERE name NOT LIKE 'sqlite_%'`)
	if err != nil {
		t.Fatal(err)
	}
	type idx struct{ name, tbl, sql string }
	var indexes []idx
	for rows.Next() {
		var typ, name, tbl string
		var text sql.NullString
		if err := rows.Scan(&typ, &name, &tbl, &text); err != nil {
			t.Fatal(err)
		}
		switch typ {
		case "table":
			out[name] = liveTable{name: name, ddl: text.String}
		case "index":
			if text.Valid {
				indexes = append(indexes, idx{name: name, tbl: tbl, sql: text.String})
			}
		}
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	rows.Close()

	for name, table := range out {
		info, err := db.Query(`PRAGMA table_info("` + name + `")`)
		if err != nil {
			t.Fatal(err)
		}
		for info.Next() {
			var cid int
			var col liveColumn
			var dflt sql.NullString
			if err := info.Scan(&cid, &col.name, &col.ctype, &col.notNull, &dflt, &col.pk); err != nil {
				t.Fatal(err)
			}
			col.dflt = dflt.String
			table.columns = append(table.columns, col)
		}
		info.Close()
		out[name] = table
	}
	for _, i := range indexes {
		table := out[i.tbl]
		table.indexes = append(table.indexes, i.sql)
		out[i.tbl] = table
	}
	for name, table := range out {
		sort.Strings(table.indexes)
		out[name] = table
	}
	return out
}

var referenceRe = regexp.MustCompile(`(?i)REFERENCES\s+"?([A-Za-z_][A-Za-z0-9_]*)"?\s*\(`)

func childMap(live map[string]liveTable) map[string][]string {
	out := map[string][]string{}
	for name, table := range live {
		for _, match := range referenceRe.FindAllStringSubmatch(table.ddl, -1) {
			out[match[1]] = append(out[match[1]], name)
		}
	}
	for parent := range out {
		sort.Strings(out[parent])
	}
	return out
}

type conversion struct {
	column string
	unit   string // "sec" | "ms"
	flag   string // "" | "d0" | "voucher"
}

type tablePlan struct {
	table       string
	conversions []conversion
}

type descriptorPlan struct {
	version     int
	dir         string
	moduleID    string
	name        string
	transformID string
	tables      []tablePlan
	// guardTables must be absent before the descriptor may run. The frozen v73
	// scope names the retired `records` table (descriptor ledger §1).
	guardTables []string
}

func vp040Plans() []descriptorPlan {
	sec := func(cols ...string) []conversion {
		out := make([]conversion, 0, len(cols))
		for _, col := range cols {
			out = append(out, conversion{column: col, unit: "sec"})
		}
		return out
	}
	ms := func(cols ...string) []conversion {
		out := make([]conversion, 0, len(cols))
		for _, col := range cols {
			out = append(out, conversion{column: col, unit: "ms"})
		}
		return out
	}
	flagged := func(unit, flag string, cols ...string) []conversion {
		out := make([]conversion, 0, len(cols))
		for _, col := range cols {
			out = append(out, conversion{column: col, unit: unit, flag: flag})
		}
		return out
	}
	appendConv := func(groups ...[]conversion) []conversion {
		var out []conversion
		for _, group := range groups {
			out = append(out, group...)
		}
		return out
	}

	return []descriptorPlan{
		{
			version: 73, dir: "corepersistence/migration", moduleID: "core.persistence",
			name: "vp040_temporal_core_persistence", transformID: "0073:vp040-temporal-core-persistence:v1",
			guardTables: []string{"records"},
			tables: []tablePlan{
				{table: "schema_migrations", conversions: sec("applied_at")},
				{table: "mail_outbox", conversions: ms("created_at")},
				{table: "mail_config", conversions: flagged("ms", "d0", "updated_at")},
			},
		},
		{
			version: 74, dir: "authsession/migration", moduleID: "core.auth-session",
			name: "vp040_temporal_authsession", transformID: "0074:vp040-temporal-authsession:v1",
			tables: []tablePlan{
				{table: "system_data_reconcile", conversions: sec("applied_at")},
				{table: "users", conversions: appendConv(
					sec("created_at", "updated_at"),
					flagged("sec", "d0", "locked_until", "last_login_failure_at"))},
				{table: "roles", conversions: sec("created_at", "updated_at")},
				{table: "permissions", conversions: sec("created_at", "updated_at")},
				{table: "menu_items", conversions: sec("created_at", "updated_at")},
				{table: "refresh_tokens", conversions: sec("expires_at", "revoked_at", "created_at")},
				{table: "email_verification_challenges", conversions: sec("expires_at", "sent_at")},
				{table: "password_recovery_challenges", conversions: sec("expires_at", "sent_at")},
				{table: "login_failures", conversions: appendConv(
					flagged("sec", "d0", "locked_until"), sec("updated_at"))},
				{table: "user_password_history", conversions: sec("created_at")},
				{table: "user_invites", conversions: sec("expires_at", "consumed_at", "revoked_at", "last_sent_at", "created_at")},
				{table: "service_credentials", conversions: sec("expires_at", "revoked_at", "last_used_at", "created_at", "updated_at")},
			},
		},
		{
			version: 75, dir: "operationlog/migration", moduleID: "core.operationlog",
			name: "vp040_temporal_operationlog", transformID: "0075:vp040-temporal-operationlog:v1",
			tables: []tablePlan{
				{table: "operation_log", conversions: ms("created_at")},
				{table: "operation_log_archive", conversions: ms("created_at", "archived_at")},
			},
		},
		{
			version: 76, dir: "jobs/migration", moduleID: "core.jobs",
			name: "vp040_temporal_jobs", transformID: "0076:vp040-temporal-jobs:v1",
			tables: []tablePlan{
				{table: "jobs", conversions: ms("lease_expires_at", "created_at", "updated_at", "finished_at", "expires_at")},
			},
		},
		{
			version: 77, dir: "datadictionary/migration", moduleID: "admin.data-dictionary",
			name: "vp040_temporal_dictionary", transformID: "0077:vp040-temporal-dictionary:v1",
			tables: []tablePlan{
				{table: "dict_types", conversions: sec("created_at", "updated_at")},
				{table: "dict_entries", conversions: sec("created_at", "updated_at")},
			},
		},
		{
			version: 78, dir: "datapermission/migration", moduleID: "admin.data-permission",
			name: "vp040_temporal_data_permission", transformID: "0078:vp040-temporal-data-permission:v1",
			tables: []tablePlan{
				{table: "data_scope_policies", conversions: sec("updated_at")},
				{table: "user_data_scopes", conversions: sec("updated_at")},
			},
		},
		{
			version: 79, dir: "logincaptcha/migration", moduleID: "admin.login-captcha",
			name: "vp040_temporal_captcha", transformID: "0079:vp040-temporal-captcha:v1",
			tables: []tablePlan{
				{table: "captcha_challenges", conversions: sec("expires_at", "created_at")},
				{table: "captcha_config", conversions: sec("created_at", "updated_at")},
			},
		},
		{
			version: 80, dir: "mfa/migration", moduleID: "admin.mfa",
			name: "vp040_temporal_mfa", transformID: "0080:vp040-temporal-mfa:v1",
			tables: []tablePlan{
				{table: "user_mfa", conversions: sec("created_at", "updated_at")},
				{table: "mfa_proofs", conversions: sec("expires_at", "created_at")},
			},
		},
		{
			version: 81, dir: "notifications/migration", moduleID: "admin.notifications",
			name: "vp040_temporal_notifications", transformID: "0081:vp040-temporal-notifications:v1",
			tables: []tablePlan{
				{table: "notifications", conversions: sec("read_at", "created_at")},
			},
		},
		{
			version: 82, dir: "recyclebin/migration", moduleID: "admin.recycle-bin",
			name: "vp040_temporal_recycle", transformID: "0082:vp040-temporal-recycle:v1",
			tables: []tablePlan{
				{table: "recycle_items", conversions: sec("deleted_at", "restored_at")},
			},
		},
		{
			version: 83, dir: "scheduledtasks/migration", moduleID: "admin.scheduled-tasks",
			name: "vp040_temporal_scheduled_tasks", transformID: "0083:vp040-temporal-scheduled-tasks:v1",
			tables: []tablePlan{
				{table: "scheduled_tasks", conversions: sec("created_at", "updated_at")},
				{table: "task_runs", conversions: sec("started_at", "finished_at", "created_at")},
			},
		},
		{
			version: 84, dir: "settings/migration", moduleID: "admin.settings",
			name: "vp040_temporal_settings", transformID: "0084:vp040-temporal-settings:v1",
			tables: []tablePlan{
				{table: "site_settings", conversions: sec("updated_at")},
			},
		},
		{
			version: 85, dir: "wallet/migration", moduleID: "admin.wallet",
			name: "vp040_temporal_wallet", transformID: "0085:vp040-temporal-wallet:v1",
			tables: []tablePlan{
				{table: "wallet_accounts", conversions: sec("created_at", "updated_at")},
				{table: "wallet_ledger_entries", conversions: sec("created_at")},
				{table: "wallet_reconciliation_runs", conversions: sec("created_at")},
				{table: "subjects", conversions: sec("created_at")},
				{table: "vouchers", conversions: appendConv(
					flagged("sec", "voucher", "expires_at", "redeemed_at"), sec("created_at", "updated_at"))},
				{table: "voucher_batches", conversions: sec("created_at", "updated_at")},
			},
		},
		{
			version: 86, dir: "channel/telegram/migration", moduleID: "channel.telegram",
			name: "vp040_temporal_telegram", transformID: "0086:vp040-temporal-telegram:v1",
			tables: []tablePlan{
				{table: "telegram_config", conversions: flagged("sec", "d0", "updated_at")},
				{table: "telegram_sessions", conversions: sec("last_message_at", "created_at", "updated_at")},
				{table: "telegram_inbound_messages", conversions: sec("received_at")},
				{table: "telegram_outbound_messages", conversions: sec("created_at", "updated_at")},
			},
		},
		{
			version: 87, dir: "digitaloffer/migration", moduleID: "biz.digital-offer",
			name: "vp040_temporal_digital_offer", transformID: "0087:vp040-temporal-digital-offer:v1",
			tables: []tablePlan{
				{table: "digital_offers", conversions: sec("created_at", "updated_at")},
				{table: "digital_purchases", conversions: sec("created_at")},
				{table: "digital_entitlements", conversions: sec("expires_at", "created_at", "updated_at")},
			},
		},
	}
}

func generateDescriptor(t *testing.T, live map[string]liveTable, children map[string][]string, plan descriptorPlan) (string, string) {
	t.Helper()
	var preflight []string
	var rebuild []string
	var verify []string
	var pgStatements []string
	var pgVerify []string

	byColumn := func(table tablePlan, column string) (conversion, bool) {
		for _, conv := range table.conversions {
			if conv.column == column {
				return conv, true
			}
		}
		return conversion{}, false
	}

	// Child inventory (D-019 §1 / A-032 F-5 precondition): a table may be
	// rebuilt with the bare four-step only when no surviving DDL references it.
	parents := map[string]bool{}
	for _, table := range plan.tables {
		if len(children[table.table]) > 0 {
			parents[table.table] = true
		}
	}
	var detached []string
	for _, table := range plan.tables {
		if parents[table.table] {
			for _, child := range children[table.table] {
				detached = appendUnique(detached, child)
			}
		}
	}
	sort.Strings(detached)
	inPlan := map[string]bool{}
	for _, table := range plan.tables {
		inPlan[table.table] = true
	}

	buildTable := func(name string, conversions []conversion, converted bool, source string) (create string, copyStmt string) {
		table := live[name]
		ddl := table.ddl
		if converted {
			ddl = convertDDL(ddl, conversions)
			assertConvertedShape(t, ddl, table, conversions)
		}
		columns := make([]string, 0, len(table.columns))
		exprs := make([]string, 0, len(table.columns))
		for _, col := range table.columns {
			columns = append(columns, `"`+col.name+`"`)
			conv, ok := conversion{}, false
			if converted {
				conv, ok = byColumn(tablePlan{table: name, conversions: conversions}, col.name)
			}
			if !ok {
				exprs = append(exprs, `"`+col.name+`"`)
				continue
			}
			exprs = append(exprs, conversionExpr(conv, col))
		}
		copyStmt = fmt.Sprintf("INSERT INTO \"%s\" (%s)\nSELECT %s\nFROM \"%s\"",
			name, strings.Join(columns, ", "), strings.Join(exprs, ", "), source)
		return ddl, copyStmt
	}

	// Phase 1–2: detach every FK child of every rebuilt parent.
	for _, child := range detached {
		rebuild = append(rebuild, fmt.Sprintf(`CREATE TEMP TABLE "%s_bak" AS SELECT * FROM "%s"`, child, child))
	}
	for _, child := range detached {
		rebuild = append(rebuild, fmt.Sprintf(`DROP TABLE "%s"`, child))
	}

	// Phase 3: rebuild every parent that is not itself a detached child
	// (rename → create → copy → drop).
	var parentsOnly, aGroup, bare []tablePlan
	for _, table := range plan.tables {
		switch {
		case isDetached(detached, table.table):
			aGroup = append(aGroup, table)
		case parents[table.table]:
			parentsOnly = append(parentsOnly, table)
		default:
			bare = append(bare, table)
		}
	}
	renameRebuild := func(entry tablePlan) {
		ddl, copyStmt := buildTable(entry.table, entry.conversions, true, entry.table+"_old")
		rebuild = append(rebuild,
			fmt.Sprintf(`ALTER TABLE "%s" RENAME TO "%s_old"`, entry.table, entry.table),
			ddl, copyStmt,
			fmt.Sprintf(`DROP TABLE "%s_old"`, entry.table))
	}
	for _, table := range parentsOnly {
		renameRebuild(table)
	}

	// Phase 4: reattach converted children (A group), then the tables with no
	// FK children, then the verbatim children whose own conversion, if any,
	// belongs to a later descriptor.
	emitChild := func(entry tablePlan, converted bool) {
		source := entry.table + "_bak"
		ddl, copyStmt := buildTable(entry.table, entry.conversions, converted, source)
		rebuild = append(rebuild, ddl, copyStmt, fmt.Sprintf(`DROP TABLE "%s_bak"`, entry.table))
	}
	for _, table := range aGroup {
		emitChild(table, true)
	}
	for _, table := range bare {
		renameRebuild(table)
	}
	var verbatim []string
	for _, child := range detached {
		if !inPlan[child] {
			verbatim = append(verbatim, child)
		}
	}
	for _, child := range verbatim {
		emitChild(tablePlan{table: child}, false)
	}

	// Phase 5: every index last (SQLite keeps index names for the renamed
	// table, so an index created before DROP <t>_old collides).
	indexTables := map[string]bool{}
	for _, table := range plan.tables {
		indexTables[table.table] = true
	}
	for _, child := range detached {
		indexTables[child] = true
	}
	var indexNames []string
	for name := range indexTables {
		indexNames = append(indexNames, name)
	}
	sort.Strings(indexNames)
	for _, name := range indexNames {
		rebuild = append(rebuild, live[name].indexes...)
	}

	// m0 + m4 + postgres statements.
	var preflightSpecs []temporalmigrate.Preflight
	var verifySpecs []temporalmigrate.Verify
	guardSpecs := make([]temporalmigrate.Guard, 0, len(plan.guardTables))
	for _, table := range plan.guardTables {
		guardSpecs = append(guardSpecs, temporalmigrate.Guard{Table: table})
	}
	for _, table := range plan.tables {
		for _, conv := range table.conversions {
			if conv.flag == "d0" || conv.flag == "voucher" {
				preflight = append(preflight, fmt.Sprintf(
					"\t{Table: %q, Column: %q, Voucher: %v},", table.table, conv.column, conv.flag == "voucher"))
				preflightSpecs = append(preflightSpecs, temporalmigrate.Preflight{
					Table: table.table, Column: conv.column, Voucher: conv.flag == "voucher"})
			}
		}
	}
	for _, table := range plan.tables {
		var cols []string
		for _, conv := range table.conversions {
			cols = append(cols, fmt.Sprintf("%q", conv.column))
		}
		spec := temporalmigrate.Verify{Table: table.table}
		for _, conv := range table.conversions {
			spec.Columns = append(spec.Columns, conv.column)
		}
		var kids []string
		if parents[table.table] {
			for _, child := range detached {
				for _, kid := range children[table.table] {
					if kid == child {
						kids = append(kids, fmt.Sprintf("%q", child))
						spec.Children = append(spec.Children, child)
					}
				}
			}
		}
		entry := fmt.Sprintf("\t{Table: %q, Columns: []string{%s}", table.table, strings.Join(cols, ", "))
		if len(kids) > 0 {
			entry += fmt.Sprintf(", Children: []string{%s}", strings.Join(kids, ", "))
		}
		verify = append(verify, entry+"},")
		verifySpecs = append(verifySpecs, spec)
	}

	for _, table := range plan.tables {
		for _, conv := range table.conversions {
			col := liveColumn{name: conv.column}
			for _, candidate := range live[table.table].columns {
				if candidate.name == conv.column {
					col = candidate
				}
			}
			body := pgExprBody(conv, conv.column)
			switch conv.flag {
			case "d0":
				pgStatements = append(pgStatements,
					fmt.Sprintf(`ALTER TABLE "%s" ALTER COLUMN "%s" DROP DEFAULT`, table.table, conv.column),
					fmt.Sprintf(`ALTER TABLE "%s" ALTER COLUMN "%s" DROP NOT NULL`, table.table, conv.column),
					fmt.Sprintf(`ALTER TABLE "%s" ALTER COLUMN "%s" TYPE timestamptz(6) USING (CASE WHEN "%s" = 0 THEN NULL ELSE %s END)`,
						table.table, conv.column, conv.column, body))
			case "voucher":
				pgStatements = append(pgStatements,
					fmt.Sprintf(`ALTER TABLE "%s" ALTER COLUMN "%s" TYPE timestamptz(6) USING (CASE WHEN "%s" IS NULL OR "%s" = 0 THEN NULL ELSE %s END)`,
						table.table, conv.column, conv.column, conv.column, body))
			default:
				if col.notNull {
					pgStatements = append(pgStatements, fmt.Sprintf(
						`ALTER TABLE "%s" ALTER COLUMN "%s" TYPE timestamptz(6) USING (%s)`, table.table, conv.column, body))
				} else {
					pgStatements = append(pgStatements, fmt.Sprintf(
						`ALTER TABLE "%s" ALTER COLUMN "%s" TYPE timestamptz(6) USING (CASE WHEN "%s" IS NULL THEN NULL ELSE %s END)`,
						table.table, conv.column, conv.column, body))
				}
			}
			nullable := conv.flag == "d0" || conv.flag == "voucher" || !col.notNull
			pgVerify = append(pgVerify, fmt.Sprintf("\t{Table: %q, Column: %q, NonNull: %t},",
				table.table, conv.column, !nullable))
		}
	}

	v := plan.version
	prefix := fmt.Sprintf("vp040V%d", v)
	var src strings.Builder
	src.WriteString(fileHeader(plan))
	fmt.Fprintf(&src, "const (\n\t%sName        = %q\n\t%sTransformID = %q\n)\n\n", prefix, plan.name, prefix, plan.transformID)
	preflightBody := "{}"
	if len(preflight) > 0 {
		preflightBody = "{\n" + strings.Join(preflight, "\n") + "\n}"
	}
	guardBody := "{}"
	if len(plan.guardTables) > 0 {
		lines := make([]string, 0, len(plan.guardTables))
		for _, table := range plan.guardTables {
			lines = append(lines, fmt.Sprintf("\t{Table: %q},", table))
		}
		guardBody = "{\n" + strings.Join(lines, "\n") + "\n}"
	}
	fmt.Fprintf(&src, "// %sGuards are the m0 preconditions on tables retired by earlier history\n// (descriptor ledger §1: v73 asserts the retired `records` table is absent).\nvar %sGuards = []temporalmigrate.Guard%s\n\n",
		prefix, prefix, guardBody)
	fmt.Fprintf(&src, "// %sPreflight is the m0 sentinel census (Root D-012 / D-015 policy).\nvar %sPreflight = []temporalmigrate.Preflight%s\n\n",
		prefix, prefix, preflightBody)
	fmt.Fprintf(&src, "// %sVerify is the m4 post-rebuild assertion set.\nvar %sVerify = []temporalmigrate.Verify{\n%s\n}\n\n",
		prefix, prefix, strings.Join(verify, "\n"))
	fmt.Fprintf(&src, "// %sRebuild is the ordered m1–m3 slice: constraint handling, table\n// rebuild with the frozen conversion expressions, then every index.\nvar %sRebuild = []string{\n%s\n}\n\n",
		prefix, prefix, quoteStatements(rebuild))
	fmt.Fprintf(&src, "// %sStatements returns the canonical checksum input (D-017: m0 → m1–m3 → m4).\nfunc %sStatements() []string {\n\treturn temporalmigrate.OrderedWithGuards(%sGuards, %sPreflight, %sRebuild, %sVerify)\n}\n\n",
		prefix, prefix, prefix, prefix, prefix, prefix)
	fmt.Fprintf(&src, "// apply%s is the SQLite Apply body.\nfunc apply%s(tx kernel.Tx) error {\n", prefix, prefix)
	fmt.Fprintf(&src, "\tif err := temporalmigrate.RunGuards(tx, %sGuards, %q); err != nil {\n\t\treturn err\n\t}\n", prefix, fmt.Sprintf("vp040 v%d sqlite guards", v))
	fmt.Fprintf(&src, "\tif err := temporalmigrate.RunPreflight(tx, %sPreflight); err != nil {\n\t\treturn err\n\t}\n", prefix)
	fmt.Fprintf(&src, "\tif err := temporalmigrate.Exec(tx, %sRebuild, %q); err != nil {\n\t\treturn err\n\t}\n", prefix, fmt.Sprintf("vp040 v%d sqlite rebuild", v))
	fmt.Fprintf(&src, "\treturn temporalmigrate.RunVerify(tx, %sVerify, %q)\n}\n\n", prefix, fmt.Sprintf("vp040 v%d sqlite verify", v))
	fmt.Fprintf(&src, "// %sPostgres is the explicit postgres m1–m2 slice (D-019 §6: no regex\n// derivation; the millisecond family uses the integer-split expression).\nvar %sPostgres = []string{\n%s\n}\n\n",
		prefix, prefix, quoteStatements(pgStatements))
	fmt.Fprintf(&src, "// %sPostgresVerify is the postgres m4 type/precision assertion set.\nvar %sPostgresVerify = []temporalmigrate.PgVerify{\n%s\n}\n\n",
		prefix, prefix, strings.Join(pgVerify, "\n"))
	fmt.Fprintf(&src, "// apply%sPostgres is the postgres Apply body.\nfunc apply%sPostgres(tx kernel.Tx) error {\n", prefix, prefix)
	fmt.Fprintf(&src, "\tif err := temporalmigrate.RunPostgresGuards(tx, %sGuards, %q); err != nil {\n\t\treturn err\n\t}\n", prefix, fmt.Sprintf("vp040 v%d postgres guards", v))
	fmt.Fprintf(&src, "\tif err := temporalmigrate.RunPreflight(tx, %sPreflight); err != nil {\n\t\treturn err\n\t}\n", prefix)
	fmt.Fprintf(&src, "\tif err := temporalmigrate.Exec(tx, %sPostgres, %q); err != nil {\n\t\treturn err\n\t}\n", prefix, fmt.Sprintf("vp040 v%d postgres convert", v))
	fmt.Fprintf(&src, "\treturn temporalmigrate.RunPostgresVerify(tx, %sPostgresVerify, %q)\n}\n\n", prefix, fmt.Sprintf("vp040 v%d postgres verify", v))
	src.WriteString("// VP040TemporalDescriptor returns the workspace-040 R2 conversion descriptor\n// (D-014 allocation, D-017 checksum convention).\nfunc VP040TemporalDescriptor() kernel.MigrationContribution {\n\treturn kernel.MigrationContribution{\n")
	fmt.Fprintf(&src, "\t\tContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: %sName},\n", prefix)
	fmt.Fprintf(&src, "\t\tVersion:              %d,\n\t\tName:                 %sName,\n", v, prefix)
	fmt.Fprintf(&src, "\t\tChecksum:             kernel.MigrationChecksum(%sStatements(), %sTransformID),\n", prefix, prefix)
	fmt.Fprintf(&src, "\t\tApply:                apply%s,\n\t\tApplyPostgres:        apply%sPostgres,\n\t}\n}\n", prefix, prefix)

	var section strings.Builder
	fmt.Fprintf(&section, "\n## v%d · `%s` · `%s`\n\n", v, plan.moduleID, plan.name)
	fmt.Fprintf(&section, "- `transform_id`: `%s`\n", plan.transformID)
	fmt.Fprintf(&section, "- **`MigrationChecksum`（真实值，D-017 单 checksum / SQLite 切片）**：`%s`\n",
		kernel.MigrationChecksum(temporalmigrate.OrderedWithGuards(guardSpecs, preflightSpecs, rebuild, verifySpecs), plan.transformID))
	fmt.Fprintf(&section, "- 表范围（descriptor 顺序）：%s\n", backtickList(tableNames(plan)))
	fmt.Fprintf(&section, "- FK 子女盘点（机械解析 v72 DDL 的 `REFERENCES`）：父表 %s；摘除并建回的子表 %s\n",
		backtickList(mapKeys(parents)), backtickList(detached))
	fmt.Fprintf(&section, "- m0 预检列：%s；m0 retired-table guard：%s\n", backtickList(preflightColumns(plan)), backtickList(plan.guardTables))
	fmt.Fprintf(&section, "- PG 语句数 %d；SQLite 语句数（m1–m3）%d；m4 断言 %d\n\n", len(pgStatements), len(rebuild), len(verify))
	section.WriteString("canonical（m1–m3）：\n\n```sql\n")
	for _, stmt := range rebuild {
		section.WriteString(stmt)
		section.WriteString(";\n")
	}
	section.WriteString("```\n\nm0（guard/预检）/ m4（校验）语句文本：\n\n```sql\n")
	for _, stmt := range append(append(append([]string{}, temporalmigrate.GuardStatements(guardSpecs)...), preflight...), verify...) {
		section.WriteString(stmt)
		section.WriteString(";\n")
	}
	section.WriteString("```\n\n<details><summary>postgres m1–m2</summary>\n\n```sql\n")
	for _, stmt := range pgStatements {
		section.WriteString(stmt)
		section.WriteString(";\n")
	}
	section.WriteString("```\n\n</details>\n")

	return src.String(), section.String()
}

func conversionExpr(conv conversion, col liveColumn) string {
	body := sqliteExprBody(conv, col.name)
	switch conv.flag {
	case "d0":
		return fmt.Sprintf("CASE WHEN %s = 0 THEN NULL ELSE %s END", col.name, body)
	case "voucher":
		return fmt.Sprintf("CASE WHEN %s IS NULL OR %s = 0 THEN NULL ELSE %s END", col.name, col.name, body)
	default:
		if col.notNull {
			return body
		}
		return fmt.Sprintf("CASE WHEN %s IS NULL THEN NULL ELSE %s END", col.name, body)
	}
}

func sqliteExprBody(conv conversion, column string) string {
	if conv.unit == "ms" {
		return fmt.Sprintf("strftime('%%Y-%%m-%%dT%%H:%%M:%%S', CASE WHEN %s >= 0 THEN %s/1000 ELSE (%s-999)/1000 END, 'unixepoch') || '.' || printf('%%03d', (%s%%1000 + 1000) %% 1000) || '000Z'",
			column, column, column, column)
	}
	return fmt.Sprintf("strftime('%%Y-%%m-%%dT%%H:%%M:%%S', %s, 'unixepoch') || '.000000Z'", column)
}

func pgExprBody(conv conversion, column string) string {
	if conv.unit == "ms" {
		return fmt.Sprintf("TIMESTAMPTZ 'epoch' + ((\"%s\" - CASE WHEN \"%s\" >= 0 THEN 0 ELSE 999 END) / 1000) * INTERVAL '1 second' + (((\"%s\" %% 1000) + 1000) %% 1000) * INTERVAL '1 millisecond'",
			column, column, column)
	}
	return fmt.Sprintf("date_trunc('microseconds', to_timestamp(\"%s\"::double precision))", column)
}

func convertDDL(ddl string, conversions []conversion) string {
	for _, conv := range conversions {
		pattern := regexp.MustCompile(`(?m)(^|,)(\s*` + regexp.QuoteMeta(conv.column) + `\s+)INTEGER\b(\s+NOT NULL\b)?(\s+DEFAULT 0\b)?`)
		// NOT NULL is preserved for the plain columns and dropped for the D0 /
		// voucher columns; DEFAULT 0 is only ever present on the D0 columns and
		// is never re-emitted.
		replacement := "${1}${2}TEXT${3}"
		if conv.flag == "d0" || conv.flag == "voucher" {
			replacement = "${1}${2}TEXT"
		}
		out := pattern.ReplaceAllString(ddl, replacement)
		if out == ddl {
			panic("no type edit applied for column " + conv.column)
		}
		if !regexp.MustCompile(`\b` + regexp.QuoteMeta(conv.column) + `\s+TEXT\b`).MatchString(out) {
			panic("type edit produced no TEXT column for " + conv.column)
		}
		ddl = out
	}
	return ddl
}

// assertConvertedShape builds the converted DDL in a scratch database and
// compares the resulting physical shape with the frozen expectation: same
// column names in the same cid order, TEXT for every converted column, and
// NOT NULL dropped for the D0 / voucher columns.
func assertConvertedShape(t *testing.T, ddl string, table liveTable, conversions []conversion) {
	t.Helper()
	scratch, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer scratch.Close()
	scratch.SetMaxOpenConns(1)
	if _, err := scratch.Exec(ddl); err != nil {
		t.Fatalf("converted DDL for %s does not execute: %v\n%s", table.name, err, ddl)
	}
	info, err := scratch.Query(`PRAGMA table_info("` + table.name + `")`)
	if err != nil {
		t.Fatal(err)
	}
	defer info.Close()
	var got []liveColumn
	for info.Next() {
		var cid int
		var col liveColumn
		var dflt sql.NullString
		if err := info.Scan(&cid, &col.name, &col.ctype, &col.notNull, &dflt, &col.pk); err != nil {
			t.Fatal(err)
		}
		col.dflt = dflt.String
		got = append(got, col)
	}
	if len(got) != len(table.columns) {
		t.Fatalf("%s: converted DDL has %d columns, want %d", table.name, len(got), len(table.columns))
	}
	for i, want := range table.columns {
		if got[i].name != want.name {
			t.Fatalf("%s: column %d = %q, want %q", table.name, i, got[i].name, want.name)
		}
		conv, converted := conversionOf(conversions, want.name)
		wantType := strings.ToUpper(want.ctype)
		wantNotNull := want.notNull
		if converted {
			wantType = "TEXT"
			if conv.flag == "d0" || conv.flag == "voucher" {
				wantNotNull = false
			}
		}
		if strings.ToUpper(got[i].ctype) != wantType {
			t.Fatalf("%s.%s: type = %q, want %q", table.name, want.name, got[i].ctype, wantType)
		}
		if got[i].notNull != wantNotNull {
			t.Fatalf("%s.%s: notnull = %v, want %v", table.name, want.name, got[i].notNull, wantNotNull)
		}
	}
}

func conversionOf(conversions []conversion, column string) (conversion, bool) {
	for _, conv := range conversions {
		if conv.column == column {
			return conv, true
		}
	}
	return conversion{}, false
}

func appendUnique(list []string, value string) []string {
	for _, item := range list {
		if item == value {
			return list
		}
	}
	return append(list, value)
}

func isDetached(detached []string, table string) bool {
	for _, item := range detached {
		if item == table {
			return true
		}
	}
	return false
}

func tableNames(plan descriptorPlan) []string {
	out := make([]string, 0, len(plan.tables))
	for _, table := range plan.tables {
		out = append(out, table.table)
	}
	return out
}

func preflightColumns(plan descriptorPlan) []string {
	var out []string
	for _, table := range plan.tables {
		for _, conv := range table.conversions {
			if conv.flag == "d0" || conv.flag == "voucher" {
				out = append(out, table.table+"."+conv.column)
			}
		}
	}
	return out
}

func mapKeys(in map[string]bool) []string {
	out := make([]string, 0, len(in))
	for key := range in {
		out = append(out, key)
	}
	sort.Strings(out)
	return out
}

func backtickList(items []string) string {
	if len(items) == 0 {
		return "（无）"
	}
	return "`" + strings.Join(items, "`, `") + "`"
}

func quoteStatements(stmts []string) string {
	var out strings.Builder
	for _, stmt := range stmts {
		out.WriteString("\t`")
		out.WriteString(stmt)
		out.WriteString("`,\n")
	}
	return strings.TrimRight(out.String(), "\n")
}

func fileHeader(plan descriptorPlan) string {
	return fmt.Sprintf(`// Code generated for workspace-040 R2 (GOAL-003 checkpoints B/C) from the
// v72-applied schema. DO NOT EDIT BY HAND.
//
// Every CREATE TABLE body below is the live sqlite_master text of the frozen
// v1–v72 history (live PRAGMA table_info column order) with only the frozen
// per-column edits of r1-c2-per-column-conversion-contract-v1.0-fc.md applied:
// the converted column becomes TEXT and, for the D0 / voucher columns, loses
// NOT NULL and DEFAULT 0. No other token changed.
//
// Checksum input (D-017): m0 preflight → m1–m3 rebuild → m4 verification. The
// postgres variant is explicit DDL and is not hashed (D-019 §6).
//
// ModuleID %s · version %d · D-014 allocation / D-017 transform id.
package migration

import (
	"github.com/magicvr/schema-ui-core/apps/api/internal/temporalmigrate"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

`, plan.moduleID, plan.version)
}

func manifestHeader(live map[string]liveTable, children map[string][]string) string {
	var out strings.Builder
	out.WriteString("---\n" +
		"id: r2-v73-v87-generated-statements-v0.1\n" +
		"doc_type: evidence-attachment\n" +
		"title: R2 v73–v87 conversion descriptor 语句清单 v0.1（落码产物）\n" +
		"status: draft\n" +
		"created: 2026-09-20\n" +
		"updated: 2026-09-20\n" +
		"parent: GOAL-003-r2-codec-and-descriptor-m1-m2\n" +
		"version: 0.1.0\n" +
		"---\n\n" +
		"# R2 v73–v87 conversion descriptor 语句清单 v0.1\n\n" +
		"> 本文件由生成器从 **v72 已 apply 库的 live `sqlite_master`** 机械导出（GOAL-003\n" +
		"> checkpoint B/C 落码产物），记录每个 descriptor 的 canonical 语句顺序（m0 预检 →\n" +
		"> m1–m3 重建 → m4 校验）与显式 PG DDL。**m1–m3 的 `CREATE TABLE` 正文 = live DDL\n" +
		"> 逐字 + 冻结的逐列类型编辑**（`INTEGER` → `TEXT`；D0/voucher 列去 `NOT NULL` /\n" +
		"> `DEFAULT 0`），其余 token 未改。\n" +
		">\n" +
		"> 缺省约定：`INSERT … SELECT` 的列清单按 live `PRAGMA table_info` 的 cid 序展开；\n" +
		"> 标识符用双引号；全部 `CREATE INDEX` 在 `DROP TABLE <t>_old` / `<t>_bak` 之后。\n\n" +
		"## 0. FK 子表盘点（F-5 先决条件，A-032 §G / 逐表附件 §5.8）\n\n" +
		"机械解析 v72 DDL 的全部 `REFERENCES <parent>(` 后，**存在 FK 子表的父表**为：\n\n")
	keys := make([]string, 0, len(children))
	for key := range children {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		fmt.Fprintf(&out, "- `%s` ← %s\n", key, backtickList(children[key]))
	}
	var childless []string
	for name := range live {
		if len(children[name]) == 0 {
			childless = append(childless, name)
		}
	}
	sort.Strings(childless)
	fmt.Fprintf(&out, "\n**无任何存留 DDL 引用（= 裸四步安全）的表（%d 张）**：%s\n", len(childless), backtickList(childless))
	return out.String()
}
