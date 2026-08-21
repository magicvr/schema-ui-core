package pagination

import "testing"

func TestBounds(t *testing.T) {
	tests := []struct {
		name           string
		page, pageSize int
		total          int
		wantStart      int
		wantEnd        int
	}{
		{name: "first page full", page: 1, pageSize: 10, total: 25, wantStart: 0, wantEnd: 10},
		{name: "last partial page", page: 3, pageSize: 10, total: 25, wantStart: 20, wantEnd: 25},
		{name: "one item", page: 1, pageSize: 100, total: 1, wantStart: 0, wantEnd: 1},
		{name: "empty total", page: 1, pageSize: 10, total: 0, wantStart: 0, wantEnd: 0},
		{name: "page beyond last returns empty", page: 4, pageSize: 10, total: 25, wantStart: 25, wantEnd: 25},
		{name: "huge page returns empty safely", page: 92233720368547760, pageSize: 100, total: 25, wantStart: 25, wantEnd: 25},
		{name: "huge pageSize one page", page: 1, pageSize: 9223372036854775807, total: 25, wantStart: 0, wantEnd: 25},
		{name: "huge pageSize page two returns empty", page: 2, pageSize: 9223372036854775807, total: 25, wantStart: 25, wantEnd: 25},
		{name: "max positive page with pageSize one", page: 9223372036854775807, pageSize: 1, total: 10, wantStart: 10, wantEnd: 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, end := Bounds(tt.page, tt.pageSize, tt.total)
			if start != tt.wantStart || end != tt.wantEnd {
				t.Fatalf("Bounds(%d,%d,%d) = (%d,%d), want (%d,%d)",
					tt.page, tt.pageSize, tt.total, start, end, tt.wantStart, tt.wantEnd)
			}
		})
	}
}

func TestOffsetMatchesBoundsStart(t *testing.T) {
	for _, tt := range []struct {
		page, pageSize, total int
	}{
		{1, 20, 100},
		{5, 20, 100},
		{6, 20, 100},
		{92233720368547760, 100, 100},
	} {
		start, _ := Bounds(tt.page, tt.pageSize, tt.total)
		if got := Offset(tt.page, tt.pageSize, tt.total); got != start {
			t.Fatalf("Offset(%d,%d,%d) = %d, want %d", tt.page, tt.pageSize, tt.total, got, start)
		}
	}
}
