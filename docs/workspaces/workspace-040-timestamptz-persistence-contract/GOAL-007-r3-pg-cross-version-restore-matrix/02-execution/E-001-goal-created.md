---
id: E-001-goal-created
doc: execution-entry
status: active
parent: GOAL-007-r3-pg-cross-version-restore-matrix
created: 2026-09-21
updated: 2026-09-21
version: 0.1.0
---

# E-001 · R3-C 立项

- **来源**：Root `D-018` §2 第 4 项 + §6 渐进立项策略；编号与 slug 经用户 2026-09-21 预确认（Root `D-019` §2：`GOAL-007-r3-pg-cross-version-restore-matrix`）。
- **前提事实**：R3-A/B（`GOAL-006`）已 `done · 3/3`（independent `A-002` 判 `fail`/3 required → 修复 → `A-004` 定向复审 `pass`/开放 required = 0 → `A-005` 关门记录）；Root 仍 `active · 2/3`。
- **范围**：跨版本 `pg_dump`/`pg_restore` 组合矩阵（逐组合 supported/unsupported 落盘）+ 判据 4 升级后恢复的有界核对 + `I-041-004` 收口。
- **非目标**：不重开 R1/R2 冻结物与 R3-A/B 的 wire 合同；不改 Store/codec/C3 Port；不引入 ORM/第三库；不把容器结果当生产就绪证据（`D-017` §3 约束②）；不做 R3-D。
- **产物**：本目标五件套 + 三个 ledger 目录。
- **进度评估**：`progress: 0/3`；组合定义（检查点 A）尚未冻结。
