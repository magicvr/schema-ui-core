---
title: S1 · API/Web 协议符合性候选矩阵
status: active
created: 2026-09-06
updated: 2026-09-06
parent: GOAL-041-w29-api-web-protocol-conformance
version: 0.1.0
---

# S1 · API/Web 协议符合性候选矩阵

> 所有条目在 S1 仅为 `collecting / pending-S2`。本表不能作为 `implementation-gap`、`upstream-protocol-gap`、`custom-extension-candidate`、`explicitly-out` 或 `excluded` 的最终 verdict；S2 必须补齐上游契约、当前实现和运行时证据，并完成 cross 方案审视。

| ID | 范围 | 当前事实 | S2 待回答 | 状态 |
|----|------|----------|-------------|------|
| C-001 | v2.9 provenance 路径与数量 | 上游 tag 有 `docs/schemas/*` 10 件及 `conformance/schemas/fixture-suite.schema.json` 1 件；本仓 provenance 将后者登记为 `docs/schemas/fixture-suite.schema.json`。字节 digest 一致，但来源路径不一致；note 同时声称 20 fixture suites（scenarios 未 vendor），实际 vendor/provenance 为 19。 | 是 provenance 文案/路径错误、允许的重定位，还是应增补同步脚本与第 20 套 fixture 的本仓实现缺口？ | collecting |
| C-002 | conformance claim generator | `generate-claim.mjs` 同一生成流程中，`report.pinnedUpstream.artifactVersion` 为 `2.8.0`，最终 `claim.protocolArtifact.artifactVersion` 为 `2.9.0`。 | report 字段是否为历史残留；若上游 claim schema 要求同一身份，应按协议修正本仓。 | collecting |
| C-003 | legacy provenance / tests | `upstream-fixtures.test.ts` 仍固定 v2.7 provenance，`stage3-fixtures.test.ts` 另固定 v2.9；两条入口并存。 | 明确 compatibility 回归与 current pin 的命名、claim 和门禁边界，排除双权威。 | collecting |
| C-004 | Manifest 协议版本 | API baseline 与 17 个 fragments 均发布 Manifest `2.7`；Web 适配器/claim 支持 2.9，35 个页面文档中 8 个声明 2.9。 | Manifest 2.7 是否是有意最低兼容合同；2.9 页面和 capability 如何在同一 Manifest 中合法协商。 | collecting |
| C-005 | capability 声明与实际使用 | `digitaloffer-entitlements` 声明 `data.route-binding` 但未发现 route expression；`digitaloffer-offers` 声明 `form.controls.readonly` 但未发现 `readOnly` 字段。 | 属于允许的保守声明、生成残留，还是应删除声明/补实现；mandatory suite 如何对应。 | collecting |
| C-006 | 静态 dogfood Manifest | Web `mvp-dogfood` 为 13 页、`admin-dogfood` 为 26 页；真实 profile 聚合为 MVP 6、Admin 22、Demo 14。fixtures 还重复保存 pageId/route/schemaUrl。 | fixtures 是历史演示快照还是运行时权威；是否应生成化、重命名或消除双来源。 | collecting |
| C-007 | 页面链路测试覆盖 | 35 页中 representative 测试仅列 10 页，并直接读取 API schema、注入 fixture fetcher。 | 哪些页面必须进入逐页 validator、`loadPageDocument → RenderPage` 和真实 handler 触达；如何定义 S5 全分母。 | collecting |
| C-008 | profile / optional runtime | 已用真实 `ForModulesWithFragments` 核到默认 profile 6/22/14 页；Telegram、Digital Offer 仅显式 custom 启用；未保留实际 HTTP Manifest 快照。 | 应以哪些 profile/显式模块组合形成验收矩阵，避免把 35 页静态宇宙写成单一生产 Manifest。 | collecting |
| C-009 | custom component 扩展面 | 15 个本地注册键：12 个用于页面 `body` custom node，1 个用于 action `content` custom node，2 个用于 `afterComponent`；当前使用键都有注册，但 key 未体现统一 namespace。 | 逐键核对上游 custom/extension 允许范围、schema/validator/capability/failure/compatibility/fixtures；必要时触发 P-004。 | collecting |
| C-010 | 未注册 custom 的失败语义 | 未知标准 node 以 `RENDER_UNKNOWN_NODE_TYPE` fail closed；未知 custom 只显示 inline fallback 文本。 | 上游是否明确未注册 custom 的 failure contract；当前行为是合法降级、实现偏差，还是上游缺少协议。 | collecting |
| C-011 | custom action allowlist | 5 个 schema custom handlers 均进入 Web 固定 allowlist；上游 `action.schema.json` 允许 `custom + handler`。 | 核验 namespace、权限、失败码、下载/预览/复制安全与 fixtures，而不是误当 custom node。 | collecting |
| C-012 | 内页与导航边界 | 多个 page 无 direct nav；`notifications.navigation.user` 为空；页面可能通过 `navigate`、通知铃或内页 route 可达。 | 逐项确认是合法 inner route/Host-owned 入口，还是 Manifest/navigation 断链。 | collecting |
| C-013 | API-only / Host-only surfaces | `admin.data-transfer`、`admin.mfa`、`admin.login-captcha` 等可能无独立 schema page；login/session/branding/shell/failure 为 Host 层。 | 结合 GOAL-004 v2.8 处置与当前代码，最终标为现有协议实现、`explicitly-out` 或新协议候选。 | collecting |
| C-014 | 历史 GOAL-004 证据边界 | 95 个 Host/App 候选在 v2.8 已处置，但当前模块、35 页、15 custom 和 v2.9 增量已扩大。 | 只复用历史处置语义和边界；不得用 95/95 或旧 fixture 绿灯替代当前逐项核验。 | collecting |

## 证据入口

- 上游分母：`S1-protocol-denominator-v2.9.md` / `S1-protocol-denominator-v2.9.json`。
- 页面控件目录：`S1-api-web-page-control-catalog.md` / `S1-api-web-page-control-inventory.json`。
- 当前源码关键链：`apps/api/internal/composition/composition.go:729-766`；`apps/web/src/app/App.tsx:530-609,633-671`；`apps/web/src/protocol/load-page.ts:69-152`；`apps/web/src/renderer/render.tsx:328-353,2877-2956`。
- 历史边界：`docs/workspaces/workspace-010-design-implementation-conformance/GOAL-004-w3-schema-host-protocol-conformance/`；历史 verdict 不作为本表的现行分类。

## S1 结论边界

- 已完成：协议 identity、机器协议分母、页面/控件静态分母、真实聚合函数的 profile 投影与候选登记。
- 未完成：候选最终分类、上游增补报告、custom 用户裁决、产品代码修正、真实 API/权限/数据/浏览器验收。
