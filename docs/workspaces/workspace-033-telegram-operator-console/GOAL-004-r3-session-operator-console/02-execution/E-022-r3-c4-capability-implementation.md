---
doc_type: goal-execution
id: E-022-r3-c4-capability-implementation
parent: GOAL-004-r3-session-operator-console
date: 2026-09-05
source: self
status: done
version: 0.1.0
---

# E-022 · R3 C4 capability 实现与非阻断覆盖（2026-09-05）

## 已发生事实

- A-037 为一次性 `subagent (gpt-5.6-sol · reasoning medium)` independent 合同审计，判定 `fail`，提出 F-037-1～F-037-4 四项 required；用户已明确改用该子代理 provider，本轮未调用 Grok。
- 已在 implementation checkpoints `cae40b3a` 与 `da9d955e` 落地独立 capability 路由：Bot API `getChatMember`、结构化 Telegram API 错误、成员状态 fail-closed 映射、channel.telegram capability service、独立 cache namespace、60 秒 absolute TTL、同键 single-flight、403 精确失效、handler/composition/provider/profile 接线、错误目录和 Admin UI 发送/retry 生命周期；同时修复 Web 表单日期边界的类型错误并保留日期约束。
- API 覆盖了 capability 权限/运行时/refresh 解析/成功载荷/错误门禁；Bot API、成员状态、缓存隔离/强刷/403 失效/真实 HTTP 403、composition route 清单均有测试。
- Web 覆盖了允许、拒绝/不可用、显式 `refresh=1`、发送/重试 request ID 和切换 chat 后旧响应丢弃；普通 10 秒 sessions/timeline 刷新未隐式触发 capability 探测。

## 验证事实

- `go test ./... -count=1`：通过。
- `go test -race ./internal/channel/telegram ./internal/handler -count=1`：通过。
- `npm test -- --run`：92 个测试文件、1213 个测试通过；Telegram 定向为 15/15。
- `git diff --check` 与两个新增 Go 文件的行尾空白检查：通过。
- 用户要求修复的 `apps/web/src/renderer/form-controls.tsx:946-947` 类型错误已由 `da9d955e` 修复：数字与 ISO 日期边界按控件类型收窄，解析层不再丢弃 datePicker 边界，并补充解析/DOM 回归测试。`npm run build`：通过（仅有既有 chunk size warning）；构建副作用改写的三个 conformance 生成文件已恢复，未进入 checkpoint。

## 审计边界

A-038 self implementation response 已将 A-037 F-037-1～F-037-4 按 `fixed` 路径形成代码、测试和接线证据；`da9d955e` 又清除了 Web 构建错误，当前构建与全量 Web 测试均通过。C4 尚未关闭，等待 GPT-5.6-sol（medium）对当前 checkpoint 的 independent implementation audit。未接受 residual，不作 overrule。
