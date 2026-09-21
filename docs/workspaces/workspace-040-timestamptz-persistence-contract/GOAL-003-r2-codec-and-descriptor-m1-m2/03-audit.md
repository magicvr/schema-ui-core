---
id: GOAL-003-r2-codec-and-descriptor-m1-m2
doc: audit
status: active
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-20
updated: 2026-09-20
version: 0.3.0
---

# 审计 · GOAL-003

## 信息就绪核对

| 核对项 | 状态 | 备注 |
|--------|------|------|
| I-041-001 | verified | codec API 已定稿并落码（`01-decision/D-001-temporal-codec-api.md`；11 个单测全绿）。A-002 F-I-006 曾指出 `00-meta` 未同步，已修正 |
| I-040-001 / I-040-003 | verified | R1 关门时已 verified（见 GOAL-002 03-audit） |

## 意见台账索引

| A-ID | 日期 | source | scope | verdict | 开放 required | 文件 |
|------|------|--------|-------|---------|---------------|------|
| A-001 | 2026-09-20 | self | 检查点 B/C（15 个 descriptor + PG 显式 DDL + 目录断言） | conditional | 1（改动 3 实现形态） | `03-audit/A-001-self-checkpoints-b-c.md` |
| A-002 | 2026-09-20 | independent | 检查点 B/C/D + `D-021` F-I-005 residual 复审（grok-build grok-4.6 · high · commit `c69ee93d`） | conditional | 1（F-I-001 `records` 断言） | `03-audit/A-002-independent-checkpoints-b-c-d.md` |
| A-003 | 2026-09-20 | self（编排器响应） | 响应 A-002 全部意见 | **pass** | 0 | `03-audit/A-003-response-to-independent-b-c-d.md` |

## A-002 · independent · 检查点 B/C/D 与 F-I-005 residual 复审（2026-09-20）

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **verdict**：conditional
- **完整意见**：[`03-audit/A-002-independent-checkpoints-b-c-d.md`](A-002-independent-checkpoints-b-c-d.md)

### 结论摘要

- m0/m4 进入 checksum **符合**台账 §2 / `D-017`；15 个哈希与编译期 descriptor 由 `TestCompiledMigrationCatalogOwnership` 锁住。
- v75 通用 F-5 执行器 **可接受等价**（合同允许「或同一通用 helper」）。
- `F-I-005` residual ①②③ **可按 `fixed` 闭合**（正式闭合记录见 GOAL-002 台账的复审条目）。
- 开放 required：v73 未实现台账点名的 `records` 不存在断言 → **已由 A-003 判 `fixed`**（m0 guard + 负例测试）。

## A-003 · self（编排器响应）· 2026-09-20

- **verdict**：**pass**（A-002 的 open required 已闭合，无未闭合必改项）
- **完整响应**：[`03-audit/A-003-response-to-independent-b-c-d.md`](A-003-response-to-independent-b-c-d.md)

| finding | 处置 |
|---------|------|
| `F-I-001`（required，v73 `records` 断言） | **fixed**：m0 guard（SQLite `sqlite_master` / PG `information_schema.tables`）+ `TestV73RefusesWhenRetiredRecordsTableIsPresent`（拒绝并回滚）；v73 checksum 重算为 `4dd07092…eb8a9f`，冻结表与语句清单附件同步 |
| `F-I-002`（m4 PRAGMA 哈希与执行不一致） | **fixed**：`VerifyStatements` 改为每 descriptor 一次，与 `RunVerify` 同构；11 个多表 descriptor checksum 重算（4 个单表 descriptor 不变） |
| `F-I-003`（缺同进程 codec↔SQL 对拍） | **fixed**：`TestCodecMatchesRealMigration`（秒族 7 值 + 毫秒族 14 值，同一测试内 codec 期望 vs 真实迁移结果） |
| `F-I-004`（生成器留树） | 保留 + 记录（16/16 字节级复现；默认 skip；仅写生成物与清单） |
| `F-I-005`（v86 `ModuleID` 展示列） | 记录更正，不重写 GOAL-002 已关门附件 |
| `F-I-006`（元数据未同步） | **fixed**（`00-meta` / 本文件信息就绪表） |

## 结论状态

检查点 A/B/C/D 均有实现事实与证据：canonical SQL 与**真实 `MigrationChecksum`** 已落盘（冻结表 + 语句清单附件 + 生成器可复现），`D-021` 的 F-I-005 residual 复审已落盘并判可 `fixed` 闭合，A-002 的开放 required 已由 A-003 闭合。**R2 整体仍未放行**：M4（`D-021` 第 9 项 Backup Port 类型表面与 provider、`rebuildOperationLog` residual、关门审计）尚未立项。
