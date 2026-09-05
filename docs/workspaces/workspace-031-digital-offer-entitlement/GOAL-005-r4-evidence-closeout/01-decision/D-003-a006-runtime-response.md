---
doc_type: goal-decision
id: D-003-a006-runtime-response
parent: GOAL-005-r4-evidence-closeout
date: 2026-09-05
status: accepted
version: 1.0.0
---

# D-003 · 响应 A-006：模块注册修复与迁移策略裁决

## 背景

A-006（independent · runtime-integration · fail · 4 required）在 Root 关门后追加了运行时集成审计。其核心事实成立：`biz.digital-offer` 的 Provider descriptor 此前未进入 `kernel.BuiltinModules()`，而 `composition.ResolvePlan` 只以该列表建 registry——生产配置路径（profile / app.modules）无法解析启用该模块，composition 的 `plan.HasModule` 分支为不可达代码。本决策为 A-006 的四个 required 给出合法闭合路径并落盘裁决。

## 裁决

### F-003（模块注册表缺失）→ fixed

`kernel.BuiltinModules()` 增补 `biz.digital-offer` descriptor（与 `modules/digitaloffer.Provider.Descriptor()` 逐字段一致：路由/页面/导航/权限/fragment/依赖/capability）。仍**不进 mvp/admin 默认 profile**——启用仅经 `app.modules`（custom preset/list）。验收：`TestDigitalOfferPlanResolution`（enabled/disabled 两种 plan 形状解析）。

### F-004（Offer Delete 缺失）→ 按既有合同收窄（无需新裁决）

D-002 §2（GOAL-002 R1 合同，v1.2.0）已冻结：「**无删除**（保留历史凭证/权益引用）」；状态机 `draft → on_sale → off_sale`，下架即生命周期终点。A-006 要求的「书面收窄」由此前已接受的合同条款构成，本次仅同步措辞：VP-031 判据 1 / 证据矩阵中的「Offer CRUD」统一表述为「Offer 生命周期管理（Create/Read/Update/Status；无删除，D-002 §2）」。不实现 Delete。

### F-005（0070 全局迁移策略）→ fixed（裁决：保留 compiled-global）

- **裁决**：保留本仓既定的 **compiled-global persistence** 策略——`modules/compiled/persistence.go` 对全部模块（wallet/telegram/account…）一致地全局应用迁移，`0070 digital_offers` 与先例完全同构。**未启用模块也会创建 schema** 是该策略的既定产品/运维语义（平台级，非本工作区可静默更改；如未来不接受，须单独立项修订迁移 runner 并 cross audit——A-006 方案 C）。
- **证据映射**：fresh/reopen/幂等 = `internal/store/migrate_test.go`（对含 0070 的完整 catalog 断言 v70 tail 与 reopen no-op）；SQLite/真实 PG 双库 = GOAL-003/004 验收矩阵经完整 catalog 建库运行（`TestPurchasePostgresAcceptance` 等）；Apply 事务性 = `internal/store` 迁移 runner 语义（Apply 与 ledger 同事务）。快照恢复明确记录为**平台部署处置手段而非自动回滚**——0070 仅正向、无 Down migration，与本仓全部迁移一致。
- **残余边界（记录）**：未启用模块的 dormant 空表存在于所有实例；触发复审条件 = 首次生产启用、schema 变更、生产部署或多实例。

### F-006（组合根验收缺失）→ fixed

新增 `internal/composition/composition_digitaloffer_test.go`：
1. `TestDigitalOfferPlanResolution`：enabled（含 channel.telegram）/disabled 两种 plan 解析；
2. `TestDigitalOfferCompositionRoot`：真实 `newAppWithOptions` Fx 图 + mux——Manifest 聚合含两页面、`/api/schema/digitaloffer-*` 认证 200 / 未认证 401 / 未知 pageId 匿名 401·认证 404、Admin 路由 401/200、公开目录 200；Telegram-disabled 语义由 plan 不含 channel.telegram 的装配路径端到端证明（应用零 Bot API 依赖启动）。

## 影响

- Root/GOAL-005 关门状态**撤回**，待 A-006 四项 closed 后经 independent closure 复审重新关门。
- 无合同语义变更（F-004 为既有条款的表述同步；F-005 为既定平台策略的显式裁决）。
