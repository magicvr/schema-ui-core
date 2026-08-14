package manifest

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

//go:embed app-manifest.json
var defaultManifest []byte

type Fragment struct {
	ModuleID string
	Raw      json.RawMessage
}

type envelope struct {
	ProtocolVersion      string            `json:"protocolVersion"`
	RequiredCapabilities []string          `json:"requiredCapabilities"`
	App                  json.RawMessage   `json:"app"`
	Pages                []json.RawMessage `json:"pages"`
	Navigation           navigation        `json:"navigation"`
}

type navigation struct {
	Top     []json.RawMessage `json:"top"`
	Sidebar []json.RawMessage `json:"sidebar"`
	User    []json.RawMessage `json:"user"`
}

func Default() ([]byte, error) {
	return Aggregate([]Fragment{{ModuleID: "core.manifest-route", Raw: defaultManifest}})
}

// ForModules projects the embedded protocol baseline through the selected
// module set. The embedded file remains a source fixture; the published bytes
// are assembled from the core fragment plus only enabled Admin fragments.
func ForModules(moduleIDs []string) ([]byte, error) {
	return ForModulesWithFragments(moduleIDs, nil)
}

// ForModulesWithFragments is the R4 C4 aggregation: the embedded baseline is
// core-only, and every standard Admin module (users/roles/settings/activity)
// contributes a manifest fragment. The function publishes the baseline as the
// core fragment and merges the enabled provider fragments.
//
// GOAL-013 D-002 §4: an optional navigation order (NodeID list) reorders the
// merged top/sidebar/user slots after aggregation. Items whose NodeID is not
// in the list keep their aggregate-relative order at the end (new modules
// never disappear). order is derived from visibleWhen "features.<nodeID>"
// expressions, falling back to id/pageRef.
func ForModulesWithFragments(moduleIDs []string, moduleFragments []Fragment, order ...[]string) ([]byte, error) {
	enabled := make(map[string]struct{}, len(moduleIDs))
	for _, rawID := range moduleIDs {
		id := strings.TrimSpace(rawID)
		if id == "" {
			return nil, fmt.Errorf("manifest: selected module id is empty")
		}
		enabled[id] = struct{}{}
	}
	if _, ok := enabled["core.manifest-route"]; !ok {
		return nil, fmt.Errorf("manifest: core.manifest-route must be enabled to publish a manifest")
	}

	var base envelope
	if err := json.Unmarshal(defaultManifest, &base); err != nil {
		return nil, fmt.Errorf("manifest: parse embedded baseline: %w", err)
	}
	coreRaw, err := json.Marshal(base)
	if err != nil {
		return nil, fmt.Errorf("manifest: encode core fragment: %w", err)
	}
	allFragments := []Fragment{{ModuleID: "core.manifest-route", Raw: coreRaw}}
	allFragments = append(allFragments, moduleFragments...)
	data, err := Aggregate(allFragments)
	if err != nil {
		return nil, err
	}
	if len(order) > 0 && len(order[0]) > 0 {
		data, err = SortNavigation(data, order[0])
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

func Aggregate(fragments []Fragment) ([]byte, error) {
	if len(fragments) == 0 {
		return nil, fmt.Errorf("manifest: at least one fragment is required")
	}
	ordered := append([]Fragment(nil), fragments...)
	sort.Slice(ordered, func(i, j int) bool { return ordered[i].ModuleID < ordered[j].ModuleID })

	var result envelope
	result.RequiredCapabilities = []string{}
	result.Pages = []json.RawMessage{}
	result.Navigation = navigation{Top: []json.RawMessage{}, Sidebar: []json.RawMessage{}, User: []json.RawMessage{}}
	capabilities := map[string]struct{}{}
	pageOwners := map[string]string{}
	navigationOwners := map[string]string{}
	fragmentOwners := map[string]struct{}{}
	var appCanonical []byte

	for _, fragment := range ordered {
		moduleID := strings.TrimSpace(fragment.ModuleID)
		if moduleID == "" {
			return nil, fmt.Errorf("manifest: fragment module id is required")
		}
		if _, exists := fragmentOwners[moduleID]; exists {
			return nil, fmt.Errorf("manifest: module %s contributes more than one fragment", moduleID)
		}
		fragmentOwners[moduleID] = struct{}{}
		var parsed envelope
		if err := json.Unmarshal(fragment.Raw, &parsed); err != nil {
			return nil, fmt.Errorf("manifest: parse fragment %s: %w", moduleID, err)
		}
		if err := rejectFragmentSecrets(moduleID, fragment.Raw); err != nil {
			return nil, err
		}
		if parsed.ProtocolVersion == "" {
			return nil, fmt.Errorf("manifest: fragment %s has no protocolVersion", moduleID)
		}
		if result.ProtocolVersion == "" {
			result.ProtocolVersion = parsed.ProtocolVersion
		} else if result.ProtocolVersion != parsed.ProtocolVersion {
			return nil, fmt.Errorf("manifest: protocolVersion conflict between fragments: %s and %s", result.ProtocolVersion, parsed.ProtocolVersion)
		}
		canonicalApp, err := canonicalJSON(parsed.App)
		if err != nil {
			return nil, fmt.Errorf("manifest: invalid app in fragment %s: %w", moduleID, err)
		}
		if appCanonical == nil {
			appCanonical = canonicalApp
			result.App = append(json.RawMessage(nil), parsed.App...)
		} else if !bytes.Equal(appCanonical, canonicalApp) {
			return nil, fmt.Errorf("manifest: app identity conflict in fragment %s", moduleID)
		}

		for _, capability := range parsed.RequiredCapabilities {
			capability = strings.TrimSpace(capability)
			if capability == "" {
				return nil, fmt.Errorf("manifest: fragment %s has an empty required capability", moduleID)
			}
			capabilities[capability] = struct{}{}
		}
		for _, page := range parsed.Pages {
			identity, err := pageIdentity(page)
			if err != nil {
				return nil, fmt.Errorf("manifest: fragment %s page: %w", moduleID, err)
			}
			if owner, exists := pageOwners[identity]; exists {
				return nil, fmt.Errorf("manifest: page %q conflicts between %s and %s", identity, owner, moduleID)
			}
			pageOwners[identity] = moduleID
			result.Pages = append(result.Pages, page)
		}
		var navErr error
		result.Navigation.Top, navErr = appendNavigation(result.Navigation.Top, parsed.Navigation.Top, moduleID, navigationOwners)
		if navErr != nil {
			return nil, navErr
		}
		result.Navigation.Sidebar, navErr = appendNavigation(result.Navigation.Sidebar, parsed.Navigation.Sidebar, moduleID, navigationOwners)
		if navErr != nil {
			return nil, navErr
		}
		result.Navigation.User, navErr = appendNavigation(result.Navigation.User, parsed.Navigation.User, moduleID, navigationOwners)
		if navErr != nil {
			return nil, navErr
		}
	}

	for capability := range capabilities {
		result.RequiredCapabilities = append(result.RequiredCapabilities, capability)
	}
	sort.Strings(result.RequiredCapabilities)
	encoded, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("manifest: encode aggregate: %w", err)
	}
	return append(encoded, '\n'), nil
}

// navigationNodeID extracts the navigation-order NodeID for one manifest
// navigation item. The canonical source is the visibleWhen "when" expression,
// which modules author as "$context.features.<nodeID> == true"; id/pageRef are
// fallbacks for items without a feature expression.
func navigationNodeID(raw json.RawMessage) string {
	var item struct {
		ID          string `json:"id"`
		PageRef     string `json:"pageRef"`
		VisibleWhen struct {
			When string `json:"when"`
		} `json:"visibleWhen"`
	}
	if err := json.Unmarshal(raw, &item); err != nil {
		return ""
	}
	if m := featureNodePattern.FindStringSubmatch(item.VisibleWhen.When); len(m) == 2 {
		return m[1]
	}
	for _, value := range []string{item.ID, item.PageRef} {
		if value = strings.TrimSpace(value); value != "" {
			return value
		}
	}
	return ""
}

// featureNodePattern matches "features.<nodeID>" inside a visibleWhen when
// expression (GOAL-013 D-002 §4).
var featureNodePattern = regexp.MustCompile(`features.([A-Za-z0-9_]+)`)

// SortNavigation reorders the top/sidebar/user slots of an aggregated manifest
// document by the given navigation order (NodeID list). Items not present in
// the list keep their relative order and are appended at the end. Returns the
// re-encoded document (same formatting as Aggregate).
func SortNavigation(data []byte, order []string) ([]byte, error) {
	rank := make(map[string]int, len(order))
	for i, id := range order {
		rank[id] = i
	}
	var doc envelope
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("manifest: parse aggregated document for navigation sort: %w", err)
	}
	sortSlot := func(items []json.RawMessage) []json.RawMessage {
		sorted := append([]json.RawMessage(nil), items...)
		sort.SliceStable(sorted, func(i, j int) bool {
			ri, iOK := rank[navigationNodeID(sorted[i])]
			rj, jOK := rank[navigationNodeID(sorted[j])]
			switch {
			case iOK && jOK:
				return ri < rj
			case iOK:
				return true // listed items precede unlisted ones
			default:
				return false
			}
		})
		return sorted
	}
	doc.Navigation.Top = sortSlot(doc.Navigation.Top)
	doc.Navigation.Sidebar = sortSlot(doc.Navigation.Sidebar)
	doc.Navigation.User = sortSlot(doc.Navigation.User)
	encoded, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("manifest: encode sorted document: %w", err)
	}
	return append(encoded, '\n'), nil
}

