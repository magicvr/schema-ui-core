---
id: A-002
goal: GOAL-022-my-wallet-self-service
title: S5 关门独立审计 · 我的钱包自服务
date: 2026-08-16
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
scope: S5 关门（成功标准 + D-002 §1~§5 + 身份隔离/只读/惰性开户/导航装配/协议门禁）
audit_type: close-out
verdict: pass
status: recorded
parent: GOAL-001-admin-functional-modules
created: 2026-08-16
updated: 2026-08-16
version: 1.0.0
---

# A-002 · S5 关门独立审计（independent · 我的钱包自服务）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high）
- **类型** / **scope**：close-out · S5 关门（成功标准 S1~S5 + D-002 §1~§5；核心安全边界 = 身份隔离与数据暴露）
- **verdict**：**pass**

## 范围与区间

- **工作区**：`workspace-011-admin-functional-modules`（`root_goal: GOAL-001-admin-functional-modules`；`canonical_scope` 匹配；`plan_refs`/`primary_plan` = `VP-011-admin-functional-modules`；`shared_materials_catalog: none`）。未读取或比较其他工作区目标状态。跨区移交仅按本目标 `00-meta` 已有 Q2 路径提及 [workspace-010 GOAL-013](../../../workspace-010-design-implementation-conformance/GOAL-013-w12-product-surface-intent/00-meta.md)。
- **已通读**：`workspace.md`、`goal-tree.md`（GOAL-022 行 3/5 · active）；本目标 `00-meta` / `01-decision.md` + D-001~D-003 / `02-execution.md` + E-001~E-004 / `03-audit.md` + A-001。
- **同区上游口径**：`GOAL-019-r3-s14-wallet-ledger/03-audit.md`（A-007 close-out pass；管理端 `wallet.read`/`wallet.adjust` 分键）；`GOAL-020-wallet-auto-account` D-001（by-owner get-or-create + `auto:true` 审计 + UNIQUE 冲突 `isNew=false`）与 A-005；`GOAL-021-wallet-deduct-frozen` D-001/D-002 与 A-002（S5 independent pass 先例）。
- **代码核对**：`handler/wallet_self.go` + `wallet_self_test.go`；`handler/wallet.go`（`accountToMap`/`entryToMap`/`resourceList` + 管理端门禁）；`modules/wallet/provider.go`；`modules/wallet/manifest/fragment.json`；`modules/wallet/schema/my-wallet.json` + `schema.go`；`kernel/provider.go` `DefaultNavigationOrder`；`kernel/profile.go` `BuiltinModules`；`kernel/navigation_order_test.go`；`composition/composition_test.go`；`modules/wallet/store/repository.go` `GetOrCreateUserAccount`；`auth/auth.go` `IdentityFrom` / `Middleware`；web fixture / `upstream-fixtures.test.ts` / `app-manifest.test.ts` / `schema-keys.structural.test.ts` / `user-menu.test.tsx` / en-US·zh-CN。
- **独立复跑（2026-08-16，本会话）**：
  - `apps/api` `go test -p 1 -count=1` · `handler`（`-run TestWalletSelf|TestWalletRoutesGates`）+ `kernel`（`-run TestDefaultNavigationOrderSnapshot`）+ `composition`（`-run TestSystemDataReconcileUsesFinalizedProfileContributions`）+ `modules/wallet/store`（`-run TestGetOrCreateUserAccount`）**全绿**（handler 6.132s / kernel 0.693s / composition 2.539s / store 1.041s）。
  - `apps/web` 本地 vitest · `user-menu.test.tsx` + `upstream-fixtures.test.ts` + `app-manifest.test.ts` + `schema-keys.structural.test.ts` + `all-module-schemas-dval.test.ts` → **5 文件 / 99 通过**。
  - 独立重算 admin fixture SHA-256 = `0efb4054d0473ce0649277e9b755b2109e473e89fef6507df0badd5d403868a0`，与 `upstream-fixtures.test.ts` L73–74 钉死值一致。
  - **未**复跑全量 `go test ./...`（E-004 声明 34 包）与全量 vitest 1038；**未**复跑 :25099 实机冒烟 / e2e / V-007/V-008。本轮以定向复跑 + 源码核对为准；波次级冒烟按工作区惯例留批末。
