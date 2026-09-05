---
doc_type: goal-audit
id: A-007-self-response-a006
parent: GOAL-005-r4-evidence-closeout
date: 2026-09-05
status: closed
version: 1.0.0
---

# A-007 · 响应 A-006（self · response）

## A-007 · 响应 A-006 运行时集成独立审计（2026-09-05）

- **source**：self（编排器响应记录，非独立审）
- **模式**：response · 响应 A-006（independent · codex local · verdict **fail** · 4 required）
- **verdict**：**pass**（作为响应记录；重新关门以 A-008 independent closure 复审为准）

### 前提承认

A-006 的核心事实**成立并接受**：Root 此前的关门依据（E-001 证据矩阵 + A-002/A-004 审计链）只覆盖了领域行为与台账投影，**未覆盖「从配置解析到组合根可达」的运行时集成**——`biz.digital-offer` 不在 `kernel.BuiltinModules()` 注册表中，生产配置路径无法启用模块，composition 分支不可达。按 P-003 撤回 Root/GOAL-005/VP-031 关门状态并整改。

### 关闭证据表

| Finding | 级别 | 闭合 | 证据 |
|---------|------|------|------|
| A-006 F-003 · `biz.digital-offer` 不在编译模块注册表，配置启用路径不可达 | high required | fixed | `kernel/profile.go` `BuiltinModules()` 增补 descriptor（与 `modules/digitaloffer.Provider.Descriptor()` 逐字段一致：7 路由 / 2 页面 / 2 导航 / 3 权限键 / digitaloffer fragment / 4 core 依赖 / StandardAdminCapabilities）；仍不进 mvp/admin 默认 profile。验收：`TestDigitalOfferPlanResolution`（enabled 含 channel.telegram 与 disabled 两种 plan 均可解析） |
| A-006 F-004 · Offer Delete 缺失，「CRUD」表述与实现不一致 | high required | closed（合同收窄，D-003 §F-004） | D-002 §2（v1.2.0）已冻结「**无删除**」：状态机 `draft → on_sale → off_sale` 即生命周期，下架保留历史凭证/权益引用。E-001 证据矩阵措辞同步为「Offer 生命周期管理（Create/Read/Update/Status；无删除，D-002 §2）」。不实现 Delete（既有合同条款构成 A-006 要求的书面收窄，无新增裁决） |
| A-006 F-005 · 0070 迁移策略未裁决 | high required | fixed（裁决：保留 compiled-global，D-003 §F-005） | D-003 记录裁决：保留本仓全部模块一致的 **compiled-global persistence**（`modules/compiled/persistence.go` 既有语义），「未启用模块亦建空表」为该策略的既定产品/运维语义；如未来不接受须单独立项改迁移 runner（A-006 方案 C，超出本工作区）。证据映射：fresh/reopen/幂等（`internal/store/migrate_test.go` 对含 0070 完整 catalog 断言 v70 tail + reopen no-op）、双库（GOAL-003/004 验收矩阵经完整 catalog 于 SQLite/真 PG 建库运行）、事务性（迁移 runner Apply+ledger 同事务）；快照 = 部署处置手段而非自动回滚（0070 仅正向，与全仓一致）；残余边界与复审触发已记录 |
| A-006 F-006 · 缺组合根到 Manifest/schema/HTTP 的真实验收 | med required | fixed | 新增 `internal/composition/composition_digitaloffer_test.go`：① `TestDigitalOfferPlanResolution`（enabled/disabled plan 解析）；② `TestDigitalOfferCompositionRoot`——真实 `newAppWithOptions` Fx 图 + mux：Manifest 聚合含 digitaloffer-offers/entitlements 两页面、`/api/schema/digitaloffer-*` 认证 200（body 含 pageId）/ 匿名 401 / 未知 pageId 匿名 401 · 认证 404、Admin 路由匿名 401 · admin 200、公开目录 200；Telegram-disabled 语义由 plan 不含 channel.telegram 的装配端到端证明（零 Bot API 依赖启动）。enabled-注册语义由 service 层真实 `telegram.NewDispatcher` 测试（`TestTelegramCommandsReply`）覆盖 |

### 治理状态

- 按本响应，Root/GOAL-005/VP-031 关门状态已**撤回**（active），投影同步；四项 required 全部 closed（fixed ×3 + contract-conformant ×1）。
- 待 A-008 independent closure 复审确认 open required = 0 后重新关门。

### 声明

本响应记录不修改 A-006 原文；`fixed` 证据以现行代码、测试与 D-003 裁决文本可核对。
