---
id: E-003
goal: GOAL-022-my-wallet-self-service
title: S2 实现完成（自服务 API + schema 页 + user-nav 入口）
date: 2026-08-16
status: active
parent: GOAL-001-admin-functional-modules
created: 2026-08-16
updated: 2026-08-16
version: 1.0.0
---

# E-003 · S2 实现

## 事实

按 D-002 实现，全部落盘：

- **自服务 API**（`apps/api/internal/handler/wallet_self.go`，新文件）：
  - `GET /api/wallet/me`：会话身份 get-or-create + 账户摘要（resourceList 信封，供 statCard 绑定）；首次开户审计 `wallet.account-create`（detail 含 `"auto":true`）。
  - `GET /api/wallet/me/entries`：会话身份 get-or-create + 本人流水分页。
  - 两个端点均 **identity-only**：不接收客户端 ownerId，owner 恒由会话推导（`selfIdentity`）；无权限键（与 /api/account/profile 同款）。
- **模块接线**（`apps/api/internal/modules/wallet/provider.go`）：Descriptor 路由/pages/navigation 同步；Register 增 WalletSelfRoutes、my-wallet PageContribution（Actions: list；DataSource: /api/wallet/me）、NavigationContribution `menu_wallet_self`（Order 2、PolicyAdminEditorViewer、无权限键）。
- **内核**（`apps/api/internal/kernel/provider.go` + `profile.go`）：DefaultNavigationOrder 插入 `menu_wallet_self`（menu_account 之后、menu_activity 之前）；BuiltinModules admin.wallet 路由/pages/navigation 同步。
- **Manifest**（`modules/wallet/manifest/fragment.json`）：pages 增 my-wallet（/my-wallet）；navigation.user 增 menu_wallet_self 项（icon wallet，visibleWhen `features.menu_wallet_self`）。
- **Schema 页**（`modules/wallet/schema/my-wallet.json` + schema.go embed/PageIDs/SchemaDocuments）：intro text + 3 statCard（总余额/可用/冻结 ← /api/wallet/me）+ 流水 table（← /api/wallet/me/entries，列复用 walletEntries labelKey）。
- **i18n**：en-US/zh-CN 增 `manifest.title.myWallet`、`manifest.nav.myWallet`、`schema.myWallet.text.intro`、`schema.myWallet.statCard.{total,available,frozen}`。
- **测试**：
  - `handler/wallet_self_test.go`（新）：`TestWalletSelfIdentityOnly`（401/editor 无权限键可读）、`TestWalletSelfAutoCreateAndIdempotency`（开户+auto 审计+幂等）、`TestWalletSelfEntriesOwnScope`（本人隔离：alice 2500 vs bob 1000 互不可见）——3/3 PASS。
  - `kernel/navigation_order_test.go`：快照 +menu_wallet_self。
  - `composition_test.go`：`TestSystemDataReconcileUsesFinalizedProfileContributions` admin menu_items 14→15。
  - web：admin fixture 增 my-wallet 页 + user-nav 项（哈希重钉 `0efb4054…`）；app-manifest.test 页面清单 +my-wallet；schema-keys SCHEMA_FILES +my-wallet.json；user-menu.test 内联 manifest 增 my-wallet（T-01 槽位断言 `Account → My wallet → Settings → Sign out`）。

## 验证（S3 · E-004 见下）

- go build ./... 全绿；`go test -count=1 ./...` 全绿（34 包 ok）。
- vitest 全量 **65 文件 / 1038 测试全绿**；tsc --noEmit 无错。
- **实机冒烟**（临时 admin profile 实例，端口 25099，临时 DB）：manifest 发布 my-wallet 页 + user-nav 顺序 `account → my-wallet → settings`；/api/schema/my-wallet 200；admin 登录后 /api/wallet/me 返回信封（ownerId=user-admin、余额 0、active）+ 恰好 1 条 auto 审计；/me/entries 空列表；admin /api/wallet/accounts 回归正常；editor（无 wallet.* 权限）可读 /me 但 /api/wallet/accounts 403。