- **covered**：S1~S5 检查点对照；D-002 §1 只读 + 身份隔离；§2 路由/identity-only/惰性开户/导航；§3 schema 页；§4 管理端分工与装配；§5 测试面；I-001/I-002；D-003 协议/门禁主张。
- **excluded**：不改 `status` / `progress` / goal-tree / `00-meta` / D-002 / 执行台账；不重开 GOAL-019/020/021 已关门项；不把 A-001 观察项（min-units）升格为必改。
- **保证等级**：L0。不得解读为第三方鉴证。

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| `/api/wallet/me` 与 `/me/entries` 只从会话身份推导 owner，永不读客户端 `ownerId` | `wallet_self.go` L35–51 / L54–81：owner 仅 `selfIdentity` → `user.ID`。`selfIdentity` L86–95 只调 `auth.IdentityFrom`；查询串仅 `/me/entries` 读 `page`/`pageSize`（L69–70）。全文件无 `ownerId`/`accountId`/`PathValue` 入参。`IdentityFrom`（`auth.go` L416–418）只读 Middleware 写入的 context；Middleware（L435–478）只验 Bearer，失败 401 |
| 返回数据只属于会话用户；无横向越权经 by-owner / entries 查询 / accountId 路径 | 自服务：`ListEntries(account.ID, …)` 的 `account` 来自 `GetOrCreateUserAccount(user.ID)`（L60、L71）。管理端旁路全部 `requirePermission`：`GET /api/wallet/by-owner/{ownerId}` L74–75 `wallet.read`；`GET …/accounts/{id}/entries` 与 `GET /api/wallet/entries?accountId=` → `walletListEntries` L296–297 `wallet.read`；调账/冻/解冻/扣减 L323 `wallet.adjust`。`TestWalletRoutesGates` L108–131：editor 对上述路径（含 by-owner、entries 两种形态、全部资金写面）**403**。`TestWalletSelfEntriesOwnScope`：alice 2500 / 「grant alice」；bob（viewer）`/me` ownerId=`user-bob`、余额 1000 |
| editor/viewer 无 `wallet.*` 只能读自己 | `TestWalletSelfIdentityOnly`：匿名双端点 401；editor（无 wallet 键）`GET /me` 200。bob=viewer `/me` 200。管理端三键仍 `PolicyAdmin`（`provider.go` L214–218）；`menu_wallet` 仍 `wallet.read`（L223–233）。E-004 冒烟：editor `/me` 可读、`/api/wallet/accounts` 403（本轮未复跑冒烟，门禁由 `TestWalletRoutesGates` 独立钉住） |
| 只读：无资金操作自服务入口；schema 页无写 action/toolbar | `WalletSelfRoutes` 仅注册两条 **GET**。`my-wallet.json` 无 `actions` / toolbar / adjust / freeze / unfreeze / deduct（全文检索无匹配）。`PageContribution` `Actions: []string{"list"}`、`DataSource: /api/wallet/me`（`provider.go` L203–213）。管理端 `wallet.json` 写 action **仍在**（`adjust`/`freeze`/`unfreeze` 等），未误删 |
| 惰性开户幂等；首次 `wallet.account-create`（`auto:true`）；重复读不重复审计 | 复用 GOAL-020 `GetOrCreateUserAccount`（`repository.go` L235–289：UNIQUE 冲突 `isNew=false` 再重读）。handler `if created` 才 `recordWalletEvent(..., EventWalletAccountCreate, …"auto":true)`（`wallet_self.go` L46–48、L65–67）。`TestWalletSelfAutoCreateAndIdempotency`：首次信封 `ownerType=user` / `ownerId=user-alice` / 余额 0；二次同一 `id`；`wallet.account-create` **恰 1** 且 detail 含 `"auto":true`。本轮该测试绿 |
| statCard 信封可绑定 | `/me` 写 `resourceList{Items:[]{accountToMap}, Total:1}`（L50）。`accountToMap`（`wallet.go` L397–409）含 `balanceTotal`/`balanceAvailable`/`balanceFrozen`，与 `my-wallet.json` 三张 statCard `valueField` 一致。renderer `render.tsx` L2156–2163：`fetchResourceList` 取 `items[0][valueField]`（与系统监控单行信封同构） |
| user-nav 槽位 account → my-wallet → settings | `DefaultNavigationOrder`：`menu_account` → `menu_wallet_self` → `menu_activity` → `menu_settings`（`kernel/provider.go` L410–415）；快照测试 L18–25 同步。fixture `app-manifest.admin.json` L363–391 同序。`user-menu.test.tsx` L111：`["Account", "My wallet", "Settings", "Sign out"]`。`menu_wallet_self`：`PolicyAdminEditorViewer` + `Permission: ""`（`provider.go` L238–248），与 `menu_account`（`modules/account/provider.go` L102–110）同款 identity-only |
| 装配加法；权限三键未增；无迁移；Profile/pin 主张成立 | Descriptor / `BuiltinModules` 均加 `GET /api/wallet/me`、`GET /api/wallet/me/entries`、`my-wallet`、`menu_wallet_self`；`Permissions` 仍三键（`provider.go` L169；`profile.go` L199）。`composition_test.go` L473–475：admin **permissions=30（不变）**、**navigation=14→15**。`CompiledPersistence` 仍 `nil`（`provider.go` L174–176）；钱包迁移最高仍 0033（无 0035）。无新 capability。本目标未改 `provenance*.json` |
| I-001 / I-002 无到期 required | 均 `required`、最晚 S1、状态 **verified**（`00-meta` / `01-decision.md` / D-002 §1/§2）。无 `deferred` required |

