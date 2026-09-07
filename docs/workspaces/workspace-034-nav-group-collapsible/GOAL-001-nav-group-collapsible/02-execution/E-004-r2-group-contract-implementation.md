---
id: E-004-r2-group-contract-implementation
doc: execution-entry
parent_goal: GOAL-001-nav-group-collapsible
status: recorded
date: 2026-09-07
created: 2026-09-07
updated: 2026-09-07
version: 0.1.0
---

# E-004 · R2 分组注册与 Manifest 聚合契约实施

## 已发生事实

1. `kernel.NavigationContribution` 增加可选 `NavigationGroup` 元数据（key/order/label/labelKey/icon）；kernel finalize 对同 key 元数据执行精确全等校验，并以 `CodeModuleNavigationGroupConflict` fail closed。
2. 当前五组的 17 个 sidebar 导航贡献已写入对应模块 Provider；top/user slot、Dashboard、通知铃面和 `dev.examples` authored Examples 组未被错误归组。
3. `internal/manifest` 增加结构化 group 映射与 sidebar 归一化：Dashboard 第一、显式 GroupOrder、未分组链接兼容、Examples 等既有协议组保留、sidebar-only mismatch fail closed；现行协议不增加 group key/id。
4. `internal/composition/composition.go` 与 `apps/api/server/serve.go` 两条 Manifest assembly 均传递结构化 group 贡献，未保留分叉路径。
5. 新增 kernel group conflict、Manifest group order/mixed-output/sidebar-only、composition default/override、optional Digital Offer/Telegram 递归 pageRef 回归；新增 en-US/zh-CN group labelKey 文案。
6. 尚未实施 R3 Shell 折叠/展开、键盘交互、直接 URL 自动展开或状态保持；未宣称 R2 之后的检查点完成。

## 验证事实

- `go test ./... -count=1`：通过。
- Web 直接 Vitest：97 files / 1332 tests 通过。
- Web 直接 `tsc -b`：通过。
- Web 直接 `vite build`：通过；仅保留既有 chunk size warning。
- `git diff --check`：通过。
- pnpm 包装入口仍受环境 `esbuild` ignored-build-script 保护提前退出；未将该环境问题误报为代码失败。

## 待审计

本 E 条目记录实施事实；R2 检查点在 self implementation audit 与项目指定的本地 grok build independent implementation audit 完成后再更新路线图/进度并创建 Git checkpoint。
