---
title: "A-004 · W25 修正结果复核（independent · F-001～F-008 关闭证据复审）"
source: independent
status: recorded
created: 2026-08-23
updated: 2026-08-23
parent: GOAL-036-w25-page-performance-guardrails
version: 0.1.0
auditor: ox-alpha（DeepSeek Harness /audit 独立会话）
scope: 复核 GOAL-036 全部修正结果——A-001 F-001～F-005 的响应闭合（A-002/E-006）、self 新增 F-006/F-007（E-007）、F-008 移交闭环（GOAL-037）；finding-closure + execution-facts；关门后复审
verdict: pass
---

# A-004 · W25 修正结果复核（2026-08-23，independent）

## 范围与区间

用户指令：`/audit 复核工作区10 目标36和目标37的修正结果`。

工作区：`workspace-010-design-implementation-conformance`（Root `GOAL-001-design-implementation-conformance`；canonical 范围匹配；`shared_materials_catalog: none`；VP-010 delivery）。

复核对象：GOAL-036 已声明的全部修正**关闭证据**（非重开方案讨论）：

- A-001（grok-4.6 · conditional）F-001～F-005 → A-002/E-006 声称 fixed；
- self 新增 F-006（顺序断言）/ F-007（预存 flake）→ E-007 声称 fixed（含归因更正）；
- F-008（wallet reconcile 竞态）→ 声称由下级 GOAL-037 根治闭环。

本轮为**独立复跑验证**：不只读台账，实际复跑定向与全量回归、逐文件核对修复落点。未复跑 Playwright 浏览器 e2e 与双栈计时（I-001/I-002 维持原关闭材料，不构成本轮新增主张）；前端 vitest 因本机工具链问题未能复跑（见「证据边界」）。

## 成果（有证据）

| 修正项 | 声称 | 本轮独立核验 | 结果 |
|--------|------|--------------|------|
| **F-001**（high required）FK 未入每条连接 | `sqliteDSNParams` 补 `_foreign_keys=on` | `store.go:57` DSN 常量实含四参数+txlock；注释引用本 finding 归因 | **属实** |
| **F-002**（med required）连接面栅栏缺 FK | `TestFileStoreEveryConnectionEnforcesForeignKeys` 多连接断言 | `store_wal_test.go:109-179`：持满池 Conn 断言每条 FK=1、非首连接上 CASCADE 删 user_roles 归零 / RESTRICT 拒删在用角色 / refresh_tokens 缺主拒插；**本轮复跑 ok** | **属实且有效** |
| **F-003**（recommended）refreshList 不丢 in-flight | 与 reloadList 对称丢弃 | `render.tsx:873-880` 先 `listInFlight.current.delete(key)` 再 bump token；`render.test.tsx:440-485` 双 gate 用例断言发出第二发请求 | **属实（代码级）** |
| **F-004**（recommended）测量脚本不入库/README 过期/路线图漏勾 | README 重写 + 路线图勾选 | `attachments/README.md` 已重写并列证据清单；`00-meta` S6 行 I-002 已勾 ✓ closed | **属实** |
| **F-005**（recommended）批删测试缺 MFA 断言 | 补 user_mfa 播种与计数 | `users_repository_test.go:145-170` 播种 + 删后 COUNT=0 断言 | **属实** |
| **F-006**（self）顺序断言假绿 | 改事件集合断言 | `operations_test.go:68-75` 注释引用本 finding；集合 + 逐条属性断言 | **属实** |
| **F-007**（self→E-007）预存 flake | 真因 = 替身时钟量化；id 加单调计数 | `testhelpers_test.go:385` `"run-test-%x-%x"`(UnixNano+atomic seq)；`scheduler.go:230-240` newRunID 单调加固 + `TestNewRunIDUniqueUnderEntropyFailure` 存在；scheduledtasks 包**本轮 ok** | **属实** |
| **F-008**（移交 GOAL-037） | 下级根治闭环 | 见 GOAL-037 `03-audit/A-002-correction-recheck-independent.md`（本轮同场复核 pass）；关键回归 `-count=100` **本轮复现 ok（exit 0）** | **属实** |

### 本轮独立复跑清单（事实）

