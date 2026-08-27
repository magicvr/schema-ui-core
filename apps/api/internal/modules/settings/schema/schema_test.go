package schema

import (
	"encoding/json"
	"testing"
)

// workspace-020 R3 (independent audit A-002 F-001 guard): the Localization
// form's PATCH bodyMapping must cover every field the form renders — the web
// renderer whitelists request bodies via bodyMapping (request-construction),
// so any field missing from the map silently never reaches PATCH.
func TestLocalizationBodyMappingCoversFormFields(t *testing.T) {
	var doc map[string]any
	if err := json.Unmarshal(SchemaDocuments()["settings"], &doc); err != nil {
		t.Fatalf("decode settings.json: %v", err)
	}

	actions, ok := doc["actions"].(map[string]any)
	if !ok {
		t.Fatal("settings.json has no actions")
	}
	updateAction, ok := actions["updateLocalization"].(map[string]any)
	if !ok {
		t.Fatal("updateLocalization action is missing")
	}
	mapping, ok := updateAction["bodyMapping"].(map[string]any)
	if !ok || len(mapping) == 0 {
		t.Fatal("updateLocalization action is missing its bodyMapping")
	}

	// Walk the document tree recursively; collect forms whose props carry
	// submitAction == "updateLocalization" and their field ids.
	var checked int
	var walk func(node any)
	walk = func(node any) {
		obj, ok := node.(map[string]any)
		if !ok {
			return
		}
		props, _ := obj["props"].(map[string]any)
		if submit, _ := props["submitAction"].(string); submit == "updateLocalization" {
			checked++
			fields, _ := props["fields"].([]any)
			if len(fields) == 0 {
				t.Error("localization form declares no fields")
				return
			}
			for _, field := range fields {
				f, _ := field.(map[string]any)
				id, _ := f["id"].(string)
				if id == "" {
					t.Error("localization field without an id")
					continue
				}
				if _, ok := mapping[id]; !ok {
					t.Errorf("localization field %q missing from updateLocalization.bodyMapping (save path drops it)", id)
				}
			}
			return
		}
		if children, ok := obj["children"].([]any); ok {
			for _, child := range children {
				walk(child)
			}
		}
	}
	walk(doc["body"])

	if checked != 1 {
		t.Fatalf("localization form occurrences = %d, want 1", checked)
	}
	// Explicit guard for the R3 addition (defaultCurrency must never regress).
	if mapping["defaultCurrency"] != "defaultCurrency" {
		t.Errorf("bodyMapping[defaultCurrency] = %v, want defaultCurrency", mapping["defaultCurrency"])
	}
}