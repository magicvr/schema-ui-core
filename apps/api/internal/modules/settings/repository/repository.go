// Package repository owns admin.settings persistence.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	settingsmigration "github.com/magicvr/schema-ui-core/apps/api/internal/modules/settings/migration"
)

// TxRunner is the platform transaction boundary consumed by settings.
type TxRunner interface {
	WithTx(context.Context, func(*sql.Tx) error) error
}

// Repository owns the site settings singleton queries.
type Repository struct {
	runner TxRunner
}

// New constructs the settings repository over the platform transaction runner.
func New(runner TxRunner) *Repository {
	return &Repository{runner: runner}
}

// SiteSettings is the persisted system-settings singleton (VP-007 S3:
// General / Branding / Localization / Appearance fields).
type SiteSettings struct {
	ID            string
	SiteTitle     string
	LogoURL       string
	LogoURLLight  string
	LogoURLDark   string
	FaviconURL    string
	DefaultLocale string
	SiteTimezone  string
	DefaultTheme  string
	// W16-F10: optional footer text.
	CopyrightText string
	ICPNumber     string
	UpdatedAt     time.Time
}

var (
	ErrInvalidSiteTitle     = errors.New("settings: site title must not be empty")
	ErrInvalidLogoURL       = errors.New("settings: invalid logo url")
	ErrInvalidDefaultLocale = errors.New("settings: invalid default locale")
	ErrInvalidDefaultTheme  = errors.New("settings: invalid default theme")
	ErrInvalidSiteTimezone  = errors.New("settings: invalid site timezone")
)

// SupportedLocales is the frozen v1 locale set (VP-007).
var SupportedLocales = []string{"zh-CN", "en-US"}

// ValidDefaultLocales are the allowed defaultLocale values.
var ValidDefaultLocales = []string{"", "auto", "zh-CN", "en-US"}

// ValidDefaultThemes are the allowed defaultTheme values.
var ValidDefaultThemes = []string{"", "auto", "light", "dark"}

// GetSiteSettings returns the singleton or the frozen defaults when it is absent.
func (r *Repository) GetSiteSettings() (*SiteSettings, error) {
	var settings *SiteSettings
	err := r.withTx("get site settings", func(tx *sql.Tx) error {
		var err error
		settings, err = getSiteSettings(tx)
		return err
	})
	return settings, err
}

// UpdateSiteSettings is the legacy two-field convenience wrapper (title + logo);
// kept for the composition recovery tests and callers that predate VP-007.
func (r *Repository) UpdateSiteSettings(siteTitle, logoURL string, now time.Time) (*SiteSettings, error) {
	return r.writeSiteSettings(&siteTitle, &logoURL, nil, nil, nil, nil, nil, nil, nil, nil, now)
}

// PatchSiteSettings updates only the supplied fields in one SQL statement.
// This prevents concurrent field-level PATCH requests from overwriting fields
// they did not submit. Empty-string values clear a field (logo/theming);
// validation errors reject the whole patch atomically.
func (r *Repository) PatchSiteSettings(
	siteTitle, logoURL, logoURLLight, logoURLDark, faviconURL, defaultLocale, siteTimezone, defaultTheme, copyrightText, icpNumber *string,
	now time.Time,
) (*SiteSettings, error) {
	return r.writeSiteSettings(siteTitle, logoURL, logoURLLight, logoURLDark, faviconURL, defaultLocale, siteTimezone, defaultTheme, copyrightText, icpNumber, now)
}

// ResetSiteSettings restores every VP-007 field to its frozen default.
func (r *Repository) ResetSiteSettings(now time.Time) (*SiteSettings, error) {
	return r.writeSiteSettings(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, now, true)
}

