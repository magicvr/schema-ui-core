---
id: E-014-release-v0.7.0-preflight
parent_goal: GOAL-001-nav-group-collapsible
doc: execution-entry
status: recorded
date: 2026-09-21
created: 2026-09-21
updated: 2026-09-21
version: 0.1.0
---

# E-014 · apps/api/v0.7.0 发布候选预检（PR #16）

## 已发生事实

1. 发布 PR [#16](https://github.com/magicvr/schema-ui-core/pull/16) 目标为 `main`，最终候选 HEAD `971ebbce5256968e6e6c238c5c089473f5fddde9`；当前记录时尚未合并。
2. 候选版本钉：Go `apps/api/v0.7.0`；npmjs 新包面为 protocol `0.2.16`、lib `0.1.15`、renderer `0.3.14`、ui `0.1.12`；shell `0.1.6` 与 theme `0.1.4` 保持当前已发布版本。此前 npm registry 读取显示四个新 patch 版本尚未发布，shell/theme 当前版本已存在。
3. `node scripts/build-lib-packages.mjs`、`node scripts/rewrite-lib-aliases.mjs` 与 `node scripts/pack-npm-packages.mjs` 成功生成六个候选 tarball；产物位于忽略目录 `apps/web/dist-lib/artifacts/`。
4. 首次 `node scripts/publish-npmjs-packages.mjs --dry-run` 暴露内部名 `@schema-ui/lib` 被错映射到 `@magicvr/lib` 的问题。修复后的 `scripts/npmjs-package-name.mjs` 将六个内部/归档命名规整至 `@magicvr/schema-ui-*`；映射测试 4/4 通过，修复后 dry-run 对四个新版本显示正确公开名，对 shell/theme 幂等跳过。
5. PR #16 最终候选 Hosted CI run [35569624987](https://github.com/magicvr/schema-ui-core/actions/runs/35569624987) 的 9/9 job 均为 `success`，包括 API、真实 PostgreSQL、四组 browser E2E、两组 Compose smoke 与 Web。
6. 本条记录时尚未发生 PR merge、npm 实际发布、annotated tag `apps/api/v0.7.0` 或 GitHub Release。CI 与本地 dry-run 不代表这些外部动作已经完成。

## 证据边界

- npm token 未写入本条记录或输出；dry-run 不上传包。
- PostgreSQL helper 网络问题的代码级事实与修复证据另记于 workspace-040 `GOAL-005/E-005`；本条聚焦发布面与版本/包名预检。
- 不改变本目标或 VP-034 的 `status` / `progress`。
