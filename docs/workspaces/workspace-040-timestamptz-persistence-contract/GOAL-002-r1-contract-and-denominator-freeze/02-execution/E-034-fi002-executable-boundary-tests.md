---
id: E-034-fi002-executable-boundary-tests
doc_type: goal-execution-entry
status: recorded
date: 2026-09-20
parent: GOAL-002-r1-contract-and-denominator-freeze
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# E-034 · F-I-002 可执行边界测试落地（D-020 裁决 A）与 A-040 响应

## 事实

1. **用户 2026-09-20 P-004 裁决两项**：
   - **D-020**：F-I-002 的「可执行负值/越界测试」**允许在一次性验证库上跑**（不触碰生产 schema，不越 C2 门禁），R2 再重定向到真实转换迁移。落盘 child `01-decision/D-020-fi002-executable-test-scope.md`。
   - 下一步优先级 = 先审 C3（F-I-004）+ 按该裁决处理 F-I-002。
2. **新增可执行验证包** `apps/api/internal/w040contracttest/`（**仅测试**，无生产代码），`contract_boundaries_test.go` 三个测试全部通过：

| 测试 | 覆盖 |
|------|------|
| `TestRebuildConversionBoundaries` | 14 个样本：历元 0、正值、**999 ms 无进位**、负毫秒 `-1/-999/-1000/-1001`、负一天、1914 年负瞬间、**公元 9999 年**、NULL 保留、D0 `0→NULL`、D0 正值、voucher `0→NULL`；全部断言 **27 字符** fixed-6；重建后 `foreign_key_check`=0、`integrity_check`=ok、`legacy_old` 已删 |
| `TestNegativeMustFailClosed` | 两个子测试：voucher 负值**必须 fail closed**（Root D-012）；普通列负值**是合法 instant**（Root D-015） |
| `TestFixedSixLexicalOrder` | fixed-6 TEXT 的词法序 = 时刻序（含负值、历元 0、公元 9999 年） |

3. **该测试自身暴露并纠正了两处缺陷**（这是「可执行测试」相对「用例 ID」的实际价值）：
   - **负值政策被过度泛化**：初稿对所有时间列一律「负值即 fail closed」，被测试拦住。核对权威决策后确认：**Root `D-012` 的负值 fail-closed 是 voucher 专属**（`D-012` 原文只提名 `vouchers.expires_at`/`redeemed_at`）；**Root `D-015` 明确「负 epoch 不是 sentinel」**，是合法 instant。已收正并把该区分固化为断言。
   - **单位错配**：年份 9999 的毫秒用例初稿误传秒值，测试立即报出 `1970-01-21` 而非期望值。
4. **验证环境**：`go build ./...` exit 0；`go vet ./internal/w040contracttest/` exit 0；`go test ./internal/store/` 抽样仍 ok（未破坏既有套件）。
5. **独立审计 A-040 完成并落盘**（`source: independent`，`auditor: grok-build (grok-4.6 · reasoning high)`）复审 C3 边界文件：verdict `conditional`，**F-I-004 收窄但不可闭合**，开放 required 仍为 **3**，**无新 required**，新增 recommended **F-I-027**。
   - **已接受**：新文件内部三类区分与双 token；`kernel` + `internal/backup` 包路径；SQLite 行号与「禁止事务内调用」；harness 规格（含文件名）作为 R1 设计足够；§1 四条现状事实成立；§7 与 `D-019` 自洽。
   - **F-I-004 仍 open 的全部 7 项剩余**（A-040 §G）：① 旧 runbook L28/L38/L41 仍用步骤 1 的 `<artifact>` 做目标形状 restore（A-029 原文矛盾一字未改）；② 冻结包权威未收口（conversion contract §0 仍写本文件「尚未落盘」；Port draft L60 仍允许 `snapshotBeforePending` 当 RecoveryPoint 源）；③ PG before/after 调用点未写到 `postgres.go`（`applyPendingPG` 无 snapshot、无 `verifyIntegrity`）；④ 调用点 1 把 A 与 C 挤成「批次前一次」，而 C 实际是循环内 per-migration；⑤ 调用点 2 失败后 `actionNoop` 无重试，B 可能永不补创；⑥ `TestLegacyArtifactMustFail` 未钉形状/合同错误类（缺文件也会绿）；⑦ 转换后的 C→B 无机械身份，仅靠文件名政策。

## 证据

- 测试文件：`apps/api/internal/w040contracttest/contract_boundaries_test.go`（`go test ./internal/w040contracttest/ -v` 全绿）。
- `D-020`：`01-decision/D-020-fi002-executable-test-scope.md`。
- A-040 全文：`03-audit/A-040-r1-independent-e033-c3-boundary-fi004.md`。
- 权威决策对位：`GOAL-001/01-decision/D-012-voucher-invalid-value-policy.md`（voucher 专属 neg fail-closed）、`D-015-negative-instant-truncation.md`（负 epoch 非 sentinel）。
- 本轮 `apps/` **新增一个仅测试包**；无生产代码变更；`go build ./...` exit 0。

## 状态评估

- **开放 required = 3**（F-I-002、F-I-004、F-I-005）；**本条未闭合任何 required**。
- **F-I-002**：其唯一实质剩余（可执行负值/越界测试）**已按 `D-020` 落地并可复跑**；是否据此闭合待 independent 复审该测试文件。
- **F-I-004**：A-040 给出 7 项收口清单，本轮未做（本轮聚焦 F-I-002 测试）；下一轮按清单逐项收口（仍不改 `apps/`）。
- **F-I-005**：需 R2 落码才能记录真实哈希，结构性依赖。
- 下一步：`/audit` 复审本轮测试文件（F-I-002）与 C3 的 7 项收口（F-I-004）。
