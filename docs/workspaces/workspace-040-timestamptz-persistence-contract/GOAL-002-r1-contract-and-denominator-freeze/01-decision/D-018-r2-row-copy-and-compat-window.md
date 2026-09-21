---
id: D-018-r2-row-copy-and-compat-window
doc: decision-entry
status: accepted
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# D-018 · R2 行拷贝机制与扫描器兼容窗口（用户 P-004 裁决）

## 决定的来源

- A-030 **F-I-002.1** 要求「写出每张受影响表的 exact SQLite rebuild DDL」；A-029 原文要求「copy rows through codec」。
- 2026-09-20 调查发现两处决定性事实（详见 `attachments/r1-c2-sqlite-rebuild-mechanism-v1.0-fc.md`）：
  1. 本项目 DSN 固定 `_foreign_keys=on` 且 `SetMaxOpenConns(1)`，`PRAGMA foreign_keys` 在事务内为 no-op，**经典 rebuild 路径不可用**；但仓库既有 rebuild（wallet / operationlog）**不用任何 pragma** 且已在生产路径运行，含被 FK 引用的表。
  2. SQLite `strftime('%f', …)` **round**（实测 `1758320000.9999` → 下一秒）且输出**变宽**（整秒仅 3 位小数），与 fixed-6 合同和 Root `D-015` 的截断规则冲突。
- 用户 2026-09-20 经 P-004 裁决两项。

## 决定

### 1. 行拷贝机制 = 纯整数 SQL 表达式 + 强制 Go codec 对拍（用户选项 C）

- 逐表 `INSERT … SELECT` 的每个时间列使用 `r1-c2-sqlite-rebuild-mechanism-v1.0-fc.md` §2 的两条**纯整数**表达式（秒族 `strftime('%Y-%m-%dT%H:%M:%S', <col>, 'unixepoch') || '.000000Z'`；毫秒族用 floor 除法 + `printf('%03d', …)` + `'000Z'`）。
- **同时**把「Go codec `FromUnix` / `FromUnixMilli` 与 SQL 结果逐行等价比对」写成 **R2 强制验收步骤**（`T-<#>-RT` 用例族）；任一不等价即 fail closed，不得放行。
- 全包**禁用** `strftime('%f', …)` 与 `/1000.0`（round + 变宽）。

### 2. 扫描器兼容窗口 = 无过渡期（用户选项 A）

- 列形状（`INTEGER` → `TEXT`）变更与仓储层扫描/绑定（`int64` → `time.Time`/`sql.NullTime`）改造在**同一发布**内完成。
- 不接受「SQLite 列已是 TEXT、旧 int64 扫描器仍在读」的并存状态；不需要双向临时解析。

## 理由

- 纯整数 SQL 使转换逻辑完全落在 `MigrationChecksum` 所覆盖的 `stmts` 切片内，与 `D-017` 的单 checksum 约定自洽，且单条 SQL 便于逐表核对。
- Go 对拍同时满足 A-029「经 codec 语义」的方向要求，且把等价性变成可执行门禁而非散文承诺。
- 无过渡期与仓库既有迁移事务语义一致（`applyMigration` 单事务 + `Apply(sqlTx)`），避免引入两套解析路径。

> **2026-09-20 更正（响应 A-032）**：本决策初稿在「理由」中曾以「仓库既有 rename 重建**含被 FK 引用的表**也能工作」作为机制可行性依据。该依据**经独立审计实测否定**（裸 rename 重建会永久改写子表 `REFERENCES`，见 `r1-c2-fk-parent-rebuild-finding-v1.0-fc.md` §F-1～F-5，并由 A-032 复现）。**该理由句作废**；本决策的**选项 C 本身不变**（行拷贝仍用纯整数 SQL + 强制 Go 对拍），只是其可行性不再依赖那句错误论证，而依赖 F-5 子女先行模式（待 P-004 定案，见 **F-I-021**）。

## 未选方案

- **B：Go codec 逐行拷贝**（严格按 A-029 字面）。转换逻辑在 Go 代码内、不进 checksum 覆盖的 `stmts`；与仓库既有 rebuild 的「单条 `INSERT SELECT`」模式分叉；且需先落 `internal/temporal` 包（属 R2 产物）。用户未选。
- **兼容窗口 B：允许过渡期双读**。需要额外双向解析与第二套测试矩阵，且延长 TEXT/int64 并存的缺陷面。用户未选。

## 影响与边界

- **已知偏差（须明示）**：A-029 字面要求「copy rows through codec」；本决定下 codec **不在拷贝路径上**，而是作为独立对拍验收。该偏差由用户书面裁决接受；independent 复审如有异议须按 P-004 回到用户。
- 本决定只确定**机制**；**不等于** exact rebuild DDL 已落盘或 F-I-002 已闭合——逐表正文仍在冻结候选阶段。
- 本决定不改变 `D-017`（单 checksum over SQLite DDL 切片）与 append-only 边界。