## 对照成功标准（S1～S5 + D-002 §1～§5）

| 标准 | 状态 | 证据 |
|------|------|------|
| S1 方案冻结（D-002；I-001/I-002） | 满足 | D-002 `accepted`；用户裁决只读 + `/my-wallet` 惰性开户；I-001/I-002 verified（E-002） |
| S2 实现（API + schema + user-nav） | 满足 | E-003 产物路径与上表代码一致；本轮定向测试绿 |
| S3 验证 | 满足（本轮定向复跑；全量/实机未复跑） | E-004 记录 go 34 包 + vitest 1038 + :25099 冒烟。本轮核验 handler/kernel/composition/store + web 5 文件/99。全量数字本轮不独立背书 |
| S4 go 判定 + 自审 | 满足（证据在 D-003 / A-001；`00-meta` 复选框未勾，见 F-002） | D-003：加法路由、三键/pin/装配不变。A-001 self **pass**（0 required） |
| S5 本独立关门审 | 本条 | 无 high/med **required**；无到期 required 信息项 |
| D-002 §1 只读 + 禁止调他人账 | 满足 | 无自服务写面；owner 仅会话；管理端旁路仍分键 403 |
| D-002 §2 路由 / identity-only / 惰性开户 / 导航 | 满足 | `/my-wallet`；双 `/me` GET；无权限键；get-or-create + auto 审计；user 槽位 |
| D-002 §3 schema 页 | 满足 | intro + 3 statCard + 流水 table；金额 min-units（A-001 F-001 观察项仍成立，不阻断） |
| D-002 §4 管理端并存 / 装配 | 满足 | `/wallet` + 三键不变；fragment pages+user-nav；DefaultNavigationOrder 插入位置正确 |
| D-002 §5 测试面 | 满足（有推荐补测，见 F-001） | 三 handler 测试覆盖 401/editor 可读/开户幂等/alice-bob 账户隔离；kernel 快照；web fixture/哈希/页面清单/双语/D-VAL/T-01 槽位 |

## Findings

### F-001 · 身份隔离测试未钉「查询参数注入」与 bob 流水腿

