package pgsetup

import (
	"context"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/pgtest"
)

func TestValidName(t *testing.T) {
	valid := []string{"schema_ui_dev", "schema_ui_test", "a", "_x", "db1"}
	for _, name := range valid {
		if !ValidName(name) {
			t.Errorf("ValidName(%q) = false, want true", name)
		}
	}
	// Everything that could escape an unquoted identifier is rejected: this is
	// the only guard between a configured name and CREATE/DROP DATABASE.
	invalid := []string{"", "Upper", "1leading", "has-dash", "has space", "a;b", `a"b`, "a'b", "a\nb"}
	for _, name := range invalid {
		if ValidName(name) {
			t.Errorf("ValidName(%q) = true, want false", name)
		}
	}
}

// TestMaintenanceCandidatesBootstrapOnEmptyServer pins the W35 defect: the
// maintenance connection must not depend on the configured database existing, or
// the tool cannot create the first database on a fresh/reset server.
func TestMaintenanceCandidatesBootstrapOnEmptyServer(t *testing.T) {
	server := Server{Role: RoleDev, Database: "schema_ui_dev"}
	got := server.MaintenanceCandidates()
	want := []string{"schema_ui_dev", "postgres", "template1"}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("MaintenanceCandidates = %v, want %v", got, want)
	}

	// A configured name that is already a maintenance database must not repeat.
	server.Database = "postgres"
	got = server.MaintenanceCandidates()
	if strings.Join(got, ",") != "postgres,template1" {
		t.Fatalf("MaintenanceCandidates(postgres) = %v, want [postgres template1]", got)
	}

	// An empty configured name still yields the always-present candidates.
	server.Database = ""
	got = server.MaintenanceCandidates()
	if strings.Join(got, ",") != "postgres,template1" {
		t.Fatalf("MaintenanceCandidates(\"\") = %v, want [postgres template1]", got)
	}
	if len(got) == 0 {
		t.Fatal("MaintenanceCandidates returned nothing")
	}
}

// TestServerForRoleUsesTheMatchingNamespace pins the per-role credential split:
// the dev database is provisioned with DB_* (the API server's own namespace) and
// the test database with PG_TEST_* (internal/pgtest's namespace). A single
// namespace would provision the wrong server on a machine that separates them.
func TestServerForRoleUsesTheMatchingNamespace(t *testing.T) {
	t.Setenv("DB_HOST", "dev.example")
	t.Setenv("DB_USER", "devuser")
	t.Setenv("DB_PASSWORD", "devpass")
	t.Setenv("DB_NAME", "dev_db")
	t.Setenv("PG_TEST_HOST", "test.example")
	t.Setenv("PG_TEST_USER", "testuser")
	t.Setenv("PG_TEST_PASSWORD", "testpass")
	t.Setenv("PG_TEST_DB", "test_db")

	dev := ServerForRole(RoleDev)
	if dev.Host != "dev.example" || dev.User != "devuser" || dev.Password != "devpass" || dev.Database != "dev_db" {
		t.Fatalf("dev server = %+v, want the DB_* namespace", dev)
	}
	test := ServerForRole(RoleTest)
	if test.Host != "test.example" || test.User != "testuser" || test.Password != "testpass" || test.Database != "test_db" {
		t.Fatalf("test server = %+v, want the PG_TEST_* namespace", test)
	}
	if DatabaseName(RoleDev) != "dev_db" || DatabaseName(RoleTest) != "test_db" {
		t.Fatalf("DatabaseName = %q/%q, want dev_db/test_db", DatabaseName(RoleDev), DatabaseName(RoleTest))
	}
}

func TestServerForRoleDefaults(t *testing.T) {
	for _, key := range []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSLMODE"} {
		t.Setenv(key, "")
	}
	for _, key := range []string{"PG_TEST_HOST", "PG_TEST_PORT", "PG_TEST_USER", "PG_TEST_PASSWORD", "PG_TEST_DB", "PG_TEST_SSLMODE"} {
		t.Setenv(key, "")
	}
	// Unset/empty must still resolve to something a server can answer on, and to
	// a database every server has.
	for _, role := range []string{RoleDev, RoleTest} {
		server := ServerForRole(role)
		if server.Host != "127.0.0.1" || server.Port != "5432" || server.SSLMode != "disable" {
			t.Fatalf("%s defaults = %+v, want 127.0.0.1:5432 sslmode=disable", role, server)
		}
		if server.Database != "postgres" {
			t.Fatalf("%s default database = %q, want postgres", role, server.Database)
		}
	}
	// Unknown roles fall back to the dev namespace rather than inventing one.
	if got := ServerForRole("whatever"); got.Role != RoleDev {
		t.Fatalf("unknown role resolved to %q, want dev", got.Role)
	}
}

