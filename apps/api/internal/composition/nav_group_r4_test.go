package composition

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/internal/config"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

type r4NavigationCase struct {
	name         string
	profile      string
	extra        []string
	withTelegram bool
	groups       map[string][]string
	wantUser     []string
	wantTop      []string
	wantExamples bool
}

func r4NavigationCases(t *testing.T) []r4NavigationCase {
	t.Helper()
	admin, err := kernel.ResolveProfile("admin", nil)
	if err != nil {
		t.Fatal(err)
	}
	withAdmin := func(items ...string) []string {
		modules := append([]string(nil), admin.Modules...)
		return append(modules, items...)
	}
	return []r4NavigationCase{
		{
			name:    "mvp",
			profile: "mvp",
			groups: map[string][]string{
				"manifest.nav.group.workspace":      {"dashboard"},
				"manifest.nav.group.identityAccess": {"users", "roles"},
			},
			wantUser: []string{"account"},
		},
		{
			name:    "admin",
			profile: "admin",
			groups: map[string][]string{
				"manifest.nav.group.workspace":      {"dashboard"},
				"manifest.nav.group.identityAccess": {"users", "roles", "data-permission"},
				"manifest.nav.group.contentData":    {"file-library", "data-dictionary"},
				"manifest.nav.group.operations":     {"activity", "system-monitoring", "scheduled-tasks", "recycle-bin"},
				"manifest.nav.group.communications": {"mail", "mail-outbox"},
				"manifest.nav.group.commerce":       {"wallet", "wallet-vouchers"},
			},
			wantUser: []string{"account", "my-wallet", "settings"},
		},
		{
			name:    "demo",
			profile: "demo",
			groups: map[string][]string{
				"manifest.nav.group.workspace":      {"dashboard"},
				"manifest.nav.group.identityAccess": {"users", "roles"},
			},
			wantUser:     []string{"account"},
			wantTop:      []string{"overview"},
			wantExamples: true,
		},
		{
			name:         "admin+telegram",
			profile:      "custom",
			extra:        withAdmin("channel.telegram"),
			withTelegram: true,
			groups: map[string][]string{
				"manifest.nav.group.workspace":      {"dashboard"},
				"manifest.nav.group.identityAccess": {"users", "roles", "data-permission"},
				"manifest.nav.group.contentData":    {"file-library", "data-dictionary"},
				"manifest.nav.group.operations":     {"activity", "system-monitoring", "scheduled-tasks", "recycle-bin"},
				"manifest.nav.group.communications": {"mail", "mail-outbox", "telegram-settings"},
				"manifest.nav.group.commerce":       {"wallet", "wallet-vouchers"},
			},
			wantUser: []string{"account", "my-wallet", "settings"},
		},
		{
			name:         "admin+digitaloffer+telegram",
			profile:      "custom",
			extra:        withAdmin("channel.telegram", "biz.digital-offer"),
			withTelegram: true,
			groups: map[string][]string{
				"manifest.nav.group.workspace":      {"dashboard"},
				"manifest.nav.group.identityAccess": {"users", "roles", "data-permission"},
				"manifest.nav.group.contentData":    {"file-library", "data-dictionary"},
				"manifest.nav.group.operations":     {"activity", "system-monitoring", "scheduled-tasks", "recycle-bin"},
				"manifest.nav.group.communications": {"mail", "mail-outbox", "telegram-settings"},
				"manifest.nav.group.commerce":       {"wallet", "wallet-vouchers", "digitaloffer-offers", "digitaloffer-entitlements", "digitaloffer-purchases"},
			},
			wantUser: []string{"account", "my-wallet", "settings"},
		},
	}
}

