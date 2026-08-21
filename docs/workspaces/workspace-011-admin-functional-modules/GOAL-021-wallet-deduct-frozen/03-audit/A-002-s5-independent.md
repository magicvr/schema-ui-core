---
id: A-002
goal: GOAL-021-wallet-deduct-frozen
title: S5 关门独立审计 · 冻结扣款
date: 2026-08-16
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
scope: S5 关门（方案+实现合并审）
audit_type: close-out
verdict: pass
status: recorded
parent: GOAL-021-wallet-deduct-frozen
created: 2026-08-16
updated: 2026-08-16
version: 1.0.0
---

# A-002 · S5 关门独立审计（独立 · 冻结扣款）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high）
- **类型** / **scope**：close-out · S5 关门（小目标合并审视：D-001 方案 + 实现；data 门禁：资金原语 + 迁移 0033/0034）
- **verdict**：**pass**

## 范围与区间

- **工作区**：`workspace-011-admin-functional-modules`（`root_goal: GOAL-001-admin-functional-modules`；`canonical_scope` 匹配；`plan_refs`/`primary_plan` = `VP-011-admin-functional-modules`；`shared_materials_catalog: none`）。未读取或比较其他工作区目标状态。
- **已通读**：`workspace.md`、`goal-tree.md`（GOAL-021 行）、本目标 `00-meta` / `01-decision.md` + D-001 / `02-execution.md` + E-001 / `03-audit.md` + A-001；同区上游 `GOAL-019-r3-s14-wallet-ledger/03-audit/A-008-industry-benchmark-independent.md` 与 `attachments/audit-A-008-independent.md`（F-001/F-002 原文）。
- **代码核对**：`modules/wallet/store/repository.go`（`Apply` `EntryDeductFrozen`、幂等比对 + RefType/RefID、`Mutate` disabled、`checkAccountChain` 回放）；`modules/wallet/migration/migration.go`（0033 rename/create/copy/drop + CHECK 超集 + UNIQUE/索引）；`modules/operationlog/migration/migration.go`（0034 事件 CHECK 重建 + `rebuildOperationLog`）+ `repository.go`（`EventWalletDeductFrozen`）；`handler/wallet.go`（deduct-frozen 端点 / `walletMutate` / `writeWalletError`）；`modules/wallet/schema/wallet.json` 与 en/zh i18n；`modules/wallet/provider.go` / `kernel/profile.go` 路由声明；`compiled/persistence.go` 注册。
- **测试核对**：`store/repository_test.go`（`TestApplyTable` deduct 行、`TestApplyDeductFrozenWithFrozenBalance`、`TestMutateDeductFrozen`、`TestMutateIdempotencyRefCompare`、既有 `TestMutateIdempotency` / disabled adjust）；`handler/wallet_deduct_test.go`；`handler/wallet_test.go` 门禁表含 deduct-frozen；`store/migrate_test.go` 0033/0034 checksum 冻结；`store/operations_test.go` applied=34。
- **独立复跑（2026-08-16）**：
  - `apps/api` `go test -p 1 -count=1` · `modules/wallet/store`（`-run TestApply|TestMutate`）+ `handler`（`-run TestWallet`）+ `store`（`-run TestCompiledMigrationCatalogOwnership|TestMigrate0005PreservesOperationLogRows`）**全绿**。
  - `apps/web` `vitest run` · `all-module-schemas-dval.test.ts` + `schema-keys.structural.test.ts` → **2 文件 / 21 通过**（含 wallet.json D-VAL 与 labelKey 双语目录）。
  - 未复跑全量 `go test ./...` 与全量 vitest / e2e / 容器冒烟（本目标无 E-003 全量回归声明可采信；本轮以定向复跑 + 源码核对为准）。
