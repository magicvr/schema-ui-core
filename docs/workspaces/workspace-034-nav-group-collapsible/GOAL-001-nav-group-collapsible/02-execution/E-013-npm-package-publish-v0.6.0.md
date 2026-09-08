---
id: E-013-npm-package-publish-v0.6.0
parent_goal: GOAL-001-nav-group-collapsible
doc: execution-entry
status: recorded
date: 2026-09-08
created: 2026-09-08
updated: 2026-09-08
version: 0.1.0
---

# E-013 · npmjs.com 包更新与消费验证

## 已发生事实

1. npm 包版本更新通过 PR #12、#13、#14 与 #15 分阶段完成；各 PR 的 `r6-basic-matrix` 均全绿，合并后的 main CI 也全绿。
   - PR #12：初次受影响包版本推进。
   - PR #13：修正 `lib`/`protocol`/`ui` 的 package root exports 入口。
   - PR #14：补齐运行时 dependencies、peerDependencies 与 protocol `schemas/` 资产。
   - PR #15：将 CLI/QUICKSTART 的包面钉到最终 npm 版本，并把 Go tag 推进到 `apps/api/v0.6.1`。
2. 最终 npmjs.com 发布版本：
   - `@magicvr/schema-ui-lib@0.1.14`
   - `@magicvr/schema-ui-protocol@0.2.15`
   - `@magicvr/schema-ui-renderer@0.3.13`
   - `@magicvr/schema-ui-ui@0.1.11`
   - `@magicvr/schema-ui-shell@0.1.6`（已存在，幂等跳过）
   - `@magicvr/schema-ui-theme@0.1.4`（已存在，幂等跳过）
3. `npm view` 已复核上述六个版本在 `https://registry.npmjs.org` 可见。
4. 干净临时消费者执行 `npm install --ignore-scripts --no-audit --no-fund`，随后成功动态导入六个包：
   - `schema-ui-lib` 导出 41 项
   - `schema-ui-protocol` 导出 22 项
   - `schema-ui-renderer` 导出 10 项
   - `schema-ui-shell` 导出 41 项
   - `schema-ui-theme` 导出 5 项
   - `schema-ui-ui` 导出 18 项
5. 初次 npm 发布后，消费者冒烟发现旧包面存在真实元数据缺口：tsc 包的 root exports、runtime dependencies 与 protocol schemas 文件集不完整。已通过 PR #13/#14 递增 patch 版本修正；最终版本在干净消费者中通过，不把中间版本作为最终推荐包面。
6. CLI 当前版本钉已更新为 `apps/api/v0.6.1`，包面为 protocol `0.2.15`、lib `0.1.14`、renderer `0.3.13`、ui `0.1.11`、shell `0.1.6`、theme `0.1.4`。

## 证据边界

- npm token 只由发布脚本从本地 gitignored `.env` 注入，未输出或提交。
- 本记录只记录 npm registry 发布与消费验证事实，不修改 Goal/VP 状态；VP-034 Vision 层仍按 `/vision` 另行处理。
