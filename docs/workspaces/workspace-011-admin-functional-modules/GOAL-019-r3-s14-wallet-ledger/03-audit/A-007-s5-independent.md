---
id: A-007
goal: GOAL-019-r3-s14-wallet-ledger
title: S5 关门独立审计 · admin.wallet
date: 2026-08-16
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
scope: S5 关门（成功标准 + D-002 落实 + 实现/验证/安全/协议 + 台账）
audit_type: close-out
verdict: pass
status: recorded
parent: GOAL-019-r3-s14-wallet-ledger
created: 2026-08-16
updated: 2026-08-16
version: 1.0.0
---

# A-007 · 独立关门审计意见（S5 · S-14 钱包/账务）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high）
- **类型** / **scope**：close-out · S5 关门（00-meta S1~S5、D-002 v1.1.0 落实、A-001~A-006 链、实现与验证证据、data/安全路径、协议/go、I-001~I-004）
- **verdict**：**pass**

## 范围与区间

- **工作区**：`workspace-011-admin-functional-modules`（`root_goal: GOAL-001-admin-functional-modules`；`canonical_scope` 匹配；`plan_refs`/`primary_plan` = `VP-011-admin-functional-modules`；`shared_materials_catalog: none`）。未读取或比较其他工作区目标状态。
- **已通读**：本目标 `00-meta`、`01-decision.md`、D-001~D-004、`02-execution.md`、E-001~E-006、`03-audit.md`、A-001~A-006；I-011-001 §4 S-14 / §7 / §8。
- **代码核对**：`modules/wallet/migration/migration.go`（0031 DDL）、`modules/wallet/migration/provider.go`、`modules/wallet/store/repository.go`（Apply / Mutate / ReconcileRun）、`modules/wallet/provider.go`、`modules/wallet/schema/wallet.json` + `wallet-entries.json`、`modules/wallet/manifest/fragment.json`、`handler/wallet.go`、`errorcatalog/errorcatalog.go`、`modules/operationlog/migration/migration.go`（0032）、`modules/operationlog/repository.go`（EventWallet*）、`composition/composition.go`、`kernel/profile.go`、`kernel/provider.go`、`modules/compiled/persistence.go`、`handler/export.go`、`handler/error_contract_test.go`、`store/migrate_test.go`。
- **独立复跑（2026-08-16）**：`apps/api` `go test -p 1 -count=1` 覆盖 `modules/wallet/store` / `handler` / `composition` / `store` / `kernel` **全绿**（errorcatalog 无测试文件）；web 定向 `schema-keys.structural` + `app-manifest` + `upstream-fixtures` **77/77** 全绿。未复跑全量 `./...` 与全量 vitest 1004（采信 E-005/E-006 记录 + 本轮定向复跑）。
- **covered**：成功标准对照、D-002 v1.1.0 apply 表/幂等/快照链/乐观锁/disabled/双层审计/0031·0032/权限三键/错误码/协议/组合根 27→30·13→14、A-004 required 实施核验、I-001~I-004、波次级 e2e/冒烟是否可接受。
- **excluded**：不改 `status` / `progress` / goal-tree / `00-meta` / D-002 正文；不跑 e2e 双 profile 与容器冒烟（见波次级事项）。
- **保证等级**：L0。不得解读为第三方鉴证。

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| apply 表与 D-002 v1.1.0 §1 一致 | `store/repository.go` `Apply` L273–309：adjust `d≠0` 且 total/available `+=d`、frozen 不变；freeze `d>0` 且 available `-=d` / frozen `+=d`、total 不变；unfreeze 反向；任一余额 &lt; 0 → `ErrInsufficient`；恒等式破 / 未知类型 / 零额 → `ErrInvalidEntry`。`TestApplyTable` 覆盖正负调账、freeze、零/负 freeze、透支、未知类型 |
| 幂等复合范围 + 同/异载荷分流 + 禁裸 key | `Mutate` L321–357：`WHERE account_id = ? AND idempotency_key = ?`；同 type/delta/memo → 返回既有；异载荷 → `ErrIdempotencyConflict`。DDL `UNIQUE (account_id, idempotency_key)`（`migration.go` L52）。`TestMutateIdempotency`：同载荷不改余额、异载荷冲突、他户可复用同一 key |
| 快照链重放：链序、末笔==当前、恒等式 | `checkAccountChain` L542–598：`ORDER BY created_at ASC, id ASC`；首笔 prev=(0,0,0)；`Apply(prev, entry)==after_*`；每笔 `after_total=after_avail+after_frozen`；末笔 vs 账户三余额。provider `newID`（`provider.go` L38–48）毫秒时间序 id，保证同秒并列按创建序 |
| 乐观锁 + WithTx 原子性 | `Mutate` L379–392：`UPDATE … version = version + 1 WHERE id = ? AND version = ?`，`affected==0` → `ErrVersionConflict`；余额更新与流水 INSERT 同 `WithTx`。`UpdateStatus` L233–253 同款。`TestMutateVersionConflict` 覆盖 stale status |
| disabled 拒写（含解冻） | `Mutate` L372–374：`Status != active` → `ErrDisabled`（入口在 Apply 之前，对 adjust/freeze/unfreeze 一律拒绝）。`TestMutateVersionConflict` / `TestWalletIdempotencyAndStatus` 覆盖 disabled 后 adjust |
| 流水不可变 | 全仓 `wallet_ledger_entries` 仅 SELECT / INSERT，无 UPDATE/DELETE（`repository.go` L327 / L394 / L443 / L550） |
| 双层审计 | 账本 entry 同事务写入；handler `recordWalletEvent`（`wallet.go` L288–297）在成功后独立记 `wallet.account-create` / `account-update` / `adjust` / `freeze` / `unfreeze` / `reconcile`（`_ = operations.RecordOperation`，与 MFA 非同事务残余一致）。常量 `operationlog/repository.go` L61–66；0032 CHECK 含六事件（`operationlog/migration/migration.go` L162） |
| 迁移 0031/0032 | wallet Version **31** checksum `bc92082f…`（`migrate_test.go` L597）；operationlog Version **32** checksum `1c27e86c…`（L598）。DDL：账户恒等式 CHECK L35、流水快照恒等式 CHECK L53、`amount_delta != 0` L41、复合 UNIQUE L52、索引 L55。`compiled/persistence.go` L34 注册 |
| 权限三键分端点 | `wallet.go`：list/entries/reconcile/runs → `wallet.read`；create/PATCH → `wallet.write`；adjust/freeze/unfreeze → `wallet.adjust`。Descriptor 声明三键（`provider.go` L149）；`testsupport/store.go` L64–67 种子 |
| 错误码冻结集 + 双语 | 八码进 `errorcatalog.go` L141–148（en+zh+messageKey）与 `error_contract_test.go` L57–58；i18n `en-US.json` / `zh-CN.json` L576–583 对齐。`TestWalletErrorCodesCataloged` 核对 catalog 完备 |
| 组合根 27→30 / 13→14 | `composition_test.go` L471–473 `wantPermissions: 30, wantNavigation: 14`；`profile.go` L90–92 ProfileAdmin += `admin.wallet`（mvp/demo 不含）；`provider.go` L403–419 DefaultNavigationOrder 14 项末项 `menu_wallet`；`navigation_order_test.go` L12–28 锁序。本轮 composition/kernel 复跑全绿 |
| 协议不越界 | 无新 capability：`wallet.json` 用既有 `actions.row.navigate`（上游 version-negotiation cases）；`wallet-entries.json` 用既有 `data.route-binding`（`conformance-claim.json` L28）。pin 仍 `v2.8.0`（`provenance-v2.8.json` / Charter）。`export.go` L126 仍只允许 `users`/`roles`，wallet 未登记导出资源 |
| I-001~I-004 均 verified，无到期 required | `00-meta` 信息表；最晚阶段均为 S1；无 `deferred`；本 scope 无开放 required 信息门禁 |
| A-004 F-001/F-002 实施层可重复核对 | 金额原语按 apply 表落地；幂等按复合 UNIQUE + 分流落地。不重开 |

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| S1 方案冻结（D-002 v1.1.0 + A-003/A-004/A-005） | 满足 | D-002 accepted v1.1.0；A-004 F-001/F-002 经 D-003 **fixed** + A-005 pass；本轮实施与表一致 |
| S2 实现（D-002 §8 1–8） | 满足（测试覆盖有 recommended 残留） | E-005；0031/0032、store/handler/组合根/schema/web 均可核对；§8.3「九端点」实现为 **10 路由**（多 `GET /api/wallet/entries` 查询变体，E-005 已记，服务 wallet-entries `data.route-binding`） |
| S3 验证（go + web；e2e 波次级） | 满足 | E-006；本轮定向 go/web 复跑全绿；e2e/冒烟见波次级事项 |
| S4 go 判定 + 自审 | 满足 | D-004 内容扩展不失效；pin/装配语义未改；A-006 self pass |
| S5 本独立关门审 | 本条 | 无 required；无到期 required 信息项 |
| A-004 F-001 apply 表可执行 | 闭合 | `Apply` + `TestApplyTable` + Mutate 集成 |
| A-004 F-002 幂等隔离 | 闭合 | 复合 UNIQUE + 带 account_id 查找 + `TestMutateIdempotency` |