- **covered**：D-001 §1–§4/§6/§7；A-008 F-001/F-002 对应实现；I-001～I-003；S1～S5 对照；data 门禁（语义 / 0033·0034 / 幂等 Ref / 端点门禁与审计 / schema·i18n·D-VAL）。
- **excluded**：不改 `status` / `progress` / goal-tree / `00-meta` / D-001 / 执行台账；不写入 GOAL-019 响应节；F-003～F-011 演进项不在本批实现范围。
- **保证等级**：L0。不得解读为第三方鉴证。

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| deduct_frozen 语义：total-=d、frozen-=d、available 不变 | `store/repository.go` `Apply` L363–371：`AmountDelta <= 0` → `ErrInvalidEntry`；否则 `total -= d`、`frozen -= d`、available 不写。`TestApplyDeductFrozenWithFrozenBalance`：基线 (100,60,40) + d=25 → **(75,60,15)**。`TestMutateDeductFrozen`：1000 入账 + 400 冻结 + 250 扣减 → **(750,600,150)** |
| d<=0 → ErrInvalidEntry；frozen<d → ErrInsufficient | L367–368 `<= 0` → `ErrInvalidEntry`；L375–377 任一余额 &lt; 0 → `ErrInsufficient`。`TestApplyTable`：「deduct frozen zero / negative / no frozen」均 `wantErr`；`TestApplyDeductFrozenWithFrozenBalance` d=41 → **`ErrInsufficient`**；`TestMutateDeductFrozen` 超额 → **`ErrInsufficient`** |
| disabled 拒绝 | `Mutate` L446–448：`Status != active` → `ErrDisabled`（在 `Apply` 之前，对全部 entry_type 一律拒绝）。`TestMutateVersionConflict` 以 **adjust** 覆盖该入口；deduct_frozen 走同一分支 |
| 快照恒等式保持 | `Apply` L378–380：`total != available+frozen` → `ErrInvalidEntry`。DDL 0033 保留 `CHECK (balance_after_total = balance_after_available + balance_after_frozen)` 与三列 `>= 0`（`migration.go` L78–89）。`checkAccountChain` L638–651 用同一 `Apply` 回放；`TestMutateDeductFrozen` 对账 **`ResultConsistent`** |
| 幂等比对纳入 RefType/RefID（同 key 异单据 → conflict；余额不变） | `Mutate` L408–409：同载荷条件为 EntryType + AmountDelta + Memo **+ RefType + RefID**；否则 `ErrIdempotencyConflict`。`TestMutateIdempotencyRefCompare`：同 key 换 `RefID` → **conflict**；`BalanceTotal` 仍为 100。既有 `TestMutateIdempotency`（空 Ref）同载荷重放不改余额 |
| 0033 重建保数据：rename → create → copy → drop + UNIQUE/CHECK/索引 | `migrateWalletLedgerDeduct` L94–114：`ALTER … RENAME TO …_old` → `CREATE`（CHECK 超集含 `deduct_frozen`，`UNIQUE (account_id, idempotency_key)`，快照恒等式）→ `INSERT … SELECT` 14 列原样 → `DROP …_old` → `CREATE INDEX idx_wallet_ledger_account`。索引在 drop 之后创建，避免 SQLite rename 后旧索引占名。Version **33** checksum 冻结 `b3135b2888dec0aa6da032121026e1ebfa07bed8bc75396bdc60166a09b3077d`（`migrate_test.go` L600） |
| 0034 事件 CHECK 加 `wallet.deduct-frozen` | `operationLogWalletDeductDDL` L181 超集含 `wallet.deduct-frozen`；`migrateOperationLogWalletDeduct` → `rebuildOperationLog`（L383–402 同款 rename/create/copy/drop/index）。Version **34** checksum 冻结 `b6b54bee8b1baff9b5c8222a6619074ef3be54f211c55bceabe96f1c3a291467`（`migrate_test.go` L601）。`operations_test.go` L54：applied **34** 末项 `operation_log_wallet_deduct` |
| 端点门禁 wallet.adjust；审计 wallet.deduct-frozen；错误码复用 | `wallet.go` L245–247 → `walletMutate`（L323 `requirePermission(..., "wallet.adjust")`）；成功后 `"wallet."+eventSuffix` = `wallet.deduct-frozen`（与 `EventWalletDeductFrozen` L68 同值）；detail 含 accountId/entryId/amountDelta（L359）。`writeWalletError` L371–380：`ErrDisabled`→`WALLET_DISABLED`、`ErrInsufficient`→`INSUFFICIENT_BALANCE`、`ErrIdempotencyConflict`→`LEDGER_IDEMPOTENCY_CONFLICT`、`ErrInvalidEntry`→`INVALID_LEDGER_ENTRY`（无新码）。`TestWalletDeductFrozenEndpoint`：750/600/150 + 超额 **409** + 审计事件存在。`TestWalletRoutesGates` 含 deduct-frozen **401/403** |
| schema / i18n 与既有页一致；D-VAL 覆盖 | `wallet.json`：行操作 `deductFrozen` + `openDeductFrozen` modal（amount/memo/idempotencyKey，`wallet.adjust`，无 `requestMapping`，`permissionIntent: edit`）与 freeze/unfreeze 同构。en/zh：`schema.wallet.submit.deductFrozen` / `schema.wallet.action.deductFrozen`。本轮 D-VAL 全模块 + schema-keys **21/21** |
| Descriptor / BuiltinModules 同步；协议不越界 | `provider.go` L161 与 `kernel/profile.go` L198 均声明 `POST /api/wallet/accounts/{id}/deduct-frozen`。权限三键未增；无新 capability；模块 Version 仍 `2.0.0`。加法路由，**不暂挂**判定与 D-001「无越界」一致（D-002 文档尚未落盘，见 F-003） |
| I-001～I-003 无到期未闭 required | 最晚阶段均为 S1；状态 **verified**（D-001 §1–§3）。本 S5 无开放 required 信息门禁、无 deferred required |
| A-008 F-001 / F-002 实现层可重复核对 | F-001：冻结桶原子扣减，available 不回露。F-002：同 key 异单据不再静默丢单。本条不改 GOAL-019 台账；闭合留痕归 /govern |