func (r *Repository) writeSiteSettings(
	siteTitle, logoURL, logoURLLight, logoURLDark, faviconURL, defaultLocale, siteTimezone, defaultTheme, copyrightText, icpNumber *string,
	now time.Time,
	reset ...bool,
) (*SiteSettings, error) {
	title := settingsmigration.DefaultSiteTitle
	titleSet := 0
	if siteTitle != nil {
		title = strings.TrimSpace(*siteTitle)
		if title == "" {
			return nil, ErrInvalidSiteTitle
		}
		titleSet = 1
	}
	logo, logoSet, err := normalizeOptionalLogo(logoURL)
	if err != nil {
		return nil, err
	}
	logoLight, logoLightSet, err := normalizeOptionalLogo(logoURLLight)
	if err != nil {
		return nil, err
	}
	logoDark, logoDarkSet, err := normalizeOptionalLogo(logoURLDark)
	if err != nil {
		return nil, err
	}
	favicon, faviconSet, err := normalizeOptionalLogo(faviconURL)
	if err != nil {
		return nil, err
	}
	locale, localeSet, err := optionalEnum(defaultLocale, ValidDefaultLocales, ErrInvalidDefaultLocale)
	if err != nil {
		return nil, err
	}
	theme, themeSet, err := optionalEnum(defaultTheme, ValidDefaultThemes, ErrInvalidDefaultTheme)
	if err != nil {
		return nil, err
	}
	// "" is accepted for legacy writes but is not a valid stored value: an
	// explicit empty locale/theme means "auto" (D6). Normalizing here keeps the
	// store, the GET projection and configuration.Validate in agreement.
	if localeSet == 1 && locale == "" {
		locale = "auto"
	}
	if themeSet == 1 && theme == "" {
		theme = "auto"
	}
	timezone, timezoneSet, err := optionalEnum(siteTimezone, nil, nil)
	if err != nil {
		return nil, err
	}
	if err := validateTimezone(timezone); err != nil {
		return nil, err
	}
	copyright := ""
	copyrightSet := 0
	if copyrightText != nil {
		copyright = strings.TrimSpace(*copyrightText)
		copyrightSet = 1
	}
	icp := ""
	icpSet := 0
	if icpNumber != nil {
		icp = strings.TrimSpace(*icpNumber)
		icpSet = 1
	}

	forceReset := len(reset) > 0 && reset[0]
	var settings *SiteSettings
	err = r.withTx("update site settings", func(tx *sql.Tx) error {
		var stmt string
		var args []any
		if forceReset {
			stmt = `UPDATE site_settings SET
			  site_title = ?, logo_url = '', logo_url_light = '', logo_url_dark = '',
			  favicon_url = '', default_locale = 'auto', site_timezone = 'auto',
			  default_theme = 'auto', copyright_text = '', icp_number = '', updated_at = ? WHERE id = 'default'`
			args = []any{settingsmigration.DefaultSiteTitle, now.Unix()}
		} else {
			stmt = `INSERT INTO site_settings (
			  id, site_title, logo_url, logo_url_light, logo_url_dark, favicon_url,
			  default_locale, site_timezone, default_theme, copyright_text, icp_number, updated_at)
			 VALUES ('default', ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			 ON CONFLICT(id) DO UPDATE SET
			   site_title = CASE WHEN ? = 1 THEN excluded.site_title ELSE site_settings.site_title END,
			   logo_url = CASE WHEN ? = 1 THEN excluded.logo_url ELSE site_settings.logo_url END,
			   logo_url_light = CASE WHEN ? = 1 THEN excluded.logo_url_light ELSE site_settings.logo_url_light END,
			   logo_url_dark = CASE WHEN ? = 1 THEN excluded.logo_url_dark ELSE site_settings.logo_url_dark END,
			   favicon_url = CASE WHEN ? = 1 THEN excluded.favicon_url ELSE site_settings.favicon_url END,
			   default_locale = CASE WHEN ? = 1 THEN excluded.default_locale ELSE site_settings.default_locale END,
			   site_timezone = CASE WHEN ? = 1 THEN excluded.site_timezone ELSE site_settings.site_timezone END,
			   default_theme = CASE WHEN ? = 1 THEN excluded.default_theme ELSE site_settings.default_theme END,
			   copyright_text = CASE WHEN ? = 1 THEN excluded.copyright_text ELSE site_settings.copyright_text END,
			   icp_number = CASE WHEN ? = 1 THEN excluded.icp_number ELSE site_settings.icp_number END,
			   updated_at = excluded.updated_at`
			args = []any{
				title, logo, logoLight, logoDark, favicon, locale, timezone, theme, copyright, icp, now.Unix(),
				titleSet, logoSet, logoLightSet, logoDarkSet, faviconSet, localeSet, timezoneSet, themeSet, copyrightSet, icpSet,
			}
		}
		if _, err := tx.Exec(stmt, args...); err != nil {
			return fmt.Errorf("upsert singleton: %w", err)
		}
		var err error
		settings, err = getSiteSettings(tx)
		return err
	})
	return settings, err
}