## Findings

### F-001 · 对账不一致路径与 operationlog 落库未被测试断言（A-006 覆盖陈述过满）

| 字段 | 值 |
|------|-----|
| level | recommended（med） |
| status | open |
| evidence | D-002 §6 要求「不一致清单与 mismatch_count」及「审计事件」。`TestReconcileConsistentAndInconsistent`（`repository_test.go` L195–243）只跑一致账本；注释自承无法经公开 API 篡改，**从未**构造 inconsistent。`wallet_test.go` `TestWalletLifecycleAndAdjustFlow` 标题写「audit events」，正文在 reconcile consistent 后结束，全仓 `*_test.go` 无 `EventWallet*` / `wallet.account-create` 断言。A-006 将「快照链重放」与「审计六事件落 operationlog」标为 handler/store 测试已覆盖 |
| closure | — |
| 影响门禁 | 不阻断本目标关门；建议补一条篡改/不一致用例 + 一条 operationlog 事件断言 |

实现可核对：`checkAccountChain` 在 snapshot 不符 / 恒等式破 / 末笔≠当前时返回 reason；handler 成功路径调用 `recordWalletEvent`。缺口是**测试未证明** mismatch 分支与事件行真正写入，不是实现缺失。

### F-002 · 权限三键隔离未做分键用例（实现按端点绑定正确）

