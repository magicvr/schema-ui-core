---
id: GOAL-002-r1-repo-layout-conventions
doc: execution
status: active
parent: GOAL-001-mvp-admin-foundation
created: 2026-07-31
updated: 2026-07-31
version: 0.2.0
---

# 执行记录 · GOAL-002

## 时间线

### 2026-07-31 · 立项

- `/govern` 用户确认 R1 脚手架取舍（见父目标 D-004）后创建本目标五件套。
- **未做**：尚未改 `docs/architecture/` 或根 README；尚未创建 `apps/*` 占位。

### 2026-07-31 · 响应 A-001 并落盘约定

- 记录 **D-002**：`apps/*` 创建权方案 A；运行入口两档成功标准；闭合 independent A-001 F-001/F-002（见 `03-audit` A-002）。
- 落盘 [docs/architecture/monorepo-layout.md](../../architecture/monorepo-layout.md)（应用 monorepo / 包管理 / 命令契约）。
- 更新 [docs/architecture/directory-layout.md](../../architecture/directory-layout.md) v0.7.0（挂 `apps/*` + 链 monorepo 约定）；同步 `skills/core` 镜像。
- 新建根 [README.md](../../../README.md)：文档入口 + 运行契约链到 monorepo 与 app README。
- **未**在本目标创建可运行 `apps/*` 工程树（服从 D-002）。

## 待办（计划 · 非完成事实）

1. 阶段自审 / 可选 `/audit` 后评估关门。
2. 003/004 README 稳定后可再收紧根 README 链接表述（已预留）。

## 进度评估

**约定文档已落盘**；成功标准必达项对照见 `00-meta` 勾选。未标 `done`（建议阶段审后再关门）。
