---
id: E-015-r6-c64-acceptance-freeze
doc: execution-entry
goal: GOAL-013-r6-old-path-removal
source: orchestrator
date: 2026-08-06
status: recorded
---

# E-015 · R6 C6.4 验收矩阵冻结与回归基线

## 已发生事实

- 完整读取现行 Charter `schema-ui-core-admin-foundation@0.2.0`、alignment、workspace
  protocol、VP-003、workspace-003/goal-tree、Root 与 GOAL-013 当前台账；工作区绑定、
  `primary_plan`、VP `vision_ref` 与 Vision Review 门禁均合法。
- 扫描现有 API/Web/E2E/Compose/smoke/workflow 与迁移恢复测试，D-004 固定八组 C6.4
  验收，并逐条映射 VP exit #1～#7。
- 工具可用性已核对：Go `1.26.0`、Node `22.17.0`、npm `10.9.2`、Docker
  client/server `29.6.2`、Compose `v5.3.1`。
- 方案冻结前运行 `apps/web` `npm test`，结果为 `481/495` 通过、14 失败。三组测试仍
  读取已删除的 `apps/api/internal/handler/fixtures/schema/`，失败均为 `ENOENT`；该事实
  已登记为 C64-V01/V03 required 修复。
- 文档一致性扫描确认 README/QUICKSTART 仍含错误 `MODULES_ENABLED`、旧 handler Schema
  路径、Web public 静态 Manifest、中心 store seed 与本地 `.env` 自动加载假设；Compose/
  workflow 只覆盖默认 mvp。这些均进入 D-004 fixed 路线。

## 状态边界

- 本条只记录方案冻结与失败基线；尚未完成任何 C6.4 required 验收。
- R6-I004 保持 `collecting`，C6.4 不勾选，GOAL-013 保持 `active / 3/4`，Root 保持
  `active / 5/6`。
- 当前结果是本地 dirty-worktree snapshot，不是 hosted CI、合并或发布证据。

## 下一步（计划）

1. 修复 Web fixture ownership 与静态 Manifest test fixture，修正文档/env/Compose，并
   固化双 Profile CI/local matrix。
2. 全量执行 D-004 C64-V01～V07，形成逐条 evidence package。
3. 依次完成 GOAL-013 self、Grok independent 与 `/govern` 响应；之后再进入 Root
   close-out。
