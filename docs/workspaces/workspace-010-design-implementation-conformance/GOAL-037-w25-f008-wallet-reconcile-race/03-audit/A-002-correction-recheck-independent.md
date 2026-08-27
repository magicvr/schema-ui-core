---
title: "A-002 · F-008 根治修正结果复核（independent）"
source: independent
status: recorded
created: 2026-08-23
updated: 2026-08-23
parent: GOAL-037-w25-f008-wallet-reconcile-race
version: 0.1.0
auditor: ox-alpha（DeepSeek Harness /audit 独立会话）
scope: 复核 GOAL-037 全部修正结果——E-002 id 单调序修复、E-003 根治（0050 数据修复迁移 + reconcile 审计原子化）；finding-closure（承接 GOAL-036 A-001 F-008）；关门后复审
verdict: pass
---

# A-002 · F-008 根治修正结果复核（2026-08-23，independent）

## 范围与区间

用户指令：`/audit 复核工作区10 目标36和目标37的修正结果`（本条为其中目标 37 部分；目标 36 部分见该目标 `03-audit/A-004-correction-recheck-independent.md`，同场复核）。

工作区绑定核对：`workspace-010-design-implementation-conformance` / Root `GOAL-001` / canonical 匹配 / `shared_materials_catalog: none`。本目标 `parent: GOAL-036-w25-page-performance-guardrails`，树与 goal-tree 一致。

复核对象：GOAL-037 声称的全部修复落点与其回归主张：

1. E-002：产品 `newID` 与测试替身 `walletStubEntryID` 同毫秒单调计数；
2. E-003 根治①：0050 `wallet_ledger_order_repair` 数据修复迁移（fail-closed）+ 台账锚点同步；
3. E-003 根治②：reconcile 成功审计原子化（job 事务内 `RecordOperationTx`）+ 失败路径可观测；
4. 回归主张：`-count=100` 全绿、全量 go 两轮全绿、vitest 1097/tsc 0。

本轮为**独立复跑验证**：逐文件读码核对 + 实际复跑定向与全量回归。未复跑浏览器 e2e/vitest（前端未改动，vitest 因本机工具链问题无法复跑，见「证据边界」）。

## 成果（有证据）

| 修正项 | 声称 | 本轮独立核验 | 结果 |
|--------|------|--------------|------|
| 产品 `newID` 同毫秒单调序 | UnixMilli(16hex)+seq(8hex)+随机尾；`entryIDSeq atomic.Uint64` & 0xFFFFFFFF | `provider.go:47-57` 实现与 D-001 公式一致；注释如实记录历史缺陷归因 | **属实** |
| 替身同构修复 | stub id 同契约 | `wallet_test.go:74-79` `walletStubEntryID` 同公式（seq & 0xFFFFFFFF + newOperationID 尾） | **属实** |
| 审计断言最终形态 | 原子化后恢复同步单次查询（E-002 轮询版已退役） | `wallet_test.go:418-443`：单次查询断言六事件全集，注释引用根治 | **属实** |
| 0050 迁移注册 | Version 50 · 方言中立 Go · 双方言接线 | `migration.go:222-232` 描述符（`Apply` + `ApplyPostgres` 均 `migrateOrderRepair*`）；checksum 锚于 `migrate_test.go:690`（`835902cb…`） | **属实** |
| 身份台账锚点同步 | v50 data-only 空表清单 | `identity.go:93` `completeFingerprintCatalogHead = 50`；`identity_test.go:101` `lockedHeadExtraTables[50] = {}`；operations/restart/migrate 测试头更新 | **属实** |
| 0050 三用例 | 乱序修复/健康 no-op/坏数据 fail-closed | `migrate_0050_test.go`（L69/L160/L173）；**本轮复跑 ok** | **属实且有效** |
| 成功审计原子化 | CommitFunc 内 `RecordOperationTx`，审计失败 → job 失败 | `jobs.go:113-146`：类型断言 `TransactionalRecorder` → 事务内写事件，错误 fail closed；`recordTerminal` 成功分支显式 no-op 防双写（L148-155） | **属实** |
| 失败/取消可观测 | 不再吞错 | `jobs.go:164-187`：best-effort 路径三处 `slog.Error` | **属实** |
| stub 与产品同构 | handler 测试替身同样原子化 | `wallet_test.go:96-120,193`（CommitFunc 内 RecordOperationTx）；`jobs_test.go:31-34` operationRecorder 实现 TransactionalRecorder | **属实** |

### 本轮独立复跑清单（事实）

| 命令（apps/api） | 结果 |
|------------------|------|
| `go test ./internal/handler/ -run TestWalletLifecycleAndAdjustFlow -count=100` | **ok**（exit 0，19.6s）——E-002/E-003 主张独立复现 |
| `go test ./internal/modules/wallet/ ./internal/modules/scheduledtasks/ -count=1` | **ok ×2**（含 `TestNewIDSameMillisecondOrdering` / `TestNewIDCrossMillisecondOrdering`） |
| `go test ./internal/store/ -run "TestMigration0050\|TestStoreIdentityFingerprint\|TestMigrations" -count=1` | **ok** |
| `go test ./... -count=1`（全量） | **40 包全 ok，0 FAIL，exit 0** |

提交链核对：`be62582`（id 序修复，3 文件）→ `dbf919d`(根治批，11 文件含 order_repair.go 266 行 / migrate_0050_test.go 198 行) → `9bc3dc5`（治理文档），与 E-002/E-003 记录一致。

## 对照成功标准

- C1 机制定性：E-001 失败帧（`replay apply failed: insufficient balance`）+ 排序假设失效因果链完整；A/B 排除记录在案。
- C2 方案冻结：D-001 采纳/拒绝取舍清楚；多实例边界按 README R5 书面保留（部署边界，非任务残余）。
- C3 实施回归：`-count=100` 主张**本轮复现成立**。
- C4 关门：A-001 self pass 的 F-002/F-003 fixed 闭合经本轮代码+复跑证实有效；I-001/I-002 closed 有据；「GOAL-037 关门后回归关门 GOAL-036」的用户书面约定已履行。

## Findings

| F-ID | 级别 | 严重度 | 内容 | 处置建议 |
|------|------|--------|------|----------|
| — | required | — | 无 | — |

无新增 required/recommended。备注两条观察（非 finding）：① `newID` 契约在产品与替身两处手工镜像，未来若再改格式存在漂移可能——现有排序单测构成概率性网络，暂不需动作；② `_txlock=immediate` 的栅栏缺口属 store 连接面（GOAL-036 A-004 F-010 承接），不在本目标范围。

## 必改项汇总

**无**。

## 与既有意见的异同

与 A-001（self 关门）：结论一致且加固——其全部 pass 判定所依赖的回归主张中，最关键的高频窗口（-count=100）与全量 go 由本轮独立复现；无虚报发现。

## 结论 + 建议给编排器/用户的下一步

**verdict: pass。**

GOAL-037 修正结果**名实相符**：机制定性、方案、两级修复（新写入单调序 + 存量数据 0050 重排 + 审计原子化）全部落地且可复现；「无遗留项」声明经核验成立。done 4/4 维持，无需任何重开动作。

建议 `/govern`：仅需知悉本复核结论；如采纳 GOAL-036 A-004 的 F-009/F-010 卫生项，可一并在编排响应中处理。

## 声明

本意见不修改 status/progress/checkpoint/goal-tree；响应与任何状态变更由 `/govern` 处理。
