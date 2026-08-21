// Package configuration owns the admin.settings runtime configuration contract.
package configuration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	settingsmigration "github.com/magicvr/schema-ui-core/apps/api/internal/modules/settings/migration"
)

const Namespace = "settings.branding"

// Branding is the VP-007 S3 startup/branding payload: General + Branding +
// Localization + Appearance defaults consumed by the shell/login/bootstrap.
type Branding struct {
	SiteTitle     string `json:"siteTitle"`
	LogoURL       string `json:"logoUrl"`
	LogoURLLight  string `json:"logoUrlLight"`
	LogoURLDark   string `json:"logoUrlDark"`
	FaviconURL    string `json:"faviconUrl"`
	DefaultLocale string `json:"defaultLocale"`
	SiteTimezone  string `json:"siteTimezone"`
	DefaultTheme  string `json:"defaultTheme"`
}

// Contribution returns the stable namespace, defaults and validation rule used
// by both runtime aggregation and the Settings change event.
func Contribution() kernel.ConfigurationContribution {
	defaults, err := json.Marshal(Branding{
		SiteTitle:     settingsmigration.DefaultSiteTitle,
		LogoURL:       "",
		LogoURLLight:  "",
		LogoURLDark:   "",
		FaviconURL:    "",
		DefaultLocale: "auto",
		SiteTimezone:  "auto",
		DefaultTheme:  "auto",
	})
	if err != nil {
		panic(err)
	}
	return kernel.ConfigurationContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: settingsmigration.ModuleID, Key: Namespace},
		Namespace:            Namespace,
		Defaults:             defaults,
		Validate:             Validate,
	}
}

// Validate enforces the version-2 branding payload without accepting unknown
// keys or a second JSON value.
func Validate(raw json.RawMessage) error {
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	var branding Branding
	if err := decoder.Decode(&branding); err != nil {
		return fmt.Errorf("decode branding: %w", err)
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return fmt.Errorf("branding must contain exactly one JSON value")
	}
	if strings.TrimSpace(branding.SiteTitle) == "" || strings.TrimSpace(branding.SiteTitle) != branding.SiteTitle {
		return fmt.Errorf("siteTitle must be non-empty and trimmed")
	}
	for _, field := range []struct {
		name  string
		value string
	}{
		{"logoUrl", branding.LogoURL},
		{"logoUrlLight", branding.LogoURLLight},
		{"logoUrlDark", branding.LogoURLDark},
		{"faviconUrl", branding.FaviconURL},
	} {
		if err := validateLogoURL(field.value); err != nil {
			return fmt.Errorf("%s: %w", field.name, err)
		}
	}
	switch branding.DefaultLocale {
	case "auto", "zh-CN", "en-US":
	default:
		return fmt.Errorf("defaultLocale must be auto, zh-CN or en-US")
	}
	switch branding.DefaultTheme {
	case "auto", "light", "dark":
	default:
		return fmt.Errorf("defaultTheme must be auto, light or dark")
	}
	if branding.SiteTimezone != "" && branding.SiteTimezone != "auto" {
		if _, err := time.LoadLocation(branding.SiteTimezone); err != nil {
			return fmt.Errorf("siteTimezone is not a valid IANA timezone: %v", err)
		}
	}
	return nil
}

func validateLogoURL(raw string) error {
	if raw == "" {
		return nil
	}
	if strings.TrimSpace(raw) != raw {
		return fmt.Errorf("must be trimmed")
	}
	if strings.HasPrefix(raw, "/") {
		// D-001 P2: backslash rejected in same-origin paths (browsers treat `\`
		// as a separator, so "/\evil.com" would resolve to the external host).
		if strings.HasPrefix(raw, "//") || strings.ContainsAny(raw, " \t\r\n\\") {
			return fmt.Errorf("same-origin path is invalid")
		}
		return nil
	}
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return fmt.Errorf("must be empty, a same-origin path, or an http(s) URL")
	}
	scheme := strings.ToLower(parsed.Scheme)
	if scheme != "http" && scheme != "https" {
		return fmt.Errorf("must use http or https")
	}
	return nil
}
