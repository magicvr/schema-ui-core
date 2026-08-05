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

// SiteSettings is the persisted branding singleton.
type SiteSettings struct {
	ID        string
	SiteTitle string
	LogoURL   string
	UpdatedAt time.Time
}

var (
	ErrInvalidSiteTitle = errors.New("settings: site title must not be empty")
	ErrInvalidLogoURL   = errors.New("settings: invalid logo url")
)

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

// UpdateSiteSettings validates and updates the singleton atomically.
func (r *Repository) UpdateSiteSettings(siteTitle, logoURL string, now time.Time) (*SiteSettings, error) {
	return r.writeSiteSettings(&siteTitle, &logoURL, now)
}

// PatchSiteSettings updates only the supplied fields in one SQL statement.
// This prevents concurrent field-level PATCH requests from overwriting fields
// they did not submit.
func (r *Repository) PatchSiteSettings(siteTitle, logoURL *string, now time.Time) (*SiteSettings, error) {
	return r.writeSiteSettings(siteTitle, logoURL, now)
}

func (r *Repository) writeSiteSettings(siteTitle, logoURL *string, now time.Time) (*SiteSettings, error) {
	title := settingsmigration.DefaultSiteTitle
	titleSet := 0
	if siteTitle != nil {
		title = strings.TrimSpace(*siteTitle)
		if title == "" {
			return nil, ErrInvalidSiteTitle
		}
		titleSet = 1
	}
	logo := ""
	logoSet := 0
	if logoURL != nil {
		var err error
		logo, err = normalizeLogoURL(*logoURL)
		if err != nil {
			return nil, err
		}
		logoSet = 1
	}

	var settings *SiteSettings
	err := r.withTx("update site settings", func(tx *sql.Tx) error {
		if _, err := tx.Exec(
			`INSERT INTO site_settings (id, site_title, logo_url, updated_at)
			 VALUES ('default', ?, ?, ?)
			 ON CONFLICT(id) DO UPDATE SET
			   site_title = CASE WHEN ? = 1 THEN excluded.site_title ELSE site_settings.site_title END,
			   logo_url = CASE WHEN ? = 1 THEN excluded.logo_url ELSE site_settings.logo_url END,
			   updated_at = excluded.updated_at`,
			title, logo, now.Unix(), titleSet, logoSet,
		); err != nil {
			return fmt.Errorf("upsert singleton: %w", err)
		}
		var err error
		settings, err = getSiteSettings(tx)
		return err
	})
	return settings, err
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
		`SELECT id, site_title, logo_url, updated_at FROM site_settings WHERE id = 'default'`,
	).Scan(&settings.ID, &settings.SiteTitle, &settings.LogoURL, &updatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return &SiteSettings{
			ID:        "default",
			SiteTitle: settingsmigration.DefaultSiteTitle,
			UpdatedAt: time.Unix(0, 0).UTC(),
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
		if strings.HasPrefix(logo, "//") || strings.ContainsAny(logo, " \t\r\n") {
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