// normalizeOptionalLogo trims + validates a logo-ish URL (empty = clear).
func normalizeOptionalLogo(raw *string) (string, int, error) {
	if raw == nil {
		return "", 0, nil
	}
	value, err := normalizeLogoURL(*raw)
	if err != nil {
		return "", 0, err
	}
	return value, 1, nil
}

// optionalEnum validates a pointer string against the allowed set (nil set =
// free-form like IANA timezone names) and returns (value, set, error).
func optionalEnum(raw *string, allowed []string, invalid error) (string, int, error) {
	if raw == nil {
		return "", 0, nil
	}
	value := strings.TrimSpace(*raw)
	if allowed != nil {
		ok := false
		for _, entry := range allowed {
			if value == entry {
				ok = true
				break
			}
		}
		if !ok {
			return "", 0, invalid
		}
	}
	return value, 1, nil
}

// validateTimezone accepts empty, "auto" or a resolvable IANA timezone name.
func validateTimezone(raw string) error {
	if raw == "" || raw == "auto" {
		return nil
	}
	if _, err := time.LoadLocation(raw); err != nil {
		return ErrInvalidSiteTimezone
	}
	return nil
}

func (r *Repository) withTx(operation string, fn func(*sql.Tx) error) error {
	if r == nil || r.runner == nil {
		return fmt.Errorf("%s: settings repository is not configured", operation)
	}
	if err := r.runner.WithTx(context.Background(), fn); err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}
	return nil
}

func getSiteSettings(row interface{ QueryRow(string, ...any) *sql.Row }) (*SiteSettings, error) {
	var settings SiteSettings
	var updatedAt int64
	err := row.QueryRow(
		`SELECT id, site_title, logo_url, logo_url_light, logo_url_dark, favicon_url,
		        default_locale, site_timezone, default_theme, copyright_text, icp_number, updated_at
		 FROM site_settings WHERE id = 'default'`,
	).Scan(
		&settings.ID, &settings.SiteTitle, &settings.LogoURL, &settings.LogoURLLight,
		&settings.LogoURLDark, &settings.FaviconURL, &settings.DefaultLocale,
		&settings.SiteTimezone, &settings.DefaultTheme, &settings.CopyrightText,
		&settings.ICPNumber, &updatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return &SiteSettings{
			ID:            "default",
			SiteTitle:     settingsmigration.DefaultSiteTitle,
			DefaultLocale: "auto",
			SiteTimezone:  "auto",
			DefaultTheme:  "auto",
			UpdatedAt:     time.Unix(0, 0).UTC(),
		}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query singleton: %w", err)
	}
	settings.UpdatedAt = time.Unix(updatedAt, 0).UTC()
	return &settings, nil
}

func normalizeLogoURL(raw string) (string, error) {
	logo := strings.TrimSpace(raw)
	if logo == "" {
		return "", nil
	}
	if strings.HasPrefix(logo, "/") {
		// D-001 P2: a backslash in a same-origin path is rejected too. Browsers
		// treat `\` as a path separator in URL parsing, so "/\evil.com" would
		// be parsed as the external host "evil.com" and fetch the attacker's
		// origin; requiring a backslash-free single-slash path keeps logo/favicon
		// strictly same-origin.
		if strings.HasPrefix(logo, "//") || strings.ContainsAny(logo, " \t\r\n\\") {
			return "", ErrInvalidLogoURL
		}
		return logo, nil
	}
	parsed, err := url.Parse(logo)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return "", ErrInvalidLogoURL
	}
	scheme := strings.ToLower(parsed.Scheme)
	if scheme != "http" && scheme != "https" {
		return "", ErrInvalidLogoURL
	}
	return parsed.String(), nil
}
