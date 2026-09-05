package schema

import (
	"encoding/json"
	"testing"
)

func TestPageDocumentsExposeSeparateSettingsAndOperatorSurfaces(t *testing.T) {
	documents := SchemaDocuments()
	if len(documents) != 2 {
		t.Fatalf("expected two Telegram page documents, got %d", len(documents))
	}

	for _, pageID := range PageIDs() {
		document, ok := documents[pageID]
		if !ok {
			t.Fatalf("missing document for page %q", pageID)
		}
		var page struct {
			Meta struct {
				PageID string `json:"pageId"`
			} `json:"meta"`
		}
		if err := json.Unmarshal(document, &page); err != nil {
			t.Fatalf("unmarshal page %q: %v", pageID, err)
		}
		if page.Meta.PageID != pageID {
			t.Fatalf("page %q has meta.pageId %q", pageID, page.Meta.PageID)
		}
	}

	var settings struct {
		Actions map[string]struct {
			Type string `json:"type"`
			URL  string `json:"url"`
		} `json:"actions"`
	}
	if err := json.Unmarshal(documents["telegram-settings"], &settings); err != nil {
		t.Fatalf("unmarshal settings page: %v", err)
	}
	if action := settings.Actions["openTelegramOperator"]; action.Type != "navigate" || action.URL != "/telegram-settings/operator" {
		t.Fatalf("unexpected operator entry action: %+v", action)
	}

	var operator struct {
		Body struct {
			Type      string `json:"type"`
			Component string `json:"component"`
			Props     struct {
				Surface string `json:"surface"`
			} `json:"props"`
		} `json:"body"`
	}
	if err := json.Unmarshal(documents["telegram-operator"], &operator); err != nil {
		t.Fatalf("unmarshal operator page: %v", err)
	}
	if operator.Body.Type != "custom" || operator.Body.Component != "telegram-admin-tab" || operator.Body.Props.Surface != "operator" {
		t.Fatalf("operator page must mount the operator surface directly: %+v", operator.Body)
	}
}