| 字段 | 值 |
|------|-----|
| level | recommended（med） |
| status | open |
| evidence | D-002 §3 / 019-F-002：金额变动专用 `wallet.adjust`，与 `wallet.write` 拆分。`TestWalletRoutesGates` 只测匿名 401 与 editor **无任何钱包键** 的 403；403 用例还漏了 `GET /api/wallet/entries`。无「仅 read 不能 adjust / 仅 write 不能调账 / 仅 adjust 不能建户」用例 |
| closure | — |

实现逐路由 `requirePermission` 与方案一致（见成果表）。不构成现网越权；分键是本目标安全主张的核心，缺隔离测试。

### F-003 · 生产 `CreateAccount` 不校验 ownerType/ownerID（测试桩有、实现无）

| 字段 | 值 |
|------|-----|
| level | recommended（low） |
| status | open |
| evidence | handler 测试桩 `wallet_test.go` L31–33：空 ownerType/ownerID → `ErrInvalidEntry`。生产 `provider.go` `Service.CreateAccount` L61–77 只默认 currency，原样插入。非法 `owner_type` 走 SQLite CHECK → `INTERNAL`（非 `INVALID_WALLET_BODY`）；空 `owner_id` 可插入（`TEXT NOT NULL` 接受 `""`） |
| closure | — |