| 命令（apps/api） | 结果 |
|------------------|------|
| `go test ./internal/store/ -run "TestSQLiteDSNPragmas\|TestFileStoreEveryConnectionEnforcesForeignKeys\|TestFileStoreWALPoolAndPragma\|TestMemoryStoreStaysSingleConnection\|TestForeignKeyEnabled" -count=1 -v` | **ok**（6 组全 PASS） |
| `go test ./internal/handler/ -run TestWalletLifecycleAndAdjustFlow -count=100` | **ok**（exit 0，19.6s）——GOAL-037 主张的 100/100 独立成立 |
| `go test ./internal/modules/wallet/ ./internal/modules/scheduledtasks/ -count=1` | **ok ×2** |
| `go test ./internal/store/ -run "TestMigration0050\|TestStoreIdentityFingerprint\|TestMigrations" -count=1` | **ok**（0050 三用例 + 锚点） |
| `go test ./... -count=1`（全量） | **40 包全 ok，0 FAIL，exit 0** |

提交链核对：`ba7d5c6`（W25 本体）→ `3eabd59`（A-001 响应）→ `20d4b5d`（F-007/txlock）→ `be62582`（F-008 id 序）→ `dbf919d`（0050+原子化）→ `9bc3dc5`（治理文档关门），与各 E 条目声称一致；工作树 clean。

## 对照成功标准

关门审计 A-003 所依赖的全部外部意见（A-001 required×2、recommended×3、self×2、F-008 承接）**关闭证据均可指回且可重复核对**。C2 后端改法自 F-001/F-002 闭合后可视为安全（本轮多连接探针级测试绿）；「防复发」主张在 store 连接面已由确定性栅栏支撑（缺口见 F-010）。I-001/I-002 维持原关闭材料，本轮不重开。

## Findings（接目标全局 F 序列）

| F-ID | 级别 | 严重度 | 内容 | 处置建议 |
|------|------|--------|------|----------|
| **F-009** | recommended | low | **台账卫生**：`01-decision.md` / `02-execution.md` / `03-audit.md`（索引）frontmatter 仍 `status: active`，而 `00-meta` 已 `done`；下级 GOAL-037 同类索引均为 `done`，姊妹目标实践不一致 | `/govern` 顺手统一为 `done`；不阻断 |
| **F-010** | recommended | low-med | **栅栏缺口**：E-007 以 `_txlock=immediate` 修 SQLITE_BUSY_SNAPSHOT 竞态（A/B 实证），但 `TestSQLiteDSNPragmas` 只钉 busy/WAL/synchronous/FK 四参数，**未钉 txlock**；静默移除将无确定性栅栏拦截，仅剩 wallet 高频用例的概率性网络。与 F-002 同一精神（连接面不变量应全部钉死） | 在 `TestSQLiteDSNPragmas` want 列表补 `_txlock=immediate` 一行即可 |

无 required。两条均不追溯影响任何已闭合 finding 的有效性。

## 必改项汇总

**无**。

## 与既有意见的异同

- 与 A-001（grok-4.6）：结论相容且更进一步——A-001 提出的两项 required 修复经本轮**复跑证实有效**；其「修改方式总评」对前端三条取舍的肯定亦与读码结论一致。
- 与 A-002/A-003（self）：响应与关门的各项 fixed 主张**无一虚报**；E-007 的归因更正（替身时钟量化而非产品回落路径）与代码现状一致。
- 新增差异点仅 F-009/F-010 两条 recommended 卫生项。

## 结论 + 建议给编排器/用户的下一步

**verdict: pass。**

GOAL-036 及其下级 GOAL-037 的全部修正结果**名实相符**：required/recommended/self findings 关闭证据充分且关键项可独立复现（含 100×wallet 生命周期与全量 go 40 包零失败）；关门（done 6/6）与回归关门时序符合用户书面约定。剩余两条 recommended 卫生项可由 `/govern` 低成本处理，无需重开门禁。

1. `/govern` 将 F-009（frontmatter 统一）与 F-010（DSN 测试补钉 `_txlock`）按 fixed 或书面 residual 处置。
2. 若需浏览器层最终背书，可选复跑一次 sqlite e2e 双 profile（非必需——机制层已被单测栅栏覆盖）。
3. 本意见不改 status/progress；GOAL-036 保持 done。

## 声明

本意见不修改 status/progress/checkpoint/goal-tree；响应与任何状态变更由 `/govern` 处理。
