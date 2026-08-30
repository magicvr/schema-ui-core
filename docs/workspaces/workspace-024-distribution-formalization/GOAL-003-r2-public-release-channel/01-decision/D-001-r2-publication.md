---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-003-r2-public-release-channel
version: 0.1.0
---

# D-001 · 发布面定案（2026-08-29 · 用户裁决落盘）

## 决策

1. **公开 scope = `@magicvr` 先行**（用户裁决 2026-08-29；初选 `@schema-ui` 因账号非同名 org 不可直接发布，变更记录见 §6）：npmjs.com 发布
   `@magicvr/schema-ui-{protocol,lib,theme,ui,renderer,shell}`；版本沿用
   （protocol/renderer 0.2.0 · lib/theme/ui/shell 0.1.0）。GH Packages
   旧包（`@magicvr/schema-ui-*` 私有）保留不删（历史消费面）。
2. **真实发布授权**（用户裁决 2026-08-29）：token = 仓库根 `.env` 的
   `npm_token` 键（`.env` 已在 `.gitignore`，不入库）。
   **凭据处置**：发布脚本运行时读取 `.env` 注入**临时** `.npmrc`
   （`//registry.npmjs.org/:_authToken=…`），仅存在于 stage 临时目录，
   成功后随 stage 删除；不落盘、不打印 token。
3. **Go 侧**：`apps/api/v0.4.0` origin tag（含 serve 面 = `server` 包 +
   CLI `serve` 子命令 + 新模板）；`go get`/`go install` 走公共 proxy。
4. **scope 迁移 = npm 消费面 breaking**：changelog 注记
   `@magicvr/schema-ui-*`（GH Packages 私有）→ 公开发布面（见 §6 变更记录）；
   冻结面 v1.3.0 的 npm 导出面以公开发布 scope 为准，peer 矩阵定稿随 R5（v1.4.0）。
5. **审计模式**：S2（真实发布）与 S4（关门）实证门禁 = independent
   （grok build 先例）；grok 不可用时 self + 用户确认（不冒充 independent）。

## 6. 变更记录（v0.1.1 · 2026-08-29 · 用户裁决）

- 初选 `@schema-ui` 公开 scope；实测 `npm whoami` = account `magicvr`，npmjs 新 scope
  `@schema-ui` 需同名组织（网页端创建 + owner 授权）——**用户裁决：先用 `@magicvr`
  scope 实发 npmjs 公开包**（零等待；GH Packages 私有同 scope 与 npmjs 公开包为不同
  registry 独立身份，无冲突）；`@schema-ui` 保留为 org 就绪后的正式化候选
  （届时新包名 + golden-field 迁移，或沿用 `@magicvr` 定稿）。
- 消费迁移说明：既有消费者移除 GH Packages registry 映射（web/.npmrc）后，
  `pnpm add @magicvr/schema-ui-*` 即 npmjs 公开版。

## 未选方案

| 项 | 未选 | 理由 |
|----|------|------|
| scope | `@magicvr/schema-ui-*` 复用 | 公开品牌辨识度低；公开/私有面同 scope 易混淆（GH Packages 私有包保留） |
| 发布 | 仅 dry-run + 决策文档（VP-023 I-003 先例） | 用户已给 token 授权真实发布；判据 #2 字面要求 registry 上传实证 |
| Tag | v0.4.0 前先合并其它 | serve 面是 v0.3.0 后唯一功能增量；独立发布保持演进面清晰 |

## 信息门禁

- I-024-001 → **verified**（用户裁决：@magicvr 定稿 · 真实发布 · npm_token 注入点）。

## 7. 定稿裁决（v0.1.2 · 2026-08-29 · 用户）

**scope 定稿 = `@magicvr`（npmjs 公开）**：用户裁决「维持当前发布的发布形式，不再改为 `@schema-ui/xxx` 形式」——§2/§6 中的 `@schema-ui` 正式化候选**取消**；计划文档（VP-024 判据 #2 · GOAL-003 · 本文件）已随之定稿。Charter 成功边界 #1 的 npm 包组占位 `@schema-ui/*` 经 editorial（VR-051）对齐为 `@magicvr/schema-ui-*`。