---
id: GOAL-005-r3-admin-shell-navigation
title: R3 · Admin 外壳与导航
status: active
parent: GOAL-001-mvp-admin-foundation
created: 2026-07-31
updated: 2026-07-31
version: 0.2.0
---

# GOAL-005 · R3 · Admin 外壳与导航

## 概述

在 `apps/web` 的 R1 React 工程骨架之上，交付由 App manifest 驱动的 Admin 外壳、导航入口与路由边界。当前工作树已包含 manifest loader、2.7 子集校验、导航投影、History 路由和 shell；目标仍保持 `active`，等待实施阶段自审与关门记录。

范围依据：Root 已将 R3 定义为“Admin 外壳与导航”；协议资料将 `D-APP` 映射为 React 侧的“装载与导航壳”，并固定 `S-09`、App Manifest schema 及 `app-manifest` / `app-navigation` fixture 的上游路径。固定资料版本为 `schema-ui-docs` artifact `2.7.0`，source commit 为 `ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b`。

## 范围边界

### 纳入

- App manifest 的来源、最小可用子集、装载入口和失败边界的决策与实现。
- Admin shell 的页面级布局、导航容器和应用入口边界。
- manifest 到导航项/路由入口的映射，以及默认路由、fallback、active-route 语义的决策与实现。
- 面向 `app-manifest` / `app-navigation` 的结构或行为验证路径，以及 shell 集成验证证据。

### 排除

- R4 的账号、权限模型与权限继承实现；父目标的 `I-PROTO-002` 仍为 R4 的开放 required 门禁。
- R5 的协议 Renderer 全量、业务域范例页和并行 `mock-app` 业务演示；父目标的 `I-PROTO-003` 仍为 R5 的开放 required 门禁。
- “完整协议支持”或完整 conformance 主张；上游 excluded reference/runner 不能单独作为兼容证明。
- 完整主题/设计系统产品化，以及与 R3 无关的业务页面实现。

## 高层路线图（P-001）

1. **契约发现与信息就绪**：**完成**；五项 `I-005-*` 已由固定 artifact、代码和测试证据验证。
2. **方案冻结**：**完成**；D-005 记录 manifest 子集、路由、导航和 shell 边界。
3. **R3 实施**：**完成**；工作树已落地 loader、Admin shell、导航投影和 History 路由，保持 R4/R5 边界。
4. **验证与关门**：**进行中**；固定 fixture、集成测试、静态 manifest loader 和构建已复核，待实施阶段 self audit 与最终台账闭合。

以上状态只记录当前可核对事实；`active` 仍表示关门审计尚未完成。

## 成功标准

- [ ] R3 必需信息项已由证据验证，或有用户书面接受的有界 residual；未知项没有被默认为已知。
- [ ] manifest 装载使用已冻结的 schema/版本/最小子集，并有可核对的无效输入或装载失败处理边界。
- [ ] Admin shell 能由已冻结的导航数据进入页面，默认路由、fallback 和 active-route 语义可通过测试或运行时证据复核。
- [ ] `app-manifest` / `app-navigation` 的验证路径已执行并记录结果；未用 excluded reference/runner 替代语义验证。
- [ ] R4 权限、R5 Renderer/范例和完整协议支持仍保持边界外；目标关门前无开放 required finding。

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-005-001 | required | App manifest 的最小 schema、来源、版本和本地校验入口是什么？ | 方案冻结 / 实施 | 方案冻结前 | 对照 `S-09`、`docs/schemas/app-manifest.schema.json`、`upstream/provenance.json` 和固定 artifact；采用仓库内 pinned copy + sha256 校验，不在运行时依赖远程请求 | verified | 不延期；后续协议升级需重新核验 | `APP_MANIFEST_PROTOCOL_VERSION=2.7`；`validateAppManifest()` / `loadAppManifest()` fail closed；本地 schema/fixture hash 与 provenance 测试固定；R3 明确不主张完整 schema/conformance |
| I-005-002 | required | manifest 条目如何映射为导航项、路由入口和层级关系？ | 方案冻结 / 实施 | 方案冻结前 | 对照固定 app-navigation fixture，验证 `top`/`sidebar`/`user`、group、pageRef/url 和可见性投影 | verified | 不延期；导航字段变化需重新核验 | `navigation.ts` 的 `projectNavigation()` 映射三 slot；group 只保留可见 child；`upstream-fixtures.test.ts` 执行 16 个 navigation cases |
| I-005-003 | required | 默认路由、未知路由和 fallback 页面/行为分别是什么？ | 方案冻结 / 实施 / 验收 | 方案冻结前 | 从固定 fixture 和 shell 集成测试收集 `/`、deep link、unknown route、manifest failure 行为 | verified | 不延期；fallback 文案或 home 规则变化需重新核验 | `resolveInitialRoute()` 将 `/` replace 到 `homePageRef`；已知 deep link 优先；未知路由保留路径并显示 shell fallback；App 集成测试覆盖返回 home |
| I-005-004 | required | active-route 的来源、匹配优先级、重定向和 URL 语义是什么？ | 方案冻结 / 实施 / 验收 | 方案冻结前 | 验证 D4a literal count、route length、declaration order、query 忽略和 History API 行为 | verified | 不延期；路由模板或 URL 规则变化需重新核验 | `matchRoute()` 与 `projectNavigation()` 提供 active；`App.tsx` 使用 `replaceState` 处理 root、`pushState` 处理站内导航；参数 pageRef 只有可绑定时提供具体 href |
| I-005-005 | required | Admin shell 的产品边界包含哪些固定区域，哪些留给业务页或后续产品化？ | 目标方案 / 实施 | 方案冻结前 | 对照 Charter、VP、R3 边界并核对 shell 集成测试；明确 R4/R5 留白 | verified | 不延期；R4/R5 进入时需新目标门禁 | 固定区域为 header/app identity、top/user nav、responsive mobile nav、desktop sidebar、main page surface 和 route fallback；page renderer、身份来源、权限产品化与完整主题系统留在 R4/R5/后续范围 |

父目标的 `I-PROTO-001=verified` 仅表示 R2 冻结范围；不解除本目标的未知项，也不代表 R3-R5 已实施或完整协议已支持。

Root `I-PROTO-004` 仍是父目标的 `open` / `non-blocking` 工程策略项，不会被本目标静默升格为独立 required。它与 `I-005-001` 的依赖已显式登记：关闭 `I-005-001` 时必须说明采用 vendor 还是 pin 远程校验，以及相应失败边界。

## 父目标

- [GOAL-001-mvp-admin-foundation](../GOAL-001-mvp-admin-foundation/00-meta.md)

## 备注

本目标的 `status: active` 表示 R3 实施事实已产生，但在阶段 self audit、required finding 闭合和关门记录完成前不表示目标已 `done`。
