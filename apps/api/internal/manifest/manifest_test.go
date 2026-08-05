package manifest

import (
	"encoding/json"
	"strings"
	"testing"

	rolesmanifest "github.com/magicvr/schema-ui-core/apps/api/internal/modules/roles/manifest"
	usersmanifest "github.com/magicvr/schema-ui-core/apps/api/internal/modules/users/manifest"
)

const manifestApp = `{"name":"Test","version":"1.0.0"}`

func TestAggregateIsDeterministicAndMergesCapabilitiesPagesAndNavigation(t *testing.T) {
	first := `{"protocolVersion":"0.1.3","requiredCapabilities":["z.cap","a.cap"],"app":` + manifestApp + `,"pages":[{"pageId":"one","route":"/one"}],"navigation":{"top":[{"id":"one","pageRef":"one"}],"sidebar":[],"user":[]}}`
	second := `{"protocolVersion":"0.1.3","requiredCapabilities":["b.cap"],"app":` + manifestApp + `,"pages":[{"pageId":"two","route":"/two"}],"navigation":{"top":[{"id":"two","pageRef":"two"}],"sidebar":[],"user":[]}}`
	got, err := Aggregate([]Fragment{{ModuleID: "z-module", Raw: []byte(second)}, {ModuleID: "a-module", Raw: []byte(first)}})
	if err != nil {
		t.Fatal(err)
	}
	var decoded struct {
		RequiredCapabilities []string `json:"requiredCapabilities"`
		Pages                []struct {
			PageID string `json:"pageId"`
		} `json:"pages"`
		Navigation struct {
			Top []json.RawMessage `json:"top"`
		} `json:"navigation"`
	}
	if err := json.Unmarshal(got, &decoded); err != nil {
		t.Fatal(err)
	}
	if strings.Join(decoded.RequiredCapabilities, ",") != "a.cap,b.cap,z.cap" {
		t.Fatalf("capabilities = %v", decoded.RequiredCapabilities)
	}
	if len(decoded.Pages) != 2 || decoded.Pages[0].PageID != "one" || decoded.Pages[1].PageID != "two" {
		t.Fatalf("pages = %+v", decoded.Pages)
	}
	if len(decoded.Navigation.Top) != 2 {
		t.Fatalf("top navigation = %s", got)
	}
}

func TestAggregateRejectsProtocolAppPageAndNavigationConflicts(t *testing.T) {
	base := func(protocol, app, page, nav string) []byte {
		return []byte(`{"protocolVersion":"` + protocol + `","requiredCapabilities":[],"app":` + app + `,"pages":[{"pageId":"` + page + `"}],"navigation":{"top":[{"id":"` + nav + `"}],"sidebar":[],"user":[]}}`)
	}
	cases := []struct {
		name string
		a, b []byte
	}{
		{name: "protocol", a: base("0.1.3", manifestApp, "one", "one"), b: base("0.2.0", manifestApp, "two", "two")},
		{name: "app", a: base("0.1.3", manifestApp, "one", "one"), b: base("0.1.3", `{"name":"Other"}`, "two", "two")},
		{name: "page", a: base("0.1.3", manifestApp, "same", "one"), b: base("0.1.3", manifestApp, "same", "two")},
		{name: "navigation", a: base("0.1.3", manifestApp, "one", "same"), b: base("0.1.3", manifestApp, "two", "same")},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := Aggregate([]Fragment{{ModuleID: "a", Raw: tc.a}, {ModuleID: "b", Raw: tc.b}}); err == nil {
				t.Fatal("expected conflict")
			}
		})
	}
}

func TestDefaultManifestIsValid(t *testing.T) {
	data, err := Default()
	if err != nil {
		t.Fatal(err)
	}
	var document map[string]any
	if err := json.Unmarshal(data, &document); err != nil {
		t.Fatal(err)
	}
	if document["protocolVersion"] == nil || document["app"] == nil {
		t.Fatalf("default manifest missing required fields: %s", data)
	}
}

func TestForModulesOnlyPublishesSelectedAdminPages(t *testing.T) {
	// R4 C3.3: users/roles are module-contributed fragments; settings/activity
	// remain baseline-projected.
	data, err := ForModulesWithFragments([]string{
		"core.server-registration",
		"core.auth-session",
		"core.manifest-route",
		"core.navigation-capability",
		"core.schema-render",
		"core.operationlog",
		"admin.roles",
		"admin.users",
	}, []Fragment{
		{ModuleID: "admin.users", Raw: usersmanifest.FragmentJSON},
		{ModuleID: "admin.roles", Raw: rolesmanifest.FragmentJSON},
	})
	if err != nil {
		t.Fatal(err)
	}
	var decoded struct {
		Pages []struct {
			PageID string `json:"pageId"`
		} `json:"pages"`
		Navigation struct {
			Sidebar []json.RawMessage `json:"sidebar"`
			User    []json.RawMessage `json:"user"`
		} `json:"navigation"`
	}
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}
	pageIDs := make(map[string]bool, len(decoded.Pages))
	for _, page := range decoded.Pages {
		pageIDs[page.PageID] = true
	}
	for _, pageID := range []string{"users", "roles"} {
		if !pageIDs[pageID] {
			t.Fatalf("selected page %q missing from %s", pageID, data)
		}
	}
	for _, pageID := range []string{"settings", "activity"} {
		if pageIDs[pageID] {
			t.Fatalf("disabled page %q leaked into %s", pageID, data)
		}
	}
	if len(decoded.Navigation.Sidebar) != 3 || len(decoded.Navigation.User) != 0 {
		t.Fatalf("navigation projection = %+v", decoded.Navigation)
	}
}

func TestAggregateRejectsDuplicateModuleFragments(t *testing.T) {
	fragment := []byte(`{"protocolVersion":"0.1.3","requiredCapabilities":[],"app":{"appId":"test"},"pages":[],"navigation":{"top":[],"sidebar":[],"user":[]}}`)
	if _, err := Aggregate([]Fragment{{ModuleID: "same", Raw: fragment}, {ModuleID: "same", Raw: fragment}}); err == nil {
		t.Fatal("duplicate module fragments must fail closed")
	}
}

// TestAggregateRejectsSecretKey verifies the public-Manifest secrecy rule:
// a fragment carrying a secret-like key fails closed (R4 C4.4).
func TestAggregateRejectsSecretKey(t *testing.T) {
	fragment := []byte(`{"protocolVersion":"0.1.3","requiredCapabilities":[],"app":{"appId":"test"},"pages":[{"pageId":"x","meta":{"token":"s3cr3t"}}],"navigation":{"top":[],"sidebar":[],"user":[]}}`)
	_, err := Aggregate([]Fragment{{ModuleID: "leaky", Raw: fragment}})
	if err == nil {
		t.Fatal("fragment with secret key must fail closed")
	}
}
