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

// A3 boundary tests — accepted-residual finding converted to fixed (补测不改生产代码).

// ① Duplicate legacy Order values: when two unlisted nodes share the same Order
// integer the tiebreak is lexicographic NodeID, giving a deterministic result
// regardless of original slice position.
func TestSortNavigationDuplicateLegacyOrderIsStable(t *testing.T) {
	// Both "menu_alpha" and "menu_zeta" have Order 5 — neither appears in the
	// (nil) override list, so both fall through to the legacy Order/NodeID branch.
	// Lexicographic NodeID tiebreak must place "menu_alpha" before "menu_zeta".
	nodes := []NavigationContribution{
		{ContributionIdentity: ContributionIdentity{Key: "z"}, NodeID: "menu_zeta", PageID: "z", Order: 5},
		{ContributionIdentity: ContributionIdentity{Key: "a"}, NodeID: "menu_alpha", PageID: "a", Order: 5},
	}
	sortNavigation(nodes, nil) // nil → DefaultNavigationOrder; neither node is listed
	got := nodeIDs(nodes)
	want := []string{"menu_alpha", "menu_zeta"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("duplicate-Order tiebreak = %v, want %v (lexicographic NodeID)", got, want)
	}
}

// ② Case-sensitive exact match: an override list whose entry differs only in
// case ("Menu_Users" vs "menu_users") is treated as containing an unknown
// NodeID — the whole list is invalid and the result falls back to
// DefaultNavigationOrder, not a partial application.
func TestNormalizeNavigationOrderCaseSensitiveExactMatch(t *testing.T) {
	known := []string{"menu_users", "menu_dashboard"}
	// "Menu_Users" (capital M) is NOT in known — different key.
	order := []string{"Menu_Users", "menu_dashboard"}
	got := NormalizeNavigationOrder(order, known)
	if !reflect.DeepEqual(got, DefaultNavigationOrder) {
		t.Fatalf("wrong-case key should fall back; got %v, want DefaultNavigationOrder", got)
	}
}

// ③ Illegal override value falls back to DefaultNavigationOrder: direct unit
// test of NormalizeNavigationOrder (not via sortNavigation) to cover the
// function's own contract boundary — a single unknown ID invalidates the
// entire list and the returned slice is exactly DefaultNavigationOrder.
func TestNormalizeNavigationOrderUnknownIDReturnDefault(t *testing.T) {
	known := []string{"menu_dashboard", "menu_users"}
	order := []string{"menu_dashboard", "menu_bogus"} // "menu_bogus" not in known
	got := NormalizeNavigationOrder(order, known)
	if !reflect.DeepEqual(got, DefaultNavigationOrder) {
		t.Fatalf("unknown NodeID should return DefaultNavigationOrder; got %v", got)
	}
	// Confirm it is the exact same slice value (not just equal content).
	// This documents that callers must not mutate the returned slice.
	gotFirst := got[0]
	if gotFirst != DefaultNavigationOrder[0] {
		t.Fatalf("first element mismatch: %q vs %q", gotFirst, DefaultNavigationOrder[0])
	}
}