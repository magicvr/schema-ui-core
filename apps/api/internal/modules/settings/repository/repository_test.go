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
	for _, logo := range []string{"//cdn.example/logo.svg", "ftp://example.com/logo.svg", "logo.svg", "/bad path.svg", `/\evil.com/logo.svg`, `/assets\evil.png`, `\evil.com`} {
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
	settings, err = repository.PatchSiteSettings(&title, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, now.Add(time.Second))
	if err != nil {
		t.Fatalf("title-only patch: %v", err)
	}
	if settings.SiteTitle != "Operations" || settings.LogoURL != "https://cdn.example/logo.svg" {
		t.Fatalf("title-only patch overwrote unsubmitted logo: %+v", settings)
	}
	logo := "/assets/logo.svg"
	settings, err = repository.PatchSiteSettings(nil, &logo, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, now.Add(2*time.Second))
	if err != nil {
		t.Fatalf("logo-only patch: %v", err)
	}
	if settings.SiteTitle != "Operations" || settings.LogoURL != "/assets/logo.svg" {
		t.Fatalf("logo-only patch overwrote unsubmitted title: %+v", settings)
	}
}

func TestRepositoryVp007FieldPatchMergeAndValidation(t *testing.T) {
	repository, _ := openSettingsRepository(t, "settings-vp007.db")
	now := time.Now().UTC().Truncate(time.Second)

	locale := "zh-CN"
	theme := "dark"
	timezone := "Asia/Shanghai"
	light := "/assets/logo-light.svg"
	dark := "/assets/logo-dark.svg"
	favicon := "/favicon.ico"
	settings, err := repository.PatchSiteSettings(nil, nil, &light, &dark, &favicon, &locale, &timezone, &theme, nil, nil, nil, nil, now)
	if err != nil {
		t.Fatalf("vp007 patch: %v", err)
	}
	if settings.DefaultLocale != "zh-CN" || settings.DefaultTheme != "dark" || settings.SiteTimezone != "Asia/Shanghai" {
		t.Fatalf("localization/appearance = %+v", settings)
	}
	if settings.LogoURLLight != "/assets/logo-light.svg" || settings.LogoURLDark != "/assets/logo-dark.svg" || settings.FaviconURL != "/favicon.ico" {
		t.Fatalf("branding fields = %+v", settings)
	}

	// Field-level merge: a locale-only patch must not touch theme/timezone.
	locale = "en-US"
	settings, err = repository.PatchSiteSettings(nil, nil, nil, nil, nil, &locale, nil, nil, nil, nil, nil, nil, now.Add(time.Second))
	if err != nil {
		t.Fatalf("locale-only patch: %v", err)
	}
	if settings.DefaultLocale != "en-US" || settings.DefaultTheme != "dark" || settings.SiteTimezone != "Asia/Shanghai" {
		t.Fatalf("locale-only patch overwrote unsubmitted fields: %+v", settings)
	}

	// Empty string clears a branding field.
	empty := ""
	settings, err = repository.PatchSiteSettings(nil, nil, &empty, nil, nil, nil, nil, nil, nil, nil, nil, nil, now.Add(2*time.Second))
	if err != nil {
		t.Fatalf("clear light logo: %v", err)
	}
	if settings.LogoURLLight != "" || settings.LogoURLDark != "/assets/logo-dark.svg" {
		t.Fatalf("clear semantics = %+v", settings)
	}

	// Validation: enums + IANA timezone.
	badLocale := "fr-FR"
	if _, err := repository.PatchSiteSettings(nil, nil, nil, nil, nil, &badLocale, nil, nil, nil, nil, nil, nil, now); !errors.Is(err, ErrInvalidDefaultLocale) {
		t.Fatalf("bad locale = %v, want ErrInvalidDefaultLocale", err)
	}
	badTheme := "neon"
	if _, err := repository.PatchSiteSettings(nil, nil, nil, nil, nil, nil, nil, &badTheme, nil, nil, nil, nil, now); !errors.Is(err, ErrInvalidDefaultTheme) {
		t.Fatalf("bad theme = %v, want ErrInvalidDefaultTheme", err)
	}
	badTimezone := "Foo/Bar"
	if _, err := repository.PatchSiteSettings(nil, nil, nil, nil, nil, nil, &badTimezone, nil, nil, nil, nil, nil, now); !errors.Is(err, ErrInvalidSiteTimezone) {
		t.Fatalf("bad timezone = %v, want ErrInvalidSiteTimezone", err)
	}
	// A rejected patch must leave the previous values untouched.
	settings, err = repository.GetSiteSettings()
	if err != nil {
		t.Fatalf("re-read: %v", err)
	}
	if settings.DefaultLocale != "en-US" || settings.DefaultTheme != "dark" {
		t.Fatalf("rejected patch mutated state: %+v", settings)
	}

	// Reset restores the frozen defaults.
	settings, err = repository.ResetSiteSettings(now.Add(3 * time.Second))
	if err != nil {
		t.Fatalf("reset: %v", err)
	}
	if settings.SiteTitle != settingsmigration.DefaultSiteTitle ||
		settings.LogoURL != "" || settings.LogoURLLight != "" || settings.LogoURLDark != "" || settings.FaviconURL != "" ||
		settings.DefaultLocale != "auto" || settings.SiteTimezone != "auto" || settings.DefaultTheme != "auto" ||
		settings.OperationLogRetentionDays != settingsmigration.DefaultOperationLogRetentionDays ||
		settings.OperationLogExpirationAction != settingsmigration.DefaultOperationLogExpirationAction {
		t.Fatalf("reset defaults = %+v", settings)
	}
}

