// GOAL-037 / F-008 根治：wallet_ledger_order_repair（0050）——一次性修复
// 既有库中"同一毫秒内多笔流水"的乱序存档。
//
// 背景：旧版 newID 的 id = UnixMilli + 随机后缀；同一毫秒内连续写入的多条
// 流水，其回放排序（created_at ASC, id ASC）由随机后缀决定，可能 ≠ 实际写入
// 序。回放（checkAccountChain）遇乱序链会把 freeze 排在入账 adjust 之前，
// Apply 反负并报 "replay apply failed: insufficient balance"（reconcile
// inconsistent）。新写入已由同毫秒单调计数修复（provider.newID，E-002）；
// 本迁移对既有数据做一次性重排：按账户回放，在首个失败点定位"同毫秒 id
// 前缀 + 同 created_at 秒"条目组，枚举组内排列取唯一满足链式 Apply 且与
// balance-after 快照一致的顺序，重写该组 id（原 ms 前缀 + 新 seq + 原随机
// 尾）。无唯一合法序 → 迁移失败（fail-closed）；迁移幂等（一次性记账）。
package migration

import (
	"context"
	"fmt"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
)

const (
	orderRepairVersion = 50
	orderRepairKey     = "wallet_ledger_order_repair"
)

type repairEntry struct {
	ID             string
	OriginalID     string
	EntryType      string
	AmountDelta    int64
	AfterTotal     int64
	AfterAvailable int64
	AfterFrozen    int64
	CreatedAt      int64
	MsPrefix       string
	RandomSuffix   string
}

// applyStep replays one entry against prev balances (mirrors store.Apply for
// the repair path: validates non-negativity and the total invariant).
func applyStep(entryType string, delta int64, total, avail, frozen int64) (int64, int64, int64, bool) {
	switch entryType {
	case "adjust":
		if delta == 0 {
			return 0, 0, 0, false
		}
		total += delta
		avail += delta
	case "freeze":
		if delta <= 0 {
			return 0, 0, 0, false
		}
		avail -= delta
		frozen += delta
	case "unfreeze":
		if delta <= 0 {
			return 0, 0, 0, false
		}
		avail += delta
		frozen -= delta
	case "deduct_frozen":
		if delta <= 0 {
			return 0, 0, 0, false
		}
		total -= delta
		frozen -= delta
	default:
		return 0, 0, 0, false
	}
	if total < 0 || avail < 0 || frozen < 0 || total != avail+frozen {
		return 0, 0, 0, false
	}
	return total, avail, frozen, true
}

// validOrder reports whether entries in this order from the given balance
// state keep every step valid and matching its row's balance-after snapshot.
func validOrder(total, avail, frozen int64, entries []repairEntry) (int64, int64, int64, bool) {
	for _, e := range entries {
		nt, na, nf, ok := applyStep(e.EntryType, e.AmountDelta, total, avail, frozen)
		if !ok || nt != e.AfterTotal || na != e.AfterAvailable || nf != e.AfterFrozen {
			return 0, 0, 0, false
		}
		total, avail, frozen = nt, na, nf
	}
	return total, avail, frozen, true
}

// permutations returns every index ordering of [0..n).
func permutations(n int) [][]int {
	if n <= 1 {
		return [][]int{{0}}
	}
	out := [][]int{}
	var gen func(prefix, remaining []int)
	gen = func(prefix, remaining []int) {
		if len(remaining) == 0 {
			cp := make([]int, len(prefix))
			copy(cp, prefix)
			out = append(out, cp)
			return
		}
		for i, r := range remaining {
			rest := make([]int, 0, len(remaining)-1)
			rest = append(rest, remaining[:i]...)
			rest = append(rest, remaining[i+1:]...)
			next := make([]int, len(prefix)+1)
			copy(next, prefix)
			next[len(prefix)] = r
			gen(next, rest)
		}
	}
	idx := make([]int, n)
	for i := range idx {
		idx[i] = i
	}
	gen(nil, idx)
	return out
}

