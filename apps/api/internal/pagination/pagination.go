// Package pagination provides overflow-safe page bounds and SQL OFFSET
// calculations shared by API list surfaces (W8 F-001).
//
// The helpers deliberately accept page/pageSize values that callers have
// already validated as positive. Even when a caller passes an extremely large
// page, the returned bounds never perform the unchecked `(page-1)*pageSize`
// multiplication that could overflow and panic a slice or produce a bogus SQL
// OFFSET.
package pagination

// Bounds returns the half-open [start,end) slice bounds for page/pageSize
// against total. It preserves the existing "beyond the last page returns an
// empty page" behavior: a page after the final page yields (total,total).
//
// total must be >= 0. Non-positive page/pageSize are treated as no items.
func Bounds(page, pageSize, total int) (int, int) {
	if total <= 0 || page < 1 || pageSize < 1 {
		return 0, 0
	}
	// lastPage is floor((total-1)/pageSize)+1, computed without overflowing
	// when total or pageSize is near MaxInt.
	lastPage := int((int64(total)-1)/int64(pageSize) + 1)
	if page > lastPage {
		return total, total
	}
	// page <= lastPage guarantees (page-1)*pageSize < total, so the
	// multiplication cannot overflow.
	start := (page - 1) * pageSize
	end := start + pageSize
	if end > total {
		end = total
	}
	return start, end
}

// Offset returns the safe SQL OFFSET for page/pageSize against total.
// It equals Bounds(...).start, so a page beyond the last page maps to total
// and yields an empty result set rather than a negative/overflowed offset.
func Offset(page, pageSize, total int) int {
	start, _ := Bounds(page, pageSize, total)
	return start
}