func TestDSNCarriesNameHostAndSSLMode(t *testing.T) {
	server := Server{
		Role: RoleDev, Host: "db.example", Port: "6543",
		User: "someone", Password: "s3cret", SSLMode: "require",
	}
	raw := server.DSN("schema_ui_dev")
	u, err := url.Parse(raw)
	if err != nil {
		t.Fatalf("DSN is not a URL: %v (%q)", err, raw)
	}
	if u.Scheme != "postgres" {
		t.Fatalf("scheme = %q, want postgres", u.Scheme)
	}
	if u.Path != "/schema_ui_dev" {
		t.Fatalf("path = %q, want /schema_ui_dev", u.Path)
	}
	if u.Host != "db.example:6543" {
		t.Fatalf("host = %q, want db.example:6543", u.Host)
	}
	if u.User.Username() != "someone" {
		t.Fatalf("user = %q, want someone", u.User.Username())
	}
	if pass, _ := u.User.Password(); pass != "s3cret" {
		t.Fatal("password not carried in the DSN")
	}
	if got := u.Query().Get("sslmode"); got != "require" {
		t.Fatalf("sslmode = %q, want require", got)
	}
}

// TestEnsureIsIdempotentAndDropRemoves exercises the provisioning against a real
// server. It is gated exactly like the other postgres integration tests: without
// PG_TEST_* it skips and says why, and it only ever touches a dedicated
// disposable database name. It uses the TEST role, whose namespace is the one the
// gate above reads.
func TestEnsureIsIdempotentAndDropRemoves(t *testing.T) {
	if pgtest.DSN() == "" {
		t.Skip("postgres test env not set (PG_TEST_*); skipping pgsetup integration test")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	server := ServerForRole(RoleTest)

	const scratch = "schema_ui_w35_scratch"
	// Cleanup is idempotent: dropping a missing database reports (false, nil)
	// rather than erroring, so re-running cleanup is safe.
	dropped, err := server.Drop(ctx, scratch)
	if err != nil {
		t.Fatalf("pre-clean drop: %v", err)
	}
	if dropped {
		t.Fatalf("pre-clean drop reported a drop although %s should not exist yet", scratch)
	}
	present, err := server.Exists(ctx, scratch)
	if err != nil {
		t.Fatalf("Exists after drop: %v", err)
	}
	if present {
		t.Fatalf("%s still exists after Drop", scratch)
	}

	created, err := server.Ensure(ctx, scratch)
	if err != nil {
		t.Fatalf("Ensure: %v", err)
	}
	if !created {
		t.Fatal("Ensure reported created=false on a missing database")
	}
	// Idempotence: the second call must not recreate or error.
	created, err = server.Ensure(ctx, scratch)
	if err != nil {
		t.Fatalf("second Ensure: %v", err)
	}
	if created {
		t.Fatal("second Ensure reported created=true (not idempotent)")
	}
	present, err = server.Exists(ctx, scratch)
	if err != nil {
		t.Fatalf("Exists after Ensure: %v", err)
	}
	if !present {
		t.Fatalf("%s missing after Ensure", scratch)
	}
	// And Drop reports the removal it performed.
	dropped, err = server.Drop(ctx, scratch)
	if err != nil {
		t.Fatalf("drop: %v", err)
	}
	if !dropped {
		t.Fatalf("Drop reported no drop although %s existed", scratch)
	}
	present, err = server.Exists(ctx, scratch)
	if err != nil {
		t.Fatalf("Exists after Drop: %v", err)
	}
	if present {
		t.Fatalf("%s still exists after Drop", scratch)
	}

	// A name that could escape the identifier shape must be refused before any
	// SQL is built, for both Ensure and Drop.
	if _, err := server.Ensure(ctx, "bad name"); err == nil {
		t.Fatal("Ensure accepted an invalid database name")
	}
	if _, err := server.Drop(ctx, "bad;name"); err == nil {
		t.Fatal("Drop accepted an invalid database name")
	}
}