func fetchR4Manifest(t *testing.T, tc r4NavigationCase) (struct {
	Navigation struct {
		Top     []json.RawMessage `json:"top"`
		Sidebar []json.RawMessage `json:"sidebar"`
		User    []json.RawMessage `json:"user"`
	} `json:"navigation"`
}, *kernel.Plan) {
	t.Helper()
	cfg := &config.Config{ProfileName: tc.profile, ModulesEnabled: tc.extra}
	plan, err := ResolvePlan(cfg)
	if err != nil {
		t.Fatalf("ResolvePlan(%s): %v", tc.name, err)
	}
	mux := s5Mux(t, plan, tc.withTelegram)
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/.well-known/schema-ui/app-manifest.json", nil))
	if response.Code != http.StatusOK {
		t.Fatalf("%s manifest status = %d: %s", tc.name, response.Code, response.Body.String())
	}
	var document struct {
		Navigation struct {
			Top     []json.RawMessage `json:"top"`
			Sidebar []json.RawMessage `json:"sidebar"`
			User    []json.RawMessage `json:"user"`
		} `json:"navigation"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &document); err != nil {
		t.Fatalf("%s manifest JSON: %v", tc.name, err)
	}
	return document, &plan
}

func TestR4NavigationGroupProfileMatrix(t *testing.T) {
	for _, tc := range r4NavigationCases(t) {
		t.Run(tc.name, func(t *testing.T) {
			document, _ := fetchR4Manifest(t, tc)
			if len(document.Navigation.Sidebar) == 0 {
				t.Fatal("sidebar must contain the Workspace group")
			}
			var firstSidebar struct {
				Label    string            `json:"label"`
				LabelKey string            `json:"labelKey"`
				Items    []json.RawMessage `json:"items"`
			}
			if err := json.Unmarshal(document.Navigation.Sidebar[0], &firstSidebar); err != nil {
				t.Fatal(err)
			}
			if firstSidebar.Label != "Workspace · WORKSPACE" || firstSidebar.LabelKey != "manifest.nav.group.workspace" || len(firstSidebar.Items) != 1 {
				t.Fatalf("first sidebar item = %+v, want Workspace group with Dashboard", firstSidebar)
			}
			var dashboard struct {
				PageRef string `json:"pageRef"`
				Label   string `json:"label"`
			}
			if err := json.Unmarshal(firstSidebar.Items[0], &dashboard); err != nil {
				t.Fatal(err)
			}
			if dashboard.PageRef != "dashboard" || dashboard.Label != "Dashboard · 01" {
				t.Fatalf("workspace dashboard = %+v, want registered 01 secondary", dashboard)
			}
			for groupKey, wantRefs := range tc.groups {
				gotRefs := manifestGroupPageRefs(document.Navigation.Sidebar, groupKey)
				if len(gotRefs) != len(wantRefs) {
					t.Fatalf("group %s refs=%v, want exact %v", groupKey, gotRefs, wantRefs)
				}
				for _, pageRef := range wantRefs {
					if !gotRefs[pageRef] {
						t.Fatalf("group %s missing %s; got=%v", groupKey, pageRef, gotRefs)
					}
				}
			}
			userRefs := collectManifestNavigationPageRefs(document.Navigation.User)
			if len(userRefs) != len(tc.wantUser) {
				t.Fatalf("user refs=%v, want %v", userRefs, tc.wantUser)
			}
			for _, pageRef := range tc.wantUser {
				if !userRefs[pageRef] {
					t.Fatalf("user slot missing %s; refs=%v", pageRef, userRefs)
				}
			}
			topRefs := collectManifestNavigationPageRefs(document.Navigation.Top)
			if len(topRefs) != len(tc.wantTop) {
				t.Fatalf("top refs=%v, want %v", topRefs, tc.wantTop)
			}
			for _, pageRef := range tc.wantTop {
				if !topRefs[pageRef] {
					t.Fatalf("top slot missing %s; refs=%v", pageRef, topRefs)
				}
			}
			examplesRefs := manifestGroupPageRefs(document.Navigation.Sidebar, "manifest.nav.examples")
			if tc.wantExamples && len(examplesRefs) == 0 {
				t.Fatal("demo Examples group must remain present")
			}
			if !tc.wantExamples && len(examplesRefs) != 0 {
				t.Fatalf("Examples group leaked into non-demo profile: %v", examplesRefs)
			}
		})
	}
}