// repairGroup rewrites one same-millisecond group into its unique valid
// order. Returns the reordered entries (with rewritten ids) and the balance
// state at the group tail. Fails closed on ambiguity or no solution.
func repairGroup(tx kernel.Tx, accountID string, prevTotal, prevAvail, prevFrozen int64, group []repairEntry) ([]repairEntry, int64, int64, int64, error) {
	if len(group) < 2 || len(group) > 6 {
		return nil, 0, 0, 0, fmt.Errorf("wallet order repair: account %s: unrepairable group of %d same-millisecond entries", accountID, len(group))
	}
	var best []int
	var tailTotal, tailAvail, tailFrozen int64
	validCount := 0
	for _, perm := range permutations(len(group)) {
		cand := make([]repairEntry, 0, len(group))
		for _, idx := range perm {
			cand = append(cand, group[idx])
		}
		gt, ga, gf, ok := validOrder(prevTotal, prevAvail, prevFrozen, cand)
		if !ok {
			continue
		}
		best = perm
		tailTotal, tailAvail, tailFrozen = gt, ga, gf
		validCount++
		if validCount > 1 {
			return nil, 0, 0, 0, fmt.Errorf("wallet order repair: account %s: ambiguous order (%d) for same-millisecond group at id %s", accountID, validCount, group[0].ID)
		}
	}
	if validCount != 1 || best == nil {
		return nil, 0, 0, 0, fmt.Errorf("wallet order repair: account %s: no valid reorder for same-millisecond group at id %s", accountID, group[0].ID)
	}
	ordered := make([]repairEntry, 0, len(group))
	for pos, idx := range best {
		e := group[idx]
		newID := fmt.Sprintf("%s%08x%s", e.MsPrefix, pos, e.RandomSuffix)
		if newID != e.ID {
			if _, err := tx.Exec(context.Background(),
				`UPDATE wallet_ledger_entries SET id = ? WHERE id = ? AND account_id = ?`,
				newID, e.OriginalID, accountID,
			); err != nil {
				return nil, 0, 0, 0, fmt.Errorf("wallet order repair: rewrite id %s: %w", e.OriginalID, err)
			}
			e.ID = newID
		}
		ordered = append(ordered, e)
	}
	return ordered, tailTotal, tailAvail, tailFrozen, nil
}

type ledgerRow struct {
	accountID string
	entry     repairEntry
}

// orderRepair replays every account's ledger in (created_at, id) order and
// repairs the first same-millisecond group at every replay failure, continuing
// until the chain replays clean. Runs inside the caller's transaction.
func orderRepair(tx kernel.Tx) error {
	rows, err := tx.Query(context.Background(),
		`SELECT account_id, id, entry_type, amount_delta, balance_after_total, balance_after_available, balance_after_frozen, created_at
		 FROM wallet_ledger_entries ORDER BY account_id, created_at, id`,
	)
	if err != nil {
		return fmt.Errorf("wallet order repair: read ledger: %w", err)
	}
	var all []ledgerRow
	for rows.Next() {
		var r ledgerRow
		e := &r.entry
		var createdAt int64
		if err := rows.Scan(&r.accountID, &e.ID, &e.EntryType, &e.AmountDelta, &e.AfterTotal, &e.AfterAvailable, &e.AfterFrozen, &createdAt); err != nil {
			_ = rows.Close()
			return fmt.Errorf("wallet order repair: scan ledger: %w", err)
		}
		e.CreatedAt = createdAt
		e.OriginalID = e.ID
		if len(e.ID) >= 16 {
			e.MsPrefix = e.ID[:16]
			e.RandomSuffix = e.ID[16:]
		}
		all = append(all, r)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("wallet order repair: iterate ledger: %w", err)
	}
	if err := rows.Close(); err != nil {
		return fmt.Errorf("wallet order repair: close rows: %w", err)
	}

	repaired := 0
	i := 0
	for i < len(all) {
		accountID := all[i].accountID
		total, avail, frozen := int64(0), int64(0), int64(0)
		for i < len(all) && all[i].accountID == accountID {
			e := all[i].entry
			nt, na, nf, ok := applyStep(e.EntryType, e.AmountDelta, total, avail, frozen)
			if ok && nt == e.AfterTotal && na == e.AfterAvailable && nf == e.AfterFrozen {
				total, avail, frozen = nt, na, nf
				i++
				continue
			}
			if e.MsPrefix == "" {
				return fmt.Errorf("wallet order repair: account %s: entry %s has no millisecond prefix", accountID, e.ID)
			}
			// Same-millisecond group: contiguous entries with the same id ms
			// prefix AND the same created_at second as the failing entry.
			j := i + 1
			for j < len(all) && all[j].accountID == accountID &&
				all[j].entry.MsPrefix == e.MsPrefix &&
				all[j].entry.CreatedAt == e.CreatedAt {
				j++
			}
			if j-i < 2 {
				// Single entry failing with no same-ms sibling: data is
				// genuinely broken, not an ordering issue.
				return fmt.Errorf("wallet order repair: account %s: entry %s cannot replay (not an ordering problem)", accountID, e.ID)
			}
			group := make([]repairEntry, 0, j-i)
			for k := i; k < j; k++ {
				group = append(group, all[k].entry)
			}
			ordered, gt, ga, gf, err := repairGroup(tx, accountID, total, avail, frozen, group)
			if err != nil {
				return err
			}
			for k := i; k < j; k++ {
				all[k].entry = ordered[k-i]
			}
			total, avail, frozen = gt, ga, gf
			repaired++
			i = j
		}
	}
	return nil
}

// migrateOrderRepair applies the 0050 repair (sqlite and postgres share the
// same Go logic; no dialect SQL involved).
func migrateOrderRepair(tx kernel.Tx) error {
	return orderRepair(tx)
}

// migrateOrderRepairPG is the postgres variant (identical logic).
func migrateOrderRepairPG(tx kernel.Tx) error {
	return orderRepair(tx)
}