// StampHomePageRef sets (or removes) app.homePageRef on an already-aggregated
// manifest document while preserving all other content. homePageRef == "" deletes
// the field (used when the enabled set contributes no pages, satisfying the web
// consumer rule that an empty page registry must not declare a home page).
//
// W1 (GOAL-002 / workspace-010): fragments no longer carry homePageRef — every
// fragment's app block is canonically equal so Aggregate's app-identity check
// passes regardless of the enabled set; the assembly layer derives and stamps the
// published home page (D-003 §1, mechanism A).
func StampHomePageRef(data []byte, homePageRef string) ([]byte, error) {
	var doc struct {
		ProtocolVersion      string            `json:"protocolVersion"`
		RequiredCapabilities []string          `json:"requiredCapabilities"`
		App                  json.RawMessage   `json:"app"`
		Pages                []json.RawMessage `json:"pages"`
		Navigation           json.RawMessage   `json:"navigation"`
	}
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("manifest: parse aggregated document for home page stamp: %w", err)
	}
	app := map[string]json.RawMessage{}
	if err := json.Unmarshal(doc.App, &app); err != nil {
		return nil, fmt.Errorf("manifest: parse app block for home page stamp: %w", err)
	}
	if homePageRef == "" {
		delete(app, "homePageRef")
	} else {
		raw, err := json.Marshal(homePageRef)
		if err != nil {
			return nil, err
		}
		app["homePageRef"] = raw
	}
	stampedApp, err := json.Marshal(app)
	if err != nil {
		return nil, fmt.Errorf("manifest: encode app block with home page: %w", err)
	}
	doc.App = stampedApp
	return json.MarshalIndent(doc, "", "  ")
}

