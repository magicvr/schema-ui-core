package manifest

import (
	"encoding/json"
	"testing"
)

func TestFragmentKeepsOperatorPageOutOfSidebar(t *testing.T) {
	var fragment struct {
		Pages []struct {
			PageID string `json:"pageId"`
			Route  string `json:"route"`
		} `json:"pages"`
		Navigation struct {
			Sidebar []struct {
				PageRef string `json:"pageRef"`
			} `json:"sidebar"`
		} `json:"navigation"`
	}
	if err := json.Unmarshal(FragmentJSON, &fragment); err != nil {
		t.Fatalf("unmarshal Telegram fragment: %v", err)
	}

	if len(fragment.Pages) != 2 {
		t.Fatalf("expected settings and operator pages, got %+v", fragment.Pages)
	}
	pageRoutes := map[string]string{}
	for _, page := range fragment.Pages {
		pageRoutes[page.PageID] = page.Route
	}
	if pageRoutes["telegram-settings"] != "/telegram-settings" {
		t.Fatalf("unexpected settings route: %+v", pageRoutes)
	}
	if pageRoutes["telegram-operator"] != "/telegram-settings/operator" {
		t.Fatalf("unexpected operator route: %+v", pageRoutes)
	}
	if len(fragment.Navigation.Sidebar) != 1 || fragment.Navigation.Sidebar[0].PageRef != "telegram-settings" {
		t.Fatalf("operator page must not be a sidebar item: %+v", fragment.Navigation.Sidebar)
	}
}
