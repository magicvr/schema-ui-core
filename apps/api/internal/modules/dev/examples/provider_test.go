package examples

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
)

// W1 (GOAL-002 / workspace-010): dev.examples owns the 8 demonstration pages.
func TestProviderPublishesExampleSchemaDocuments(t *testing.T) {
	provider := New()
	descriptor := provider.Descriptor()
	plan := kernel.Plan{
		Modules:      []kernel.Module{descriptor},
		Capabilities: []kernel.Capability{kernel.CapabilitySchema, kernel.CapabilityValidation, kernel.CapabilityNavigation, kernel.CapabilityManifest},
	}
	set, err := kernel.RegisterContributions(context.Background(), plan, []kernel.Provider{provider})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"admin-list-batch", "data-display", "data-table", "form-controls", "form-with-reactions", "form-with-upload", "overview", "search-form-table"}
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
	if len(set.Fragments) != 1 || set.Fragments[0].FragmentID != "examples" {
		t.Fatalf("fragments = %+v, want one examples fragment", set.Fragments)
	}
}

func TestProviderRegisterNoRoutesPermissionsOrNavigation(t *testing.T) {
	descriptor := New().Descriptor()
	if len(descriptor.Contributions.Routes) != 0 || len(descriptor.Contributions.Permissions) != 0 || len(descriptor.Contributions.Navigation) != 0 {
		t.Fatalf("dev.examples must be a horizontal demo module without routes/permissions/system-data nav: %+v", descriptor.Contributions)
	}
}