| 字段 | 值 |
|------|-----|
| level | recommended（low） |
| status | open |
| evidence | `TestWalletSelfEntriesOwnScope`（`wallet_self_test.go` L120–180）断言 alice 流水与 bob **`/me` 余额**，**未** `GET` bob `/me/entries`，也**未**用 `?ownerId=` / `?accountId=<他户>` 打 `/me` 与 `/me/entries`。`user-menu.test.tsx` L130–138 移动抽屉只排除 Account/Settings，未断言不含 “My wallet”。实现可核对：handler 从不读这些查询键（`wallet_self.go` L35–81）；管理端旁路 403 已由 `TestWalletRoutesGates` 覆盖 |
| closure | — |
| 影响门禁 | **不阻断 S5**。安全结论来自代码路径 + 已有隔离/门禁用例，不是「未测即不存在」。建议后续补：忽略 `ownerId`/`accountId` 查询参数；bob entries 隔离；抽屉不含 my-wallet |

### F-002 · `00-meta` 派生进度 frontmatter 与检查点未与 S4 证据对齐

| 字段 | 值 |
|------|-----|
| level | recommended（low） |
| status | open |
| evidence | `00-meta.md` frontmatter `progress: 1/5`，正文与 `goal-tree.md` 为 **3/5**；S4/S5 复选框未勾。S4 事实已在 D-003 + A-001。本条**不**把台账滞后写成实现缺失；审计不得改这些字段 |
| closure | — |
| 影响门禁 | 不阻断。`/govern` 关门时应勾 S4/S5、按 P-001 重算 `5/5` 并同步 goal-tree |

## 必改项汇总

无 required / 必改项。

| ID | 级别 | 一句话 |
|----|------|--------|
| — | — | 无 |

recommended：F-001（补测查询注入 / bob 流水 / 抽屉）；F-002（关门时同步检查点与 progress）。

## 与既有意见的异同

| 意见 | 关系 |
|------|------|
| A-001（self · S2–S4 pass） | **同意**身份隔离、只读、惰性开户、装配与 go 判定。本条在独立复跑 + 管理端旁路逐条核对后维持同向 **pass**，不构成 P-004.2 冲突 |
| A-001 F-001（non-blocking · min-units） | 同意保留为观察项。`my-wallet.json` `format: plain` + i18n 已标注单位；F-011 前端金额格式化仍 deferred。不升格 |
| GOAL-020 A-003/A-005 | get-or-create UNIQUE + `isNew=false` + `auto:true` 审计被本目标复用；本条未重开 |
| GOAL-021 A-002（S5 independent pass 先例） | 同区同 auditor 关门口径：定向复跑 + 源码核对；无 required 即可放行；台账卫生记 recommended。本条沿用 |

无与本条 verdict 相反的相关意见冲突（P-004.2 不触发）。

## 结论 + 关门放行条件

**verdict: pass。** 自服务双端点身份隔离成立（owner 只来自会话；管理端 by-owner/entries/accountId 写读旁路仍要 `wallet.read`/`wallet.adjust`）；只读面无资金操作入口；惰性开户与 auto 审计幂等可核对；user-nav 槽位与 DefaultNavigationOrder/fixture/T-01 一致；D-003「不加不减」与组合根 30 权限 / 15 导航一致。无 high required，无到期 required 信息项。

**关门放行条件（给 /govern，本条不改状态）**

1. 无开放 required finding 需先闭合 → **可关门**（`status: done` 由编排器执行）。
2. 编排器应：勾选 S4/S5、按 P-001 将 progress 重算为 `5/5`、同步 `00-meta` frontmatter 与 `goal-tree.md`（消化 F-002）。
3. F-001 recommended 不阻断关门；可后续补测或接受为残余。
4. 波次级 e2e / V-007/V-008 仍按工作区惯例留批末，不构成本目标身份隔离门禁缺口。

建议用户下一步：`/govern` 响应 A-002 并办理关门。

## 声明

本意见不修改 status / progress / 方案正文 / goal-tree。响应由 /govern 处理。保证等级 L0。
