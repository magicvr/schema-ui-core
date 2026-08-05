// Package configuration owns the admin.settings runtime configuration contract.
package configuration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	settingsmigration "github.com/magicvr/schema-ui-core/apps/api/internal/modules/settings/migration"
)

const Namespace = "settings.branding"

type Branding struct {
	SiteTitle string `json:"siteTitle"`
	LogoURL   string `json:"logoUrl"`
}

// Contribution returns the stable namespace, defaults and validation rule used
// by both runtime aggregation and the Settings change event.
func Contribution() kernel.ConfigurationContribution {
	defaults, err := json.Marshal(Branding{SiteTitle: settingsmigration.DefaultSiteTitle, LogoURL: ""})
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

// Validate enforces the version-1 branding payload without accepting unknown
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
	if err := validateLogoURL(branding.LogoURL); err != nil {
		return err
	}
	return nil
}

func validateLogoURL(raw string) error {
	if raw == "" {
		return nil
	}
	if strings.TrimSpace(raw) != raw {
		return fmt.Errorf("logoUrl must be trimmed")
	}
	if strings.HasPrefix(raw, "/") {
		if strings.HasPrefix(raw, "//") || strings.ContainsAny(raw, " \t\r\n") {
			return fmt.Errorf("logoUrl same-origin path is invalid")
		}
		return nil
	}
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return fmt.Errorf("logoUrl must be empty, a same-origin path, or an http(s) URL")
	}
	scheme := strings.ToLower(parsed.Scheme)
	if scheme != "http" && scheme != "https" {
		return fmt.Errorf("logoUrl must use http or https")
	}
	return nil
}
