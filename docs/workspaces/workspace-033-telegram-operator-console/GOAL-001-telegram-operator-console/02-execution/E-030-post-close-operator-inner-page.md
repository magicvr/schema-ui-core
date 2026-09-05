---
doc_type: goal-execution
id: E-030-post-close-operator-inner-page
parent: GOAL-001-telegram-operator-console
date: 2026-09-05
source: self
status: done
version: 0.1.0
---

# E-030 · Telegram 人工会话入口与内页分离（2026-09-05）

## 已发生事实

- Root R4 关门后，按用户对 Telegram 通道页交互边界的修改要求，保留
  `telegram-settings` 作为侧栏页，并新增非侧栏的 `telegram-operator` 内页，路由为
  `/telegram-settings/operator`。
- 设置页 schema 只提供配置表单与“进入人工会话”导航入口；人工会话页 schema 以
  `surface: "operator"` 挂载原有 `telegram-admin-tab`。
- `telegram-admin-tab` 按 surface 隔离配置与 operator 内容：设置页不渲染会话、时间线、
  composer、重试或运行态 `captured_messages_count`，也不发起 operator sessions polling
  或 lease；这些能力只在 operator 内页出现。
- App 为 operator 内页补充设置页父级 breadcrumb 与返回路径；provider、kernel profile、
  manifest/schema 契约和中英文文案同步更新，并补充 manifest/schema/provider、App 导航和
  UI surface 测试。
- 代码变更已在 checkpoint `6a94ba28fad08de43d3b88a129c5dcbcd0b18ccb` 提交，提交信息为
  `feat(telegram): move operator chat to inner page`。

## 审计与验证事实

- A-005 independent 在修复前发现 required F-001（设置页仍显示运行态计数）及
  recommended F-002（契约 identity/resource/action 与负向 surface 断言不足）；原始意见保留。
- A-006 已按 `fixed` 路径响应：计数移入 operator section；设置页断言无计数、无 operator
  内容和无 operator 请求；provider 测试补齐 settings/operator 的 `ModuleID`、`Key`、`Owner`、
  `Resources`、`Actions` 与 datasource。
- A-007 `subagent (gpt-5.6-sol · reasoning medium)` 对 checkpoint 当前代码独立复审为
  `pass`，`open_required: 0`；A-008 完成最终 self response。未调用 Grok。
- Web 定向测试：3 个文件、38 个测试通过；Web 全量 `npm test`：92 个文件、1214 个测试通过。
- Web `npm run build` 通过，1870 个模块转换完成；仅保留既有 chunk size warning，构建生成的
  conformance projection 已恢复，不纳入本次提交。
- API `go test ./...`（工作目录 `apps/api`）通过；Telegram/manifest/schema、composition、
  kernel 定向测试亦通过。`git diff --check` 与新增文件尾随空白检查通过。

## 状态边界

本条是 Root 关门后的局部 UI 结构修正，不改变 Root 六项成功标准、progress 或 status，
因此不重新打开 GOAL-001/workspace-033；VP-033 仍保持 `active`。
