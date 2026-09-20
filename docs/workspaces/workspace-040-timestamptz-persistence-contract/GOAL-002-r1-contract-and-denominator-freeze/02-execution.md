---
id: GOAL-002-r1-contract-and-denominator-freeze
doc: execution
status: active
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# 执行记录 · GOAL-002

## 执行索引

| E-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|--------|------|
| E-001 | 2026-09-20 | R1 用户方案裁决已落盘 | recorded | `02-execution/E-001-user-decisions-recorded.md` |
| E-002 | 2026-09-20 | 时间列与单位 inventory v0.1 | recorded | `02-execution/E-002-time-column-inventory-v0.1.md` |
| E-003 | 2026-09-20 | 完整时间列 inventory v0.2 | recorded | `02-execution/E-003-inventory-complete-v0.2.md` |
| E-004 | 2026-09-20 | 公共时间输出合同用户裁决 | recorded | `02-execution/E-004-public-wire-decision-recorded.md` |
| E-005 | 2026-09-20 | compiled catalog v72 口径纠正 | recorded | `02-execution/E-005-catalog-72-correction.md` |
| E-006 | 2026-09-20 | wire 输入兼容用户裁决 | recorded | `02-execution/E-006-wire-input-compat-decision.md` |
| E-007 | 2026-09-20 | 公共时间 wire inventory v0.1 | recorded | `02-execution/E-007-public-wire-inventory-v0.1.md` |
| E-008 | 2026-09-20 | Backup SPI/Service 用户裁决 | recorded | `02-execution/E-008-backup-spi-user-decision.md` |
| E-009 | 2026-09-20 | C2/C3 guardrails 草案 | recorded | `02-execution/E-009-c2-c3-guardrails-draft.md` |
| E-010 | 2026-09-20 | Backup/RecoveryPoint Port 公共边界裁决 | recorded | `02-execution/E-010-backup-port-surface-user-decision.md` |
| E-011 | 2026-09-20 | C2 精度、D0 与 PG backup provider 用户裁决 | recorded | `02-execution/E-011-c2-precision-zero-backup-decisions.md` |
| E-012 | 2026-09-20 | wire 非 DB 例外范围用户裁决 | recorded | `02-execution/E-012-wire-nondb-exceptions-decision.md` |
| E-013 | 2026-09-20 | Backup Port 最小方法裁决 | recorded | `02-execution/E-013-backup-port-methods-decision.md` |
| E-014 | 2026-09-20 | schema_migrations owner 用户裁决 | recorded | `02-execution/E-014-schema-ledger-owner.md` |
| E-015 | 2026-09-20 | minimal Backup/RecoveryPoint Port contract 草案 | recorded | `02-execution/E-015-backup-port-contract-draft.md` |
| E-016 | 2026-09-20 | voucher / monotonic policy 用户裁决 | recorded | `02-execution/E-016-voucher-monotonic-policy-decisions.md` |
| E-017 | 2026-09-20 | C2 90 列 codec/NULL mapping matrix 草案 | recorded | `02-execution/E-017-c2-column-matrix-draft.md` |
| E-018 | 2026-09-20 | predicate/index/check 矩阵草案 | recorded | `02-execution/E-018-predicate-index-matrix-draft.md` |
| E-019 | 2026-09-20 | v73+ module conversion owner allocation draft | recorded | `02-execution/E-019-v73-owner-allocation-draft.md` |
| E-020 | 2026-09-20 | module-owned migration specification draft | recorded | `02-execution/E-020-owner-migration-spec-draft.md` |
| E-021 | 2026-09-20 | runtime read/write and predicate spec draft | recorded | `02-execution/E-021-readwrite-predicate-spec-draft.md` |
| E-022 | 2026-09-20 | backup/restore-to-new-db runbook draft | recorded | `02-execution/E-022-backup-restore-runbook-draft.md` |
| E-023 | 2026-09-20 | v73+ append-only test rewrite checklist draft | recorded | `02-execution/E-023-v73-test-rewrite-checklist-draft.md` |
| E-024 | 2026-09-20 | v73 allocation / negative truncation policy | recorded | `02-execution/E-024-v73-negative-truncation-decision.md` |
| E-025 | 2026-09-20 | C2 冻结候选第一批：90 列逐列转换合同 + 谓词 exact SQL 单表 + descriptor 台账 | recorded | `02-execution/E-025-c2-freeze-candidate-batch1.md` |
| E-026 | 2026-09-20 | D-015 秒/毫秒族措辞收口与 D-012/D-013 编号限定（裁决 B） | recorded | `02-execution/E-026-d015-seconds-family-and-d012-scoping.md` |
| E-027 | 2026-09-20 | A-030 独立审计响应与缺陷修正（含 `#5` 锁谓词方向改正） | recorded | `02-execution/E-027-a030-response-and-defect-fixes.md` |
| E-028 | 2026-09-20 | 逐表 rebuild DDL 前置调查：发现 FK 父表重建阻塞 | recorded | `02-execution/E-028-fk-parent-rebuild-blocker.md` |
| E-029 | 2026-09-20 | A-032 响应与 FK 重建模式定案（P-004 裁决 F-5） | recorded | `02-execution/E-029-fk-rebuild-mode-decided.md` |
| E-030 | 2026-09-20 | 逐表 exact rebuild DDL 正文落盘（含 v74 F-5 编排整链实测） | recorded | `02-execution/E-030-per-table-rebuild-ddl-drafted.md` |
| E-031 | 2026-09-20 | A-034 响应：补 v78 与 site_settings DDL，起草 PG 显式 DDL | recorded | `02-execution/E-031-v78-and-site-settings-ddl-added.md` |
| E-032 | 2026-09-20 | A-036 响应：四项闭合、PG v78 与 dict_entries 补齐 | recorded | `02-execution/E-032-a036-response.md` |
| E-033 | 2026-09-20 | A-038 响应与 C3 备份/回滚边界首次落盘 | recorded | `02-execution/E-033-c3-boundary-drafted.md` |
| E-034 | 2026-09-20 | F-I-002 可执行边界测试落地（D-020）与 A-040 响应 | recorded | `02-execution/E-034-fi002-executable-boundary-tests.md` |

> 索引行序自本版起为 **E-001 → E-034 严格递增**（A-029/A-030 **F-I-019** 状态维持 closed）。
>
> **Git checkpoint**：本轮（E-025/E-026）落盘后提交 `2d0734a1` — `govern(workspace-040): C2 freeze candidate batch 1 (90-column contract, predicate SQL, descriptor ledger)`。只暂存显式 owned paths（15 个文件，全部位于 workspace-040 与 Root `D-015`），未使用 `git add -A`；`apps/**` 无变更。commit hash 不作为审计或验收证据，仅用于可追溯。

## 事实边界

> 方案、inventory、转换设计与审计事实必须在发生后追加；不得把用户裁决或计划写成已实施迁移。