func canonicalJSON(raw json.RawMessage) ([]byte, error) {
	if len(bytes.TrimSpace(raw)) == 0 {
		return nil, fmt.Errorf("value is required")
	}
	var value any
	if err := json.Unmarshal(raw, &value); err != nil {
		return nil, err
	}
	return json.Marshal(value)
}

// secretKeyNames are keys that a public, login-before Manifest must never carry
// (freeze package §5: "public Manifest 登录前可读但不得包含 secret、token 或用户
// 个性化信息"). Fragments containing them fail closed.
var secretKeyNames = []string{"password", "token", "secret", "authorization", "apikey", "api_key"}

// rejectFragmentSecrets recursively scans a fragment's JSON for secret-like
// keys and rejects it, so module-authored manifest content cannot leak secrets
// into the public document (R4 C4.4).
func rejectFragmentSecrets(moduleID string, raw []byte) error {
	var value any
	if err := json.Unmarshal(raw, &value); err != nil {
		return nil // structural validation happens elsewhere
	}
	var walk func(any) error
	walk = func(node any) error {
		switch v := node.(type) {
		case map[string]any:
			for key, child := range v {
				lower := strings.ToLower(strings.TrimSpace(key))
				for _, secret := range secretKeyNames {
					if lower == secret {
						return fmt.Errorf("manifest: fragment %s contains secret key %q (public manifest secrecy)", moduleID, key)
					}
				}
				if err := walk(child); err != nil {
					return err
				}
			}
		case []any:
			for _, child := range v {
				if err := walk(child); err != nil {
					return err
				}
			}
		}
		return nil
	}
	return walk(value)
}

func pageIdentity(raw json.RawMessage) (string, error) {
	var page struct {
		PageID string `json:"pageId"`
		Route  string `json:"route"`
	}
	if err := json.Unmarshal(raw, &page); err != nil {
		return "", err
	}
	identity := strings.TrimSpace(page.PageID)
	if identity == "" {
		identity = strings.TrimSpace(page.Route)
	}
	if identity == "" {
		return "", fmt.Errorf("pageId or route is required")
	}
	return identity, nil
}

func appendNavigation(dst, items []json.RawMessage, moduleID string, owners map[string]string) ([]json.RawMessage, error) {
	for _, item := range items {
		identity, err := navigationIdentity(item)
		if err != nil {
			return nil, fmt.Errorf("manifest: fragment %s navigation: %w", moduleID, err)
		}
		if owner, exists := owners[identity]; exists {
			return nil, fmt.Errorf("manifest: navigation %q conflicts between %s and %s", identity, owner, moduleID)
		}
		owners[identity] = moduleID
		dst = append(dst, item)
	}
	return dst, nil
}

func navigationIdentity(raw json.RawMessage) (string, error) {
	var item struct {
		ID      string `json:"id"`
		PageRef string `json:"pageRef"`
		Label   string `json:"label"`
	}
	if err := json.Unmarshal(raw, &item); err != nil {
		return "", err
	}
	for _, value := range []string{item.ID, item.PageRef, item.Label} {
		if value = strings.TrimSpace(value); value != "" {
			return value, nil
		}
	}
	return "", fmt.Errorf("id, pageRef or label is required")
}
