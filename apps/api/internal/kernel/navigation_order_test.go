package kernel

import (
	"reflect"
	"testing"
)

// GOAL-013 D-002 §2: the product-frozen default navigation order. This test
// locks the exact NodeID sequence — maintainers must update both the constant
// and this snapshot when a new module lands (I-003 closed).
func TestDefaultNavigationOrderSnapshot(t *testing.T) {
	want := []string{
		"menu_dashboard",
		"menu_users",
		"menu_roles",
		// S-14 (GOAL-019, user 2026-08-16): wallet directly below Roles.
		"menu_wallet",
		"menu_account",
		// GOAL-022 (D-002 §2): my-wallet self-service — user slot between
		// 个人中心 and 设置.
		"menu_wallet_self",
		"menu_activity",
		"menu_settings",
		"menu_notifications",
		"menu_files",
		"menu_dictionary",
		"menu_monitoring",
		"menu_scheduled_tasks",
		"menu_recycle_bin",
		"menu_data_permission",
	}
	if !reflect.DeepEqual(DefaultNavigationOrder, want) {
		t.Fatalf("DefaultNavigationOrder = %v, want %v", DefaultNavigationOrder, want)
	}
}

func navNodes() []NavigationContribution {
	return []NavigationContribution{
		{ContributionIdentity: ContributionIdentity{Key: "n1"}, NodeID: "menu_users", PageID: "users", Order: 1},
		{ContributionIdentity: ContributionIdentity{Key: "n2"}, NodeID: "menu_dashboard", PageID: "dashboard", Order: 0},
		{ContributionIdentity: ContributionIdentity{Key: "n3"}, NodeID: "menu_newbie", PageID: "newbie", Order: 99},
		{ContributionIdentity: ContributionIdentity{Key: "n4"}, NodeID: "menu_roles", PageID: "roles", Order: 2},
	}
}

func nodeIDs(nodes []NavigationContribution) []string {
	ids := make([]string, 0, len(nodes))
	for _, n := range nodes {
		ids = append(ids, n.NodeID)
	}
	return ids
}

// The default list wins over legacy Order values; unknown nodes append at the
// end in legacy Order/NodeID order.
func TestSortNavigationDefaultOrder(t *testing.T) {
	nodes := navNodes()
	sortNavigation(nodes, nil)
	want := []string{"menu_dashboard", "menu_users", "menu_roles", "menu_newbie"}
	if got := nodeIDs(nodes); !reflect.DeepEqual(got, want) {
		t.Fatalf("order = %v, want %v", got, want)
	}
}

// A full custom list is honored exactly as given.
func TestSortNavigationCustomOrder(t *testing.T) {
	nodes := navNodes()
	order := []string{"menu_roles", "menu_dashboard", "menu_users", "menu_newbie"}
	sortNavigation(nodes, order)
	want := []string{"menu_roles", "menu_dashboard", "menu_users", "menu_newbie"}
	if got := nodeIDs(nodes); !reflect.DeepEqual(got, want) {
		t.Fatalf("order = %v, want %v", got, want)
	}
}

// Unknown NodeIDs in a provided order invalidate the whole list: fall back to
// the default order (with a warning), never a partial application.
func TestSortNavigationInvalidOrderFallsBackToDefault(t *testing.T) {
	nodes := navNodes()
	order := []string{"menu_users", "menu_bogus"}
	sortNavigation(nodes, order)
	want := []string{"menu_dashboard", "menu_users", "menu_roles", "menu_newbie"}
	if got := nodeIDs(nodes); !reflect.DeepEqual(got, want) {
		t.Fatalf("order = %v, want default %v", got, want)
	}
}

// A custom order need not cover every node: unlisted nodes append at the end
// (legacy Order/NodeID fallback) so new modules never disappear.
func TestSortNavigationPartialCustomOrderAppendsRest(t *testing.T) {
	nodes := navNodes()
	order := []string{"menu_roles", "menu_dashboard"}
	sortNavigation(nodes, order)
	want := []string{"menu_roles", "menu_dashboard", "menu_users", "menu_newbie"}
	if got := nodeIDs(nodes); !reflect.DeepEqual(got, want) {
		t.Fatalf("order = %v, want %v", got, want)
	}
}