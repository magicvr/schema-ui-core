---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-003-r2-public-release-channel
version: 0.1.0
---

# E-001 · 目标建立（2026-08-29）

1. **立项**：承接 Root 纲领 R2（公开发布通道 · VP-024 判据 #2 · go 后清单 ② · R1 残余 2 发布前提）；goal-tree 同步。
2. **发布面核验（设计依据）**：
   - `scripts/publish-npm-packages.mjs` 已参数化（`PUBLISH_SCOPE` / `NPM_REGISTRY`），GH Packages 直发路径可用；npmjs 路径需注入 `//registry.npmjs.org/:_authToken` 型临时 .npmrc；
   - 仓库根 `.env`（`.gitignore` 已排除）含 `github_token` / `npm_token` 两键（值未读取）；
   - `apps/web/dist-lib/artifacts/` 六份 tgz 就绪（schema-ui-{lib,protocol,renderer,shell,theme,ui}）。
3. **信息门禁**：I-024-001 verified（用户裁决：`@schema-ui` scope · 真实发布 · token 注入点）。