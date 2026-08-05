package configuration

import (
	"encoding/json"
	"testing"
)

func TestContributionDefaultsAndValidation(t *testing.T) {
	contribution := Contribution()
	if contribution.Namespace != Namespace || contribution.Key != Namespace {
		t.Fatalf("identity = %s/%s", contribution.Key, contribution.Namespace)
	}
	if err := contribution.Validate(json.RawMessage(contribution.Defaults)); err != nil {
		t.Fatalf("defaults: %v", err)
	}
	for name, raw := range map[string]string{
		"empty title":   `{"siteTitle":"","logoUrl":""}`,
		"padded title":  `{"siteTitle":" Admin ","logoUrl":""}`,
		"unknown key":   `{"siteTitle":"Admin","logoUrl":"","secret":"x"}`,
		"bad scheme":    `{"siteTitle":"Admin","logoUrl":"javascript:alert(1)"}`,
		"protocol path": `{"siteTitle":"Admin","logoUrl":"//example.test/logo.svg"}`,
		"trailing JSON": `{"siteTitle":"Admin","logoUrl":""} {}`,
	} {
		t.Run(name, func(t *testing.T) {
			if err := Validate(json.RawMessage(raw)); err == nil {
				t.Fatal("invalid branding should fail")
			}
		})
	}
}
