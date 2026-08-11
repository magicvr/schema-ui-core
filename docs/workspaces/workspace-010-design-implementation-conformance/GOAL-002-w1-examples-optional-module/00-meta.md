---
id: GOAL-002-w1-examples-optional-module
title: W1 · 范例/演示产品面可选模块化
status: done
parent: GOAL-001-design-implementation-conformance
created: 2026-08-11
updated: 2026-08-11
version: 0.1.0
progress: 6/6
---

# GOAL-002 · W1 · 范例/演示产品面可选模块化

## 概述

本子目标是 VP-010 / workspace-010 的**首波整改**：把当前挂在 `core.schema-render` 上的**协议范例/演示页面与 Examples 导航**，从「伪核心、双 Profile 强制、依赖图绑死」纠正为与标准功能模块同构的**可配置启用/注销（可选模块）**形态，使生产场景默认不暴露演示产品面，开发/演示场景可显式启用。

### 已核实的符合性缺口（as-designed vs as-built）

| # | 设计意图（摘要） | as-built 偏差 |
|---|------------------|---------------|
| G1 | 标准能力按 Profile / `modules.enabled` 选择；生产不需要的面应可关 | `core.schema-render` 在 mvp/admin **默认强制**；组合根 **无条件** `schemarendermodule.New()` |
| G2 | 模块贡献自有 Manifest fragment；baseline 不塞业务/演示页 | 范例 pages + Examples nav 写在 `core.manifest-route` baseline `app-manifest.json` |
| G3 | 真业务模块不依赖演示包 | `admin.users` / `roles` / `settings` 与 `core.manifest-route` **DependsOn** `core.schema-render` |
| G4 | home 指向可交付产品面或可配置首页 | `homePageRef: overview` 绑在演示 overview 页 |
| G5 | Web Shell 不硬编码业务路由（已基本满足） | 运行时干净；问题在供给/装配侧，不在 Renderer 中央注册 |

范例 pageId 集合：`overview`、`data-table`、`admin-list-batch`、`data-display`、`search-form-table`、`form-controls`、`form-with-reactions`、`form-with-upload`。  
无独立 demo API 模块；部分范例复用真实 API（如 `/api/users`）——注销页面即可，不必虚构 demo 后端。

## 成功标准

- [x] **S1 · 能力拆分**：演示页/Examples 导航从「强制 core 能力」拆出；真 Schema 贡献协议/能力不再与 demo page 集合绑死同一不可关模块语义（`dev.examples` 持有 8 页；`core.schema-render` 仅能力壳；E-004）
- [x] **S2 · 可选模块形态**：存在可 `plan.HasModule` 开关的模块（命名已冻结 = **`dev.examples`**）；自有 schema + manifest fragment；组合根按启用集装配（`composition.go` HasModule + dev/examples 包；E-004）
- [x] **S3 · 依赖剪枝**：`admin.*` 与 `core.manifest-route` 不再 DependsOn 演示模块；禁用演示模块时 Admin 仍可启动并发布 Manifest（`profile.go` BuiltinModules；E-004）
- [x] **S4 · Profile 默认**：生产向 `mvp`/`admin` **默认不启用**演示模块；开发/dogfood 可通过显式 modules 列表或专用 profile 启用（`profileDefaults` 不含 `dev.examples`；E-004）
- [x] **S5 · 产品面与 home**：禁用时 Manifest 无 Examples 组、无 8 范例 pageId、schema 404；`homePageRef` = 首个启用的 admin 功能页（无则任意首页；无页则省略）；`dev.examples` 启用时 = `overview`（`deriveHomePageRef` + `StampHomePageRef`；`TestManifestHomePageRefDerivation`；E-004）
- [x] **S6 · 回归与 go 接口**：API/Web 相关测试与双 Profile 烟测通过（`go test` 全绿；web 746 测试；mvp/admin e2e 3+3 passed）；记录对 VP-008 `go` 消费有效性的影响（**已暂挂**，E-004 §go）；架构/playbook 核验无需回贴（E-004）

## 高层路线图（P-001）

1. **方案冻结**：模块 id、homePageRef 策略、Profile 默认、是否保留 overview 为可选首页、测试分母调整。 **（2026-08-11 完成 · D-002 + 实施冻结附录 D-003，cross 审计 R1–R4 闭合）**  
2. **拆分与迁移**：页面文档/fragment 归属；baseline 瘦身；BuiltinModules / providers / DependsOn。 **（2026-08-11 完成 · E-004）**  
3. **回归**：composition/manifest/profile 测试 + web 代表路径；更新 i18n 仅作非阻断清理（可 residual）。 **（2026-08-11 完成 · E-004：go/web 全绿 + 双 Profile e2e）**  
4. **go 影响留痕**：矩阵变更说明 + freshness/重验证指针。 **（2026-08-11 完成 · E-004 §go 暂挂触发；恢复证据清单）**  
5. **波次审计**：self（必要）+ 若触及模块矩阵/装配语义则按风险 `cross`/`independent`（P-004）。 **（待办 · D-003 已定 cross）**

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 可选模块稳定 id 与是否进入 compiled 候选默认集 | 方案冻结 | 方案 | 用户确认或方案 D 裁决 | verified | — | `dev.examples`；进 compiled 候选、默认 Profile 不含（用户 2026-08-11 确认；D-002） |
| I-002 | required | 禁用演示后 `homePageRef` 目标（users / 首个 admin 页 / 可配置） | 方案冻结 | 方案 | 用户确认 | verified | — | 首个启用的 admin 功能页（用户 2026-08-11 确认；D-002） |
| I-003 | required | mvp/admin 是否**默认关闭**演示（强烈建议是） | 方案/Profile | 方案 | 用户确认 | verified | — | 默认关；dogfood/`APP_MODULES_ENABLED` 显式开（用户 2026-08-11 确认；D-002） |
| I-004 | non-blocking | web i18n 范例 key 是否本波删除 | 验收整洁度 | 验收 | 可 residual 留 key | open | deferred：不阻断功能；owner=本波；复核=S6 | 可接受残留死 key |

## 父目标

- [GOAL-001-design-implementation-conformance](../GOAL-001-design-implementation-conformance/00-meta.md)

## 台账布局

本目标从首条记录起使用 `01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger；索引与目录条目共同构成正式记录。
