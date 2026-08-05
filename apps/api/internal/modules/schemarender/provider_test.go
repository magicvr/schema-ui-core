package schemarender

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
)

func TestProviderPublishesCoreSchemaDocuments(t *testing.T) {
	provider := New()
	descriptor := provider.Descriptor()
	plan := kernel.Plan{
		Modules:      []kernel.Module{descriptor},
		Capabilities: []kernel.Capability{kernel.CapabilitySchema, kernel.CapabilityValidation},
	}
	set, err := kernel.RegisterContributions(context.Background(), plan, []kernel.Provider{provider})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"data-table", "form-controls", "form-with-reactions", "overview", "search-form-table"}
	got := make([]string, 0, len(set.Pages))
	for _, page := range set.Pages {
		got = append(got, page.PageID)
		var document struct {
			Meta struct {
				PageID string `json:"pageId"`
			} `json:"meta"`
		}
		if err := json.Unmarshal(page.Document, &document); err != nil {
			t.Fatalf("%s document: %v", page.PageID, err)
		}
		if page.Owner != ModuleID || document.Meta.PageID != page.PageID {
			t.Fatalf("page = %+v meta.pageId=%q", page, document.Meta.PageID)
		}
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("pages = %v, want %v", got, want)
	}
}