func TestRepositoryOperationLogRetentionPatch(t *testing.T) {
	repository, _ := openSettingsRepository(t, "settings-retention.db")
	now := time.Now().UTC().Truncate(time.Second)

	settings, err := repository.GetSiteSettings()
	if err != nil {
		t.Fatalf("defaults: %v", err)
	}
	if settings.OperationLogRetentionDays != 90 || settings.OperationLogExpirationAction != "archive" {
		t.Fatalf("default retention = %+v", settings)
	}

	days := 30
	action := "delete"
	settings, err = repository.PatchSiteSettings(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, &days, &action, now)
	if err != nil {
		t.Fatalf("retention patch: %v", err)
	}
	if settings.OperationLogRetentionDays != 30 || settings.OperationLogExpirationAction != "delete" {
		t.Fatalf("patched retention = %+v", settings)
	}

	title := "Kept"
	settings, err = repository.PatchSiteSettings(&title, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, now.Add(time.Second))
	if err != nil {
		t.Fatalf("title-only: %v", err)
	}
	if settings.SiteTitle != "Kept" || settings.OperationLogRetentionDays != 30 || settings.OperationLogExpirationAction != "delete" {
		t.Fatalf("title-only overwrote retention: %+v", settings)
	}

	badDays := 0
	if _, err := repository.PatchSiteSettings(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, &badDays, nil, now); !errors.Is(err, ErrInvalidRetentionDays) {
		t.Fatalf("days 0 = %v", err)
	}
	badDays = 4000
	if _, err := repository.PatchSiteSettings(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, &badDays, nil, now); !errors.Is(err, ErrInvalidRetentionDays) {
		t.Fatalf("days 4000 = %v", err)
	}
	badAction := "compress"
	if _, err := repository.PatchSiteSettings(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, &badAction, now); !errors.Is(err, ErrInvalidExpirationAction) {
		t.Fatalf("bad action = %v", err)
	}
}

// D6 hardening: an explicit empty string for defaultLocale/defaultTheme is
// normalized to "auto" instead of being stored as "" (which configuration.Validate
// would reject and the GET projection would leak).
func TestRepositoryEmptyLocaleThemeNormalizedToAuto(t *testing.T) {
	repository, _ := openSettingsRepository(t, "settings-d6.db")
	now := time.Now().UTC().Truncate(time.Second)

	empty := ""
	settings, err := repository.PatchSiteSettings(nil, nil, nil, nil, nil, &empty, nil, &empty, nil, nil, nil, nil, now)
	if err != nil {
		t.Fatalf("empty locale/theme patch: %v", err)
	}
	if settings.DefaultLocale != "auto" {
		t.Fatalf("DefaultLocale = %q, want auto", settings.DefaultLocale)
	}
	if settings.DefaultTheme != "auto" {
		t.Fatalf("DefaultTheme = %q, want auto", settings.DefaultTheme)
	}

	// The stored row agrees with the projection.
	reRead, err := repository.GetSiteSettings()
	if err != nil {
		t.Fatalf("re-read: %v", err)
	}
	if reRead.DefaultLocale != "auto" || reRead.DefaultTheme != "auto" {
		t.Fatalf("stored = %+v, want locale/theme auto", reRead)
	}
}