账户表 CHECK/UNIQUE 仍挡住非法枚举与重复持有方，不是余额完整性缺口。建议生产与桩对齐：空/非法 → 400。

### F-004 · 对账工具栏按 `wallet.write` 显隐，API 门禁是 `wallet.read`

| 字段 | 值 |
|------|-----|
| level | recommended（low） |
| status | open |
| evidence | D-002 §3：`POST /api/wallet/reconcile` 与 `GET …/runs` 为 `wallet.read`。`handler/wallet.go` L153 / L175 确实 `wallet.read`。`wallet.json` L198–217：表格 `edit` = `wallet.write`，toolbar `reconcile` 的 `permissionIntent` 为 `edit` |
| closure | — |

仅 read 的操作者可调 API，UI 看不到对账按钮。不是越权，是呈现比方案更严。后续把 intent 改到 read（或单独 permission 键）即可。

## 必改项汇总

无 required。无到期且影响本 scope 的 required 信息项。

## 波次级事项（可接受）

| 项 | 本目标 | 先例与说明 |
|----|--------|------------|
| 双 profile e2e | 未跑 | 00-meta S3 / E-006 / A-006 已写「归 S5 波次」。与 GOAL-016/017 同款推迟。`apps/web/e2e` 本轮未检索到 wallet 专属断言（证据：web 定向测试不含 e2e） |
| 容器冒烟 V-007/V-008 | 未跑 | R3 批末统一。不得用本 pass 代替波次证据 |

**本目标关门可接受**这些项留到 R3 第四批收尾统一验证；不构成 required。批末必须补跑，失败则回流。

## 与既有意见的异同

- **A-006**（self · pass）：同意 S2–S4 主体落地、apply 表/幂等/乐观锁/组合根/协议可进 S5。本意见用代码行号与定向复跑独立核对后维持 pass。**不同意**「对账链不一致 + 六事件」已被测试覆盖的字面（F-001 recommended）。
- **A-004 F-001/F-002 required**：本轮按实现复核为 **fixed**（与 A-005 / D-003 一致），不重开。
- **A-005** recommended 残留（§8 26→29、§1 旧 UNIQUE 句、结论状态过时）：`03-audit.md` 响应记录称已勘误；实施按 §6 27→30 / 13→14，不采用旧基数。本条不重开文案残留。
- **A-002** 019-F-001/F-002（空引用 / 权限拆分）：D-002 §1/§3 已裁定且落地（`ref_type`/`ref_id` 可 NULL；三键分端点）。不重开。
- **A-001 / A-003**：立项与 S1 self，与本 close-out 无冲突。
- 无意见冲突需 P-004。

## 结论 + 建议给编排器/用户的下一步

**verdict: pass**。S-14 作为 data 门禁目标，D-002 v1.1.0 的金额原语、幂等隔离、快照链、乐观锁、disabled 拒写（含解冻）、双层审计、0031/0032、权限三键与协议边界均有可核对实现；A-004 required 闭合可重复核对。无 high/med required；I-001~I-004 均 verified。

**可关门**（`status: done` 与 progress 5/5 由 `/govern` 执行；本意见不改状态）。

建议 `/govern`：

1. 响应本意见：记录 0 required；F-001~F-004 recommended 可带入后续加固或批末补测，不阻断关门。
2. 将 S5 检查点勾选、progress 重算为 5/5、goal-tree 同步。`workspace.md` R3 行仍写本目标 `0/5`（L50），权威在 goal-tree；关门时一并改指针更干净，非门禁。
3. R3 第四批收尾统一跑 e2e 双 profile + V-007/V-008；失败回流，不把本 pass 当波次证据。

勿用 `progress: 4/5` 作为放行或拒绝依据。

## 声明

本意见不修改 `status` / `progress` / goal-tree / `00-meta` / D-002 正文。响应由 `/govern` 处理。保证等级 L0。
