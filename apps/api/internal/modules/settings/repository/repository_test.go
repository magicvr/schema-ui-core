package repository

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"testing"
	"time"

	settingsmigration "github.com/magicvr/schema-ui-core/apps/api/internal/modules/settings/migration"
	"github.com/magicvr/schema-ui-core/apps/api/internal/store"
	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
)

func openSettingsRepository(t *testing.T, name string) (*Repository, *store.Store) {
	t.Helper()
	platform, err := testsupport.OpenStore(filepath.Join(t.TempDir(), name), "admin", "hash", true)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = platform.Close() })
	return New(platform), platform
}

func TestRepositoryDefaultsAndMissingRowInsert(t *testing.T) {
	repository, platform := openSettingsRepository(t, "settings-default.db")
	if err := platform.WithTx(context.Background(), func(tx *sql.Tx) error {
		_, err := tx.Exec(`DELETE FROM site_settings WHERE id = 'default'`)
		return err
	}); err != nil {
		t.Fatalf("delete singleton: %v", err)
	}

	settings, err := repository.GetSiteSettings()
	if err != nil {
		t.Fatalf("GetSiteSettings: %v", err)
	}
	if settings.ID != "default" || settings.SiteTitle != settingsmigration.DefaultSiteTitle || settings.LogoURL != "" || !settings.UpdatedAt.Equal(time.Unix(0, 0).UTC()) {
		t.Fatalf("defaults = %+v", settings)
	}

	now := time.Date(2026, 8, 6, 9, 30, 0, 0, time.UTC)
	settings, err = repository.UpdateSiteSettings("  Operations  ", " /assets/logo.svg ", now)
	if err != nil {
		t.Fatalf("UpdateSiteSettings: %v", err)
	}
	if settings.SiteTitle != "Operations" || settings.LogoURL != "/assets/logo.svg" || !settings.UpdatedAt.Equal(now) {
		t.Fatalf("updated settings = %+v", settings)
	}
}

func TestRepositoryValidationAndUpdate(t *testing.T) {
	repository, _ := openSettingsRepository(t, "settings-validation.db")
	now := time.Now().UTC().Truncate(time.Second)
	if _, err := repository.UpdateSiteSettings("  ", "", now); !errors.Is(err, ErrInvalidSiteTitle) {
		t.Fatalf("empty title = %v, want ErrInvalidSiteTitle", err)
	}
	for _, logo := range []string{"//cdn.example/logo.svg", "ftp://example.com/logo.svg", "logo.svg", "/bad path.svg"} {
		if _, err := repository.UpdateSiteSettings("Admin", logo, now); !errors.Is(err, ErrInvalidLogoURL) {
			t.Fatalf("logo %q = %v, want ErrInvalidLogoURL", logo, err)
		}
	}
	settings, err := repository.UpdateSiteSettings("Admin Console", "https://cdn.example/logo.svg", now)
	if err != nil {
		t.Fatalf("absolute logo update: %v", err)
	}
	if settings.SiteTitle != "Admin Console" || settings.LogoURL != "https://cdn.example/logo.svg" {
		t.Fatalf("settings = %+v", settings)
	}
	title := "Operations"
	settings, err = repository.PatchSiteSettings(&title, nil, now.Add(time.Second))
	if err != nil {
		t.Fatalf("title-only patch: %v", err)
	}
	if settings.SiteTitle != "Operations" || settings.LogoURL != "https://cdn.example/logo.svg" {
		t.Fatalf("title-only patch overwrote unsubmitted logo: %+v", settings)
	}
	logo := "/assets/logo.svg"
	settings, err = repository.PatchSiteSettings(nil, &logo, now.Add(2*time.Second))
	if err != nil {
		t.Fatalf("logo-only patch: %v", err)
	}
	if settings.SiteTitle != "Operations" || settings.LogoURL != "/assets/logo.svg" {
		t.Fatalf("logo-only patch overwrote unsubmitted title: %+v", settings)
	}
}
