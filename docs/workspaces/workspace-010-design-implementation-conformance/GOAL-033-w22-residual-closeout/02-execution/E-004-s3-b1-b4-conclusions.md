---
title: E-004 · S3 复核结论：B1–B4 架构债判定与回写记录
status: recorded
created: 2026-08-23
updated: 2026-08-23
parent: GOAL-033-w22-residual-closeout
version: 0.1.0
---

# E-004 · S3 复核结论（2026-08-23）：B1–B4

代码级只读核查（证据均为文件:行级摘录）后判定如下；已兑现项当轮完成源台账回写。

## C8 · B1 R4-I004 operationlog —— **拆分判定**

| 半边 | 判定 | 证据 |
|------|------|------|
| retention | **已兑现 → fixed/verified** | `apps/api/internal/modules/operationlog/retention.go:21-95`（ApplyRetention + StartRetentionSweep，默认 90 天/archive、1–3650 可配）；workspace-012 GOAL-008（done 2026-08-19）交付 site_settings 两列迁移 + sweeper 接线（composition.go:747），有测试 |
| append best-effort | **未兑现（by design）→ 建议续期 residual** | `handler/audit.go:52-54` fire-and-forget 仅 slog；TransactionalRecorder（RecordOperationTx）为敏感路径提供可选 fail-closed，非全局强制 |

处置：retention 半边在源台账（W3·GOAL-006 C1-I003 / GOAL-005 D-003 residual 表）回写兑现注记；append 半边续期需用户书面确认新 review date（建议 **2027-02-01**，触发=合规零日志丢失要求 / 引入幂等重试机制 / 安全审计发现审计缺口）→ 已列入关门 P-004 批次。

## C9 · B2 F-003b document 字节发布 —— **已兑现 → fixed**

- `PageContribution.Document` 字段随贡献集注入字节；`composition.go:528-579` `RegisterSchemas(mux, set.Pages)`，577 行注释「R6 C6.3: finalized page contributions own both metadata and document bytes; the handler has no static document or owner fallback」；`handler/schema.go` 无静态合并兜底。
- 回写完成：[W3 Root A-011](../../../../workspace-003-modular-admin-architecture/GOAL-001-modular-admin-architecture/03-audit/A-011-a010-cohesion-response.md) F-003 行追加「2026-08-23 兑现复核：fixed」。

## C10 · B3 C4-004 PolicyID/Visibility allowlist —— **未兑现 → 建议续期**

- 现状仍为 `validDottedIdentifier` 最小点标识符语法（kernel/contribution.go:275-300）；表达式 `||`/`&&` 被测试明确拒绝（provider_test.go:335-343）；R5（GOAL-012 E-002/E-004）登记 pending 未实施。
- 处置：续期需用户书面确认新 review date（建议 **2026-12-01**，触发=R6 清理完成后首个需要多条件 Visibility 表达式的业务模块 / 安全审计要求策略表达能力）→ 已列入关门 P-004 批次；确认后在 GOAL-010 A-003 与 GOAL-012 E-004 补续期行。

## C11 · B4 C5-002 双 Profile Start/Ready 矩阵 —— **已兑现 → fixed**

- `TestDualProfileLifecycleMatrix`（lifecycle_test.go:80-179，mvp/admin × success/start-fail/ready-fail/stop-fail）+ `TestDualProfileContractMatrix`（provider_test.go:514+）交付完整自动化矩阵。
- 回写完成：[W3 GOAL-011 A-004](../../../../workspace-003-modular-admin-architecture/GOAL-011-r4-c5-acceptance/03-audit/A-004-r4-c5-acceptance-response.md) F-IND-C5-002 行追加「2026-08-23 兑现复核：fixed」。

## 进度影响

C8/C9/C10/C11 复核结论落盘 → 累计 **11/18**。待办：A1/A2/A5+A6/H3（代理进行中）、B1/B3 续期用户追认（关门批次）、S5 回归 + 关门审计。
