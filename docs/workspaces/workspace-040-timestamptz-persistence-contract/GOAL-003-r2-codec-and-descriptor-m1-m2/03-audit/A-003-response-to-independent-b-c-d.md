---
id: A-003-response-to-independent-b-c-d
doc: audit-entry
status: active
parent: GOAL-003-r2-codec-and-descriptor-m1-m2
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
source: self
verdict: pass
---

# A-003（self · 编排器响应）· GOAL-003 检查点 B/C/D

- **source**: self（编排器汇总响应）
- **日期**: 2026-09-20
- **scope**: 响应本目标 A-002（independent · grok-build grok-4.6 · high · `conditional`）的全部意见
- **verdict**: `pass`（A-002 的 open required 已按 `fixed` 闭合；无未闭合必改项）

## 必改项闭合

| finding | 级别 | 处置 | 证据 |
|---------|------|------|------|
| A-002 `F-I-001`（v73 缺「retired `records` 不存在」断言） | required | **fixed** | v73 新增 m0 guard（SQLite：`SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='records'`；PG：`information_schema.tables` 同类），`RunGuards` / `RunPostgresGuards` 在 Apply 最前 fail closed；v73 checksum 重算为 `4dd07092330cb3b344143f49ec57edd1f89e92cf2786446c0635c4a320eb8a9f`，冻结表与 `attachments/r2-v73-v87-generated-statements-v0.1.md` 同步。**可执行证据**：`internal/w040contracttest/migration_boundaries_test.go::TestV73RefusesWhenRetiredRecordsTableIsPresent`（人造 `records` → v73 拒绝、v73 事务回滚、`records.updated_at`/`schema_migrations.applied_at` 仍为 INTEGER） |

## recommended 项处置

| finding | 处置 | 证据 / 说明 |
|---------|------|-------------|
| `F-I-002`（m4 的 `PRAGMA foreign_key_check` / `integrity_check` 在哈希里每表一次、执行只一次） | **fixed** | `VerifyStatements` 改为**每 descriptor 一次**，与 `RunVerify` 完全同构；11 个多表 descriptor 的 checksum 随之重算并同步冻结表（4 个单表 descriptor 的 checksum 不变，恰好印证改动的局部性） |
| `F-I-003`（D-018 的 `T-#-RT` 是抽样硬编码，缺同进程 codec↔SQL 对拍） | **fixed** | 新增 `TestCodecMatchesRealMigration`：秒族 7 值 + 毫秒族 14 值（含 ±1/±999/±1000/±1001、1758320000123、-1758320000123、公元 9999、公元 1 下界）在**同一次测试**里断言 `temporal.MustFormat(temporal.FromUnix/FromUnixMilli(v))` 等于真实 v73–v87 迁移后的扫描值 |
| `F-I-004`（生成器留生产树） | **保持 + 加固说明** | 有意保留：重跑生成器对同一份 v1–v72 历史**字节级复现** 15 个 descriptor + 清单（16/16 SHA-256 相同，已在 `02-execution.md` E-003 记录）。默认 skip（须显式 `VP040_GENERATE=1`），文件头已写明其为工具而非普通测试；risk 是「误开环境变量」，但生成器只写 `vp040_temporal.go` 与清单附件，且任何改动都会被冻结 checksum 测试立即抓住 |
| `F-I-005`（v86 `ModuleID` 台账展示列） | **记录更正，不重写已冻结附件** | 该列不在台账 §4 的「R2 唯一允许输入」内，且运行时 `ModuleID` 由 provider 机械决定。**不**改动 GOAL-002 已关门附件，改为在本响应 + `E-003` + 语句清单附件中显式标注更正（避免重写已冻结载体，同时消除对照漂移）。若后续需要单点权威，由新决策回写该展示列 |
| `F-I-006`（`00-meta` 仍写 I-041-001 collecting） | **fixed** | `00-meta.md` 与 `03-audit.md` 的信息就绪表同步为 `verified` |

## `D-021` residual 复审（本次触发）

- 触发条件满足：R2 首次记录 v73+ 真实哈希（本批 15 个）。
- A-002 判定 residual ①②③ 均有可核对实现证据、可走 `fixed` 闭合；编排器接受该复审结论，并在 **GOAL-002（residual 归属目标）** 的审计台账留下闭合记录（不修改 GOAL-002 的 `D-021` 决策正文字面）。
- **明确边界**：residual 闭合**不**等于 R2 放行，**不**等于 M4（Backup Port）完成。

## 未闭合项

无。A-002 的 open required = 1 已 `fixed`；其余为 recommended 且均已处置或记录。GOAL-003 检查点 D 据此可判定完成（canonical SQL + 真实 checksum 已落盘、`D-021` 复审已落盘）。
