---
doc_type: goal-audit
id: A-017-r3-c3-contract-self
parent: GOAL-004-r3-session-operator-console
date: 2026-09-04
source: self
auditor: Codex
audit_type: stage-gate
scope: R3 C3 API、权限、运行时门禁、v69 outbound 状态/幂等/显式重试合同；对照 VP-033、D-002、D-008 与现有 Telegram provider/runtime/store 接缝；不审 C3 生产实现
verdict: pass
open_required: 0
version: 0.1.0
---

# A-017 · R3 C3 API/权限/运行时合同自审（2026-09-04）

## 审视结论

`D-009-r3-c3-operator-console-contract` 已把 C3 的可实施边界写成可核对合同，
可以放行到生产代码实现；本条不把合同通过写成 C3 已完成，也不替代后续 Grok
independent 审计。

## 对照证据

- D-002 已冻结专用 `telegram.operator.read/write`、未绑定入口、同步 sender、
  `pending → sent/failed`、客户端 `request_id` 幂等和显式重试；D-008 已冻结每次
  重试新建 request/记录并以 `retry_of` 指向原始请求（D-002 L20–30；D-008 L15–21）。
- VP-033 只要求未绑定且连接成功时列出 bot 实际收到的会话文本并允许管理员代 bot
  发言；首波仍只文本、无历史/FSM/群发/多 bot（VP-033 L33–47、L58–59、L96–98）。
- 当前运行时提供 `ConnectionStatus` 的 `running`/receiver/`BotID` 状态，dispatcher
  提供 `HasBusinessHandlers()` 占用位；Provider 已有鉴权非 Public 路由接缝，后续
  C3 会在此处新增 operator routes（`runtime.go` L23–43、L271–289；`dispatcher.go`
  L27–38；`provider.go` L78–145）。
- D-009 明确了 v69 双方言 outbound 表、root 级 pending 唯一边界、统一入站/出站
  成绩单、服务端 bot scope、分页溢出保护、失败留痕和 token/secret/raw JSON 不落盘，
  并列出 SQLite/gated PG、权限、并发、失败顺序和 race 验证分母（D-009 L23–107）。

## 门禁判定

- **用户裁决已满足**：唯一影响重试身份的方案选择已通过裁决工具落盘为 D-008；
  D-009 没有把未选方案混入合同。
- **信息门禁**：I-033-021/022 的 C3 合同已就绪；其“实现+测试+independent”
  证据仍未宣称完成，须在代码实现后补齐。I-033-009/010 的 UI timer 与
  `getChatMember` 反馈仍属于 C4，不被 C3 合同越界关闭。
- **范围隔离**：C3 不新增 operator page/navigation/fragment，不扩张 kernel
  Telegram port，不启用历史回灌或后台自动重试。
- **required findings**：本条无 required finding，`open_required: 0`。下一门禁仍须
  由 Grok `grok-4.6 · reasoning high` 独立审视；independent 失败不得由本条降级。

本条不修改目标 status/progress，不接受 residual，不 overrule，C3 生产实现获准
开始但 C3 检查点仍为待开始。
