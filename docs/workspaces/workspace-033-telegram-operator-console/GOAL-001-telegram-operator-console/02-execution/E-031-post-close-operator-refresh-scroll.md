---
doc_type: goal-execution
id: E-031-post-close-operator-refresh-scroll
parent: GOAL-001-telegram-operator-console
date: 2026-09-05
source: self
status: done
version: 0.1.0
---

# E-031 · 人工会话刷新稳定性与滚动布局（2026-09-05）

## 已发生事实

- 按用户反馈修正 Telegram 人工会话内页的两个 UI 问题：后台轮询刷新期间不再用
  loading 分支替换已有消息；消息较多或会话较多时不再依赖页面无限增长。
- `loadTimeline` 在后台刷新开始时保留已有 timeline；已有消息继续渲染，仅在标题旁显示轻量
  `Refreshing…` 状态。切换到实际不同的 chat 时才清空旧 timeline，并保留既有 stale-response
  守卫与请求去重。
- operator 面板建立视口上限、纵向 flex 和 `min-h-0` 高度链；sessions nav 与 message list
  分别使用独立纵向滚动，message list 占用剩余空间，composer 使用 `shrink-0` 固定在底部。
  长 chat ID、时间戳、状态文案和 textarea 分别增加换行、截断或高度约束。
- 新增后台刷新 pending 时保留旧消息的回归测试，并补充滚动容器、min-height、flex 和
  composer 固定布局契约断言；英文/中文刷新文案同步落盘。
- 代码变更已在 checkpoint `dc7ac5e5be1c98042f9ef5da453892edcc12331b` 提交，提交信息为
  `fix(telegram): stabilize operator refresh and scrolling`。

## 审计与验证事实

- A-009 `subagent (gpt-5.6-sol · reasoning medium)` 对当前 checkpoint 独立复审为
  `pass`、`open_required: 0`，确认后台消息保留、会话切换竞态防护、sessions/message 独立
  滚动和 composer 固定约束均有源码与测试证据。
- A-009 仅保留一个非阻断验证边界：单元测试验证了布局 class 契约，但没有真实浏览器的像素级
  clientHeight/scrollHeight 测量；这不构成当前 required 门禁。
- Web 定向测试：2 个文件、28 个测试通过；Web 全量 `npm test`：92 个文件、1215 个测试通过。
- Web `npm run build` 通过，1870 个模块转换完成；仅有既有 chunk size warning，构建自动改写的
  conformance projection 已恢复，不纳入本次提交。
- API `go test ./...`（工作目录 `apps/api`）通过；`go test ./internal/docscheck` 亦通过。
  `git diff --check` 通过，Git 工作树在代码提交与治理记录提交后保持干净。

## 状态边界

本条是 Root 关门后的局部 UI 稳定性与可用性修正，不改变 Root 六项成功标准、progress 或
status，因此不重新打开 GOAL-001/workspace-033；VP-033 仍保持 `active`。
