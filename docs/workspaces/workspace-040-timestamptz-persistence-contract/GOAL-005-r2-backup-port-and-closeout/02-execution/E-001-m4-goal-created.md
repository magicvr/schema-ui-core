---
id: E-001-m4-goal-created
doc: execution-entry
status: active
parent: GOAL-005-r2-backup-port-and-closeout
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# E-001 · M4 立项

- **来源**：用户 2026-09-20 P-004 裁决 —— 立项 `GOAL-005-r2-backup-port-and-closeout`（AGENTS §11：slug 经用户确认）。
- **前提事实**：`GOAL-003`（M1/M2）`done · 4/4`；`GOAL-004`（M3）`done · 3/3`；`D-021` 的 F-I-005 residual 经 independent 复审判定可 `fixed` 闭合（`GOAL-002/03-audit/A-048` 留痕）；Root `D-017` 关闭 `I-041-003`（常驻 PostgreSQL 15.4 可用于 M3/M4 的可执行验证；破坏性 migration 只可作用于一次性/专用测试 database）。
- **产物**：本目标五件套 + 三个 ledger 目录。
- **进度评估**：`progress: 0/3`（A/B/C 均未开始）；Backup Port 尚无任何代码。
- **待用户裁决**：`I-041-006`（部分升级是否需批级快照/整批回滚）—— 将在 B 检查点前以 P-004 提问。