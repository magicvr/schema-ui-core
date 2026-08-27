---
id: E-003
title: R2 测试验证 + 关门（GOAL-003 done 3/3）
date: 2026-08-27
status: done
---

# E-003 · R2 关门（2026-08-27）

## 事实

1. **全量回归**：`go test ./...`（apps/api）**全绿**（exit 0，全部包 ok；含 config 新增 7 子测）。
2. **自审**：A-001 self `pass`（0 required；F-001 recommended → R3 承接进程级 harness）。
3. **关门**：GOAL-003 `done · 3/3`（C1 方案冻结 → C2 实施 → C3 测试与关门）。
4. **Root progress 校正留痕**：进度口径 = 已关门纲领阶段 / 3——R1→1/3、R2→**2/3**、R3→3/3。R1 关门时 Root/goal-tree 曾误记 2/3（当时应为 1/3）；R2 关门后现值 2/3 与实际一致；历史快照不回改，此条供后来读者校正。

## 验证 / 后续

- R3（GOAL-004）：合同 §8 harness 进程级核对——clean drain exit 0 / timeout exit 1 / 重启 reclaim（A-001 F-001 承接）；SQLite + PG 双方言；compose stop 路径核对。