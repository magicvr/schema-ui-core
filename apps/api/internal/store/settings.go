// Site branding settings (GOAL-013): singleton row with required site title and
// optional logo URL text (no upload plugin — empty logo means "do not display").
package store

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	settingsmigration "github.com/magicvr/schema-ui-core/apps/api/internal/modules/settings/migration"
)

// SiteSettings is the persisted branding singleton.
type SiteSettings struct {
	ID        string
	SiteTitle string
	LogoURL   string
	UpdatedAt time.Time
}

// GetSiteSettings returns the default branding row, or the frozen defaults when
// the table is empty (should not happen after 0007 seed).
func (s *Store) GetSiteSettings() (*SiteSettings, error) {
	var (
		st        SiteSettings
		updatedAt int64
	)
	err := s.db.QueryRow(
		`SELECT id, site_title, logo_url, updated_at FROM site_settings WHERE id = 'default'`,
	).Scan(&st.ID, &st.SiteTitle, &st.LogoURL, &updatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return &SiteSettings{
			ID:        "default",
			SiteTitle: settingsmigration.DefaultSiteTitle,
			LogoURL:   "",
			UpdatedAt: time.Unix(0, 0).UTC(),
		}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get site settings: %w", err)
	}
	st.UpdatedAt = time.Unix(updatedAt, 0).UTC()
	return &st, nil
}

// UpdateSiteSettings patches the singleton. siteTitle must be non-empty after
// trim; logoURL may be empty (no logo) or a validated http(s)/same-origin path.
func (s *Store) UpdateSiteSettings(siteTitle, logoURL string, now time.Time) (*SiteSettings, error) {
	title := strings.TrimSpace(siteTitle)
	if title == "" {
		return nil, ErrInvalidSiteTitle
	}
	logo, err := normalizeLogoURL(logoURL)
	if err != nil {
		return nil, err
	}
	updatedAt := now.Unix()
	result, err := s.db.Exec(
		`UPDATE site_settings SET site_title = ?, logo_url = ?, updated_at = ? WHERE id = 'default'`,
		title, logo, updatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("update site settings: %w", err)
	}
	n, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("site settings rows affected: %w", err)
	}
	if n == 0 {
		// Fresh edge: insert if missing (defensive after partial migrate).
		if _, err := s.db.Exec(
			`INSERT INTO site_settings (id, site_title, logo_url, updated_at) VALUES ('default', ?, ?, ?)`,
			title, logo, updatedAt,
		); err != nil {
			return nil, fmt.Errorf("insert site settings: %w", err)
		}
	}
	return s.GetSiteSettings()
}

// ErrInvalidSiteTitle and ErrInvalidLogoURL are domain sentinels for settings.
var (
	ErrInvalidSiteTitle = errors.New("store: site title must not be empty")
	ErrInvalidLogoURL   = errors.New("store: invalid logo url")
)

// normalizeLogoURL accepts empty (hide logo), absolute http(s) URLs, or a
// single-slash same-origin path. Rejects javascript:, data:, and protocol-relative.
func normalizeLogoURL(raw string) (string, error) {
	logo := strings.TrimSpace(raw)
	if logo == "" {
		return "", nil
	}
	if strings.HasPrefix(logo, "/") {
		if strings.HasPrefix(logo, "//") {
			return "", ErrInvalidLogoURL
		}
		if strings.ContainsAny(logo, " \t\r\n") {
			return "", ErrInvalidLogoURL
		}
		return logo, nil
	}
	u, err := url.Parse(logo)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return "", ErrInvalidLogoURL
	}
	scheme := strings.ToLower(u.Scheme)
	if scheme != "http" && scheme != "https" {
		return "", ErrInvalidLogoURL
	}
	return u.String(), nil
}
