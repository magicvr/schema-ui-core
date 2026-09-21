package composition

import (
	"path/filepath"
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/internal/backup"
	"github.com/magicvr/schema-ui-core/apps/api/internal/config"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

// TestRecoveryArtifactsDirIsAbsolute pins the dev-startup defect found on
// 2026-09-21: the shipped config sets a RELATIVE db.path ("./data/schema-ui.db"),
// so recoveryArtifactsDir returned "data/recovery" and the postgres provider
// handed it to `docker run -v data\recovery:/vp040-backup`, which docker rejects
// ("includes invalid characters for a local volume name ... use absolute path").
// Startup then failed at the class-A rollback artifact, so the dev launcher could
// not come up whenever the configured dialect was postgres.
func TestRecoveryArtifactsDirIsAbsolute(t *testing.T) {
	cases := []struct {
		name string
		cfg  *config.Config
		want string // expected final path element, or "" for an empty result
	}{
		{
			name: "relative sqlite path",
			cfg:  &config.Config{DBPath: filepath.Join("data", "schema-ui.db")},
			want: "recovery",
		},
		{
			name: "relative dot-slash sqlite path",
			cfg:  &config.Config{DBPath: filepath.Join(".", "data", "schema-ui.db")},
			want: "recovery",
		},
		{
			name: "absolute sqlite path stays absolute",
			cfg:  &config.Config{DBPath: filepath.Join(t.TempDir(), "schema-ui.db")},
			want: "recovery",
		},
		{
			name: "dsn-only postgres path",
			cfg:  &config.Config{DBDSN: "postgres://sa@127.0.0.1:5432/postgres?sslmode=disable"},
			want: "", // keyed directory; asserted absolute below
		},
		{
			name: "nothing configured",
			cfg:  &config.Config{},
			want: "",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := recoveryArtifactsDir(tc.cfg)
			if tc.cfg.DBPath == "" && tc.cfg.DBDSN == "" {
				if got != "" {
					t.Fatalf("recoveryArtifactsDir(%+v) = %q, want empty", tc.cfg, got)
				}
				return
			}
			if got == "" {
				t.Fatalf("recoveryArtifactsDir(%+v) = empty, want a directory", tc.cfg)
			}
			if !filepath.IsAbs(got) {
				t.Fatalf("recoveryArtifactsDir(%+v) = %q, want an ABSOLUTE path "+
					"(docker bind mounts reject a relative source)", tc.cfg, got)
			}
			if tc.want != "" && filepath.Base(got) != tc.want {
				t.Fatalf("recoveryArtifactsDir(%+v) = %q, want base %q", tc.cfg, got, tc.want)
			}
		})
	}
}

// TestRecoveryWiringHandsPgProviderAnAbsoluteWorkDir closes the loop: the
// provider is the component that actually builds the `docker run -v` argument,
// so its WorkDir must be absolute for a relative db.path.
func TestRecoveryWiringHandsPgProviderAnAbsoluteWorkDir(t *testing.T) {
	cfg := &config.Config{
		DBPath:                filepath.Join("data", "schema-ui.db"),
		DBDSN:                 "postgres://sa@127.0.0.1:5432/postgres?sslmode=disable",
		DBClientDockerNetwork: "schema-ui_default",
	}
	port, creator, dir := recoveryWiring(kernel.DialectPostgres, cfg)
	if port == nil {
		t.Fatal("postgres wiring returned a nil recovery-point port")
	}
	if creator == nil {
		t.Fatal("postgres wiring returned a nil rollback creator")
	}
	provider, ok := creator.(backup.PgProvider)
	if !ok {
		t.Fatalf("rollback creator is %T, want backup.PgProvider", creator)
	}
	if !filepath.IsAbs(provider.WorkDir) {
		t.Fatalf("PgProvider.WorkDir = %q, want absolute (docker -v source)", provider.WorkDir)
	}
	if !filepath.IsAbs(dir) {
		t.Fatalf("wired artifact dir = %q, want absolute", dir)
	}
	if provider.ClientImage != backup.DefaultPgClientImage {
		t.Fatalf("PgProvider.ClientImage = %q, want %q", provider.ClientImage, backup.DefaultPgClientImage)
	}
	if provider.ClientDockerNetwork != "schema-ui_default" {
		t.Fatalf("PgProvider.ClientDockerNetwork = %q, want configured network", provider.ClientDockerNetwork)
	}

	// A configuration with neither path nor DSN must disable the anchors rather
	// than inventing a directory (the store records that fact).
	if p, c, d := recoveryWiring(kernel.DialectPostgres, &config.Config{}); p != nil || c != nil || d != "" {
		t.Fatalf("empty config wiring = (%v, %v, %q), want all disabled", p, c, d)
	}
	// SQLite keeps a service but no rollback creator (the runner writes its own
	// VACUUM INTO artifacts).
	p, c, d := recoveryWiring(kernel.DialectSQLite, &config.Config{DBPath: filepath.Join("data", "schema-ui.db")})
	if p == nil || c != nil || d == "" {
		t.Fatalf("sqlite wiring = (%v, %v, %q), want (service, nil, dir)", p, c, d)
	}
}
