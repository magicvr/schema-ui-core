package manifest

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
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

// ForModulesWithFragments is the R4 C3.3 aggregation: baseline projection for
// still-central modules (settings/activity) plus module-contributed manifest
// fragments (users/roles providers). The embedded baseline no longer contains
// the users/roles pages or their admin navigation; those arrive as fragments.
func ForModulesWithFragments(moduleIDs []string, moduleFragments []Fragment) ([]byte, error) {
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
	// R4 C3.3: users/roles are module-contributed fragments; settings/activity
	// remain baseline-projected until C4 migrates them.
	adminModules := map[string]string{
		"settings": "admin.settings",
		"activity": "admin.activity",
	}
	core := base
	core.Pages = []json.RawMessage{}
	core.Navigation = navigation{Top: append([]json.RawMessage(nil), base.Navigation.Top...), Sidebar: []json.RawMessage{}, User: []json.RawMessage{}}
	modulePages := map[string]json.RawMessage{}
	for _, page := range base.Pages {
		var identity struct {
			PageID string `json:"pageId"`
		}
		if err := json.Unmarshal(page, &identity); err != nil {
			return nil, fmt.Errorf("manifest: parse baseline page: %w", err)
		}
		if moduleID, ok := adminModules[identity.PageID]; ok {
			modulePages[moduleID] = page
			continue
		}
		core.Pages = append(core.Pages, page)
	}

	moduleNavigation := map[string]navigation{}
	for _, item := range base.Navigation.Sidebar {
		var group struct {
			Label string            `json:"label"`
			Icon  string            `json:"icon"`
			Items []json.RawMessage `json:"items"`
		}
		if err := json.Unmarshal(item, &group); err != nil {
			return nil, fmt.Errorf("manifest: parse baseline sidebar: %w", err)
		}
		if group.Label != "Admin" || len(group.Items) == 0 {
			core.Navigation.Sidebar = append(core.Navigation.Sidebar, item)
			continue
		}
		for _, navItem := range group.Items {
			var ref struct {
				PageRef string `json:"pageRef"`
			}
			if err := json.Unmarshal(navItem, &ref); err != nil {
				return nil, fmt.Errorf("manifest: parse baseline admin navigation: %w", err)
			}
			moduleID, ok := adminModules[ref.PageRef]
			if !ok {
				return nil, fmt.Errorf("manifest: no module owns admin page %q", ref.PageRef)
			}
			if _, selected := enabled[moduleID]; selected {
				nav := moduleNavigation[moduleID]
				nav.Sidebar = append(nav.Sidebar, navItem)
				moduleNavigation[moduleID] = nav
			}
		}
	}
	for _, item := range base.Navigation.User {
		var ref struct {
			PageRef string `json:"pageRef"`
		}
		if err := json.Unmarshal(item, &ref); err != nil {
			return nil, fmt.Errorf("manifest: parse baseline user navigation: %w", err)
		}
		moduleID, ok := adminModules[ref.PageRef]
		if !ok {
			return nil, fmt.Errorf("manifest: no module owns user page %q", ref.PageRef)
		}
		if _, selected := enabled[moduleID]; selected {
			nav := moduleNavigation[moduleID]
			nav.User = append(nav.User, item)
			moduleNavigation[moduleID] = nav
		}
	}

	coreRaw, err := json.Marshal(core)
	if err != nil {
		return nil, fmt.Errorf("manifest: encode core fragment: %w", err)
	}
	allFragments := []Fragment{{ModuleID: "core.manifest-route", Raw: coreRaw}}
	selectedIDs := make([]string, 0, len(modulePages))
	for moduleID := range modulePages {
		if _, selected := enabled[moduleID]; selected {
			selectedIDs = append(selectedIDs, moduleID)
		}
	}
	sort.Strings(selectedIDs)
	for _, moduleID := range selectedIDs {
		page := modulePages[moduleID]
		moduleFragment := envelope{
			ProtocolVersion:      base.ProtocolVersion,
			RequiredCapabilities: append([]string(nil), base.RequiredCapabilities...),
			App:                  append(json.RawMessage(nil), base.App...),
			Pages:                []json.RawMessage{page},
			Navigation:           moduleNavigation[moduleID],
		}
		raw, err := json.Marshal(moduleFragment)
		if err != nil {
			return nil, fmt.Errorf("manifest: encode module fragment %s: %w", moduleID, err)
		}
		allFragments = append(allFragments, Fragment{ModuleID: moduleID, Raw: raw})
	}
	allFragments = append(allFragments, moduleFragments...)
	return Aggregate(allFragments)
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