## 对照成功标准（S1～S5）

| 标准 | 状态 | 证据 |
|------|------|------|
| S1 方案冻结（D-001） | 满足 | D-001 `accepted`：语义 / 幂等 / 0033·0034 / 端点与前端 / 演进登记 / 未选方案。I-001～I-003 verified |
| S2 实现（store/迁移/handler/schema/i18n/测试） | 满足（台账未记 E-002） | 上表代码与测试均可核对；`00-meta` S2 未勾、执行索引仅 E-001（见 F-003） |
| S3 验证 | 满足（本轮定向复跑；全量未复跑） | store/handler/migrate 定向 go 全绿；D-VAL + schema-keys 21/21。无 E-003 全量声明 |
| S4 go 影响 + 自审 | 实质满足（D-002 / 自审条目未落盘） | 加法路由、权限/pin/装配语义未改。本条合并审视覆盖实质判定；文档见 F-003 |
| S5 本独立关门审 | 本条 | 无 high/med **required**；无到期 required 信息项 |
| A-008 F-001 冻结扣款原语 | 实现可核对 | Apply + 端点 + 审计 + 0033/0034 |
| A-008 F-002 幂等 RefType/RefID | 实现可核对 | Mutate 比对 + `TestMutateIdempotencyRefCompare` |

## Findings

### F-001 · 0033「有流水后升级」保数据测试未落地

| 字段 | 值 |
|------|-----|
| level | recommended（med） |
| status | open |
| evidence | D-001 §6 要求「重建保数据（有流水后升级）」。全仓 `*_test.go` 无「先插入 `wallet_ledger_entries` 再跨 0033 升级并断言行/约束」用例。`TestCompiledMigrationCatalogOwnership` 只冻 checksum；`TestMigrate0005PreservesOperationLogRows` 保的是 **operation_log** 行穿过含 0034 的全量升级，**不是** 钱包流水行穿过 0033。实现 SQL（L94–114）可核对：列清单完整、CHECK 为超集、UNIQUE/索引重建、事务内失败回滚 |
| closure | — |
| 影响门禁 | **不阻断 S5**。数据门禁的实现路径可读且与 operationlog rebuild 先例同构；缺口是测试未钉住既有流水行 |

