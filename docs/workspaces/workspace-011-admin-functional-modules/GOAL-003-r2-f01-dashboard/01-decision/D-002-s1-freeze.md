---
id: D-002
goal: GOAL-003-r2-f01-dashboard
title: S1 · 方案冻结 — admin.dashboard 模块（生产 home 仪表盘 + 必办-1 协议对照 + 必办-3 装配声明）
date: 2026-08-14
status: accepted
parent: GOAL-003-r2-f01-dashboard
created: 2026-08-14
updated: 2026-08-14
version: 0.1.0
---

# D-002 · S1 方案冻结（F-01 仪表盘/控制台 · 生产 Profile home）

> 依据：I-011-001 `3 F-01、`8 必办-1/必办-3；GOAL-003 00-meta 边界与 S1 门禁。

## 1. 模块设计：`admin.dashboard`

| 项 | 冻结值 |
|----|--------|
| 模块 ID | `admin.dashboard` |
| 依赖 | core.auth-session / core.navigation-capability / core.schema-render / core.operationlog |
| 能力 | StandardAdminCapabilities() |
| 页面 | `dashboard`（`/dashboard` 路由） |
| 导航 | `menu_dashboard`（sidebar，Order 0——首位）；可见性 PolicyAdminEditorViewer |
| 持久化 | 无迁移 |
| 权限键 | 无新增（页面数据源复用 users.read / roles.read；无管理操作） |

**页面内容**（schema 驱动，纯既有节点）：
- text 欢迎语（`dashboard 控制台`）
- grid（3 列）× 2 个 statCard：
  - 用户总数：`dataSource: /api/users`，`valueField: total`
  - 角色总数：`dataSource: /api/roles`，`valueField: total`
- **指标数据源选择（I-003 关闭）**：复用既有列表端点 envelope 的 `total`（statCard 渲染器原生支持 envelope 字段）——**不新增后端端点**，全部观众角色可读（users.read/roles.read 为 PolicyAdminEditorViewer）。不引入 API 聚合端点（避免新契约面）。

## 2. 必办-1 · 独立协议对照（v2.8.0 面）

| 对照对象 | 证据 | 结论 |
|----------|------|------|
| `grid-dashboard`（protocol-inventory `2.5 信息性场景） | 非语义权威（`docs/05-scenarios/grid-dashboard.md` 为范例候选） | **呈现自由**：不构成 dashboard 语义约束 |
| node.schema.json / page.schema.json | 无 dashboard 专用节点；`statCard`/`grid`/`section` 为 registry 显示节点（data-display 范例同构） | 页面用既有 registry 节点表达，无协议面新增 |
| 上游样例 `user-profile-*` / `order-list-*` | 非语义权威 | 与 dashboard 无关 |

**处置**：dashboard 为**呈现自由**（信息性场景无语义约束）+ **fail-open 留痕**（statCard 数据源加载失败 → 卡片级错误态，页面不整体崩溃；本页数据源均为既有端点，观众可读）。文档引用本对照。

## 3. 必办-3 · Profile 内容扩展声明（防误触「不改变 Profile 默认集」门禁）

- **`admin.dashboard` 加入 mvp + admin 默认启用集**（demo 经 mvp 继承）。定性：**Profile 内容扩展**——生产面 home 仪表盘是 F-01 交付目标本身（I-011-001 `3 明确「进入 mvp/admin 默认启用集」），经既有模块贡献机制（provider + fragment + reconcile）落地；**不改装配语义**：Manifest 聚合规则、协议 pin、共同门禁、capability 语义零改动。
- **`adminFunctionalOrder` 头部插入 `admin.dashboard`** → `deriveHomePageRef` 在 mvp/admin 下推导 home = `dashboard`（这是 F-01 的目标语义：生产 Profile 以仪表盘为 home）。demo 仍为 `overview`（dev.examples 优先规则不变）。
- 相应更新 `scripts/smoke.sh` SM-007 页面集：mvp = dashboard users roles account；admin = dashboard users roles settings activity account；demo = overview users roles account（+dashboard 不参与必需集断言？——demo 含 dashboard 页面，纳入 required）。

## 4. 路由与页面装配

- fragment：pages `dashboard`（titleKey `manifest.title.dashboard`）+ sidebar `menu_dashboard`（Order 0，`visibleWhen: $context.features.menu_dashboard == true`，`labelKey: manifest.nav.dashboard`，icon `dashboard`——iconRegistry 已有）。
- 页面 schema meta capabilities：`app.manifest` / `app.navigation`（statCard 属 registry 显示类型，无额外 capability——data-display 范例同构）。
- i18n：`manifest.title.dashboard` / `manifest.nav.dashboard` / `schema.dashboard.*`（intro/statCard labels）en/zh 同步。

## 5. 必办核对（I-011-001 `8）

| 必办 | 适用 | 处置 |
|------|------|------|
| **必办-1（协议对照）** | **适用** | **✅ `2** |
| **必办-3（home 装配声明）** | **适用** | **✅ `3**（贡献机制 + adminFunctionalOrder + 装配语义不变声明） |
| 必办-2/4/5 | 其它目标 | 不适用 |

## 6. 未选方案（留痕）

- 不新增后端聚合端点（`/api/dashboard/summary` 之类）——R2 用既有列表端点 envelope；聚合/图表中心归 R4 B-02。
- 不做 chart 节点（B-02 范畴）；R2 仅 statCard。
- 不在 demo 改变 home 语义（overview 保持）。

## 7. 实现范围（S2 清单）

1. `modules/dashboard/`：provider + fragment + schema（dashboard.json）。
2. kernel：`profileDefaults`（mvp/admin += admin.dashboard）+ BuiltinModules 描述符。
3. composition：provider 装配 + `adminFunctionalOrder` 头部插入。
4. smoke.sh SM-007 页面集；kernel_test/composition_test 计数更新。
5. i18n 键。
6. 测试：Go（composition home 推导 dashboard、permissions/nav 计数）+ Web（schema 结构校验含 dashboard 页、i18n 键）+ 回归。