### F-002 · 部分边界只断言「出错 / 409」，未钉精确哨兵与码体

| 字段 | 值 |
|------|-----|
| level | recommended（low） |
| status | open |
| evidence | `TestApplyTable` deduct zero/negative/no-frozen 仅 `wantErr`，未区分 `ErrInvalidEntry` vs `ErrInsufficient`。disabled 拒写仅用 adjust 覆盖（`repository_test.go` L187）。`TestWalletDeductFrozenEndpoint` 超额只断言 HTTP 409，注释写 `INSUFFICIENT_BALANCE` 但未读 body `code`。实现：`Apply` L367–377、`Mutate` L446–448、`writeWalletError` L371–380 可核对 |
| closure | — |
| 影响门禁 | 不阻断。建议补：d=0/`ErrInvalidEntry`、disabled+`EntryDeductFrozen`、handler 码体 `INSUFFICIENT_BALANCE` |

### F-003 · 执行台账 / S2–S5 检查点 / D-002 尚未落盘

| 字段 | 值 |
|------|-----|
| level | recommended（low） |
| status | open |
| evidence | 执行索引仅 E-001（立项）。`00-meta` S2–S5 未勾、`progress: 1/5`；`goal-tree.md` GOAL-021 行为 **0/5**（与 meta 不一致）。S4 计划的 D-002「不暂挂」文件不存在。本审以代码与定向测试为证据，**不**把台账滞后写成实现缺失 |
| closure | — |
| 影响门禁 | 不阻断本独立审对实现的 pass。编排器关门前须补 E-002～E-005、勾选检查点、写 D-002（或不暂挂一句话）、同步 goal-tree；审计不得改这些字段 |

## 必改项汇总

无 required / 必改项。

| ID | 级别 | 一句话 |
|----|------|--------|
| — | — | 无 |

recommended：F-001（0033 有流水升级用例）、F-002（精确哨兵/码体）、F-003（台账与 D-002）。

## 与既有意见的异同

| 意见 | 关系 |
|------|------|
| A-001（self · 立项 pass） | 同意立项 scaffold。本条 scope 是 S5 方案+实现，不重开立项项 |
| A-008 F-001（GOAL-019 · required · 缺 deduct_frozen） | **实现已对应**：Apply/端点/审计/0033·0034 可重复核对。本条不改 GOAL-019 `status`；建议 /govern 在 GOAL-019 响应节记 **fixed** |
| A-008 F-002（GOAL-019 · required · 幂等漏 Ref） | **实现已对应**：比对含 RefType/RefID；异单据 conflict 且余额不变。同上建议 GOAL-019 **fixed** |
| A-008 F-003～F-011 | 本批不实现。D-001 §5 已登记触发条件；不纳入本 scope 必改 |
| A-007（GOAL-019 close-out pass） | 不撤回。本条是对其后 A-008 响应批次的关门审 |

无与本条 verdict 相反的相关意见冲突（P-004.2 不触发）。

## 结论 + 关门放行条件

**verdict: pass。** D-001 冻结语义、0033/0034 重建、幂等 Ref 比对、端点门禁/审计/复用错误码、schema/i18n/D-VAL 均有可重复核对证据；A-008 F-001/F-002 在本目标实现层已落地。无 high required，无到期 required 信息项。

**关门放行条件（给 /govern，本条不改状态）**

1. 无开放 required finding 需先闭合。
2. 编排器应：补 E-002～E-005 事实、勾选 S2–S5、落盘 D-002（加法路由 / 不暂挂）、同步 `00-meta` 与 `goal-tree` progress、在 GOAL-019 响应 A-008 F-001/F-002 为 **fixed**。
3. F-001～F-003 recommended 不阻断关门；可后续补测或接受为残余。
4. 波次级 e2e / 容器冒烟仍按工作区惯例留批末，不构成本目标 data 门禁缺口。

建议用户下一步：`/govern` 响应 A-002 并办理关门。

## 声明

本意见不修改 status / progress / 方案正文 / goal-tree。响应由 /govern 处理。保证等级 L0。
